package handlers

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *service.Service
	isPayedMode bool
}

func NewHandler(services *service.Service, isPayMode bool) *Handler {
	return &Handler{
		service:     services,
		isPayedMode: isPayMode,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/notify", h.HandleNotifyWebhook)
	router.GET("/auth", h.HandleAuthRedirect)
	router.POST("/subscription-purchased", h.HandleSubscriptionPurchased)

	return router
}
