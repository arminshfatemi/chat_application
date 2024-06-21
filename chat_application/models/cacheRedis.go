package models

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// connecting to the redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to Redis")
	return client

}
