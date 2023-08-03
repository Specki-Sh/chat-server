package repository

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisTokenRepository struct {
	redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) use_case.TokenStorage {
	return &redisTokenRepository{
		redis: redisClient,
	}
}

func (r *redisTokenRepository) SetInvalidRefreshToken(ctx context.Context, userID entity.ID, refreshToken string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", refreshToken, "invalid")
	if err := r.redis.Set(ctx, key, userID, expiresIn).Err(); err != nil {
		return fmt.Errorf("Could not SET refresh token to redis for userID: %d: %v\n", userID, err)
	}
	return nil
}

func (r *redisTokenRepository) InvalidRefreshTokenExists(ctx context.Context, refreshToken string) (bool, error) {
	key := fmt.Sprintf("%s:%s", refreshToken, "invalid")
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("Could not check if refresh token exists in redis for refreshToken: %s: %v\n", refreshToken, err)
	}
	return result == 1, nil
}
