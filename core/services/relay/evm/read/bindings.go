package read

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
)

type Reader interface {
	BatchCall(address common.Address, params, retVal any) (Call, error)
	GetLatestValueWithHeadData(ctx context.Context, addr common.Address, confidence primitives.ConfidenceLevel, params, returnVal any) (*commontypes.Head, error)
	QueryKey(context.Context, common.Address, query.KeyFilter, query.LimitAndSort, any) ([]commontypes.Sequence, error)

	Bind(context.Context, ...common.Address) error
	Unbind(context.Context, ...common.Address) error
	SetCodec(commontypes.RemoteCodec)

	Register(context.Context) error
	Unregister(context.Context) error
}

type BindingsRegistry struct {
	// dependencies
	batch BatchCaller

	// private state
	mu               sync.RWMutex
	contractBindings map[string]*contractBinding
	contractLookup   *lookup
}

func NewBindingsRegistry() *BindingsRegistry {
	return &BindingsRegistry{
		contractBindings: make(map[string]*contractBinding),
		contractLookup:   newLookup(),
	}
}

func (b *BindingsRegistry) SetBatchCaller(caller BatchCaller) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.batch = caller
}

func (b *BindingsRegistry) HasContractBinding(contractName string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.contractBindings[contractName]

	return ok
}

// GetReader should only be called after Chain Reader is started.
func (b *BindingsRegistry) GetReader(readIdentifier string) (Reader, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	values, ok := b.contractLookup.getContractForReadName(readIdentifier)
	if !ok {
		return nil, "", fmt.Errorf("%w: %w", commontypes.ErrInvalidType, newMissingReadIdentifierErr(readIdentifier))
	}

	cb, cbExists := b.contractBindings[values.contract]
	if !cbExists {
		return nil, "", fmt.Errorf("%w: %w", commontypes.ErrInvalidType, newMissingContractErr(readIdentifier, values.contract))
	}

	binding, rbExists := cb.GetReaderNamed(values.readName)
	if !rbExists {
		return nil, "", fmt.Errorf("%w: %w", commontypes.ErrInvalidType, newMissingReadNameErr(readIdentifier, values.contract, values.readName))
	}

	return binding, values.address, nil
}

func (b *BindingsRegistry) AddReader(contractName, readName string, rdr Reader) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if binding, ok := rdr.(*EventBinding); ok {
		// unwrap codec type naming for event data words and topics to be used by lookup for Querying by Value Comparators
		// For e.g. "params.contractName.eventName.IndexedTopic" -> "eventName.IndexedTopic"
		// or "params.contractName.eventName.someFieldInData" -> "eventName.someFieldInData"
		for name := range binding.eventTypes {
			split := strings.Split(name, ".")
			if len(split) < 3 || split[1] != contractName {
				return fmt.Errorf("%w: invalid event type name %s", commontypes.ErrInvalidType, name)
			}

			b.contractLookup.addReadNameForContract(contractName, strings.Join(split[2:], "."))
		}
	}

	b.contractLookup.addReadNameForContract(contractName, readName)

	cb, cbExists := b.contractBindings[contractName]
	if !cbExists {
		cb = newContractBinding(contractName)
		b.contractBindings[contractName] = cb
	}

	cb.AddReaderNamed(readName, rdr)

	return nil
}

// Bind binds contract addresses to contract bindings and read bindings.
// Bind also registers the common contract polling filter and eventBindings polling filters.
func (b *BindingsRegistry) Bind(ctx context.Context, reg Registrar, bindings []commontypes.BoundContract) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, binding := range bindings {
		contract, exists := b.contractBindings[binding.Name]
		if !exists {
			return fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newMissingContractErr("binding contract", binding.Name))
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

func (b *BindingsRegistry) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var batchCall BatchCall

	for binding, contractBatch := range request {
		cb := b.contractBindings[binding.Name]

		for idx := range contractBatch {
			req := contractBatch[idx]

			values, ok := b.contractLookup.getContractForReadName(binding.ReadIdentifier(req.ReadName))
			if !ok {
				return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newMissingReadNameErr(binding.ReadIdentifier(req.ReadName), binding.Name, req.ReadName))
			}

			rdr, exists := cb.GetReaderNamed(values.readName)
			if !exists {
				return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newMissingReadNameErr(binding.ReadIdentifier(req.ReadName), binding.Name, req.ReadName))
			}

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
		for _, methodResult := range contractResult {
			binding := commontypes.BoundContract{
				Name:    contractName,
				Address: methodResult.Address,
			}

			brr := commontypes.BatchReadResult{
				ReadName: methodResult.MethodName,
			}

			brr.SetResult(methodResult.ReturnValue, methodResult.Err)

			if _, ok := batchGetLatestValuesResults[binding]; !ok {
				batchGetLatestValuesResults[binding] = make(commontypes.ContractBatchResults, 0)
			}

			batchGetLatestValuesResults[binding] = append(batchGetLatestValuesResults[binding], brr)
		}
	}

	return batchGetLatestValuesResults, err
}

// Unbind unbinds contract addresses to contract bindings and read bindings.
func (b *BindingsRegistry) Unbind(ctx context.Context, reg Registrar, bindings []commontypes.BoundContract) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, binding := range bindings {
		contract, exists := b.contractBindings[binding.Name]
		if !exists {
			return newMissingContractErr("unbinding contract", binding.Name)
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

func (b *BindingsRegistry) RegisterAll(ctx context.Context, reg Registrar) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		if err := errors.Join(cb.RegisterReaders(ctx), cb.Register(ctx, reg)); err != nil {
			return err
		}
	}

	return nil
}

func (b *BindingsRegistry) UnregisterAll(ctx context.Context, reg Registrar) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		if err := errors.Join(cb.UnregisterReaders(ctx), cb.Unregister(ctx, reg)); err != nil {
			return err
		}
	}

	return nil
}

func (b *BindingsRegistry) SetCodecAll(codec commontypes.RemoteCodec) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, cb := range b.contractBindings {
		cb.SetCodecAll(codec)
	}
}

func (b *BindingsRegistry) SetFilter(name string, filter logpoller.Filter) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	contract, ok := b.contractBindings[name]
	if !ok {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newMissingContractErr("set filter", name))
	}

	contract.SetFilter(filter)

	return nil
}

func (b *BindingsRegistry) ReadTypeIdentifier(readName string, forEncoding bool) string {
	values, ok := b.contractLookup.getContractForReadName(readName)
	if !ok {
		return ""
	}

	return codec.WrapItemType(values.contract, values.readName, forEncoding)
}

// confidenceToConfirmations matches predefined chain agnostic confidence levels to predefined EVM finality.
func confidenceToConfirmations(
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations,
	confidenceLevel primitives.ConfidenceLevel,
) (evmtypes.Confirmations, error) {
	confirmations, exists := confirmationsMapping[confidenceLevel]
	if !exists {
		return 0, fmt.Errorf("%w: missing mapping for confidence level: %s", commontypes.ErrInvalidConfig, confidenceLevel)
	}

	return confirmations, nil
}
