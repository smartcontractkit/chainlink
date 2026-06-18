package forwarder_test

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	changeset3 "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestConfigureForwardersSeq(t *testing.T) {
	h, donConfig := setupForwarderTest(t, false)
	env := h.Runtime.Environment()

	b := optest.NewBundle(t)
	deps := forwarder.ConfigureSeqDeps{
		Env: &env,
	}
	input := forwarder.ConfigureSeqInput{
		DON:        donConfig,
		Qualifier:  "test-configure-forwarder",
		MCMSConfig: nil,
		Chains:     map[uint64]struct{}{}, //  Empty means all chains
	}

	// Execute the ConfigureSeq operation directly
	output, err := operations.ExecuteSequence(b, forwarder.ConfigureSeq, deps, input)
	require.NoError(t, err, "ConfigureSeq should execute successfully")
	require.NotNil(t, output, "ConfigureSeq should return output")
	require.NotNil(t, output.Output.Config, "should have configuration")
}

func TestConfigureForwarders(t *testing.T) {
	h, donConfig := setupForwarderTest(t, false)

	// Test the durable pipeline wrapper
	t.Log("Starting configure changeset application...")
	task := runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
		DON:        donConfig,
		Qualifier:  "test-configure-forwarder",
		MCMSConfig: nil, // Not using MCMS for this test
		Chains:     map[uint64]struct{}{h.RegistrySelector: {}},
	})
	err := h.Runtime.Exec(task)
	require.NoError(t, err, "changeset apply failed")
	out := h.Runtime.State().Outputs[task.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	t.Logf("Configure changeset applied successfully")

	// Verify the changeset output
	require.NotNil(t, out.Reports, "reports should be present")
	require.Empty(t, out.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")
}

func TestConfigureForwarders_WithMCMS(t *testing.T) {
	h, donConfig := setupForwarderTest(t, true)
	registryChainSel := h.RegistrySelector

	// Test the durable pipeline wrapper
	t.Log("Starting configure changeset application with MCMS...")
	task := runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
		DON:       donConfig,
		Qualifier: "test-configure-forwarder",
		MCMSConfig: &contracts.MCMSConfig{
			MinDelay: 10 * time.Second,
			TimelockQualifierPerChain: map[uint64]string{
				registryChainSel: "",
			},
		},
		Chains: map[uint64]struct{}{registryChainSel: {}},
	})
	err := h.Runtime.Exec(task)
	require.NoError(t, err, "changeset apply failed")
	out := h.Runtime.State().Outputs[task.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	t.Logf("Configure changeset with MCMS applied successfully")

	// Verify the changeset output
	require.NotNil(t, out.Reports, "reports should be present")
	require.NotEmpty(t, out.MCMSTimelockProposals, "should have MCMS proposals when using MCMS")
}

func TestConfigureForwarders_SpecificChains(t *testing.T) {
	// This test needs a custom setup to deploy to multiple chains first
	h := test.NewTestHarness(t)
	env := h.Runtime.Environment()
	registryChainSel := h.RegistrySelector

	// Get all available chain selectors for multi-chain deployment
	allChains := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		allChains = append(allChains, chainSel)
	}

	// Deploy Keystone Forwarder contracts to ALL chains (unlike the helper which deploys to one)
	b := optest.NewBundle(t)
	deps := forwarder.DeploySequenceDeps{
		Env: &env,
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

	// We need to create a new runtime from the updated environment
	h.Runtime = runtime.NewFromEnvironment(env)

	// Apply the changeset to configure only specific chains
	t.Log("Starting configure changeset application for specific chains...")
	task := runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
		DON:        donConfig,
		Qualifier:  "test-configure-specific-chains",
		MCMSConfig: nil,
		Chains:     specificChains, // Only configure for registry chain
	})
	err = h.Runtime.Exec(task)
	require.NoError(t, err, "changeset apply failed")
	out := h.Runtime.State().Outputs[task.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	t.Logf("Configure changeset for specific chains applied successfully")

	// Verify the changeset output
	require.NotNil(t, out.Reports, "reports should be present")
	require.Empty(t, out.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")
	require.NotEmpty(t, out.Reports, "should have at least one report for the configured chain")
}

func TestConfigureForwarders_SameQualifierAcrossChains(t *testing.T) {
	h := test.NewTestHarness(t)
	env := h.Runtime.Environment()

	allChains := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		allChains = append(allChains, chainSel)
	}

	b := optest.NewBundle(t)
	deps := forwarder.DeploySequenceDeps{
		Env: &env,
	}
	input := forwarder.DeploySequenceInput{
		Targets:   allChains,
		Qualifier: "test-configure-shared-qualifier",
	}

	got, err := operations.ExecuteSequence(b, forwarder.DeploySequence, deps, input)
	require.NoError(t, err)

	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
	require.NotEmpty(t, got.Output.Datastore)

	env.DataStore = got.Output.Datastore

	donConfig := forwarder.DonConfiguration{
		Name:    "testDONSharedQualifier",
		ID:      4,
		F:       1,
		Version: 1,
		NodeIDs: env.NodeIDs,
	}

	// We need to create a new runtime from the updated environment
	h.Runtime = runtime.NewFromEnvironment(env)

	task := runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
		DON:        donConfig,
		Qualifier:  "test-configure-shared-qualifier",
		MCMSConfig: nil,
		Chains:     map[uint64]struct{}{},
	})
	err = h.Runtime.Exec(task)
	require.NoError(t, err, "changeset apply failed")
	out := h.Runtime.State().Outputs[task.ID()]

	require.NotNil(t, out, "changeset output should not be nil")
	require.NotEmpty(t, out.Reports, "should have reports for configured chains")
}

// setupForwarderTest is a helper function to reduce duplication in configure tests
func setupForwarderTest(t *testing.T, enableMCMS bool) (*test.Harness, forwarder.DonConfiguration) {
	hopts := []test.HarnessOpt{}
	if enableMCMS {
		hopts = append(hopts, test.WithMCMS())
	}

	// Setup test environment
	h := test.NewTestHarness(t, hopts...)
	env := h.Runtime.Environment()
	registryChainSel := h.RegistrySelector

	// Deploy Keystone Forwarder contracts to the test chains
	deps := forwarder.DeploySequenceDeps{
		Env: &env,
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

	// We need to create a new runtime from the updated environment
	env.DataStore = ds.Seal()
	h.Runtime = runtime.NewFromEnvironment(env)

	// Transfer ownership to MCMS if enabled
	if enableMCMS {
		// We need to transfer forwarder ownership to MCMS
		t.Log("Transferring forwarder ownership to MCMS...")
		err = h.Runtime.Exec(
			runtime.ChangesetTask(
				cldf.CreateLegacyChangeSet(acceptForwarderOwnershipProposal),
				&changeset3.AcceptAllOwnershipRequest{
					ChainSelector: registryChainSel,
					MinDelay:      0,
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
		)
		require.NoError(t, err, "failed to transfer forwarder ownership to MCMS")
		t.Log("Forwarder ownership transferred to MCMS successfully")
	}

	// Create test DON configuration
	donConfig := forwarder.DonConfiguration{
		Name:    "testDON",
		ID:      1,
		F:       1,
		Version: 1,
		NodeIDs: env.NodeIDs,
	}

	return h, donConfig
}

// acceptForwarderOwnershipProposal is a test-only variant of
// changeset3.AcceptAllOwnershipsProposal scoped to KeystoneForwarder contracts.
//
// This is used to accept ownership of only the KeystoneForwarder contracts to MCMS.
func acceptForwarderOwnershipProposal(e cldf.Environment, req *changeset3.AcceptAllOwnershipRequest) (cldf.ChangesetOutput, error) {
	var forwarders []common.Address
	for _, addr := range e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSelector)) {
		if addr.Type == datastore.ContractType(changeset3.KeystoneForwarder) {
			forwarders = append(forwarders, common.HexToAddress(addr.Address))
		}
	}
	return mcmschangesets.TransferToMCMSWithTimelockV2(e, mcmschangesets.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{req.ChainSelector: forwarders},
		MCMSConfig:       cldfproposalutils.TimelockConfig{MinDelay: req.MinDelay},
	})
}
