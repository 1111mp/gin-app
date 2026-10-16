// Package server implements RabbitMQ RPC server.
package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	rmqrpc "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// CallHandler -.
type CallHandler func(context.Context, *amqp.Delivery) (any, error)

// CallHandlerTyped -.
func CallHandlerTyped[Req, Resp any](
	fn func(context.Context, Req) (Resp, error),
) CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req Req
		if err := json.Unmarshal(d.Body, &req); err != nil {
			return nil, fmt.Errorf("rmq_rpc - %s - json.Unmarshal: %w", d.Type, err)
		}

		resp, err := fn(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// Server -.
type Server struct {
	ctx             context.Context
	eg              *errgroup.Group
	cancel          context.CancelFunc
	reconnectCtx    context.Context
	reconnectCancel context.CancelFunc
	stopOnce        sync.Once

	conn   *rmqrpc.Connection
	router map[string]CallHandler
	stop   chan struct{}
	notify chan error

	timeout time.Duration

	logger logger.Logger
}

// New -.
func New(url, serverExchange string, l logger.Logger, opts ...Option) (*Server, error) {
	baseCtx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(baseCtx)
	reconnectCtx, reconnectCancel := context.WithCancel(context.Background())

	group.SetLimit(1) // Run only one goroutine

	cfg := rmqrpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		ctx:             ctx,
		cancel:          cancel,
		reconnectCtx:    reconnectCtx,
		reconnectCancel: reconnectCancel,
		eg:              group,
		conn:            rmqrpc.New(serverExchange, cfg),
		router:          make(map[string]CallHandler),
		stop:            make(chan struct{}),
		notify:          make(chan error, 1),
		timeout:         _defaultTimeout,
		logger:          l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnectContext(s.reconnectCtx)
	if err != nil {
		cancel()
		reconnectCancel()
		return nil, fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err)
	}

	return s, nil
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		defer close(s.notify)

		err := s.handleMessages()
		if err != nil {
			select {
			case <-s.stop:
				return nil
			default:
			}

			s.notify <- err

			return err
		}

		return nil
	})

	s.logger.Info("rmq_rpc server - Server - Started")
}

// RegisterRouter -.
func (s *Server) RegisterRouter(method string, handler CallHandler) {
	if method == "" {
		s.logger.Warn("rmq_rpc server - RegisterRouter - attempt to register empty method name")
		return
	}
	if handler == nil {
		s.logger.Warn(
			"rmq_rpc server - RegisterRouter - attempt to register nil handler",
			"method", method,
		)

		return
	}

	if _, exists := s.router[method]; exists {
		s.logger.Warn(
			"rmq_rpc server - RegisterRouter - overwriting existing handler",
			"method", method,
		)
	}

	s.router[method] = handler
	s.logger.Debug(
		"rmq_rpc server - RegisterRouter - registered RPC handler",
		"method", method,
	)
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown(ctx context.Context) error {
	s.stopOnce.Do(func() { close(s.stop); s.reconnectCancel() })

	done := make(chan error, 1)

	go func() { done <- s.eg.Wait() }()

	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	s.cancel()

	closeErr := s.conn.Close(ctx)
	if errors.Is(err, context.Canceled) && ctx.Err() == nil {
		err = nil
	}

	return errors.Join(err, closeErr)
}

func (s *Server) handleMessages() error {
	for {
		select {
		case <-s.stop:
			return nil
		default:
		}

		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-s.stop:
			return nil
		case d, opened := <-s.conn.Messages():
			if !opened {
				err := s.reconnect()
				if err != nil {
					return err
				}

				break
			}

			s.serveCall(&d)
		}
	}
}

func (s *Server) reconnect() error {
	return s.conn.AttemptConnectContext(s.reconnectCtx)
}

func (s *Server) serveCall(d *amqp.Delivery) {
	defer func() {
		if s.ctx.Err() != nil {
			if err := d.Nack(false, true); err != nil {
				s.logger.Error(err, "rmq_rpc server - Server - canceled call nack")
			}

			return
		}

		s.ack(d, false)
	}()

	callHandler, ok := s.router[d.Type]
	if !ok {
		s.publish(d, nil, rmqrpc.ErrBadHandler.Error())

		return
	}

	response, err := callHandler(s.ctx, d)
	if s.ctx.Err() != nil {
		return
	}

	if err != nil {
		s.publish(d, nil, rmqrpc.ErrInternalServer.Error())

		s.logger.Error(err, "rmq_rpc server - Server - serveCall - callHandler")

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - serveCall - json.Marshal")
	}

	s.publish(d, body, rmqrpc.Success)
}

func (s *Server) ack(d *amqp.Delivery, multiple bool) {
	err := d.Ack(multiple)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - ack - d.Ack")
	}
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.Publisher().Publish(
		d.ReplyTo,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		},
	)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - publish - s.conn.Channel.Publish")
	}
}
