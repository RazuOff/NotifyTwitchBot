package repository

import (
	"errors"
	"log"
	"strconv"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type Chat struct {
	ID              int64                         `json:"id"`
	TwitchID        string                        `json:"twitch_id"`
	UserAccessToken twitchmodels.UserAccessTokens `json:"user_accessToken"`
	UUID            string                        `json:"uuid"`
}

var Chats []*Chat

func AddChat(chatId int64) {
	if !ChatExists(chatId) {
		Chats = append(Chats, &Chat{ID: chatId})
	}
}

func DeleteChat(chatID int64) (exists bool) {
	for index, chat := range Chats {
		if chat.ID == chatID {
			Chats = append(Chats[:index], Chats[index+1:]...)
			return true
		}
	}
	return false
}

func SetToken(chatID int64, token twitchmodels.UserAccessTokens) error {
	chat, exists := GetChat(chatID)
	if !exists {
		return errors.New("Chat does not exists")
	}
	chat.UserAccessToken = token

	log.Print("CHAT_ID:" + strconv.FormatInt(chatID, 10) + "TOKEN:")
	log.Println(token)

	return nil
}

func SetTwitchID(chatID int64, twitchID string) error {
	chat, exists := GetChat(chatID)
	if !exists {
		return errors.New("Chat does not exists")
	}
	log.Printf("%d Chat sets twitchID = %s", chat.ID, twitchID)
	chat.TwitchID = twitchID
	return nil
}

func SetUUID(chatID int64, uuid string) error {
	chat, exists := GetChat(chatID)
	if !exists {
		return errors.New("Chat does not exists")
	}
	log.Printf("%d Chat sets token = %s", chat.ID, uuid)
	chat.UUID = uuid
	return nil
}

func GetChatByUUID(uuid string) (*Chat, error) {
	for _, chat := range Chats {
		if chat.UUID == uuid {
			chat.UUID = ""
			return chat, nil
		}
	}

	return &Chat{}, errors.New("uuid not found")
}

func GetChatByTwitchID(twitchID string) (*Chat, error) {
	for _, chat := range Chats {
		if chat.TwitchID == twitchID {
			return chat, nil
		}
	}

	return &Chat{}, errors.New("TwitchID not found")
}

func GetChat(chatID int64) (chat *Chat, exists bool) {
	for _, chat := range Chats {
		if chat.ID == chatID {
			return chat, true
		}
	}
	return &Chat{}, false
}

func ChatExists(chatID int64) bool {
	for _, chat := range Chats {
		if chat.ID == chatID {
			return true
		}
	}

	return false
}
