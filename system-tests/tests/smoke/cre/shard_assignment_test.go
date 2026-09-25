package cre

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	ring_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	stvault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteManualShardAssignmentTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for manual assignment test")

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"

	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "org_test_manual")

	shardLeaderDON := getShardZeroDon(t, testEnv)
	shardZeroDonID := uint32(shardLeaderDON.ID) //nolint:gosec // G115: overflow is unrealistic

	var shardOneDON *cre.Don
	for _, don := range shardDONs {
		if don.ID != shardLeaderDON.ID {
			shardOneDON = don
			break
		}
	}
	require.NotNil(t, shardOneDON, "Expected to find a second shard DON")
	shardOneDonID := uint32(shardOneDON.ID) //nolint:gosec // G115: overflow is unrealistic

	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_org_assignment]
  org_test_manual = [%d]
`, shardOneDonID, shardZeroDonID)

	proposeAndApproveShardAssignmentJob(t, testEnv, shardLeaderDON, shardAssignmentTOML, testLogger)

	const numWorkflows = 5
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowName := fmt.Sprintf("manualshard%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed workflows for manual shard assignment test")

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = shardZeroDonID
	}

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 3 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on correct shards")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on correct shards (manual-only mode)")
}

func ExecuteManualShardAssignmentBothSpecs(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for manual assignment test")

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"

	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "test_org_manual")

	shardLeaderDON := getShardZeroDon(t, testEnv)
	shardZeroDonID := uint32(shardLeaderDON.ID) //nolint:gosec // G115: overflow is unrealistic

	var shardOneDON *cre.Don
	for _, don := range shardDONs {
		if don.ID != shardLeaderDON.ID {
			shardOneDON = don
			break
		}
	}
	require.NotNil(t, shardOneDON, "Expected to find a second shard DON")
	shardOneDonID := uint32(shardOneDON.ID) //nolint:gosec // G115: overflow is unrealistic

	// The static default points at shard zero, so the workflow only lands on shard one if the
	// per-org entry is what routed it.
	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_org_assignment]
  non_existing_org= [%d]
`, shardZeroDonID, shardOneDonID)

	// Every shard resolves ownership from its own copy of the spec, so both the shard that must
	// run the workflow and the shard that must not need it.
	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, testLogger)
	}

	// The CRE limits job and the shard assignment job are both cresettings jobs; a node keeps one
	// slot per config_type, so the two must run side by side. Propose a limits job for the org the
	// workflows run under and check both are active before relying on the assignment below.
	limits := t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Org("test_org_manual", `
[PerOrg]
WorkflowExecutionConcurrencyLimit = '42'`))
	for _, don := range shardDONs {
		requireCRESettingsAndShardAssignmentJobsActive(t, don, limits.AppliedHash(don.Name))
	}

	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "manualshard0", &workflowConfig, workflowFileLocation)
	workflowIDs := []string{workflowID}
	workflowToShardIndex := map[string]uint32{workflowID: shardZeroDonID}
	testLogger.Info().Str("workflowID", workflowID).Msg("Deployed workflow for manual shard assignment test")

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 3 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on correct shards")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on correct shards (manual-only mode)")
}

// ExecuteManualShardAssignmentWithEVMLogTriggerTest exercises manual-only shard assignment on a
// topology where the same capability is hosted by several capability DON shards and routed by DON
// family (configs/workflow-sharded-capabilities-don.toml).
//
// It differs from ExecuteManualShardAssignmentTest in three ways:
//   - the workflow is pinned with per_owner_assignment instead of per_org_assignment, so no linking
//     service is involved and the highest-precedence branch of the manual resolver is covered;
//   - the workflow is driven by an EVM log trigger instead of cron, and the workflow shards host no
//     EVM capability of their own, so the trigger must be served by a remote capability DON;
//   - it asserts which capability DON served that trigger: only the capability shard sharing a DON
//     family with the workflow shard may ack the trigger event, the other one must stay idle.
//
// There is no Ring OCR here: SetupSharding is never called, so no ShardConfig contract, no Ring
// OCR3 contract and no Ring jobs exist. Nodes fall back to the default "manual-only" assignment
// mode (core/services/chainlink/config_sharding.go), which resolves ownership purely from the
// shard-assignment job spec.
func ExecuteManualShardAssignmentWithEVMLogTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for manual assignment log trigger test")

	shardLeaderDON := getShardZeroDon(t, testEnv)
	shardLeaderDonID := uint32(shardLeaderDON.ID) //nolint:gosec // G115: overflow is unrealistic

	// Pin to a non-leader shard on purpose. The leader is where the static default below sends
	// anything the per-owner entry does not match, so an assignment that silently did not take
	// effect cannot look like a pass.
	nonLeaderDONs := slices.DeleteFunc(slices.Clone(shardDONs), func(don *cre.Don) bool {
		return don.ID == shardLeaderDON.ID
	})
	require.NotEmpty(t, nonLeaderDONs, "Expected to find a non-leader shard DON")
	targetDON := nonLeaderDONs[0]
	targetDonID := uint32(targetDON.ID) //nolint:gosec // G115: overflow is unrealistic

	// Both capability shards host the same EVM chain and the workflow shards host none, so the
	// only thing that can decide who serves the log trigger is the DON family the workflow shard
	// shares with one of them. servingCapDON is that shard; idleCapDONs must never be involved.
	targetNodeSet := nodeSetForDON(t, testEnv, targetDON)
	servingCapDON, idleCapDONs := evmCapabilityDONsByFamily(t, testEnv, targetNodeSet)
	logTriggerChain := evmChainEnabledOnNodeSet(t, testEnv, servingCapDON)
	logTriggerChainID := logTriggerChain.CtfOutput().ChainID

	workflowOwner := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.MustGetRootPrivateKey()
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	testLogger.Info().
		Str("shardLeaderDON", shardLeaderDON.Name).
		Uint32("shardLeaderDonID", shardLeaderDonID).
		Str("targetDON", targetDON.Name).
		Uint32("targetDonID", targetDonID).
		Str("logTriggerChainID", logTriggerChainID).
		Str("workflowOwner", workflowOwnerAddress).
		Str("servingCapDON", servingCapDON.Name).
		Strs("idleCapDONs", nodeSetNames(idleCapDONs)).
		Msg("Manual shard assignment with EVM log trigger")

	workflowConfig, msgEmitter := configureEVMLogTriggerWorkflow(t, testLogger, logTriggerChain)
	expectedUserLog := "Data for manual shard assignment log trigger chain " + logTriggerChainID

	emitCtx, emitCancelFn := context.WithCancel(t.Context())
	defer emitCancelFn()
	startEVMLogTriggerEventEmitter(emitCtx, t, testLogger, logTriggerChainID, logTriggerChain, msgEmitter, expectedUserLog)

	// per_owner_assignment takes precedence over both per_org_assignment and the static default,
	// so the static default points at the leader to prove the per-owner entry is what routed.
	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, shardLeaderDonID, workflowOwnerAddress, targetDonID)

	// Every shard resolves ownership from its own copy of the spec, so both the shard that must
	// run the workflow and the shard that must not need it.
	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, testLogger)
	}

	// One workflow is enough: manual resolution keys off the owner alone, ignoring the workflow ID
	// (resolveManual in core/services/workflows/shardownership/resolver.go), so extra workflows
	// would re-evaluate the same per_owner_assignment branch at the cost of another WASM compile.
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "manualshard-evmlogtrigger", &workflowConfig, "./evm/logtrigger/main.go")
	workflowIDs := []string{workflowID}
	workflowToShardIndex := map[string]uint32{workflowID: targetDonID}
	testLogger.Info().Str("workflowID", workflowID).Msg("Deployed workflow for manual shard assignment log trigger test")

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 4 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	// waitForAllWorkflowsExecuted only counts a workflow once it is seen on its expected shard,
	// so a log from the leader shard is reported as a mismatch rather than accepted.
	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Workflow did not execute on the manually assigned shard")
	testLogger.Info().Msg("EVM log trigger workflow executed on the manually assigned shard")

	// The ack is logged by the node that hosts the trigger capability, so the container it comes
	// from names the capability DON that served the trigger. The ack can trail the user log, hence
	// the retry on the positive side; by the time it lands, an ack on the out-of-family shard would
	// already be in its logs.
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, servingCapDON.Name, triggerEventACKLogNeedle, 2*time.Minute, 5*time.Second)
	for _, idle := range idleCapDONs {
		t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, idle.Name, triggerEventACKLogNeedle)
	}
	testLogger.Info().Str("servingCapDON", servingCapDON.Name).Msg("Log trigger was served by the in-family capability DON only")
}

// nodeSetForDON returns the node set that backs don. Node sets carry the DON family and
// chain-capability information the Don snapshot does not.
func nodeSetForDON(t *testing.T, testEnv *ttypes.TestEnvironment, don *cre.Don) *cre.NodeSet {
	t.Helper()

	for _, ns := range testEnv.Config.NodeSets {
		if ns.Name == don.Name {
			return ns
		}
	}

	require.FailNowf(t, "node set not found", "failed to find the node set for DON %s", don.Name)
	return nil
}

// evmCapabilityDONsByFamily splits the capability DONs that host an EVM capability into the single
// one sharing a DON family with workflowNodeSet - the one the launcher's family-overlap filter must
// route to - and the rest, which must never serve that workflow's triggers.
func evmCapabilityDONsByFamily(t *testing.T, testEnv *ttypes.TestEnvironment, workflowNodeSet *cre.NodeSet) (inFamily *cre.NodeSet, outOfFamily []*cre.NodeSet) {
	t.Helper()

	for _, ns := range testEnv.Config.NodeSets {
		if !slices.Contains(ns.DONTypes, cre.CapabilitiesDON) {
			continue
		}

		enabledChainIDs, err := ns.GetEnabledChainIDsForCapability(cre.EVMCapability)
		require.NoErrorf(t, err, "failed to get EVM-enabled chain IDs for DON %s", ns.Name)
		if len(enabledChainIDs) == 0 {
			continue
		}

		if sharesDonFamily(workflowNodeSet, ns) {
			require.Nilf(t, inFamily, "DONs %s and %s both host EVM in a family of %s, routing would be ambiguous", ns.Name, nodeSetName(inFamily), workflowNodeSet.Name)
			inFamily = ns
			continue
		}
		outOfFamily = append(outOfFamily, ns)
	}

	require.NotNilf(t, inFamily, "no capability DON hosts an EVM capability in a DON family of %s", workflowNodeSet.Name)
	require.NotEmptyf(t, outOfFamily, "every EVM capability DON shares a family with %s, so routing cannot be told apart from a broadcast", workflowNodeSet.Name)
	return inFamily, outOfFamily
}

func sharesDonFamily(a, b *cre.NodeSet) bool {
	return slices.ContainsFunc(a.DonFamilies, func(family string) bool {
		return slices.Contains(b.DonFamilies, family)
	})
}

func nodeSetName(ns *cre.NodeSet) string {
	if ns == nil {
		return ""
	}
	return ns.Name
}

func nodeSetNames(nodeSets []*cre.NodeSet) []string {
	names := make([]string, 0, len(nodeSets))
	for _, ns := range nodeSets {
		names = append(names, ns.Name)
	}
	return names
}

// evmChainEnabledOnNodeSet returns a deployed blockchain whose EVM capability is enabled on the
// given node set. It reads the node set rather than Don.GetEnabledChainIDsForCapability, because
// the node set builds its chain-capability index on demand while the Don only carries a snapshot.
func evmChainEnabledOnNodeSet(t *testing.T, testEnv *ttypes.TestEnvironment, nodeSet *cre.NodeSet) blockchains.Blockchain {
	t.Helper()

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(cre.EVMCapability)
	require.NoErrorf(t, err, "failed to get EVM-enabled chain IDs for DON %s", nodeSet.Name)
	require.NotEmptyf(t, enabledChainIDs, "DON %s has no EVM chain enabled, it cannot serve a log trigger", nodeSet.Name)

	for _, chainID := range enabledChainIDs {
		for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
			if bcOutput.ChainID() == chainID {
				return bcOutput
			}
		}
	}

	require.FailNowf(t, "no deployed blockchain matches the EVM chains enabled on the DON",
		"DON %s has EVM chains %v enabled, none of which is deployed in this environment", nodeSet.Name, enabledChainIDs)
	return nil
}

func ExecuteRingOCROverridesTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for override test")

	var shardZero *cre.Don
	for _, don := range shardDONs {
		if don.Metadata().IsShardLeader() {
			shardZero = don
			break
		}
	}
	require.NotNil(t, shardZero, "Expected to find shard zero DON")

	var shardOne *cre.Don
	for _, don := range shardDONs {
		if don.ID != shardZero.ID {
			shardOne = don
			break
		}
	}
	require.NotNil(t, shardOne, "Expected to find a second shard DON")
	shardZeroDonID := uint32(shardZero.ID) //nolint:gosec // G115: overflow is unrealistic
	shardOneDonID := uint32(shardOne.ID)   //nolint:gosec // G115: overflow is unrealistic

	topology, tErr := cre.NewTopology(testEnv.Config.NodeSets, *testEnv.Config.Infra, testEnv.Config.CapabilityConfigs)
	require.NoError(t, tErr, "Failed to recreate topology")

	err := sharding.SetupSharding(t.Context(), sharding.SetupShardingInput{
		Logger:   testLogger,
		CreEnv:   testEnv.CreEnvironment,
		Topology: topology,
		Dons:     testEnv.Dons,
	})
	if err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Ring jobs already exist (from previous run), continuing...")
		} else {
			require.NoError(t, err, "SetupSharding failed")
		}
	} else {
		testLogger.Info().Msg("SetupSharding completed successfully")
	}

	waitForRingOracleHealthy(t, shardZero)

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "org_test_override")

	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = true

[per_org_assignment]
  org_test_override = [%d]
`, shardZeroDonID, shardOneDonID)

	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, testLogger)
	}

	const numWorkflows = 5
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowName := fmt.Sprintf("override-shard%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed workflows for ringocr-with-overrides test")

	var rpcHost string
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name == "shard0" && nodeSet.Out != nil && len(nodeSet.Out.CLNodes) > 0 {
			externalURL := nodeSet.Out.CLNodes[0].Node.ExternalURL
			parsedURL, parseErr := url.Parse(externalURL)
			require.NoError(t, parseErr, "Failed to parse ExternalURL")
			rpcHost = parsedURL.Hostname()
			break
		}
	}
	require.NotEmpty(t, rpcHost, "Failed to find shard0 node set to extract RPC host")

	shardOrchClient := newShardOrchestratorClient(t, rpcHost+":60051")

	testLogger.Info().Msg("Reporting shard status to ALL nodes' Arbiters...")
	initializeAllArbiterStates(t, testEnv, shardZero, len(shardDONs))

	testLogger.Info().Msg("Diagnostic: Verifying store connection (direct registration)...")
	verifyStoreConnection(t, shardOrchClient)

	testLogger.Info().Msg("Diagnostic: Verifying Ring OCR rounds are completing...")
	waitForRingOCRRounds(t, shardOrchClient)

	testLogger.Info().Msg("Waiting for workflows to be registered via Ring OCR...")
	waitForWorkflowsRegistered(t, shardOrchClient, workflowIDs)

	resp, err := shardOrchClient.GetWorkflowShardMapping(t.Context(), &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	testLogger.Info().Interface("mappings", resp.Mappings).Msg("Ring OCR workflow mappings")
	require.Len(t, resp.Mappings, len(workflowIDs), "All deployed workflows should be mapped")

	overrideShard := shardOneDonID

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = overrideShard
	}

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 5 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on correct shards")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on correct shards (ringocr-with-overrides mode)")
}

func proposeAndApproveShardAssignmentJob(t *testing.T, testEnv *ttypes.TestEnvironment, targetDON *cre.Don, shardAssignmentTOML string, testLogger zerolog.Logger) {
	t.Helper()

	jobInput := ring_ops.ProposeShardAssignmentJobInput{
		Domain:          offchain.ProductLabel,
		Environment:     testEnv.CreEnvironment.CldfEnvironment.Name,
		DONName:         targetDON.Name,
		ShardAssignment: shardAssignmentTOML,
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: targetDON.Name},
		},
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: "shard-assignment"},
	}

	report, err := operations.ExecuteOperation(
		testEnv.CreEnvironment.CldfEnvironment.OperationsBundle,
		ring_ops.ProposeShardAssignmentJob,
		ring_ops.ProposeShardAssignmentJobDeps{Env: *testEnv.CreEnvironment.CldfEnvironment},
		jobInput,
	)
	if err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Shard assignment job already exists (from previous run), continuing...")
			return
		}
		require.NoError(t, err, "Failed to propose shard assignment job")
	}

	if err := jobs.Approve(t.Context(), testEnv.CreEnvironment.CldfEnvironment.Offchain, testEnv.Dons, report.Output.Specs); err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Shard assignment job already approved (from previous run), continuing...")
			return
		}
		require.NoError(t, err, "Failed to approve shard assignment job")
	}

	testLogger.Info().Msg("Shard assignment job proposed and approved")
}

// requireCRESettingsAndShardAssignmentJobsActive asserts that every worker node of don has two
// distinct cresettings jobs approved and running: the shard assignment job and the CRE limits job
// whose spec carries settingsHash.
func requireCRESettingsAndShardAssignmentJobsActive(t *testing.T, don *cre.Don, settingsHash string) {
	t.Helper()
	require.NotEmpty(t, settingsHash, "no CRE settings were applied to DON %q", don.Name)

	workers, err := don.Workers()
	require.NoError(t, err, "failed to get worker nodes of DON %q", don.Name)

	for _, node := range workers {
		require.EventuallyWithTf(t, func(c *assert.CollectT) {
			jd, jdErr := node.Clients.GQLClient.GetJobDistributor(t.Context(), node.JobDistributorDetails.JDID)
			if !assert.NoError(c, jdErr, "failed to get job distributor") {
				return
			}

			// The node resolves a proposal's jobID from its jobs table by external job ID, so a
			// non-empty jobID means the job was actually created. ListJobs can't be used here: the
			// GQL client can't unmarshal cresettings job specs.
			var limitsJobID, shardJobID string
			for _, proposal := range jd.JobProposals {
				spec := proposal.LatestSpec
				if !strings.Contains(spec.Definition, `type = "cresettings"`) || string(spec.Status) != "APPROVED" {
					continue
				}
				switch {
				case strings.Contains(spec.Definition, `config_type = "shard_assignment"`):
					shardJobID = proposal.JobID
				case strings.Contains(spec.Definition, fmt.Sprintf("hash = %q", settingsHash)):
					limitsJobID = proposal.JobID
				}
			}
			if !assert.NotEmpty(c, shardJobID, "no running job for an approved shard assignment proposal") ||
				!assert.NotEmpty(c, limitsJobID, "no running job for an approved CRE limits proposal with hash %s", settingsHash) {
				return
			}
			assert.NotEqual(c, shardJobID, limitsJobID, "limits and shard assignment jobs must be distinct")
		}, 2*time.Minute, 3*time.Second, "node %s of DON %q does not run both CRE limits and shard assignment jobs", node.Name, don.Name)
	}

	framework.L.Info().Str("don", don.Name).Int("nodes", len(workers)).Msg("CRE limits and shard assignment jobs are active on all worker nodes")
}
