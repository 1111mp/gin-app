package httpserver

import (
	"net"
	"time"
)

// option -.
type Option func(*serverImpl)

// Port -.
func Port(port string) Option {
	return func(s *serverImpl) {
		s.Address = net.JoinHostPort("0.0.0.0", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *serverImpl) {
		s.readTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *serverImpl) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *serverImpl) {
		s.shutdownTimeout = timeout
	}
}
