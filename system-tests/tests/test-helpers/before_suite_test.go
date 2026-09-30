package helpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func TestCreateEnvironment_ParallelSafe(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, envconfig.StateDirname)
	require.NoError(t, os.MkdirAll(stateDir, 0o755))
	stateFile := filepath.Join(stateDir, envconfig.LocalCREStateFilename)
	require.NoError(t, os.WriteFile(stateFile, []byte(""), 0o600))

	testConfig := &ttypes.TestConfig{
		RelativePathToRepoRoot: tmpDir,
		EnvironmentDirPath:     tmpDir,
		EnvironmentConfigPath:  filepath.Join(tmpDir, "test.toml"),
		EnvironmentStateFile:   stateFile,
	}

	// createEnvironment must not panic when called after t.Parallel(),
	// and must not mutate process-level CTF_CONFIGS.
	originalCTFConfigs := os.Getenv("CTF_CONFIGS")
	createEnvironment(t, testConfig)
	require.Equal(t, originalCTFConfigs, os.Getenv("CTF_CONFIGS"))
}

func TestCreateEnvironmentIfNotExists_SubprocessEnv(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	environmentDir := filepath.Join(tmpDir, "core", "scripts", "cre", "environment")
	require.NoError(t, os.MkdirAll(environmentDir, 0o755))

	binDir := filepath.Join(tmpDir, "system-tests", "tests", "bin")
	require.NoError(t, os.MkdirAll(binDir, 0o755))

	capturedEnvFile := filepath.Join(tmpDir, "captured_env.txt")
	scriptContent := fmt.Sprintf("#!/bin/sh\nprintf \"%%s\" \"$CTF_CONFIGS\" > %q\n", capturedEnvFile)
	binPath := filepath.Join(binDir, "cre-env")
	require.NoError(t, os.WriteFile(binPath, []byte(scriptContent), 0o600))
	require.NoError(t, os.Chmod(binPath, 0o700))

	expectedConfigPath := filepath.Join(tmpDir, "custom-config.toml")
	testConfig := &ttypes.TestConfig{
		RelativePathToRepoRoot: tmpDir,
		EnvironmentDirPath:     environmentDir,
		EnvironmentConfigPath:  expectedConfigPath,
		EnvironmentStateFile:   filepath.Join(environmentDir, envconfig.LocalCREStateFilename),
	}

	originalEnv := os.Getenv("CTF_CONFIGS")
	err := createEnvironmentIfNotExists(context.Background(), testConfig)
	require.NoError(t, err)

	// Verify child process received CTF_CONFIGS from cmd.Env
	captured, readErr := os.ReadFile(capturedEnvFile)
	require.NoError(t, readErr)
	require.Equal(t, expectedConfigPath, strings.TrimSpace(string(captured)))

	// Verify parent process environment was not mutated
	require.Equal(t, originalEnv, os.Getenv("CTF_CONFIGS"))
}
