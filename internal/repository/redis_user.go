package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

func NewUserCacheRepository(redis *redis.Client) use_case.UserCacheStorage {
	return &UserCacheRepository{
		redis: redis,
	}
}

type UserCacheRepository struct {
	redis *redis.Client
}

func (r *UserCacheRepository) SetUserData(
	ctx context.Context,
	secretCode string,
	userData *entity.UserData,
) error {
	key := secretCode
	data, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("UserCacheRepository.SetUserData: %w", err)
	}
	if err := r.redis.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("UserCacheRepository.SetUserData: %w", err)
	}
	return nil
}

func (r *UserCacheRepository) GetUserData(
	ctx context.Context,
	secretCode string,
) (*entity.UserData, error) {
	key := secretCode
	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("UserCacheRepository.GetUserData: %w", err)
	}
	var userData entity.UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return nil, fmt.Errorf("UserCacheRepository.GetUserData: %w", err)
	}
	return &userData, nil
}

func (r *UserCacheRepository) DeleteUserData(ctx context.Context, secretCode string) error {
	key := secretCode
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("UserCacheRepository.DeleteUserData: %w", err)
	}
	return nil
}
