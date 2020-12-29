package telegram

import (
	"log"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	telebot "gopkg.in/tucnak/telebot.v2"
)

// TgBot represents telegram bot view
type TgBot struct {
	Bot        *telebot.Bot
	Controller *controller.Controller
}

// NewTelegramBot constructor of TelegramBot
func NewTelegramBot(token string, cnt *controller.Controller) *TgBot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &TgBot{
		Bot:        bot,
		Controller: cnt,
	}
}

// SetupHandlers to default values
func SetupHandlers(tb *TgBot) {
	tb.Bot.Handle("/start", start(tb))
	tb.Bot.Handle("/list", list(tb))
	tb.Bot.Handle("/help", help(tb))
	tb.Bot.Handle("/add", add(tb))
	tb.Bot.Handle("/rm", del(tb))
}
