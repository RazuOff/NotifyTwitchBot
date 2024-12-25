package app

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	router.POST("/notify", handlers.HandleNotifyWebhook)
	router.GET("/auth", handlers.HandleAuthRedirect)
	go func() { router.Run() }()

	twitch.Init()

	twitch.TwitchAPI.SubscribeToTwitchEvent()
}
