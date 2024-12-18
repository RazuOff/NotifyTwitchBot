package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Проверяем, является ли это challenge-запросом
	if challenge, ok := data["challenge"]; ok {
		c.String(http.StatusOK, challenge.(string))
		return
	}

	// Обрабатываем событие (например, начало стрима)
	fmt.Println("Received webhook event:", string(body))

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
