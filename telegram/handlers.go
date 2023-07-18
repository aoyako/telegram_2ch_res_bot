package telegram

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/aoyako/telegram_2ch_res_bot/logic"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// /start endpoint
func start(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {
			log.Println(err)
		}

		help(tb)(m)
	}
}

// /subs endpoint
func subs(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		subs, err := tb.Controller.Subscription.GetSubsByChatID(m.Chat.ID)
		if err != nil {
			log.Println(err)
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}
		result := fmt.Sprintf("Your subs:%s", marshallSubs(subs, true))
		tb.Bot.Send(m.Sender, result)
	}
}

// /list endpoint
func list(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		subs := tb.Controller.Subscription.GetAllDefaultSubs()
		result := fmt.Sprintf("Available subs:%s", marshallSubs(subs, true))
		tb.Bot.Send(m.Sender, result)
	}
}

// /clist endpoint
func cleverList(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		subs := tb.Controller.Subscription.GetAllDefaultSubs()
		result := fmt.Sprintf("Available subs:%s", marshallSubs(subs, false))
		tb.Bot.Send(m.Sender, result)
	}
}

// /help endpoint
func help(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		tb.Bot.Send(m.Sender, HelpMessage, telebot.ModeMarkdown)
	}
}

// /add endpoint
func create(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.AddNew(m.Chat.ID, args)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// /subscribe endpoint
func subscribe(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.Subscription.Subscribe(m.Chat.ID, args)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// /create_default endpoint
func createDefault(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.Create(m.Chat.ID, args)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// /rm endpoint
func deleleSub(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.Subscription.Remove(m.Chat.ID, args)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad index")
			return
		}

		tb.Bot.Send(m.Sender, "OK")
	}
}

// /rm_default endpoint
func removeDefault(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		args, err := parseCommand(m.Text)
		if err != nil {
			tb.Bot.Send(m.Sender, "Bad request")
			return
		}

		err = tb.Controller.Subscription.RemoveDefault(m.Chat.ID, args)
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
		return "", errors.New("bad request")
	}

	return args[1], nil
}

// Format []logic.Publication to string
func marshallSubs(subs []logic.Publication, displayAlias bool) string {
	result := ""
	for id, sub := range subs {
		if sub.Alias != "" && displayAlias {
			result = fmt.Sprintf("%s\n%d: %s", result, id+1, sub.Alias)
		} else {
			result = fmt.Sprintf("%s\n%d: %s", result, id+1, marshallSub(sub))
		}
	}
	return result
}

// Format logic.Publication to string
func marshallSub(sub logic.Publication) string {
	return fmt.Sprintf("/%s %s %s", sub.Board, sub.Type, sub.Tags)
}
