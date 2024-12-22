package cache

import (
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

// NewCache create a new redis client
func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redzone-redis:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0, // use default DB
	})

	log.Debug("Connected to redis.")
	return rdb
}
