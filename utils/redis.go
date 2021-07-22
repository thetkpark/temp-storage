package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func getRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func SetTokenFileData(ctx context.Context, token string, fileData FileMetadata) error {
	rdb := getRedisClient()
	jsonString, err := json.Marshal(fileData)
	if err != nil {
		return fmt.Errorf("json.marshall: %v", err)
	}
	err = rdb.SetEX(ctx, token, jsonString, time.Hour*72).Err()
	return err
}

func GetFileDataFromToken(ctx context.Context, token string) (FileMetadata, error) {
	rdb := getRedisClient()
	var fileData FileMetadata
	data, err := rdb.Get(ctx, token).Result()
	if err != nil {
		return fileData, fmt.Errorf("rdb.Get: %v", err)
	}
	err = json.Unmarshal([]byte(data), &fileData)
	if err != nil {
		return fileData, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return fileData, nil
}

type FileMetadata struct {
	FileName   string `json:"fileName"`
	Key        string `json:"key"`
	ObjectName string `json:"objectName"`
}
