package handlers

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/notify", h.HandleNotifyWebhook)
	router.GET("/auth", h.HandleAuthRedirect)

	return router
}
