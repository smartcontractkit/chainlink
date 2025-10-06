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
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
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
	Config                   *envconfig.Config
	TestConfig               *TestConfig
	EnvArtifact              *environment.EnvArtifact
	Logger                   zerolog.Logger
	CreEnvironment           *cre.Environment
	WrappedBlockchainOutputs []*cre.WrappedBlockchainOutput
}

func SetupTestEnvironmentWithConfig(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()

	createEnvironment(t, tconf, flags...)
	in := getEnvironmentConfig(t)
	envArtifact := getEnvironmentArtifact(t, tconf.RelativePathToRepoRoot)
	creEnvironment, wrappedBlockchainOutputs, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, err, "failed to load environment")

	return &ttypes.TestEnvironment{
		Config:                   in,
		TestConfig:               tconf,
		EnvArtifact:              envArtifact,
		Logger:                   framework.L,
		CreEnvironment:           creEnvironment,
		WrappedBlockchainOutputs: wrappedBlockchainOutputs,
	}
}

func GetDefaultTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return GetTestConfig(t, "/configs/workflow-don.toml")
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

	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")
	return in
}

func getEnvironmentArtifact(t *testing.T, relativePathToRepoRoot string) *environment.EnvArtifact {
	t.Helper()

	envArtifact, artErr := environment.ReadEnvArtifact(environment.MustEnvArtifactAbsPath(relativePathToRepoRoot))
	require.NoError(t, artErr, "failed to read environment artifact")
	return envArtifact
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
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
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
