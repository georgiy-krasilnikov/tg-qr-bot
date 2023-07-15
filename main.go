package main

import (
	log "github.com/sirupsen/logrus"

	"tg_qr_bot/bot"
)

func main() {
	if err := bot.TelegramBot(); err != nil {
		log.WithError(err).Fatal("failed to launch tg bot")
	}
}
