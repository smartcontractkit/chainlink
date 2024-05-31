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
type bindings map[string]*contractBindings

type FilterRegisterer struct {
	pollingFilter logpoller.Filter
	filterLock    sync.Mutex
	isRegistered  bool
}

type contractBindings struct {
	// FilterRegisterer is used to manage polling filter registration.
	FilterRegisterer
	// key is read name
	bindings map[string]readBinding
}

func (b bindings) GetReadBinding(contractName, readName string) (readBinding, error) {
	rb, rbExists := b[contractName]
	if !rbExists {
		return nil, fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidType, contractName)
	}

	reader, readerExists := rb.bindings[readName]
	if !readerExists {
		return nil, fmt.Errorf("%w: no readName named %s in contract %s", commontypes.ErrInvalidType, readName, contractName)
	}
	return reader, nil
}

func (b bindings) AddReadBinding(contractName, readName string, rb readBinding) {
	rbs, rbsExists := b[contractName]
	if !rbsExists {
		rbs = &contractBindings{}
		b[contractName] = rbs
	}
	rbs.bindings[readName] = rb
}

func (b bindings) Bind(ctx context.Context, logPoller logpoller.LogPoller, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		rbs, rbsExist := b[bc.Name]
		if !rbsExist {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}

		rbs.pollingFilter.Addresses = evmtypes.AddressArray{common.HexToAddress(bc.Address)}
		rbs.pollingFilter.Name = logpoller.FilterName(bc.Name+"."+uuid.NewString(), bc.Address)

		// we are changing contract address reference, so we need to unregister old filters
		if err := rbs.Unregister(ctx, logPoller); err != nil {
			return err
		}

		for _, r := range rbs.bindings {
			r.Bind(bc)
		}

		// if contract event filters aren't already registered then they will on startup
		// if they are already registered then we are overriding them because contract binding (address) has changed
		if rbs.isRegistered {
			if err := rbs.Register(ctx, logPoller); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b bindings) ForEach(ctx context.Context, fn func(context.Context, *contractBindings) error) error {
	for _, rbs := range b {
		if err := fn(ctx, rbs); err != nil {
			return err
		}
	}
	return nil
}

// Register registers polling filters.
func (rb *contractBindings) Register(ctx context.Context, logPoller logpoller.LogPoller) error {
	rb.filterLock.Lock()
	defer rb.filterLock.Unlock()

	rb.isRegistered = true

	if logPoller.HasFilter(rb.pollingFilter.Name) {
		return nil
	}

	if err := logPoller.RegisterFilter(ctx, rb.pollingFilter); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	for _, binding := range rb.bindings {
		if err := binding.Register(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Unregister unregisters polling filters.
func (rb *contractBindings) Unregister(ctx context.Context, logPoller logpoller.LogPoller) error {
	rb.filterLock.Lock()
	defer rb.filterLock.Unlock()

	if !logPoller.HasFilter(rb.pollingFilter.Name) {
		return nil
	}

	if err := logPoller.UnregisterFilter(ctx, rb.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	for _, binding := range rb.bindings {
		if err := binding.Unregister(ctx); err != nil {
			return err
		}
	}
	return nil
}
