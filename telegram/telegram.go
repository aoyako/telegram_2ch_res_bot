package telegram

import (
	"log"
	"strings"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/downloader"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/xfrr/goffmpeg/transcoder"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// TgBot represents telegram bot view
type TgBot struct {
	Bot        *telebot.Bot
	Controller *controller.Controller
	Downloader *downloader.Downloader
}

// NewTelegramBot constructor of TelegramBot
func NewTelegramBot(token string, cnt *controller.Controller, d *downloader.Downloader) *TgBot {
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
		Downloader: d,
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

func (tb *TgBot) Send(users []*logic.User, path, caption string) {
	if len(users) == 0 {
		return
	}
	// img := &telebot.Video{File: telebot.FromDisk("vid.mp4")}
	// fmt.Println(path)
	var file telebot.Sendable
	if strings.HasSuffix(path, ".mp4") {
		file = &telebot.Video{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".gif") {
		file = &telebot.Photo{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".webm") {
		tb.Downloader.Save(path)
		defer func() {
			err := tb.Downloader.Free(path)
			if err != nil {
				log.Println(err)
			}
			err = tb.Downloader.Free(strings.TrimSuffix(path, ".webm") + ".mp4")
			if err != nil {
				log.Println(err)
			}
		}()

		trans := new(transcoder.Transcoder)

		vidPath := tb.Downloader.Get(path)
		newVidPath := strings.TrimSuffix(vidPath, ".webm") + ".mp4"

		err := trans.Initialize(vidPath, newVidPath)
		if err != nil {
			return
		}
		done := trans.Run(false)
		err = <-done
		if err != nil {
			log.Println(err)
			return
		}

		file = &telebot.Video{File: telebot.FromDisk(newVidPath), Caption: caption}
	}

	for _, user := range users {
		tb.Bot.Send(&telebot.User{
			ID: int(user.ChatID),
		}, file)
	}
}
