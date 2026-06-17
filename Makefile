.PHONY: run worker test test-integration tidy swagger migrate-up

run:
	go run ./cmd/api

worker:
	go run ./cmd/worker

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

tidy:
	go mod tidy

swagger:
	swag init -g cmd/api/main.go -o docs/swagger

migrate-up:
	goose -dir migrations postgres "$$DATABASE_DSN" up
