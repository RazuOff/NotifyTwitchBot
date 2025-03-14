package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

// Getting an app token
func (api *TwitchAPI) getOAuthToken(ctx context.Context) (twitchmodels.OAuthResponse, error) {
	client := &http.Client{}
	apiUrl := "https://id.twitch.tv/oauth2/token"

	payload := url.Values{}

	payload.Set("client_id", api.clientId)
	payload.Set("client_secret", api.appToken)
	payload.Set("grant_type", "client_credentials")

	req, _ := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer([]byte(payload.Encode())))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

func (api *TwitchAPI) CreateAuthLink(chatID int64, state string) string {

	apiURL := "https://id.twitch.tv/oauth2/authorize"
	query := url.Values{}
	query.Set("client_id", api.clientId)
	redirectURI := api.serverURL + "/auth"
	query.Set("redirect_uri", redirectURI)
	query.Set("response_type", "code")
	query.Set("scope", "user:read:follows")
	query.Set("state", state)
	query.Set("force_verify", "true")

	parsedURL, _ := url.Parse(apiURL)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}

func (api *TwitchAPI) RefreshUserAccessToken(ctx context.Context, refreshToken string) (twitchmodels.UserAccessToken, error) {
	apiURL := "https://id.twitch.tv/oauth2/token"
	reqBody := url.Values{}
	reqBody.Set("client_id", api.clientId)
	reqBody.Set("client_secret", api.appToken)
	reqBody.Set("grant_type", "refresh_token")
	reqBody.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(reqBody.Encode()))
	if err != nil {
		log.Print("RefreshUserAccessToken - request err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("RefreshUserAccessToken - responce err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("RefreshUserAccessToken - ReadAll err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}

	if resp.StatusCode >= 400 {
		var errorDesc map[string]string
		json.Unmarshal(rawData, &errorDesc)
		log.Print("RefreshUserAccessToken -err=" + string(rawData))
		return twitchmodels.UserAccessToken{}, fmt.Errorf(errorDesc["message"])
	}

	var userAccessTokens twitchmodels.UserAccessToken
	if err = json.Unmarshal(rawData, &userAccessTokens); err != nil {
		log.Print("RefreshUserAccessToken - Unmarshal err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}

	userAccessTokens.CreatedAt = time.Now()

	return userAccessTokens, nil
}

func (api *TwitchAPI) GetUserAccessToken(code string) (twitchmodels.UserAccessToken, error) {
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
		return twitchmodels.UserAccessToken{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("GetAuthTokens - responce err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("GetAuthTokens - ReadAll err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}

	var userAccessTokens twitchmodels.UserAccessToken
	if err = json.Unmarshal(rawData, &userAccessTokens); err != nil {
		log.Print("GetAuthTokens - Unmarshal err=" + err.Error())
		return twitchmodels.UserAccessToken{}, err
	}

	return userAccessTokens, nil
}
