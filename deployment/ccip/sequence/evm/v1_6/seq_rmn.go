package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/maps"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipsharedops "github.com/smartcontractkit/chainlink/deployment/ccip/operation"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	SetRMNRemoteConfigSequence = operations.NewSequence(
		"SetRMNRemoteConfigSequence",
		semver.MustParse("1.0.0"),
		"Set RMNRemoteConfig based on ActiveDigest from RMNHome for evm chain(s)",
		func(b operations.Bundle, deps opsutil.ConfigureDependencies, input SetRMNRemoteConfig) (opsutil.OpOutput, error) {
			finalOutput := &opsutil.OpOutput{}
			err := input.Validate(deps.Env, deps.CurrentState)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to validate input: %w", err)
			}
			for chainSelector, config := range input.RMNRemoteConfigs {
				opInput := ccipops.SetRMNRemoteConfig{
					RMNRemoteConfig: config,
					ChainSelector:   chainSelector,
					MCMSConfig:      input.MCMSConfig,
				}

				report, err := operations.ExecuteOperation(b, ccipops.SetRMNRemoteConfigOp, deps, opInput)
				if err != nil {
					return report.Output, fmt.Errorf("failed to set RMNRemoteConfig for chain %d: %w", chainSelector, err)
				}
				if err := finalOutput.Merge(report.Output); err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to merge output for chain %d: %w", chainSelector, err)
				}
			}
			// if the MCMSConfig is not nil, we need to aggregate the proposals
			if len(finalOutput.Proposals) > 0 {
				report, err := operations.ExecuteOperation(b, ccipsharedops.PostOpsAggregateProposals, deps, ccipsharedops.PostOpsInput{
					MCMSConfig: input.MCMSConfig,
					Proposals:  finalOutput.Proposals,
				})
				if err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to aggregate proposals: %w", err)
				}
				b.Logger.Infow("Generated proposal for RMNRemoteConfig", "chains", maps.Keys(input.RMNRemoteConfigs))
				return opsutil.OpOutput{
					Proposals:                  report.Output,
					DescribedTimelockProposals: finalOutput.DescribedTimelockProposals,
				}, err
			}
			return *finalOutput, nil
		})
)

type SetRMNRemoteConfig struct {
	RMNRemoteConfigs map[uint64]ccipops.RMNRemoteConfig `json:"rmnRemoteConfigs"`
	MCMSConfig       *proposalutils.TimelockConfig      `json:"mcmsConfig,omitempty"`
}

func (c SetRMNRemoteConfig) Validate(env cldf.Environment, state stateview.CCIPOnChainState) error {
	for chainSelector, config := range c.RMNRemoteConfigs {
		err := stateview.ValidateChain(env, state, chainSelector, c.MCMSConfig)
		if err != nil {
			return err
		}
		chain := env.BlockChains.EVMChains()[chainSelector]
		if state.MustGetEVMChainState(chainSelector).RMNRemote == nil {
			return fmt.Errorf("RMNRemote not found for chain %s", chain.String())
		}
		err = commoncs.ValidateOwnership(
			env.GetContext(), c.MCMSConfig != nil,
			chain.DeployerKey.From, state.MustGetEVMChainState(chainSelector).Timelock.Address(),
			state.MustGetEVMChainState(chainSelector).RMNRemote,
		)
		if err != nil {
			return fmt.Errorf("failed to validate ownership for chain %d: %w", chainSelector, err)
		}
		for i := 0; i < len(config.Signers)-1; i++ {
			if config.Signers[i].NodeIndex >= config.Signers[i+1].NodeIndex {
				return fmt.Errorf("signers must be in ascending order of nodeIndex, but found %d >= %d", config.Signers[i].NodeIndex, config.Signers[i+1].NodeIndex)
			}
		}

		//nolint:gosec // G115
		if len(config.Signers) < 2*int(config.F)+1 {
			return fmt.Errorf("signers count (%d) must be greater than or equal to %d", len(config.Signers), 2*config.F+1)
		}
	}

	return nil
}
