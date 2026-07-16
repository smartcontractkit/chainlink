package changeset

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const testSelectorHex = "0x4b9f5c20"

func TestSetCallAllowed_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	validChain := func() types.SetCallAllowedChainConfig {
		return types.SetCallAllowedChainConfig{
			TargetAddress: testAddr2,
			Selector:      testSelectorHex,
			Allowed:       true,
		}
	}

	tests := []struct {
		name      string
		cfg       types.SetCallAllowedInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.SetCallAllowedInput{Chains: map[uint64]types.SetCallAllowedChainConfig{}},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "negative MCMS delay",
			cfg: types.SetCallAllowedInput{
				Chains:     map[uint64]types.SetCallAllowedChainConfig{selector: validChain()},
				MCMSConfig: &cldfproposalutils.TimelockConfig{MinDelay: -1},
			},
			wantError: true,
			errorMsg:  "MCMS minimum delay cannot be negative",
		},
		{
			name: "unknown chain selector",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{math.MaxUint64: validChain()},
			},
			wantError: true,
			errorMsg:  "unknown chain selector",
		},
		{
			name: "chain not in environment",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{selectorOther: validChain()},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "invalid target address",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{
					selector: {
						TargetAddress: "not-an-address",
						Selector:      testSelectorHex,
					},
				},
			},
			wantError: true,
			errorMsg:  "targetAddress is not a valid hex address",
		},
		{
			name: "invalid selector - not hex",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{
					selector: {
						TargetAddress: testAddr2,
						Selector:      "0xzzzz",
					},
				},
			},
			wantError: true,
			errorMsg:  "invalid selector",
		},
		{
			name: "invalid selector - wrong length",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{
					selector: {
						TargetAddress: testAddr2,
						Selector:      "0x1234",
					},
				},
			},
			wantError: true,
			errorMsg:  "selector must be exactly 4 bytes",
		},
		{
			name: "valid",
			cfg: types.SetCallAllowedInput{
				Chains: map[uint64]types.SetCallAllowedChainConfig{selector: validChain()},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := SetCallAllowedChangeSet.VerifyPreconditions(*env, tt.cfg)
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

func TestSetCallAllowedChangeSet(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	// The AutomationReceiver validates that setCallAllowed targets have deployed
	// code (reverts with TargetHasNoCode otherwise), so use a real contract
	// address as the target. The timelock is already deployed by the MCMS setup.
	timelockAddr, err := GetContractAddress(rt.State().DataStore, selector, commontypes.RBACTimelock)
	require.NoError(t, err)

	// Deploy the AutomationReceiver so its address is recorded in the datastore.
	deployCfg := types.DeployAutomationReceiverInput{
		Chains: map[uint64]types.AutomationReceiverChainConfig{
			selector: {
				ForwarderAddress: testAddr1,
				TargetAddress:    timelockAddr,
				Selector:         testSelectorHex,
			},
		},
	}
	require.NoError(t, rt.Exec(runtime.ChangesetTask(DeployAutomationReceiverChangeSet, deployCfg)))

	receiverAddr, err := GetContractAddress(rt.State().DataStore, selector, cldf.ContractType(types.AutomationReceiverContractType))
	require.NoError(t, err)

	cfg := types.SetCallAllowedInput{
		Chains: map[uint64]types.SetCallAllowedChainConfig{
			selector: {
				TargetAddress: testAddr2,
				Selector:      testSelectorHex,
				Allowed:       true,
			},
		},
	}
	require.NoError(t, SetCallAllowedChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	task := runtime.ChangesetTask(SetCallAllowedChangeSet, cfg)
	require.NoError(t, rt.Exec(task))

	out := rt.State().Outputs[task.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "AutomationReceiver SetCallAllowed")
	require.Len(t, prop.Operations, 1)
	require.Len(t, prop.Operations[0].Transactions, 1)
	tx := prop.Operations[0].Transactions[0]
	require.Contains(t, tx.Tags, "setCallAllowed")
	require.Equal(t, receiverAddr, tx.To)
}

func TestSetCallAllowed_VerifyPreconditions_rejectsEmptyChains(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	err = SetCallAllowedChangeSet.VerifyPreconditions(rt.Environment(), types.SetCallAllowedInput{
		Chains: map[uint64]types.SetCallAllowedChainConfig{},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "chains must not be empty")
}

func TestSetCallAllowed_Apply_withoutReceiverInDatastore(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	// MCMS infra is present, but the AutomationReceiver was never deployed, so its
	// address cannot be resolved from the datastore.
	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	cfg := types.SetCallAllowedInput{
		Chains: map[uint64]types.SetCallAllowedChainConfig{
			selector: {
				TargetAddress: testAddr2,
				Selector:      testSelectorHex,
				Allowed:       true,
			},
		},
	}
	require.NoError(t, SetCallAllowedChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	_, err = SetCallAllowedChangeSet.Apply(rt.Environment(), cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "AutomationReceiver")
}
