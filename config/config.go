package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		// PG      PG
		// GRPC    GRPC
		Metrics Metrics
		Swagger Swagger
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	// Log -.
	Log struct {
		Dir   string `env:"LOG_DIR,required"`
		Level string `env:"LOG_LEVEL,required"`
	}

	// // PG -.
	// PG struct {
	// 	PoolMax int    `env:"PG_POOL_MAX,required"`
	// 	URL     string `env:"PG_URL,required"`
	// }

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
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
