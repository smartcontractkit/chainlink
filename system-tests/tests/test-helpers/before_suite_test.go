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

// newFakeCreEnv installs a fake precompiled cre-env binary that records the CTF_CONFIGS
// value it was started with, and returns a test config pointing at it plus the capture file path.
func newFakeCreEnv(t *testing.T) (*ttypes.TestConfig, string) {
	t.Helper()

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

	testConfig := &ttypes.TestConfig{
		RelativePathToRepoRoot: tmpDir,
		EnvironmentDirPath:     environmentDir,
		EnvironmentConfigPath:  filepath.Join(tmpDir, "custom-config.toml"),
		EnvironmentStateFile:   filepath.Join(environmentDir, envconfig.LocalCREStateFilename),
	}

	return testConfig, capturedEnvFile
}

func readCapturedCTFConfigs(t *testing.T, capturedEnvFile string) string {
	t.Helper()

	captured, readErr := os.ReadFile(capturedEnvFile)
	require.NoError(t, readErr)
	return strings.TrimSpace(string(captured))
}

func TestCreateEnvironmentIfNotExists_SubprocessEnv(t *testing.T) {
	t.Setenv("CTF_CONFIGS", "")

	testConfig, capturedEnvFile := newFakeCreEnv(t)

	err := createEnvironmentIfNotExists(context.Background(), testConfig)
	require.NoError(t, err)

	// Verify child process received the test config path via cmd.Env
	require.Equal(t, testConfig.EnvironmentConfigPath, readCapturedCTFConfigs(t, capturedEnvFile))

	// Verify parent process environment was not mutated
	require.Empty(t, os.Getenv("CTF_CONFIGS"))
}

func TestCreateEnvironmentIfNotExists_UserCTFConfigsTakesPrecedence(t *testing.T) {
	userConfigPath := filepath.Join(t.TempDir(), "user-config.toml")
	t.Setenv("CTF_CONFIGS", userConfigPath)

	testConfig, capturedEnvFile := newFakeCreEnv(t)

	err := createEnvironmentIfNotExists(context.Background(), testConfig)
	require.NoError(t, err)

	// A CTF_CONFIGS value set by the caller overrides the test's default topology
	require.Equal(t, userConfigPath, readCapturedCTFConfigs(t, capturedEnvFile))
	require.Equal(t, userConfigPath, os.Getenv("CTF_CONFIGS"))
}
