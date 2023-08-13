package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type redisTokenRepository struct {
	redis *redis.Client
}

func NewTokenCacheRepository(redisClient *redis.Client) use_case.TokenStorage {
	return &redisTokenRepository{
		redis: redisClient,
	}
}

func (r *redisTokenRepository) SetInvalidRefreshToken(
	ctx context.Context,
	userID entity.ID,
	refreshToken string,
	expiresIn time.Duration,
) error {
	key := fmt.Sprintf("%s:%s", refreshToken, "invalid")
	if err := r.redis.Set(ctx, key, userID, expiresIn).Err(); err != nil {
		return fmt.Errorf("redisTokenRepository.SetInvalidRefreshToken: %w",
			fmt.Errorf("could not SET refresh token to redis for userID: %d: %w", userID, err))
	}
	return nil
}

func (r *redisTokenRepository) InvalidRefreshTokenExists(
	ctx context.Context,
	refreshToken string,
) (bool, error) {
	key := fmt.Sprintf("%s:%s", refreshToken, "invalid")
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf(
			"redisTokenRepository.SetInvalidRefreshToken: %w",
			fmt.Errorf(
				"could not check if refresh token exists in redis for refreshToken: %s: %w",
				refreshToken,
				err,
			),
		)
	}
	return result == 1, nil
}
