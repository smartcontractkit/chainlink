package evm

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// key is contract name
type bindings map[string]*contractBindings

type contractBindings struct {
	// contractFilter is used to filter over all events or any subset of events with same filtering parameters.
	// if an event is present in the contract filter, it can't define its own filter in the event binding.
	contractFilter            logpoller.Filter
	filterLock                sync.Mutex
	areEventFiltersRegistered bool
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

		rbs.contractFilter.Addresses = append(rbs.contractFilter.Addresses, common.HexToAddress(bc.Address))
		rbs.contractFilter.Name = logpoller.FilterName(bc.Name, bc.Address)

		if err := rbs.UnregisterEventFilters(ctx, logPoller); err != nil {
			return err
		}

		// if contract event filter isn't already registered then it will be by chain reader on startup
		// if it is already registered then we are overriding filters registered on startup
		if rbs.areEventFiltersRegistered {
			return rbs.RegisterEventFilters(ctx, logPoller)
		}

		for _, r := range rbs.bindings {
			r.Bind(bc)
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

func (rb *contractBindings) RegisterEventFilters(ctx context.Context, logPoller logpoller.LogPoller) error {
	rb.filterLock.Lock()
	defer rb.filterLock.Unlock()

	rb.areEventFiltersRegistered = true

	if logPoller.HasFilter(rb.contractFilter.Name) {
		return nil
	}

	if err := logPoller.RegisterFilter(ctx, rb.contractFilter); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return nil
}

func (rb *contractBindings) UnregisterEventFilters(ctx context.Context, logPoller logpoller.LogPoller) error {
	rb.filterLock.Lock()
	defer rb.filterLock.Unlock()

	if !logPoller.HasFilter(rb.contractFilter.Name) {
		return nil
	}

	if err := logPoller.UnregisterFilter(ctx, rb.contractFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return nil
}
