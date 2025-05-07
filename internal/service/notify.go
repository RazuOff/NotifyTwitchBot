package service

import (
	"fmt"
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/apperrors"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

type NotifyService struct {
	repository repository.Follows
	bot        View
	twitchAPI  *twitch.TwitchAPI
}

func NewNotifyService(repo repository.Follows, bot View, api *twitch.TwitchAPI) *NotifyService {
	return &NotifyService{repository: repo, bot: bot, twitchAPI: api}
}

func (service *NotifyService) SendNotify(broadcasterUserID string) *apperrors.AppError {
	streamInfo, err := service.twitchAPI.GetStreamInfo(broadcasterUserID)
	if err != nil {
		return apperrors.NewAppError("NOT FOUND", "GetStreamInfo error:", err)
	}

	log.Println("Received stream info:", streamInfo)

	chats, err := service.repository.GetChatsByFollow(streamInfo.BroadcasterID)
	if err != nil {
		return apperrors.NewAppError("DB ERROR", "GetStreamInfo error:", err)
	}

	for _, chat := range chats {
		link := fmt.Sprintf("https://www.twitch.tv/%s", streamInfo.BroadcasterLogin)
		service.bot.SendMessage(chat.ID, fmt.Sprintf("%s START STREAM!!\n\n%s\n\n%s", streamInfo.BroadcasterName, streamInfo.Title, link))
	}

	return nil
}


