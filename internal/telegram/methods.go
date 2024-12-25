package telegram

import (
	"log"
	"strconv"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageToAll(message string) error {

	for _, chat := range repository.Chats {

		msg := tgbotapi.NewMessage(chat.ID, message)
		if _, err := Bot.Send(msg); err != nil {
			log.Print(strconv.FormatInt(chat.ID, 10) + " chat - message  not sended because of " + err.Error())
			continue
		}
	}

	return nil
}

func SendMessage(chatID int64, message string) error {

	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := Bot.Send(msg); err != nil {
		log.Print(strconv.FormatInt(chatID, 10) + " chat - message  not sended because of " + err.Error())
	}

	return nil
}
