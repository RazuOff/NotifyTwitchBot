package twitchmodels

import "time"

type WebhookData struct {
	Data         []Webhook `json:"data"`
	Total        int       `json:"total"`
	TotalCost    int       `json:"total_cost"`
	MaxTotalCost int       `json:"max_total_cost"`
}

type Webhook struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	CreatedAt time.Time `json:"created_at"`
	Transport Transport `json:"transport"`
	Cost      int       `json:"cost"`
}

type Condition struct {
	UserID string `json:"user_id"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}
