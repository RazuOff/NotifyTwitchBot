package repository

import (
	"fmt"
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"gorm.io/gorm/clause"
)

func AddChat(chatId int64) error {

	chat := models.Chat{ID: chatId}

	if err := postgre.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&chat).Error; err != nil {
		log.Printf("Error while inserting chat: %v", err)
		return fmt.Errorf("AddChat error: %w", err)
	}

	return nil
}

func DeleteChat(chatID int64) error {

	var chat models.Chat
	if err := postgre.DB.Preload("Follows").First(&chat, chatID).Error; err != nil {
		return err
	}

	// Получить IDs связанных Follow перед удалением
	var followIDs []string
	for _, follow := range chat.Follows {
		followIDs = append(followIDs, follow.ID)
	}

	// Удалить запись Chat и связи в chat_follows
	if err := postgre.DB.Delete(&chat).Error; err != nil {
		return err
	}

	if err := postgre.DB.Delete(&models.Follow{}, "id IN ? AND NOT EXISTS (SELECT 1 FROM chat_follows WHERE chat_follows.follow_id = id)", followIDs).Error; err != nil {
		return err
	}

	//ДОБАВИТЬ ОТМЕНУ ПОДПИСИ НА WEBHOOK

	return nil
}

func SetToken(chat *models.Chat, token twitchmodels.UserAccessTokens) error {

	chat.UserAccessToken = token

	if err := postgre.DB.Save(chat).Error; err != nil {
		log.Printf("SetToken save error")
		return err
	}
	log.Printf("CHAT_ID:  %d TOKEN: %+v", chat.ID, token)

	return nil
}

func SetTwitchID(chat *models.Chat, twitchID string) error {

	chat.TwitchID = twitchID
	if err := postgre.DB.Save(chat).Error; err != nil {
		log.Printf("SetTwitchID save error")
		return err
	}
	log.Printf("%d Chat sets twitchID = %s", chat.ID, twitchID)

	return nil
}

func SetUUID(chatID int64, uuid string) error {
	chat, err := GetChat(chatID)
	if err != nil {
		log.Printf("SetUUID error")
		return err
	}

	if chat == nil {
		return fmt.Errorf("сhat not found")
	}

	chat.UUID = uuid
	if err := postgre.DB.Save(chat).Error; err != nil {
		log.Printf("SetUUID save error")
		return err
	}
	log.Printf("%d Chat sets token = %s", chat.ID, uuid)

	return nil
}

func GetChatByUUID(uuid string) (*models.Chat, error) {

	var chat models.Chat

	if err := postgre.DB.Where("uuid = ?", uuid).Find(&chat).Error; err != nil {
		log.Printf("GetChatByUUID error")
		return nil, err
	}

	return &chat, nil
}

func GetChatByTwitchID(twitchID string) (*models.Chat, error) {
	var chat models.Chat

	if err := postgre.DB.Where("twitch_id = ?", twitchID).Find(&chat).Error; err != nil {
		log.Printf("GetChatByTwitchID error")
		return nil, err
	}

	return &chat, nil
}

func GetChat(chatID int64) (*models.Chat, error) {

	var chat models.Chat

	if err := postgre.DB.Where("id = ?", chatID).Find(&chat).Error; err != nil {
		log.Printf("GetChat error")
		return nil, err
	}

	return &chat, nil
}
