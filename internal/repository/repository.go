package repository

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"gorm.io/gorm"
)

type Chats interface {
	GenerateStateForChat(chatID int64) (string, error)
	AddChat(chatId int64) error
	DeleteChat(chatID int64) error
	DeleteUUID(chat *models.Chat) error
	SetToken(chat *models.Chat, token twitchmodels.UserAccessToken) error
	SetTwitchID(chat *models.Chat, twitchID string) error
	SetUUID(chatID int64, uuid string) error
	GetChatByUUID(uuid string) (*models.Chat, error)
	GetChatByTwitchID(twitchID string) (*models.Chat, error)
	GetChat(chatID int64) (*models.Chat, error)
}

type Follows interface {
	GetChatsByFollow(followID string) ([]models.Chat, error)
	GetFollow(id string) (*models.Follow, error)
	GetUnSubedFollows() ([]models.Follow, error)
	SaveFollow(follow *models.Follow) error
	UpdateSubID(followID string, subID string) error
	AddFollow(chatID int64, follow models.Follow) error
}

type Repository struct {
	Chats
	Follows
}

func NewRepository(db *gorm.DB, twitchAPI *twitch.TwitchAPI) *Repository {
	return &Repository{Chats: NewChatPostgre(db, twitchAPI), Follows: NewFollowsPostgre(db)}
}
