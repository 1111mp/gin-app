package grpcclient

import (
	"fmt"

	"github.com/1111mp/gin-app/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_defaultAddr = ":80"
)

type Client struct {
	conn    *pbgrpc.ClientConn
	address string

	logger logger.Logger
}

func New(l logger.Logger, opts ...Option) (*Client, error) {
	c := &Client{
		address: _defaultAddr,
		logger:  l,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	conn, err := pbgrpc.NewClient(c.address, pbgrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatal("gRPC Client - init error - grpc.NewClient", err)
		return nil, err
	}

	c.conn = conn
	l.Info("gRPC Client - initialized", "address", c.address)

	return c, nil
}

func (c *Client) Shutdown() error {
	c.logger.Info("gRPC Client - starting graceful shutdown")

	if err := c.conn.Close(); err != nil {
		c.logger.Error("gRPC Client - shutdown error", err)
		return fmt.Errorf("gRPC client close error: %w", err)
	}

	c.logger.Info("gRPC Client - shutdown completed")
	return nil
}

func (c *Client) GetConn() *pbgrpc.ClientConn {
	return c.conn
}
