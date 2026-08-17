package matrix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanDir_ASTDiscovery(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

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

	testNames, err := ScanDir(tempDir, `^TestCRE_.*_E2E$`)
	require.NoError(t, err)

	expected := []string{
		"TestCRE_V2_Aptos_Suite_E2E",
		"TestCRE_V2_Solana_Write_E2E",
		"TestCRE_V2_Suite_Bucket_A_E2E",
	}
	require.Equal(t, expected, testNames)
}

func TestBuildMatrix(t *testing.T) {
	t.Parallel()
	testNames := []string{
		"TestCRE_V2_Aptos_Suite_E2E",
		"TestCRE_V2_Solana_Write_E2E",
		"TestCRE_V2_Suite_Bucket_A_E2E",
	}

	entries := BuildMatrix(testNames, "12345", "2", "cpu=16/ram=64")
	require.Len(t, entries, 3)

	require.Equal(t, "TestCRE_V2_Aptos_Suite_E2E", entries[0].TestName)
	require.Equal(t, "TestCRE_V2_Aptos_Suite_E2E", entries[0].TestID)
	require.Equal(t, "runs-on=12345-0-2/cpu=16/ram=64", entries[0].RunsOn)
	require.Equal(t, "workflow-gateway-aptos", entries[0].Topology)
	require.Equal(t, "configs/workflow-gateway-don-aptos.toml", entries[0].Configs)

	require.Equal(t, "TestCRE_V2_Solana_Write_E2E", entries[1].TestName)
	require.Equal(t, "TestCRE_V2_Solana_Write_E2E", entries[1].TestID)
	require.Equal(t, "runs-on=12345-1-2/cpu=16/ram=64", entries[1].RunsOn)
	require.Equal(t, "workflow", entries[1].Topology)
	require.Equal(t, "configs/workflow-don-solana.toml", entries[1].Configs)

	require.Equal(t, "TestCRE_V2_Suite_Bucket_A_E2E", entries[2].TestName)
	require.Equal(t, "TestCRE_V2_Suite_Bucket_A_E2E", entries[2].TestID)
	require.Equal(t, "runs-on=12345-2-2/cpu=16/ram=64", entries[2].RunsOn)
	require.Equal(t, "workflow-gateway-capabilities", entries[2].Topology)
	require.Equal(t, "configs/workflow-gateway-capabilities-don.toml", entries[2].Configs)
}

func TestScanDir_Empty(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	sampleContent := `package sample_test
func NotATest() {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0600))

	_, err := ScanDir(tempDir, `^TestCRE_.*_E2E$`)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no matching test functions found")
}

func TestCRESmokeSuite_RealDirectory(t *testing.T) {
	t.Parallel()
	// Find repo root from current test file
	repoRoot, err := filepath.Abs("../../../..")
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

	// Verify unit tests are excluded
	for _, name := range testNames {
		require.NotEqual(t, "TestVaultStaticTopologies_LoadExpectedConfig", name)
		require.NotEqual(t, "TestMustMintVaultJWTForRequest_UsesRawRequestDigest", name)
		require.NotEqual(t, "Test_Upgrade_Suite", name)
	}
}
