package forwarder_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	changeset3 "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
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

func TestConfigureForwarders_WithMCMS(t *testing.T) {
	envWrapper, donConfig := setupForwarderTest(t, true)
	env := envWrapper.Env
	registryChainSel := envWrapper.RegistrySelector

	// Test the durable pipeline wrapper
	t.Log("Starting configure changeset application with MCMS...")
	changesetOutput, err := forwarder.ConfigureForwarders{}.Apply(*env, forwarder.ConfigureSeqInput{
		DON: donConfig,
		MCMSConfig: &contracts.MCMSConfig{
			MinDelay: 10 * time.Second,
			TimelockQualifierPerChain: map[uint64]string{
				registryChainSel: "",
			},
		},
		Chains: map[uint64]struct{}{registryChainSel: {}},
	})
	require.NoError(t, err, "changeset with MCMS apply failed")
	require.NotNil(t, changesetOutput, "changeset output with MCMS should not be nil")
	t.Logf("Configure changeset with MCMS applied successfully")

	// Verify the changeset output
	require.NotNil(t, changesetOutput.Reports, "reports should be present")
	require.NotEmpty(t, changesetOutput.MCMSTimelockProposals, "should have MCMS proposals when using MCMS")
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
	deps := forwarder.DeploySequenceDeps{
		Env: env,
	}
	input := forwarder.DeploySequenceInput{
		Targets:   []uint64{registryChainSel},
		Qualifier: "test-configure-forwarder",
	}

	got, err := operations.ExecuteSequence(env.OperationsBundle, forwarder.DeploySequence, deps, input)
	require.NoError(t, err)

	// Check that the deployment succeeded
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
	require.NotEmpty(t, got.Output.Datastore)

	ds := datastore.NewMemoryDataStore()

	prevDS := env.DataStore
	require.NoError(t, ds.Merge(prevDS), "failed to merge existing datastore")
	require.NoError(t, ds.Merge(got.Output.Datastore), "failed to merge output datastore")

	// Try and transfer ownership to MCMS if enabled
	if enableMCMS {
		env.DataStore = got.Output.Datastore // temporary override to perform ownership transfer of only the forwarder to MCMS
		// We need to transfer forwarder ownership to MCMS
		t.Log("Transferring forwarder ownership to MCMS...")
		resultEnv, mcmsErr := changeset.Apply(t, *env, changeset.Configure(
			cldf.CreateLegacyChangeSet(changeset3.AcceptAllOwnershipsProposal),
			&changeset3.AcceptAllOwnershipRequest{
				ChainSelector: registryChainSel,
				MinDelay:      0,
			},
		))
		require.NoError(t, mcmsErr, "failed to transfer forwarder ownership to MCMS")
		t.Log("Forwarder ownership transferred to MCMS successfully")

		require.NoError(t, ds.Merge(env.DataStore), "failed to merge existing datastore")
		require.NoError(t, ds.Merge(resultEnv.DataStore), "failed to merge existing datastore")
	}

	envWrapper.Env.DataStore = ds.Seal()

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
