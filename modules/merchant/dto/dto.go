package dto

type CreateMerchantRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	WebhookURL string `json:"webhook_url"`
}

type UpdateWebhookRequest struct {
	WebhookURL string `json:"webhook_url"`
}

type MerchantResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	WebhookURL string `json:"webhook_url"`
}

type APIKeyResponse struct {
	APIKey string `json:"api_key"`
}

type SecretKeyResponse struct {
	SecretKey string `json:"secret_key"`
}
