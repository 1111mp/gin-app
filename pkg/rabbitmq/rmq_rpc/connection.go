package rmqrpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionSetupTimeout = 3 * time.Second

// Config -.
type Config struct {
	URL      string
	WaitTime time.Duration
	Attempts int
}

// Connection -.
type Connection struct {
	mu               sync.RWMutex
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Delivery   <-chan amqp.Delivery
}

// New -.
func New(consumerExchange string, cfg Config) *Connection {
	conn := &Connection{
		ConsumerExchange: consumerExchange,
		Config:           cfg,
	}

	return conn
}

// AttemptConnect -.
func (c *Connection) AttemptConnect() error { return c.AttemptConnectContext(context.Background()) }

func (c *Connection) AttemptConnectContext(ctx context.Context) error {
	var err error

	for i := c.Attempts; i > 0; i-- {
		attemptCtx, cancel := context.WithTimeout(ctx, connectionSetupTimeout)
		err = c.connect(attemptCtx)

		cancel()

		if err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)

		if i == 1 {
			break
		}

		timer := time.NewTimer(c.WaitTime)
		select {
		case <-ctx.Done():
			timer.Stop()

			return ctx.Err()
		case <-timer.C:
		}
	}

	if err != nil {
		return fmt.Errorf("rmq_rpc - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Connection) connect(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var (
		raw        net.Conn
		stopCancel func() bool
	)

	conn, err := amqp.DialConfig(c.URL, amqp.Config{Dial: func(network, addr string) (net.Conn, error) {
		d := net.Dialer{}

		var err error

		raw, err = d.DialContext(ctx, network, addr)
		if err == nil {
			stopCancel = context.AfterFunc(ctx, func() {
				_ = raw.Close()
			})
		}

		return raw, err
	}})

	if stopCancel != nil {
		defer stopCancel()
	}

	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	success := false

	defer func() {
		if !success {
			//nolint:errcheck // Best-effort forced cleanup; return the original connection setup error.
			_ = conn.CloseDeadline(time.Now())
		}
	}()

	channel, delivery, err := c.consume(conn)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	old := c.Connection
	c.Connection, c.Channel, c.Delivery = conn, channel, delivery
	c.mu.Unlock()

	if old != nil {
		//nolint:errcheck // Force-close the superseded connection without delaying reconnect.
		_ = old.CloseDeadline(time.Now())
	}

	success = true

	return nil
}

func (c *Connection) Messages() <-chan amqp.Delivery {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Delivery
}

func (c *Connection) Publisher() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Channel
}

func (c *Connection) Close(ctx context.Context) error {
	c.mu.RLock()
	conn := c.Connection
	c.mu.RUnlock()

	if conn == nil || conn.IsClosed() {
		return nil
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(time.Second)
	}

	err := conn.CloseDeadline(deadline)
	if errors.Is(err, amqp.ErrClosed) {
		return nil
	}

	return err
}

func (c *Connection) consume(conn *amqp.Connection) (channel *amqp.Channel, delivery <-chan amqp.Delivery, err error) {
	channel, err = conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("c.Connection.Channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("c.Connection.Channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("channel.QueueDeclare: %w", err)
	}

	err = channel.QueueBind(
		queue.Name,
		"",
		c.ConsumerExchange,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("channel.QueueBind: %w", err)
	}

	delivery, err = channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("channel.Consume: %w", err)
	}

	return channel, delivery, nil
}
