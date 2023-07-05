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

// SetRefreshToken stores a refresh token with an expiry time
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID entity.ID, tokenID entity.ID, expiresIn time.Duration) error {
	// We'll store userID with token id, so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		return fmt.Errorf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
	}
	return nil
}

// DeleteRefreshToken used to delete old  refresh tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID entity.ID, tokenID entity.ID) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		return fmt.Errorf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		return fmt.Errorf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
	}

	return nil
}

// DeleteUserRefreshTokens looks for all tokens beginning with
// userID and scans to delete them in a non-blocking fashion
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID entity.ID) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.redis.Del(ctx, iter.Val()).Err(); err != nil {
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		failCount++
	}

	if failCount > 0 {
		return fmt.Errorf("Failed to delete %d refresh tokens\n", failCount)
	}

	return nil
}
