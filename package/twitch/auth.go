package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"

	"github.com/google/uuid"
)

func (api *twitchAPI) getOAuthToken() (twitchmodels.OAuthResponse, error) {
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
		return twitchmodels.OAuthResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return twitchmodels.OAuthResponse{}, err
	}

	var auth twitchmodels.OAuthResponse
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return twitchmodels.OAuthResponse{}, err
	}

	log.Println("GetOAuthToken Response Status:", resp.Status)
	return auth, nil
}

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

func (api *twitchAPI) GetUserAccessToken(code string) (twitchmodels.UserAccessTokens, error) {
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
		return twitchmodels.UserAccessTokens{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("GetAuthTokens - responce err=" + err.Error())
		return twitchmodels.UserAccessTokens{}, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("GetAuthTokens - ReadAll err=" + err.Error())
		return twitchmodels.UserAccessTokens{}, err
	}

	var userAccessTokens twitchmodels.UserAccessTokens
	if err = json.Unmarshal(rawData, &userAccessTokens); err != nil {
		log.Print("GetAuthTokens - Unmarshal err=" + err.Error())
		return twitchmodels.UserAccessTokens{}, err
	}

	return userAccessTokens, nil
}
