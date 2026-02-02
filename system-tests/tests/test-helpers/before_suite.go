package helpers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

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
