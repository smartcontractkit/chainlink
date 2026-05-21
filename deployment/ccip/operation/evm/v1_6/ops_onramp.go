package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	opsevm "github.com/smartcontractkit/cld-changesets/pkg/family/evm/operations"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

var (
	DeployOnRampOp = opsevm.NewEVMDeployOperation(
		"DeployOnRamp",
		semver.MustParse("1.0.0"),
		"Deploys OnRamp 1.6 contract on the specified evm chain",
		shared.OnRamp,
		onramp.OnRampMetaData,
		&opsevm.ContractOpts{
			Version:          &deployment.Version1_6_0,
			EVMBytecode:      common.FromHex(onramp.OnRampBin),
			ZkSyncVMBytecode: onramp.ZkBytecode,
		},
		func(input DeployOnRampInput) []any {
			return []any{
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
			}
		},
	)

	OnRampApplyDestChainConfigUpdatesOp = opsevm.NewEVMCallOperation(
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
