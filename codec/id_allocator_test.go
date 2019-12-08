package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIDAllocator(t *testing.T) {
	var id IDAllocator

	assert.True(t, id.Require() == 0)
	assert.True(t, id.Require() == 1)
	assert.True(t, id.Require() == 2)
	assert.True(t, id.Require() == 3)
	assert.True(t, id.Require() == 4)

	id.Release(1)
	id.Release(3)
	id.Release(2)

	assert.True(t, id.Require() == 1)
	assert.True(t, id.Require() == 2)
	assert.True(t, id.Require() == 3)
	assert.True(t, id.Require() == 5)
}
