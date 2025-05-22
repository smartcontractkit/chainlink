package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
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
		func(b operations.Bundle, deps opsutil.OpDependencies, input DeployOffRampInput) (common.Address, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := deps.AddressBook
			chain := e.Chains[input.Chain]
			chainState, chainExists := state.Chains[input.Chain]
			if !chainExists {
				return common.Address{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", chain.String())
			}
			if chainState.RMNProxy == nil {
				e.Logger.Errorw("RMNProxy not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("rmn proxy not found for chain %s", chain.String())
			}
			if chainState.FeeQuoter == nil {
				e.Logger.Errorw("FeeQuoter not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("fee quoter not found for chain %s", chain.String())
			}
			if chainState.NonceManager == nil {
				e.Logger.Errorw("NonceManager not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("nonce manager not found for chain %s", chain.String())
			}
			if chainState.TokenAdminRegistry == nil {
				e.Logger.Errorw("TokenAdminRegistry not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("token admin registry not found for chain %s", chain.String())
			}
			offRampContract := chainState.OffRamp
			if offRampContract == nil {
				offRamp, err := cldf.DeployContract(e.Logger, chain, ab,
					func(chain cldf.Chain) cldf.ContractDeploy[*offramp.OffRamp] {
						offRampAddr, tx2, offRamp, err2 := offramp.DeployOffRamp(
							chain.DeployerKey,
							chain.Client,
							offramp.OffRampStaticConfig{
								ChainSelector:        chain.Selector,
								GasForCallExactCheck: input.Params.GasForCallExactCheck,
								RmnRemote:            chainState.RMNProxy.Address(),
								NonceManager:         chainState.NonceManager.Address(),
								TokenAdminRegistry:   chainState.TokenAdminRegistry.Address(),
							},
							offramp.OffRampDynamicConfig{
								FeeQuoter:                               chainState.FeeQuoter.Address(),
								PermissionLessExecutionThresholdSeconds: input.Params.PermissionLessExecutionThresholdSeconds,
								MessageInterceptor:                      input.Params.MessageInterceptor,
							},
							[]offramp.OffRampSourceChainConfigArgs{},
						)
						return cldf.ContractDeploy[*offramp.OffRamp]{
							Address: offRampAddr, Contract: offRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_6_0), Err: err2,
						}
					})
				if err != nil {
					e.Logger.Errorw("Failed to deploy offramp", "chain", chain.String(), "err", err)
					return common.Address{}, err
				}
				return offRamp.Address, nil
			}
			e.Logger.Infow("offramp already deployed", "chain", chain.String(), "addr", chainState.OffRamp.Address)
			return common.Address{}, nil
		})
)

type DeployOffRampInput struct {
	Chain  uint64
	Params OffRampParams
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
