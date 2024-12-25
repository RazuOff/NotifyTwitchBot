package handlers

import (
	"log"
	"net/http"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	"github.com/gin-gonic/gin"
)

func HandleAuthRedirect(c *gin.Context) {

	if err := c.Query("error"); err != "" {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
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
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	repository.SetToken(chat.ID, userAccessToken)

	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
}
