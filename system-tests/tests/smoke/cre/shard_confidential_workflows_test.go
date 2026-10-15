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

	shardedConfidentialSchedule          = "*/30 * * * * *"
	shardedConfidentialNamespace         = "main"
	shardedConfidentialExecTimeout       = 5 * time.Minute
	confidentialWorkflowSuccessLogNeedle = `"msg":"Workflow execution finished successfully"`

	shardZeroEnclaveBaseCID = 16
	shardOneEnclaveBaseCID  = 18
)

var (
	shardZeroEnclaveHTTPPorts   = []string{"8080", "8081"}
	shardZeroEnclaveConfigPorts = []string{"8082", "8083"}
	shardOneEnclaveHTTPPorts    = []string{"8084", "8085"}
	shardOneEnclaveConfigPorts  = []string{"8086", "8087"}
)

// ExecuteShardManualAssignmentConfidentialWorkflowsTest pins the owner to shard 0
// (the shard the confidential relay serves) with the static default pointing at
// shard 1, so a silently missed per_owner entry can only fail the test.
func ExecuteShardManualAssignmentConfidentialWorkflowsTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	ccRoot := confidentialComputeRoot(t)

	shards := mustSharedVaultShardPair(t, testEnv, shardedConfidentialTopologyHint)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "manual-sharded-confidential-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	enclaveHost := confidentialEnclaveHostAddr(testhelpers.UseFakeEnclave())
	hostServices := startConfidentialHostServices(t, enclaveHost)
	hostServices.targetGateway(t, testEnv)

	setupShardedConfidentialEnclaves(t, testEnv, testLogger, ccRoot, shards, fixture)

	assignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d]
`, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx)
	proposeSharedVaultAssignment(t, testEnv, assignmentTOML)

	workflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, hostServices.storageSvc, enclaveHost,
		"manual-shard0-confidential", secretKey, secretValue)

	waitForConfidentialWorkflowExecution(t, testEnv, testLogger, workflowID, shardedConfidentialExecTimeout, shards.shardZeroDON.Name)
	assertNoConfidentialExecutionOnNodeset(t, testEnv, workflowID, shards.shardOneDON.Name)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle)
}

// ExecuteShardFailoverConfidentialWorkflowsTest runs [shard 0, shard 1] with the
// secret leg on shard 0, then swaps to [shard 1, shard 0] with the secret leg
// dropped: the enclave -> relay secret path is bound to the first shard today,
// so a shard-1 execution cannot fetch secrets through shard-0's relay nodes.
func ExecuteShardFailoverConfidentialWorkflowsTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	ccRoot := confidentialComputeRoot(t)

	shards := mustSharedVaultShardPair(t, testEnv, shardedConfidentialTopologyHint)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-sharded-confidential-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	enclaveHost := confidentialEnclaveHostAddr(testhelpers.UseFakeEnclave())
	hostServices := startConfidentialHostServices(t, enclaveHost)
	hostServices.targetGateway(t, testEnv)

	setupShardedConfidentialEnclaves(t, testEnv, testLogger, ccRoot, shards, fixture)

	primaryAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, shards.shardZeroIdx, shards.shardOneIdx, shards.workflowOwner, shards.shardZeroIdx, shards.shardOneIdx)
	proposeSharedVaultAssignment(t, testEnv, primaryAssignmentTOML)

	workflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, hostServices.storageSvc, enclaveHost,
		"failover-shard0-confidential", secretKey, secretValue)
	waitForConfidentialWorkflowExecution(t, testEnv, testLogger, workflowID, shardedConfidentialExecTimeout, shards.shardZeroDON.Name)
	assertNoConfidentialExecutionOnNodeset(t, testEnv, workflowID, shards.shardOneDON.Name)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)

	swapWorkflowID := deployShardedConfidentialWorkflow(t, testEnv, testLogger, hostServices.storageSvc, enclaveHost,
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

type confidentialHostServices struct {
	gwProxy    *deferredGatewayProxy
	storageSvc *fakeStorageService
}

func startConfidentialHostServices(t *testing.T, enclaveHost string) *confidentialHostServices {
	t.Helper()

	gwProxy := newDeferredGatewayProxy(t, confidentialGatewayProxyPort)
	storageAddr, storageSvc := startFakeStorageService(t, enclaveHost)

	// Both env vars below configure the enclave *host servers*, not this repo:
	// the harness launches them with this process's environment inherited
	// (testhelpers.MustSetupEnclaveWithEnv appends to os.Environ())
	// REQUIRE_BFT_QUORUM makes each host demand a 2f+1 BFT supermajority of node
	// signatures instead of f+1 (enclave/nitro/host), matching the relay's
	// requireBFTQuorum = true in the topology TOML.
	t.Setenv("REQUIRE_BFT_QUORUM", "true")
	// The vsock outbound proxy blocks loopback and local interface IPs by
	// default; the fake storage service and artifact server are host-local, so
	// opt the test host servers into allowing local destinations.
	t.Setenv("OUTBOUND_ALLOW_LOCAL_FOR_TESTS", "true")
	t.Setenv("ENCLAVE_SETTINGS", fmt.Sprintf(
		`{"storageKey":%q,"storageServiceUrl":%q,"storageServiceTls":false,"gatewayUrl":%q}`,
		confidentialStorageKeyHex,
		storageAddr,
		fmt.Sprintf("http://%s:%d", enclaveHost, confidentialGatewayProxyPort),
	))

	return &confidentialHostServices{gwProxy: gwProxy, storageSvc: storageSvc}
}

func (s *confidentialHostServices) targetGateway(t *testing.T, testEnv *ttypes.TestEnvironment) string {
	t.Helper()
	gatewayURL := confidentialGatewayURL(t, testEnv)
	require.NoError(t, s.gwProxy.SetTarget(gatewayURL), "failed to set gateway proxy target")
	return gatewayURL
}

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
			// The enclave URLs are dialed from the node containers, so they need
			// the Docker-reachable host, not the test host's own address.
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

	artifacts := buildAndServeConfidentialWorkflow(t, shardedConfidentialWorkflowSrc, workflowName, configJSON, enclaveHost)
	storageSvc.setArtifactURL(artifacts.BinaryFilename, artifacts.BinaryURL)

	copyWorkflowArtifactsToContainers(t, testEnv, artifacts)

	return registerConfidentialWorkflow(t, testEnv, testLogger, workflowName, artifacts)
}

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
