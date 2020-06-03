package utils

import "time"

// Set an item to the cache
func (u *Utils) Set(k string, x interface{}) {
	u.cache.SetDefault(k, x)
}

// SetWithTTL an item to the cache with TTL
func (u *Utils) SetWithTTL(k string, x interface{}, ttlSeconds int) {
	ttl := time.Duration(ttlSeconds) * time.Second
	u.cache.Set(k, x, ttl)
}

// Get returns an object containing Found (bool) and Value
func (u *Utils) Get(k string) struct {
	Found bool
	Value interface{}
} {
	v, found := u.cache.Get(k)

	return struct {
		Found bool
		Value interface{}
	}{
		Found: found,
		Value: v,
	}
}

// Delete an item to the cache
func (u *Utils) Delete(k string) {
	u.cache.Delete(k)
}
