package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/automation-cre/generated/latest/automation_receiver"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// readReceiverIdentity reads the on-chain expectedAuthor + expectedWorkflowName from the
// AutomationReceiver deployed at receiverAddr on the given chain. Shared by the AR deploy
// and combined EthBalMon+AR deploy tests to assert the inbound identity guard was configured.
func readReceiverIdentity(t *testing.T, env cldf.Environment, selector uint64, receiverAddr string) (common.Address, [10]byte) {
	t.Helper()
	chain := env.BlockChains.EVMChains()[selector]
	ar, err := automation_receiver.NewAutomationReceiver(common.HexToAddress(receiverAddr), chain.Client)
	require.NoError(t, err)
	author, err := ar.GetExpectedAuthor(nil)
	require.NoError(t, err)
	name, err := ar.GetExpectedWorkflowName(nil)
	require.NoError(t, err)
	return author, name
}

func TestDeployAutomationReceiver_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		cfg       types.DeployAutomationReceiverInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{}},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "chain not in environment",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selectorOther: {ForwarderAddress: testAddr1, TargetAddress: testAddr2},
			}},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "invalid forwarder address",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: "not-an-address", TargetAddress: testAddr2},
			}},
			wantError: true,
			errorMsg:  "forwarderAddress is not a valid hex address",
		},
		{
			name: "author without name",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, TargetAddress: testAddr2, ExpectedAuthor: testAddr2},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor and expectedWorkflowName must be set together",
		},
		{
			name: "name without author",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, TargetAddress: testAddr2, ExpectedWorkflowName: testWorkflowName},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor and expectedWorkflowName must be set together",
		},
		{
			name: "invalid author address",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, TargetAddress: testAddr2, ExpectedAuthor: "not-an-address", ExpectedWorkflowName: testWorkflowName},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor is not a valid hex address",
		},
		{
			name: "valid without identity",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, TargetAddress: testAddr2},
			}},
			wantError: false,
		},
		{
			name: "valid with identity",
			cfg: types.DeployAutomationReceiverInput{Chains: map[uint64]types.AutomationReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, TargetAddress: testAddr2, ExpectedAuthor: testAddr2, ExpectedWorkflowName: testWorkflowName},
			}},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := DeployAutomationReceiverChangeSet.VerifyPreconditions(*env, tt.cfg)
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

// TestDeployAutomationReceiver_SetsIdentity exercises the SetExpectedWorkflowIdentityOperation
// (deployer-owned direct op) through the standalone AutomationReceiver deploy flow: when
// expectedAuthor + expectedWorkflowName are provided, they must be written on-chain before
// ownership is transferred to the Timelock.
func TestDeployAutomationReceiver_SetsIdentity(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	// setCallAllowed requires the target to have deployed code (reverts with TargetHasNoCode
	// otherwise); the timelock deployed by the MCMS setup is a convenient real contract.
	timelockAddr, err := GetContractAddress(rt.State().DataStore, selector, commontypes.RBACTimelock)
	require.NoError(t, err)

	cfg := types.DeployAutomationReceiverInput{
		Chains: map[uint64]types.AutomationReceiverChainConfig{
			selector: {
				ForwarderAddress:     testAddr1,
				TargetAddress:        timelockAddr,
				ExpectedAuthor:       testAddr2,
				ExpectedWorkflowName: testWorkflowName,
			},
		},
	}
	require.NoError(t, DeployAutomationReceiverChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	out, err := DeployAutomationReceiverChangeSet.Apply(rt.Environment(), cfg)
	require.NoError(t, err)

	receiverAddr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.AutomationReceiverContractType))
	require.NoError(t, err)

	author, name := readReceiverIdentity(t, rt.Environment(), selector, receiverAddr)
	require.Equal(t, common.HexToAddress(testAddr2), author, "expectedAuthor must be set on-chain")
	require.NotEqual(t, [10]byte{}, name, "expectedWorkflowName must be set on-chain")
}

// TestDeployAutomationReceiver_WithoutIdentity_LeavesGuardUnset confirms the op is skipped when
// no identity is provided (the guard stays unconfigured).
func TestDeployAutomationReceiver_WithoutIdentity_LeavesGuardUnset(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	timelockAddr, err := GetContractAddress(rt.State().DataStore, selector, commontypes.RBACTimelock)
	require.NoError(t, err)

	cfg := types.DeployAutomationReceiverInput{
		Chains: map[uint64]types.AutomationReceiverChainConfig{
			selector: {
				ForwarderAddress: testAddr1,
				TargetAddress:    timelockAddr,
			},
		},
	}
	out, err := DeployAutomationReceiverChangeSet.Apply(rt.Environment(), cfg)
	require.NoError(t, err)

	receiverAddr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.AutomationReceiverContractType))
	require.NoError(t, err)

	author, name := readReceiverIdentity(t, rt.Environment(), selector, receiverAddr)
	require.Equal(t, common.Address{}, author, "expectedAuthor must remain unset when no identity is provided")
	require.Equal(t, [10]byte{}, name, "expectedWorkflowName must remain unset when no identity is provided")
}
