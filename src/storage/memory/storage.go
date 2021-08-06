package memory

import (
	"fmt"
	"time"

	"github.com/away-team/go-cache/src/cache"
	pcache "github.com/patrickmn/go-cache"
)

type storage struct {
	cache *pcache.Cache
}

// NewStorage returns a memory based cache.Storage
func NewStorage(defaultExpiration time.Duration, cleanupInterval time.Duration) cache.Storage {
	return &storage{cache: pcache.New(defaultExpiration, cleanupInterval)}
}

// Get will fetch the cached value
func (s *storage) Get(key string, value interface{}) error {
	result, ok := s.cache.Get(key)
	if !ok {
		return fmt.Errorf("Error cache key (%s) not found", key)
	}

	cache.SetValue(value, result)
	return nil
}

// Set will set a value in the cache for the expiration duration
func (s *storage) Set(key string, value interface{}, expiration time.Duration) error {
	s.cache.Set(key, value, expiration)
	return nil
}

// Delete will remove a value from the cache
func (s *storage) Delete(key string) error {
	s.cache.Delete(key)
	return nil
}
