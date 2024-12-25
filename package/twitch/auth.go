package twitch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/package/models"
	"github.com/google/uuid"
)

func (api *twitchAPI) CreateAuthLink(chatID int64) string {

	apiURL := "https://id.twitch.tv/oauth2/authorize"
	query := url.Values{}
	query.Set("client_id", api.clientId)
	redirectURI := api.serverURL + "/auth"
	query.Set("redirect_uri", redirectURI)
	query.Set("response_type", "code")
	query.Set("scope", "user:read:follows")
	query.Set("state", generateState(chatID))
	query.Set("force_verify", "true")

	parsedURL, _ := url.Parse(apiURL)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}

func generateState(chatID int64) string {
	uuid := uuid.New().String()
	repository.SetUUID(chatID, uuid)
	return uuid
}

func (api *twitchAPI) GetUserAccessToken(code string) (models.UserAccessTokens, error) {
	client := &http.Client{}
	apiURL := "https://id.twitch.tv/oauth2/token"

	reqBody := url.Values{}
	reqBody.Set("client_id", api.clientId)
	reqBody.Set("client_secret", api.appToken)
	reqBody.Set("code", code)
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("redirect_uri", api.serverURL+"/auth")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(reqBody.Encode()))
	if err != nil {
		log.Print("GetAuthTokens - request err=" + err.Error())
		return models.UserAccessTokens{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("GetAuthTokens - responce err=" + err.Error())
		return models.UserAccessTokens{}, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("GetAuthTokens - ReadAll err=" + err.Error())
		return models.UserAccessTokens{}, err
	}

	var userAccessTokens models.UserAccessTokens
	if err = json.Unmarshal(rawData, &userAccessTokens); err != nil {
		log.Print("GetAuthTokens - Unmarshal err=" + err.Error())
		return models.UserAccessTokens{}, err
	}

	return userAccessTokens, nil
}
