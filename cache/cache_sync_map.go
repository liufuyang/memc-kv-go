package cache

import "sync"

// A cache implementation using sync.Map
type syncMapCache struct {
	defaultTtlSeconds uint
	cacheMap          sync.Map
}

func NewSyncMapCache(defaultTtlSeconds uint) *syncMapCache {
	return &syncMapCache{defaultTtlSeconds: defaultTtlSeconds}
}

func (c *syncMapCache) Get(key interface{}) interface{} {
	vRaw, ok := c.cacheMap.Load(key)
	if v, isCV := vRaw.(cacheValue); !ok || !isCV || v.IsTimeout() {
		return nil
	}
	v, _ := vRaw.(cacheValue)
	return v.value
}

func (c *syncMapCache) Set(key interface{}, value interface{}) {
	c.cacheMap.Store(key, newCacheValue(value, c.defaultTtlSeconds))
}

func (c *syncMapCache) SetWithTTl(key interface{}, value interface{}, ttlSeconds uint) {
	c.cacheMap.Store(key, newCacheValue(value, ttlSeconds))
}

func (c *syncMapCache) Del(key interface{}) {
	c.cacheMap.Delete(key)
}

func (c *syncMapCache) Size() int {
	// TODO - use an async counter to get size estimation instead of using heavy range
	size := 0
	c.cacheMap.Range(func(key, value interface{}) bool {
		size = size + 1
		return true
	})
	return size
}

func (c *syncMapCache) Range(f func(key, value interface{}) bool) {
	c.cacheMap.Range(f)
}

func (c *syncMapCache) Name() string {
	return "sync_map_cache"
}
