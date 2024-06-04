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

type contractBinding struct {
	// FilterRegisterer is used to manage polling filter registration for the contact wide event filter.
	FilterRegisterer
	// key is read name method, event or event keys used for queryKey.
	readBindings map[string]readBinding
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

func (b bindings) Bind(ctx context.Context, logPoller logpoller.LogPoller, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		cb, cbExists := b[bc.Name]
		if !cbExists {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}

		cb.pollingFilter.Addresses = evmtypes.AddressArray{common.HexToAddress(bc.Address)}
		cb.pollingFilter.Name = logpoller.FilterName(bc.Name+"."+uuid.NewString(), bc.Address)

		// we are changing contract address reference, so we need to unregister old filters if they exist
		if err := cb.Unregister(ctx, logPoller); err != nil {
			return err
		}

		for _, rb := range cb.readBindings {
			rb.Bind(bc)
		}

		// if contract event filters aren't already registered then they will on startup
		// if they are already registered then we are overriding them because contract binding (address) has changed
		if cb.isRegistered {
			if err := cb.Register(ctx, logPoller); err != nil {
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

// Register registers polling filters.
func (cb *contractBinding) Register(ctx context.Context, logPoller logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	cb.isRegistered = true

	if logPoller.HasFilter(cb.pollingFilter.Name) {
		return nil
	}

	if err := logPoller.RegisterFilter(ctx, cb.pollingFilter); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	for _, rb := range cb.readBindings {
		if err := rb.Register(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Unregister unregisters polling filters.
func (cb *contractBinding) Unregister(ctx context.Context, logPoller logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if !logPoller.HasFilter(cb.pollingFilter.Name) {
		return nil
	}

	if err := logPoller.UnregisterFilter(ctx, cb.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	for _, rb := range cb.readBindings {
		if err := rb.Unregister(ctx); err != nil {
			return err
		}
	}
	return nil
}
