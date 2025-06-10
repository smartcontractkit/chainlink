package v1_6

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

var (
	DeployOffRampOp = operations.NewOperation(
		"DeployOffRamp",
		semver.MustParse("1.0.0"),
		"Deploys OffRamp 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input DeployOffRampInput) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain
			offRamp, err := cldf.DeployContract(b.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*offramp.OffRamp] {
					var (
						offRampAddr common.Address
						tx2         *types.Transaction
						offRamp     *offramp.OffRamp
						err2        error
					)
					if chain.IsZkSyncVM {
						offRampAddr, _, offRamp, err2 = offramp.DeployOffRampZk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
							offramp.OffRampStaticConfig{
								ChainSelector:        chain.Selector,
								GasForCallExactCheck: input.Params.GasForCallExactCheck,
								RmnRemote:            input.RmnRemote,
								NonceManager:         input.NonceManager,
								TokenAdminRegistry:   input.TokenAdminRegistry,
							},
							offramp.OffRampDynamicConfig{
								FeeQuoter:                               input.FeeQuoter,
								PermissionLessExecutionThresholdSeconds: input.Params.PermissionLessExecutionThresholdSeconds,
								MessageInterceptor:                      input.Params.MessageInterceptor,
							},
							[]offramp.OffRampSourceChainConfigArgs{},
						)
					} else {
						offRampAddr, tx2, offRamp, err2 = offramp.DeployOffRamp(
							chain.DeployerKey,
							chain.Client,
							offramp.OffRampStaticConfig{
								ChainSelector:        chain.Selector,
								GasForCallExactCheck: input.Params.GasForCallExactCheck,
								RmnRemote:            input.RmnRemote,
								NonceManager:         input.NonceManager,
								TokenAdminRegistry:   input.TokenAdminRegistry,
							},
							offramp.OffRampDynamicConfig{
								FeeQuoter:                               input.FeeQuoter,
								PermissionLessExecutionThresholdSeconds: input.Params.PermissionLessExecutionThresholdSeconds,
								MessageInterceptor:                      input.Params.MessageInterceptor,
							},
							[]offramp.OffRampSourceChainConfigArgs{},
						)
					}
					return cldf.ContractDeploy[*offramp.OffRamp]{
						Address: offRampAddr, Contract: offRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_6_0), Err: err2,
					}
				})
			if err != nil {
				b.Logger.Errorw("Failed to deploy offramp", "chain", chain.String(), "err", err)
				return common.Address{}, err
			}
			return offRamp.Address, nil
		})

	OffRampApplySourceChainConfigUpdatesOp = opsutil.NewEVMCallOperation(
		"OffRampApplySourceChainConfigUpdatesOp",
		semver.MustParse("1.0.0"),
		"Applies updates to source chain configurations stored on the OffRamp contract",
		offramp.OffRampABI,
		shared.OffRamp,
		offramp.NewOffRamp,
		func(offRamp *offramp.OffRamp, opts *bind.TransactOpts, input []offramp.OffRampSourceChainConfigArgs) (*types.Transaction, error) {
			return offRamp.ApplySourceChainConfigUpdates(opts, input)
		},
	)
)

type DeployOffRampInput struct {
	Chain              uint64
	Params             OffRampParams
	FeeQuoter          common.Address
	RmnRemote          common.Address
	NonceManager       common.Address
	TokenAdminRegistry common.Address
}

type OffRampParams struct {
	GasForCallExactCheck                    uint16
	PermissionLessExecutionThresholdSeconds uint32
	MessageInterceptor                      common.Address
}

func (c OffRampParams) Validate(ignoreGasForCallExactCheck bool) error {
	if !ignoreGasForCallExactCheck && c.GasForCallExactCheck == 0 {
		return errors.New("GasForCallExactCheck is 0")
	}
	if c.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("PermissionLessExecutionThresholdSeconds is 0")
	}
	return nil
}

func DefaultOffRampParams() OffRampParams {
	return OffRampParams{
		GasForCallExactCheck:                    uint16(5000),
		PermissionLessExecutionThresholdSeconds: uint32(globals.PermissionLessExecutionThreshold.Seconds()),
	}
}
