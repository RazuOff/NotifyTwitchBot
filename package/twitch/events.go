package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (api *twitchAPI) SubscribeToTwitchEvent(broadcasterID string) error {
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
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+api.OAuth.Access_token)
	req.Header.Set("Client-Id", api.clientId)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("SubscribeToTwitchEvent Response bodsy:", string(body))
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	if data["error"] != "" {
		return fmt.Errorf("%s message: %s", data["error"].(string), data["message"].(string))
	}

	log.Println("SubscribeToTwitchEvent Response Status:", resp.Status)
	return nil
}
