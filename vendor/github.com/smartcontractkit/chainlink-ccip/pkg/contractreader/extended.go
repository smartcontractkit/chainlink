package contractreader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/internal/libs/slicelib"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// Extended version of a ContractReader.
type Extended interface {
	types.ContractReader
	GetBindings(contractName string) []ExtendedBoundContract
}

type ExtendedBoundContract struct {
	BoundAt time.Time
	Binding types.BoundContract
}

// extendedContractReader is an extended version of the contract reader.
type extendedContractReader struct {
	types.ContractReader
	contractBindingsByName map[string][]ExtendedBoundContract
	mu                     *sync.RWMutex
}

func NewExtendedContractReader(baseContractReader types.ContractReader) Extended {
	return &extendedContractReader{
		ContractReader:         baseContractReader,
		contractBindingsByName: make(map[string][]ExtendedBoundContract),
		mu:                     &sync.RWMutex{},
	}
}

func (e *extendedContractReader) Bind(ctx context.Context, allBindings []types.BoundContract) error {
	validBindings := slicelib.Filter(allBindings, func(b types.BoundContract) bool { return !e.bindingExists(b) })
	if len(validBindings) == 0 {
		return nil
	}

	err := e.ContractReader.Bind(ctx, validBindings)
	if err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	for _, binding := range validBindings {
		e.contractBindingsByName[binding.Name] = append(e.contractBindingsByName[binding.Name], ExtendedBoundContract{
			BoundAt: time.Now(),
			Binding: binding,
		})
	}

	return nil
}

func (e *extendedContractReader) GetBindings(contractName string) []ExtendedBoundContract {
	e.mu.RLock()
	defer e.mu.RUnlock()

	bindings, exists := e.contractBindingsByName[contractName]
	if !exists {
		return []ExtendedBoundContract{}
	}
	return bindings
}

func (e *extendedContractReader) bindingExists(b types.BoundContract) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, boundContracts := range e.contractBindingsByName {
		for _, boundContract := range boundContracts {
			if boundContract.Binding.Key() == b.Key() {
				return true
			}
		}
	}
	return false
}

// Interface compliance check
var _ Extended = (*extendedContractReader)(nil)
