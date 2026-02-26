package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7" 
)

func initRedisProperty() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,

	
		PoolSize:     20,            
		MinIdleConns: 5,                
		IdleTimeout:  5 * time.Minute, 
		PoolTimeout:  30 * time.Second, 
	})

	// Check connection
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
