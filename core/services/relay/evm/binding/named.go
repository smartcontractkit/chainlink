package binding

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type Reader interface {
	BatchCall(address common.Address, params, retVal any) (Call, error)
	GetLatestValue(ctx context.Context, addr common.Address, confidence primitives.ConfidenceLevel, params, returnVal any) error
	QueryKey(context.Context, common.Address, query.KeyFilter, query.LimitAndSort, any) ([]commontypes.Sequence, error)

	Bind(context.Context, ...common.Address) error
	Unbind(context.Context, ...common.Address) error
	SetCodec(commontypes.RemoteCodec)

	Register(context.Context) error
	Unregister(context.Context) error
}

type NamedBindings struct {
	// dependencies
	batch BatchCaller

	// private state
	mu               sync.RWMutex
	contractBindings map[string]*contractBinding
	contractLookup   *lookup
}

func NewNamedBindings() *NamedBindings {
	return &NamedBindings{
		contractBindings: make(map[string]*contractBinding),
		contractLookup:   newLookup(),
	}
}

func (b *NamedBindings) SetBatchCaller(caller BatchCaller) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.batch = caller
}

func (b *NamedBindings) HasContractBinding(contractName string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.contractBindings[contractName]

	return ok
}

// TODO: GetReader needs to accept a readName and do a mapping to bound contracts
func (b *NamedBindings) GetReader(readName string) (Reader, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// TODO: get contract name from readName using reverseLookup
	values, ok := b.contractLookup.getContractForReadName(readName)
	if !ok {
		return nil, "", fmt.Errorf("%w: no reader for read name %s", commontypes.ErrInvalidType, readName)
	}

	// GetReadBindings should only be called after Chain Reader init.
	cb, cbExists := b.contractBindings[values.contract]
	if !cbExists {
		return nil, "", fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidType, values.contract)
	}

	binding, rbExists := cb.GetReaderNamed(values.method)
	if !rbExists {
		return nil, "", fmt.Errorf("%w: no reader named %s in contract %s", commontypes.ErrInvalidType, values.method, values.contract)
	}

	return binding, values.address, nil
}

func (b *NamedBindings) AddReader(contractName, methodName string, rdr Reader) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.contractLookup.addMethodForContract(contractName, methodName)

	cb, cbExists := b.contractBindings[contractName]
	if !cbExists {
		cb = newContractBinding(contractName)
		b.contractBindings[contractName] = cb
	}

	cb.AddReaderNamed(methodName, rdr)
}

// Bind binds contract addresses to contract bindings and read bindings.
// Bind also registers the common contract polling filter and eventBindings polling filters.
func (b *NamedBindings) Bind(ctx context.Context, reg Registrar, bindings []commontypes.BoundContract) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, binding := range bindings {
		contract, exists := b.contractBindings[binding.Name]
		if !exists {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, binding.Name)
		}

		b.contractLookup.bindAddressForContract(binding.Name, binding.Address)

		if err := errors.Join(
			contract.Bind(ctx, reg, common.HexToAddress(binding.Address)),
			contract.BindReaders(ctx, common.HexToAddress(binding.Address)),
		); err != nil {
			return err
		}
	}

	return nil
}

func (b *NamedBindings) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var batchCall BatchCall

	for contractName, contractBatch := range request {
		cb := b.contractBindings[contractName]

		for i := range contractBatch {
			req := contractBatch[i]

			values, ok := b.contractLookup.getContractForReadName(req.ReadIdentifier)
			if !ok {
				return nil, fmt.Errorf("%w: no method for read name %s", commontypes.ErrInvalidType, req.ReadIdentifier)
			}

			rdr, exists := cb.GetReaderNamed(values.method)
			if !exists {
				continue
			}

			// TODO: need address for batch call
			call, err := rdr.BatchCall(common.HexToAddress(values.address), req.Params, req.ReturnVal)
			if err != nil {
				return nil, err
			}

			batchCall = append(batchCall, call)
		}
	}

	results, err := b.batch.BatchCall(ctx, 0, batchCall)
	if err != nil {
		return nil, err
	}

	// reconstruct results from batchCall and filteredLogs into common type while maintaining order from request.
	batchGetLatestValuesResults := make(commontypes.BatchGetLatestValuesResult)
	for contractName, contractResult := range results {
		batchGetLatestValuesResults[contractName] = commontypes.ContractBatchResults{}
		for _, methodResult := range contractResult {
			brr := commontypes.BatchReadResult{
				ReadIdentifier: types.BoundContract{
					Address: methodResult.Address,
					Name:    contractName,
				}.ReadIdentifier(methodResult.MethodName),
			}

			brr.SetResult(methodResult.ReturnValue, methodResult.Err)
			batchGetLatestValuesResults[contractName] = append(batchGetLatestValuesResults[contractName], brr)
		}
	}

	return batchGetLatestValuesResults, err
}

// Unbind unbinds contract addresses to contract bindings and read bindings.
func (b *NamedBindings) Unbind(ctx context.Context, reg Registrar, bindings []commontypes.BoundContract) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, binding := range bindings {
		contract, exists := b.contractBindings[binding.Name]
		if !exists {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, binding.Name)
		}

		b.contractLookup.unbindAddressForContract(binding.Name, binding.Address)

		if err := errors.Join(
			contract.Unbind(ctx, reg, common.HexToAddress(binding.Address)),
			contract.UnbindReaders(ctx, common.HexToAddress(binding.Address)),
		); err != nil {
			return err
		}
	}

	return nil
}

func (b *NamedBindings) RegisterAll(ctx context.Context, reg Registrar) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		if err := errors.Join(cb.RegisterReaders(ctx), cb.Register(ctx, reg)); err != nil {
			return err
		}
	}

	return nil
}

func (b *NamedBindings) UnregisterAll(ctx context.Context, reg Registrar) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		if err := errors.Join(cb.UnregisterReaders(ctx), cb.Unregister(ctx, reg)); err != nil {
			return err
		}
	}

	return nil
}

func (b *NamedBindings) SetCodecAll(codec commontypes.RemoteCodec) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		cb.SetCodecAll(codec)
	}
}

func (b *NamedBindings) WithFilter(name string, filter logpoller.Filter) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if contract, ok := b.contractBindings[name]; ok {
		contract.WithFilter(filter)
	}
}

func (b *NamedBindings) ReadTypeIdentifier(readName string, forEncoding bool) string {
	values, ok := b.contractLookup.getContractForReadName(readName)
	if !ok {
		return ""
	}

	return WrapItemType(values.contract, values.method, forEncoding)
}

// confidenceToConfirmations matches predefined chain agnostic confidence levels to predefined EVM finality.
func confidenceToConfirmations(confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations, confidenceLevel primitives.ConfidenceLevel) (evmtypes.Confirmations, error) {
	confirmations, exists := confirmationsMapping[confidenceLevel]
	if !exists {
		return 0, fmt.Errorf("missing mapping for confidence level: %s", confidenceLevel)
	}
	return confirmations, nil
}
