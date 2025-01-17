package twitch

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type twitchAPI struct {
	clientId  string
	appToken  string
	serverURL string
	OAuth     twitchmodels.OAuthResponse
	mutex     *sync.RWMutex
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
	twitchAPI := &twitchAPI{clientId: clientId, appToken: appToken, serverURL: server_url, mutex: &sync.RWMutex{}}
	TwitchAPI = twitchAPI

	go TwitchAPI.updateOAuthToken()
}

func (api *twitchAPI) updateOAuthToken() {
	for {

		if api.OAuth.ExpiresIn < 1000 {
			api.mutex.Lock()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			OAuth, err := api.getOAuthToken(ctx)
			if err != nil {
				log.Println(err)

				cancel()
				api.mutex.Unlock()
				time.Sleep(1 * time.Second)
				continue
			}

			api.OAuth = OAuth

			api.mutex.Unlock()
			cancel()
		} else {
			api.OAuth.ExpiresIn--
			time.Sleep(1 * time.Second)
		}
	}
}
