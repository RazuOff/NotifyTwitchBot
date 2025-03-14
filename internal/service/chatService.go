package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type ChatService struct {
	twitchAPI  *twitch.TwitchAPI
	repository repository.Chats
}

func NewChatService(api *twitch.TwitchAPI, repo *repository.Repository) *ChatService {
	return &ChatService{twitchAPI: api, repository: repo.Chats}
}

func (service *ChatService) GetChatUserAccessToken(chat *models.Chat) (*twitchmodels.UserAccessToken, error) {
	if chat.UserAccessToken == nil {
		return nil, fmt.Errorf("UserAccessToken is nil")
	}

	token, isExpired := chat.GetUserAccessToken()

	if isExpired {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		token, err := service.twitchAPI.RefreshUserAccessToken(ctx, token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("RefreshUserAccessToken err: %w", err)
		}

		if err := service.repository.SetToken(chat, token); err != nil {
			return nil, fmt.Errorf("SetToken err: %w", err)
		}

		time.Sleep(time.Second * 1)
	}

	return token, nil
}
