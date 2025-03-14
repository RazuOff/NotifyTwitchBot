package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

func (api *TwitchAPI) GetAccountFollows(twitchID string, userAccessToken *twitchmodels.UserAccessToken) ([]twitchmodels.FollowInfo, error) {
	apiURL := "https://api.twitch.tv/helix/channels/followed"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return []twitchmodels.FollowInfo{}, fmt.Errorf("GetAccountFollows error: %w", err)
	}
	queries := url.Values{}
	queries.Set("user_id", twitchID)
	queries.Set("first", "100")

	req.URL.RawQuery = queries.Encode()
	req.Header.Set("Authorization", "Bearer "+userAccessToken.AccessToken)
	req.Header.Set("Client-Id", api.clientId)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []twitchmodels.FollowInfo{}, fmt.Errorf("GetAccountFollows error: %w", err)
	}
	defer resp.Body.Close()

	rawData, _ := io.ReadAll(resp.Body)
	log.Println("GetAccountFollows get data:" + string(rawData))
	var data map[string]json.RawMessage

	if err = json.Unmarshal(rawData, &data); err != nil {
		return []twitchmodels.FollowInfo{}, fmt.Errorf("GetAccountFollows error: %w", err)
	}

	var followInfo []twitchmodels.FollowInfo
	if err = json.Unmarshal([]byte(data["data"]), &followInfo); err != nil {
		return []twitchmodels.FollowInfo{}, fmt.Errorf("GetAccountFollows error: %w", err)
	}

	return followInfo, nil

}

func (api *TwitchAPI) GetAccountClaims(token *twitchmodels.UserAccessToken) (twitchmodels.DeafultAccountClaims, error) {
	apiURL := "https://id.twitch.tv/oauth2/userinfo"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Print("GetAccountClaims error")
		return twitchmodels.DeafultAccountClaims{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("GetAccountClaims error")
		return twitchmodels.DeafultAccountClaims{}, err
	}
	defer resp.Body.Close()

	rawClaims, _ := io.ReadAll(resp.Body)
	log.Print("GetAccountClaims get data:" + string(rawClaims))
	var claims twitchmodels.DeafultAccountClaims

	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		log.Print("GetAccountClaims error")
		return twitchmodels.DeafultAccountClaims{}, err
	}

	return claims, nil
}
