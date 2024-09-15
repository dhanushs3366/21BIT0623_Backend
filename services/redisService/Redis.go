package redisservice

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var EXPIRATION_TIME time.Duration = 5 * time.Minute

func GetNewRedisClient() (*RedisClient, error) {
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	client := RedisClient{client: redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST,
		Password: REDIS_PASSWORD,
		DB:       0,
	})}

	err := client.client.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}
	return &client, nil
}

// cache metadata
// key-> user:userID:file:fileID
// value -> metadata obj
func (r *RedisClient) Add(key string, value *models.FileMetaData) error {
	// Marshal the FileMetaData struct to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(context.Background(), key, jsonValue, EXPIRATION_TIME).Err()
	return err
}
