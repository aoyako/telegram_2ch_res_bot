package telegram

import "github.com/aoyako/telegram_2ch_res_bot/logic"

type Sender interface {
	Send(user []*logic.User, path, caption string)
}
