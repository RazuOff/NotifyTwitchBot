package twitch

import (
	"context"
	"log"
	"sync"
	"time"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

type TwitchAPI struct {
	clientId  string
	appToken  string
	serverURL string
	OAuth     twitchmodels.OAuthResponse
	mutex     *sync.RWMutex
}

func NewTiwtchAPI(clientID string, appToken string, serverURL string) *TwitchAPI {
	twitchAPI := &TwitchAPI{clientId: clientID, appToken: appToken, serverURL: serverURL, mutex: &sync.RWMutex{}}

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
