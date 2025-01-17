package service

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/apperrors"
	"github.com/gin-gonic/gin"
)

type Redirect interface {
	GetChatFromRedirect(state string) (*models.Chat, error)
	SetUserAccessToken(code string, chat *models.Chat) error
	SetTwitchID(chat *models.Chat) error
	SubscribeToAllStreamUps(chat *models.Chat) (int, error)
	HandleAuthError(chatID int64, text string, c *gin.Context, err error)
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

type Service struct {
	Redirect
	Notify
	View
}

func NewService(repository *repository.Repository) *Service {
	service := Service{View: NewTelegramView(repository)}
	service.Redirect = NewRedirectService(repository, service.View)
	service.Notify = NewNotifyService(repository, service.View)
	return &service
}
