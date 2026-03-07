package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	goredis "github.com/redis/go-redis/v9"
)

const (
	__defaultMaxPoolSize = 10
)

// Redis -.
type Redis struct {
	maxPoolSize int

	Client *goredis.Client
}

// New -.
func New(url string, opts ...Option) (*Redis, error) {
	rdb := &Redis{
		maxPoolSize: __defaultMaxPoolSize,
	}

	// Custom options
	for _, opt := range opts {
		opt(rdb)
	}

	config, err := goredis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - goredis.ParseConfig: %w", err)
	}

	config.PoolSize = rdb.maxPoolSize

	rdb.Client = goredis.NewClient(config)

	// Ping to verify connection
	if err := rdb.Client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - rdb.Client.Ping: %w", err)
	}

	// OpenTelemetry instrumentation
	if err := errors.Join(redisotel.InstrumentTracing(rdb.Client), redisotel.InstrumentMetrics(rdb.Client)); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - redisotel.Instrumentation: %w", err)
	}

	return rdb, nil
}

// Close -.
func (rdb *Redis) Close() {
	if rdb.Client != nil {
		rdb.Client.Close()
	}
}
