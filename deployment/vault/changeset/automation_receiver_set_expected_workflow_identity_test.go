package changeset

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/automation-cre/generated/latest/automation_receiver"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const testWorkflowName = "eth-balance-monitor-test"

func TestSetExpectedWorkflowIdentity_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	validChain := func() types.SetExpectedWorkflowIdentityChainConfig {
		return types.SetExpectedWorkflowIdentityChainConfig{
			ExpectedAuthor:       testAddr2,
			ExpectedWorkflowName: testWorkflowName,
		}
	}

	tests := []struct {
		name      string
		cfg       types.SetExpectedWorkflowIdentityInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.SetExpectedWorkflowIdentityInput{Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{}},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "negative MCMS delay",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains:     map[uint64]types.SetExpectedWorkflowIdentityChainConfig{selector: validChain()},
				MCMSConfig: &cldfproposalutils.TimelockConfig{MinDelay: -1},
			},
			wantError: true,
			errorMsg:  "MCMS minimum delay cannot be negative",
		},
		{
			name: "unknown chain selector",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{math.MaxUint64: validChain()},
			},
			wantError: true,
			errorMsg:  "unknown chain selector",
		},
		{
			name: "missing author",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{
					selector: {ExpectedWorkflowName: testWorkflowName},
				},
			},
			wantError: true,
			errorMsg:  "both expectedAuthor and expectedWorkflowName are required",
		},
		{
			name: "missing workflow name",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{
					selector: {ExpectedAuthor: testAddr2},
				},
			},
			wantError: true,
			errorMsg:  "both expectedAuthor and expectedWorkflowName are required",
		},
		{
			name: "invalid author address",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{
					selector: {
						ExpectedAuthor:       "not-an-address",
						ExpectedWorkflowName: testWorkflowName,
					},
				},
			},
			wantError: true,
			errorMsg:  "expectedAuthor is not a valid hex address",
		},
		{
			name: "chain not in environment",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{selectorOther: validChain()},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "valid",
			cfg: types.SetExpectedWorkflowIdentityInput{
				Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{selector: validChain()},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := SetExpectedWorkflowIdentityChangeSet.VerifyPreconditions(*env, tt.cfg)
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

func TestSetExpectedWorkflowIdentityChangeSet(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	// The AutomationReceiver's setCallAllowed requires deployed code on the target
	// (reverts with TargetHasNoCode otherwise), so use the already-deployed timelock
	// as the target when recording the receiver in the datastore.
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

	cfg := types.SetExpectedWorkflowIdentityInput{
		Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{
			selector: {
				ExpectedAuthor:       testAddr2,
				ExpectedWorkflowName: testWorkflowName,
			},
		},
	}
	require.NoError(t, SetExpectedWorkflowIdentityChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	task := runtime.ChangesetTask(SetExpectedWorkflowIdentityChangeSet, cfg)
	require.NoError(t, rt.Exec(task))

	out := rt.State().Outputs[task.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "AutomationReceiver SetExpectedWorkflowIdentity")
	require.Len(t, prop.Operations, 1)
	// Two transactions in the batch: setExpectedAuthor + setExpectedWorkflowName.
	require.Len(t, prop.Operations[0].Transactions, 2)
	require.Contains(t, prop.Operations[0].Transactions[0].Tags, "setExpectedAuthor")
	require.Contains(t, prop.Operations[0].Transactions[1].Tags, "setExpectedWorkflowName")
	require.Equal(t, receiverAddr, prop.Operations[0].Transactions[0].To)
	require.Equal(t, receiverAddr, prop.Operations[0].Transactions[1].To)

	// Decode the generated calldata to confirm it is exactly
	// setExpectedAuthor(expectedAuthor) + setExpectedWorkflowName(expectedWorkflowName).
	arABI, err := automation_receiver.AutomationReceiverMetaData.GetAbi()
	require.NoError(t, err)

	authorTx := prop.Operations[0].Transactions[0]
	authorMethod := arABI.Methods["setExpectedAuthor"]
	require.Equal(t, authorMethod.ID, authorTx.Data[:4])
	authorArgs, err := authorMethod.Inputs.Unpack(authorTx.Data[4:])
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress(testAddr2), authorArgs[0])

	nameTx := prop.Operations[0].Transactions[1]
	nameMethod := arABI.Methods["setExpectedWorkflowName"]
	require.Equal(t, nameMethod.ID, nameTx.Data[:4])
	nameArgs, err := nameMethod.Inputs.Unpack(nameTx.Data[4:])
	require.NoError(t, err)
	require.Equal(t, testWorkflowName, nameArgs[0])
}

func TestSetExpectedWorkflowIdentity_Apply_withoutReceiverInDatastore(t *testing.T) {
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

	cfg := types.SetExpectedWorkflowIdentityInput{
		Chains: map[uint64]types.SetExpectedWorkflowIdentityChainConfig{
			selector: {
				ExpectedAuthor:       testAddr2,
				ExpectedWorkflowName: testWorkflowName,
			},
		},
	}
	require.NoError(t, SetExpectedWorkflowIdentityChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	_, err = SetExpectedWorkflowIdentityChangeSet.Apply(rt.Environment(), cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "AutomationReceiver")
}
