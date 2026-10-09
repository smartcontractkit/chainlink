package cre

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-confidential-compute/tests/testhelpers"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	shardedConfidentialManualConfigPath   = "/configs/workflow-gateway-sharded-confidential-workflows-manual.toml"
	shardedConfidentialFailoverConfigPath = "/configs/workflow-gateway-sharded-confidential-workflows-failover.toml"

	shardedConfidentialTopologyHint = "configs/workflow-gateway-sharded-confidential-workflows-manual.toml / " +
		"configs/workflow-gateway-sharded-confidential-workflows-failover.toml"

	shardedConfidentialWorkflowSrc = "./confidentialvaultsecretcron"

	shardedConfidentialSchedule    = "*/30 * * * * *"
	shardedConfidentialNamespace   = "main"
	shardedConfidentialExecTimeout = 5 * time.Minute

	confidentialWorkflowSuccessLogNeedle = `"msg":"Workflow execution finished successfully"`

	shardZeroEnclaveBaseCID = 16
	shardOneEnclaveBaseCID  = 18
)

// One enclave group per shard. Group shard 0 on the harness's default ports,
// group shard 1 on the next block, so both fit the CI runner's gateway egress
// allowlist (-e 8080..8087). Distinct BaseCIDs keep the two groups apart on
// Nitro-capable hosts.
var (
	shardZeroEnclaveHTTPPorts   = []string{"8080", "8081"}
	shardZeroEnclaveConfigPorts = []string{"8082", "8083"}
	shardOneEnclaveHTTPPorts    = []string{"8084", "8085"}
	shardOneEnclaveConfigPorts  = []string{"8086", "8087"}
)

// ExecuteShardManualAssignmentConfidentialWorkflowsTest covers manual-only shard
// assignment for a confidential workflow with a single shared vault: the
// per_owner entry pins the workflow to shard 0, only that shard syncs it, and
// its engine executes the workflow inside shard 0's enclaves, fetching the
// shared vault secret through the confidential relay.
//
// The static default points at shard 1 on purpose: shard 1's enclave group is
// configured but its relay path cannot serve shard-1-bound enclaves (the relay
// nodes verify the enclave's signers against their own DON), so a silently
// missed per_owner entry lands the workflow on shard 1 where no execution can
// ever succeed - a visible timeout failure, never a false pass.
func ExecuteShardManualAssignmentConfidentialWorkflowsTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	ccRoot := confidentialComputeRoot(t)

	shards := mustSharedVaultShardPair(t, testEnv, shardedConfidentialTopologyHint)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "manual-sharded-confidential-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	enclaveHost := confidentialEnclaveHostAddr(testhelpers.UseFakeEnclave())
	storageSvc := startConfidentialHostServices(t, testEnv, enclaveHost)

	setupShardedConfidentialEnclaves(t, testEnv, testLogger, ccRoot, shards, fixture)

	assignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx)
	proposeSharedVaultAssignment(t, testEnv, assignmentTOML)

	workflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, storageSvc, enclaveHost,
		"manual-shard0-confidential", secretKey, secretValue)

	// The workflow returns an error unless GetSecret returns exactly the stored
	// secret value, so this success log proves the whole chain: shard-0 sync,
	// engine -> enclave execution, enclave -> relay -> shared vault fetch.
	waitForConfidentialWorkflowExecution(t, testEnv, testLogger, workflowID, shardedConfidentialExecTimeout, shards.shardZeroDON.Name)
	assertNoConfidentialExecutionOnNodeset(t, testEnv, workflowID, shards.shardOneDON.Name)

	// The failover gate is closed here, so shard 1 must not even hold a standby
	// engine for the workflow.
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle)
}

// ExecuteShardFailoverConfidentialWorkflowsTest covers shard failover for a
// confidential workflow with a single shared vault: the ordered per_owner list
// assigns shard 0 as primary and shard 1 as secondary, both shards sync the
// workflow, shard 0 executes it inside its enclaves while shard 1 keeps a
// standby engine that denies and caches triggers, and re-proposing the list
// swapped moves execution to shard 1's enclaves.
//
// The secret leg is dropped for the swap phase: the enclave -> relay secret
// path is bound to the first shard today (the gateway's confidential handler
// fans out to that shard's nodes, and each relay node verifies the enclave's
// signers against its own DON), so a shard-1 execution cannot fetch secrets
// through shard-0's relay nodes. The secret path is proven on the phase-1
// primary instead; the swap phase proves execution itself moved.
func ExecuteShardFailoverConfidentialWorkflowsTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	ccRoot := confidentialComputeRoot(t)

	shards := mustSharedVaultShardPair(t, testEnv, shardedConfidentialTopologyHint)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-sharded-confidential-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	enclaveHost := confidentialEnclaveHostAddr(testhelpers.UseFakeEnclave())
	storageSvc := startConfidentialHostServices(t, testEnv, enclaveHost)

	setupShardedConfidentialEnclaves(t, testEnv, testLogger, ccRoot, shards, fixture)

	// Phase 1: ordered assignment [shard 0, shard 1] - shard 0 is the primary.
	primaryAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, shards.shardZeroIdx, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx, shards.shardOneIdx)
	proposeSharedVaultAssignment(t, testEnv, primaryAssignmentTOML)

	workflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, storageSvc, enclaveHost,
		"failover-shard0-confidential", secretKey, secretValue)
	waitForConfidentialWorkflowExecution(t, testEnv, testLogger, workflowID, shardedConfidentialExecTimeout, shards.shardZeroDON.Name)
	assertNoConfidentialExecutionOnNodeset(t, testEnv, workflowID, shards.shardOneDON.Name)

	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)

	// Phase 2: swap. The fresh workflow is deployed BEFORE re-proposing so the
	// swap itself must move it, and drops the secret leg (see the function
	// comment for why). Its success on shard 1 proves the confidential execution
	// path - engine, per-shard enclaves, on-chain routing - followed the swap.
	swapWorkflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, storageSvc, enclaveHost,
		"failover-swap-shard1-confidential", "", "")

	secondaryAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, shards.shardOneIdx, shards.shardZeroIdx, shards.workflowOwner, shards.shardOneIdx, shards.shardZeroIdx)
	proposeSharedVaultAssignment(t, testEnv, secondaryAssignmentTOML)

	waitForConfidentialWorkflowExecution(t, testEnv, testLogger, swapWorkflowID, shardedConfidentialExecTimeout, shards.shardOneDON.Name)
	assertNoConfidentialExecutionOnNodeset(t, testEnv, swapWorkflowID, shards.shardZeroDON.Name)

	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
}

// startConfidentialHostServices stands up the host-side services the enclaves
// need before they start, and bakes their addresses into the enclave settings:
// the deferred gateway proxy (the real gateway URL is only known once the
// environment is up) and the fake CRE storage service the enclaves fetch the
// workflow binary from. It also sets the enclave host env vars the harness
// passes to the enclave host servers, so it must run before any enclave group.
func startConfidentialHostServices(t *testing.T, testEnv *ttypes.TestEnvironment, enclaveHost string) *fakeStorageService {
	t.Helper()

	gwProxy := newDeferredGatewayProxy(t, confidentialGatewayProxyPort)
	storageAddr, storageSvc := startFakeStorageService(t, enclaveHost)

	t.Setenv("REQUIRE_BFT_QUORUM", "true")
	t.Setenv("OUTBOUND_ALLOW_LOCAL_FOR_TESTS", "true")
	t.Setenv("ENCLAVE_SETTINGS", fmt.Sprintf(
		`{"storageKey":%q,"storageServiceUrl":%q,"storageServiceTls":false,"gatewayUrl":%q}`,
		confidentialStorageKeyHex,
		storageAddr,
		fmt.Sprintf("http://%s:%d", enclaveHost, confidentialGatewayProxyPort),
	))

	gatewayURL := confidentialGatewayURL(t, testEnv)
	require.NoError(t, gwProxy.SetTarget(gatewayURL), "failed to set gateway proxy target")

	return storageSvc
}

// setupShardedConfidentialEnclaves starts one enclave group per shard, configures
// each group's signers to its shard's DON with the vault master public key, and
// publishes each group in that shard DON's own on-chain capability config. Both
// shards' capabilities refresh on the same ticker, so one settle window covers
// both publishes.
func setupShardedConfidentialEnclaves(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	ccRoot string,
	shards sharedVaultShardPair,
	fixture *sharedVaultShardFixture,
) {
	t.Helper()

	startShardEnclaveGroup := func(httpPorts, configPorts []string, baseCID int) *testhelpers.LocalEnclaveResult {
		return testhelpers.SetupLocalEnclaves(t, testhelpers.LocalEnclaveSetupConfig{
			RepoRoot: ccRoot,
			AppName:  confidentialWorkflowsApp,
			// The enclave URLs are dialed from the node containers, so they need the
			// Docker-reachable host (framework.HostDockerInternal: host.docker.internal
			// locally, the bridge IP on CI runners), not the test host's own address.
			// The config-plane URLs the test dials stay loopback.
			HostIP:          strings.TrimPrefix(framework.HostDockerInternal(), "http://"),
			HTTPPorts:       httpPorts,
			ConfigHTTPPorts: configPorts,
			BaseCID:         baseCID,
			Region:          confidentialEnclaveRegion,
		})
	}

	shardZeroGroup := startShardEnclaveGroup(shardZeroEnclaveHTTPPorts, shardZeroEnclaveConfigPorts, shardZeroEnclaveBaseCID)
	shardOneGroup := startShardEnclaveGroup(shardOneEnclaveHTTPPorts, shardOneEnclaveConfigPorts, shardOneEnclaveBaseCID)

	configureEnclaves(t, testLogger, shardZeroGroup.ConfigURLs, shards.shardZeroDON, fixture.vaultPublicKey)
	configureEnclaves(t, testLogger, shardOneGroup.ConfigURLs, shards.shardOneDON, fixture.vaultPublicKey)

	publishEnclaves(t, testEnv, testLogger, shards.shardZeroDON.Name, shardZeroGroup.Enclaves)
	publishEnclaves(t, testEnv, testLogger, shards.shardOneDON.Name, shardOneGroup.Enclaves)
}

// deployShardedConfidentialWorkflow compiles, serves, copies and registers the
// confidentialvaultsecretcron workflow with confidential attributes. An empty
// secretKey deploys it without the secret leg (see the workflow's package
// comment). It returns the workflow ID.
func deployShardedConfidentialWorkflow(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	storageSvc *fakeStorageService,
	enclaveHost string,
	workflowBaseName string,
	secretKey string,
	expectedSecretValue string,
) string {
	t.Helper()

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	configJSON := fmt.Sprintf(`{"schedule":%q,"secret_namespace":%q,"secret_key":%q,"expected_secret_value":%q}`,
		shardedConfidentialSchedule, shardedConfidentialNamespace, secretKey, expectedSecretValue)

	// The artifact URL must be reachable by whoever fetches the binary: the
	// registration download from the test process, and the enclaves through the
	// storage service's returned URL. Both run on the test host, so the enclave
	// host (loopback for fake enclaves, the wireguard host IP for real ones)
	// covers them; the syncer reads the binary from disk instead.
	artifacts := buildAndServeConfidentialWorkflow(t, shardedConfidentialWorkflowSrc, workflowName, configJSON, enclaveHost)
	storageSvc.setArtifactURL(artifacts.BinaryFilename, artifacts.BinaryURL)

	copyWorkflowArtifactsToContainers(t, testEnv, artifacts)

	return registerConfidentialWorkflow(t, testEnv, testLogger, workflowName, artifacts)
}

// assertNoConfidentialExecutionOnNodeset fails the test if any container of the
// named DON logged a successful execution of the workflow. Only an executing
// engine emits that line, so finding it on a shard the assignment did not pick
// means the routing is broken.
func assertNoConfidentialExecutionOnNodeset(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string, donName string) {
	t.Helper()

	containers := confidentialWorkflowDONContainers(testEnv, donName)
	require.NotEmpty(t, containers, "no containers found for DON %s", donName)

	for _, container := range containers {
		out, err := exec.CommandContext(t.Context(), "docker", "logs", "--tail", "10000", container).CombinedOutput()
		require.NoError(t, err, "failed to read logs of container %s", container)
		for line := range bytes.SplitSeq(out, []byte{'\n'}) {
			if bytes.Contains(line, []byte(confidentialWorkflowSuccessLogNeedle)) && bytes.Contains(line, []byte(workflowID)) {
				t.Fatalf("workflow %s executed on DON %s (container %s), which the assignment must not have picked",
					workflowID, donName, container)
			}
		}
	}
}
