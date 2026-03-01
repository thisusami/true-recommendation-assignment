package adapters

import (
	"time"

	"github.com/go-redis/redis/v7"
)

func InitRedisProperty(connectionString string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: password,
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
