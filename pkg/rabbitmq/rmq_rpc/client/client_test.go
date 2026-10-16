package client

import (
	"context"
	"testing"
	"time"

	rmqrpc "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestShutdownWithoutNotifyReceiver(t *testing.T) {
	t.Parallel()

	base, cancel := context.WithCancel(t.Context())
	defer cancel()

	reconnectCtx, reconnectCancel := context.WithCancel(t.Context())
	defer reconnectCancel()

	group, ctx := errgroup.WithContext(base)
	deliveries := make(chan amqp.Delivery)
	close(deliveries)
	client := &Client{
		ctx: ctx, cancel: cancel, eg: group, reconnectCtx: reconnectCtx, reconnectCancel: reconnectCancel,
		conn: &rmqrpc.Connection{Config: rmqrpc.Config{URL: "invalid://", Attempts: 1}, Delivery: deliveries},
		stop: make(chan struct{}), notify: make(chan error, 1),
	}
	client.start()
	// errgroup cancels ctx after the consumer exits; nobody receives notify.
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("consumer blocked reporting failure")
	}

	require.Len(t, client.notify, 1)

	stopCtx, stopCancel := context.WithTimeout(t.Context(), time.Second)
	defer stopCancel()

	require.Error(t, client.Shutdown(stopCtx))
	require.NoError(t, stopCtx.Err())
}
