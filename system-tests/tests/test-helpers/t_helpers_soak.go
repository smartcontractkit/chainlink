package helpers

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

type workflowCleanupEntry struct {
	name       string
	configPath string
	wasmPath   string
	sethClient *seth.Client
}

// deleteWorkflowUnsafe is the soak/NTimes-only counterpart to deleteWorkflows: same local + registry
// teardown and retry policy, but no deleteWorkflowsMu so independent deploy keys can delete in parallel.
func deleteWorkflowUnsafe(
	ctx context.Context,
	uniqueWorkflowName string,
	workflowConfigFilePath string,
	compressedWorkflowWasmPath string,
	workflowRegistryAddress common.Address,
	version *semver.Version,
	sethClient *seth.Client,
) error {
	testLogger := framework.L
	testLogger.Info().Msgf("Deleting workflow artifacts (%s) after test.", uniqueWorkflowName)
	if localEnvErr := creworkflow.RemoveWorkflowArtifactsFromLocalEnv(workflowConfigFilePath, compressedWorkflowWasmPath); localEnvErr != nil {
		return errors.Wrap(localEnvErr, "failed to remove workflow artifacts from local environment")
	}

	retryErr := retry.Do(func() error {
		return creworkflow.DeleteWithContract(ctx, sethClient, workflowRegistryAddress, version, uniqueWorkflowName)
	}, retry.Attempts(3), retry.Delay(1*time.Second), retry.DelayType(retry.BackOffDelay), retry.RetryIf(func(err error) bool {
		return strings.Contains(err.Error(), ReentrancySentryOOGError)
	}), retry.OnRetry(func(n uint, err error) {
		testLogger.Error().Msgf("Error deleting workflow '%s': %s", uniqueWorkflowName, err.Error())
	}))
	if retryErr != nil {
		return errors.Wrapf(retryErr, "failed to delete workflow '%s' from registry", uniqueWorkflowName)
	}
	testLogger.Info().Msgf("Workflow '%s' deleted successfully from the registry.", uniqueWorkflowName)
	return nil
}

// deleteWorkflowsNTimesCleanup deletes workflows registered by CompileAndDeployWorkflowNTimes.
// Parallel across deploy keys (one goroutine per key, matching registration); serial within each key
// to avoid nonce collisions on the same Seth client.
func deleteWorkflowsNTimesCleanup(
	t *testing.T,
	byKey [][]workflowCleanupEntry,
	workflowRegistryAddress common.Address,
	version *semver.Version,
) {
	t.Helper()
	if len(byKey) == 0 {
		return
	}

	var eg errgroup.Group
	eg.SetLimit(envVarOrDefault("CRE_TEST_DEPLOY_MAX_PARALLEL", defaultDeployMaxParallel))
	for _, entries := range byKey {
		if len(entries) == 0 {
			continue
		}
		eg.Go(func() error {
			for _, e := range entries {
				ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
				if err := deleteWorkflowUnsafe(ctx, e.name, e.configPath, e.wasmPath, workflowRegistryAddress, version, e.sethClient); err != nil {
					cancel()
					return err
				}
				cancel()
			}
			return nil
		})
	}
	require.NoError(t, eg.Wait(), "parallel workflow cleanup")
}

// CompileAndDeployWorkflowNTimes compiles the workflow once, provisions numKeys funded signers via
// ConfigureAdditionalWorkflowSigners, then registers numDeployments workflows in parallel across keys
// (serial per key). Artifacts are packed into one tarball per DON. Requires numDeployments >= 2;
// use CompileAndDeployWorkflow for a single deployment.
func CompileAndDeployWorkflowNTimes[T WorkflowConfig](t *testing.T,
	testEnv *ttypes.TestEnvironment, testLogger zerolog.Logger,
	workflowNameFn func(i int) string,
	workflowConfigFn func(i int) *T,
	workflowFileLocation string,
	numDeployments, numKeys int,
	opts ...CompileAndDeployWorkflowOpt,
) []string {
	t.Helper()
	require.Greater(t, numDeployments, 1, "use CompileAndDeployWorkflow for a single deployment")
	require.NotNil(t, workflowConfigFn, "workflowConfigFn is required")
	require.Positive(t, numKeys, "numKeys must be positive")

	name0 := workflowNameFn(0)
	require.GreaterOrEqual(t, len(name0), 10, "workflow name for compile must be at least 10 characters")

	cfg := compileAndDeployWorkflowCfg{
		artifactCopyDONTypes: []cre.CapabilityFlag{cre.WorkflowDON},
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	testLogger.Info().
		Int("num_deployments", numDeployments).
		Int("num_keys", numKeys).
		Int("max_parallel", envVarOrDefault("CRE_TEST_DEPLOY_MAX_PARALLEL", defaultDeployMaxParallel)).
		Str("workflow_file_location", workflowFileLocation).
		Msg("compiling once and registering multiple workflows")

	sharedEnv := getOrCreateSharedEnvironment(t, testEnv.TestConfig)
	if testEnv.Execution == nil {
		testEnv.Execution = &ttypes.ExecutionContext{}
	}
	if testEnv.Execution.TestID == "" {
		testEnv.Execution.TestID = deriveExecutionTestID(t)
	}

	deployKeys := configureAdditionalWorkflowSigners(t, sharedEnv, testEnv, numKeys)
	require.Len(t, deployKeys, numKeys, "deploy keys count must match numKeys")

	artifactDir := workflowArtifactsDir(t, testEnv)
	workflowDONs := selectArtifactTargetDONs(testEnv, cfg.artifactCopyDONTypes)
	wasmPaths, configPaths := buildNTimesWorkflowArtifacts(t, testLogger, workflowNameFn, workflowConfigFn, workflowFileLocation, artifactDir, numDeployments)
	copyNTimesArtifactsToDONs(t, testLogger, workflowDONs, artifactDir, wasmPaths, configPaths, numDeployments)

	registryChainSelector := testEnv.CreEnvironment.Blockchains[0].ChainSelector()
	workflowRegistryAddress := crecontracts.MustGetAddressRefFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		registryChainSelector,
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain type")

	registryAddr := common.HexToAddress(workflowRegistryAddress.Address)
	registryVersion := workflowRegistryAddress.Version
	donID := testEnv.Dons.MustWorkflowDON().ID

	ids := registerNTimesWorkflowsParallel(
		t.Context(), t, testLogger, workflowNameFn, workflowFileLocation,
		wasmPaths, configPaths, numDeployments, numKeys, deployKeys,
		registryAddr, registryVersion, registryChainSelector, donID, cfg.attributes,
	)

	cleanupByKey := buildNTimesCleanupByKey(workflowNameFn, configPaths, wasmPaths, numDeployments, numKeys, deployKeys)
	t.Cleanup(func() {
		deleteWorkflowsNTimesCleanup(t, cleanupByKey, registryAddr, registryVersion)
	})
	return ids
}

// configureAdditionalWorkflowSigners creates numSigners distinct funded keys and registry-chain seth clients,
// authorizes each on the v2 workflow registry when needed, and returns them. It does not mutate testEnv
// blockchains or CLDF (use configurePerTestExecutionContext when the test should own the env with one signer).
func configureAdditionalWorkflowSigners(t *testing.T, sharedEnv *ttypes.TestEnvironment, testEnv *ttypes.TestEnvironment, numSigners int) []ttypes.PerTestDeployKey {
	t.Helper()
	require.Positive(t, numSigners, "numSigners must be at least 1")

	registryChainSelector := testEnv.CreEnvironment.Blockchains[0].ChainSelector()
	evmChains := make(map[uint64]*evm.Blockchain)
	for _, bcOutput := range sharedEnv.CreEnvironment.Blockchains {
		evmChain, ok := bcOutput.(*evm.Blockchain)
		if !ok {
			continue
		}
		evmChains[evmChain.ChainSelector()] = evmChain
	}

	out := make([]ttypes.PerTestDeployKey, 0, numSigners)
	for keyIdx := range numSigners {
		ownerAddress, privateKey, addrErr := crecrypto.GenerateNewKeyPair()
		require.NoError(t, addrErr, "failed to generate workflow signer key pair")
		privateKeyHex := hex.EncodeToString(gethcrypto.FromECDSA(privateKey))

		var registryClient *seth.Client

		for chainSelector, evmChain := range evmChains {
			wsURL := evmChain.WSURL()
			require.NotEmptyf(t, wsURL, "missing WS URL for chain selector %d", chainSelector)

			perTestClient, clientErr := seth.NewClientBuilder().
				WithRpcUrl(wsURL).
				WithPrivateKeys([]string{privateKeyHex}).
				WithProtections(false, false, seth.MustMakeDuration(time.Second)).
				Build()
			require.NoErrorf(t, clientErr, "failed to create per-signer seth client for selector %d", chainSelector)

			rootSignerNonceLock.Lock()
			require.NoError(
				t,
				evmChain.Fund(t.Context(), ownerAddress.Hex(), perTestEVMFundingAmountWei),
				"failed to fund workflow signer %s on chain selector %d",
				ownerAddress.Hex(),
				chainSelector,
			)
			rootSignerNonceLock.Unlock()

			if evmChain.ChainSelector() == registryChainSelector {
				registryClient = perTestClient
			}
		}

		require.NotNil(t, registryClient, "failed to build registry chain seth client for signer %d", keyIdx)
		authorizePerTestWorkflowSignerIfNeeded(t, sharedEnv, ownerAddress)
		out = append(out, ttypes.PerTestDeployKey{
			OwnerAddress:   ownerAddress,
			RegistryClient: registryClient,
		})
	}

	return out
}

func buildNTimesWorkflowArtifacts[T WorkflowConfig](
	t *testing.T,
	testLogger zerolog.Logger,
	workflowNameFn func(i int) string,
	workflowConfigFn func(i int) *T,
	workflowFileLocation, artifactDir string,
	numDeployments int,
) (wasmPaths []string, configPaths []string) {
	t.Helper()

	configPaths = make([]string, numDeployments)
	for i := range numDeployments {
		configPaths[i] = workflowConfigFactory(t, testLogger, workflowNameFn(i), workflowConfigFn(i), artifactDir)
	}

	name0 := workflowNameFn(0)
	compressedWasmPath, compileErr := creworkflow.CompileWorkflowToDir(t.Context(), workflowFileLocation, name0, artifactDir)
	require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)
	testLogger.Info().Msg("Workflow compiled successfully.")

	wasmBytes, readWasmErr := os.ReadFile(compressedWasmPath)
	require.NoError(t, readWasmErr, "failed to read compiled workflow wasm artifact")

	wasmPaths = make([]string, numDeployments)
	wasmPaths[0] = compressedWasmPath
	for i := range numDeployments {
		wasmPaths[i] = writeNTimesWasmCopy(t, artifactDir, workflowNameFn(i), wasmBytes)
	}
	return wasmPaths, configPaths
}

func writeNTimesWasmCopy(t *testing.T, artifactDir, workflowName string, wasmBytes []byte) string {
	t.Helper()
	dst := filepath.Join(artifactDir, workflowName+".br.b64")
	require.NoError(t, os.WriteFile(dst, wasmBytes, 0o644), "failed to write per-workflow wasm copy for %q", workflowName) //nolint:gosec // G306: test artifact
	absDst, absErr := filepath.Abs(dst)
	require.NoError(t, absErr, "failed to resolve wasm path for %q", workflowName)
	return absDst
}

func copyNTimesArtifactsToDONs(
	t *testing.T,
	testLogger zerolog.Logger,
	workflowDONs []*cre.Don,
	artifactDir string,
	wasmPaths, configPaths []string,
	numDeployments int,
) {
	t.Helper()

	pathsForTar := make([]string, 0, numDeployments*2)
	for i := range numDeployments {
		pathsForTar = append(pathsForTar, wasmPaths[i])
		if configPaths[i] != "" {
			pathsForTar = append(pathsForTar, configPaths[i])
		}
	}
	bundlePath := filepath.Join(artifactDir, "workflow-artifacts-bundle.tar")
	require.NoError(t, creworkflow.WriteArtifactTarball(bundlePath, pathsForTar), "failed to write workflow artifact tarball")
	testLogger.Info().Str("bundle", bundlePath).Int("files", len(pathsForTar)).Msg("Copying workflow artifact tarball to DONs")
	for _, don := range workflowDONs {
		copyErr := creworkflow.CopyAndExtractTarballToDockerContainers(
			creworkflow.DefaultWorkflowTargetDir,
			ns.NodeNamePrefix(don.Name),
			bundlePath,
		)
		require.NoError(t, copyErr, "failed to copy and extract artifact tarball for DON %q", don.Name)
	}
}

func registerNTimesWorkflowsParallel(
	ctx context.Context,
	t *testing.T,
	testLogger zerolog.Logger,
	workflowNameFn func(i int) string,
	workflowFileLocation string,
	wasmPaths, configPaths []string,
	numDeployments, numKeys int,
	deployKeys []ttypes.PerTestDeployKey,
	registryAddr common.Address,
	registryVersion *semver.Version,
	registryChainSelector, donID uint64,
	attributes []byte,
) []string {
	t.Helper()

	ids := make([]string, numDeployments)
	var eg errgroup.Group
	eg.SetLimit(envVarOrDefault("CRE_TEST_DEPLOY_MAX_PARALLEL", defaultDeployMaxParallel))

	for keyIdx := range numKeys {
		eg.Go(func() error {
			sc := deployKeys[keyIdx].RegistryClient
			for i := keyIdx; i < numDeployments; i += numKeys {
				name := workflowNameFn(i)
				wfRegCfg := &WorkflowRegistrationConfig{
					WorkflowName:            name,
					WorkflowLocation:        workflowFileLocation,
					ConfigFilePath:          configPaths[i],
					CompressedWasmPath:      wasmPaths[i],
					WorkflowRegistryAddr:    registryAddr,
					WorkflowRegistryVersion: registryVersion,
					ChainID:                 registryChainSelector,
					DonID:                   donID,
					ContainerTargetDir:      creworkflow.DefaultWorkflowTargetDir,
					SethClient:              sc,
					Attributes:              attributes,
				}

				workflowID, regErr := registerWorkflowErr(ctx, wfRegCfg, sc)
				if regErr != nil {
					return fmt.Errorf("register workflow %q (index %d): %w", name, i, regErr)
				}
				ids[i] = workflowID
				testLogger.Info().Str("workflow_id", workflowID).Str("workflow_name", name).Int("index", i).Msg("Workflow registered successfully")
			}
			return nil
		})
	}
	require.NoError(t, eg.Wait(), "parallel workflow registration")
	return ids
}

func buildNTimesCleanupByKey(
	workflowNameFn func(i int) string,
	configPaths, wasmPaths []string,
	numDeployments, numKeys int,
	deployKeys []ttypes.PerTestDeployKey,
) [][]workflowCleanupEntry {
	cleanupByKey := make([][]workflowCleanupEntry, numKeys)
	for keyIdx := range numKeys {
		sc := deployKeys[keyIdx].RegistryClient
		for i := keyIdx; i < numDeployments; i += numKeys {
			cleanupByKey[keyIdx] = append(cleanupByKey[keyIdx], workflowCleanupEntry{
				name:       workflowNameFn(i),
				configPath: configPaths[i],
				wasmPath:   wasmPaths[i],
				sethClient: sc,
			})
		}
	}
	return cleanupByKey
}
