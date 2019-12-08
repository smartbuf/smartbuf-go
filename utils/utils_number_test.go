package utils

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestIntToUint(t *testing.T) {
	var arr = []int64{math.MaxInt64, 0, 1, -1, math.MinInt64}

	for _, i := range arr {
		i2 := IntToUint(i)
		assert.Truef(t, i == UintToInt(i2), "%v, %v", i, i2)
	}

	for i := 0; i < 100000; i++ {
		num := rand.Int63()
		x := IntToUint(num)
		val := UintToInt(x)

		assert.True(t, num == val)
	}
}

// BenchmarkIntToUint-12    	1000000000	         0.536 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIntToUint(b *testing.B) {
	num := rand.Int63()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := IntToUint(num)
		UintToInt(x)
	}
}
