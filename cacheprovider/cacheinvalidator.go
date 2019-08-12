package cacheprovider

import (
	"sync"
	"time"
)

var instance *CacheInvalidator
var once sync.Once

type ICacheInvalidator interface {
	Invalidate(hash string, m map[string]int) bool
	GetRemainingSeconds(hash string) float64
}

type CacheInvalidator struct {
	ICacheInvalidator
	m             map[string]time.Time
	timeInSeconds float64
	maxCount      int
}

func (ci CacheInvalidator) Invalidate(hash string, m map[string]int) bool {
	if _, ok := ci.m[hash]; !ok {
		ci.m[hash] = time.Now()
		return false
	}
	t := ci.m[hash]
	dur := time.Now().Sub(t)
	if m[hash] <= ci.maxCount {
		// update timestamp only if not throttled
		ci.m[hash] = time.Now()
	}
	if dur.Seconds() > ci.timeInSeconds {
		m[hash] = 0
		return true
	}
	return false
}

func (ci CacheInvalidator) GetRemainingSeconds(hash string) float64 {
	t := ci.m[hash]
	dur := time.Now().Sub(t)
	return ci.timeInSeconds - dur.Seconds()
}

func NewCacheInvalidator(timeinSeconds float64, maxCount int) *CacheInvalidator {
	var invalidator CacheInvalidator
	invalidator.m = make(map[string]time.Time)
	invalidator.timeInSeconds = timeinSeconds
	invalidator.maxCount = maxCount
	return &invalidator
}

func GetCacheInvalidatorSingleton(timeinSeconds float64, maxCount int) *CacheInvalidator {
	singletonFunc := func() {
		var invalidator CacheInvalidator
		invalidator.m = make(map[string]time.Time)
		invalidator.timeInSeconds = timeinSeconds
		invalidator.maxCount = maxCount
		instance = &invalidator
	}
	once.Do(singletonFunc)
	return instance
}
