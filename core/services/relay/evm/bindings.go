package evm

import (
	"context"
	"fmt"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// bindings manage all contract bindings, key is contract name.

type bindings struct {
	contractBindings map[string]*contractBinding
	BatchCaller
}

func (b bindings) GetReadBinding(contractName, readName string) (readBinding, error) {
	// GetReadBindings should only be called after Chain Reader init.
	cb, cbExists := b.contractBindings[contractName]
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
	cb, cbExists := b.contractBindings[contractName]
	if !cbExists {
		cb = &contractBinding{
			name:         contractName,
			readBindings: make(map[string]readBinding),
		}
		b.contractBindings[contractName] = cb
	}
	cb.readBindings[readName] = rb
}

func (b bindings) SetBatchCaller(batchCaller BatchCaller) {
	b.BatchCaller = batchCaller
}

// Bind binds contract addresses to contract bindings and read bindings.
// Bind also registers the common contract polling filter and eventBindings polling filters.
func (b bindings) Bind(ctx context.Context, lp logpoller.LogPoller, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		cb, cbExists := b.contractBindings[bc.Name]
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

func (b bindings) BatchGetLatestValue(ctx context.Context, batchRequests BatchGetLatestValueRequests) error {
	var batchCall BatchCall
	for contractName, contractBatch := range batchRequests {
		cb := b.contractBindings[contractName]
		for i := range contractBatch {
			req := contractBatch[i]
			switch rb := cb.readBindings[req.readName].(type) {
			case *methodBinding:
				batchCall = append(batchCall, Call{
					contractAddress: rb.address,
					contractName:    cb.name,
					methodName:      rb.method,
					params:          req.params,
					returnVal:       req.params,
				})
				// results here will have chain specific method names.

			case *eventBinding:
				// TODO Use FilteredLogs to batch? This isn't a priority right now, but should get implemented at some point.
			}
		}
	}

	_, err := b.BatchCall(ctx, 0, batchCall)
	return err
}

func (b bindings) ForEach(ctx context.Context, fn func(context.Context, *contractBinding) error) error {
	for _, cb := range b.contractBindings {
		if err := fn(ctx, cb); err != nil {
			return err
		}
	}
	return nil
}
