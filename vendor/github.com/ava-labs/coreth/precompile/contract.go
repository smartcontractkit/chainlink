// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	selectorLen = 4
)

type RunStatefulPrecompileFunc func(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)

// PrecompileAccessibleState defines the interface exposed to stateful precompile contracts
type PrecompileAccessibleState interface {
	GetStateDB() StateDB
	GetBlockContext() BlockContext
	NativeAssetCall(caller common.Address, input []byte, suppliedGas uint64, gasGost uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)
}

// BlockContext defines an interface that provides information to a stateful precompile
// about the block that activates the upgrade. The precompile can access this information
// to initialize its state.
type BlockContext interface {
	Number() *big.Int
	Timestamp() *big.Int
}

// ChainContext defines an interface that provides information to a stateful precompile
// about the chain configuration. The precompile can access this information to initialize
// its state.
type ChainConfig interface {
	// Note: None of the existing stateful precompiles currently access chain config information
	// in Configure so this interface is empty.
}

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	SetCode(common.Address, []byte)

	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64

	GetBalance(common.Address) *big.Int
	AddBalance(common.Address, *big.Int)
	SubBalance(common.Address, *big.Int)

	SubBalanceMultiCoin(common.Address, common.Hash, *big.Int)
	AddBalanceMultiCoin(common.Address, common.Hash, *big.Int)
	GetBalanceMultiCoin(common.Address, common.Hash) *big.Int

	CreateAccount(common.Address)
	Exist(common.Address) bool
}

// StatefulPrecompiledContract is the interface for executing a precompiled contract
type StatefulPrecompiledContract interface {
	// Run executes the precompiled contract.
	Run(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)
}

// statefulPrecompileFunction defines a function implemented by a stateful precompile
type statefulPrecompileFunction struct {
	// selector is the 4 byte function selector for this function
	// This should be calculated from the function signature using CalculateFunctionSelector
	selector []byte
	// execute is performed when this function is selected
	execute RunStatefulPrecompileFunc
}

// newStatefulPrecompileFunction creates a stateful precompile function with the given arguments
//nolint:unused,deadcode
func newStatefulPrecompileFunction(selector []byte, execute RunStatefulPrecompileFunc) *statefulPrecompileFunction {
	return &statefulPrecompileFunction{
		selector: selector,
		execute:  execute,
	}
}

// statefulPrecompileWithFunctionSelectors implements StatefulPrecompiledContract by using 4 byte function selectors to pass
// off responsibilities to internal execution functions.
// Note: because we only ever read from [functions] there no lock is required to make it thread-safe.
type statefulPrecompileWithFunctionSelectors struct {
	fallback  *statefulPrecompileFunction
	functions map[string]*statefulPrecompileFunction
}

// newStatefulPrecompileWithFunctionSelectors generates new StatefulPrecompile using [functions] as the available functions and [fallback]
// as an optional fallback if there is no input data. Note: the selector of [fallback] will be ignored, so it is required to be left empty.
//nolint:unused,deadcode
func newStatefulPrecompileWithFunctionSelectors(fallback *statefulPrecompileFunction, functions []*statefulPrecompileFunction) StatefulPrecompiledContract {
	// Ensure that if a fallback is present, it does not have a mistakenly populated function selector.
	if fallback != nil && len(fallback.selector) != 0 {
		panic(fmt.Errorf("fallback function cannot specify non-zero length function selector"))
	}

	// Construct the contract and populate [functions].
	contract := &statefulPrecompileWithFunctionSelectors{
		fallback:  fallback,
		functions: make(map[string]*statefulPrecompileFunction),
	}
	for _, function := range functions {
		_, exists := contract.functions[string(function.selector)]
		if exists {
			panic(fmt.Errorf("cannot create stateful precompile with duplicated function selector: %q", function.selector))
		}
		contract.functions[string(function.selector)] = function
	}

	return contract
}

// Run selects the function using the 4 byte function selector at the start of the input and executes the underlying function on the
// given arguments.
func (s *statefulPrecompileWithFunctionSelectors) Run(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	// If there is no input data present, call the fallback function if present.
	if len(input) == 0 && s.fallback != nil {
		return s.fallback.execute(accessibleState, caller, addr, nil, suppliedGas, readOnly)
	}

	// Otherwise, an unexpected input size will result in an error.
	if len(input) < selectorLen {
		return nil, suppliedGas, fmt.Errorf("missing function selector to precompile - input length (%d)", len(input))
	}

	// Use the function selector to grab the correct function
	selector := input[:selectorLen]
	functionInput := input[selectorLen:]
	function, ok := s.functions[string(selector)]
	if !ok {
		return nil, suppliedGas, fmt.Errorf("invalid function selector %#x", selector)
	}

	return function.execute(accessibleState, caller, addr, functionInput, suppliedGas, readOnly)
}
