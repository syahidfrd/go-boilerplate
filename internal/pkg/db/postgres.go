package db

import (
	"fmt"

	"github.com/syahidfrd/go-boilerplate/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgres creates a new PostgreSQL database connection using GORM
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DataSourceName()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConn)

	return db, nil
}

// AutoMigrate runs database migrations for the provided models
func AutoMigrate(db *gorm.DB, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	return nil
}
