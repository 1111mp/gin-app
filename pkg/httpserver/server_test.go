package httpserver

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type httpReply struct {
	status int
	err    error
}

type activeRequest struct {
	server            *serverImpl
	release, finished chan struct{}
	reply             chan httpReply
}

func startActiveRequest(t *testing.T) *activeRequest {
	t.Helper()

	request := &activeRequest{release: make(chan struct{}), finished: make(chan struct{}), reply: make(chan httpReply, 1)}
	entered := make(chan struct{})
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		close(entered)

		select {
		case <-request.release:
			c.Status(http.StatusOK)
		case <-c.Request.Context().Done():
		}

		close(request.finished)
	})
	server, ok := New(router).(*serverImpl)
	require.True(t, ok)

	request.server = server
	listenerConfig := net.ListenConfig{}
	listener, err := listenerConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	served := make(chan error, 1)

	go func() { served <- request.server.srv.Serve(listener) }()

	t.Cleanup(func() {
		select {
		case <-request.release:
		default:
			close(request.release)
		}

		require.NoError(t, request.server.srv.Close())
		require.ErrorIs(t, <-served, http.ErrServerClosed)
	})

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+listener.Addr().String(), http.NoBody)
	require.NoError(t, err)

	go func() {
		client := &http.Client{Timeout: 3 * time.Second}

		response, err := client.Do(req)
		if err != nil {
			request.reply <- httpReply{err: err}

			return
		}

		request.reply <- httpReply{status: response.StatusCode, err: response.Body.Close()}
	}()

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("request did not start")
	}

	return request
}

func TestShutdownDrainsActiveRequest(t *testing.T) {
	t.Parallel()
	request := startActiveRequest(t)
	stopped := make(chan error, 1)

	go func() { stopped <- request.server.Shutdown(t.Context()) }()

	select {
	case err := <-stopped:
		t.Fatalf("shutdown returned with request active: %v", err)
	case <-time.After(10 * time.Millisecond):
	}

	close(request.release)
	require.NoError(t, <-stopped)

	reply := <-request.reply
	require.NoError(t, reply.err)
	require.Equal(t, http.StatusOK, reply.status)
}

func TestShutdownCancelsActiveRequest(t *testing.T) {
	t.Parallel()

	for _, parentDeadline := range []bool{false, true} {
		t.Run(map[bool]string{false: "local_deadline", true: "parent_deadline"}[parentDeadline], func(t *testing.T) {
			t.Parallel()
			request := startActiveRequest(t)
			request.server.shutdownTimeout = 50 * time.Millisecond
			ctx := t.Context()

			if parentDeadline {
				request.server.shutdownTimeout = time.Second

				var cancel context.CancelFunc

				ctx, cancel = context.WithTimeout(ctx, 50*time.Millisecond)
				defer cancel()
			}

			require.ErrorIs(t, request.server.Shutdown(ctx), context.DeadlineExceeded)

			select {
			case <-request.finished:
			case <-time.After(time.Second):
				t.Fatal("request context was not canceled")
			}

			<-request.reply
		})
	}
}
