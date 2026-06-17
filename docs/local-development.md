# Local Development

Requirements:

- Go 1.24+
- Docker and Docker Compose
- Goose CLI for manual migrations

Start everything:

```bash
cp .env.example .env
docker compose up --build
```

Run locally against Docker dependencies:

```bash
docker compose up postgres redis rabbitmq jaeger
go run ./cmd/api
go run ./cmd/worker
```

Run tests:

```bash
go test ./...
go test -tags=integration ./...
```
