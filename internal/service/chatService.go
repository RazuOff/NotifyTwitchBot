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

		newToken, err := service.twitchAPI.RefreshUserAccessToken(ctx, token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("RefreshUserAccessToken err: %w", err)
		}

		if err := service.repository.SetToken(chat, newToken); err != nil {
			return nil, fmt.Errorf("SetToken err: %w", err)
		}
		chat.UserAccessToken = &newToken

		return &newToken, nil
	}

	return token, nil
}

func (service *ChatService) SetUserAccessToken(code string, chat *models.Chat) error {
	userAccessToken, err := service.twitchAPI.GetUserAccessToken(code)
	if err != nil {
		return fmt.Errorf("SetUserAccessToken error: %w", err)
	}

	if err := service.repository.SetToken(chat, userAccessToken); err != nil {
		return fmt.Errorf("SetUserAccessToken error: %w", err)
	}

	return nil
}

func (service *ChatService) SetTwitchID(chat *models.Chat) error {

	token, err := service.GetChatUserAccessToken(chat)
	if err != nil {
		return fmt.Errorf("SetTwitchID err: %w", err)
	}

	claims, err := service.twitchAPI.GetAccountClaims(token)
	if err != nil {
		return fmt.Errorf("SetTwitchID error: %w", err)
	}

	if err = service.repository.SetTwitchID(chat, claims.Sub); err != nil {
		return fmt.Errorf("SetTwitchID error: %w", err)
	}
	return nil
}
