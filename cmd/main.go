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
	bot, apicnt, duration := initialize.App()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go bot.Bot.Start()
	go initialize.StartPolling(apicnt, duration)

	log.Println("Started")

	<-quit
	log.Println("Quit")
}
