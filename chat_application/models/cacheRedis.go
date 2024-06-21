package models

import (
	"context"
	"encoding/json"
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

// GetRecentMessagesCache will get the recent messages if the cache exists
func GetRecentMessagesCache(redisClient *redis.Client, roomName string) ([]MessageJson, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	value, err := redisClient.Get(ctx, "cache:recentMessages:"+roomName).Result()
	if err != nil {
		log.Println("CheckCacheExists: ", err)
		return []MessageJson{}, err
	}
	log.Println("cache:", value)

	var retrievedMessages []MessageJson
	err = json.Unmarshal([]byte(value), &retrievedMessages)
	if err != nil {
		log.Fatalf("Could not deserialize data: %v", err)
		return []MessageJson{}, err
	}

	log.Println("END: ", retrievedMessages)
	return retrievedMessages, nil
}

// CacheRecentMessages is going to cache given recent messages for the given room
func CacheRecentMessages(redisClient *redis.Client, roomName string, recentMessages *[]MessageJson) error {
	marshaled, err := json.Marshal(*recentMessages)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = redisClient.Set(ctx, "cache:recentMessages:"+roomName, marshaled, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
