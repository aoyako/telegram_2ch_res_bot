package initialize

import (
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/dvach"
)

// StartPolling starts file sending
func StartPolling(api *dvach.APIController, minutes uint64) {
	for {
		go api.InitiateSending()
		<-time.After(time.Duration(minutes) * time.Minute)
	}
}
