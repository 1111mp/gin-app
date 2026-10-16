package server

import (
	"context"
	"testing"
	"time"

	rmqrpc "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestTypedHandlerReceivesCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	handler := CallHandlerTyped(func(ctx context.Context, req struct{}) (struct{}, error) { return req, ctx.Err() })
	_, err := handler(ctx, &amqp.Delivery{Body: []byte("{}")})
	require.ErrorIs(t, err, context.Canceled)
}

func startActiveHandler(t *testing.T, ignoreCancel bool) (server *Server, released chan struct{}, completed <-chan struct{}) {
	t.Helper()
	base, cancel := context.WithCancel(t.Context())
	reconnectCtx, reconnectCancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	t.Cleanup(reconnectCancel)

	group, ctx := errgroup.WithContext(base)
	server = &Server{ctx: ctx, cancel: cancel, eg: group, reconnectCtx: reconnectCtx, reconnectCancel: reconnectCancel, stop: make(chan struct{}), conn: &rmqrpc.Connection{}}
	release, entered, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})

	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	group.Go(func() error {
		defer close(finished)

		close(entered)

		if ignoreCancel {
			<-release
		} else {
			select {
			case <-ctx.Done():
			case <-release:
			}
		}

		return nil
	})
	<-entered

	return server, release, finished
}

func TestShutdownDrainsActiveHandler(t *testing.T) {
	t.Parallel()
	server, release, _ := startActiveHandler(t, false)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	stopped := make(chan error, 1)

	go func() { stopped <- server.Shutdown(ctx) }()

	select {
	case err := <-stopped:
		t.Fatalf("returned before handler finished: %v", err)
	case <-time.After(10 * time.Millisecond):
	}

	close(release)
	require.NoError(t, <-stopped)
}

func TestShutdownCancelsActiveHandler(t *testing.T) {
	t.Parallel()
	server, _, finished := startActiveHandler(t, false)

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()

	require.ErrorIs(t, server.Shutdown(ctx), context.DeadlineExceeded)

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("handler was not canceled")
	}
}

func TestShutdownBoundsUncooperativeHandler(t *testing.T) {
	t.Parallel()
	server, _, _ := startActiveHandler(t, true)

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()

	stopped := make(chan error, 1)

	go func() { stopped <- server.Shutdown(ctx) }()

	select {
	case err := <-stopped:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Fatal("shutdown hung")
	}
}
