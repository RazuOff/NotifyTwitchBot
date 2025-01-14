package models

import twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"

type Chat struct {
	ID              int64                         `json:"id"`
	TwitchID        string                        `json:"twitch_id"`
	UserAccessToken twitchmodels.UserAccessTokens `gorm:"type:jsonb" json:"user_accessToken"`
	UUID            string                        `json:"uuid"`
	Follows         []Follow                      `gorm:"many2many:chat_follows;constraint:OnDelete:CASCADE;"`
}

type Follow struct {
	ID              string `gorm:"unique;not null" json:"id"`
	BroadcasterName string `json:"broadcaster_name"`
	Subscribtion_id string `json:"subscribtion_id"`
	Chats           []Chat `gorm:"many2many:chat_follows;constraint:OnDelete:CASCADE;"`
}
