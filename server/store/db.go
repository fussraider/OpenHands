package store

import (
	"log/slog"
	"openhands-go/server/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}

	// Auto Migrate models
	err = db.AutoMigrate(&models.ConversationInfo{})
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
		return err
	}

	DB = db
	return nil
}
