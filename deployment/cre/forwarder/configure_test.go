package forwarder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestConfigureForwardersSeq(t *testing.T) {
	envWrapper, donConfig := setupForwarderTest(t, false)
	env := envWrapper.Env

	b := optest.NewBundle(t)
	deps := forwarder.ConfigureSeqDeps{
		Env: env,
	}
	input := forwarder.ConfigureSeqInput{
		DON:        donConfig,
		MCMSConfig: nil,
		Chains:     map[uint64]struct{}{}, // Empty means all chains
	}

	// Execute the ConfigureSeq operation directly
	output, err := operations.ExecuteSequence(b, forwarder.ConfigureSeq, deps, input)
	require.NoError(t, err, "ConfigureSeq should execute successfully")
	require.NotNil(t, output, "ConfigureSeq should return output")
	require.NotNil(t, output.Output.Config, "should have configuration")
}

func TestConfigureForwarders(t *testing.T) {
	envWrapper, donConfig := setupForwarderTest(t, false)
	env := envWrapper.Env
	registryChainSel := envWrapper.RegistrySelector

	// Test the durable pipeline wrapper
	t.Log("Starting configure changeset application...")
	changesetOutput, err := forwarder.ConfigureForwarders{}.Apply(*env, forwarder.ConfigureSeqInput{
		DON:        donConfig,
		MCMSConfig: nil, // Not using MCMS for this test
		Chains:     map[uint64]struct{}{registryChainSel: {}},
	})
	require.NoError(t, err, "changeset apply failed")
	require.NotNil(t, changesetOutput, "changeset output should not be nil")
	t.Logf("Configure changeset applied successfully")

	// Verify the changeset output
	require.NotNil(t, changesetOutput.Reports, "reports should be present")
	require.Empty(t, changesetOutput.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")
}

func TestConfigureForwarders_SpecificChains(t *testing.T) {
	// This test needs a custom setup to deploy to multiple chains first
	envWrapper := test.SetupEnvV2(t, false)
	env := envWrapper.Env
	registryChainSel := envWrapper.RegistrySelector

	// Get all available chain selectors for multi-chain deployment
	allChains := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		allChains = append(allChains, chainSel)
	}

	// Deploy Keystone Forwarder contracts to ALL chains (unlike the helper which deploys to one)
	b := optest.NewBundle(t)
	deps := forwarder.DeploySequenceDeps{
		Env: env,
	}
	input := forwarder.DeploySequenceInput{
		Targets:   allChains,
		Qualifier: "test-configure-specific-chains",
	}

	got, err := operations.ExecuteSequence(b, forwarder.DeploySequence, deps, input)
	require.NoError(t, err)

	// Check that deployment to all chains succeeded
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
	require.NotEmpty(t, got.Output.Datastore)

	env.DataStore = got.Output.Datastore

	// Create test DON configuration
	donConfig := forwarder.DonConfiguration{
		Name:    "testDONSpecific",
		ID:      3,
		F:       1,
		Version: 1,
		NodeIDs: env.NodeIDs,
	}

	// Configure only for the registry chain (specific chain selection)
	specificChains := map[uint64]struct{}{
		registryChainSel: {},
	}

	// Apply the changeset to configure only specific chains
	t.Log("Starting configure changeset application for specific chains...")
	changesetOutput, err := forwarder.ConfigureForwarders{}.Apply(*env, forwarder.ConfigureSeqInput{
		DON:        donConfig,
		MCMSConfig: nil,
		Chains:     specificChains, // Only configure for registry chain
	})
	require.NoError(t, err, "changeset apply failed")
	require.NotNil(t, changesetOutput, "changeset output should not be nil")
	t.Logf("Configure changeset for specific chains applied successfully")

	// Verify the changeset output
	require.NotNil(t, changesetOutput.Reports, "reports should be present")
	require.Empty(t, changesetOutput.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")
	require.NotEmpty(t, changesetOutput.Reports, "should have at least one report for the configured chain")
}

// setupForwarderTest is a helper function to reduce duplication in configure tests
func setupForwarderTest(t *testing.T, enableMCMS bool) (*test.EnvWrapperV2, forwarder.DonConfiguration) {
	// Setup test environment
	envWrapper := test.SetupEnvV2(t, enableMCMS)
	env := envWrapper.Env
	registryChainSel := envWrapper.RegistrySelector

	// Deploy Keystone Forwarder contracts to the test chains
	b := optest.NewBundle(t)
	deps := forwarder.DeploySequenceDeps{
		Env: env,
	}
	input := forwarder.DeploySequenceInput{
		Targets:   []uint64{registryChainSel},
		Qualifier: "test-configure-forwarder",
	}

	got, err := operations.ExecuteSequence(b, forwarder.DeploySequence, deps, input)
	require.NoError(t, err)

	// Check that the deployment succeeded
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
	require.NotEmpty(t, got.Output.Datastore)

	// Update environment with deployed contracts
	env.DataStore = got.Output.Datastore

	// Create test DON configuration
	donConfig := forwarder.DonConfiguration{
		Name:    "testDON",
		ID:      1,
		F:       1,
		Version: 1,
		NodeIDs: env.NodeIDs,
	}

	return envWrapper, donConfig
}
