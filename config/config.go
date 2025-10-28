package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// ConfigInter -.
	ConfigInterface interface {
		App() App
		HTTP() HTTP
		JWT() JWT
		Log() Log
		PG() PG
		// GRPC() GRPC
		Metrics() Metrics
		Swagger() Swagger
	}

	// Config -.
	Config struct {
		app  App
		http HTTP
		jwt  JWT
		log  Log
		pg   PG
		// grpc    GRPC
		metrics Metrics
		swagger Swagger
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port       string `env:"HTTP_PORT,required"`
		CookieName string `env:"HTTP_COOKIE_NAME,required"`
	}

	// JWT -.
	JWT struct {
		SECRET string `env:"JWT_SECRET,required"`
	}

	// Log -.
	Log struct {
		Dir   string `env:"LOG_DIR,required"`
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// // GRPC -.
	// GRPC struct {
	// 	Port string `env:"GRPC_PORT,required"`
	// }

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (ConfigInterface, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

// Implementation of ConfigInter methods
func (c *Config) App() App         { return c.app }
func (c *Config) HTTP() HTTP       { return c.http }
func (c *Config) JWT() JWT         { return c.jwt }
func (c *Config) Log() Log         { return c.log }
func (c *Config) PG() PG           { return c.pg }
func (c *Config) Metrics() Metrics { return c.metrics }
func (c *Config) Swagger() Swagger { return c.swagger }

// func (c *Config) GRPC() GRPC    { return c.grpc }
