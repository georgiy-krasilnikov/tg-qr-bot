package bot

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	"tg_qr_bot/qr"
)

var tgBotToken = "5964576520:AAGWq_FksatTyahDc8Qlu-PmuVAVM3Y93BI"

func TelegramBot() error {
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		return fmt.Errorf("failed to create new bot API: %s", err.Error())
	}

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		return fmt.Errorf("failed to create channel fo updates: %s", err.Error())
	}

	for upd := range updates {
		if upd.Message == nil {
			continue
		}

		if reflect.TypeOf(upd.Message.Text).Kind() == reflect.String && upd.Message.Text != "" {

			switch upd.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Отправь мне ссылку, и ты увидишь, что произойдет:)")
				_, err := bot.Send(msg)
				if err != nil {
					return fmt.Errorf("failed to send message: %s", err.Error())
				}

			default:
				_, err := url.ParseRequestURI(upd.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Это не ссылка!:(")
					_, err := bot.Send(msg)
					if err != nil {
						return fmt.Errorf("failed to send message: %s", err.Error())
					}
				} else {
					msg, fileName, err := CreateQRMessage(upd)
					if err != nil {
						return fmt.Errorf("failed to create message with QR file: %s", err.Error())
					}
					_, err = bot.Send(msg)
					if err != nil {
						return fmt.Errorf("failed to send message: %s", err.Error())
					}
					if err = os.Remove(fileName); err != nil {
						return fmt.Errorf("failed to delete QR file: %s", err.Error())
					}
				}
			}
		}
	}
	return nil
}

func CreateQRMessage(upd tgbotapi.Update) (tgbotapi.PhotoConfig, string, error) {
	fileName, err := qr.CreateQR(upd.Message.Text)
	if err != nil {
		return tgbotapi.PhotoConfig{}, "", fmt.Errorf("failed to create QR file with this url: %s", err.Error())
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return tgbotapi.PhotoConfig{}, "", fmt.Errorf("failed to read QR file: %s", err.Error())
	}

	bytes := tgbotapi.FileBytes{Name: fileName, Bytes: file}
	msg := tgbotapi.NewPhotoUpload(upd.Message.Chat.ID, bytes)

	if upd.Message.From.UserName != "" {
		msg.Caption = "@" + upd.Message.From.UserName + ", вот твой QR код 👆"
	} else {
		msg.Caption = "Вот твой QR код 👆"
	}

	return msg, fileName, nil
}
