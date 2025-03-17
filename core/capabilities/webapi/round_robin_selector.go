package webapi

import (
	"sync"

	"github.com/pkg/errors"
)

var ErrNoGateways = errors.New("no gateways available")

type RoundRobinSelector struct {
	items []string
	index int
	mu    sync.Mutex
}

func NewRoundRobinSelector(items []string) *RoundRobinSelector {
	return &RoundRobinSelector{
		items: items,
		index: 0,
	}
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
