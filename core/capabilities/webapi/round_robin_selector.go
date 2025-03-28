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

// WithRandomStart starts selection at a random index.
func WithRandomStart() func(*RoundRobinSelector) {
	return func(rrs *RoundRobinSelector) {
		start := rand.Intn(len(rrs.items))
		rrs.index = start
	}
}

func NewRoundRobinSelector(items []string, opts ...func(*RoundRobinSelector)) *RoundRobinSelector {
	rrs := &RoundRobinSelector{
		items: items,
		index: 0,
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
