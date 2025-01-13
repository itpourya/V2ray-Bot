package cache

import (
	"context"
	"github.com/itpourya/Haze/config"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

// NewCache create a new redis client
func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.CACHE_ADR,
		Password: config.CACHE_PASSWORD,
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

	return rdb
}
