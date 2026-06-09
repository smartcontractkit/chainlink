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

func TestValidateEthBalMonTransferOwnershipConfig(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		cfg       types.EthBalMonTransferOwnershipInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.EthBalMonTransferOwnershipInput{Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{}},
			wantError: true,
			errorMsg:  "no chains provided",
		},
		{
			name: "unknown chain selector",
			cfg: types.EthBalMonTransferOwnershipInput{
				Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
					math.MaxUint64: {NewOwner: testAddr1},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "chain not in environment",
			cfg: types.EthBalMonTransferOwnershipInput{
				Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
					selectorOther: {NewOwner: testAddr1},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "invalid new owner",
			cfg: types.EthBalMonTransferOwnershipInput{
				Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
					selector: {NewOwner: "not-hex"},
				},
			},
			wantError: true,
			errorMsg:  "newOwner",
		},
		{
			name: "zero new owner",
			cfg: types.EthBalMonTransferOwnershipInput{
				Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
					selector: {NewOwner: zeroAddr},
				},
			},
			wantError: true,
			errorMsg:  "cannot be zero address",
		},
		{
			name: "valid",
			cfg: types.EthBalMonTransferOwnershipInput{
				Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
					selector: {NewOwner: testAddr1},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateEthBalMonTransferOwnershipConfig(t.Context(), *env, tt.cfg)
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

func TestEthBalMonTransferOwnershipChangeset(t *testing.T) {
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

	cfg := types.EthBalMonTransferOwnershipInput{
		Chains: map[uint64]types.EthBalMonTransferOwnershipChainConfig{
			selector: {NewOwner: testAddr2},
		},
	}
	require.NoError(t, EthBalMonTransferOwnership.VerifyPreconditions(rt.Environment(), cfg))

	transferTask := runtime.ChangesetTask(EthBalMonTransferOwnership, cfg)
	require.NoError(t, rt.Exec(transferTask))

	out := rt.State().Outputs[transferTask.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "EthBalMon transferOwnership")
	require.Len(t, prop.Operations, 1)
	require.Len(t, prop.Operations[0].Transactions, 1)
	require.Contains(t, prop.Operations[0].Transactions[0].Tags, "transferOwnership")
}
