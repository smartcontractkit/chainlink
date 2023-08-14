package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSnoozedRoots(t *testing.T) {
	c := NewSnoozedRootsInMem()

	k1 := [32]byte{1}
	k2 := [32]byte{2}
	t1 := time.Now().Add(-24 * time.Hour)
	t2 := time.Now()

	// element should not exist
	_, ok := c.Get(k1)
	assert.False(t, ok)

	// after an element is set it should exist
	c.Set(k1, t1)
	v, ok := c.Get(k1)
	assert.True(t, ok)
	assert.Equal(t, t1, v)

	// other key should not exist
	_, ok = c.Get(k2)
	assert.False(t, ok)

	// after setting some other key it should exist
	c.Set(k2, t2)
	v, ok = c.Get(k2)
	assert.True(t, ok)
	assert.Equal(t, t2, v)

	// and other elements should not be affected
	v, ok = c.Get(k1)
	assert.True(t, ok)
	assert.Equal(t, t1, v)
}
