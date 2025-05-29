package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
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
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input DeployOnRampInput) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain
			onRamp, err := cldf.DeployContract(b.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*onramp.OnRamp] {
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
