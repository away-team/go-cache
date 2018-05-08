package cache

import (
	"fmt"
	"time"

	"github.com/golang/groupcache/singleflight"
)

// Logger can log using Printf
type Logger interface {
	Printf(format string, v ...interface{})
}

// Cache can GetAndLoad items
type Cache interface {
	GetAndLoad(key string, loader func() (interface{}, error)) (interface{}, error)
}

type cache struct {
	backend         Storage
	expiration      time.Duration
	group           *singleflight.Group
	ignoreSetErrors bool
	log             Logger
}

// New returns a new *Cache
func New(backend Storage, expiration time.Duration, ignoreSetErrors bool, log Logger) Cache {
	return &cache{backend: backend, group: &singleflight.Group{}, expiration: expiration, log: log, ignoreSetErrors: ignoreSetErrors}
}

// GetAndLoad will fetch the value from the cache, if it is missing it will load it into cache and return it.
// Concurrent calls to GetAndLoad for the same key will wait for the initial call to load the object
// to prevent cache runs.
func (c *cache) GetAndLoad(key string, loader func() (interface{}, error)) (interface{}, error) {
	// if its in the cache return it
	result, err := c.backend.Get(key)
	if err == nil {
		return result, nil
	}

	// otherwise fetch & load using singleflight
	result, err = c.group.Do(key, func() (interface{}, error) {
		//fetch the value using the provided loader
		value, err := loader()
		if err != nil {
			return nil, fmt.Errorf("Error loading value for key (%s): %v", key, err)
		}

		// set the value in the backend
		err = c.backend.Set(key, value, c.expiration)
		if err != nil {
			// log the error and ignore it.
			if c.log != nil {
				c.log.Printf("Error setting value in cache for key (%s): %v", key, err)
			}
			// only return an error if we aren't configured to ignore them
			if !c.ignoreSetErrors {
				return nil, fmt.Errorf("Error setting value in cache for key (%s): %v", key, err)
			}
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
