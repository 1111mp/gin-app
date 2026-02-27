package grpcclient

import (
	"net"
)

// Option -.
type Option func(*Client)

// Port -.
func Port(port string) Option {
	return func(s *Client) {
		s.address = net.JoinHostPort("", port)
	}
}
