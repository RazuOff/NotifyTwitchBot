package telegram

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const START_COMMAND = "start"
const LOGIN_COMMAND = "login"

func HandleUpdates(updates tgbotapi.UpdatesChannel) {
	go func() {
		for update := range updates {
			if update.Message != nil {
				if update.Message.IsCommand() {
					handleCommands(update.Message)
				} else {
					handleNotCommand(update.Message)
				}
			}

		}
	}()
}

func handleCommands(message *tgbotapi.Message) {
	switch message.Command() {
	case START_COMMAND:
		handleStartCommand(message)
	case LOGIN_COMMAND:
		handleLoginCommand(message)
	default:

	}
}

func handleStartCommand(message *tgbotapi.Message) {
	text := "Привет\nБот находится на стадии тестирования поэтому может сломаться\n\nДля работы необходимо ввести команду /login"
	SendMessage(message.Chat.ID, text)
}

func handleLoginCommand(message *tgbotapi.Message) {
	repository.AddChat(message.Chat.ID)
	SendMessage(message.Chat.ID, "Пройди по ссылке ниже и зайди в аккаунт твича")
	SendMessage(message.Chat.ID, twitch.TwitchAPI.CreateAuthLink(message.Chat.ID))
}

func handleNotCommand(message *tgbotapi.Message) {
	SendMessage(message.Chat.ID, "Введите команду для работы\n\nДоступные команды можно увидеть нажав кнопку \"Меню\"")
}
