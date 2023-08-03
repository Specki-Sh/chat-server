package repository

import (
	"context"
	"encoding/json"

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
		return err
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *UserCacheRepository) GetUserData(ctx context.Context, secretCode string) (*entity.UserData, error) {
	key := secretCode
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var userData entity.UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return nil, err
	}
	return &userData, nil
}

func (r *UserCacheRepository) DeleteUserData(ctx context.Context, secretCode string) error {
	key := secretCode
	return r.client.Del(ctx, key).Err()
}
