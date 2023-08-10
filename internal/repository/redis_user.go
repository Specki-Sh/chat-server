package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"chat-server/internal/domain/entity"
)

type UserCacheRepository struct {
	client *redis.Client
}

func (r *UserCacheRepository) SetUserData(ctx context.Context, secretCode string, userData *entity.UserData) error {
	key := secretCode
	data, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("UserCacheRepository.SetUserData: %w", err)
	}
	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("UserCacheRepository.SetUserData: %w", err)
	}
	return nil
}

func (r *UserCacheRepository) GetUserData(ctx context.Context, secretCode string) (*entity.UserData, error) {
	key := secretCode
	data, err := r.client.Get(ctx, key).Bytes()
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
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("UserCacheRepository.DeleteUserData: %w", err)
	}
	return nil
}
