package service

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/apperrors"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
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

type Service struct {
	Redirect
	Notify
	View
	Debug
}

func NewService(repository *repository.Repository, twitchAPI *twitch.TwitchAPI) *Service {
	service := Service{View: NewTelegramView(repository, twitchAPI)}
	service.Redirect = NewRedirectService(repository, service.View, twitchAPI)
	service.Notify = NewNotifyService(repository, service.View, twitchAPI)
	service.Debug = NewDebugConsoleService(twitchAPI)
	return &service
}
