package ratelimiter

import (
	"net/http"
	"testing"

	cp "github.com/justinbarias/goratelimiter/cacheprovider"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type MockedCacheInvalidator struct {
	mock.Mock
	cp.ICacheInvalidator
}

type MockHttpHandler struct {
	handler http.Handler
}

func (mh MockHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (ci MockedCacheInvalidator) Invalidate(hash string, m map[string]int) bool {
	args := ci.Called(hash, m)
	return args.Bool(0)
}

func (ci MockedCacheInvalidator) GetRemainingSeconds(hash string) float64 {
	args := ci.Called(hash)
	return float64(args.Int(0))
}

type MockedCacheProvider struct {
	mock.Mock
	cp.CacheProvider
}

func (p MockedCacheProvider) Get(hash string) (value int, err error) {
	args := p.Called(hash)
	return args.Int(0), args.Error(1)
}

func (p MockedCacheProvider) Put(hash string) {
	p.Called(hash)
}

func TestRateLimiter_LimitIfExceedsThreshold(t *testing.T) {
	//arrange
	timeInSeconds := 2.0
	requestThreshold := 2
	requestHash := "some-hash"
	mockedInvalidator := MockedCacheInvalidator{}
	mockedInvalidator.On("Invalidate", requestHash, mock.Anything).Return(false)
	mockedInvalidator.On("GetRemainingSeconds", requestHash).Return(1.0)
	mockedCacheProvider := MockedCacheProvider{}
	mockedCacheProvider.On("Put", requestHash)

	sut := RateLimiter{
		handler:          MockHttpHandler{},
		maxCount:         requestThreshold,
		timeoutInSeconds: timeInSeconds,
		cache:            &mockedCacheProvider,
		invalidator:      &mockedInvalidator,
	}

	testCases := map[string]struct {
		providerResult int
		result         bool
	}{
		"LessThanThreshold":    {providerResult: requestThreshold + 1, result: true},
		"GreaterThanThreshold": {providerResult: requestThreshold + 1, result: false},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mockedCacheProvider.On("Get", requestHash).Return(testCase.providerResult, nil)
			//act
			result, _ := sut.Limit(requestHash)
			//assert
			assert.Equal(t, true, result)
		})
	}
}
