package handlers

import (
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	notifyService   service.Notify
	redirectService service.Redirect
	viewService     service.View
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{notifyService: services.Notify, redirectService: services.Redirect, viewService: services.View}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/notify", h.HandleNotifyWebhook)
	router.GET("/auth", h.HandleAuthRedirect)

	return router
}
