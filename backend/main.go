package main

import (
	"log"

	"github.com/ys-systemrun/require_manager/backend/internal/config"
	"github.com/ys-systemrun/require_manager/backend/internal/database"
	"github.com/ys-systemrun/require_manager/backend/internal/router"
	"github.com/ys-systemrun/require_manager/backend/internal/seed"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := seed.Run(db); err != nil {
		log.Printf("warning: seeding failed: %v", err)
	}

	r := router.New(db, cfg)

	addr := ":" + cfg.ServerPort
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
