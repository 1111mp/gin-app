// Package grpcserver implements HTTP server.
package grpcserver

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"golang.org/x/sync/errgroup"
	pbgrpc "google.golang.org/grpc"
)

const (
	_defaultAddr   = ":80"
	forceStopGrace = 500 * time.Millisecond
)

// Server -.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App     *pbgrpc.Server
	notify  chan error
	address string

	logger logger.Logger
}

// New -.
func New(l logger.Logger, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &Server{
		ctx:     ctx,
		eg:      group,
		App:     pbgrpc.NewServer(),
		notify:  make(chan error, 1),
		address: _defaultAddr,
		logger:  l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		defer close(s.notify)

		var lc net.ListenConfig

		ln, err := lc.Listen(s.ctx, "tcp", s.address)
		if err != nil && !errors.Is(err, pbgrpc.ErrServerStopped) {
			s.notify <- err

			return err
		}

		err = s.App.Serve(ln)
		if err != nil && !errors.Is(err, pbgrpc.ErrServerStopped) {
			s.notify <- err

			return err
		}

		return nil
	})

	s.logger.Info("grpc server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan error, 1)

	go func() {
		s.App.GracefulStop()

		done <- s.eg.Wait()
	}()

	select {
	case err := <-done:
		if errors.Is(err, context.Canceled) {
			return nil
		}

		return err
	case <-ctx.Done():
		// Stop can also wait behind GracefulStop when a handler ignores cancellation.
		go s.App.Stop()

		timer := time.NewTimer(forceStopGrace)
		defer timer.Stop()

		select {
		case <-done:
		case <-timer.C:
		}

		return ctx.Err()
	}
}
