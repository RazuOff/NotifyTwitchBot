package twitch

import (
	"log"
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
