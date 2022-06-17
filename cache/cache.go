package cache

import "time"

type Cache interface {
	Get(key interface{}) interface{}
	Set(key interface{}, value interface{})
	SetWithTTl(key interface{}, value interface{}, ttlSeconds uint)
	Del(key interface{})
	Size() int
	Name() string
	Range(f func(key, value interface{}) bool)
}

type cacheValue struct {
	value     interface{}
	timestamp time.Time
}

func newCacheValue(value interface{}, ttlSeconds uint) cacheValue {
	return cacheValue{
		value:     value,
		timestamp: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
}

func (cv *cacheValue) IsTimeout() bool {
	return cv.timestamp.Before(time.Now())
}
