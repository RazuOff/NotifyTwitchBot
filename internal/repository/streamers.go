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

func (repository *StreamersGORM) GetStreamerByID(broadcasterID string) (*models.StreamerAccount, error) {

	var streamerAccount models.StreamerAccount
	if err := repository.DB.Take(&streamerAccount, "id = ?", broadcasterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("GetStreamerByID error: %w", err)
	}

	return &streamerAccount, nil
}

func (repository *StreamersGORM) UpdateStreamer(streamer *models.StreamerAccount) error {
	result := repository.DB.Model(streamer).Updates(*streamer)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("streamer not found")
	}
	return nil
}

func (repository *StreamersGORM) CreateStreamer(streamer *models.StreamerAccount) error {
	return repository.DB.FirstOrCreate(streamer).Error
}
