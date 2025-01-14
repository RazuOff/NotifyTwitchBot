package twitch

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

func (api *twitchAPI) GetStreamInfo(id string) (twitchmodels.StreamInfo, error) {
	client := &http.Client{}
	apiUrl := "https://api.twitch.tv/helix/channels"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return twitchmodels.StreamInfo{}, err
	}

	q := req.URL.Query()
	q.Add("broadcaster_id", id)
	req.URL.RawQuery = q.Encode()

	api.mutex.Lock()
	defer api.mutex.Unlock()

	req.Header.Set("Authorization", "Bearer "+api.OAuth.Access_token)
	req.Header.Set("Client-Id", api.clientId)

	tryNumber := 3
	for i := 0; i < tryNumber; i++ {
		resp, err := client.Do(req)
		if err != nil {
			return twitchmodels.StreamInfo{}, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return twitchmodels.StreamInfo{}, err
		}

		log.Println("GetStreamInfo Response body:", string(body))

		var result struct {
			Data []twitchmodels.StreamInfo `json:"data"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return twitchmodels.StreamInfo{}, err
		}

		if len(result.Data) != 0 {
			data := result.Data[0]
			return data, nil
		} else {
			log.Default().Print("Dont get stream info. Trying enother time")
			resp.Body.Close()
			time.Sleep(1 * time.Second)
		}

	}

	return twitchmodels.StreamInfo{}, errors.New("stream info is nil")
}
