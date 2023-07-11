package simcache

import (
	"sync"
	"time"
)

// Cache holds any items of type T that are cleared after a given TTL.
// The cache clears any expired items upon any retrieval operation.
type Cache[T any] struct {
	*cache[T]
}

// New creates an empty Cache where the TTL for item's added will be set to the given duration.
func New[T any](defaultTTL time.Duration) *Cache[T] {
	items := make(map[string]item[T])
	return &Cache[T]{cache: &cache[T]{
		items:      items,
		defaultTTL: defaultTTL,
		mutex:      &sync.RWMutex{},
	}}
}

// Add inserts the item T into the cache for a given key if no item has been already added with the same key.
// It returns false if the item was not added due to an existing item with the same key being there.
// It returns true if the item was added successfully.
func (c *cache[T]) Add(key string, value T, ttl ...time.Duration) bool {
	expiration := calculateExpiration(c.defaultTTL, ttl...)
	c.mutex.RLock()
	_, found := c.items[key]
	if found {
		c.mutex.RUnlock()
		return false
	}

	c.mutex.RUnlock()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = item[T]{
		value:      value,
		expiration: expiration,
	}
	return true
}

// Set replaces the value in the cache for a given key. If no such key exists, it adds it to the cache.
// If no duration, or a value of 0, is specified it uses the default TTL when the cache was made.
// Only the first duration given is used when multiple are passed in.
func (c *cache[T]) Set(key string, value T, ttl ...time.Duration) {
	expiration := calculateExpiration(c.defaultTTL, ttl...)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = item[T]{
		value:      value,
		expiration: expiration,
	}
}

// Get returns the value in the cache for a given key and if it was found. If no such key exists, the returned bool will be false.
func (c *cache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	i, found := c.items[key]
	if !found {
		c.mutex.RUnlock()
		return i.value, false
	}

	if i.expired() {
		c.mutex.RUnlock()
		c.Delete(key)
		return i.value, false
	}
	c.mutex.RUnlock()
	return i.value, true
}

// Delete removes the item from the cache for the given key.
func (c *cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// Items returns a copy of the cache's map that holds type T.
func (c *cache[T]) Items() map[string]T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	items := make(map[string]T, len(c.items))
	for k, i := range c.items {
		if i.expired() {
			c.mutex.RUnlock()
			c.Delete(k)
			c.mutex.RLock()
			continue
		}
		items[k] = i.value
	}
	return items
}

// Keys returns a slice of the cache's keys.
func (c *cache[T]) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var keys []string
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of the cache's values of type T.
func (c *cache[T]) Values() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var values []T
	for k, i := range c.items {
		if i.expired() {
			c.mutex.RUnlock()
			c.Delete(k)
			c.mutex.RLock()
			continue
		}
		values = append(values, i.value)
	}
	return values
}

// Purge removes all expired items from the cache.
func (c *Cache[T]) Purge() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	count := 0
	for k, i := range c.items {
		if i.expired() {
			c.mutex.RUnlock()
			c.Delete(k)
			c.mutex.RLock()
			count++
		}
	}
	return count
}

type item[T any] struct {
	value      T
	expiration time.Time
}

func (i *item[T]) expired() bool {
	return time.Now().UTC().After(i.expiration)
}

type cache[T any] struct {
	items      map[string]item[T]
	defaultTTL time.Duration
	mutex      *sync.RWMutex
}

func calculateExpiration(defaultTTL time.Duration, ttl ...time.Duration) time.Time {
	t := time.Now().Add(defaultTTL).UTC()
	givenValidTTL := len(ttl) > 0 && ttl[0] > 0
	if givenValidTTL {
		t = time.Now().Add(ttl[0]).UTC()
	}
	return t
}
