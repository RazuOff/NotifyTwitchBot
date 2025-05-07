package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
)

type StreamerSubscriptionService struct {
	repository   *repository.Repository
	subscription Subscription
	validator    Validate
}

func NewStreamerSubscriptionService(r *repository.Repository, s Subscription, v Validate) *StreamerSubscriptionService {
	return &StreamerSubscriptionService{repository: r, subscription: s, validator: v}
}

func (service *StreamerSubscriptionService) BuyStreamerSub(id string, days *int) error {
	dbStreamer, err := service.repository.GetStreamerByID(id)

	if err != nil {
		return err
	}

	now := time.Now()

	if dbStreamer == nil {
		if err := service.repository.CreateStreamer(&models.StreamerAccount{ID: id, SubedAt: &now, SubDaysDuration: days}); err != nil {
			return err
		}
		return nil
	}

	endDate := dbStreamer.SubedAt.AddDate(0, 0, *dbStreamer.SubDaysDuration)

	diff := endDate.Sub(now.UTC())
	daysLeft := int(diff.Hours()) / 24
	remainder := diff - time.Duration(daysLeft)*24*time.Hour

	subedAt := time.Now().Add(remainder)
	subDaysDuration := *days + daysLeft

	streamer := models.StreamerAccount{
		ID:              id,
		SubedAt:         &subedAt,
		SubDaysDuration: &subDaysDuration,
	}

	if err := service.repository.UpdateStreamer(&streamer); err != nil {
		return err
	}

	follow, err := service.repository.GetFollow(streamer.ID)
	if err != nil {
		return fmt.Errorf("BuyStreamerSub err: %w", err)
	}

	if follow == nil || follow.Subscribtion_id != "" {
		return nil
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

func (service *StreamerSubscriptionService) GetStreamerInfo(broadcasterID string) (*models.StreamerAccount, error) {
	info, err := service.repository.GetStreamerByID(broadcasterID)

	if err != nil {
		log.Printf("GetStreamerInfo err: %s\n", err.Error())
		return nil, err
	}

	return info, nil
}
