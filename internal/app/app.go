// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/router"
	"github.com/1111mp/gin-app/pkg/httpserver"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/1111mp/gin-app/pkg/redis"
	"github.com/1111mp/gin-app/pkg/state"
)

// Run creates objects via constructors.
func Run(cfg config.ConfigInterface) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	logger := logger.New(cfg.Log().Dir, cfg.Log().Level)

	// postgres
	pg, err := postgres.New(cfg.PG().URL, postgres.MaxPoolSize(cfg.PG().PoolMax))
	if err != nil {
		logger.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// redis
	rdb, err := redis.New(cfg.Redis().URL, redis.MaxPoolSize(cfg.Redis().PoolMax))
	if err != nil {
		logger.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer rdb.Close()

	appState := &state.AppState{
		PG:    pg,
		Redis: rdb,
	}

	// HTTP Server
	httpServer := httpserver.New(logger, httpserver.Port(cfg.HTTP().Port))
	router.NewRouter(httpServer.App, cfg, appState, logger)

	// Start server
	httpServer.Start()

	// Wait for interrupt signal to gracefully shutdown the server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
