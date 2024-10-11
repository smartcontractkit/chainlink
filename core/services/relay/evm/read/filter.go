package read

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

type Registrar interface {
	HasFilter(string) bool
	RegisterFilter(context.Context, logpoller.Filter) error
	UnregisterFilter(context.Context, string) error
}

type syncedFilter struct {
	// internal state properties
	mu     sync.RWMutex
	filter logpoller.Filter
}

func newSyncedFilter() *syncedFilter {
	return &syncedFilter{}
}

func (r *syncedFilter) Register(ctx context.Context, registrar Registrar) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !registrar.HasFilter(r.filter.Name) {
		if err := registrar.RegisterFilter(ctx, r.filter); err != nil {
			return FilterError{
				Err:    fmt.Errorf("%w: %s", types.ErrInternal, err.Error()),
				Action: "register",
				Filter: r.filter,
			}
		}
	}

	return nil
}

func (r *syncedFilter) Unregister(ctx context.Context, registrar Registrar) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !registrar.HasFilter(r.filter.Name) {
		return nil
	}

	if err := registrar.UnregisterFilter(ctx, r.filter.Name); err != nil {
		return FilterError{
			Err:    fmt.Errorf("%w: %s", types.ErrInternal, err.Error()),
			Action: "unregister",
			Filter: r.filter,
		}
	}

	return nil
}

func (r *syncedFilter) SetFilter(filter logpoller.Filter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.filter = filter
}

func (r *syncedFilter) SetName(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.filter.Name = name
}

func (r *syncedFilter) AddAddress(address common.Address) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.filter.Addresses = append(r.filter.Addresses, address)
}

func (r *syncedFilter) RemoveAddress(address common.Address) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var addrIdx int
	for idx, addr := range r.filter.Addresses {
		if addr.Hex() == address.Hex() {
			addrIdx = idx
		}
	}

	r.filter.Addresses[addrIdx] = r.filter.Addresses[len(r.filter.Addresses)-1]
	r.filter.Addresses = r.filter.Addresses[:len(r.filter.Addresses)-1]
}

func (r *syncedFilter) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.filter.Addresses)
}

func (r *syncedFilter) HasEventSigs() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.filter.EventSigs) > 0 && len(r.filter.Addresses) > 0
}
