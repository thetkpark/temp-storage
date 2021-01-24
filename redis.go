package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func getRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func SetURLAndToken(ctx context.Context, token string, url string) error {
	rdb := getRedisClient()
	err := rdb.SetEX(ctx, token, url, time.Hour*72).Err()
	return err
}

func GetURLFromToken(ctx context.Context, token string) (string, error) {
	rdb := getRedisClient()
	url, err := rdb.Get(ctx, token).Result()
	if err != nil {
		return "", fmt.Errorf("rdb.Get: %v", err)
	}
	return url, nil
}
