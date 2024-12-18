package app

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	router.POST("/notify", handlers.HandleNotifyWebhook)
	defer router.Run()

	twitchAPI := twitch.Init()

	authToken, err := twitchAPI.GetOAuthToken()
	if err != nil {
		log.Fatal(err.Error())
	}

	twitchAPI.SubscribeToTwitchEvent(authToken.Access_token)
}
