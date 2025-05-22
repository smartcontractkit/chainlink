package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

var (
	DeployOnRampOp = operations.NewOperation(
		"DeployOnRamp",
		semver.MustParse("1.0.0"),
		"Deploys OnRamp 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.OpDependencies, input DeployOnRampInput) (common.Address, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := deps.AddressBook
			chain := e.Chains[input.ChainSelector]
			chainState, chainExists := state.Chains[chain.Selector]
			if !chainExists {
				return common.Address{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", chain.String())
			}
			onRampContract := chainState.OnRamp

			if onRampContract == nil {
				onRamp, err := cldf.DeployContract(b.Logger, chain, ab,
					func(chain cldf.Chain) cldf.ContractDeploy[*onramp.OnRamp] {
						onRampAddr, tx2, onRamp, err2 := onramp.DeployOnRamp(
							chain.DeployerKey,
							chain.Client,
							onramp.OnRampStaticConfig{
								ChainSelector:      chain.Selector,
								RmnRemote:          input.RmnRemote,
								NonceManager:       input.NonceManager,
								TokenAdminRegistry: input.TokenAdminRegistry,
							},
							onramp.OnRampDynamicConfig{
								FeeQuoter:     input.FeeQuoter,
								FeeAggregator: input.FeeAggregator,
							},
							[]onramp.OnRampDestChainConfigArgs{},
						)
						return cldf.ContractDeploy[*onramp.OnRamp]{
							Address: onRampAddr, Contract: onRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0), Err: err2,
						}
					})
				if err != nil {
					b.Logger.Errorw("Failed to deploy onramp", "chain", chain.String(), "err", err)
					return common.Address{}, err
				}
				return onRamp.Address, nil
			}
			b.Logger.Infow("onramp already deployed", "chain", chain.String(), "addr", chainState.OnRamp.Address)
			return common.Address{}, nil
		})
)

type DeployOnRampInput struct {
	ChainSelector      uint64
	TokenAdminRegistry common.Address
	NonceManager       common.Address
	RmnRemote          common.Address
	FeeQuoter          common.Address
	FeeAggregator      common.Address
}
