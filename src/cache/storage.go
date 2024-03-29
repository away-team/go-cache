package cache

import "time"

// Storage provides Get, Set, and Delete to a cache backend
type Storage interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}
