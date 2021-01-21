package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
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

// MigrateDatabase migrates database
func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&logic.User{}, &logic.Admin{}, &logic.Publication{}, &logic.Info{})

	if err != nil {
		log.Fatalf("Error migrating database")
	}

	var count int64
	db.Find(&logic.Info{}).Count(&count)
	if count == 0 {
		db.Create(&logic.Info{
			LastPost: uint64(time.Now().Unix()),
		})
	} else {
		var info logic.Info
		db.Find(&logic.Info{}).First(&info)
		info.LastPost = uint64(time.Now().Unix())
		db.Save(info)
	}
}
