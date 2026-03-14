package dao

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// InitDB opens a SQLite database at the given path, runs auto-migrations,
// and returns the *gorm.DB instance.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(&models.User{}, &Task{}, &models.Post{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
