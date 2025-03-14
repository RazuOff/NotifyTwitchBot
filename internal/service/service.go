package service

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/apperrors"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type Redirect interface {
	GetChatFromRedirect(state string) (*models.Chat, error)
	SetUserAccessToken(code string, chat *models.Chat) error
	SetTwitchID(chat *models.Chat) error
	SubscribeToAllStreamUps(chat *models.Chat) (int, error)
	HandleAuthError(chatID int64, text string, err error)
}

type Notify interface {
	SendNotify(broadcasterUserID string) *apperrors.AppError
}

type View interface {
	StartHandlingMessages()
	SendMessage(chatID int64, message string)
	handleStartCommand(chatID int64)
	handleLoginCommand(chatID int64)
	handleFollowsCommand(chatID int64)
	handleExitCommand(chatID int64)
	handleNotCommand(chatID int64)
}

type Debug interface {
	HandleInput()
	printSubs(input string) bool
	deleteSubs(input string) bool
}

type Chat interface {
	GetChatUserAccessToken(chat *models.Chat) (*twitchmodels.UserAccessToken, error)
}

type Service struct {
	Redirect
	Notify
	View
	Debug
	Chat
}

func NewService(repository *repository.Repository, twitchAPI *twitch.TwitchAPI, telegramToken string) *Service {
	service := Service{
		Chat: NewChatService(twitchAPI, repository),
	}

	service.View = NewTelegramView(repository, service.Chat, twitchAPI, telegramToken)
	service.Redirect = NewRedirectService(repository, service.View, service.Chat, twitchAPI)
	service.Notify = NewNotifyService(repository, service.View, twitchAPI)
	service.Debug = NewDebugConsoleService(twitchAPI)
	return &service
}
