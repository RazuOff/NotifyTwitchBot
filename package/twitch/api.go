package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/package/models"
)

type twitchAPI struct {
	clientId  string
	appToken  string
	serverURL string
	OAuth     models.OAuthResponse
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

func (api *twitchAPI) getOAuthToken() (models.OAuthResponse, error) {
	client := &http.Client{}
	apiUrl := "https://id.twitch.tv/oauth2/token"

	payload := url.Values{}
	payload.Set("client_id", api.clientId)
	payload.Set("client_secret", api.appToken)
	payload.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(payload.Encode())))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return models.OAuthResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.OAuthResponse{}, err
	}

	var auth models.OAuthResponse
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return models.OAuthResponse{}, err
	}

	log.Println("GetOAuthToken Response Status:", resp.Status)
	return auth, nil
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
