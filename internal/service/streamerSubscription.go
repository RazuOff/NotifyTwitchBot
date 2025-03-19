package service

import (
	"context"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
)

type StreamerSubscriptionService struct {
	repository   repository.Streamers
	subscription Subscription
	validator    Validate
}

func NewStreamerSubscriptionService(r repository.Streamers, s Subscription, v Validate) *StreamerSubscriptionService {
	return &StreamerSubscriptionService{repository: r, subscription: s, validator: v}
}

func (service *StreamerSubscriptionService) BuyStreamerSub(streamer models.StreamerAccount) error {
	dbStreamer, err := service.repository.GetStreamerByID(streamer.ID)

	if err != nil {
		return err
	}

	if dbStreamer == nil {
		if err := service.repository.CreateStreamer(&streamer); err != nil {
			return err
		}
		return nil
	}

	if err := service.repository.UpdateStreamer(&streamer); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)

	defer cancel()
	if err := service.subscription.SubscribeStreamUPEvent(ctx, dbStreamer.ID); err != nil {
		return err
	}

	return nil
}

func (service *StreamerSubscriptionService) UnsubNotPayedStreamer(broadcasterID string) (bool, error) {

	isStreamerPayed, err := service.validator.IsSubscriptionActive(broadcasterID)
	if err != nil {
		return false, err
	}

	if !isStreamerPayed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := service.subscription.UnsubscribeStreamUPEvent(ctx, broadcasterID); err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}
