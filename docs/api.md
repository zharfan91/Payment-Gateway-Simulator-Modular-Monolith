# API Guide

All API routes use `/api/v1`.

1. Register and log in with `/auth/register` or `/auth/login`.
2. Create a merchant with `POST /merchants`.
3. Generate API credentials with `POST /merchants/api-keys` and `POST /merchants/secret-keys`.
4. Create payments using `POST /payments` with:
   - `X-API-Key`
   - optional `X-Signature` HMAC SHA256 over the raw JSON body
   - optional `Idempotency-Key`
5. Simulate status changes using the dashboard JWT routes:
   - `POST /payments/:id/simulate-success`
   - `POST /payments/:id/simulate-failed`
   - `POST /payments/:id/refund`
6. Register webhook URL with `POST /webhooks`.
7. Withdraw available balance with `POST /settlements`.

Money values are integer minor units. For example, `500000` with the fixed `2.9%` fee produces fee `14500` and net `485500`.
