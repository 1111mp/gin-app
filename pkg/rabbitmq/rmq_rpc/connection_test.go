package rmqrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReconnectWaitIsCancelable(t *testing.T) {
	t.Parallel()

	connection := New("test", Config{URL: "invalid://", WaitTime: time.Minute, Attempts: 10})

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()

	result := make(chan error, 1)

	go func() { result <- connection.AttemptConnectContext(ctx) }()

	select {
	case err := <-result:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Fatal("retry ignored cancellation")
	}
}

func TestHandshakeIsCancelable(t *testing.T) {
	t.Parallel()

	listenerConfig := net.ListenConfig{}
	listener, err := listenerConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, listener.Close()) })

	accepted := make(chan net.Conn, 1)

	go func() {
		conn, err := listener.Accept()
		if err == nil {
			accepted <- conn
		}
	}()

	connection := New("test", Config{URL: "amqp://guest:guest@" + listener.Addr().String() + "/", Attempts: 1})

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	result := make(chan error, 1)

	go func() { result <- connection.AttemptConnectContext(ctx) }()

	select {
	case conn := <-accepted:
		t.Cleanup(func() { require.NoError(t, conn.Close()) })
	case <-time.After(time.Second):
		t.Fatal("no connection")
	}

	cancel()

	select {
	case err := <-result:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("handshake ignored cancellation")
	}
}
