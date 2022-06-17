package cache

import (
	"time"
)

// Provide vacuum capability on Cache
// For example:
//   var c Cache = NewStdMapCache(defaultTtlSeconds)
//   c = NewVacuumedCache(c)
type vacuumedCache struct {
	Cache
	vacuumCycleSleepMs   time.Duration
	vacuumKeyLoopSleepNs time.Duration
}

func NewVacuumedStdMapCache(defaultTtlSeconds uint) Cache {
	var c Cache = NewStdMapCache(defaultTtlSeconds)
	c = newVacuumedCache(c)
	return c
}

func NewVacuumedSyncMapCache(defaultTtlSeconds uint) Cache {
	var c Cache = NewSyncMapCache(defaultTtlSeconds)
	c = newVacuumedCache(c)
	return c
}

func newVacuumedCache(cache Cache) *vacuumedCache {
	v := &vacuumedCache{Cache: cache,
		vacuumCycleSleepMs:   time.Duration(1000) * time.Millisecond,
		vacuumKeyLoopSleepNs: time.Duration(1000) * time.Nanosecond,
	}
	v.startVacuum()
	return v
}

func (vc *vacuumedCache) startVacuum() {
	go func() {
		for {
			vc.Range(vc.vacuumFunc)
			time.Sleep(vc.vacuumCycleSleepMs) // sleep here for some time to reduce loop frequency
		}
	}()
}

func (vc *vacuumedCache) vacuumFunc(key, value interface{}) bool {
	cValue, _ := value.(cacheValue)
	if cValue.IsTimeout() {
		vc.Del(key)
	}
	time.Sleep(vc.vacuumKeyLoopSleepNs) // sleep here for some time to reduce loop frequency
	return true
}

func (vc *vacuumedCache) Name() string {
	return "vacuumed_" + vc.Cache.Name() // using vc.Cache.Name() instead vc.Name() to avoid recursion
}
