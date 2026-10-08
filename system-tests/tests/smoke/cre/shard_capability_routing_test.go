package cre

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

/*
Don2Don sharded-capability routing scenarios on the sharded-capabilities topology
(configs/workflow-sharded-capabilities-don.toml):

	workflow-1-zone-a                 (shard 0, families zone-a + zone-a_shard-0) - no EVM capability
	workflow-1-zone-a-shard-1         (shard 1, families zone-a + zone-a_shard-1) - no EVM capability
	chain-capabilities-zone-a-shard-0  (families zone-a_shard-0)                   - hosts evm:<selector>
	chain-capabilities-zone-a-shard-1  (families zone-a_shard-1)                   - hosts evm:<selector>
	chain-capabilities-zone-a          (families zone-a_shard-0 + zone-a_shard-1)  - hosts vault only

The two capability shards host the SAME EVM capability ID, so the only thing that can decide
which of them serves a given workflow shard is the launcher's DON-family overlap filter
(core/capabilities/launcher.go). Every assertion below therefore answers one question: which
capability DON was selected, and did anything ever route across shard families.

ExecuteShardedCapabilityCallIsolationTest (scenario 1) drives the remote EXECUTABLE path:
each phase pins the workflow owner to one shard via per_owner_assignment (manual-only mode,
no Ring OCR) and runs the evmread workflow, whose EVM calls (balance read + report write)
must be served by the in-family capability shard only.

ExecuteDon2DonDiscoveryRoutingTest (scenario 2) drives the remote TRIGGER path with the
logtrigger workflow and proves, on both ends of the don2don hop, which capability DON
discovery selected for the workflow: the workflow shard's RegisterTrigger call (client
side, carries the workflow ID) and the Event ACK on the serving capability DON (server
side), with the idle same-capability-ID DON held to zero.

Both scenarios read their per-phase evidence from live container log followers started with
Tail="0" (see startShardNodesetLogWatcher), so lines recorded by earlier subtests on this
shared environment cannot leak into a phase's absence assertions. Routing-table evidence
that is static for the topology (the executable client SetConfig lines naming remote DONs)
is instead read from the full historical logs, where absence is equally valid: a workflow
shard's routing table is built from DON families, not from shard assignment, so no subtest
can ever have added an out-of-family route to it.
*/

const (
	// executableRoutingLogNeedle matches the executable client SetConfig line
	// (core/capabilities/remote/executable/client.go), which carries
	// remoteDONName/remoteDONID naming the DON this node routes requests to.
	executableRoutingLogNeedle = "remoteDONName"
	// addRemoteCapabilityLogNeedle matches the launcher "addRemoteCapabilityV2"
	// line, which carries the remote capability's ID but not the DON name, so it
	// is combined with the capability-ID needle.
	addRemoteCapabilityLogNeedle = "addRemoteCapabilityV2"
	// serveExecutableCapLogNeedle matches the launcher "added new remote
	// executable server" line on the capability DON, which carries the
	// capability ID and method it exposes.
	serveExecutableCapLogNeedle = "added new remote executable server"
	// evmCapabilityIDLogNeedle matches the registered EVM capability ID
	// ("evm:ChainSelector:<selector>") inside launcher/client routing lines.
	evmCapabilityIDLogNeedle = "evm:ChainSelector:"
	// executableRequestReceivedLogNeedle matches the executable server
	// "received request" line (core/capabilities/remote/executable/server.go),
	// logged by the capability DON that serves each remote executable call; it
	// carries capabilityId, method and callerDonId.
	executableRequestReceivedLogNeedle = "received request"
	// registerTriggerLogNeedle matches the trigger subscriber "RegisterTrigger
	// called" line (core/capabilities/remote/trigger_subscriber.go), logged by
	// the workflow node when it registers its trigger with a capability DON; it
	// carries donId, workflowID and triggerID.
	registerTriggerLogNeedle = "RegisterTrigger called"
	// duplicateInFamilyCapWarnNeedle matches the launcher warning that fires when
	// discovery would route to the lowest DON ID because several in-family DONs
	// host the same capability - the silent-mis-routing sentinel this topology
	// must never trigger.
	duplicateInFamilyCapWarnNeedle = "multiple in-family capability DONs host the same capability"
	// evmReadWorkflowUserLog is the user log the evmread workflow emits once its
	// remote EVM calls all returned, proving the executing shard's don2don hop.
	evmReadWorkflowUserLog = "Read workflow test case passed"
	// maxNodesetWatcherLines bounds the lines a watcher records per nodeset so a
	// chatty container cannot grow the test's memory without limit.
	maxNodesetWatcherLines = 20000
)

// shardedCapabilitiesTopology carries the topology guard output shared by both routing
// scenarios: the two workflow shards, their shard indices (the terms the shard-assignment
// TOML is authored in), the in-family EVM capability DON of each shard, and the workflow
// owner / EVM chain the workflows run against.
type shardedCapabilitiesTopology struct {
	shardZeroDON *cre.Don
	shardOneDON  *cre.Don
	shardZeroIdx uint32
	shardOneIdx  uint32

	// capForShardZero/capForShardOne are the EVM capability shards the launcher's
	// family-overlap filter must route each workflow shard to.
	capForShardZero *cre.NodeSet
	capForShardOne  *cre.NodeSet

	// workflowOwner is the lowercase root-key address; per_owner_assignment keys
	// are normalized this way and every workflow in these scenarios is owned by it.
	workflowOwner string

	// chain is an EVM chain enabled on both capability shards, so both phases can
	// run their workflows against the same chain.
	chain blockchains.Blockchain
}

// mustShardedCapabilitiesTopology verifies the running environment is the
// sharded-capabilities topology these scenarios require, and returns its fixtures. The
// guard is meant to run FIRST, so a mismatched topology (e.g. CI's per-test topology
// mapping sending this test to the default config) fails in seconds with an actionable
// message instead of failing minutes later with a cryptic count.
func mustShardedCapabilitiesTopology(t *testing.T, testEnv *ttypes.TestEnvironment) shardedCapabilitiesTopology {
	t.Helper()

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	if len(shardDONs) < 2 {
		require.FailNowf(t, "wrong topology for the sharded capability routing tests",
			"expected at least 2 shard DONs, found %d. This test requires configs/workflow-sharded-capabilities-don.toml: "+
				"2 workflow shard DONs without the EVM capability + 2 capability shard DONs hosting the same EVM capability, "+
				"one per shard family. If this fails in CI, the test is missing from the per-test topology mapping "+
				"(.github/workflows/cre-system-tests.yaml PER_TEST_TOPOLOGIES_JSON and tools/ci/internal/matrix/system.go "+
				"defaultCRESmokePerTestTopologies) and ran against the default config. Running DONs: %v",
			len(shardDONs), donNames(testEnv))
	}

	shardZeroDON := getShardZeroDon(t, testEnv)
	shardOneDON := getOtherShardDON(t, testEnv, shardZeroDON)

	capForShardZero, idleForShardZero := evmCapabilityDONsByFamily(t, testEnv, nodeSetForDON(t, testEnv, shardZeroDON))
	capForShardOne, idleForShardOne := evmCapabilityDONsByFamily(t, testEnv, nodeSetForDON(t, testEnv, shardOneDON))

	// Mutual isolation guard: each workflow shard's out-of-family EVM capability DONs
	// must be exactly the other workflow shard's in-family one. Anything else means the
	// families overlap (routing cannot be told apart from a broadcast) or a third
	// capability shard exists that neither shard routes to (extra DONs, weaker proof).
	require.Lenf(t, idleForShardZero, 1, "workflow shard %s has %d out-of-family EVM capability DONs, expected exactly the other shard's in-family one",
		shardZeroDON.Name, len(idleForShardZero))
	require.Equalf(t, capForShardOne.Name, idleForShardZero[0].Name,
		"workflow shard %s's idle EVM capability DON %s is not the other shard's in-family DON %s - family layout is not mutually isolating",
		shardZeroDON.Name, idleForShardZero[0].Name, capForShardOne.Name)
	require.Lenf(t, idleForShardOne, 1, "workflow shard %s has %d out-of-family EVM capability DONs, expected exactly the other shard's in-family one",
		shardOneDON.Name, len(idleForShardOne))
	require.Equalf(t, capForShardZero.Name, idleForShardOne[0].Name,
		"workflow shard %s's idle EVM capability DON %s is not the other shard's in-family DON %s - family layout is not mutually isolating",
		shardOneDON.Name, idleForShardOne[0].Name, capForShardZero.Name)

	chain := evmChainEnabledOnNodeSet(t, testEnv, capForShardZero)
	capForShardOneChains, err := capForShardOne.GetEnabledChainIDsForCapability(cre.EVMCapability)
	require.NoErrorf(t, err, "failed to get EVM-enabled chain IDs for DON %s", capForShardOne.Name)
	require.Containsf(t, capForShardOneChains, chain.ChainID(),
		"EVM chain %d is not enabled on capability DON %s; both phases need the same chain", chain.ChainID(), capForShardOne.Name)

	workflowOwner := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.MustGetRootPrivateKey()
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	shardZeroIdx := uint32(shardZeroDON.Metadata().ShardIndex) //nolint:gosec // G115: overflow is unrealistic
	shardOneIdx := uint32(shardOneDON.Metadata().ShardIndex)   //nolint:gosec // G115: overflow is unrealistic

	framework.L.Info().
		Str("shardZeroDON", shardZeroDON.Name).
		Uint32("shardZeroIdx", shardZeroIdx).
		Str("shardOneDON", shardOneDON.Name).
		Uint32("shardOneIdx", shardOneIdx).
		Str("capForShardZero", capForShardZero.Name).
		Str("capForShardOne", capForShardOne.Name).
		Str("workflowOwner", workflowOwnerAddress).
		Str("chain", chain.CtfOutput().ChainID).
		Msg("Sharded capability routing topology")

	return shardedCapabilitiesTopology{
		shardZeroDON:    shardZeroDON,
		shardOneDON:     shardOneDON,
		shardZeroIdx:    shardZeroIdx,
		shardOneIdx:     shardOneIdx,
		capForShardZero: capForShardZero,
		capForShardOne:  capForShardOne,
		workflowOwner:   workflowOwnerAddress,
		chain:           chain,
	}
}

// getOtherShardDON returns the first shard DON that is not the given one.
func getOtherShardDON(t *testing.T, testEnv *ttypes.TestEnvironment, shard *cre.Don) *cre.Don {
	t.Helper()

	for _, don := range testEnv.Dons.DonsWithFlag(cre.ShardDON) {
		if don.ID != shard.ID {
			return don
		}
	}
	require.FailNowf(t, "no second shard DON found", "expected a shard DON other than %s", shard.Name)
	return nil
}

// proposeShardedCapabilitiesAssignment proposes and approves the shard-assignment job
// pinning the workflow owner to assignedShard on every shard DON, with the static default
// pointing at defaultShard: if the per_owner entry silently failed to apply, the workflows
// would land on the default shard and the phase's execution assertion would fail, so a
// routing that "looks right" cannot pass by accident.
func proposeShardedCapabilitiesAssignment(t *testing.T, testEnv *ttypes.TestEnvironment, topo shardedCapabilitiesTopology, assignedShard *cre.Don, assignedShardIdx, defaultShardIdx uint32) {
	t.Helper()

	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, defaultShardIdx, topo.workflowOwner, assignedShardIdx)

	// Under local load the JD occasionally drops idle node streams, so wait for full
	// shard connectivity before proposing (see awaitShardDONsConnectedToJD).
	awaitShardDONsConnectedToJD(t, testEnv)

	for _, don := range testEnv.Dons.DonsWithFlag(cre.ShardDON) {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, framework.L)
	}

	framework.L.Info().
		Str("assignedShard", assignedShard.Name).
		Uint32("assignedShardIdx", assignedShardIdx).
		Uint32("defaultShardIdx", defaultShardIdx).
		Msg("Sharded capability routing assignment proposed")
}

// requireShardedCapabilityRoutingTables asserts the static client-side routing table of
// each workflow shard from the full historical logs. No single line names both a remote
// DON and its capability ID (the SetConfig line carries remoteDONName/remoteDONID, the
// launcher lines carry the capability ID), so the route to the in-family EVM capability
// DON is proven by three shard-local facts: the workflow shard built an executable client
// for the in-family DON, it added the EVM capability as a remote capability, and the
// in-family DON exposes the EVM capability as a remote executable server. The
// out-of-family capability DON must appear in no routing entry at all.
//
// These lines are emitted at registry sync and depend only on DON families, not on shard
// assignment, so reading them from the full historical logs is valid: no subtest on this
// shared environment can ever have added an out-of-family route to a workflow shard.
func requireShardedCapabilityRoutingTables(t *testing.T, testEnv *ttypes.TestEnvironment, topo shardedCapabilitiesTopology) {
	t.Helper()

	for _, tc := range []struct {
		workflowShard *cre.Don
		inFamily      *cre.NodeSet
		outOfFamily   *cre.NodeSet
	}{
		{topo.shardZeroDON, topo.capForShardZero, topo.capForShardOne},
		{topo.shardOneDON, topo.capForShardOne, topo.capForShardZero},
	} {
		routingLines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, tc.workflowShard.Name, executableRoutingLogNeedle)
		require.Truef(t, logLinesMatchAny(routingLines, tc.inFamily.Name),
			"workflow shard %s has no executable routing entry (%s) naming its in-family capability DON %s; the launcher never built the don2don route",
			tc.workflowShard.Name, executableRoutingLogNeedle, tc.inFamily.Name)

		addCapLines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, tc.workflowShard.Name, addRemoteCapabilityLogNeedle)
		require.Truef(t, logLinesMatchAny(addCapLines, evmCapabilityIDLogNeedle),
			"workflow shard %s never added the EVM capability (%s) as a remote capability", tc.workflowShard.Name, evmCapabilityIDLogNeedle)

		serveLines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, tc.inFamily.Name, serveExecutableCapLogNeedle)
		require.Truef(t, logLinesMatchAny(serveLines, evmCapabilityIDLogNeedle),
			"capability DON %s never exposed the EVM capability (%s) as a remote executable server", tc.inFamily.Name, evmCapabilityIDLogNeedle)

		outOfFamilyLines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, tc.workflowShard.Name, tc.outOfFamily.Name)
		var violations []string
		for _, line := range outOfFamilyLines {
			if strings.Contains(line, executableRoutingLogNeedle) {
				violations = append(violations, line)
			}
		}
		require.Emptyf(t, violations,
			"workflow shard %s built an executable route to the out-of-family capability DON %s; cross-shard routing must not exist",
			tc.workflowShard.Name, tc.outOfFamily.Name)
	}

	// The same-capability-ID ambiguity sentinel must never have fired on either workflow
	// shard: this topology hosts the EVM capability exactly once per shard family, and the
	// warning is what would surface a family layout that silently routes by lowest DON ID.
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, topo.shardZeroDON.Name, duplicateInFamilyCapWarnNeedle)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, topo.shardOneDON.Name, duplicateInFamilyCapWarnNeedle)
}

// logLinesMatchAny reports whether any of the lines contains all needles.
func logLinesMatchAny(lines []string, needles ...string) bool {
	for _, line := range lines {
		if lineContainsAll(line, needles...) {
			return true
		}
	}
	return false
}

// lineContainsAll reports whether the line contains all needles.
func lineContainsAll(line string, needles ...string) bool {
	for _, needle := range needles {
		if !strings.Contains(line, needle) {
			return false
		}
	}
	return true
}

// ExecuteShardedCapabilityCallIsolationTest covers scenario 1: a workflow pinned to a
// shard calls the EVM capability over don2don, and only that shard's in-family capability
// shard may serve the calls - the other same-capability-ID shard must stay idle, and the
// calling shard must hold no route to it. Both shard directions take a turn, with a
// quiescence barrier between the phases so the reassigned workflow's engine cannot leak
// requests from the previous phase into the next phase's absence assertions.
func ExecuteShardedCapabilityCallIsolationTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	topo := mustShardedCapabilitiesTopology(t, testEnv)
	requireShardedCapabilityRoutingTables(t, testEnv, topo)

	evmChain, ok := topo.chain.(*evm.Blockchain)
	require.Truef(t, ok, "topology chain is not an EVM blockchain, got %T", topo.chain)

	// --- Phase A: owner pinned to shard 0; its EVM calls go to capability shard 0 only ---
	proposeShardedCapabilitiesAssignment(t, testEnv, topo, topo.shardZeroDON, topo.shardZeroIdx, topo.shardOneIdx)

	watcher := startShardNodesetLogWatcher(t, testEnv, shardedRoutingWatcherNodesets(topo))

	workflowA := deployShardedEVMReadWorkflow(t, testEnv, evmChain, "shardcap-isolation-shard0")
	testLogger.Info().Str("workflowID", workflowA).Msg("Phase A: evmread workflow deployed for the shard-0 assignment")

	awaitShardedWorkflowExecution(t, testEnv, []string{workflowA}, topo.shardZeroDON, evmReadWorkflowUserLog, 4*time.Minute)

	// Serving proof: the in-family capability shard served the remote executable
	// EVM requests; the other same-capability-ID shard must have received none.
	requireWatcherEventuallyLogs(t, watcher, topo.capForShardZero.Name, 2*time.Minute, executableRequestReceivedLogNeedle, evmCapabilityIDLogNeedle)
	requireWatcherNeverLogs(t, watcher, topo.capForShardOne.Name, executableRequestReceivedLogNeedle, evmCapabilityIDLogNeedle)
	testLogger.Info().Msg("Phase A: workflow on shard 0 called the EVM capability on shard 0's capability DON only")

	// --- Phase B: owner re-pinned to shard 1; the calls must move with it ---
	proposeShardedCapabilitiesAssignment(t, testEnv, topo, topo.shardOneDON, topo.shardOneIdx, topo.shardZeroIdx)

	// Quiescence barrier: after the reassignment, shard 0's engine must stop issuing
	// requests to capability shard 0. Waiting for the request stream to go quiet (and
	// recording its final count) makes the phase-B absence check immune to in-flight
	// phase-A requests while still failing loudly if the old shard keeps ownership.
	baseline := awaitWatcherLogQuiescence(t, watcher, topo.capForShardZero.Name, 45*time.Second, 3*time.Minute, executableRequestReceivedLogNeedle, evmCapabilityIDLogNeedle)
	testLogger.Info().Int("baseline", baseline).Msg("Phase B: capability shard 0 request stream is quiet")

	workflowB := deployShardedEVMReadWorkflow(t, testEnv, evmChain, "shardcap-isolation-shard1")
	testLogger.Info().Str("workflowID", workflowB).Msg("Phase B: evmread workflow deployed for the shard-1 assignment")

	awaitShardedWorkflowExecution(t, testEnv, []string{workflowB}, topo.shardOneDON, evmReadWorkflowUserLog, 4*time.Minute)

	requireWatcherEventuallyLogs(t, watcher, topo.capForShardOne.Name, 2*time.Minute, executableRequestReceivedLogNeedle, evmCapabilityIDLogNeedle)
	require.Equalf(t, baseline, watcher.lineCount(topo.capForShardZero.Name, executableRequestReceivedLogNeedle, evmCapabilityIDLogNeedle),
		"capability DON %s received new executable requests after the owner was reassigned to shard 1; cross-shard call or stale shard-0 ownership",
		topo.capForShardZero.Name)
	testLogger.Info().Msg("Phase B: workflow on shard 1 called the EVM capability on shard 1's capability DON only; shard 0's capability DON stayed idle")
}

// ExecuteDon2DonDiscoveryRoutingTest covers scenario 2: the same EVM capability ID is
// hosted by two capability DONs, and each workflow shard's don2don discovery must select
// exactly the one sharing its shard family. The selection is proven on both ends of the
// hop: the workflow shard's RegisterTrigger call (client side, carries the workflow ID)
// and the Event ACK on the serving capability DON (server side), with the idle
// same-capability-ID DON held to zero. Both shard directions take a turn.
func ExecuteDon2DonDiscoveryRoutingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	topo := mustShardedCapabilitiesTopology(t, testEnv)
	requireShardedCapabilityRoutingTables(t, testEnv, topo)

	chainID := topo.chain.CtfOutput().ChainID
	workflowConfig, msgEmitter := configureEVMLogTriggerWorkflow(t, testLogger, topo.chain)
	expectedUserLog := "Data for don2don discovery routing chain " + chainID

	emitCtx, emitCancelFn := context.WithCancel(t.Context())
	defer emitCancelFn()
	startEVMLogTriggerEventEmitter(emitCtx, t, testLogger, chainID, topo.chain, msgEmitter, expectedUserLog)

	// --- Phase A: owner pinned to shard 0; discovery must select capability shard 0 ---
	proposeShardedCapabilitiesAssignment(t, testEnv, topo, topo.shardZeroDON, topo.shardZeroIdx, topo.shardOneIdx)

	watcher := startShardNodesetLogWatcher(t, testEnv, shardedRoutingWatcherNodesets(topo))

	workflowA := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger,
		t_helpers.UniqueWorkflowName(testEnv, "don2don-routing-shard0"), &workflowConfig, "./evm/logtrigger/main.go")
	testLogger.Info().Str("workflowID", workflowA).Msg("Phase A: logtrigger workflow deployed for the shard-0 assignment")

	awaitShardedWorkflowExecution(t, testEnv, []string{workflowA}, topo.shardZeroDON, expectedUserLog, 4*time.Minute)

	// Client-side selection proof: the assigned shard registered the workflow's trigger
	// with a capability DON (the RegisterTrigger line carries donId + workflowID); the
	// other shard must never register it.
	requireWatcherEventuallyLogs(t, watcher, topo.shardZeroDON.Name, 2*time.Minute, registerTriggerLogNeedle, workflowA)
	requireWatcherNeverLogs(t, watcher, topo.shardOneDON.Name, registerTriggerLogNeedle, workflowA)

	// Server-side selection proof: the in-family capability DON acked the trigger events;
	// the idle same-capability-ID DON must have acked none.
	requireWatcherEventuallyLogs(t, watcher, topo.capForShardZero.Name, 2*time.Minute, triggerEventACKLogNeedle)
	requireWatcherNeverLogs(t, watcher, topo.capForShardOne.Name, triggerEventACKLogNeedle)
	testLogger.Info().Msg("Phase A: discovery selected capability shard 0 for the shard-0 workflow, on both ends of the don2don hop")

	// --- Phase B: owner re-pinned to shard 1; the selection must move with it ---
	proposeShardedCapabilitiesAssignment(t, testEnv, topo, topo.shardOneDON, topo.shardOneIdx, topo.shardZeroIdx)

	// Quiescence barrier: shard 0's engine must stop acking on capability shard 0 once
	// the owner moved. The baseline count makes the phase-B absence check immune to
	// in-flight phase-A acks while still failing loudly on stale shard-0 ownership.
	baseline := awaitWatcherLogQuiescence(t, watcher, topo.capForShardZero.Name, 45*time.Second, 3*time.Minute, triggerEventACKLogNeedle)
	testLogger.Info().Int("baseline", baseline).Msg("Phase B: capability shard 0 trigger ack stream is quiet")

	workflowB := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger,
		t_helpers.UniqueWorkflowName(testEnv, "don2don-routing-shard1"), &workflowConfig, "./evm/logtrigger/main.go")
	testLogger.Info().Str("workflowID", workflowB).Msg("Phase B: logtrigger workflow deployed for the shard-1 assignment")

	awaitShardedWorkflowExecution(t, testEnv, []string{workflowB}, topo.shardOneDON, expectedUserLog, 4*time.Minute)

	requireWatcherEventuallyLogs(t, watcher, topo.shardOneDON.Name, 2*time.Minute, registerTriggerLogNeedle, workflowB)
	require.Equalf(t, baseline, watcher.lineCount(topo.capForShardZero.Name, triggerEventACKLogNeedle),
		"capability DON %s acked new trigger events after the owner was reassigned to shard 1; cross-shard trigger routing or stale shard-0 ownership",
		topo.capForShardZero.Name)
	requireWatcherEventuallyLogs(t, watcher, topo.capForShardOne.Name, 2*time.Minute, triggerEventACKLogNeedle)
	testLogger.Info().Msg("Phase B: discovery selected capability shard 1 for the shard-1 workflow, on both ends of the don2don hop")
}

// deployShardedEVMReadWorkflow deploys one cron-scheduled evmread workflow against the
// given EVM chain. Its balance read and report write are remote executable calls served
// by the workflow shard's in-family EVM capability DON over don2don, which is exactly the
// hop scenario 1 isolates.
func deployShardedEVMReadWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, evmChain *evm.Blockchain, workflowBaseName string) string {
	t.Helper()

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	workflowConfig := configureEVMReadWorkflow(t, framework.L, evmChain, evm_config.TestCaseEVMReadBalance, workflowName)
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L, workflowName, &workflowConfig, "./evm/evmread/main.go")
}

// awaitShardedWorkflowExecution blocks until every workflowID logs the expected user log
// from a node of expectedShardDON. User logs from any other shard are counted as
// mismatches by waitForAllWorkflowsExecuted and fail the test unless the workflow also
// executes on the expected shard, so a workflow that ran on the wrong shard cannot pass.
func awaitShardedWorkflowExecution(t *testing.T, testEnv *ttypes.TestEnvironment, workflowIDs []string, expectedShardDON *cre.Don, expectedUserLog string, timeout time.Duration) {
	t.Helper()
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), timeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	// nodeP2PIDToShardIndex reports each node's real DON ID, so the "expected" side of
	// the comparison must be DON IDs too, even though the shard-assignment TOML is
	// authored in shard-index terms.
	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = uint32(expectedShardDON.ID) //nolint:gosec // G115: overflow is unrealistic
	}
	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, timeout)
	require.Lenf(t, executedWorkflows, len(workflowIDs), "workflows %v did not execute on shard %s (DON ID %d)", workflowIDs, expectedShardDON.Name, expectedShardDON.ID)
}

// shardedRoutingWatcherNodesets lists the nodesets whose live logs the routing scenarios
// watch: both workflow shards (client-side routing decisions) and both capability shards
// (server-side serving proof).
func shardedRoutingWatcherNodesets(topo shardedCapabilitiesTopology) []string {
	return []string{topo.shardZeroDON.Name, topo.shardOneDON.Name, topo.capForShardZero.Name, topo.capForShardOne.Name}
}

// shardNodesetLogWatcher records live container log lines per nodeset. Followers are
// started with Tail="0" and Follow=true, so only lines written after the watcher started
// are recorded: evidence emitted by earlier subtests on this shared environment cannot
// leak into a phase's absence assertions, and absence checks are trustworthy as soon as
// the positive side has landed.
type shardNodesetLogWatcher struct {
	mu      sync.Mutex
	lines   map[string][]string
	cancel  context.CancelFunc
	readers []io.ReadCloser
	wg      sync.WaitGroup
}

// startShardNodesetLogWatcher starts log followers for every container of the given
// nodesets and records their lines under the nodeset name. The watcher is stopped by
// subtest cleanup.
func startShardNodesetLogWatcher(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetNames []string) *shardNodesetLogWatcher {
	t.Helper()

	wanted := make(map[string]struct{}, len(nodesetNames))
	for _, name := range nodesetNames {
		wanted[name] = struct{}{}
	}

	containerToNodeset := make(map[string]string)
	for _, nodeSet := range testEnv.Config.NodeSets {
		if _, ok := wanted[nodeSet.Name]; !ok || nodeSet.Out == nil {
			continue
		}
		for _, clNode := range nodeSet.Out.CLNodes {
			if name := clNode.Node.ContainerName; name != "" {
				containerToNodeset[name] = nodeSet.Name
			}
		}
	}
	require.NotEmptyf(t, containerToNodeset, "no containers found for nodesets %v", nodesetNames)

	// Follow=true is required so lines written after the followers start are captured;
	// Tail="0" skips historical lines so earlier subtests on this shared environment
	// cannot contaminate this phase's absence assertions.
	logsOpts := framework.CTFContainersLogsOpts()
	logsOpts.Follow = true
	logsOpts.Tail = "0"
	logstream, err := framework.StreamContainerLogs(framework.CTFContainersListOpts(), logsOpts)
	require.NoError(t, err, "failed to stream container logs for the sharded routing watcher")

	w := &shardNodesetLogWatcher{lines: make(map[string][]string, len(nodesetNames))}
	scanCtx, cancel := context.WithCancel(t.Context())
	w.cancel = cancel

	for containerName, reader := range logstream {
		nodesetName, ok := containerToNodeset[containerName]
		if !ok {
			_ = reader.Close()
			continue
		}
		// Snapshot readers so stop can close them: closing the reader is what unblocks a
		// goroutine stuck in scanner.Scan() on a Follow=true stream; context cancel
		// alone is not enough.
		w.readers = append(w.readers, reader)
		w.wg.Go(func() {
			w.scan(scanCtx, containerName, nodesetName, reader)
		})
	}

	t.Cleanup(w.stop)
	return w
}

func (w *shardNodesetLogWatcher) scan(ctx context.Context, containerName, nodesetName string, reader io.ReadCloser) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		if ctx.Err() != nil {
			return
		}
		line := scanner.Text()
		w.mu.Lock()
		if len(w.lines[nodesetName]) < maxNodesetWatcherLines {
			w.lines[nodesetName] = append(w.lines[nodesetName], line)
		}
		w.mu.Unlock()
	}
	if err := scanner.Err(); err != nil && !isExpectedLogStreamCloseErr(ctx, err) {
		framework.L.Error().Err(err).Str("container", containerName).Msg("error reading container logs in the sharded routing watcher")
	}
}

func (w *shardNodesetLogWatcher) stop() {
	w.cancel()
	for _, r := range w.readers {
		_ = r.Close()
	}
	w.wg.Wait()
}

// lineCount returns how many recorded lines of nodesetName contain all needles.
func (w *shardNodesetLogWatcher) lineCount(nodesetName string, needles ...string) int {
	w.mu.Lock()
	defer w.mu.Unlock()

	count := 0
	for _, line := range w.lines[nodesetName] {
		if lineContainsAll(line, needles...) {
			count++
		}
	}
	return count
}

// hasLine reports whether nodesetName recorded any line containing all needles.
func (w *shardNodesetLogWatcher) hasLine(nodesetName string, needles ...string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, line := range w.lines[nodesetName] {
		if lineContainsAll(line, needles...) {
			return true
		}
	}
	return false
}

// requireWatcherEventuallyLogs waits for nodesetName to record a line containing all needles.
func requireWatcherEventuallyLogs(t *testing.T, w *shardNodesetLogWatcher, nodesetName string, timeout time.Duration, needles ...string) {
	t.Helper()

	require.Eventuallyf(t, func() bool {
		return w.hasLine(nodesetName, needles...)
	}, timeout, 2*time.Second, "nodeset %s never logged a line containing %v", nodesetName, needles)
}

// requireWatcherNeverLogs requires that nodesetName never recorded a line containing all
// needles. It is checked after the positive side of the phase has landed, so the
// counterpart evidence had time to appear on the nodeset under test.
func requireWatcherNeverLogs(t *testing.T, w *shardNodesetLogWatcher, nodesetName string, needles ...string) {
	t.Helper()

	require.Falsef(t, w.hasLine(nodesetName, needles...),
		"nodeset %s logged a line containing %v but must never do so", nodesetName, needles)
}

// awaitWatcherLogQuiescence waits until nodesetName stops recording new lines containing
// all needles for stableFor, and returns the stable count. Callers hold on to the
// returned baseline and require it unchanged later, which turns "no new lines since the
// barrier" into an assertion that fails loudly if the reassignment left stale owners
// behind.
func awaitWatcherLogQuiescence(t *testing.T, w *shardNodesetLogWatcher, nodesetName string, stableFor, timeout time.Duration, needles ...string) int {
	t.Helper()

	lastCount := w.lineCount(nodesetName, needles...)
	lastChange := time.Now()

	require.Eventuallyf(t, func() bool {
		count := w.lineCount(nodesetName, needles...)
		if count != lastCount {
			lastCount = count
			lastChange = time.Now()
			return false
		}
		return time.Since(lastChange) >= stableFor
	}, timeout, 2*time.Second,
		"nodeset %s did not go quiet on %v within %s; the previous assignment's owner is likely still issuing traffic",
		nodesetName, needles, timeout)

	return lastCount
}
