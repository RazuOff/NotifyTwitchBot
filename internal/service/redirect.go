package service

import (
	"fmt"
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
)

type RedirectService struct {
	repository *repository.Repository
	bot        View
}

func NewRedirectService(repo *repository.Repository, bot View, chatService Chat) *RedirectService {
	return &RedirectService{repository: repo, bot: bot}
}

func (service *RedirectService) HandleAuthError(chatID int64, text string, err error) {
	log.Println(err.Error())
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
