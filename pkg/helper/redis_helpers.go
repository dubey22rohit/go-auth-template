package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func StoreInRedis(redisClient *redis.Client, prefix string, hash string, userID uuid.UUID, expiration time.Duration) error {
	ctx := context.Background()
	err := redisClient.Set(
		ctx,
		fmt.Sprintf("%s%s", prefix, userID),
		hash,
		expiration,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetFromRedis(redisClient *redis.Client, key string) (*string, error) {
	ctx := context.Background()

	hash, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &hash, nil
}
