package soltxm

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

func TestBatchSplit(t *testing.T) {
	list := []int{}
	for i := 0; i < 100; i++ {
		list = append(list, i)
	}

	runs := []struct {
		name      string
		input     []int
		max       int // max per batch
		num       int // expected number of batches
		lastLen   int // expected number in last batch
		expectErr bool
	}{
		{"max=1", list, 1, len(list), 1, false},
		{"max=25", list, 25, 4, 25, false},
		{"max=33", list, 33, 4, 1, false},
		{"max=87", list, 87, 2, 13, false},
		{"max=len", list, len(list), 1, 100, false},
		{"max=len+1", list, len(list) + 1, 1, len(list), false}, // max exceeds len of list
		{"zero-list", []int{}, 1, 1, 0, false},                  // zero length list
		{"zero-max", list, 0, 0, 0, true},                       // zero as max input
	}

	for _, r := range runs {
		t.Run(r.name, func(t *testing.T) {
			batch, err := BatchSplit(r.input, r.max)
			if r.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, r.num, len(batch)) // check number of batches

			temp := []int{}
			for i := 0; i < len(batch); i++ {
				expectedLen := r.max
				if i == len(batch)-1 {
					expectedLen = r.lastLen // expect last batch to be less than max
				}
				assert.Equal(t, expectedLen, len(batch[i])) // check length of batch

				temp = append(temp, batch[i]...)
			}
			// assert order has not changed when list is reconstructed
			assert.Equal(t, r.input, temp)

		})
	}

}
