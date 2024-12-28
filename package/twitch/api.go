package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type twitchAPI struct {
	clientId  string
	appToken  string
	serverURL string
	OAuth     twitchmodels.OAuthResponse
}

var TwitchAPI *twitchAPI

func Init() {
	clientId, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("CLIENT_ID env parametr not found!")
	}
	appToken, exists := os.LookupEnv("APP_TOKEN")
	if !exists {
		log.Fatal("APP_TOKEN env parametr not found!")
	}
	server_url, exists := os.LookupEnv("SERVER_URL")
	if !exists {
		log.Fatal("SERVER_URL env parametr not found!")
	}
	twitchAPI := &twitchAPI{clientId: clientId, appToken: appToken, serverURL: server_url}

	OAuth, err := twitchAPI.getOAuthToken()
	if err != nil {
		log.Fatal(err)
	}
	twitchAPI.OAuth = OAuth
	TwitchAPI = twitchAPI
}

func (api *twitchAPI) SubscribeToTwitchEvent() {
	client := &http.Client{}
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	// Тело запроса
	payload := map[string]interface{}{
		"type":    "stream.online",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": "87791915", // Замените на ID стримера
		},
		"transport": map[string]string{
			"method":   "webhook",
			"callback": api.serverURL + "/notify",
			"secret":   "s3cRe7asas",
		},
	}

	// Преобразуем тело в JSON
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+api.OAuth.Access_token) // Замените на токен
	req.Header.Set("Client-Id", api.clientId)                         // Замените на ваш Client ID
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println("SubscribeToTwitchEvent Response bodsy:", string(body))
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	log.Println("SubscribeToTwitchEvent Response Status:", resp.Status)
}
