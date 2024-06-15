package evm

import (
	"context"
	"fmt"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// bindings manage all contract bindings, key is contract name.
type bindings map[string]*contractBinding

func (b bindings) GetReadBinding(contractName, readName string) (readBinding, error) {
	// GetReadBindings should only be called after Chain Reader init.
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

// AddReadBinding adds read bindings. Calling this outside of Chain Reader init is not thread safe.
func (b bindings) AddReadBinding(contractName, readName string, rb readBinding) {
	cb, cbExists := b[contractName]
	if !cbExists {
		cb = &contractBinding{
			name:         contractName,
			readBindings: make(map[string]readBinding),
		}
		b[contractName] = cb
	}
	cb.readBindings[readName] = rb
}

// Bind binds contract addresses to contract bindings and read bindings.
// Bind also registers the common contract polling filter and eventBindings polling filters.
func (b bindings) Bind(ctx context.Context, lp logpoller.LogPoller, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		cb, cbExists := b[bc.Name]
		if !cbExists {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}

		if err := cb.Bind(ctx, lp, bc); err != nil {
			return err
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
