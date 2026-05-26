package changeset

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestValidateSetKeeperRegistryAddressConfig(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		cfg       types.EthBalMonSetKeeperRegistryAddressInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.EthBalMonSetKeeperRegistryAddressInput{Chains: map[uint64]types.SetKeeperRegistryChainConfig{}},
			wantError: true,
			errorMsg:  "no chains provided",
		},
		{
			name: "unknown chain selector",
			cfg: types.EthBalMonSetKeeperRegistryAddressInput{
				Chains: map[uint64]types.SetKeeperRegistryChainConfig{
					math.MaxUint64: {NewKeeperRegistryAddress: testAddr1},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "chain not in environment",
			cfg: types.EthBalMonSetKeeperRegistryAddressInput{
				Chains: map[uint64]types.SetKeeperRegistryChainConfig{
					selectorOther: {NewKeeperRegistryAddress: testAddr1},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "invalid keeper registry address",
			cfg: types.EthBalMonSetKeeperRegistryAddressInput{
				Chains: map[uint64]types.SetKeeperRegistryChainConfig{
					selector: {NewKeeperRegistryAddress: "not-hex"},
				},
			},
			wantError: true,
			errorMsg:  "new_keeper_registry_address",
		},
		{
			name: "zero keeper registry address",
			cfg: types.EthBalMonSetKeeperRegistryAddressInput{
				Chains: map[uint64]types.SetKeeperRegistryChainConfig{
					selector: {NewKeeperRegistryAddress: zeroAddr},
				},
			},
			wantError: true,
			errorMsg:  "cannot be zero address",
		},
		{
			name: "valid",
			cfg: types.EthBalMonSetKeeperRegistryAddressInput{
				Chains: map[uint64]types.SetKeeperRegistryChainConfig{
					selector: {NewKeeperRegistryAddress: testAddr1},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateSetKeeperRegistryAddressConfig(t.Context(), *env, tt.cfg)
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

func TestSetKeeperRegistrySequence_noChains(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	_, err := operations.ExecuteSequence(b, SetKeeperRegistrySequence, VaultDeps{}, EthBalMonSetKeeperRegistryAddressSequenceInput{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no chains provided")
}

func TestSetKeeperRegistryAddressChangeset(t *testing.T) {
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

	cfg := types.EthBalMonSetKeeperRegistryAddressInput{
		Chains: map[uint64]types.SetKeeperRegistryChainConfig{
			selector: {NewKeeperRegistryAddress: testAddr2},
		},
	}
	require.NoError(t, SetKeeperRegistryAddress.VerifyPreconditions(rt.Environment(), cfg))

	keeperTask := runtime.ChangesetTask(SetKeeperRegistryAddress, cfg)
	require.NoError(t, rt.Exec(keeperTask))

	out := rt.State().Outputs[keeperTask.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "EthBalMon SetKeeperRegistryAddress")
	require.Len(t, prop.Operations, 1)
	require.Len(t, prop.Operations[0].Transactions, 1)
	require.Contains(t, prop.Operations[0].Transactions[0].Tags, "setKeeperRegistryAddress")
}
