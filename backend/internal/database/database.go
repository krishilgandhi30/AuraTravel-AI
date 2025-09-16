package database

import (
	"fmt"
	"log"
	"time"

	"auratravel-backend/internal/config"
	"auratravel-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize sets up the database connection and runs migrations
func Initialize() error {
	cfg := config.GetConfig()

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DatabaseHost,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabasePort,
		cfg.DatabaseSSLMode,
	)

	// Set up GORM logger
	var gormLogger logger.Interface
	if cfg.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database connection established and migrations completed")
	return nil
}

// runMigrations automatically migrates the database schema
func runMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.TravelPreferences{},
		&models.EmergencyContact{},
		&models.Trip{},
		&models.TripPreferences{},
		&models.TripCollaborator{},
		&models.Itinerary{},
		&models.DayPlan{},
		&models.Activity{},
		&models.Meal{},
		&models.Accommodation{},
		&models.Transportation{},
		&models.Recommendation{},
		&models.SearchHistory{},
		&models.AnalyticsEvent{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
