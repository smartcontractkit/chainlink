package cre

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
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
)

// TestConfig holds common test specific configurations related to the test execution
// These configurations are not meant to impact the actual test logic
type TestConfig struct {
	EnvironmentConfigPath   string
	EnvironmentDirPath      string
	EnvironmentArtifactPath string
	BeholderConfigPath      string
}

// TestEnvironment holds references to the main test components
type TestEnvironment struct {
	Config                   *envconfig.Config
	TestConfig               *TestConfig
	EnvArtifact              environment.EnvArtifact
	Logger                   zerolog.Logger
	FullCldEnvOutput         *cre.FullCLDEnvironmentOutput
	WrappedBlockchainOutputs []*cre.WrappedBlockchainOutput
}

// setupTestEnvironment initializes the common test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	defaultTestConfig := getDefaultTestConfig(t)
	createEnvironment(t, defaultTestConfig)
	in := getEnvironmentConfig(t)
	envArtifact := getEnvironmentArtifact(t)
	fullCldEnvOutput, wrappedBlockchainOutputs, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, err, "failed to load environment")

	return &TestEnvironment{
		Config:                   in,
		TestConfig:               defaultTestConfig,
		EnvArtifact:              envArtifact,
		Logger:                   framework.L,
		FullCldEnvOutput:         fullCldEnvOutput,
		WrappedBlockchainOutputs: wrappedBlockchainOutputs,
	}
}

func getDefaultTestConfig(t *testing.T) *TestConfig {
	t.Helper()

	return &TestConfig{
		EnvironmentDirPath:      "../../../../core/scripts/cre/environment",
		EnvironmentConfigPath:   "../../../../core/scripts/cre/environment/configs/workflow-don.toml",
		EnvironmentArtifactPath: "../../../../core/scripts/cre/environment/env_artifact/env_artifact.json",
		BeholderConfigPath:      "../../../../core/scripts/cre/environment/configs/chip-ingress-cache.toml",
	}
}

func getEnvironmentConfig(t *testing.T) *envconfig.Config {
	t.Helper()

	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")
	return in
}

func getEnvironmentArtifact(t *testing.T) environment.EnvArtifact {
	t.Helper()

	var envArtifact environment.EnvArtifact
	artFile, err := os.ReadFile(os.Getenv("ENV_ARTIFACT_PATH"))
	require.NoError(t, err, "failed to read artifact file")

	err = json.Unmarshal(artFile, &envArtifact)
	require.NoError(t, err, "failed to unmarshal artifact file")
	return envArtifact
}

func createEnvironment(t *testing.T, testConfig *TestConfig) {
	t.Helper()

	confErr := setConfigurationIfMissing(testConfig.EnvironmentConfigPath, testConfig.EnvironmentArtifactPath)
	require.NoError(t, confErr, "failed to set configuration")

	createErr := createEnvironmentIfNotExists(testConfig.EnvironmentDirPath)
	require.NoError(t, createErr, "failed to create environment")

	// transform the config file to the cache file, so that we can use the cached environment
	cachedConfigFile, cacheErr := ctfConfigToCacheFile()
	require.NoError(t, cacheErr, "failed to get cached config file")

	setErr := os.Setenv("CTF_CONFIGS", cachedConfigFile)
	require.NoError(t, setErr, "failed to set CTF_CONFIGS env var")
}

func setConfigurationIfMissing(configName, envArtifactPath string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	if os.Getenv("ENV_ARTIFACT_PATH") == "" {
		err := os.Setenv("ENV_ARTIFACT_PATH", envArtifactPath)
		if err != nil {
			return errors.Wrap(err, "failed to set ENV_ARTIFACT_PATH env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}

func createEnvironmentIfNotExists(environmentDir string) error {
	cachedConfigFile, cacheErr := ctfConfigToCacheFile()
	if cacheErr != nil {
		return errors.Wrap(cacheErr, "failed to get cached config file")
	}

	if _, err := os.Stat(cachedConfigFile); os.IsNotExist(err) {
		framework.L.Info().Str("cached_config_file", cachedConfigFile).Msg("Cached config file does not exist, starting environment...")
		cmd := exec.Command("go", "run", ".", "env", "start")
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

func ctfConfigToCacheFile() (string, error) {
	configFile := os.Getenv("CTF_CONFIGS")
	if configFile == "" {
		return "", errors.New("CTF_CONFIGS env var is not set")
	}

	if strings.HasSuffix(configFile, "-cache.toml") {
		return configFile, nil
	}

	split := strings.Split(configFile, ",")
	return strings.ReplaceAll(split[0], ".toml", "") + "-cache.toml", nil
}
