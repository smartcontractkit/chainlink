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
		func(b operations.Bundle, deps opsutil.OpDependencies, input uint64) (common.Address, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := deps.AddressBook
			chain := e.Chains[input]
			chainState, chainExists := state.Chains[chain.Selector]
			if !chainExists {
				return common.Address{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", chain.String())
			}
			onRampContract := chainState.OnRamp
			if chainState.RMNProxy == nil {
				e.Logger.Errorw("RMNProxy not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("rmn proxy not found for chain %s, deploy the prerequisites first", chain.String())
			}
			if chainState.FeeQuoter == nil {
				e.Logger.Errorw("FeeQuoter not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("fee quoter not found for chain %s, needed for onRamp deployment", chain.String())
			}
			if chainState.NonceManager == nil {
				e.Logger.Errorw("NonceManager not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("nonce manager not found for chain %s, needed for onRamp deployment", chain.String())
			}
			if chainState.TokenAdminRegistry == nil {
				e.Logger.Errorw("TokenAdminRegistry not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("token admin registry not found for chain %s, needed for onRamp deployment", chain.String())
			}
			if onRampContract == nil {
				var feeAggregator common.Address
				// if the fee aggregator is not set, use the deployer key address
				// this is to ensure that feeAggregator is not set to zero address, otherwise there is chance of
				// fund loss when WithdrawFeeToken is called on OnRamp
				if chainState.FeeAggregator != (common.Address{}) {
					feeAggregator = chainState.FeeAggregator
				} else {
					feeAggregator = chain.DeployerKey.From
				}
				onRamp, err := cldf.DeployContract(e.Logger, chain, ab,
					func(chain cldf.Chain) cldf.ContractDeploy[*onramp.OnRamp] {
						onRampAddr, tx2, onRamp, err2 := onramp.DeployOnRamp(
							chain.DeployerKey,
							chain.Client,
							onramp.OnRampStaticConfig{
								ChainSelector:      chain.Selector,
								RmnRemote:          chainState.RMNProxy.Address(),
								NonceManager:       chainState.NonceManager.Address(),
								TokenAdminRegistry: chainState.TokenAdminRegistry.Address(),
							},
							onramp.OnRampDynamicConfig{
								FeeQuoter:     chainState.FeeQuoter.Address(),
								FeeAggregator: feeAggregator,
							},
							[]onramp.OnRampDestChainConfigArgs{},
						)
						return cldf.ContractDeploy[*onramp.OnRamp]{
							Address: onRampAddr, Contract: onRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0), Err: err2,
						}
					})
				if err != nil {
					e.Logger.Errorw("Failed to deploy onramp", "chain", chain.String(), "err", err)
					return common.Address{}, err
				}
				return onRamp.Address, nil
			}
			e.Logger.Infow("onramp already deployed", "chain", chain.String(), "addr", chainState.OnRamp.Address)
			return common.Address{}, nil
		})
)
