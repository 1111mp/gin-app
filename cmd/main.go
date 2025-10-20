package main

import (
	"log"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
