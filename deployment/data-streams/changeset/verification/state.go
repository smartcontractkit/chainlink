package verification

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment"
)

type VerifierProxyState struct {
	VerifierProxy *verifier_proxy_v0_5_0.VerifierProxy
}

func maybeLoadVerifierProxyState(e cldf.Environment, chainSel uint64, contractAddr string) (*VerifierProxyState, error) {
	if err := utils.ValidateContract(e, chainSel, contractAddr, types.VerifierProxy, deployment.Version0_5_0); err != nil {
		return nil, err
	}
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel) // This should never happen due to validation
	}
	vp, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to load VerifierProxy contract on chain %s (chain selector %d): %w", chain.Name(), chain.Selector, err)
	}

	return &VerifierProxyState{
		VerifierProxy: vp,
	}, nil
}
func loadVerifierState(e cldf.Environment, chainSel uint64, contractAddr string) (*verifier_v0_5_0.Verifier, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	if err := utils.ValidateContract(e, chainSel, contractAddr, types.Verifier, deployment.Version0_5_0); err != nil {
		return nil, err
	}

	conf, err := verifier_v0_5_0.NewVerifier(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
