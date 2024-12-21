package cache

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// NewCache create a new redis client
func NewCache() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Can't load cache environment file", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_ADR"),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       0, // use default DB
	})

	log.Debug("Connected to redis.")
	return rdb
}
