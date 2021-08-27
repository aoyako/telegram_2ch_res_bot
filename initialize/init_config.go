package initialize

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/downloader"
	"github.com/aoyako/telegram_2ch_res_bot/dvach"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
	"github.com/aoyako/telegram_2ch_res_bot/telegram"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// App initializes application
func App() (*telegram.TgBot, *dvach.APIController, uint64) {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config file: %s", err.Error())
	}

	var db *gorm.DB
	var err error
	for {
		db, err = storage.NewPostgresDB(storage.Config{
			Host:     viper.GetString("db.host"),
			Port:     os.Getenv("DB_PORT"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("DB_PASSWORD"),
		})

		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}

		break
	}

	if err != nil {
		log.Fatalf("Error creating database: %s", err.Error())
	}

	admins := storage.InitDatabase{
		Admin: stringToInt64Slice(viper.GetStringSlice("tg.admin_id")),
	}

	requestURL := &dvach.RequestURL{
		AllThreadsURL: viper.GetString("dapi.all"),
		ThreadURL:     viper.GetString("dapi.thread"),
		ResourceURL:   viper.GetString("dapi.resource"),
	}

	Storage := storage.NewStorage(db, &admins)
	controller := controller.NewController(Storage)

	bot := telegram.NewTelegramBot(os.Getenv("BOT_TOKEN"), controller, downloader.NewDownloader(
		viper.GetString("disk.path"),
		viper.GetUint64("disk.size")))
	requester := dvach.NewRequester(requestURL)
	apicnt := dvach.NewAPIController(controller, bot, requester)

	telegram.SetupHandlers(bot)
	storage.MigrateDatabase(db)

	return bot, apicnt, viper.GetUint64("polling.time")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func stringToInt64Slice(data []string) []int64 {
	result := make([]int64, len(data))
	for key := range data {
		tmp, err := strconv.ParseInt(data[key], 10, 64)
		if err != nil {
			panic(err)
		}
		result[key] = tmp
	}
	return result
}
