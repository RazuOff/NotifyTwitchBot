package service

import (
	"log"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const START_COMMAND = "start"
const LOGIN_COMMAND = "login"
const FOLLOW_COMMAND = "follows"
const EXIT_COMMAND = "exit"

type TelegramView struct {
	repository repository.Chats
	api        *tgbotapi.BotAPI
}

func NewTelegramView(repo repository.Chats) *TelegramView {

	botToken, exists := os.LookupEnv("TELEGRAM_API_TOKEN")
	if !exists {
		log.Fatal("Telegram view token not found in .env file")
	}

	view, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	return &TelegramView{repository: repo, api: view}
}

func (view *TelegramView) StartHandlingMessages() {

	view.api.Debug = false

	log.Printf("Authorized on account %s", view.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := view.api.GetUpdatesChan(u)
	view.handleUpdates(updates)
}

func (view *TelegramView) handleUpdates(updates tgbotapi.UpdatesChannel) {
	go func() {
		for update := range updates {
			if update.Message != nil {
				if update.Message.IsCommand() {
					go view.handleCommands(update.Message)
				} else {
					go view.handleNotCommand(update.Message.Chat.ID)
				}
			}

		}
	}()
}

func (view *TelegramView) handleCommands(message *tgbotapi.Message) {
	switch message.Command() {
	case START_COMMAND:
		view.handleStartCommand(message.Chat.ID)
	case LOGIN_COMMAND:
		view.handleLoginCommand(message.Chat.ID)
	case FOLLOW_COMMAND:
		view.handleFollowsCommand(message.Chat.ID)
	case EXIT_COMMAND:
		view.handleExitCommand(message.Chat.ID)
	default:

	}
}

func (view *TelegramView) handleStartCommand(chatID int64) {
	text := "Привет\nБот находится на стадии тестирования поэтому может сломаться\n\nДля работы необходимо ввести команду /login"
	view.SendMessage(chatID, text)
}

func (view *TelegramView) handleLoginCommand(chatID int64) {

	chat, err := view.repository.GetChat(chatID)
	if err != nil {
		log.Printf("handleLoginCommand error: %s", err)
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}

	if chat != nil && chat.TwitchID != "" && chat.UserAccessToken.AccessToken != "" {
		view.SendMessage(chatID, "Вы уже вошли в аккаунт")
		return
	}

	if err := view.repository.AddChat(chatID); err != nil {
		log.Printf("handleLoginCommand error: %s", err.Error())
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}
	view.SendMessage(chatID, "Пройди по ссылке ниже и зайди в аккаунт твича")

	uuid, err := view.repository.GenerateStateForChat(chatID)
	if err != nil {
		log.Printf("handleLoginCommand error: %s", err.Error())
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}

	view.SendMessage(chatID, twitch.TwitchAPI.CreateAuthLink(chatID, uuid))
}

func (view *TelegramView) handleFollowsCommand(chatID int64) {
	chat, err := view.repository.GetChat(chatID)
	if err != nil {
		log.Printf("handleFollowsCommand error")
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}

	if chat == nil {
		view.SendMessage(chatID, "Для начала войди в аккаунт")
		return
	}
	follows, err := twitch.TwitchAPI.GetAccountFollows(chat.TwitchID, chat.UserAccessToken)
	if err != nil {
		view.SendMessage(chatID, "Что-то пошло не так, попробуй позже(")
		log.Println("GetAccountFollows error=" + err.Error())
		return
	}

	for _, follow := range follows {
		view.SendMessage(chatID, follow.BroadcasterName)
	}
}

func (view *TelegramView) handleExitCommand(chatID int64) {
	chat, err := view.repository.GetChat(chatID)
	if err != nil {
		log.Printf("handleExitCommand error")
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}

	if chat.ID == 0 {
		view.SendMessage(chatID, "Для начала войди в аккаунт")
		return
	}
	if err := view.repository.DeleteChat(chat.ID); err != nil {
		log.Printf("handleExitCommand error")
		view.SendMessage(chatID, "У нас сломалась БД(")
		return
	}

	view.SendMessage(chatID, "Вы успешно вышли из аккаунта!")
}

func (view *TelegramView) handleNotCommand(chatID int64) {
	view.SendMessage(chatID, "Введите команду для работы\n\nДоступные команды можно увидеть нажав кнопку \"Меню\"")
}

func (view *TelegramView) SendMessage(chatID int64, message string) {

	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := view.api.Send(msg); err != nil {
		log.Printf("%d chat - message not sended, err: %s", chatID, err.Error())
	}
}
