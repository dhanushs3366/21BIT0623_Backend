package redisservice

import (
	"context"
	"encoding/json"
	"log"
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
func (r *RedisClient) Add(key string, value []models.FileMetaData) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(context.Background(), key, jsonValue, EXPIRATION_TIME).Err()
	if err != nil {
		return err
	}
	log.Printf("Added %s to cache\n", key)
	return nil
}

func (r *RedisClient) Get(Key string) ([]models.FileMetaData, error) {
	jsonValue, err := r.client.Get(context.Background(), Key).Result()

	if err != nil {
		return nil, err
	}
	var metadata []models.FileMetaData
	err = json.Unmarshal([]byte(jsonValue), &metadata)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
