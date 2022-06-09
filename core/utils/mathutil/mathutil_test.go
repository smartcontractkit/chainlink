package mathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	// Happy path
	assert.Equal(t, 3, Max(3, 2, 1))
	// Single element
	assert.Equal(t, 3, Max(3))
	// Signed
	assert.Equal(t, -1, Max(-2, -1))
	// Uint64
	assert.Equal(t, uint64(2), Max(uint64(0), uint64(2)))
	// String
	assert.Equal(t, "c", Max("a", []string{"b", "c"}...))
}

func TestMin(t *testing.T) {
	// Happy path
	assert.Equal(t, 1, Min(3, 2, 1))
	// Single element
	assert.Equal(t, 3, Min(3))
	// Signed
	assert.Equal(t, -2, Min(-2, -1))
	// Uint64
	assert.Equal(t, uint64(0), Min(uint64(0), uint64(2)))
	// String
	assert.Equal(t, "a", Min("a", []string{"b", "c"}...))
}
