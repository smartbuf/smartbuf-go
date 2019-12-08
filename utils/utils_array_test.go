package utils

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestDescFastSort(t *testing.T) {
	for x := 0; x < 100; x++ {
		arr := make([]int, x*int(rand.Uint32()%20)+1)
		for i := 0; i < len(arr); i++ {
			arr[i] = rand.Int()
		}

		t.Log(len(arr))

		DescFastSort(arr[:])

		for i := 0; i < len(arr)-1; i++ {
			assert.True(t, arr[i] >= arr[i+1])
		}
	}
}
