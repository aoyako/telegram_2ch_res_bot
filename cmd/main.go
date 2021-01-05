package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aoyako/telegram_2ch_res_bot/initialize"
)

func main() {
	log.Println("Starting...")
	bot, _, _ := initialize.App()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go bot.Bot.Start()
	// go initialize.StartPolling(apicnt, duration)

	log.Println("Started")

	<-quit
	log.Println("Quit")
}
