package telegram

import (
	"log"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/downloader"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/xfrr/goffmpeg/transcoder"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// MessageSender defines interface for bot-sender
type MessageSender interface {
	Send(r telebot.Recipient, value interface{}, args ...interface{}) (*telebot.Message, error)
	Handle(interface{}, interface{})
	Start()
}

// TgBot represents telegram bot view
type TgBot struct {
	Bot        MessageSender
	Controller *controller.Controller
	Downloader *downloader.Downloader
}

// NewTelegramBot constructor of TelegramBot
func NewTelegramBot(token string, cnt *controller.Controller, d *downloader.Downloader) *TgBot {
	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// If token is empty, do not send request
	// Developers of telebot lib made "offline" mode unaccessible
	// so reflect and unsafe is used to change that field
	if token == "" {
		rs := reflect.ValueOf(settings)
		rs2 := reflect.New(rs.Type()).Elem()
		rs2.Set(rs)
		rsf := rs2.FieldByName("offline")
		rsf = reflect.NewAt(rsf.Type(), unsafe.Pointer(rsf.UnsafeAddr())).Elem()
		rsf.SetBool(true)

		settings = rs2.Interface().(telebot.Settings)
	}

	bot, err := telebot.NewBot(settings)

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
	tb.Bot.Handle("/clist", cleverList(tb))
	tb.Bot.Handle("/help", help(tb))

	tb.Bot.Handle("/subs", subs(tb))
	tb.Bot.Handle("/create", create(tb))
	tb.Bot.Handle("/rm", deleleSub(tb))
	tb.Bot.Handle("/subscribe", subscribe(tb))

	tb.Bot.Handle("/create_default", createDefault(tb))
	tb.Bot.Handle("/rm_default", removeDefault(tb))
}

// Send files to users
func (tb *TgBot) Send(users []*logic.User, path, caption string) {
	if len(users) == 0 {
		return
	}
	var file telebot.Sendable
	if strings.HasSuffix(path, ".mp4") {
		file = &telebot.Video{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".gif") {
		file = &telebot.Photo{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".webm") {
		err := tb.Downloader.Save(path)
		if err != nil {
			log.Println(err)
			return
		}
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

		newVidPath, err := convertWebmToMp4(tb.Downloader, path)
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

func convertWebmToMp4(d *downloader.Downloader, path string) (string, error) {
	trans := new(transcoder.Transcoder)

	vidPath := d.Get(path)
	newVidPath := strings.TrimSuffix(vidPath, ".webm") + ".mp4"

	err := trans.Initialize(vidPath, newVidPath)
	if err != nil {
		return "", err
	}
	done := trans.Run(false)
	err = <-done
	if err != nil {
		log.Println(err)
		return "", err
	}

	return newVidPath, nil
}
