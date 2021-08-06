package cache

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/groupcache/singleflight"
)

// Logger can log using Printf
type Logger interface {
	Printf(format string, v ...interface{})
}

// Cache can GetAndLoad, Get, Set, and Delete items by key
type Cache interface {
	GetAndLoad(key string, value interface{}, loader func() (interface{}, error)) error
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration *time.Duration) error
	Delete(key string) error
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
func (c *cache) GetAndLoad(key string, value interface{}, loader func() (interface{}, error)) (err error) {
	// if its in the cache return it
	err = c.backend.Get(key, value)
	if err == nil {
		return nil
	}

	// otherwise fetch & load using singleflight
	result, err := c.group.Do(key, func() (interface{}, error) {
		//fetch the value using the provided loader
		v, err := loader()
		if err != nil {
			return nil, fmt.Errorf("Error loading value for key (%s): %v", key, err)
		}

		// set the value in the backend
		err = c.backend.Set(key, v, c.expiration)
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
		return v, nil
	})
	if err != nil {
		return err
	}

	SetValue(value, result)
	return nil
}

func (c *cache) Get(key string, value interface{}) error {
	return c.backend.Get(key, value)
}

func (c *cache) Set(key string, value interface{}, expiration *time.Duration) error {
	exp := c.expiration
	if expiration != nil {
		exp = *expiration
	}
	return c.backend.Set(key, value, exp)
}

func (c *cache) Delete(key string) error {
	return c.backend.Delete(key)
}

// SetValue will write result to value
func SetValue(value interface{}, result interface{}) {
	// use reflection to set value
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("Invalid value passed. Must be a non-nil pointer")
	}
	resultVal := reflect.ValueOf(result)
	rvElem := rv.Elem()
	if resultVal.Kind() != rvElem.Kind() && rvElem.Kind() != reflect.Interface {
		panic("Invalid value passed. Must be a pointer to loader's return type")
	}
	rvElem.Set(reflect.ValueOf(result))
}
