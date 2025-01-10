package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/internal/telegram"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	"github.com/gin-gonic/gin"
)

func HandleAuthRedirect(c *gin.Context) {

	if err := c.Query("error"); err != "" {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	state := c.Query("state")
	if state == "" {
		log.Println("state does not exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "state does not exists"})
		return
	}
	log.Printf("Get status = %s", state)

	chat, err := repository.GetChatByUUID(state)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusForbidden, err)
		return
	}

	code := c.Query("code")
	userAccessToken, err := twitch.TwitchAPI.GetUserAccessToken(code)
	if err != nil {
		handleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	if err := repository.SetToken(chat, userAccessToken); err != nil {
		handleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	if !handleSettingTwitchID(chat, c) {
		return
	}

	if err := subscribeToAllStreamUps(chat, c); err != nil {
		return
	}

	telegram.SendMessage(chat.ID, "Вы успешно вошли в аккаунт!\nТеперь вам будут прихожить уведомления о начале стримов")

	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
}

func subscribeToAllStreamUps(chat *models.Chat, c *gin.Context) error {
	follows, err := twitch.TwitchAPI.GetAccountFollows(chat)
	if err != nil {
		log.Printf("subscribeToAllStreamUps error: %s", err.Error())
		return err
	}

	for _, f := range follows {
		if err := repository.AddFollow(chat.ID, models.Follow{ID: f.BroadcasterID, BroadcasterName: f.BroadcasterName}); err != nil {
			log.Printf("subscribeToAllStreamUps error: %s", err.Error())
			return err
		}
	}

	allFollows, err := repository.GetUnSubedFollows()
	if err != nil {
		log.Printf("subscribeToAllStreamUps error: %s", err.Error())
		return err
	}

	subsError := 0
	var apiError error
	for _, f := range allFollows {
		if err := twitch.TwitchAPI.SubscribeToTwitchEvent(f.ID); err != nil {
			apiError = err
			subsError++
		} else {
			f.IsSubscribed = true
			if err := repository.SaveFollow(&f); err != nil {
				log.Printf("subscribeToAllStreamUps error: %s", err.Error())
				return err
			}
		}
	}

	if subsError != 0 {
		telegram.SendMessage(chat.ID, fmt.Sprintf("%d фоллоу не получилось подписать на оповещения(", subsError))
		c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
		return fmt.Errorf("%d: subscriptions are not completed\nlast error: %w", subsError, apiError)
	}

	return nil
}

func handleSettingTwitchID(chat *models.Chat, c *gin.Context) bool {

	claims, err := twitch.TwitchAPI.GetAccountClaims(chat.UserAccessToken)
	if err != nil {
		handleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return false
	}

	if err = repository.SetTwitchID(chat, claims.Sub); err != nil {
		handleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return false
	}
	return true
}

func handleAuthError(chatID int64, text string, c *gin.Context, err error) {
	log.Println(err.Error())
	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
	if err := repository.DeleteChat(chatID); err != nil {

		log.Printf("HandleAuthRedirect error: %s", err.Error())
	}
	telegram.SendMessage(chatID, text)
}
