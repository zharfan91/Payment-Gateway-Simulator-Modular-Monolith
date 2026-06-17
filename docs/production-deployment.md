# Production Deployment Guide

- Build immutable container images for `cmd/api` and `cmd/worker`.
- Run Goose migrations before starting new application tasks.
- Store JWT secrets, database credentials, Redis credentials, RabbitMQ credentials, and webhook secrets in a secret manager.
- Use managed PostgreSQL with backups and point-in-time recovery.
- Use managed Redis for JWT blacklist, refresh tokens, idempotency keys, rate limiting, and payment cache.
- Use durable RabbitMQ queues and configure dead-letter queues for notification retries.
- Export OpenTelemetry traces to the platform observability stack.
- Terminate TLS at a load balancer or ingress and pass request IDs through `X-Request-ID`.
