package webapi

import (
	"math/rand"
	"sync"

	"github.com/pkg/errors"
)

var ErrNoGateways = errors.New("no gateways available")

type RoundRobinSelector struct {
	items []string
	index int
	mu    sync.Mutex
}

// WithFixedStart option will set the start point for the round robin selector.
func WithFixedStart() func(*RoundRobinSelector) {
	return func(rrs *RoundRobinSelector) {
		rrs.index = 0
	}
}

// NewRoundRobinSelector creates a selector that will select a string from a list using the round robin strategy.
// By default the index starts on a random start point in the list.
func NewRoundRobinSelector(items []string, opts ...func(*RoundRobinSelector)) *RoundRobinSelector {
	var index int
	if len(items) > 0 {
		index = rand.Intn(len(items)) //nolint:gosec // No need for crpto secure randomness to select an index
	}
	rrs := &RoundRobinSelector{
		items: items,
		index: index,
	}

	for _, opt := range opts {
		opt(rrs)
	}

	return rrs
}

func (r *RoundRobinSelector) NextGateway() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.items) == 0 {
		return "", ErrNoGateways
	}

	item := r.items[r.index]
	r.index = (r.index + 1) % len(r.items)
	return item, nil
}
