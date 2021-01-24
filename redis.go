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

func SetURLAndToken(token string, url string) error {
	var ctx = context.Background()
	rdb := getRedisClient()
	err := rdb.SetEX(ctx, token, url, time.Hour*72).Err()
	return err
}

func GetURLFromToken(token string) (string, error) {
	var ctx = context.Background()
	rdb := getRedisClient()
	url, err := rdb.Get(ctx, token).Result()
	if err != nil {
		return "", fmt.Errorf("rdb.Get: %v", err)
	}
	return url, nil
}
