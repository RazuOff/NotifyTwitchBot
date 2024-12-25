package telegram

import (
	"log"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func StartTelegramBot() {

	Bot = Init()

	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil { // If we got a message
				// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				repository.AddChat(update.Message.Chat.ID)
				SendMessage(update.Message.Chat.ID, twitch.TwitchAPI.CreateAuthLink(update.Message.Chat.ID))
				// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				// msg.ReplyToMessageID = update.Message.MessageID

				// Bot.Send(msg)
			}
		}
	}()
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
