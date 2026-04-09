package helpers

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	perTestEVMFundingAmountWei uint64 = 1_000_000_000_000_000_000 // 1 ETH
)

type sharedEnvironmentEntry struct {
	once sync.Once
	env  *ttypes.TestEnvironment
	err  error
}

var (
	sharedEnvMu         sync.Mutex
	sharedEnvironments  = make(map[string]*sharedEnvironmentEntry)
	rootSignerNonceLock sync.Mutex
)

// SetupTestEnvironmentWithConfig creates a test environment backed by the shared
// root deployer/signer context (legacy behavior).
//
// Use this for admin/control-plane and ownership-sensitive tests (e.g. V1
// registry paths, sharding tests) where operations are expected to be executed
// by the root owner.
//
// For parallel workflow-plane tests, use SetupTestEnvironmentWithPerTestKeys.
func SetupTestEnvironmentWithConfig(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()
	return setupTestEnvironmentWithConfigMode(t, tconf, false, flags...)
}

// SetupTestEnvironmentWithPerTestKeys creates a test environment that uses a
// dedicated per-test funded key (seth + CLDF deployer key) for on-chain writes.
//
// Use this for workflow-plane tests that may run in parallel and should avoid
// nonce collisions and shared-key state coupling (e.g. V2 suite bucket flows).
//
// Do NOT use this for admin/control-plane tests that require owner-only access
// on shared contracts (e.g. ShardConfig ownership-admin operations). Use
// SetupTestEnvironmentWithConfig for those.
func SetupTestEnvironmentWithPerTestKeys(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()
	return setupTestEnvironmentWithConfigMode(t, tconf, true, flags...)
}

func setupTestEnvironmentWithConfigMode(t *testing.T, tconf *ttypes.TestConfig, usePerTestKeys bool, flags ...string) *ttypes.TestEnvironment {
	t.Helper()

	sharedEnv := getOrCreateSharedEnvironment(t, tconf, flags...)
	testEnv := cloneSharedEnvironmentForTest(sharedEnv, tconf)
	if usePerTestKeys {
		testEnv.Execution = configurePerTestExecutionContext(t, sharedEnv, testEnv)
	}

	t.Cleanup(func() {
		if t.Failed() {
			framework.L.Warn().Msg("Test failed - checking for panics in Docker containers...")
			foundPanics := infra.CheckContainersForPanics(framework.L, 100)
			if !foundPanics {
				framework.L.Warn().Msgf("No panic patterns detected in Docker container logs")
				infra.PrintFailedContainerLogs(framework.L, 30)
			}
		}
	})

	return testEnv
}

func getOrCreateSharedEnvironment(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()

	key := sharedEnvironmentKey(tconf, flags)
	sharedEnvMu.Lock()
	entry, ok := sharedEnvironments[key]
	if !ok {
		entry = &sharedEnvironmentEntry{}
		sharedEnvironments[key] = entry
	}
	sharedEnvMu.Unlock()

	entry.once.Do(func() {
		createEnvironment(t, tconf, flags...)
		in := getEnvironmentConfig(t)
		creEnvironment, dons, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in)
		if err != nil {
			entry.err = err
			return
		}

		entry.env = &ttypes.TestEnvironment{
			Config:         in,
			TestConfig:     tconf,
			Logger:         framework.L,
			CreEnvironment: creEnvironment,
			Dons:           dons,
		}
	})

	require.NoError(t, entry.err, "failed to load environment")
	require.NotNil(t, entry.env, "shared test environment was not initialized")
	return entry.env
}

func sharedEnvironmentKey(tconf *ttypes.TestConfig, flags []string) string {
	sortedFlags := slices.Clone(flags)
	slices.Sort(sortedFlags)
	return strings.Join(append([]string{tconf.EnvironmentConfigPath}, sortedFlags...), "|")
}

func cloneSharedEnvironmentForTest(sharedEnv *ttypes.TestEnvironment, tconf *ttypes.TestConfig) *ttypes.TestEnvironment {
	clonedCRE := *sharedEnv.CreEnvironment
	if sharedEnv.CreEnvironment.CldfEnvironment != nil {
		clonedCRE.CldfEnvironment = cloneCldfEnvironmentForTest(sharedEnv.CreEnvironment.CldfEnvironment)
	}
	clonedCRE.Blockchains = make([]blockchains.Blockchain, len(sharedEnv.CreEnvironment.Blockchains))
	copy(clonedCRE.Blockchains, sharedEnv.CreEnvironment.Blockchains)

	return &ttypes.TestEnvironment{
		Config:         sharedEnv.Config,
		TestConfig:     tconf,
		Logger:         framework.L,
		CreEnvironment: &clonedCRE,
		Dons:           sharedEnv.Dons,
	}
}

func configurePerTestExecutionContext(t *testing.T, sharedEnv *ttypes.TestEnvironment, testEnv *ttypes.TestEnvironment) *ttypes.ExecutionContext {
	t.Helper()

	ownerAddress, privateKey, addrErr := crecrypto.GenerateNewKeyPair()
	require.NoError(t, addrErr, "failed to generate per-test key pair")
	privateKeyHex := hex.EncodeToString(gethcrypto.FromECDSA(privateKey))

	testID := deriveExecutionTestID(t)
	execCtx := &ttypes.ExecutionContext{
		TestID:       testID,
		OwnerAddress: ownerAddress,
	}

	registryChainSelector := testEnv.CreEnvironment.Blockchains[0].ChainSelector()
	rootEVMChains := make(map[uint64]*evm.Blockchain)
	for _, bcOutput := range sharedEnv.CreEnvironment.Blockchains {
		evmChain, ok := bcOutput.(*evm.Blockchain)
		if !ok {
			continue
		}
		rootEVMChains[evmChain.ChainSelector()] = evmChain
	}

	for i, bcOutput := range testEnv.CreEnvironment.Blockchains {
		evmChain, ok := bcOutput.(*evm.Blockchain)
		if !ok {
			continue
		}

		rootChain, exists := rootEVMChains[evmChain.ChainSelector()]
		require.Truef(t, exists, "missing shared EVM chain for selector %d", evmChain.ChainSelector())
		wsURL := rootChain.WSURL()
		require.NotEmptyf(t, wsURL, "missing WS URL for chain selector %d", evmChain.ChainSelector())

		perTestClient, clientErr := seth.NewClientBuilder().
			WithRpcUrl(wsURL).
			WithPrivateKeys([]string{privateKeyHex}).
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		require.NoErrorf(t, clientErr, "failed to create per-test seth client for selector %d", evmChain.ChainSelector())

		rootSignerNonceLock.Lock()
		require.NoError(
			t,
			rootChain.Fund(t.Context(), ownerAddress.Hex(), perTestEVMFundingAmountWei),
			"failed to fund per-test owner %s on chain selector %d",
			ownerAddress.Hex(),
			evmChain.ChainSelector(),
		)
		rootSignerNonceLock.Unlock()

		testEnv.CreEnvironment.Blockchains[i] = evmChain.CloneWithSethClient(perTestClient)
		deployerKey, txOptsErr := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(perTestClient.ChainID))
		require.NoErrorf(t, txOptsErr, "failed to create deployer key for chain selector %d", evmChain.ChainSelector())
		deployerKey.Context = t.Context() //nolint:fatcontext // false-positive
		require.NoErrorf(
			t,
			setCldfEVMDeployerKey(testEnv.CreEnvironment.CldfEnvironment, evmChain.ChainSelector(), deployerKey),
			"failed to set per-test CLDF deployer key for selector %d",
			evmChain.ChainSelector(),
		)

		if evmChain.ChainSelector() == registryChainSelector {
			execCtx.ChainClient = perTestClient
		}
	}

	authorizePerTestWorkflowSignerIfNeeded(t, sharedEnv, ownerAddress)
	return execCtx
}

func deriveExecutionTestID(t *testing.T) string {
	base := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "-")
	base = strings.ReplaceAll(base, " ", "-")
	base = strings.ReplaceAll(base, ":", "-")
	base = strings.ReplaceAll(base, ".", "-")
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano()%100000)
}

func authorizePerTestWorkflowSignerIfNeeded(t *testing.T, sharedEnv *ttypes.TestEnvironment, signer common.Address) {
	t.Helper()

	registryAddressRef := crecontracts.MustGetAddressRefFromDataStore(
		sharedEnv.CreEnvironment.CldfEnvironment.DataStore,
		sharedEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.WorkflowRegistry.String(),
		sharedEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	if registryAddressRef.Version == nil || registryAddressRef.Version.Major() != 2 {
		return
	}

	rootRegistryChain, ok := sharedEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)
	require.True(t, ok, "expected registry chain to be EVM")

	registry, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(registryAddressRef.Address), rootRegistryChain.SethClient.Client)
	require.NoError(t, err, "failed to instantiate workflow registry v2 contract")

	allowed, err := registry.IsAllowedSigner(rootRegistryChain.SethClient.NewCallOpts(), signer)
	require.NoError(t, err, "failed to check signer allowlist status")
	if allowed {
		return
	}

	rootSignerNonceLock.Lock()
	defer rootSignerNonceLock.Unlock()

	_, err = rootRegistryChain.SethClient.Decode(registry.UpdateAllowedSigners(rootRegistryChain.SethClient.NewTXOpts(), []common.Address{signer}, true))
	require.NoError(t, err, "failed to authorize per-test signer")
}

func GetDefaultTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return GetTestConfig(t, "/configs/workflow-gateway-capabilities-don.toml")
}

func GetTestConfig(t *testing.T, configPath string) *ttypes.TestConfig {
	relativePathToRepoRoot := "../../../../"
	environmentDirPath := filepath.Join(relativePathToRepoRoot, "core/scripts/cre/environment")

	return &ttypes.TestConfig{
		RelativePathToRepoRoot: relativePathToRepoRoot,
		EnvironmentDirPath:     environmentDirPath,
		EnvironmentConfigPath:  filepath.Join(environmentDirPath, configPath), // change to your desired config, if you want to use another topology
		EnvironmentStateFile:   filepath.Join(environmentDirPath, envconfig.StateDirname, envconfig.LocalCREStateFilename),
		ChipIngressGRPCPort:    chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT,
	}
}

func getEnvironmentConfig(t *testing.T) *envconfig.Config {
	t.Helper()

	// we call our own Load function because it executes a couple of crucial extra input transformations
	in := &envconfig.Config{}
	err := in.Load(os.Getenv("CTF_CONFIGS"))
	require.NoError(t, err, "couldn't load environment state")
	return in
}

func createEnvironment(t *testing.T, testConfig *ttypes.TestConfig, flags ...string) {
	t.Helper()

	confErr := setConfigurationIfMissing(testConfig.EnvironmentConfigPath)
	require.NoError(t, confErr, "failed to set configuration")

	createErr := createEnvironmentIfNotExists(t.Context(), testConfig.RelativePathToRepoRoot, testConfig.EnvironmentDirPath, flags...)
	require.NoError(t, createErr, "failed to create environment")

	setErr := os.Setenv("CTF_CONFIGS", envconfig.MustLocalCREStateFileAbsPath(testConfig.RelativePathToRepoRoot))
	require.NoError(t, setErr, "failed to set CTF_CONFIGS env var")
}

func setConfigurationIfMissing(configName string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}

func createEnvironmentIfNotExists(ctx context.Context, relativePathToRepoRoot, environmentDir string, flags ...string) error {
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		framework.L.Info().Str("CTF_CONFIGS", os.Getenv("CTF_CONFIGS")).Str("local CRE state file", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)).Msg("Local CRE state file does not exist, starting environment...")

		args := []string{"run", ".", "env", "start"}
		args = append(args, flags...)

		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start environment")
		}
	}

	return nil
}

func cloneCldfEnvironmentForTest(src *cldf.Environment) *cldf.Environment {
	cloned := *src
	chainCopies := make([]cldf_chain.BlockChain, 0)
	for _, chain := range src.BlockChains.All() {
		chainCopies = append(chainCopies, chain)
	}
	cloned.BlockChains = cldf_chain.NewBlockChainsFromSlice(chainCopies)

	memStore := datastore.NewMemoryDataStore()
	if err := memStore.Merge(src.DataStore); err == nil {
		cloned.DataStore = memStore.Seal()
	}
	return &cloned
}

func setCldfEVMDeployerKey(env *cldf.Environment, chainSelector uint64, deployerKey *bind.TransactOpts) error {
	if env == nil {
		return errors.New("cldf environment is nil")
	}

	chainCopies := make([]cldf_chain.BlockChain, 0)
	found := false
	for selector, chain := range env.BlockChains.All() {
		if selector != chainSelector {
			chainCopies = append(chainCopies, chain)
			continue
		}

		evmChain, ok := chain.(cldf_evm.Chain)
		if !ok {
			return fmt.Errorf("chain selector %d is not EVM", chainSelector)
		}
		evmChain.DeployerKey = deployerKey
		chainCopies = append(chainCopies, evmChain)
		found = true
	}
	if !found {
		return fmt.Errorf("chain selector %d not found in CLDF environment", chainSelector)
	}

	env.BlockChains = cldf_chain.NewBlockChainsFromSlice(chainCopies)
	return nil
}
