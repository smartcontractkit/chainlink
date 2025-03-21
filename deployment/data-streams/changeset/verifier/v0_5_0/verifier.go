package v0_5_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
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

	tv, found := addresses[contractAddr]
	if !found {
		return nil, fmt.Errorf(
			"unable to find Verifier contract on chain %s (selector %d)",
			chain.Name(),
			chain.Selector,
		)
	}

	if tv.Type != types.Verifier || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf(
			"unexpected contract type %s for Verifier on chain %s (selector %d)",
			tv,
			chain.Name(),
			chain.Selector,
		)
	}

	conf, err := verifier_v0_5_0.NewVerifier(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
