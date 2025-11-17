package v1_6

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

var (
	DeployOffRampOp = opsutil.NewEVMDeployOperation(
		"DeployOffRamp",
		semver.MustParse("1.0.0"),
		"Deploys OffRamp 1.6 contract on the specified evm chain",
		shared.OffRamp,
		offramp.OffRampMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_6_0,
			EVMBytecode:      common.FromHex(offramp.OffRampBin),
			ZkSyncVMBytecode: offramp.ZkBytecode,
		},
		func(input DeployOffRampInput) []any {
			return []any{
				offramp.OffRampStaticConfig{
					ChainSelector:        input.Chain,
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
			}
		},
	)

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
