package app

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	router.POST("/notify", handlers.HandleNotifyWebhook)
	router.GET("/auth", handlers.HandleAuthRedirect)

	go func() { router.Run() }()

	if err := postgre.Init(); err != nil {
		log.Print(err.Error())
	}
	twitch.Init()
}
