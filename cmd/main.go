package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config file: %s", err.Error())
	}

	db, err := storage.NewPostgresDB(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	storage.MigrateDatabase(db)

	// storage := storage.NewStorage(db)
	// storage.User.Register(&logic.User{})
	// usr, _ := storage.User.GetByChatID(0)
	// storage.Subscription.Add(usr, &logic.Publication{UserID: 1})

	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&logic.User{}, &logic.Publication{}, &logic.Info{})

	fmt.Println("hello world")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
