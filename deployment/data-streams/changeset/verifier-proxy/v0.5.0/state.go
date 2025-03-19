package v0_5_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

type VerifierProxyState struct {
	VerifierProxy *verifier_proxy_v0_5_0.VerifierProxy
}

func maybeLoadVerifierProxyState(e deployment.Environment, chainSel uint64, contractAddr string) (*VerifierProxyState, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, fmt.Errorf("unable to load existing addresses for chain %d: %w", chainSel, err)
	}
	tv, ok := addresses[contractAddr]
	if !ok {
		return nil, fmt.Errorf("unable to find VerifierProxy contract on chain %s (chain selector %d)", chain.Name(), chain.Selector)
	}

	if tv.Type != types.VerifierProxy || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf("unexpected contract type %s for VerifierProxy on chain %s (chain selector %d)", tv, chain.Name(), chain.Selector)
	}

	vp, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to load VerifierProxy contract on chain %s (chain selector %d): %w", chain.Name(), chain.Selector, err)
	}

	return &VerifierProxyState{
		VerifierProxy: vp,
	}, nil
}
