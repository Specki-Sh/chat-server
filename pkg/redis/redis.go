package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

// client is the global Redis client that will be used throughout the application.
var client *redis.Client

// Config contains the settings for connecting to Redis.
type Config struct {
	Addr     string // The address of the Redis server
	Password string // The password to use when connecting to the Redis server
	DB       int    // The Redis database number
}

// initRedis creates a new Redis client using the settings from Config and returns it.
func initRedis(cfg Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client
}

// StartRedisConnection creates a new connection to Redis and stores it in the global client variable.
func StartRedisConnection(cfg Config) {
	client = initRedis(cfg)
}

// GetRedisConn returns the global Redis client for use throughout the application.
func GetRedisConn() *redis.Client {
	return client
}

// CloseRedisConnection closes the connection to Redis and returns an error if one occurred.
func CloseRedisConnection() error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("error occurred on redis connection closing: %s", err.Error())
	}
	return nil
}
