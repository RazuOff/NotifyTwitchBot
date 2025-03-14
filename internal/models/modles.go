package models

import (
	"time"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

// Нужно разделить модель на DBModel и просто Model
type Chat struct {
	ID              int64                         `json:"id"`
	TwitchID        string                        `json:"twitch_id"`
	UserAccessToken *twitchmodels.UserAccessToken `gorm:"type:jsonb;default:null" json:"user_accessToken"`
	UUID            string                        `json:"uuid"`
	Follows         []Follow                      `gorm:"many2many:chat_follows;constraint:OnDelete:CASCADE;"`
}

func (chat *Chat) GetUserAccessToken() (*twitchmodels.UserAccessToken, bool) {
	buffer := 120 * time.Second
	currentTime := time.Now()
	tokenExpiry := chat.UserAccessToken.CreatedAt.Add(time.Duration(chat.UserAccessToken.ExpiresIn) * time.Second)
	isTokenNearExpiry := currentTime.After(tokenExpiry.Add(-buffer))

	if !isTokenNearExpiry {
		return chat.UserAccessToken, false
	}
	return chat.UserAccessToken, true

}

type Follow struct {
	ID              string `gorm:"unique;not null" json:"id"`
	BroadcasterName string `json:"broadcaster_name"`
	Subscribtion_id string `json:"subscribtion_id"`
	Chats           []Chat `gorm:"many2many:chat_follows;constraint:OnDelete:CASCADE;"`
}
