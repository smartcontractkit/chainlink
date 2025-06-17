package webapi

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundRobinSelector(t *testing.T) {
	gateways := []string{"gateway1", "gateway2", "gateway3"}
	rr := NewRoundRobinSelector(gateways, WithFixedStart())

	expectedOrder := []string{"gateway1", "gateway2", "gateway3", "gateway1", "gateway2", "gateway3"}

	for i, expected := range expectedOrder {
		got, err := rr.NextGateway()
		require.NoError(t, err, "unexpected error on iteration %d", i)
		assert.Equal(t, expected, got, "unexpected gateway at iteration %d", i)
	}
}

// TestNewRoundRobinSelector_NotAlwaysZero ensures that the start index is not always 0
// because it starts at a random index.  Does not prove randomness but proves that
// the start value is not always the first index.
func TestNewRoundRobinSelector_NotAlwaysZero(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	numInstances := 100 // Check a reasonable number of instances
	alwaysZero := true

	for i := 0; i < numInstances; i++ {
		selector := NewRoundRobinSelector(items)
		if selector.index != 0 {
			alwaysZero = false
			break
		}
	}

	assert.False(t, alwaysZero)
}

func TestNewRoundRobinSelector_AlwaysZero(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	numInstances := 100 // Check a reasonable number of instances
	alwaysZero := true

	for i := 0; i < numInstances; i++ {
		selector := NewRoundRobinSelector(items, WithFixedStart())
		if selector.index != 0 {
			alwaysZero = false
			break
		}
	}

	assert.True(t, alwaysZero)
}

func TestRoundRobinSelector_Empty(t *testing.T) {
	rr := NewRoundRobinSelector([]string{}, WithFixedStart())

	_, err := rr.NextGateway()
	assert.ErrorIs(t, err, ErrNoGateways, "expected ErrNoGateways when no gateways are available")
}

func TestRoundRobinSelector_Concurrency(t *testing.T) {
	gateways := []string{"gateway1", "gateway2", "gateway3"}
	rr := NewRoundRobinSelector(gateways, WithFixedStart())

	var wg sync.WaitGroup
	numRequests := 100
	results := make(chan string, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gw, err := rr.NextGateway()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			results <- gw
		}()
	}

	wg.Wait()
	close(results)

	counts := make(map[string]int)
	for result := range results {
		counts[result]++
	}

	expectedCount := numRequests / len(gateways)
	for _, gateway := range gateways {
		assert.InDelta(t, expectedCount, counts[gateway], 1, "unexpected request distribution for %s", gateway)
	}
}
