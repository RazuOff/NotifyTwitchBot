package twitch

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/internal/config"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type TwitchAPI struct {
	clientId  string
	appToken  string
	serverURL string
	OAuth     twitchmodels.OAuthResponse
	mutex     *sync.RWMutex
}

func NewTiwtchAPI(config *config.Config) *TwitchAPI {
	twitchAPI := &TwitchAPI{clientId: config.ClientID, appToken: config.AppToken, serverURL: config.ServerURL, mutex: &sync.RWMutex{}}
	go twitchAPI.UpdateOAuthToken()
	return twitchAPI
}

func (api *TwitchAPI) UpdateOAuthToken() {
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
