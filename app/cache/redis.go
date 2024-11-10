package cache

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// NewCache creates a new redis client
func NewCache() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load cache environment file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_ADR"),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       0, // use default DB
	})
	return rdb
}
