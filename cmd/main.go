package main

import (
	"log"
	"os"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/dvach"
	"github.com/aoyako/telegram_2ch_res_bot/telegram"

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

	// db.AutoMigrate(&logic.User{}, &logic.Publication{}, &logic.Info{})

	// storage.MigrateDatabase(db)

	Storage := storage.NewStorage(db)
	storage.MigrateDatabase(db)
	controller := controller.NewController(Storage)
	// storage.User.Register(&logic.User{})
	// usr, _ := storage.User.GetByChatID(0)
	// storage.Subscription.Add(usr, &logic.Publication{UserID: 1})

	if err != nil {
		log.Fatal(err.Error())
	}

	// fmt.Println("hello world")
	bot := telegram.NewTelegramBot(os.Getenv("BOT_TOKEN"), controller)
	telegram.SetupHandlers(bot)
	// bot.Bot.Start()

	apicnt := dvach.NewAPIController(controller)
	apicnt.SetSenderFunc(bot.GetMessageFunc())
	apicnt.InitiateSending()

	bot.Bot.Start()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
