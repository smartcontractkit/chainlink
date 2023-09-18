package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSnoozedRoots(t *testing.T) {
	c := NewSnoozedRoots(1*time.Minute, 1*time.Minute)

	k1 := [32]byte{1}
	k2 := [32]byte{2}

	// return false for non existing element
	snoozed := c.IsSnoozed(k1)
	assert.False(t, snoozed)

	// after an element is marked as executed it should be snoozed
	c.MarkAsExecuted(k1)
	snoozed = c.IsSnoozed(k1)
	assert.True(t, snoozed)

	// after snoozing an element it should be snoozed
	c.Snooze(k2)
	snoozed = c.IsSnoozed(k2)
	assert.True(t, snoozed)
}

func TestEvictingElements(t *testing.T) {
	c := newSnoozedRoots(1*time.Millisecond, 1*time.Hour, 1*time.Millisecond, 1*time.Millisecond)

	k1 := [32]byte{1}
	c.Snooze(k1)

	time.Sleep(10 * time.Millisecond)

	assert.False(t, c.IsSnoozed(k1))
}
