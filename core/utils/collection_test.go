package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
