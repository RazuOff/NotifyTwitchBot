package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatsPostgre struct {
	DB        *gorm.DB
	twitchAPI *twitch.TwitchAPI
}

func NewChatPostgre(db *gorm.DB, api *twitch.TwitchAPI) *ChatsPostgre {
	return &ChatsPostgre{DB: db, twitchAPI: api}
}

func (repository *ChatsPostgre) GenerateStateForChat(chatID int64) (string, error) {
	uuid := uuid.New().String()
	if err := repository.SetUUID(chatID, uuid); err != nil {
		return "", fmt.Errorf("GenerateState error %w", err)
	}
	return uuid, nil
}

func (repository *ChatsPostgre) AddChat(chatId int64) error {
	chat := models.Chat{ID: chatId}

	if err := repository.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&chat).Error; err != nil {
		log.Printf("Error while inserting chat: %v", err)
		return fmt.Errorf("AddChat error: %w", err)
	}

	return nil
}

// Delete chats and follows that are not used anymore
func (repository *ChatsPostgre) DeleteChat(chatID int64) error {
	tx := repository.DB.Begin()

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
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, follow := range unusedFollows {
		wg.Add(1)
		go func(f models.Follow) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			if err := repository.twitchAPI.DeleteEventSub(ctx, f.Subscribtion_id); err != nil {
				log.Printf("DeleteChat error: %s", err.Error())
			} else {
				mutex.Lock()
				followsToDelete = append(followsToDelete, f)
				mutex.Unlock()
			}
			cancel()
		}(follow)
	}

	wg.Wait()
	if err := tx.Delete(&followsToDelete).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteChat error: %w", err)
	}

	tx.Commit()

	return nil
}

func (repository *ChatsPostgre) DeleteUUID(chat *models.Chat) error {
	chat.UUID = ""

	if err := repository.DB.Save(chat).Error; err != nil {
		log.Printf("DeleteUUID save error")
		return err
	}

	return nil
}

func (repository *ChatsPostgre) SetToken(chat *models.Chat, token twitchmodels.UserAccessToken) error {
	chat.UserAccessToken = &token

	if err := repository.DB.Save(chat).Error; err != nil {
		log.Printf("SetToken save error")
		return err
	}
	log.Printf("CHAT_ID:  %d TOKEN: %+v", chat.ID, token)

	return nil
}

func (repository *ChatsPostgre) SetTwitchID(chat *models.Chat, twitchID string) error {
	chat.TwitchID = twitchID
	if err := repository.DB.Save(chat).Error; err != nil {
		log.Printf("SetTwitchID save error")
		return err
	}
	log.Printf("%d Chat sets twitchID = %s", chat.ID, twitchID)

	return nil
}

func (repository *ChatsPostgre) SetUUID(chatID int64, uuid string) error {
	chat, err := repository.GetChat(chatID)
	if err != nil {
		log.Printf("SetUUID error")
		return err
	}

	if chat == nil {
		return fmt.Errorf("сhat not found")
	}

	chat.UUID = uuid
	if err := repository.DB.Save(chat).Error; err != nil {
		log.Printf("SetUUID save error")
		return err
	}
	log.Printf("%d Chat sets token = %s", chat.ID, uuid)

	return nil
}

func (repository *ChatsPostgre) GetChatByUUID(uuid string) (*models.Chat, error) {

	var chat models.Chat

	if err := repository.DB.Where("uuid = ?", uuid).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("GetChatByUUID error")
		return nil, err
	}

	return &chat, nil
}

func (repository *ChatsPostgre) GetChatByTwitchID(twitchID string) (*models.Chat, error) {
	var chat models.Chat

	if err := repository.DB.Where("twitch_id = ?", twitchID).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("GetChatByTwitchID error")
		return nil, err
	}

	return &chat, nil
}

func (repository *ChatsPostgre) GetChat(chatID int64) (*models.Chat, error) {
	var chat models.Chat

	if err := repository.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("GetChat error")
		return nil, err
	}

	return &chat, nil
}
