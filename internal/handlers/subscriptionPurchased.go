package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleSubscriptionPurchased(c *gin.Context) {
	var streamer models.StreamerAccount

	c.BindJSON(streamer)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.service.SubscribeStreamUPEvent(ctx, streamer.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
