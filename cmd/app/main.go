package main

import (
	"log"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/app"
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
