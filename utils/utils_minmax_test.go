package utils

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestMinInt(t *testing.T) {
	assert.True(t, MinInt(math.MaxInt64) == math.MaxInt64)
	assert.True(t, MinInt(math.MaxInt64, 0) == 0)
	assert.True(t, MinInt(1, 2, 3, 10, -1, -100, -200, 300) == -200)

	arr := make([]int, 1000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Int()
	}

	min := MinInt(arr...)

	for _, v := range arr {
		assert.True(t, v >= min)
	}
}

func TestMaxInt(t *testing.T) {
	assert.True(t, MaxInt(math.MinInt64) == math.MinInt64)
	assert.True(t, MaxInt(math.MinInt64, 0) == 0)

	arr := make([]int, 1000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Int()
	}

	max := MaxInt(arr...)

	for _, v := range arr {
		assert.True(t, v <= max)
	}
}
