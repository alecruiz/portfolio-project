// pkg/database/redis.go
package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client connection
func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Application will continue without Redis caching")
		return nil
	}

	log.Println("Successfully connected to Redis")
	return client
}
