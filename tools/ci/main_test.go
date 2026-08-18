package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/cmd"
)

func TestRootHelp(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "Universal CI CLI for test matrix discovery")
	require.Contains(t, stdout.String(), "matrix")
	require.Contains(t, stdout.String(), "testshard")
	require.Contains(t, stdout.String(), "changelog")
}

func TestMatrixSubcommand(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	sampleContent := `package sample_test
import "testing"
func TestCRE_Example_E2E(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"matrix", "--dir=" + tempDir, "--run-id=99", "--attempt=1"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.Contains(t, stdout.String(), `"test_name":"TestCRE_Example_E2E"`)
	require.Contains(t, stdout.String(), `"runs_on":"runs-on=99-0-1/`)
}

func TestTestshardSubcommand(t *testing.T) {
	t.Parallel()
	stdin := strings.NewReader("pkgA\npkgB\npkgC\n")
	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(stdin, &stdout, &stderr)
	rootCmd.SetArgs([]string{"testshard", "list", "--shard-count=2", "--shard-index=0"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, strings.TrimSpace(stdout.String()))
}

func TestChangelogSubcommand(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "package.json")
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")

	require.NoError(t, os.WriteFile(pkgPath, []byte(`{"version": "1.0.0"}`), 0600))
	require.NoError(t, os.WriteFile(changelogPath, []byte("# Changelog\n\n## 1.0.0\n\n- [#added] Feature X\n\n## 0.9.0\n"), 0600))

	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"changelog", "format", "--changelog=" + changelogPath, "--package-json=" + pkgPath, "--github-output=false"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "Formatted changelog for version 1.0.0")
}

func TestMatrixSuiteSubcommand(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	sampleContent := `package sample_test
import "testing"
func TestCCIP_GasPriceUpdatesWriteFrequency_E2E(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"matrix", "--suite=ccip", "--dir=" + tempDir, "--run-id=10", "--attempt=1"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.Contains(t, stdout.String(), `"test_name":"TestCCIP_GasPriceUpdatesWriteFrequency_E2E"`)
}

func TestMatrixRegressionSuiteUsesSuiteDefaults(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	sampleContent := `package sample_test
import "testing"
func Test_CRE_V2_Foo_Regression(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"matrix", "--suite=cre-regression", "--dir=" + tempDir, "--run-id=10", "--attempt=1"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	require.Contains(t, stdout.String(), `"test_name":"Test_CRE_V2_Foo_Regression"`)
}

func TestMatrixSetupSubcommand(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"matrix", "setup", "--ccip=true", "--cre-mixed-env=true", "--run-id=10", "--attempt=1"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	require.Contains(t, stdout.String(), `"ccip-matrix":`)
	require.Contains(t, stdout.String(), `"cre-mixed-env-matrix":`)
}

func TestGatingSubcommand(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", outputFile)
	t.Setenv("GITHUB_STEP_SUMMARY", "")
	t.Setenv("EVENT_NAME", "pull_request")
	t.Setenv("REF_NAME", "feature/x")
	t.Setenv("REF_TYPE", "branch")
	t.Setenv("CRE_CHANGES", "true")
	t.Setenv("CCIP_CHANGES", "false")

	var stdout, stderr bytes.Buffer
	rootCmd := cmd.NewRootCmd(nil, &stdout, &stderr)
	rootCmd.SetArgs([]string{"gating"})

	require.NoError(t, rootCmd.Execute())

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	output := string(content)
	require.Contains(t, output, "cre-should-run=true")
	require.Contains(t, output, "cre-with-regression=true")
	require.Contains(t, output, "cre-run-mixed-env=true")
	require.Contains(t, output, "ccip-should-run=false")
	require.Contains(t, output, "build-core-image=true")
	require.Contains(t, output, "build-plugins-image=true")
}
