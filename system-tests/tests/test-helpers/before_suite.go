package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// TestConfig holds common test specific configurations related to the test execution
// These configurations are not meant to impact the actual test logic
type TestConfig struct {
	RelativePathToRepoRoot   string
	EnvironmentConfigPath    string
	EnvironmentDirPath       string
	EnvironmentStateFile     string
	EnvironmentArtifactPaths string
	BeholderStateFile        string
}

// TestEnvironment holds references to the main test components
type TestEnvironment struct {
	Config         *envconfig.Config
	TestConfig     *TestConfig
	Logger         zerolog.Logger
	CreEnvironment *cre.Environment
	Blockchains    []blockchains.Blockchain
}

func SetupTestEnvironmentWithConfig(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()

	createEnvironment(t, tconf, flags...)
	in := getEnvironmentConfig(t)
	creEnvironment, dons, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in)
	require.NoError(t, err, "failed to load environment")

	t.Cleanup(func() {
		if t.Failed() {
			framework.L.Warn().Msg("Test failed - checking for panics in Docker containers...")
			foundPanics := infra.CheckContainersForPanics(framework.L, 100)
			if !foundPanics {
				var lastLines uint64 = 30
				framework.L.Warn().Msgf("No panic patterns detected in Docker container logs. Displaying last %d lines of logs for debugging:", lastLines)
				infra.PrintFailedContainerLogs(framework.L, lastLines)
			}
		}
	})

	return &ttypes.TestEnvironment{
		Config:         in,
		TestConfig:     tconf,
		Logger:         framework.L,
		CreEnvironment: creEnvironment,
		Dons:           dons,
	}
}

func GetDefaultTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return GetTestConfig(t, "/configs/workflow-gateway-don.toml")
}

func GetTestConfig(t *testing.T, configPath string) *ttypes.TestConfig {
	relativePathToRepoRoot := "../../../../"
	environmentDirPath := filepath.Join(relativePathToRepoRoot, "core/scripts/cre/environment")

	return &ttypes.TestConfig{
		RelativePathToRepoRoot: relativePathToRepoRoot,
		EnvironmentDirPath:     environmentDirPath,
		EnvironmentConfigPath:  filepath.Join(environmentDirPath, configPath), // change to your desired config, if you want to use another topology
		EnvironmentStateFile:   filepath.Join(environmentDirPath, envconfig.StateDirname, envconfig.LocalCREStateFilename),
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

	createErr := createEnvironmentIfNotExists(testConfig.RelativePathToRepoRoot, testConfig.EnvironmentDirPath, flags...)
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

func createEnvironmentIfNotExists(relativePathToRepoRoot, environmentDir string, flags ...string) error {
	stateExists := envconfig.LocalCREStateFileExists(relativePathToRepoRoot)
	freshEnv := os.Getenv("FRESH_ENV") == "true"

	// If FRESH_ENV=true and state exists, stop environment and delete state to force recreation
	if freshEnv && stateExists {
		framework.L.Info().Msg("FRESH_ENV=true - stopping existing environment and deleting state...")

		// Stop existing environment first
		stopCmd := exec.Command("go", "run", ".", "env", "stop")
		stopCmd.Dir = environmentDir
		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr
		_ = stopCmd.Run() // Ignore errors - environment might not be running

		// Remove all test containers (including Anvil/blockchain containers) to ensure fresh state
		// This is critical to avoid reusing old blockchain state with incorrect contract data
		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			framework.L.Warn().Msgf("Failed to remove test containers (this may be expected if containers don't exist): %s", removeErr)
		} else {
			framework.L.Info().Msg("Removed all test containers (including blockchain/Anvil) to ensure fresh state")
		}

		// Delete state directory
		stateDir := filepath.Join(relativePathToRepoRoot, envconfig.StateDirname)
		if err := os.RemoveAll(stateDir); err != nil {
			return errors.Wrap(err, "failed to delete state directory")
		}
		stateExists = false
	}

	if !stateExists {
		framework.L.Info().Str("CTF_CONFIGS", os.Getenv("CTF_CONFIGS")).Str("local CRE state file", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)).Msg("Local CRE state file does not exist, starting environment...")

		args := []string{"run", ".", "env", "start"}
		args = append(args, flags...)

		cmd := exec.Command("go", args...)
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
