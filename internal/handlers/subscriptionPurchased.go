package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleSubscriptionPurchased(c *gin.Context) {

	type body struct {
		ID              string `gorm:"unique;not null" json:"id"`
		SubDaysDuration *int   `json:"sub_dayDuration"`
	}

	var streamer body

	if err := c.BindJSON(&streamer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := h.service.BuyStreamerSub(streamer.ID, streamer.SubDaysDuration); err != nil {
		log.Printf("HandleSubscriptionPurchased error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
