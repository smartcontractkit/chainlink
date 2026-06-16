package changeset

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestValidateEthBalMonAcceptOwnershipConfig(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		cfg       types.EthBalMonAcceptOwnershipInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.EthBalMonAcceptOwnershipInput{Chains: []uint64{}},
			wantError: true,
			errorMsg:  "no chains provided",
		},
		{
			name: "unknown chain selector",
			cfg: types.EthBalMonAcceptOwnershipInput{
				Chains: []uint64{math.MaxUint64},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "chain not in environment",
			cfg: types.EthBalMonAcceptOwnershipInput{
				Chains: []uint64{selectorOther},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "valid",
			cfg: types.EthBalMonAcceptOwnershipInput{
				Chains: []uint64{selector},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateEthBalMonAcceptOwnershipConfig(t.Context(), *env, tt.cfg)
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEthBalMonAcceptOwnershipChangeset(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	deployCfg := types.DeployEthBalMonInput{
		Chains: map[uint64]types.DeployEthBalMonChainConfig{
			selector: {SetKeeperRegistryAddress: testAddr1},
		},
	}
	deployTask := runtime.ChangesetTask(DeployEthBalMonChangeSet, deployCfg)
	require.NoError(t, rt.Exec(deployTask))

	cfg := types.EthBalMonAcceptOwnershipInput{
		Chains: []uint64{selector},
	}
	require.NoError(t, EthBalMonAcceptOwnership.VerifyPreconditions(rt.Environment(), cfg))

	acceptTask := runtime.ChangesetTask(EthBalMonAcceptOwnership, cfg)
	require.NoError(t, rt.Exec(acceptTask))

	out := rt.State().Outputs[acceptTask.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "EthBalMon acceptOwnership")
	require.Len(t, prop.Operations, 1)
	require.Len(t, prop.Operations[0].Transactions, 1)
	require.Contains(t, prop.Operations[0].Transactions[0].Tags, "acceptOwnership")
}

func TestEthBalMonAcceptOwnership_Apply_withoutEthBalMonInDatastore(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	cfg := types.EthBalMonAcceptOwnershipInput{
		Chains: []uint64{selector},
	}
	require.NoError(t, EthBalMonAcceptOwnership.VerifyPreconditions(rt.Environment(), cfg))

	_, err = EthBalMonAcceptOwnership.Apply(rt.Environment(), cfg)
	require.Error(t, err)
}
