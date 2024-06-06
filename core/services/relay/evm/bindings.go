package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// bindings manage all contract bindings, key is contract name.
type bindings map[string]*contractBinding

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

// Bind binds contract addresses to contract binding and read bindings.
// Bind also registers the common contract polling filter and eventBindings polling filters.
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
