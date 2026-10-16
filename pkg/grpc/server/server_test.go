package grpcserver

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

type blockingHealth struct {
	grpc_health_v1.UnimplementedHealthServer
	entered, release chan struct{}
	ignoreCancel     bool
}

func (h *blockingHealth) Check(ctx context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	close(h.entered)

	if h.ignoreCancel {
		<-h.release
	} else {
		select {
		case <-h.release:
		case <-ctx.Done():
		}
	}

	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, ctx.Err()
}

func startActiveRPC(t *testing.T, ignoreCancel bool) (server *Server, release chan struct{}, outcome <-chan error) {
	t.Helper()

	server = New(logger.New("", "error"))
	listener := bufconn.Listen(1024 * 1024)
	handler := &blockingHealth{entered: make(chan struct{}), release: make(chan struct{}), ignoreCancel: ignoreCancel}
	grpc_health_v1.RegisterHealthServer(server.App, handler)
	server.eg.Go(func() error { return server.App.Serve(listener) })
	t.Cleanup(func() {
		select {
		case <-handler.release:
		default:
			close(handler.release)
		}

		server.App.Stop()
	})

	conn, err := grpc.NewClient("passthrough:///test", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)

	result := make(chan error, 1)

	go func() {
		_, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		result <- err
	}()

	select {
	case <-handler.entered:
	case <-ctx.Done():
		t.Fatal("RPC did not start")
	}

	return server, handler.release, result
}

func TestShutdownDrainsActiveRPC(t *testing.T) {
	t.Parallel()
	server, release, result := startActiveRPC(t, false)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	stopped := make(chan error, 1)

	go func() { stopped <- server.Shutdown(ctx) }()

	select {
	case err := <-stopped:
		t.Fatalf("returned before RPC finished: %v", err)
	case <-time.After(10 * time.Millisecond):
	}

	close(release)
	require.NoError(t, <-stopped)
	require.NoError(t, <-result)
}

func TestShutdownCancelsActiveRPC(t *testing.T) {
	t.Parallel()

	for _, ignoreCancel := range []bool{false, true} {
		t.Run(map[bool]string{false: "honors_cancel", true: "ignores_cancel"}[ignoreCancel], func(t *testing.T) {
			t.Parallel()
			server, _, result := startActiveRPC(t, ignoreCancel)

			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()

			stopped := make(chan error, 1)

			go func() { stopped <- server.Shutdown(ctx) }()

			select {
			case err := <-stopped:
				require.ErrorIs(t, err, context.DeadlineExceeded)
			case <-time.After(time.Second):
				t.Fatal("shutdown hung")
			}

			select {
			case <-result:
			case <-time.After(time.Second):
				t.Fatal("RPC connection was not closed")
			}
		})
	}
}
