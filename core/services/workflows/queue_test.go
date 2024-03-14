package workflows

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestQueue(t *testing.T) {
	ctx := testutils.Context(t)
	q := newQueue[int]()
	var wg sync.WaitGroup
	q.start(ctx, &wg)

	ints := []int{1, 2, 3, 4, 5}
	for _, i := range ints {
		q.in <- i
	}

	got := []int{}
	for i := 0; i < 5; i++ {
		got = append(got, <-q.out)
	}

	assert.Equal(t, ints, got)

	assert.Len(t, q.out, 0)
}
