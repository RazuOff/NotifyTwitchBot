package repository

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/gorm"
)

type FollowsPostgre struct {
	DB *gorm.DB
}

func NewFollowsPostgre(db *gorm.DB) *FollowsPostgre {
	return &FollowsPostgre{DB: db}
}

func (repository *FollowsPostgre) GetChatsByFollow(followID string) ([]models.Chat, error) {

	var chats []models.Chat

	if err := repository.DB.
		Joins("JOIN chat_follows ON chat_follows.chat_id = chats.id").
		Where("chat_follows.follow_id = ?", followID).
		Preload("Follows").
		Find(&chats).Error; err != nil {
		log.Printf("GetChatsByFollow error: %v", err)
		return nil, err
	}

	return chats, nil
}

func (repository *FollowsPostgre) GetFollow(id string) (*models.Follow, error) {

	var follow models.Follow

	if err := repository.DB.Where("id = ?", id).Find(&follow).Error; err != nil {
		log.Printf("GetFollow error")
		return nil, err
	}

	return &follow, nil
}

func (repository *FollowsPostgre) GetUnSubedFollows() ([]models.Follow, error) {

	var follows []models.Follow

	if err := repository.DB.Where("subscribtion_id = ?", "").Find(&follows).Error; err != nil {
		log.Printf("GetFollow error")
		return nil, err
	}

	return follows, nil
}

func (repository *FollowsPostgre) SaveFollow(follow *models.Follow) error {

	if err := repository.DB.Save(follow).Error; err != nil {
		log.Printf("SaveFollow  error")
		return err
	}

	return nil
}

func (repository *FollowsPostgre) UpdateSubID(followID string, subID string) error {
	if err := repository.DB.Model(&models.Follow{}).Where("id = ?", followID).Update("subscribtion_id", subID).Error; err != nil {
		log.Print("UpdateSubscribtionId error")
		return err
	}
	return nil
}

func (repository *FollowsPostgre) AddFollow(chatID int64, follow models.Follow) error {

	if err := repository.DB.FirstOrCreate(&follow).Error; err != nil {
		log.Printf("Error creating or finding follow: %v", err)
		return err
	}

	var chat models.Chat
	if err := repository.DB.First(&chat, chatID).Error; err != nil {
		log.Printf("Error finding chat with ID %d: %v", chatID, err)
		return err
	}

	if err := repository.DB.Model(&chat).Association("Follows").Append(&follow); err != nil {
		log.Printf("Error associating follow with chat: %v", err)
		return err
	}

	return nil
}
