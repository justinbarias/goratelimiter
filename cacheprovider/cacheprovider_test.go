package cacheprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

func TestCacheProvider_PutsEntryInCache(t *testing.T) {
	//Arrange
	hash := "some-hash"
	timeInSeconds := 2.0
	sut := NewInMemoryCache(timeInSeconds)

	//Act
	sut.Put(hash)
	result, _ := sut.Get(hash)

	//Assert
	assert.Equal(t, 1, result)

}

func TestCacheProvider_ReturnErrorWhenHashDoesntExist(t *testing.T) {
	//Arrange
	hash := "some-hash"
	timeInSeconds := 2.0
	sut := NewInMemoryCache(timeInSeconds)

	//Act
	_, err := sut.Get(hash)

	//Assert
	assert.NotNil(t, err)

}

type MockedCacheInvalidator struct {
	mock.Mock
	ICacheInvalidator
}

func (ci MockedCacheInvalidator) Invalidate(hash string, m map[string]int) bool {
	args := ci.Called(hash, m)
	return args.Bool(0)
}

func (ci MockedCacheInvalidator) GetRemainingSeconds(hash string) float64 {
	args := ci.Called(hash)
	return float64(args.Int(0))
}

func TestCacheProvider_GetCallsInvalidate(t *testing.T) {
	//Arrange
	hash := "some-hash"
	timeInSeconds := 2.0
	mockedInvalidator := MockedCacheInvalidator{}
	mockedInvalidator.On("Invalidate", hash, mock.Anything).Return(false)
	sut := InMemoryCache{
		m:             make(map[string]int),
		timeInSeconds: timeInSeconds,
		invalidator:   &mockedInvalidator,
	}

	//Act
	sut.Put(hash) //put so we can get something
	sut.Get(hash)

	//Assert
	mockedInvalidator.AssertExpectations(t)

}

func TestCacheProvider_PutCallsInvalidate(t *testing.T) {
	//Arrange
	hash := "some-hash"
	timeInSeconds := 2.0
	mockedInvalidator := MockedCacheInvalidator{}
	mockedInvalidator.On("Invalidate", hash, mock.Anything).Return(false)
	sut := InMemoryCache{
		m:             make(map[string]int),
		timeInSeconds: timeInSeconds,
		invalidator:   &mockedInvalidator,
	}

	//Act
	sut.Put(hash) //put so we can get something

	//Assert
	mockedInvalidator.AssertExpectations(t)

}
