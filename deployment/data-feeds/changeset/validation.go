package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
)

func ValidateCacheForChain(env deployment.Environment, chainSelector uint64, cacheAddress common.Address) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain state")
	}
	if chainState.DataFeedsCache == nil {
		return errors.New("DataFeedsCache not found in on chain state")
	}
	_, ok = chainState.DataFeedsCache[cacheAddress]
	if !ok {
		return errors.New("contract not found in on chain state")
	}
	return nil
}
