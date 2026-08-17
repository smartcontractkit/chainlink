package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanDir_ASTDiscovery(t *testing.T) {
	tempDir := t.TempDir()

	// Sample test file with mixed E2E tests, unit tests, helper funcs, and examples
	sampleContent := `package sample_test

import "testing"

func TestCRE_V2_Suite_Bucket_A_E2E(t *testing.T) {}
func TestCRE_V2_Aptos_Suite_E2E(t *testing.T) {}
func TestVaultStaticTopologies_LoadExpectedConfig(t *testing.T) {}
func TestMustMintVaultJWTForRequest_UsesRawRequestDigest(t *testing.T) {}
func HelperFunc() {}
func ExampleCRE_E2E() {}
func TestCRE_V2_Solana_Write_E2E(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	// Another file in same package
	anotherContent := `package sample_test

import "testing"

func TestCRE_V2_Sharding_E2E(t *testing.T) {}
func NotATest() {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "another_test.go"), []byte(anotherContent), 0600))

	testNames, err := ScanDir(tempDir, `^TestCRE_.*_E2E$`)
	require.NoError(t, err)

	expected := []string{
		"TestCRE_V2_Aptos_Suite_E2E",
		"TestCRE_V2_Sharding_E2E",
		"TestCRE_V2_Solana_Write_E2E",
		"TestCRE_V2_Suite_Bucket_A_E2E",
	}
	require.Equal(t, expected, testNames)
}

func TestBuildMatrix(t *testing.T) {
	testNames := []string{
		"TestCRE_Aptos_E2E",
		"TestCRE_Solana_E2E",
	}

	entries := BuildMatrix(testNames, "12345", "2", "cpu=16/ram=64/family=m7i+m8i/spot=co")
	require.Len(t, entries, 2)

	require.Equal(t, "TestCRE_Aptos_E2E", entries[0].TestName)
	require.Equal(t, "TestCRE_Aptos_E2E", entries[0].TestID)
	require.Equal(t, "runs-on=12345-0-2/cpu=16/ram=64/family=m7i+m8i/spot=co", entries[0].RunsOn)

	require.Equal(t, "TestCRE_Solana_E2E", entries[1].TestName)
	require.Equal(t, "TestCRE_Solana_E2E", entries[1].TestID)
	require.Equal(t, "runs-on=12345-1-2/cpu=16/ram=64/family=m7i+m8i/spot=co", entries[1].RunsOn)
}

func TestCLI_Run_StdoutJSON(t *testing.T) {
	tempDir := t.TempDir()
	sampleContent := `package sample_test
import "testing"
func TestCRE_Feature_E2E(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := []string{
		"--dir=" + tempDir,
		"--pattern=^TestCRE_.*_E2E$",
		"--run-id=999",
		"--attempt=1",
		"--runner=cpu=16/ram=64",
	}

	err := run(args, &stdout, &stderr)
	require.NoError(t, err)

	var entries []MatrixEntry
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &entries))
	require.Len(t, entries, 1)
	require.Equal(t, "TestCRE_Feature_E2E", entries[0].TestName)
	require.Equal(t, "runs-on=999-0-1/cpu=16/ram=64", entries[0].RunsOn)
}

func TestCLI_Run_GithubOutput(t *testing.T) {
	tempDir := t.TempDir()
	sampleContent := `package sample_test
import "testing"
func TestCRE_Feature_E2E(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	outputFile := filepath.Join(tempDir, "github_output.txt")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := []string{
		"--dir=" + tempDir,
		"--pattern=^TestCRE_.*_E2E$",
		"--run-id=888",
		"--attempt=3",
		"--runner=cpu=16/ram=64",
		"--github-output",
	}

	err := run(args, &stdout, &stderr)
	require.NoError(t, err)

	outputContent, readErr := os.ReadFile(outputFile)
	require.NoError(t, readErr)
	require.Contains(t, string(outputContent), "matrix=[")
	require.Contains(t, string(outputContent), "TestCRE_Feature_E2E")
}

func TestCLI_InvalidArgs(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Invalid regex
	err := run([]string{"--pattern=[invalid"}, &stdout, &stderr)
	require.Error(t, err)

	// Nonexistent directory
	err = run([]string{"--dir=/nonexistent/path/here"}, &stdout, &stderr)
	require.Error(t, err)
}

func TestCRESmokeSuite_RealDirectory(t *testing.T) {
	// Find repo root from current test file
	repoRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	smokeDir := filepath.Join(repoRoot, "system-tests/tests/smoke/cre")
	if _, statErr := os.Stat(smokeDir); statErr != nil {
		t.Skipf("skipping: %s not accessible", smokeDir)
	}

	testNames, err := ScanDir(smokeDir, `^TestCRE_.*_E2E$`)
	require.NoError(t, err)
	require.NotEmpty(t, testNames)

	// Verify all discovered test names have _E2E suffix
	for _, name := range testNames {
		require.Contains(t, name, "TestCRE_")
		require.True(t, filepath.Ext(name) == "" && name[len(name)-4:] == "_E2E", "test name %s must end with _E2E", name)
	}

	// Verify unit tests are not discovered
	for _, name := range testNames {
		require.NotEqual(t, "TestVaultStaticTopologies_LoadExpectedConfig", name)
		require.NotEqual(t, "TestMustMintVaultJWTForRequest_UsesRawRequestDigest", name)
		require.NotEqual(t, "Test_Upgrade_Suite", name)
	}
}
