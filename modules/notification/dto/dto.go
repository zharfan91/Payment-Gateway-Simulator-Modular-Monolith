package dto

type RegisterWebhookRequest struct {
	WebhookURL string `json:"webhook_url"`
}

type TestWebhookRequest struct {
	Event string `json:"event"`
}
