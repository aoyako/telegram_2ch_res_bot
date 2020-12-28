package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB constructor of gorm.DB databse with postgresql database
func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	connStr := formatPostgresConfig(cfg)
	return gorm.Open(postgres.Open(connStr), &gorm.Config{})
}

// Formats config struct to meet gorm's expectations
func formatPostgresConfig(cfg Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
}
