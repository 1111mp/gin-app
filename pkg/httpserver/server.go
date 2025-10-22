package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 8 * time.Second
)

// Server -.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App    *gin.Engine
	srv    *http.Server
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &Server{
		ctx:             ctx,
		eg:              group,
		App:             nil,
		srv:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	app := gin.New()

	s.App = app

	s.srv = &http.Server{
		Addr:           s.address,
		Handler:        s.App.Handler(),
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.notify <- err
			close(s.notify)
			return err
		}
		return nil
	})

	s.logger.Info("http server - Listening on ", s.address)
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(err, "httpserver - Server - Shutdown - s.srv.Shutdown")
		shutdownErrors = append(shutdownErrors, err)
	}

	// Wait for all goroutines to finish and get any error
	if err := s.eg.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(err, "httpserver - Server - Shutdown - s.eg.Wait")
		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("http server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
