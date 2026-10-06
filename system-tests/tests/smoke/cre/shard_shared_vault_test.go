package cre

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	vaultsecretstypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/vault_secrets/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

// Static shard assignment across two shards with a single shared vault DON.
//
// Both scenarios run on the same topology shape (see
// configs/workflow-gateway-sharded-shared-vault-{manual,failover}.toml):
//
//	workflow-1-zone-a         (shard 0) - no vault capability of its own
//	workflow-1-zone-a-shard-1 (shard 1) - no vault capability of its own
//	chain-capabilities-zone-a          - the ONLY vault DON, in both shard families
//	bootstrap-gateway
//
// Neither workflow shard hosts a vault, so every secret a workflow reads is served
// by the one shared vault DON over the workflow-DON -> vault-DON remote capability
// hop - whichever shard the assignment routes the workflow to.
//
// ExecuteManualShardAssignmentSharedVaultTest covers manual-only mode with the
// failover gate closed: only the assigned shard syncs and runs the workflow.
// ExecuteShardFailoverSharedVaultTest covers the same assignment mode with
// ShardingFailoverEnabled open: the ordered per_owner list makes both shards sync
// the workflow, the primary executes while the secondary keeps a standby engine
// that denies and caches triggers, and re-proposing the list swapped moves
// execution to the other shard without losing access to the same secret.

const (
	sharedVaultWorkflowFileLocation = "../../../../core/scripts/cre/environment/examples/workflows/vault_secrets/main.go"
	sharedVaultSchedule             = "*/30 * * * * *"
	sharedVaultNamespace            = "main"
	// The log the vault_secrets example emits per execution; the secret value in
	// it is the observable proof that the fetched secret decrypted on the shard
	// that executed the workflow.
	sharedVaultUserLogPrefix = "Vault secret fetched: "
	// What a secondary shard logs when the failover gate is open and it denies
	// (and caches) a trigger it does not own (ShardFailoverManager.cacheEvent).
	sharedVaultTriggerCachedLogNeedle = "secondary shard: cached trigger event for failover"
)

// sharedVaultShardFixture carries the one-off vault setup shared by both scenarios:
// a funded seth client, the workflow owner (the test root key, linked in the workflow
// registry), and the gateway URL + vault public key needed to create secrets.
type sharedVaultShardFixture struct {
	gatewayURL     string
	parsedVaultKey *tdh2easy.PublicKey
	sc             *seth.Client
	owner          string // workflow owner, checksummed hex; also the vault SecretIdentifier.Owner
	wfRegistry     *workflow_registry_v2_wrapper.WorkflowRegistry
}

func setupSharedVaultShardFixture(t *testing.T, testEnv *ttypes.TestEnvironment) *sharedVaultShardFixture {
	t.Helper()

	// Same prerequisites as the vault DON suite: wait for the vault DON's DKG to
	// produce its result packages, fetch the vault public key through the gateway,
	// and inject public key + threshold into the vault@1.0.0 capability config so
	// workflow nodes can encrypt secret requests and verify responses.
	ensureVaultDKGResultPackages(t, testEnv)
	gatewayURL := mustVaultGatewayURL(t, testEnv)
	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	updateVaultCapabilityConfigInRegistry(t, testEnv, vaultPublicKey)

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	// The workflow owner signs (allowlisted) vault requests, so it must be linked
	// in the workflow registry before the gateway accepts its secret creations.
	wfRegistryAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	wfRegistry, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegistryAddr), sc.Client)
	require.NoError(t, err, "failed to create workflow registry wrapper")
	requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegistryAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])

	return &sharedVaultShardFixture{
		gatewayURL:     gatewayURL.String(),
		parsedVaultKey: mustVaultPublicKey(t, vaultPublicKey),
		sc:             sc,
		owner:          sc.MustGetRootKeyAddress().Hex(),
		wfRegistry:     wfRegistry,
	}
}

// createSharedVaultSecret stores value in the shared vault under a fresh key in the
// "main" namespace, authorized by the workflow owner's allowlisted signature.
func (f *sharedVaultShardFixture) createSharedVaultSecret(t *testing.T, value string) string {
	t.Helper()

	secretKey := uniqueVaultSecretID("sharedvault")
	encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner(value, f.parsedVaultKey, f.sc.MustGetRootKeyAddress())
	require.NoError(t, err, "failed to encrypt secret for the shared vault")

	auth := newAllowlistVaultRequestAuth(f.owner, f.sc, f.wfRegistry)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretKey, f.owner, f.gatewayURL, []string{sharedVaultNamespace})
	return secretKey
}

// deploySharedVaultWorkflow deploys one cron-scheduled vault_secrets workflow that
// reads the given secret from the shared vault on every run. The registered name is
// unique per test run (UniqueWorkflowName), so re-running against the same shared
// environment yields fresh workflow IDs instead of deduplicating against the
// previous run's executions.
func deploySharedVaultWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowBaseName, secretKey string) string {
	t.Helper()

	workflowConfig := vaultsecretstypes.WorkflowConfig{
		Schedule:        sharedVaultSchedule,
		SecretNamespace: sharedVaultNamespace,
		SecretKey:       secretKey,
	}
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L, t_helpers.UniqueWorkflowName(testEnv, workflowBaseName), &workflowConfig, sharedVaultWorkflowFileLocation)
}

// awaitSharedVaultWorkflowExecution blocks until every workflowID logs the user log
// carrying secretValue from a node of expectedShardDON. User logs from any other
// shard are counted as mismatches by waitForAllWorkflowsExecuted and fail the test
// unless the workflow also executes on the expected shard.
func awaitSharedVaultWorkflowExecution(t *testing.T, testEnv *ttypes.TestEnvironment, workflowIDs []string, expectedShardDON *cre.Don, secretValue string, timeout time.Duration) {
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

	// nodeP2PIDToShardIndex reports each node's real DON ID, so the "expected"
	// side of the comparison must be DON IDs too, even though the shard-assignment
	// TOML is authored in shard-index terms.
	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = uint32(expectedShardDON.ID) //nolint:gosec // G115: overflow is unrealistic
	}
	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, sharedVaultUserLogPrefix+secretValue, timeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "workflows %v did not execute on shard %s (DON ID %d)", workflowIDs, expectedShardDON.Name, expectedShardDON.ID)
}

// sharedVaultShardPair holds the two shard DONs of the topology with their shard
// indices (the terms the shard-assignment TOML is authored in).
type sharedVaultShardPair struct {
	shardZeroDON  *cre.Don
	shardZeroIdx  uint32
	shardOneDON   *cre.Don
	shardOneIdx   uint32
	workflowOwner string // lowercase owner hex; per_owner_assignment keys are normalized this way
}

// mustSharedVaultShardPair verifies the running environment is the sharded
// shared-vault topology this test requires, and returns its two shard DONs. It is
// meant to run FIRST, before the minutes-long vault setup: a mismatched topology
// (e.g. CI's per-test topology mapping sending this test to the default config)
// then fails in seconds with an actionable message instead of failing after the
// vault setup with a cryptic count.
func mustSharedVaultShardPair(t *testing.T, testEnv *ttypes.TestEnvironment) sharedVaultShardPair {
	t.Helper()

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	if len(shardDONs) < 2 {
		require.FailNowf(t, "wrong topology for the sharded shared-vault tests",
			"expected at least 2 shard DONs, found %d. This test requires the sharded shared-vault topology "+
				"(configs/workflow-gateway-sharded-shared-vault-manual.toml / configs/workflow-gateway-sharded-shared-vault-failover.toml: "+
				"2 workflow shard DONs without the vault capability + 1 shared vault DON). If this fails in CI, the test is missing from the "+
				"per-test topology mapping (.github/workflows/cre-system-tests.yaml PER_TEST_TOPOLOGIES_JSON and "+
				"tools/ci/internal/matrix/system.go defaultCRESmokePerTestTopologies) and ran against the default config. "+
				"Running DONs: %v", len(shardDONs), donNames(testEnv))
	}

	// The vault must be hosted by exactly one NON-shard capabilities DON - the shared
	// vault both shards route to. A vault on a shard DON is a different (per-shard
	// vault) topology where secret continuity across shards cannot hold.
	sharedVaultDONs := slices.DeleteFunc(slices.Clone(testEnv.Dons.DonsWithFlag(cre.VaultCapability)), func(don *cre.Don) bool {
		return don.HasFlag(cre.ShardDON)
	})
	if len(sharedVaultDONs) != 1 {
		require.FailNowf(t, "wrong topology for the sharded shared-vault tests",
			"expected exactly 1 shared vault DON (a non-shard capabilities DON hosting the vault), found %d. "+
				"This test requires the sharded shared-vault topology (configs/workflow-gateway-sharded-shared-vault-manual.toml / "+
				"configs/workflow-gateway-sharded-shared-vault-failover.toml). Running DONs: %v", len(sharedVaultDONs), donNames(testEnv))
	}

	shardZeroDON := getShardZeroDon(t, testEnv)
	nonLeaderDONs := slices.DeleteFunc(slices.Clone(shardDONs), func(don *cre.Don) bool {
		return don.ID == shardZeroDON.ID
	})
	require.NotEmpty(t, nonLeaderDONs, "Expected to find a second shard DON")

	return sharedVaultShardPair{
		shardZeroDON: shardZeroDON,
		shardZeroIdx: uint32(shardZeroDON.Metadata().ShardIndex), //nolint:gosec // G115: overflow is unrealistic
		shardOneDON:  nonLeaderDONs[0],
		shardOneIdx:  uint32(nonLeaderDONs[0].Metadata().ShardIndex), //nolint:gosec // G115: overflow is unrealistic
	}
}

// donNames lists the running DON names for topology-mismatch diagnostics.
func donNames(testEnv *ttypes.TestEnvironment) []string {
	dons := testEnv.Dons.List()
	names := make([]string, 0, len(dons))
	for _, don := range dons {
		names = append(names, don.Name)
	}
	return names
}

// awaitShardDONsConnectedToJD waits until every worker node of every shard DON reports
// connected to the job distributor. Under local load the JD occasionally drops idle node
// streams (wsrpc "Read error: i/o timeout"); nodes reconnect within seconds, but
// proposing a shard-assignment job to a node mid-reconnect fails with
// "node is not connected", so each proposal round waits for full shard connectivity first.
func awaitShardDONsConnectedToJD(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.Eventually(t, func() bool {
		for _, don := range shardDONs {
			workers, err := don.Workers()
			if err != nil {
				return false
			}
			nodes, err := offchain.FetchNodesFromJD(t.Context(), testEnv.CreEnvironment.CldfEnvironment.Offchain, offchain.TargetDONFilter{
				Key:   offchain.FilterKeyDONName,
				Value: don.Name,
			}.ToListFilter())
			if err != nil || len(nodes) != len(workers) {
				return false
			}
			for _, node := range nodes {
				if !node.IsConnected {
					return false
				}
			}
		}
		return true
	}, 2*time.Minute, 3*time.Second, "shard DON worker nodes did not all report connected to the job distributor")
}

// proposeSharedVaultAssignment proposes and approves the shard-assignment job with
// the given TOML on every shard DON: each shard resolves ownership from its own
// copy of the spec, so both the shard that must run the workflow and the shard that
// must not need it.
func proposeSharedVaultAssignment(t *testing.T, testEnv *ttypes.TestEnvironment, shardAssignmentTOML string) {
	t.Helper()

	awaitShardDONsConnectedToJD(t, testEnv)

	for _, don := range testEnv.Dons.DonsWithFlag(cre.ShardDON) {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, framework.L)
	}
}

// ExecuteManualShardAssignmentSharedVaultTest covers manual-only assignment with a
// single shared vault: the per_owner entry routes the workflow to exactly one shard
// at a time, and both shards get their turn reading the SAME secret from the SAME
// shared vault DON.
func ExecuteManualShardAssignmentSharedVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Topology guard first (see mustSharedVaultShardPair): a mismatched environment
	// fails in seconds with an actionable message, before the vault setup waits.
	shards := mustSharedVaultShardPair(t, testEnv)

	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "manual-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	// Phase 1: pin the owner to shard 1, with the static default pointing at shard 0
	// so an assignment that silently did not take effect cannot look like a pass -
	// the workflow only lands on shard 1 if the per_owner entry routed it there.
	shardOneAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, shards.shardZeroIdx, shards.workflowOwner, shards.shardOneIdx)

	proposeSharedVaultAssignment(t, testEnv, shardOneAssignmentTOML)

	workflowID := deploySharedVaultWorkflow(t, testEnv, "manual-shared-vault-shard1", secretKey)
	testLogger.Info().Str("workflowID", workflowID).Msg("Phase 1: workflow deployed for the shard-1 assignment")

	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, 3*time.Minute)
	testLogger.Info().Msg("Phase 1: workflow executed on shard 1 and read the secret from the shared vault")

	// The failover gate is closed in this topology, so the shard the assignment did
	// NOT pick never syncs the workflow at all - it must hold no standby engine and
	// therefore never log a cached trigger event.
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle)

	// Phase 2: re-propose with the owner pinned to shard 0 and the static default
	// flipped to shard 1. A fresh workflow (fresh execution IDs; the engine store
	// deduplicates the phase-1 workflow's executions) must now execute on shard 0
	// and read the SAME secret from the same shared vault DON.
	shardZeroAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx)

	proposeSharedVaultAssignment(t, testEnv, shardZeroAssignmentTOML)

	workflowID = deploySharedVaultWorkflow(t, testEnv, "manual-shared-vault-shard0", secretKey)
	testLogger.Info().Str("workflowID", workflowID).Msg("Phase 2: workflow deployed for the shard-0 assignment")

	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, 3*time.Minute)
	testLogger.Info().Msg("Phase 2: workflow executed on shard 0 and read the same secret from the shared vault")

	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle)
}

// ExecuteShardFailoverSharedVaultTest covers shard failover with a single shared
// vault: the ordered per_owner list assigns a primary and a secondary shard, the
// primary executes while the secondary keeps a standby engine that denies and caches
// every trigger, and re-proposing the list with the entries swapped moves execution
// to the secondary - which keeps reading the SAME secret from the SAME shared vault
// DON, because secret continuity does not depend on which shard executes.
func ExecuteShardFailoverSharedVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Topology guard first (see mustSharedVaultShardPair): a mismatched environment
	// fails in seconds with an actionable message, before the vault setup waits.
	shards := mustSharedVaultShardPair(t, testEnv)

	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	// Phase 1: ordered assignment [shard 0, shard 1] - position 0 is the primary,
	// position 1 the secondary.
	primaryAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, shards.shardZeroIdx, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx, shards.shardOneIdx)

	proposeSharedVaultAssignment(t, testEnv, primaryAssignmentTOML)

	const numWorkflows = 2
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowIDs = append(workflowIDs, deploySharedVaultWorkflow(t, testEnv, fmt.Sprintf("failover-shared-vault-%d", i), secretKey))
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Phase 1: workflows deployed for the [shard 0, shard 1] assignment")

	awaitSharedVaultWorkflowExecution(t, testEnv, workflowIDs, shards.shardZeroDON, secretValue, 3*time.Minute)
	testLogger.Info().Msg("Phase 1: workflows executed on the primary shard (shard 0) and read the secret from the shared vault")

	// The secondary shard holds a standby engine for every workflow: it receives the
	// same triggers, is denied by the shard ownership check, and caches each event
	// for a potential replay. Assert the caching so a silently-missing standby fails
	// loudly instead of surfacing only when a real failover is needed.
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: secondary shard (shard 1) is denying and caching triggers as expected")

	// Phase 2: swap the assignment. Fresh workflows are deployed BEFORE re-proposing
	// (like ExecuteFailoverManualSwapTest) so the swap itself must move them: they
	// sync on both shards under the current assignment, and only the re-proposed
	// order can decide who executes. The fresh names also give clean execution IDs -
	// the store deduplicates the phase-1 workflows' executions.
	for i := range numWorkflows {
		workflowIDs[i] = deploySharedVaultWorkflow(t, testEnv, fmt.Sprintf("failover-shared-vault-swap-%d", i), secretKey)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Phase 2: fresh workflows deployed for the swap")

	secondaryAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, shards.shardOneIdx, shards.shardZeroIdx, shards.workflowOwner, shards.shardOneIdx, shards.shardZeroIdx)

	proposeSharedVaultAssignment(t, testEnv, secondaryAssignmentTOML)

	awaitSharedVaultWorkflowExecution(t, testEnv, workflowIDs, shards.shardOneDON, secretValue, 5*time.Minute)
	testLogger.Info().Msg("Phase 2: workflows executed on the new primary shard (shard 1) and read the SAME secret from the shared vault")

	// After the swap the old primary (shard 0) is the secondary: it must hold a
	// standby engine again and cache the triggers it is denied. (Cache DRAINING on
	// the new secondary is a known limitation today: the communicator validates
	// senders against the primary DON registered at engine start, so after a swap
	// the old primary's cache is not drained by the new primary's status updates
	// until its engine restarts. Cached events simply expire; execution is not
	// affected. Only the caching itself is asserted here.)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 2: old primary (shard 0) is denying and caching triggers as the new secondary")
}
