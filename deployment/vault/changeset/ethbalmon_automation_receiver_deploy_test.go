package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestDeployEthBalMonWithReceiver_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	// These cases all fail before the MCMS/timelock datastore check (empty chains, chain not in
	// environment, and the input-level forwarder/identity validations), so a plain env is enough.
	tests := []struct {
		name      string
		cfg       types.DeployEthBalMonWithReceiverInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{}},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "chain not in environment",
			cfg: types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
				selectorOther: {ForwarderAddress: testAddr1},
			}},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "invalid forwarder address",
			cfg: types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
				selector: {ForwarderAddress: "not-an-address"},
			}},
			wantError: true,
			errorMsg:  "forwarderAddress is not a valid hex address",
		},
		{
			name: "author without name",
			cfg: types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, ExpectedAuthor: testAddr2},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor and expectedWorkflowName must be set together",
		},
		{
			name: "name without author",
			cfg: types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, ExpectedWorkflowName: testWorkflowName},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor and expectedWorkflowName must be set together",
		},
		{
			name: "invalid author address",
			cfg: types.DeployEthBalMonWithReceiverInput{Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
				selector: {ForwarderAddress: testAddr1, ExpectedAuthor: "not-an-address", ExpectedWorkflowName: testWorkflowName},
			}},
			wantError: true,
			errorMsg:  "expectedAuthor is not a valid hex address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := DeployEthBalMonWithReceiverChangeSet.VerifyPreconditions(*env, tt.cfg)
			require.Error(t, err)
			if tt.errorMsg != "" {
				require.Contains(t, err.Error(), tt.errorMsg)
			}
		})
	}
}

func TestDeployEthBalMonWithReceiver_VerifyPreconditions_validWithMCMS(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)
	setupMCMSInfrastructure(t, rt, []uint64{selector})

	cfg := types.DeployEthBalMonWithReceiverInput{
		Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
			selector: {
				ForwarderAddress:     testAddr1,
				ExpectedAuthor:       testAddr2,
				ExpectedWorkflowName: testWorkflowName,
			},
		},
	}
	require.NoError(t, DeployEthBalMonWithReceiverChangeSet.VerifyPreconditions(rt.Environment(), cfg))
}

// TestDeployEthBalMonWithReceiver_SetsIdentity exercises the optional identity block in the
// combined deploy sequence: when expectedAuthor + expectedWorkflowName are provided, the
// AutomationReceiver's inbound guard must be configured on-chain by the end of the deploy.
func TestDeployEthBalMonWithReceiver_SetsIdentity(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	cfg := types.DeployEthBalMonWithReceiverInput{
		Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
			selector: {
				ForwarderAddress:     testAddr1,
				ExpectedAuthor:       testAddr2,
				ExpectedWorkflowName: testWorkflowName,
			},
		},
	}
	require.NoError(t, DeployEthBalMonWithReceiverChangeSet.VerifyPreconditions(rt.Environment(), cfg))

	out, err := DeployEthBalMonWithReceiverChangeSet.Apply(rt.Environment(), cfg)
	require.NoError(t, err)

	// Both contracts recorded in the datastore.
	receiverAddr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.AutomationReceiverContractType))
	require.NoError(t, err)
	ebmAddr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.EthBalMonContractType))
	require.NoError(t, err)
	require.NotEmpty(t, ebmAddr)

	// Accept-ownership proposal for EthBalMon is produced.
	require.NotEmpty(t, out.MCMSTimelockProposals)

	// Identity guard set on-chain.
	author, name := readReceiverIdentity(t, rt.Environment(), selector, receiverAddr)
	require.Equal(t, common.HexToAddress(testAddr2), author, "expectedAuthor must be set on-chain")
	require.NotEqual(t, [10]byte{}, name, "expectedWorkflowName must be set on-chain")
}

// TestDeployEthBalMonWithReceiver_WithoutIdentity confirms the identity block is skipped when no
// identity is provided; the receiver still deploys but its guard stays unconfigured.
func TestDeployEthBalMonWithReceiver_WithoutIdentity(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	cfg := types.DeployEthBalMonWithReceiverInput{
		Chains: map[uint64]types.DeployEthBalMonWithReceiverChainConfig{
			selector: {ForwarderAddress: testAddr1},
		},
	}
	out, err := DeployEthBalMonWithReceiverChangeSet.Apply(rt.Environment(), cfg)
	require.NoError(t, err)

	receiverAddr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.AutomationReceiverContractType))
	require.NoError(t, err)

	author, name := readReceiverIdentity(t, rt.Environment(), selector, receiverAddr)
	require.Equal(t, common.Address{}, author, "expectedAuthor must remain unset when no identity is provided")
	require.Equal(t, [10]byte{}, name, "expectedWorkflowName must remain unset when no identity is provided")
}
