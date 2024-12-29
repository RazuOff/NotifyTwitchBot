package telegram

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const START_COMMAND = "start"
const LOGIN_COMMAND = "login"
const FOLLOW_COMMAND = "follows"

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
	case FOLLOW_COMMAND:
		handleFollowsCommand(message)
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

func handleFollowsCommand(message *tgbotapi.Message) {
	chat, exists := repository.GetChat(message.Chat.ID)
	if !exists {
		SendMessage(message.Chat.ID, "Для начала войди в аккаунт")
		return
	}
	follows, err := twitch.TwitchAPI.GetAccountFollows(chat.TwitchID)
	if err != nil {
		SendMessage(message.Chat.ID, "Что-то пошло не так, попробуй позже(")
		log.Println("GetAccountFollows error=" + err.Error())
		return
	}

	for _, follow := range follows {
		SendMessage(message.Chat.ID, follow.BroadcasterLogin)
	}
}

func handleNotCommand(message *tgbotapi.Message) {
	SendMessage(message.Chat.ID, "Введите команду для работы\n\nДоступные команды можно увидеть нажав кнопку \"Меню\"")
}
