package rates

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var redisInitialized bool

func InitRedis() {
	if redisInitialized {
		log.Println("Redis is already initialized. Skipping reinitialization.")
		return
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	options := &redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
	}

	log.Printf("Connecting to Redis at %s...\n", redisAddress)

	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	RedisClient = client
	redisInitialized = true
}
