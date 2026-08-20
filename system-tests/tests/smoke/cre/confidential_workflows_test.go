package cre

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-confidential-compute/tests/testhelpers"
	cctypes "github.com/smartcontractkit/chainlink-confidential-compute/types"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/confidentialcompute"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	// confidentialWorkflowsConfigPath is the topology this test runs against,
	// relative to the CRE environment directory.
	confidentialWorkflowsConfigPath = "/configs/workflow-gateway-capabilities-don-confidential-workflows.toml"

	// confidentialWorkflowName is the on-chain workflow name.
	confidentialWorkflowName = "confidential-workflows-e2e"

	// confidentialEchoURL is the outbound target the workflow fetches from inside
	// the enclave. The enclave's default egress policy allows public HTTPS.
	confidentialEchoURL = "https://postman-echo.com/post"

	// confidentialSecretName is the vault secret the workflow reads via GetSecret.
	confidentialSecretName = "MOCK_SECRET"

	// confidentialVaultThreshold matches the 4-node F=1 vault DON.
	confidentialVaultThreshold = 1

	// confidentialDONPageLimit bounds the getDONs page read when locating a
	// capability's DON. The topology has a handful of DONs, so one page covers it.
	confidentialDONPageLimit = 100
)

// Test_CRE_V2_ConfidentialWorkflows_Relay exercises the confidential workflows
// engine path end to end:
//
//	syncer -> ConfidentialModule -> confidential-workflows capability -> enclave
//	-> WASM (cron trigger) -> GetSecret (remote dispatch to the vault DON via the
//	confidential relay) + http.SendRequest (intercepted and executed in-enclave)
//
// The enclaves run locally. On a Nitro-capable host they are real Nitro
// enclaves; elsewhere (including the CRE CI runners) the harness falls back to
// fake enclaves, which run the same binaries as local processes over emulated
// vsock. Attestation validation is relaxed only in the fake case.
//
// The chain-write leg of the workflow (ReportFromDon + evm.WriteReport) is left
// disabled here: it needs a deployed report receiver, and the secret and HTTP
// legs are what prove the relay and enclave routing work.
func Test_CRE_V2_ConfidentialWorkflows_Relay(t *testing.T) {
	testLogger := framework.L

	t.Run("Confidential Workflows Relay - "+topology, func(t *testing.T) {
		// Resolved first so the test skips before standing up anything when the
		// confidential-compute checkout is not available.
		ccRoot := confidentialComputeRoot(t)

		fake := testhelpers.UseFakeEnclave()
		testLogger.Info().Bool("fakeEnclaves", fake).Msg("Starting confidential workflows relay test")

		tconf := t_helpers.GetTestConfig(t, confidentialWorkflowsConfigPath)
		t.Setenv("CTF_CONFIGS", tconf.EnvironmentConfigPath)

		// 1. Stand up the host-side services the enclaves need to reach. Their
		//    addresses have to be known before the enclaves start, because they are
		//    baked into the settings the enclave receives at startup. The gateway
		//    URL is not known until the CRE environment is up, hence the proxy.
		gwProxy := newDeferredGatewayProxy(t, confidentialGatewayProxyPort)
		enclaveHost := confidentialEnclaveHostAddr(fake)
		storageAddr, storageSvc := startFakeStorageService(t, enclaveHost)
		// Both env vars below configure the enclave *host servers*, not this repo:
		// the harness launches them with this process's environment inherited
		// (testhelpers.MustSetupEnclaveWithEnv appends to os.Environ())
		// REQUIRE_BFT_QUORUM makes each host demand a 2f+1 BFT supermajority of node
		// signatures instead of f+1 (enclave/nitro/host), matching the relay's
		// requireBFTQuorum = true in the topology TOML.
		t.Setenv("REQUIRE_BFT_QUORUM", "true")
		t.Setenv("ENCLAVE_SETTINGS", fmt.Sprintf(
			`{"storageKey":%q,"storageServiceUrl":%q,"storageServiceTls":false,"gatewayUrl":%q}`,
			confidentialStorageKeyHex,
			storageAddr,
			fmt.Sprintf("http://%s:%d", enclaveHost, confidentialGatewayProxyPort),
		))

		// 2. Start the enclaves. This is the whole point of depending on
		//    chainlink-confidential-compute's harness from this repository.
		enclaveCfg := testhelpers.DefaultLocalEnclaveSetupConfig(ccRoot, confidentialWorkflowsApp)
		enclaveCfg.Region = confidentialEnclaveRegion
		enclaves := testhelpers.SetupLocalEnclaves(t, enclaveCfg)
		t.Cleanup(enclaves.CleanupAll)
		testLogger.Info().
			Str("hostIP", enclaves.HostIP).
			Int("count", len(enclaves.Enclaves)).
			Msg("Local enclaves ready")

		// 3. Build the environment the standard way. The capability and the relay
		//    come from the standard feature set, configured by the topology TOML,
		//    so nothing here has to be injected as Go values.
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, tconf)

		// 4. Publish the enclave list to the capability's on-chain registry config,
		//    which the relay handler reads to decide where to route. The capability
		//    registers with an empty list and refreshes from the registry on a timer,
		//    so this can land after the environment is already running.
		publishEnclaves(t, testEnv, testLogger, enclaves.Enclaves)

		// 5. Point the proxy at the real gateway now that it exists.
		gatewayURL := confidentialGatewayURL(t, testEnv)
		require.NoError(t, gwProxy.SetTarget(gatewayURL), "failed to set gateway proxy target")
		testLogger.Info().Str("gatewayURL", gatewayURL).Msg("Gateway proxy target set")

		// 6. The vault DON only serves its public key once DKG has produced a
		//    result package on every worker. Fetching before then times out, and
		//    how long DKG takes tracks how loaded the runner is.
		ensureVaultDKGResultPackages(t, testEnv)

		// 6a. The engine's pre-enclave secret fetch reads VaultPublicKey and
		//     Threshold from the vault capability's registry config, which is
		//     registered empty. Without this, GetSecret fails inside the workflow.
		vaultPublicKey := injectVaultPublicKey(t, testEnv, testLogger, gatewayURL)

		// 6b. The enclaves boot with no signer set and no master public key, so they
		//     reject every compute request until this lands.
		configureEnclaves(t, testEnv, testLogger, enclaves.ConfigURLs, vaultPublicKey)

		// 6c. Store the secret the workflow reads via GetSecret. Without it the
		//     request really reaches the vault DON and really comes back empty.
		storeConfidentialWorkflowSecret(t, testEnv, testLogger, gatewayURL, vaultPublicKey,
			confidentialSecretName, "s3cret-from-vault")

		// 7. Compile and serve the workflow from the confidential-compute checkout.
		//    ConsumerAddress is left empty, which the workflow treats as "skip the
		//    chain-write leg".
		configJSON := fmt.Sprintf(`{"echo_url":%q}`, confidentialEchoURL)
		artifacts := buildAndServeConfidentialWorkflow(t, ccRoot, configJSON, testhelpers.DetectHostIP())
		testLogger.Info().
			Str("binaryURL", artifacts.BinaryURL).
			Str("configURL", artifacts.ConfigURL).
			Msg("Workflow artifacts served")

		// The artifact server binds 0.0.0.0, but the enclave reaches it at a
		// different host than Docker does, so swap only the host portion.
		parsed, pErr := url.Parse(artifacts.BinaryURL)
		require.NoError(t, pErr, "parsing workflow binary URL")
		storageSvc.setURL(fmt.Sprintf("http://%s:%s%s", enclaveHost, parsed.Port(), parsed.Path))

		// 8. The syncer reads the binary and config from disk (see the topology's
		//    CRE.WorkflowFetcher override), so copy them into the containers.
		copyWorkflowArtifactsToContainers(t, testEnv, artifacts)

		// 9. Register the workflow as confidential and wait for a successful
		//    execution. The workflow returns an error if either GetSecret or the
		//    in-enclave HTTP fetch fails, so a successful execution implies the
		//    whole relay + enclave path worked.
		workflowID := registerConfidentialWorkflow(t, testEnv, testLogger, artifacts)
		waitForConfidentialWorkflowExecution(t, testEnv, testLogger, workflowID, 5*time.Minute)

		testLogger.Info().Msg("Confidential workflows relay E2E passed")
	})
}

// confidentialEnclaveHostAddr is the address the enclave reaches host-local test
// servers at: loopback for fake enclaves (ordinary local processes), the wg0 host
// IP for real Nitro enclaves.
func confidentialEnclaveHostAddr(fake bool) string {
	if fake {
		return "localhost"
	}
	return "100.64.0.3"
}

// confidentialComputeRoot resolves the chainlink-confidential-compute checkout the
// enclave harness builds and runs the enclave from. The test is skipped when it is
// not set, so a local `go test ./...` does not fail on a missing checkout.
func confidentialComputeRoot(t *testing.T) string {
	t.Helper()

	root := os.Getenv("CONFIDENTIAL_COMPUTE_ROOT")
	if root == "" {
		t.Skip("CONFIDENTIAL_COMPUTE_ROOT must point at a chainlink-confidential-compute checkout; " +
			"CI sets this from the confidential-workflows gitRef in plugins/plugins.public.yaml")
	}

	abs, err := filepath.Abs(root)
	require.NoError(t, err, "resolving CONFIDENTIAL_COMPUTE_ROOT")
	return abs
}

// confidentialGatewayURL builds the externally reachable gateway URL.
func confidentialGatewayURL(t *testing.T, testEnv *ttypes.TestEnvironment) string {
	t.Helper()

	require.NotEmpty(t, testEnv.Dons.GatewayConnectors.Configurations, "no gateway connector configurations")
	incoming := testEnv.Dons.GatewayConnectors.Configurations[0].Incoming
	host := incoming.Host
	if host == "" {
		host = testhelpers.DetectHostIP()
	}
	return fmt.Sprintf("%s://%s:%d%s", incoming.Protocol, host, incoming.ExternalPort, incoming.Path)
}

// publishEnclaves writes the enclave list into the capability's on-chain
// registry config and waits for the capability to pick it up.
//
// The capability refreshes from the registry on a ticker
// (DefaultEnclaveRefreshIntervalSeconds, 10s), so this waits two intervals
// rather than one: a single interval races a refresh that started just before
// the transaction landed and therefore read the old config.
func publishEnclaves(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	enclaves []cctypes.Enclave,
) {
	t.Helper()

	ctx := t.Context()

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	capRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.CapabilitiesRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
		"",
	)

	config, err := confidentialcompute.MarshalRegistryConfig(enclaves)
	require.NoError(t, err, "failed to encode enclave list for the registry")

	donName := donNameForCapability(t, sethClient, capRegAddr, confidentialWorkflowsApp)

	require.NoError(t,
		crelib.UpdateDONCapabilityConfig(ctx, sethClient, capRegAddr, donName, confidentialWorkflowsApp, config),
		"failed to publish enclave list to the capabilities registry",
	)

	testLogger.Info().
		Str("don", donName).
		Int("count", len(enclaves)).
		Dur("wait", confidentialEnclaveRefreshWait).
		Msg("Published enclave list; waiting for the capability to refresh from the registry")
	time.Sleep(confidentialEnclaveRefreshWait)
}

// donNameForCapability returns the registry name of the DON providing the named
// capability. The registry derives DON names from the topology, so the name is
// resolved rather than assumed: a wrong name makes getDONByName revert with an
// opaque custom error, whereas this reports which DONs actually exist.
func donNameForCapability(
	t *testing.T,
	sethClient *seth.Client,
	capabilitiesRegistryAddr string,
	capabilityName string,
) string {
	t.Helper()

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capabilitiesRegistryAddr), sethClient.Client,
	)
	require.NoError(t, err, "failed to create capabilities registry wrapper")

	allDONs, err := capReg.GetDONs(&bind.CallOpts{Context: t.Context()}, big.NewInt(0), big.NewInt(confidentialDONPageLimit))
	require.NoError(t, err, "failed to list DONs from the capabilities registry")

	names := make([]string, 0, len(allDONs))
	for i := range allDONs {
		names = append(names, allDONs[i].Name)
		// The registry keys capabilities as "name@version".
		for _, capabilityConfig := range allDONs[i].CapabilityConfigurations {
			if strings.HasPrefix(capabilityConfig.CapabilityId, capabilityName+"@") {
				return allDONs[i].Name
			}
		}
	}

	require.FailNowf(t, "capability is not registered on any DON",
		"no DON provides capability %q; DONs present: %s", capabilityName, strings.Join(names, ", "))

	return ""
}

// injectVaultPublicKey writes the vault DON's DKG public key and threshold into
// the vault capability's registry config.
func injectVaultPublicKey(t *testing.T, testEnv *ttypes.TestEnvironment, testLogger zerolog.Logger, gatewayURL string) string {
	t.Helper()

	ctx := t.Context()
	vaultPublicKey, err := creworkflow.FetchVaultPublicKey(ctx, gatewayURL)
	require.NoError(t, err, "failed to fetch vault public key from gateway")

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	capRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.CapabilitiesRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
		"",
	)

	vaultDON, _, err := crelib.GetVaultCapabilityDON(ctx, sethClient, capRegAddr)
	require.NoError(t, err, "failed to locate vault capability DON in registry")

	require.NoError(t,
		creworkflow.UpdateVaultCapabilityConfig(ctx, sethClient, capRegAddr, vaultDON, vaultPublicKey, confidentialVaultThreshold),
		"failed to inject VaultPublicKey/Threshold into the vault capability config")
	testLogger.Info().Msg("Injected VaultPublicKey + Threshold into the vault capability config")

	return vaultPublicKey
}

// copyWorkflowArtifactsToContainers copies the compiled binary and its config into
// every workflow DON container so the syncer's file fetcher can read them.
func copyWorkflowArtifactsToContainers(t *testing.T, testEnv *ttypes.TestEnvironment, artifacts confidentialWorkflowArtifacts) {
	t.Helper()

	for _, don := range testEnv.Dons.List() {
		if !don.HasFlag(crelib.WorkflowDON) {
			continue
		}
		for _, filename := range []string{confidentialWorkflowBinaryFilename, confidentialWorkflowConfigFilename} {
			require.NoError(t,
				creworkflow.CopyArtifactsToDockerContainers(
					creworkflow.DefaultWorkflowTargetDir,
					ns.NodeNamePrefix(don.Name),
					filepath.Join(artifacts.ArtifactDir, filename),
				),
				"failed to copy %s to the %s DON containers", filename, don.Name)
		}
	}
}

// registerConfidentialWorkflow registers the workflow on-chain with confidential
// attributes and returns its workflow ID.
func registerConfidentialWorkflow(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	artifacts confidentialWorkflowArtifacts,
) string {
	t.Helper()

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	wfRegistryRef := crecontracts.MustGetAddressRefFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)

	// The confidential attribute is what routes execution into the enclave rather
	// than running the WASM on the workflow DON.
	attributes := []byte(`{"confidential":true}`)
	configURL := artifacts.ConfigURL

	workflowID, err := creworkflow.RegisterWithContract(
		context.Background(),
		sethClient,
		common.HexToAddress(wfRegistryRef.Address),
		wfRegistryRef.Version,
		0, // donID unused for v2
		testEnv.Dons.MustWorkflowDON().DonFamily,
		confidentialWorkflowName,
		workflowTag,
		artifacts.BinaryURL,
		&configURL,
		nil, // no secrets URL
		attributes,
		nil, // keep the HTTP URL on-chain; the enclave fetches the binary itself
	)
	require.NoError(t, err, "failed to register confidential workflow")
	testLogger.Info().Str("workflowID", workflowID).Msg("Confidential workflow registered")

	t.Cleanup(func() {
		_ = creworkflow.DeleteWithContract(
			context.Background(),
			sethClient,
			common.HexToAddress(wfRegistryRef.Address),
			wfRegistryRef.Version,
			confidentialWorkflowName,
		)
	})

	return workflowID
}

// waitForConfidentialWorkflowExecution waits for the engine to log a successful
// execution for this workflow. The engine emits that line once per successful
// trigger execution, not for the Subscribe-phase call at engine startup, so
// finding it means the cron trigger fired and the whole enclave path succeeded.
func waitForConfidentialWorkflowExecution(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	workflowID string,
	timeout time.Duration,
) {
	t.Helper()

	containers := confidentialWorkflowDONContainers(testEnv)
	require.NotEmpty(t, containers, "no workflow DON containers found to scrape")

	needleMsg := []byte(`"msg":"Workflow execution finished successfully"`)
	needleID := []byte(workflowID)
	testLogger.Info().
		Str("workflowID", workflowID).
		Strs("containers", containers).
		Msg("Waiting for a successful workflow execution")

	deadline := time.Now().Add(timeout)
	for {
		for _, name := range containers {
			out, _ := exec.CommandContext(t.Context(), "docker", "logs", "--tail", "10000", name).CombinedOutput()
			for line := range bytes.SplitSeq(out, []byte{'\n'}) {
				if bytes.Contains(line, needleMsg) && bytes.Contains(line, needleID) {
					testLogger.Info().Str("container", name).Msg("Found successful execution log")
					return
				}
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out after %s waiting for a successful execution of workflow %s\n%s",
				timeout, workflowID, confidentialExecutionDiagnostics(t, containers))
		}
		time.Sleep(5 * time.Second)
	}
}

// confidentialExecutionDiagnostics summarises why no execution succeeded, so a
// timeout reports the node-side cause instead of only the absence of a success
// log. Returns one deduplicated line per distinct error.
func confidentialExecutionDiagnostics(t *testing.T, containers []string) string {
	t.Helper()

	needles := [][]byte{
		[]byte("Workflow Engine initialization failed"),
		[]byte("Workflow execution failed"),
		[]byte("failed to get regions from"),
		[]byte("no compatible capability found"),
	}

	seen := map[string]bool{}
	var found []string
	for _, name := range containers {
		out, _ := exec.CommandContext(t.Context(), "docker", "logs", "--tail", "10000", name).CombinedOutput()
		for line := range bytes.SplitSeq(out, []byte{'\n'}) {
			for _, needle := range needles {
				if !bytes.Contains(line, needle) {
					continue
				}
				// Key on the message alone; every node logs the same failure.
				key := string(needle)
				if !seen[key] {
					seen[key] = true
					found = append(found, fmt.Sprintf("  [%s] %s", name, truncateForLog(line, 400)))
				}
				break
			}
		}
	}

	if len(found) == 0 {
		return "no workflow engine or execution errors found in the node logs"
	}

	return "node-side errors:\n" + strings.Join(found, "\n")
}

// truncateForLog shortens a log line so a failure message stays readable.
func truncateForLog(line []byte, maxLen int) string {
	if len(line) <= maxLen {
		return string(line)
	}

	return string(line[:maxLen]) + "... (truncated)"
}

// confidentialWorkflowDONContainers returns the chainlink container names for
// every nodeset whose DON carries the workflow DON flag.
func confidentialWorkflowDONContainers(testEnv *ttypes.TestEnvironment) []string {
	workflowDONNames := map[string]bool{}
	for _, don := range testEnv.Dons.List() {
		if don.HasFlag(crelib.WorkflowDON) {
			workflowDONNames[don.Name] = true
		}
	}

	var names []string
	for _, nodeSet := range testEnv.Config.NodeSets {
		if !workflowDONNames[nodeSet.Name] || nodeSet.Out == nil {
			continue
		}
		for _, cl := range nodeSet.Out.CLNodes {
			if cl == nil || cl.Node == nil || cl.Node.ContainerName == "" {
				continue
			}
			names = append(names, cl.Node.ContainerName)
		}
	}
	return names
}
