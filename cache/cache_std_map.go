package cache

import "sync"

// A cache implementation using standard go map
// a sync.RWMutex has to be used to deal with data race
type stdMapCache struct {
	sync.RWMutex      // type embedding
	defaultTtlSeconds uint
	cacheMap          map[interface{}]cacheValue
}

func NewStdMapCache(defaultTtlSeconds uint) *stdMapCache {
	return &stdMapCache{defaultTtlSeconds: defaultTtlSeconds,
		cacheMap: make(map[interface{}]cacheValue),
	}
}

func (c *stdMapCache) Get(key interface{}) interface{} {
	defer c.RUnlock()
	c.RLock()
	v, ok := c.cacheMap[key]
	if !ok || v.IsTimeout() {
		return nil
	}
	return v.value
}

func (c *stdMapCache) Set(key interface{}, value interface{}) {
	defer c.Unlock()
	c.Lock()
	c.cacheMap[key] = newCacheValue(value, c.defaultTtlSeconds)
}

func (c *stdMapCache) SetWithTTl(key interface{}, value interface{}, ttlSeconds uint) {
	defer c.Unlock()
	c.Lock()
	c.cacheMap[key] = newCacheValue(value, ttlSeconds)
}

func (c *stdMapCache) Del(key interface{}) {
	defer c.Unlock()
	c.Lock()
	delete(c.cacheMap, key)
}

func (c *stdMapCache) Size() int {
	return len(c.cacheMap)
}

func (c *stdMapCache) Range(f func(key, value interface{}) bool) {
	defer c.RUnlock()
	c.RLock()
	for k, v := range c.cacheMap {
		c.RUnlock()
		f(k, v)
		c.RLock()
	}
}

func (c *stdMapCache) Name() string {
	return "std_map_cache"
}
