package fee_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func LoadFeeManagerState(
	e deployment.Environment,
	chainSel uint64,
	contractAddr string,
) (*fee_manager_v0_5_0.FeeManager, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}

	tv, found := addresses[contractAddr]
	if !found {
		return nil, fmt.Errorf(
			"unable to find FeeManager contract on chain %s (selector %d)",
			chain.Name(),
			chain.Selector,
		)
	}

	if tv.Type != types.FeeManager || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf(
			"unexpected contract type %s for Verifier on chain %s (selector %d)",
			tv,
			chain.Name(),
			chain.Selector,
		)
	}

	conf, err := fee_manager_v0_5_0.NewFeeManager(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
