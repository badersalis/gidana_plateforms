package database

import (
	"log"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	cfg := config.App

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.DatabaseURL != "" {
		DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	} else {
		DB, err = gorm.Open(sqlite.Open(cfg.DBPath), gormCfg)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	migrate()
}

func migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.PropertyImage{},
		&models.Rental{},
		&models.Review{},
		&models.Favorite{},
		&models.Alert{},
		&models.Wallet{},
		&models.Transaction{},
		&models.SearchHistory{},
		&models.Conversation{},
		&models.Message{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration completed")
}
