package aptos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstWorkingAptosCLI_SkipsInvalidBinary(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	badCLI := filepath.Join(tempDir, "bad-aptos")
	goodCLI := filepath.Join(tempDir, "good-aptos")

	require.NoError(t, os.WriteFile(badCLI, []byte("#!/bin/sh\necho 'not an aptos cli' >&2\nexit 1\n"), 0o755))
	require.NoError(t, os.WriteFile(goodCLI, []byte("#!/bin/sh\necho 'aptos 7.2.0'\n"), 0o755))

	got, err := firstWorkingAptosCLI([]string{badCLI, goodCLI})
	require.NoError(t, err)
	require.Equal(t, goodCLI, got)
}

func TestEnsureHostAptosCLI_PrependsResolvedBinary(t *testing.T) {
	tempDir := t.TempDir()
	goodCLI := filepath.Join(tempDir, "aptos")
	require.NoError(t, os.WriteFile(goodCLI, []byte("#!/bin/sh\necho 'aptos 7.2.0'\n"), 0o755))

	t.Setenv(aptosCLIPathEnvVar, goodCLI)
	t.Setenv("PATH", "/usr/bin")

	require.NoError(t, ensureHostAptosCLI())
	require.Equal(t, tempDir, filepath.SplitList(os.Getenv("PATH"))[0])
}
