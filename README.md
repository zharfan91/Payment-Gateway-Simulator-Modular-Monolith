# Payment Gateway Simulator

A production-style modular monolith payment gateway simulator built with Go, Fiber, GORM, PostgreSQL, Redis, RabbitMQ, Swagger, Prometheus, OpenTelemetry, Jaeger, and Docker Compose.

## Quick Start

```bash
cp .env.example .env
docker compose up --build
```

API:

- App: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- RabbitMQ: http://localhost:15672
- Jaeger: http://localhost:16686

## Documentation

- [Architecture](docs/architecture.md)
- [ERD](docs/erd.md)
- [API Guide](docs/api.md)
- [Local Development](docs/local-development.md)
- [Production Deployment](docs/production-deployment.md)
