package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0, // use default DB
	})

	return rdb
}
