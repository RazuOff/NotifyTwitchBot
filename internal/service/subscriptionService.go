package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/config"
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

const (
	MAX_EVENTGO = 20
)

type SubscriptionService struct {
	repository        *repository.Repository
	config            *config.Config
	twitchAPI         *twitch.TwitchAPI
	chatService       Chat
	validationService Validate
}

func NewSubscrpitionService(repo *repository.Repository, api *twitch.TwitchAPI, chat Chat, validationService Validate, conf *config.Config) *SubscriptionService {

	return &SubscriptionService{repository: repo, twitchAPI: api, chatService: chat, validationService: validationService, config: conf}

}

func (service *SubscriptionService) SubscribeToAllStreamUps(chat *models.Chat) (int, int, error) {

	token, err := service.chatService.GetChatUserAccessToken(chat)
	if err != nil {
		return 0, 0, fmt.Errorf("subscribeToAllStreamUps err: %w", err)
	}

	follows, err := service.twitchAPI.GetAccountFollows(chat.TwitchID, token)
	if err != nil {
		return 0, 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
	}

	for _, f := range follows {
		if err := service.repository.AddFollow(chat.ID, models.Follow{ID: f.BroadcasterID, BroadcasterName: f.BroadcasterName}); err != nil {
			return 0, 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
		}

		if err := service.repository.CreateStreamer(&models.StreamerAccount{ID: f.BroadcasterID}); err != nil {
			return 0, 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
		}

	}

	allFollows, err := service.repository.GetUnSubedFollows()
	if err != nil {
		return 0, 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
	}

	notPayedStreamers := 0
	subsError := 0
	var apiError error
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, MAX_EVENTGO)
	errChan := make(chan error)
	defer close(errChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, f := range allFollows {
		wg.Add(1)

		if service.config.PayModeOn {
			payed, err := service.validationService.IsSubscriptionActive(f.ID)
			if err != nil {
				wg.Done()
				return 0, 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
			}

			if !payed {
				notPayedStreamers++
				wg.Done()
				continue
			}

			go service.subscribeToTwitchEvents(ctx, cancel, sem, &wg, &mu, errChan, &subsError, &apiError, f)
			continue
		}

		go service.subscribeToTwitchEvents(ctx, cancel, sem, &wg, &mu, errChan, &subsError, &apiError, f)
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		return 0, 0, <-errChan
	default:
	}

	if subsError != 0 {
		return 0, subsError, fmt.Errorf("%d: subscriptions are not completed\nlast error: %w", subsError, apiError)
	}

	return notPayedStreamers, 0, nil
}

func (service *SubscriptionService) subscribeToTwitchEvents(ctx context.Context, cancel context.CancelFunc, sem chan struct{}, wg *sync.WaitGroup, mu *sync.Mutex, errChan chan<- (error), subsError *int, apiError *error, follow models.Follow) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	select {
	case <-ctx.Done():
		return
	default:
	}

	goctx, gocancel := context.WithTimeout(ctx, time.Second*10)
	defer gocancel()

	id, err := service.twitchAPI.SubscribeToTwitchEvent(goctx, follow.ID)

	if err != nil {
		mu.Lock()
		*subsError += 1
		*apiError = err
		mu.Unlock()
		return
	}

	if err := service.repository.UpdateSubID(follow.ID, id); err != nil {
		{
			errorContext, cancel := context.WithTimeout(ctx, time.Second*2)
			defer cancel()
			if err := service.twitchAPI.DeleteEventSub(errorContext, id); err != nil {
				log.Println("Failed to undo event sub, wrote it into network collection...")

				//Записываем в очередь сетевых вызовов которые необходимо выполнить

			}
		}

		select {
		case errChan <- fmt.Errorf("subscribeToAllStreamUps error: %w", err):
			cancel()
		default:
		}

	}
}

func (service *SubscriptionService) SubscribeStreamUPEvent(ctx context.Context, broadcasterID string) error {
	id, err := service.twitchAPI.SubscribeToTwitchEvent(ctx, broadcasterID)

	if err != nil {
		return err
	}

	if err := service.repository.UpdateSubID(broadcasterID, id); err != nil {
		errorContext, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		if err := service.twitchAPI.DeleteEventSub(errorContext, id); err != nil {
			log.Println("(TODO!!)Failed to undo event sub, wrote it into network collection...")

			//Записываем в очередь сетевых вызовов которые необходимо выполнить

			return err
		}

		return err
	}

	return nil
}

func (service *SubscriptionService) UnsubscribeStreamUPEvent(ctx context.Context, broadcasterID string) error {

	follow, err := service.repository.GetFollow(broadcasterID)
	if err != nil {
		return err
	}

	if err := service.repository.UpdateSubID(broadcasterID, ""); err != nil {
		return err
	}

	if err := service.twitchAPI.DeleteEventSub(ctx, follow.Subscribtion_id); err != nil {

		//Опять в очередь вызовов

		return err
	}

	return nil
}
