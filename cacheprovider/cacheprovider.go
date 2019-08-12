package cacheprovider

import (
	"errors"
)

type CacheProvider interface {
	Put(hash string)
	Get(hash string) (value int, err error)
}

type InMemoryCache struct {
	CacheProvider
	m             map[string]int
	timeInSeconds float64
	invalidator   ICacheInvalidator
}

func (p InMemoryCache) Get(hash string) (value int, err error) {
	if p.m[hash] == 0 {
		return 0, errors.New("Hash does not exist")
	}
	p.invalidator.Invalidate(hash, p.m)

	return p.m[hash], nil
}

func (p InMemoryCache) Put(hash string) {
	val, err := p.Get(hash)
	if err != nil {
		val = 0
	}
	if !p.invalidator.Invalidate(hash, p.m) {
		p.m[hash] = val + 1
	}
}

func NewInMemoryCache(timeInSeconds float64) *InMemoryCache {
	var cache InMemoryCache
	cache.m = make(map[string]int)
	cache.invalidator = GetCacheInvalidatorSingleton(timeInSeconds, 10)
	cache.timeInSeconds = timeInSeconds
	return &cache
}
