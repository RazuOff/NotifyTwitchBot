package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleNotifyWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var data map[string]json.RawMessage
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var challange string
	err = json.Unmarshal([]byte(data["challenge"]), &challange)
	if err == nil {
		fmt.Println("challange sucsessfull take reqest " + challange)
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, challange)
		return
	}

	var event twitchmodels.Event
	err = json.Unmarshal([]byte(data["event"]), &event)
	if err != nil {
		fmt.Println("Error decoding Event:", err)
		return
	}
	log.Println("Received webhook event:", event)

	if err := h.services.Notify.SendNotify(event.BroadcasterUserID); err != nil {
		switch err.Code {
		case "NOT FOUND":
			log.Print(err)
			c.JSON(http.StatusNotFound, gin.H{"status": err})
			return
		case "DB ERROR":
			log.Printf("HandleNotifyWebhook error: %s" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
