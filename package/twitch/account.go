package twitch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

func (api *twitchAPI) GetAccountFollows(accountID string) ([]twitchmodels.FollowInfo, error) {
	apiURL := "https://api.twitch.tv/helix/channels/followed"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return []twitchmodels.FollowInfo{}, err
	}
	queries := url.Values{}
	queries.Set("user_id", accountID)
	queries.Set("first", "100")

	req.URL.RawQuery = queries.Encode()

	currentChat, err := repository.GetChatByTwitchID(accountID)
	if err != nil {
		return []twitchmodels.FollowInfo{}, err
	}

	req.Header.Set("Authorization", "Bearer "+currentChat.UserAccessToken.AccessToken)
	req.Header.Set("Client-Id", api.clientId)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []twitchmodels.FollowInfo{}, err
	}
	defer resp.Body.Close()

	rawData, _ := io.ReadAll(resp.Body)
	log.Println("GetAccountFollows get data:" + string(rawData))
	var data map[string]json.RawMessage

	if err = json.Unmarshal(rawData, &data); err != nil {
		return []twitchmodels.FollowInfo{}, err
	}

	var followInfo []twitchmodels.FollowInfo
	if err = json.Unmarshal([]byte(data["data"]), &followInfo); err != nil {
		return []twitchmodels.FollowInfo{}, err
	}

	return followInfo, nil

}

func (api *twitchAPI) GetAccountClaims(token *twitchmodels.UserAccessTokens) (twitchmodels.DeafultAccountClaims, error) {
	apiURL := "https://id.twitch.tv/oauth2/userinfo"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return twitchmodels.DeafultAccountClaims{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return twitchmodels.DeafultAccountClaims{}, err
	}
	defer resp.Body.Close()

	rawClaims, _ := io.ReadAll(resp.Body)
	log.Print("GetAccountClaims get data:" + string(rawClaims))
	var claims twitchmodels.DeafultAccountClaims

	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return twitchmodels.DeafultAccountClaims{}, err
	}

	return claims, nil
}
