package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("redis ping err:", err)
	} else {
		log.Println("redis ping success")
	}
}
