package cache

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

// NewCache create a new redis client
func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Error("Failed to connect to Redis: %v", err)
	} else {
		log.Info("Successfully connected to Redis!")
	}

	log.Info("Connected to redis.")
	return rdb
}
