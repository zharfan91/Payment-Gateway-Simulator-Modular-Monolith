package dto

type AuditLogResponse struct {
	ID       string `json:"id"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}
