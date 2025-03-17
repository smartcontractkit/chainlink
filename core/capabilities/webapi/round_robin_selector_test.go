package webapi

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundRobinSelector(t *testing.T) {
	gateways := []string{"gateway1", "gateway2", "gateway3"}
	rr := NewRoundRobinSelector(gateways)

	expectedOrder := []string{"gateway1", "gateway2", "gateway3", "gateway1", "gateway2", "gateway3"}

	for i, expected := range expectedOrder {
		got, err := rr.NextGateway()
		require.NoError(t, err, "unexpected error on iteration %d", i)
		assert.Equal(t, expected, got, "unexpected gateway at iteration %d", i)
	}
}

func TestRoundRobinSelector_Empty(t *testing.T) {
	rr := NewRoundRobinSelector([]string{})

	_, err := rr.NextGateway()
	assert.ErrorIs(t, err, ErrNoGateways, "expected ErrNoGateways when no gateways are available")
}

func TestRoundRobinSelector_Concurrency(t *testing.T) {
	gateways := []string{"gateway1", "gateway2", "gateway3"}
	rr := NewRoundRobinSelector(gateways)

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
