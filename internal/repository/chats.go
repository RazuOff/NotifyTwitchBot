package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func GenerateStateForChat(chatID int64) (string, error) {
	uuid := uuid.New().String()
	if err := SetUUID(chatID, uuid); err != nil {
		return "", fmt.Errorf("GenerateState error %w", err)
	}
	return uuid, nil
}

func AddChat(chatId int64) error {
	chat := models.Chat{ID: chatId}

	if err := postgre.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&chat).Error; err != nil {
		log.Printf("Error while inserting chat: %v", err)
		return fmt.Errorf("AddChat error: %w", err)
	}

	return nil
}

// Delete chats and follows that are not used anymore
func DeleteChat(chatID int64) error {
	tx := postgre.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var chat models.Chat
	if err := tx.Preload("Follows").First(&chat, chatID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteChat error: %w", err)
	}

	var followIDs []string
	for _, follow := range chat.Follows {
		followIDs = append(followIDs, follow.ID)
	}

	if err := tx.Delete(&chat).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteChat error: %w", err)
	}

	var unusedFollows []models.Follow
	if err := tx.Find(&unusedFollows, "id IN ? AND NOT EXISTS (SELECT 1 FROM chat_follows WHERE chat_follows.follow_id = id)", followIDs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteChat error: %w", err)
	}

	if len(unusedFollows) == 0 {
		tx.Commit()
		return nil
	}

	var followsToDelete []models.Follow
	for _, follow := range unusedFollows {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		if err := twitch.TwitchAPI.DeleteEventSub(ctx, follow.Subscribtion_id); err != nil {
			log.Printf("DeleteChat error: %s", err.Error())
		} else {
			followsToDelete = append(followsToDelete, follow)
		}
		cancel()
	}

	if err := tx.Delete(&followsToDelete).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteChat error: %w", err)
	}

	tx.Commit()

	return nil
}

func DeleteUUID(chat *models.Chat) error {
	chat.UUID = ""

	if err := postgre.DB.Save(chat).Error; err != nil {
		log.Printf("DeleteUUID save error")
		return err
	}

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
