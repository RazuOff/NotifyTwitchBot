package telegram

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func StartTelegramBot() {

	Bot = Init()

	Bot.Debug = false

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)
	HandleUpdates(updates)
}

func Init() *tgbotapi.BotAPI {
	botToken, exists := os.LookupEnv("TELEGRAM_API_TOKEN")
	if !exists {
		log.Fatal("Telegram bot token not found in .env file")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	return bot
}
