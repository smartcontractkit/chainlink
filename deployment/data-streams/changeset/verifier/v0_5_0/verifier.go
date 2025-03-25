package v0_5_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	verifier_v0_5_0 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

func loadVerifierState(
	e deployment.Environment,
	chainSel uint64,
	contractAddr string,
) (*verifier_v0_5_0.Verifier, error) {
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

	conf, found := chainState.Verifiers[common.HexToAddress(contractAddr)]

	if !found {
		return nil, fmt.Errorf(
			"unable to find Verifier contract on chain %s (selector %d, address %s)",
			chain.Name(),
			chain.Selector,
			contractAddr,
		)
	}

	return conf, nil
}
