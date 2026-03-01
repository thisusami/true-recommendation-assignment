package repositories

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/thisusami/true-recommendation-assignment/src/util"
)

type Caching struct {
	Client *redis.Client
}

func (c *Caching) Get(key string) (string, error) {
	util.Request("Request Redis cache", key)
	result, err := c.Client.Get(key).Result()
	if err != nil {
		util.Error("Redis Get Error", err.Error())
		return "", err
	}
	util.Response("Response Redis cache", result, 0)
	return result, nil
}
func (c *Caching) Set(key string, value interface{}, expiration time.Duration) error {
	util.Request("Set Redis cache", key)
	err := c.Client.Set(key, value, expiration).Err()
	if err != nil {
		util.Error("Redis Set Error", err.Error())
		return err
	}
	util.Response("Set Redis cache", key, 0)
	return nil
}
func NewCaching(client *redis.Client) *Caching {
	return &Caching{
		Client: client,
	}
}
