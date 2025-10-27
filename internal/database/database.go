package database

import (
	"ac-ai/internal/config"
	"ac-ai/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established.")

	// Автоматическая миграция
	err = db.AutoMigrate(
		&models.User{},
		&models.FinancialProfile{},
		&models.ScoringApplication{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database migrated.")
	return db, nil
}