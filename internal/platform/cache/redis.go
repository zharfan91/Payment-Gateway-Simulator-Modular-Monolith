package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	return client, client.Ping(context.Background()).Err()
}
