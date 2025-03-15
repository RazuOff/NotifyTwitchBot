package service

import (
	"fmt"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
)

type ValidateService struct {
	repository repository.Streamers
}

func NewValidateService(streamers repository.Streamers) *ValidateService {
	return &ValidateService{repository: streamers}
}

func (service *ValidateService) IsSubscriptionActive(streamerID string) (bool, error) {
	streamer, err := service.repository.GetStreamerByID(streamerID)
	if err != nil {
		return false, fmt.Errorf("CheckStreamer error %w", err)
	}

	if streamer == nil {
		return false, nil
	}

	expirationTime := streamer.SubedAt.Add(time.Duration(streamer.SubDaysDuration) * 24 * time.Hour)
	return time.Now().Before(expirationTime), nil

}
