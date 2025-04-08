package reward_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
)

func loadRewardManagerState(
	e deployment.Environment,
	chainSel uint64,
	contractAddr string,
) (*rewardManager.RewardManager, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}

	chainState, err := changeset.LoadChainState(e.Logger, chain, addresses)
	if err != nil {
		e.Logger.Errorw("Failed to load chain state", "err", err)
		return nil, err
	}

	conf, found := chainState.RewardManagers[common.HexToAddress(contractAddr)]

	if !found {
		return nil, fmt.Errorf(
			"unable to find RewardManager contract on chain %s (selector %d, address %s)",
			chain.Name(),
			chain.Selector,
			contractAddr,
		)
	}

	return conf, nil
}
