package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

var (
	DeployOnRampOp = opsutil.NewEVMDeployOperation(
		"DeployOnRamp",
		semver.MustParse("1.0.0"),
		"Deploys OnRamp 1.6 contract on the specified evm chain",
		cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0),
		opsutil.VMDeployers[DeployOnRampInput]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, input DeployOnRampInput) (common.Address, *types.Transaction, error) {
				addr, tx, _, err := onramp.DeployOnRamp(
					opts,
					backend,
					onramp.OnRampStaticConfig{
						ChainSelector:      input.ChainSelector,
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
				return addr, tx, err
			},
			DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, input DeployOnRampInput) (common.Address, error) {
				addr, _, _, err := onramp.DeployOnRampZk(
					opts,
					client,
					wallet,
					backend,
					onramp.OnRampStaticConfig{
						ChainSelector:      input.ChainSelector,
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
				return addr, err
			},
		},
	)

	OnRampApplyDestChainConfigUpdatesOp = opsutil.NewEVMCallOperation(
		"OnRampApplyDestChainConfigUpdatesOp",
		semver.MustParse("1.0.0"),
		"Applies updates to destination chain configurations stored on the OnRamp contract",
		onramp.OnRampABI,
		shared.OnRamp,
		onramp.NewOnRamp,
		func(onRamp *onramp.OnRamp, opts *bind.TransactOpts, input []onramp.OnRampDestChainConfigArgs) (*types.Transaction, error) {
			return onRamp.ApplyDestChainConfigUpdates(opts, input)
		},
	)
)

type DeployOnRampInput struct {
	ChainSelector      uint64
	TokenAdminRegistry common.Address
	NonceManager       common.Address
	RmnRemote          common.Address
	FeeQuoter          common.Address
	FeeAggregator      common.Address
}
