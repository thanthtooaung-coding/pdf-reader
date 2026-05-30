package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

func Connect(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetimeMin) * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if cfg.DBAutoMigrate {
		log.Info("Running GORM auto-migrate")
		if err := autoMigrate(db); err != nil {
			return nil, fmt.Errorf("auto-migrate: %w", err)
		}
	}

	log.Info("Database connection established")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Workspace{},
		&models.File{},
		&models.Comment{},
		&models.AIJob{},
	)
}
