package reward_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment"
)

func loadRewardManagerState(
	e cldf.Environment,
	chainSel uint64,
	contractAddr string,
) (*rewardManager.RewardManager, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	if err := utils.ValidateContract(e, chainSel, contractAddr, types.RewardManager, deployment.Version0_5_0); err != nil {
		return nil, err
	}

	conf, err := rewardManager.NewRewardManager(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
