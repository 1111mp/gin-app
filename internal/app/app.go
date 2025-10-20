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
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Dir, cfg.Log.Level)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	router.NewRouter(httpServer.App, l)

	// Start server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Infof("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
