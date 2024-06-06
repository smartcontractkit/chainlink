package evm

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// key is contract name
type bindings map[string]*contractBinding

type FilterRegisterer struct {
	pollingFilter logpoller.Filter
	filterLock    sync.Mutex
	isRegistered  bool
}

// contractBinding stores read bindings and manages the common contract event filter.
type contractBinding struct {
	// FilterRegisterer is used to manage polling filter registration for the common contract filter.
	// The common contract filter should be used by events that share filtering args.
	FilterRegisterer
	// key is read name method, event or event keys used for queryKey.
	readBindings map[string]readBinding
	bound        bool
}

func (b bindings) GetReadBinding(contractName, readName string) (readBinding, error) {
	cb, cbExists := b[contractName]
	if !cbExists {
		return nil, fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidType, contractName)
	}

	rb, rbExists := cb.readBindings[readName]
	if !rbExists {
		return nil, fmt.Errorf("%w: no readName named %s in contract %s", commontypes.ErrInvalidType, readName, contractName)
	}
	return rb, nil
}

func (b bindings) AddReadBinding(contractName, readName string, rb readBinding) {
	cb, cbExists := b[contractName]
	if !cbExists {
		cb = &contractBinding{readBindings: make(map[string]readBinding)}
		b[contractName] = cb
	}
	cb.readBindings[readName] = rb
}

// Bind binds contract addresses and creates event binding filters and the common contract filter.
func (b bindings) Bind(ctx context.Context, lp logpoller.LogPoller, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		cb, cbExists := b[bc.Name]
		if !cbExists {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}

		cb.pollingFilter.Addresses = evmtypes.AddressArray{common.HexToAddress(bc.Address)}
		cb.pollingFilter.Name = logpoller.FilterName(bc.Name+"."+uuid.NewString(), bc.Address)
		cb.bound = true

		// we are changing contract address reference, so we need to unregister old filter it exists
		if err := cb.Unregister(ctx, lp); err != nil {
			return err
		}

		// if contract event filter isn't already registered then it will be on startup
		// if its already registered then we are overriding it because(address) has changed
		if cb.isRegistered {
			if err := cb.Register(ctx, lp); err != nil {
				return err
			}
		}

		for _, rb := range cb.readBindings {
			if err := rb.Bind(ctx, bc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b bindings) ForEach(ctx context.Context, fn func(context.Context, *contractBinding) error) error {
	for _, cb := range b {
		if err := fn(ctx, cb); err != nil {
			return err
		}
	}
	return nil
}

// Register registers the common contract filter.
func (cb *contractBinding) Register(ctx context.Context, lp logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if cb.bound && len(cb.pollingFilter.EventSigs) > 0 && !lp.HasFilter(cb.pollingFilter.Name) {
		if err := lp.RegisterFilter(ctx, cb.pollingFilter); err != nil {
			return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
		}
		cb.isRegistered = true
	}

	return nil
}

// Unregister unregisters the common contract filter.
func (cb *contractBinding) Unregister(ctx context.Context, lp logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if !lp.HasFilter(cb.pollingFilter.Name) {
		return nil
	}

	if err := lp.UnregisterFilter(ctx, cb.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	cb.isRegistered = false

	return nil
}
