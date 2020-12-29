package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aoyako/telegram_2ch_res_bot/logic"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// /start endpoint
func start(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(uint64(m.Chat.ID))
		if err != nil {

		}

		help(tb)(m)
	}
}

// /list endpoint
func list(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		subs, err := tb.Controller.Subscription.GetSubsByChatID(uint64(m.Chat.ID))
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
		}
		result := marshallSubs(subs)
		tb.Bot.Send(m.Sender, result)
	}
}

// /help endpoint
func help(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		tb.Bot.Send(m.Sender, HELP_MESSAGE, telebot.ModeMarkdownV2)
	}
}

// /add endpoint
func add(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.Add(uint64(m.Chat.ID), args)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// /del endpoint
func del(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		num, err := strconv.Atoi(args)
		num--
		if err != nil || num < 0 {
			tb.Bot.Send(m.Sender, "Bad index")
			return
		}

		err = tb.Controller.Subscription.Remove(uint64(m.Chat.ID), uint(num))
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad index")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// Format command as ([comand_name] [command_text])
func parseCommand(cmd string) (string, error) {
	separator := regexp.MustCompile(` `)
	args := separator.Split(cmd, 2)
	if len(args) != 2 {
		return "", errors.New("Bad request")
	}

	return args[1], nil
}

// Format []logic.Publication to string
func marshallSubs(subs []logic.Publication) string {
	result := "Current subscriptions:"
	for id, sub := range subs {
		result = fmt.Sprintf("%s\n%d: %s", result, id+1, marshallSub(sub))
	}
	return result
}

// Format logic.Publication to string
func marshallSub(sub logic.Publication) string {
	return fmt.Sprintf("/%s %s %s", sub.Board, sub.Type, sub.Tags)
}
