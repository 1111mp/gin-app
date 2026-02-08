package httpserver

import (
	"context"
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
	Shutdown() error
	GetAddress() string
}

// serverImpl -.
type serverImpl struct {
	srv *http.Server

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

	s.srv = &http.Server{
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
func (s *serverImpl) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

// GetAddress -.
func (s *serverImpl) GetAddress() string {
	return s.Address
}
