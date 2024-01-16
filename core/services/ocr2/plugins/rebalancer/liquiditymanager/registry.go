package liquiditymanager

import (
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type Registry struct {
	liquidityManagers map[models.NetworkSelector]models.Address
	mu                *sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		liquidityManagers: make(map[models.NetworkSelector]models.Address),
		mu:                &sync.RWMutex{},
	}
}

func (r *Registry) Add(net models.NetworkSelector, addr models.Address) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.liquidityManagers[net] = addr
}

func (r *Registry) Get(net models.NetworkSelector) (models.Address, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	addr, exists := r.liquidityManagers[net]
	return addr, exists
}

func (r *Registry) GetAll() map[models.NetworkSelector]models.Address {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cp := make(map[models.NetworkSelector]models.Address, len(r.liquidityManagers))
	for k, v := range r.liquidityManagers {
		cp[k] = v
	}
	return cp
}
