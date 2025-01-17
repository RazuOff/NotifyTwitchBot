package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"github.com/gin-gonic/gin"
)

type RedirectService struct {
	repository *repository.Repository
	bot        View
}

func NewRedirectService(repo *repository.Repository, bot View) *RedirectService {
	return &RedirectService{repository: repo, bot: bot}
}

// FIX THIS
func (service *RedirectService) HandleAuthError(chatID int64, text string, c *gin.Context, err error) {
	log.Println(err.Error())
	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
	if err := service.repository.DeleteChat(chatID); err != nil {
		log.Printf("HandleAuthRedirect error: %s", err.Error())
	}
	service.bot.SendMessage(chatID, text)
}

func (service *RedirectService) GetChatFromRedirect(state string) (*models.Chat, error) {

	chat, err := service.repository.GetChatByUUID(state)
	if err != nil {
		return nil, fmt.Errorf("GetChatFromRedirect error: %w", err)
	}

	if chat.ID == 0 {
		return nil, fmt.Errorf("GetChatFromRedirect error: UUID not valid ")
	}

	if err := service.repository.DeleteUUID(chat); err != nil {
		return chat, fmt.Errorf("GetChatFromRedirect error: %w", err)
	}

	return chat, nil
}

func (service *RedirectService) SetUserAccessToken(code string, chat *models.Chat) error {
	userAccessToken, err := twitch.TwitchAPI.GetUserAccessToken(code)
	if err != nil {
		return fmt.Errorf("SetUserAccessToken error: %w", err)
	}

	if err := service.repository.SetToken(chat, userAccessToken); err != nil {
		return fmt.Errorf("SetUserAccessToken error: %w", err)
	}

	return nil
}

func (service *RedirectService) SetTwitchID(chat *models.Chat) error {

	claims, err := twitch.TwitchAPI.GetAccountClaims(chat.UserAccessToken)
	if err != nil {
		return fmt.Errorf("SetTwitchID error: %w", err)
	}

	if err = service.repository.SetTwitchID(chat, claims.Sub); err != nil {
		return fmt.Errorf("SetTwitchID error: %w", err)
	}
	return nil
}

func (service *RedirectService) GetChatUserAccessToken(chat *models.Chat) (twitchmodels.UserAccessToken, error) {
	token, isExpired := chat.GetUserAccessToken()

	if isExpired {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		token, err := twitch.TwitchAPI.RefreshUserAccessToken(ctx, token.RefreshToken)
		cancel()
		if err != nil {
			return twitchmodels.UserAccessToken{}, fmt.Errorf("subscribeToAllStreamUps err: %w", err)
		}

		if err := service.repository.SetToken(chat, token); err != nil {
			return twitchmodels.UserAccessToken{}, fmt.Errorf("subscribeToAllStreamUps err: %w", err)
		}
	}

	return token, nil
}

func (service *RedirectService) SubscribeToAllStreamUps(chat *models.Chat) (int, error) {

	token, err := service.GetChatUserAccessToken(chat)
	if err != nil {
		return 0, fmt.Errorf("subscribeToAllStreamUps err: %w", err)
	}

	follows, err := twitch.TwitchAPI.GetAccountFollows(chat.TwitchID, token)
	if err != nil {
		return 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
	}

	for _, f := range follows {
		if err := service.repository.AddFollow(chat.ID, models.Follow{ID: f.BroadcasterID, BroadcasterName: f.BroadcasterName}); err != nil {
			return 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
		}
	}

	allFollows, err := service.repository.GetUnSubedFollows()
	if err != nil {
		return 0, fmt.Errorf("subscribeToAllStreamUps error: %w", err)
	}

	subsError := 0
	var apiError error
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	defer close(errChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, f := range allFollows {
		wg.Add(1)
		go func(follow models.Follow) {

			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			goctx, gocancel := context.WithTimeout(ctx, time.Second*10)
			defer gocancel()

			id, err := twitch.TwitchAPI.SubscribeToTwitchEvent(goctx, follow.ID)

			if err != nil {
				mu.Lock()
				subsError++
				apiError = err
				mu.Unlock()
				return
			}

			follow.Subscribtion_id = id
			if err := service.repository.SaveFollow(&follow); err != nil {
				select {
				case errChan <- fmt.Errorf("subscribeToAllStreamUps error: %w", err):
					cancel()
				default:
				}

			}

		}(f)
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		return 0, <-errChan
	default:
	}

	if subsError != 0 {
		return subsError, fmt.Errorf("%d: subscriptions are not completed\nlast error: %w", subsError, apiError)
	}

	return 0, nil
}
