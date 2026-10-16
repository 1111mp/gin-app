package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	_defaultAddr            = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 8 * time.Second
)

// Server -.
type Server interface {
	Start() error
	Shutdown(context.Context) error
	GetAddress() string
}

// serverImpl -.
type serverImpl struct {
	srv    *http.Server
	cancel context.CancelFunc

	Address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

// New -.
func New(handler *gin.Engine, opts ...Option) Server {
	s := &serverImpl{
		srv:             nil,
		Address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.srv = &http.Server{
		BaseContext:    func(net.Listener) context.Context { return ctx },
		Addr:           s.Address,
		Handler:        handler,
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

// Start -.
func (s *serverImpl) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown -.
func (s *serverImpl) Shutdown(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, s.shutdownTimeout)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	s.cancel()

	if err != nil {
		return errors.Join(err, s.srv.Close())
	}

	return nil
}

// GetAddress -.
func (s *serverImpl) GetAddress() string {
	return s.Address
}
