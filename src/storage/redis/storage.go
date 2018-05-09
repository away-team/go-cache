package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/away-team/go-cache/src/cache"
	rcache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

type storage struct {
	cache *rcache.Codec
}

// NewStorage returns a memory based cache.Storage
func NewStorage(ring *redis.Ring) cache.Storage {
	c := &rcache.Codec{
		Redis:     ring,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
	return &storage{cache: c}
}

// Get will fetch the cached value
func (s *storage) Get(key string, value interface{}) error {
	err := s.cache.Get(key, value)
	if err != nil {
		return fmt.Errorf("Error cache key (%s) not found: %v", key, err)
	}

	return nil
}

// Set will set a value in the cache for the expiration duration
func (s *storage) Set(key string, value interface{}, expiration time.Duration) error {
	return s.cache.Set(rcache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiration,
	})
}
