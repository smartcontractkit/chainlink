package cre

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	chain_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	job_pkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	cronCapabilityID = "cron-trigger@1.0.0"

	// offchainRegistryEVMDeltaStage matches the local EVM feature (block time + 1s).
	offchainRegistryEVMDeltaStage = 500*time.Millisecond + 1*time.Second
)

// Test_CRE_V2_OffchainCapabilitiesRegistry runs on a topology where the workflow DON has
// [Capabilities.Local] UseOffchainRegistry = true, so the locally launched cron and EVM capabilities take
// their spec_config only from a capabilities_registry cresettings job delivered via JD.
//
// The offchain config lowers cron's fastestScheduleIntervalSeconds from the 30s default to 5s. A
// workflow scheduled every 5s registers and executes only if that config reached the capability. It
// also carries the EVM chain entries: the workflow DON allowlists its EVM chains for registry launch,
// and an allowlisted capability missing from the payload would launch with an empty config.
//
//nolint:paralleltest // the test owns the topology it runs on
func Test_CRE_V2_OffchainCapabilitiesRegistry(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-capabilities-don-offchain-registry.toml"),
	)
	ExecuteOffchainCapabilitiesRegistryTest(t, testEnv)
}

// Test_CRE_V2_OffchainCapabilitiesRegistry_ChainCapabilities runs the same topology and proves the EVM
// chain capability launches from the offchain payload alone. With the gate on, the node ignores the
// onchain spec_config, so an EVM read that succeeds shows the payload's chainId, network and
// creForwarderAddress reached the capability. The entry is built with the production payload-entry
// helper (deployment/cre/jobs), not the local feature's own config builder.
//
//nolint:paralleltest // the test owns the topology it runs on
func Test_CRE_V2_OffchainCapabilitiesRegistry_ChainCapabilities(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-capabilities-don-offchain-registry.toml"),
	)
	ExecuteOffchainChainCapabilitiesTest(t, testEnv)
}

func ExecuteOffchainCapabilitiesRegistryTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowDON := testEnv.Dons.MustWorkflowDON()

	cronConfig, err := jobs.SpecConfigCapability(map[string]any{"fastestScheduleIntervalSeconds": 5})
	require.NoError(t, err)
	donCapabilities := map[string]*capabilitiespb.CapabilityConfig{cronCapabilityID: cronConfig}
	for _, entry := range evmChainCapabilityEntries(t, testEnv, workflowDON) {
		donCapabilities[entry.CapabilityID] = entry.Config
	}

	// Versions must strictly increase across runs against the same environment.
	version := uint64(time.Now().Unix()) //nolint:gosec // G115: positive
	payload, err := jobs.CapabilitiesRegistryPayload(testEnv.CreEnvironment.CldfEnvironment, version, map[string]map[string]*capabilitiespb.CapabilityConfig{
		workflowDON.Name: donCapabilities,
	})
	require.NoError(t, err, "failed to build offchain capabilities registry payload")
	require.NoError(t, jobs.ProposeAndApproveCapabilitiesRegistry(t.Context(), testEnv.CreEnvironment.CldfEnvironment, testEnv.Dons, workflowDON, payload))
	testLogger.Info().Str("don", workflowDON.Name).Msg("Capabilities registry job proposed and approved")

	userLogsCh := make(chan *workfloweventsv2.WorkflowUserLog, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowName := t_helpers.UniqueWorkflowName(testEnv, "offchaincapreg")
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName,
		&crontypes.WorkflowConfig{Schedule: "*/5 * * * * *"},
		"../../../../core/scripts/cre/environment/examples/workflows/cron/main.go",
	)

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	testLogger.Info().Msg("Offchain capabilities registry test completed")
}

// ExecuteOffchainChainCapabilitiesTest delivers a payload with the DON's EVM chain entries and runs an
// EVM read through the registry-launched capability.
func ExecuteOffchainChainCapabilitiesTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowDON := testEnv.Dons.MustWorkflowDON()

	entries := evmChainCapabilityEntries(t, testEnv, workflowDON)
	require.NotEmpty(t, entries, "the workflow DON has no EVM chain enabled")
	allowlisted := make([]string, 0, len(entries))
	chainCapabilities := make(map[string]*capabilitiespb.CapabilityConfig, len(entries)+1)
	for _, entry := range entries {
		allowlisted = append(allowlisted, entry.CapabilityID)
		chainCapabilities[entry.CapabilityID] = entry.Config
	}

	cronConfig, err := jobs.SpecConfigCapability(map[string]any{"fastestScheduleIntervalSeconds": 5})
	require.NoError(t, err)
	version := uint64(time.Now().Unix()) //nolint:gosec // G115: positive

	// An allowlisted chain missing from the payload would launch with an empty config. The completeness
	// check is what catches it before delivery; prove it does for a payload that carries only cron.
	cronOnly, err := jobs.CapabilitiesRegistryPayload(testEnv.CreEnvironment.CldfEnvironment, version, map[string]map[string]*capabilitiespb.CapabilityConfig{
		workflowDON.Name: {cronCapabilityID: cronConfig},
	})
	require.NoError(t, err)
	require.Error(t, job_pkg.CheckAllowlistCompleteness(cronOnly, workflowDON.Name, allowlisted),
		"a payload without the allowlisted EVM chains must be reported incomplete")

	chainCapabilities[cronCapabilityID] = cronConfig
	payload, err := jobs.CapabilitiesRegistryPayload(testEnv.CreEnvironment.CldfEnvironment, version, map[string]map[string]*capabilitiespb.CapabilityConfig{
		workflowDON.Name: chainCapabilities,
	})
	require.NoError(t, err)
	require.NoError(t, job_pkg.CheckAllowlistCompleteness(payload, workflowDON.Name, allowlisted), "the full payload covers every allowlisted chain")

	require.NoError(t, jobs.ProposeAndApproveCapabilitiesRegistry(t.Context(), testEnv.CreEnvironment.CldfEnvironment, testEnv.Dons, workflowDON, payload))
	testLogger.Info().Str("don", workflowDON.Name).Strs("chains", allowlisted).Msg("Chain capabilities delivered through the offchain registry payload")

	ExecuteEVMReadTestForCases(t, testEnv, []evm_config.TestCase{evm_config.TestCaseEVMReadBalance})
}

// evmChainCapabilityEntries builds the offchain payload entry of every EVM chain enabled on don, with the
// production payload-entry helper, the chain selector derived from the chain ID and the CRE forwarder
// address from the environment's datastore.
func evmChainCapabilityEntries(t *testing.T, testEnv *ttypes.TestEnvironment, don *cre.Don) []chain_jobs.ChainCapabilityEntry {
	t.Helper()

	var chainIDs []uint64
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name != don.Name {
			continue
		}
		var err error
		chainIDs, err = nodeSet.GetEnabledChainIDsForCapability(cre.EVMCapability)
		require.NoError(t, err, "failed to get the EVM chains enabled on DON %s", don.Name)
	}

	entries := make([]chain_jobs.ChainCapabilityEntry, 0, len(chainIDs))
	for _, chainID := range chainIDs {
		selector, selErr := chainselectors.SelectorFromChainId(chainID)
		require.NoError(t, selErr, "no chain selector for chain ID %d", chainID)

		forwarder, addrErr := testEnv.CreEnvironment.CldfEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
			selector,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		))
		require.NoError(t, addrErr, fmt.Sprintf("no CRE forwarder address for chain %d", chainID))

		entry, entryErr := chain_jobs.EVMCapabilityEntry(chain_jobs.EVMCapabilityEntryInput{
			ChainSelector:    selector,
			ForwarderAddress: forwarder.Address,
			DeltaStage:       offchainRegistryEVMDeltaStage,
		})
		require.NoError(t, entryErr, "failed to build the EVM payload entry for chain %d", chainID)
		entries = append(entries, entry)
	}
	return entries
}
