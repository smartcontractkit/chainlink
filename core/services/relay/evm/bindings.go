package evm

import (
	"context"
	"fmt"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// key is contract name
type contractBindings map[string]readBindings

type readBindings struct {
	// contractFilter is used to filter over all events or any subset of events with same filtering parameters.
	// if an event is present in the contract filter, it can't define its own filter in the event binding.
	contractFilter logpoller.Filter
	// key is read name
	bindings map[string]readBinding
}

func (b contractBindings) GetReadBinding(contractName, readName string) (readBinding, error) {
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

func (b contractBindings) AddReadBinding(contractName, readName string, reader readBinding) {
	rbs, rbsExists := b[contractName]
	if !rbsExists {
		rbs = readBindings{}
		b[contractName] = rbs
	}
	rbs.bindings[readName] = reader
}

func (b contractBindings) Bind(ctx context.Context, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		rbs, rbsExist := b[bc.Name]
		if !rbsExist {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}
		for _, r := range rbs.bindings {
			if err := r.Bind(ctx, bc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b contractBindings) ForEach(ctx context.Context, fn func(readBinding, context.Context) error) error {
	for _, rbs := range b {
		for _, rb := range rbs.bindings {
			if err := fn(rb, ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
