package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	SetRMNRemoteConfigSequence = operations.NewSequence(
		"SetRMNRemoteConfigSequence",
		semver.MustParse("1.0.0"),
		"Set RMNRemoteConfig based on ActiveDigest from RMNHome for evm chain(s)",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, inputs map[uint64]opsutil.EVMCallInput[rmn_remote.RMNRemoteConfig]) (map[uint64][]opsutil.EVMCallOutput, error) {
			out := make(map[uint64][]opsutil.EVMCallOutput, len(inputs))

			for chainSelector, input := range inputs {
				if _, ok := chains[chainSelector]; !ok {
					return nil, fmt.Errorf("chain with selector %d not defined in dependencies", chainSelector)
				}

				report, err := operations.ExecuteOperation(b, ccipops.SetRMNRemoteConfigOp, chains[chainSelector], input)
				if err != nil {
					return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to set RMNRemoteConfig for chain %d: %w", chainSelector, err)
				}
				out[chainSelector] = []opsutil.EVMCallOutput{report.Output}
			}

			return out, nil
		})

	SetRMNRemoteOnRMNProxySequence = operations.NewSequence(
		"SetRMNRemoteOnRMNProxySequece",
		semver.MustParse("1.0.0"),
		"Setting SetRMNRemote on RMNProxy across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input SetRMNRemoteOnRMNProxySequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))

			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.SetRMNRemoteOnRMNProxyOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute SetRMNRemoteOnRMNProxyOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)

type SetRMNRemoteOnRMNProxySequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[common.Address] `json:"updatesByChain"`
}

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
