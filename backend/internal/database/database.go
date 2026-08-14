package database

import (
	"log"
	"time"

	"github.com/ys-systemrun/require_manager/backend/internal/config"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the database connection, retrying while PostgreSQL warms up.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormCfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	for attempt := 1; attempt <= 15; attempt++ {
		db, err = gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
		if err == nil {
			break
		}
		log.Printf("database not ready (attempt %d/15): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate runs AutoMigrate for all registered models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(models.AllModels()...)
}
