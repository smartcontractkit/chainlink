package liquiditymanager

import (
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type Registry struct {
	rebalancers map[models.NetworkSelector]models.Address
	mu          *sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		rebalancers: make(map[models.NetworkSelector]models.Address),
		mu:          &sync.RWMutex{},
	}
}

func (r *Registry) Add(net models.NetworkSelector, addr models.Address) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rebalancers[net] = addr
}

func (r *Registry) Get(net models.NetworkSelector) (models.Address, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	addr, exists := r.rebalancers[net]
	return addr, exists
}

func (r *Registry) GetAll() map[models.NetworkSelector]models.Address {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cp := make(map[models.NetworkSelector]models.Address, len(r.rebalancers))
	for k, v := range r.rebalancers {
		cp[k] = v
	}
	return cp
}
