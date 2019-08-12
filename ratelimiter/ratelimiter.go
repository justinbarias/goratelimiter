package ratelimiter

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"

	"github.com/justinbarias/goratelimiter/cacheprovider"
)

type RateLimiter struct {
	handler          http.Handler
	maxCount         int
	timeoutInSeconds float64
	cache            cacheprovider.CacheProvider
	invalidator      cacheprovider.ICacheInvalidator
}

func computeHash(r *http.Request) string {
	encoded := b64.StdEncoding.EncodeToString([]byte(r.RemoteAddr + r.RequestURI))

	fmt.Printf("Remote address %s, Request URI %s", r.RemoteAddr, r.RequestURI)
	return encoded
}

// Limit - returns True if request is throttled, False otherwise
func (rl *RateLimiter) Limit(requestHash string) (throttle bool, err error) {
	rl.cache.Put(requestHash)
	result, err := rl.cache.Get(requestHash)
	if err != nil {
		// log error
		return
	}
	if result <= rl.maxCount {
		return false, nil
	}
	return true, nil

}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestHash := computeHash(r)
	result, err := rl.Limit(requestHash)
	if err != nil {
		// log error
		return
	}
	if !result {
		rl.handler.ServeHTTP(w, r)
	} else {
		http.Error(w, fmt.Sprintf("Too many requests. Try again in %f seconds", rl.invalidator.GetRemainingSeconds(requestHash)), http.StatusTooManyRequests)
		return
	}
}

// NewRateLimiter - creates new instance of RateLimiter
func NewRateLimiter(handlerToWrap http.Handler, maxCount int, timeoutInSeconds float64) *RateLimiter {
	return &RateLimiter{
		handler:     handlerToWrap,
		maxCount:    maxCount,
		invalidator: cacheprovider.GetCacheInvalidatorSingleton(timeoutInSeconds, maxCount),
		cache:       cacheprovider.NewInMemoryCache(timeoutInSeconds),
	}
}
