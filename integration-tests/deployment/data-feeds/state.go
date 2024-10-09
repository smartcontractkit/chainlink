package data_feeds

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gnosis_safe_1_3_0"
)

type DataFeedChainState struct {
	// TODO: Add multiple aggregators keyed by address
	Aggregator     *offchainaggregator.OffchainAggregator
	AggregatorAddr common.Address // Doesn't support Address()
	// Might be multiple per chain?
	Safe     *gnosis_safe_1_3_0.GnosisSafe130
	Timelock *owner_helpers.RBACTimelock
}

type DataFeedsState struct {
	Chains map[uint64]DataFeedChainState
}

func LoadChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (DataFeedChainState, error) {
	state := DataFeedChainState{}
	for address, tnv := range addresses {
		switch tnv.String() {
		case deployment.NewTypeAndVersion("OffchainAggregator", deployment.Version1_0_0).String():
			aggregator, err := offchainaggregator.NewOffchainAggregator(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Aggregator = aggregator
			state.AggregatorAddr = common.HexToAddress(address)
		case deployment.NewTypeAndVersion("GnosisSafe", deployment.Version1_0_0).String():
			safe, err := gnosis_safe_1_3_0.NewGnosisSafe130(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Safe = safe
		case deployment.NewTypeAndVersion("Timelock", deployment.Version1_0_0).String():
			timelock, err := owner_helpers.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Timelock = timelock
		default:
			// Ignore unknown type
			return state, fmt.Errorf("unknown type %s", tnv.String())
		}
	}
	return state, nil
}

func LoadOnchainState(e deployment.Environment, ab deployment.AddressBook) (DataFeedsState, error) {
	state := DataFeedsState{
		Chains: make(map[uint64]DataFeedChainState),
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := ab.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainState(chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = chainState
	}
	return state, nil
}
