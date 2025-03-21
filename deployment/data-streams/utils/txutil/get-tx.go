package txutil

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

// HasContractAddress is the generic interface for anything
// that knows which contract address to use.
type HasContractAddress interface {
	GetContractAddress() common.Address
}

// ContractLoader is a function that, given an Environment, chain selector,
// and contract address, returns the strongly-typed contract “state”
type ContractLoader[S any] func(
	e deployment.Environment,
	chainSelector uint64,
	contractAddr string,
) (*S, error)

// ContractMethod is a function that, given a loaded contract state (pointer to S)
// and some config T, returns a transaction.
type ContractMethod[S any, T any] func(
	*S,
	T,
) (*goEthTypes.Transaction, error)

// GetTxs is the generic method that:
//
// Iterates over chainSelector => []T
// Loads the contract using the given loader
// Calls the method function to build the transaction
// Accumulates them into a slice of PreparedTx.
func GetTxs[S any, T HasContractAddress](
	e deployment.Environment,
	contractType string,
	configsByChain map[uint64][]T,
	loader ContractLoader[S],
	method ContractMethod[S, T],
) ([]*PreparedTx, error) {
	var preparedTxs []*PreparedTx

	for chainSelector, cfgs := range configsByChain {
		for _, cfg := range cfgs {
			contractAddr := cfg.GetContractAddress().Hex()

			state, err := loader(e, chainSelector, contractAddr)
			if err != nil {
				return nil, fmt.Errorf("failed to load contract state: %w", err)
			}

			tx, err := method(state, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to build transaction: %w", err)
			}

			preparedTxs = append(preparedTxs, &PreparedTx{
				Tx:            tx,
				ChainSelector: chainSelector,
				ContractType:  contractType,
			})
		}
	}
	return preparedTxs, nil
}
