package handlers

import (
	"log"
	"net/http"

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
	log.Println(repository.Chats)
	chat, err := repository.GetChatByUUID(state)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusForbidden, err)
		return
	}

	code := c.Query("code")
	userAccessToken, err := twitch.TwitchAPI.GetUserAccessToken(code)
	if err != nil {
		HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	if err := repository.SetToken(chat.ID, userAccessToken); err != nil {
		HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	telegram.SendMessage(chat.ID, "Вы успешно вошли в аккаунт!\nТеперь вам будут прихожить уведомления о начале стримов")

	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
}

func HandleAuthError(chatID int64, text string, c *gin.Context, err error) {
	log.Println(err)
	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
	telegram.SendMessage(chatID, text)
}
