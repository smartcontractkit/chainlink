package cre

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	vaultsecretcron_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecretcron/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	sharedVaultWorkflowFileLocation = "./vaultsecretcron/main.go"
	sharedVaultSchedule             = "*/30 * * * * *"
	sharedVaultNamespace            = "main"
	sharedVaultUserLogPrefix        = "Vault secret fetched: "

	sharedVaultTriggerCachedLogNeedle = "secondary shard: cached trigger event for failover"
)

type sharedVaultShardFixture struct {
	gatewayURL     string
	vaultPublicKey string // hex; the raw form enclave configs need
	parsedVaultKey *tdh2easy.PublicKey
	sc             *seth.Client
	owner          string
	wfRegistry     *workflow_registry_v2_wrapper.WorkflowRegistry
}

func setupSharedVaultShardFixture(t *testing.T, testEnv *ttypes.TestEnvironment) *sharedVaultShardFixture {
	t.Helper()

	ensureVaultDKGResultPackages(t, testEnv)
	gatewayURL := mustVaultGatewayURL(t, testEnv)
	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	updateVaultCapabilityConfigInRegistry(t, testEnv, vaultPublicKey)

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

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
		vaultPublicKey: vaultPublicKey,
		parsedVaultKey: mustVaultPublicKey(t, vaultPublicKey),
		sc:             sc,
		owner:          sc.MustGetRootKeyAddress().Hex(),
		wfRegistry:     wfRegistry,
	}
}

func (f *sharedVaultShardFixture) createSharedVaultSecret(t *testing.T, value string) string {
	t.Helper()

	secretKey := uniqueVaultSecretID("sharedvault")
	encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner(value, f.parsedVaultKey, f.sc.MustGetRootKeyAddress())
	require.NoError(t, err, "failed to encrypt secret for the shared vault")

	auth := newAllowlistVaultRequestAuth(f.owner, f.sc, f.wfRegistry)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretKey, f.owner, f.gatewayURL, []string{sharedVaultNamespace})
	return secretKey
}

func deploySharedVaultWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowBaseName, secretKey string) string {
	t.Helper()

	workflowConfig := vaultsecretcron_config.Config{
		Schedule:        sharedVaultSchedule,
		SecretNamespace: sharedVaultNamespace,
		SecretKey:       secretKey,
	}
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L, t_helpers.UniqueWorkflowName(testEnv, workflowBaseName), &workflowConfig, sharedVaultWorkflowFileLocation)
}

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

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = uint32(expectedShardDON.ID) //nolint:gosec // G115: overflow is unrealistic
	}
	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, sharedVaultUserLogPrefix+secretValue, timeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "workflows %v did not execute on shard %s (DON ID %d)", workflowIDs, expectedShardDON.Name, expectedShardDON.ID)
}

type sharedVaultShardPair struct {
	shardZeroDON  *cre.Don
	shardZeroIdx  uint32
	shardOneDON   *cre.Don
	shardOneIdx   uint32
	workflowOwner string
}

func mustSharedVaultShardPair(t *testing.T, testEnv *ttypes.TestEnvironment, topologyHint string) sharedVaultShardPair {
	t.Helper()

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	if len(shardDONs) < 2 {
		require.FailNowf(t, "wrong topology for the sharded shared-vault tests",
			"expected at least 2 shard DONs, found %d. This test requires the topology pair %s: "+
				"2 workflow shard DONs without the vault capability + 1 shared vault DON. If this fails in CI, the test is missing from the "+
				"per-test topology mapping (.github/workflows/cre-system-tests.yaml PER_TEST_TOPOLOGIES_JSON and "+
				"tools/ci/internal/matrix/system.go defaultCRESmokePerTestTopologies) and ran against the default config. "+
				"Running DONs: %v", len(shardDONs), topologyHint, donNames(testEnv))
	}

	sharedVaultDONs := slices.DeleteFunc(slices.Clone(testEnv.Dons.DonsWithFlag(cre.VaultCapability)), func(don *cre.Don) bool {
		return don.HasFlag(cre.ShardDON)
	})
	if len(sharedVaultDONs) != 1 {
		require.FailNowf(t, "wrong topology for the sharded shared-vault tests",
			"expected exactly 1 shared vault DON (a non-shard capabilities DON hosting the vault), found %d. "+
				"This test requires the topology pair %s. Running DONs: %v", len(sharedVaultDONs), topologyHint, donNames(testEnv))
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

func donNames(testEnv *ttypes.TestEnvironment) []string {
	dons := testEnv.Dons.List()
	names := make([]string, 0, len(dons))
	for _, don := range dons {
		names = append(names, don.Name)
	}
	return names
}

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

func proposeSharedVaultAssignment(t *testing.T, testEnv *ttypes.TestEnvironment, shardAssignmentTOML string) {
	t.Helper()

	awaitShardDONsConnectedToJD(t, testEnv)

	for _, don := range testEnv.Dons.DonsWithFlag(cre.ShardDON) {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, framework.L)
	}
}

func ExecuteManualShardAssignmentSharedVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Topology guard first (see mustSharedVaultShardPair): a mismatched environment
	// fails in seconds with an actionable message, before the vault setup waits.
	shards := mustSharedVaultShardPair(t, testEnv,
		"configs/workflow-gateway-sharded-shared-vault-manual.toml / configs/workflow-gateway-sharded-shared-vault-failover.toml")

	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "manual-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

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

	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle)

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

func ExecuteShardFailoverSharedVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Topology guard first (see mustSharedVaultShardPair): a mismatched environment
	// fails in seconds with an actionable message, before the vault setup waits.
	shards := mustSharedVaultShardPair(t, testEnv,
		"configs/workflow-gateway-sharded-shared-vault-manual.toml / configs/workflow-gateway-sharded-shared-vault-failover.toml")

	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

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

	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: secondary shard (shard 1) is denying and caching triggers as expected")

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

	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 2: old primary (shard 0) is denying and caching triggers as the new secondary")
}
