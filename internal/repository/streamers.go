package repository

import (
	"errors"
	"fmt"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/gorm"
)

type StreamersGORM struct {
	DB *gorm.DB
}

func NewStreamersGORMRepository(db *gorm.DB) *StreamersGORM {
	return &StreamersGORM{DB: db}
}

func (repository StreamersGORM) GetStreamerByID(broadcasterID string) (*models.StreamerAccount, error) {

	var streamerAccount models.StreamerAccount
	if err := repository.DB.Take(&streamerAccount, "id = ?", broadcasterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("GetStreamerByID error: %w", err)
	}

	return &streamerAccount, nil
}
