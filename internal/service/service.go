package service

import (
	"context"
	"sync"

	"github.com/RazuOff/NotifyTwitchBot/internal/config"
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/apperrors"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type Subscription interface {
	SubscribeToAllStreamUps(chat *models.Chat) (notPayedStreamers int, errCount int, err error)
	SubscribeStreamUPEvent(ctx context.Context, broadcasterID string) error
	UnsubscribeStreamUPEvent(ctx context.Context, broadcasterID string) error
	subscribeToTwitchEvents(ctx context.Context, cancel context.CancelFunc, sem chan struct{}, wg *sync.WaitGroup, mu *sync.Mutex, errChan chan<- (error), subsError *int, apiError *error, follow models.Follow)
}

type StreamerSubscription interface {
	BuyStreamerSub(id string, days *int) error
	UnsubNotPayedStreamer(broadcasterID string) (done bool, err error)
	GetStreamerInfo(broadcasterID string) (*models.StreamerAccount, error)
}

type Validate interface {
	IsSubscriptionActive(streamerID string) (bool, error)
}

type Chat interface {
	GetChatUserAccessToken(chat *models.Chat) (*twitchmodels.UserAccessToken, error)
	SetUserAccessToken(code string, chat *models.Chat) error
	SetTwitchID(chat *models.Chat) error
}

type Debug interface {
	HandleInput()
	printSubs(input string) bool
	deleteSubs(input string) bool
}

type Redirect interface {
	GetChatFromRedirect(state string) (*models.Chat, error)
	HandleAuthError(chatID int64, text string, err error)
}

type Notify interface {
	SendNotify(broadcasterUserID string) *apperrors.AppError
	//NotifyAboutNotPayedStreamers() error
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
	Subscription
	Chat
	StreamerSubscription
	Validate
	Debug

	Redirect
	Notify
	View
}

func NewService(repository *repository.Repository, twitchAPI *twitch.TwitchAPI, config *config.Config) *Service {
	service := Service{
		Chat: NewChatService(twitchAPI, repository),
	}
	service.Validate = NewValidateService(repository.Streamers)

	service.Subscription = NewSubscrpitionService(repository, twitchAPI, service.Chat, service.Validate, config)
	service.StreamerSubscription = NewStreamerSubscriptionService(repository, service.Subscription, service.Validate)
	service.View = NewTelegramView(repository, service.Chat, twitchAPI, config.TelegramToken)
	service.Redirect = NewRedirectService(repository, service.View, service.Chat)
	service.Notify = NewNotifyService(repository, service.View, twitchAPI)
	service.Debug = NewDebugConsoleService(twitchAPI)
	return &service
}
