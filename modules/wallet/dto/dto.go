package dto

type WalletResponse struct {
	MerchantID string `json:"merchant_id"`
	Balance    int64  `json:"balance"`
}
