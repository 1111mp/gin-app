package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config interface {
		App() App
		HTTP() HTTP
		JWT() JWT
		Log() Log
		PG() PG
		Redis() Redis
		// GRPC() GRPC
		RMQ() RMQ
		Metrics() Metrics
		Swagger() Swagger
		Github() Github
		Google() Google
	}

	// configImpl -.
	configImpl struct {
		AppData   App
		HTTPData  HTTP
		JWTData   JWT
		LogData   Log
		PGData    PG
		RedisData Redis
		// GRPCData    GRPC
		RMQData     RMQ
		MetricsData Metrics
		SwaggerData Swagger
		GoogleData  Google
		GithubData  Github
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

	Redis struct {
		PoolMax int    `env:"REDIS_POOL_MAX"`
		URL     string `env:"REDIS_URL,required"`
	}

	// // GRPC -.
	// GRPC struct {
	// 	Port string `env:"GRPC_PORT,required"`
	// }

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	Google struct {
		ClientID     string `env:"GOOGLE_CLIENT_ID,required"`
		ClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
		RedirectURL  string `env:"GOOGLE_REDIRECT_URL,required"`
	}

	Github struct {
		ClientID     string `env:"GITHUB_CLIENT_ID,required"`
		ClientSecret string `env:"GITHUB_CLIENT_SECRET,required"`
		RedirectURL  string `env:"GITHUB_REDIRECT_URL,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg := &configImpl{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

// Implementation of ConfigInter methods
func (c *configImpl) App() App         { return c.AppData }
func (c *configImpl) HTTP() HTTP       { return c.HTTPData }
func (c *configImpl) JWT() JWT         { return c.JWTData }
func (c *configImpl) Log() Log         { return c.LogData }
func (c *configImpl) PG() PG           { return c.PGData }
func (c *configImpl) Redis() Redis     { return c.RedisData }
func (c *configImpl) Metrics() Metrics { return c.MetricsData }
func (c *configImpl) Swagger() Swagger { return c.SwaggerData }
func (c *configImpl) Google() Google   { return c.GoogleData }
func (c *configImpl) Github() Github   { return c.GithubData }

// func (c *configImpl) GRPC() GRPC    { return c.GRPCData }
func (c *configImpl) RMQ() RMQ { return c.RMQData }
