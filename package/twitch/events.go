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

	twitchmodels "github.com/RazuOff/NotifyTwitchBot/package/twitch/models"
)

func (api *twitchAPI) SubscribeToTwitchEvent(ctx context.Context, broadcasterID string) (string, error) {
	client := &http.Client{}
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	payload := map[string]interface{}{
		"type":    "stream.online",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": broadcasterID,
		},
		"transport": map[string]string{
			"method":   "webhook",
			"callback": api.serverURL + "/notify",
			"secret":   "s3cRe7asas",
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Client-Id", api.clientId)
	req.Header.Set("Content-Type", "application/json")

	api.mutex.RLock()
	req.Header.Set("Authorization", "Bearer "+api.OAuth.Access_token)
	resp, err := client.Do(req)
	api.mutex.RUnlock()

	if err != nil {
		return "", fmt.Errorf("SubscribeToTwitchEvent error %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("SubscribeToTwitchEvent error %w", err)
	}

	if resp.StatusCode != 202 {
		var errorDesc map[string]string
		json.Unmarshal(body, &errorDesc)
		return "", fmt.Errorf("SubscribeToTwitchEvent error %s %s", resp.Status, errorDesc["message"])
	}

	log.Println("SubscribeToTwitchEvent Response bodsy:", string(body))

	var data twitchmodels.WebhookData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("SubscribeToTwitchEvent error %w", err)
	}

	return data.Data[0].ID, nil
}

func (api *twitchAPI) DeleteEventSub(ctx context.Context, eventID string) error {
	apiURL := "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequestWithContext(ctx, "DELETE", apiURL, nil)
	if err != nil {
		return fmt.Errorf("DeleteEventSub error %w", err)
	}

	api.mutex.RLock()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.OAuth.Access_token))
	req.Header.Set("Client-Id", api.clientId)

	query := url.Values{}
	query.Set("id", eventID)
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)

	api.mutex.RUnlock()

	if err != nil {
		return fmt.Errorf("DeleteEventSub error %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		var errorDesc map[string]string
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &errorDesc)
		return fmt.Errorf("DeleteEventSub error %s %s", resp.Status, errorDesc["message"])
	}

	log.Printf("EventSub id= %s has been deleted", eventID)
	return nil
}

func (api *twitchAPI) GetAllSubs(ctx context.Context) ([]string, error) {
	apiURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)

	api.mutex.RLock()

	req.Header.Set("Authorization", "Bearer "+api.OAuth.Access_token)
	req.Header.Set("Client-Id", api.clientId)
	client := http.Client{}
	resp, err := client.Do(req)

	api.mutex.RUnlock()

	if err != nil {
		return nil, fmt.Errorf("GetAllEvents error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Data []struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("GetAllEvents error: %w", err)
	}

	var subIDs []string
	for _, sub := range data.Data {
		if sub.Type == "stream.online" {
			subIDs = append(subIDs, sub.Id)
		}
	}

	return subIDs, nil
}
