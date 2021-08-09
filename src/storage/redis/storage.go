package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/away-team/go-cache/src/cache"
	rcache "github.com/go-redis/cache/v8"
	redis "github.com/go-redis/redis/v8"
)

type storage struct {
	cache *rcache.Cache
}

// NewStorage returns a memory based cache.Storage
func NewStorage(ring *redis.Ring) cache.Storage {
	c := rcache.New(&rcache.Options{
		Redis:     ring,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	})
	return &storage{cache: c}
}

// Get will fetch the cached value
func (s *storage) Get(key string, value interface{}) error {
	err := s.cache.Get(context.Background(), key, value)
	if err != nil {
		return fmt.Errorf("Error cache key (%s) not found: %v", key, err)
	}

	return nil
}

// Set will set a value in the cache for the expiration duration
func (s *storage) Set(key string, value interface{}, expiration time.Duration) error {
	return s.cache.Set(&rcache.Item{
		Key:   key,
		Value: value,
		TTL:   expiration,
	})
}

// Delete will remove a value from the cache
func (s *storage) Delete(key string) error {
	s.cache.Delete(context.Background(), key)
	return nil
}
