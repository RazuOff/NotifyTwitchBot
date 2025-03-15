package handlers

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	notifyService       service.Notify
	redirectService     service.Redirect
	chatService         service.Chat
	viewService         service.View
	subscriptionService service.Subscription
	isPayedMode         bool
}

func NewHandler(services *service.Service, isPayMode bool) *Handler {
	return &Handler{
		notifyService:       services.Notify,
		redirectService:     services.Redirect,
		chatService:         services.Chat,
		viewService:         services.View,
		subscriptionService: services.Subscription,
		isPayedMode:         isPayMode,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/notify", h.HandleNotifyWebhook)
	router.GET("/auth", h.HandleAuthRedirect)

	return router
}
