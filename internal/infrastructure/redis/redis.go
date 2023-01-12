package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

func NewConnection(redisURL string) (*redis.Client, error) {

	url, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(url)
	err = testConnection(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func testConnection(client *redis.Client) error {
	ret, err := client.Ping().Result()
	if ret != "PONG" {
		return errors.New("connection test failed. Check connection to Redis")
	}
	return err
}
