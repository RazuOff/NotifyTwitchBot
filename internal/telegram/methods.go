package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(chatID int64, message string) {

	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := Bot.Send(msg); err != nil {
		log.Printf("%d chat - message not sended, err: %s", chatID, err.Error())
		//log.Print(strconv.FormatInt(chatID, 10) + " chat - message  not sended because of " + err.Error())
	}
}
