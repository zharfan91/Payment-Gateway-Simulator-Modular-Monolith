package dto

type CreatePaymentRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	MerchantID   string `json:"merchant_id"`
	ExternalID   string `json:"external_id"`
	Amount       int64  `json:"amount"`
	Fee          int64  `json:"fee"`
	NetAmount    int64  `json:"net_amount"`
	Status       string `json:"status"`
	PaymentToken string `json:"payment_token"`
	PaymentURL   string `json:"payment_url"`
}
