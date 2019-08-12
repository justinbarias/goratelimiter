package cacheprovider

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheInvalidator_DoestNotInvalidateOnFirst(t *testing.T) {
	//Arrange
	m := make(map[string]int)
	hash := "some-hash"
	sut := NewCacheInvalidator(10, 10)

	//Act
	result := sut.Invalidate(hash, m)

	//Assert
	assert.Equal(t, false, result)

}

func TestCacheInvalidator_InvalidateWhenReachedThreshold(t *testing.T) {
	//Arrange
	m := make(map[string]int)
	hash := "some-hash"
	timeInSeconds := 2.0 //set to a low value
	rateLimit := 2       //set to a low value
	sut := NewCacheInvalidator(timeInSeconds, rateLimit)
	sut.Invalidate(hash, m)
	sut.Invalidate(hash, m)
	time.Sleep(time.Second * time.Duration(int64(timeInSeconds)))
	//Act
	result := sut.Invalidate(hash, m)

	//Assert
	assert.Equal(t, true, result)

}

func TestCacheInvalidator_GetsCorrectRemainingTime(t *testing.T) {
	//Arrange
	m := make(map[string]int)
	hash := "some-hash"
	timeInSeconds := 2.0
	rateLimit := 2
	sut := NewCacheInvalidator(timeInSeconds, rateLimit)
	sut.Invalidate(hash, m)
	sut.Invalidate(hash, m)
	time.Sleep(time.Second * 1)
	//Act
	result := sut.GetRemainingSeconds(hash)
	t.Logf("Remaining time = %f", result)
	//Assert
	assert.True(t, result <= 1 && result > 0)
}
