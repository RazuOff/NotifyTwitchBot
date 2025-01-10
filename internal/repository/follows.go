package repository

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
)

func GetChatsByFollow(followID string) ([]models.Chat, error) {

	var chats []models.Chat

	if err := postgre.DB.Preload("Follows", "id = ?", followID).Find(&chats).Error; err != nil {
		log.Printf("GetChatsByFollow error")
		return nil, err
	}

	return chats, nil
}

func GetFollow(id string) (*models.Follow, error) {

	var follow models.Follow

	if err := postgre.DB.Where("id = ?", id).Find(&follow).Error; err != nil {
		log.Printf("GetFollow error")
		return nil, err
	}

	return &follow, nil
}

func GetUnSubedFollows() ([]models.Follow, error) {

	var follows []models.Follow

	if err := postgre.DB.Where("is_subscribed = ?", false).Find(&follows).Error; err != nil {
		log.Printf("GetFollow error")
		return nil, err
	}

	return follows, nil
}

func SaveFollow(follow *models.Follow) error {

	if err := postgre.DB.Save(follow).Error; err != nil {
		log.Printf("SaveFollow  error")
		return err
	}

	return nil
}

func AddFollow(chatID int64, follow models.Follow) error {

	if err := postgre.DB.FirstOrCreate(&follow).Error; err != nil {
		log.Printf("Error creating or finding follow: %v", err)
		return err
	}

	var chat models.Chat
	if err := postgre.DB.First(&chat, chatID).Error; err != nil {
		log.Printf("Error finding chat with ID %d: %v", chatID, err)
		return err
	}

	if err := postgre.DB.Model(&chat).Association("Follows").Append(&follow); err != nil {
		log.Printf("Error associating follow with chat: %v", err)
		return err
	}

	return nil
}
