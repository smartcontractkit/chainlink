package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazyLoad(t *testing.T) {
	var clientwg sync.WaitGroup

	tc := func() (int, error) {
		clientwg.Done()
		return 10, nil
	}

	// Get should only request a client once, use cached afterward
	t.Run("get", func(t *testing.T) {
		clientwg.Add(1) // expect one call to get client
		c := NewLazyLoad(tc)
		rw, err := c.Get()
		assert.NoError(t, err)
		assert.NotNil(t, rw)
		assert.NotNil(t, c.state)

		// used cached client
		rw, err = c.Get()
		assert.NoError(t, err)
		assert.NotNil(t, rw)
		clientwg.Wait()
	})

	// Clear removes the cached client, should refetch
	t.Run("clear", func(t *testing.T) {
		clientwg.Add(2) // expect two calls to get client

		c := NewLazyLoad(tc)
		rw, err := c.Get()
		assert.NotNil(t, rw)
		assert.NoError(t, err)

		c.Reset()

		rw, err = c.Get()
		assert.NotNil(t, rw)
		assert.NoError(t, err)
		clientwg.Wait()
	})

	// Race checks a race condition of Getting and Clearing a new client
	t.Run("race", func(t *testing.T) {
		clientwg.Add(1) // expect one call to get client

		c := NewLazyLoad(tc)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			rw, err := c.Get()
			assert.NoError(t, err)
			assert.NotNil(t, rw)
			wg.Done()
		}()
		go func() {
			c.Reset()
			wg.Done()
		}()
		wg.Wait()
		clientwg.Wait()
	})
}
