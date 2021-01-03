package telegram

import "github.com/aoyako/telegram_2ch_res_bot/logic"

// Sender can send files to users
type Sender interface {
	Send(user []*logic.User, path, caption string)
}
