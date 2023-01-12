package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type cacheStorage struct {
	client *redis.Client
}

func NewCacheStorage(db *redis.Client) *cacheStorage {
	return &cacheStorage{
		client: db,
	}
}

func (cs *cacheStorage) SetKey(key, value string, expiration int) error {
	key = strings.ReplaceAll(key, ":", "")
	return cs.client.Set(key, value, time.Duration(int64(expiration))*time.Minute).Err()
}

func (cs *cacheStorage) GetKey(key string) (string, error) {
	key = strings.ReplaceAll(key, ":", "")
	return cs.client.Get(key).Result()
}

func (cs *cacheStorage) DelKey(key string) (int64, error) {
	key = strings.ReplaceAll(key, ":", "")
	return cs.client.Del(key).Result()
}
