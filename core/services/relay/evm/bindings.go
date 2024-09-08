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

func (b bindings) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	var batchCall BatchCall
	toChainAgnosticMethodName := make(map[string]string)
	for contractName, contractBatch := range request {
		cb := b.contractBindings[contractName]
		for i := range contractBatch {
			req := contractBatch[i]
			switch rb := cb.readBindings[req.ReadName].(type) {
			case *methodBinding:
				toChainAgnosticMethodName[rb.method] = req.ReadName
				batchCall = append(batchCall, Call{
					ContractAddress: rb.address,
					ContractName:    cb.name,
					MethodName:      rb.method,
					Params:          req.Params,
					ReturnVal:       req.ReturnVal,
				})
				// results here will have chain specific method names.
			case *eventBinding:
				// TODO Use FilteredLogs to batch? This isn't a priority right now, but should get implemented at some point.
				return nil, fmt.Errorf("%w: events are not yet supported in batch get latest values", commontypes.ErrInvalidType)
			default:
				return nil, fmt.Errorf("%w: missing read binding type for contract: %s read: %s", commontypes.ErrInvalidType, contractName, req.ReadName)
			}
		}
	}

	results, err := b.BatchCall(ctx, 0, batchCall)
	if err != nil {
		return nil, err
	}

	// reconstruct results from batchCall and filteredLogs into common type while maintaining order from request.
	batchGetLatestValuesResults := make(commontypes.BatchGetLatestValuesResult)
	for contractName, contractResult := range results {
		batchGetLatestValuesResults[contractName] = commontypes.ContractBatchResults{}
		for _, methodResult := range contractResult {
			brr := commontypes.BatchReadResult{ReadName: toChainAgnosticMethodName[methodResult.MethodName]}
			brr.SetResult(methodResult.ReturnValue, methodResult.Err)
			batchGetLatestValuesResults[contractName] = append(batchGetLatestValuesResults[contractName], brr)
		}
	}

	return batchGetLatestValuesResults, err
}

func (b bindings) ForEach(ctx context.Context, fn func(context.Context, *contractBinding) error) error {
	for _, cb := range b.contractBindings {
		if err := fn(ctx, cb); err != nil {
			return err
		}
	}
	return nil
}
