package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/internal/telegram"

	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
	"github.com/gin-gonic/gin"
)

func HandleNotifyWebhook(c *gin.Context) {
	// Считываем тело запроса
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Парсим JSON
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

	streamInfo, err := twitch.TwitchAPI.GetStreamInfo(event.BroadcasterUserID)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"status": err})
		return
	}

	log.Println("Received stream info:", streamInfo)

	chats, err := repository.GetChatsByFollow(streamInfo.BroadcasterID)
	if err != nil {
		log.Println("HandleNotifyWebhook error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, chat := range chats {
		telegram.SendMessage(chat.ID, streamInfo.BroadcasterName+" START STREAM!!")
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
