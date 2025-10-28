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
)

// Run creates objects via constructors.
func Run(cfg config.ConfigInterface) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	l := logger.New(cfg.Log().Dir, cfg.Log().Level)

	// Repository
	pg, err := postgres.New(cfg.PG().URL, postgres.MaxPoolSize(cfg.PG().PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP().Port))
	router.NewRouter(httpServer.App, cfg, pg, l)

	// Start server
	httpServer.Start()

	// Wait for interrupt signal to gracefully shutdown the server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Infof("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
