package matrix

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildMatrix_ExpandsTopologiesForTest(t *testing.T) {
	t.Parallel()

	entries := BuildMatrix([]string{"Test_CRE_V2_Suite_Bucket_B"}, "12345", "2", "cpu=16/ram=64")

	require.Len(t, entries, 4)

	topologies := make([]string, 0, len(entries))
	for idx, entry := range entries {
		require.Equal(t, "Test_CRE_V2_Suite_Bucket_B", entry.TestName)
		require.Equal(t, strconv.Itoa(idx), entry.TestID)
		require.Equal(t, fmt.Sprintf("runs-on=12345-%d-2/cpu=16/ram=64", idx), entry.RunsOn)
		topologies = append(topologies, entry.Topology)
	}

	require.Equal(t, []string{
		"workflow-gateway-capabilities",
		"workflow-gateway-capabilities-vault-jwt_auth-enabled",
		"workflow-gateway-capabilities-vault-optimizations-enabled",
		"workflow-gateway-capabilities-vault-stall-purge",
	}, topologies)
}

func TestScanDir_ExcludesTestMain(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	sampleContent := `package sample_test
import "testing"
func TestMain(m *testing.M) {}
func Test_CRE_V2_RealTest(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	testNames, err := ScanDir(tempDir, `^(Test|Example).*`)
	require.NoError(t, err)
	require.Equal(t, []string{"Test_CRE_V2_RealTest"}, testNames)
}

func TestScanDir_Empty(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	sampleContent := `package sample_test
func NotATest() {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	_, err := ScanDir(tempDir, `^Test_CRE_.*`)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no matching test functions found")
}

func TestCRESmokeSuite_RealDirectory(t *testing.T) {
	t.Parallel()
	repoRoot, err := filepath.Abs("../../../..")
	require.NoError(t, err)

	smokeDir := filepath.Join(repoRoot, "system-tests/tests/smoke/cre")
	if _, statErr := os.Stat(smokeDir); statErr != nil {
		t.Skipf("skipping: %s not accessible", smokeDir)
	}

	testNames, err := ScanDir(smokeDir, `^(Test_CRE_|TestCRE_).*`)
	require.NoError(t, err)
	require.NotEmpty(t, testNames)

	// Verify unit tests are excluded
	for _, name := range testNames {
		require.NotEqual(t, "TestVaultStaticTopologies_LoadExpectedConfig", name)
		require.NotEqual(t, "TestMustMintVaultJWTForRequest_UsesRawRequestDigest", name)
		require.NotEqual(t, "Test_Upgrade_Suite", name)
		require.NotEqual(t, "TestMain", name)
	}
}

func TestBuildSuite_CRESmoke(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	sampleContent := `package sample_test
import "testing"
func Test_CRE_V2_Aptos_Suite(t *testing.T) {}
func Test_CRE_V2_Suite_Bucket_A(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	entries, err := BuildSuiteMatrix(SuiteCRESmoke, SuiteOptions{
		Dir:        tempDir,
		RunID:      "123",
		RunAttempt: "1",
	})
	require.NoError(t, err)
	require.Len(t, entries, 2)
	require.Equal(t, "Test_CRE_V2_Aptos_Suite", entries[0].TestName)
	require.Equal(t, "workflow-gateway-aptos", entries[0].Topology)
	require.Equal(t, "configs/workflow-gateway-don-aptos.toml", entries[0].Configs)
	require.Contains(t, entries[0].RunsOn, "family=m7i+m8i")
}

func TestBuildSuite_CRERegression(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	sampleContent := `package sample_test
import "testing"
func Test_CRE_V2_Stellar_Regression(t *testing.T) {}
func Test_CRE_V2_Standard_Regression(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	entries, err := BuildSuiteMatrix(SuiteCRERegression, SuiteOptions{
		Dir:        tempDir,
		RunID:      "123",
		RunAttempt: "1",
	})
	require.NoError(t, err)
	require.Len(t, entries, 2)
	require.Equal(t, "Test_CRE_V2_Standard_Regression", entries[0].TestName)
	require.Equal(t, "configs/workflow-gateway-capabilities-don.toml", entries[0].Configs)
	require.Equal(t, "Test_CRE_V2_Stellar_Regression", entries[1].TestName)
	require.Equal(t, "configs/workflow-gateway-don-stellar.toml", entries[1].Configs)
}

func TestBuildSuite_CREMixedEnv(t *testing.T) {
	t.Parallel()
	entries, err := BuildSuiteMatrix(SuiteCREMixedEnv, SuiteOptions{
		RunID:      "123",
		RunAttempt: "1",
	})
	require.NoError(t, err)
	require.Len(t, entries, 5)
	require.Equal(t, "Test_CRE_V2_Suite_Bucket_A", entries[0].TestName)
	require.Contains(t, entries[0].RunsOn, "family=m7i+m8i")
}

func TestBuildSuite_CCIP(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	sampleContent := `package sample_test
import "testing"
func Test_CCIPGasPriceUpdatesWriteFrequency(t *testing.T) {}
func TestRMN_GlobalCurseTwoMessagesOnTwoLanes(t *testing.T) {}
func TestDeleteCCIPJobs(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	entries, err := BuildSuiteMatrix(SuiteCCIP, SuiteOptions{
		Dir:        tempDir,
		RunID:      "123",
		RunAttempt: "1",
	})
	require.NoError(t, err)
	require.Len(t, entries, 3)

	require.Equal(t, "TestDeleteCCIPJobs", entries[0].TestName)
	require.Equal(t, 20, entries[0].JobTimeout)

	require.Equal(t, "TestRMN_GlobalCurseTwoMessagesOnTwoLanes", entries[1].TestName)
	require.Equal(t, "master-amd6416f5d86", entries[1].RMNRageProxyVersion)
	require.Equal(t, "master-amd64-10b42b2", entries[1].RMNAFN2ProxyVersion)

	require.Equal(t, "Test_CCIPGasPriceUpdatesWriteFrequency", entries[2].TestName)
	require.Equal(t, "15m", entries[2].Timeout)
	require.Equal(t, "SIMULATED_1,SIMULATED_2", entries[2].SelectedNetwork)
	require.Contains(t, entries[2].RunsOn, "family=r6i+r7i+r8i")
}

func TestGenerateSetupMatrices(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	sampleContent := `package sample_test
import "testing"
func Test_CRE_V2_Suite_Bucket_A(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample_test.go"), []byte(sampleContent), 0o600))

	matrices, err := GenerateSetupMatrices(SetupOptions{
		RunID:         "999",
		RunAttempt:    "2",
		CRESmoke:      true,
		CRESmokeDir:   tempDir,
		CRERegression: false,
		CREMixedEnv:   true,
		CCIP:          true,
	})
	require.NoError(t, err)
	require.Contains(t, matrices, "cre-matrix")
	require.Contains(t, matrices, "cre-mixed-env-matrix")
	require.Contains(t, matrices, "ccip-matrix")
	require.NotContains(t, matrices, "cre-regression-matrix")
	require.Len(t, matrices["cre-matrix"], 1)
	require.Len(t, matrices["cre-mixed-env-matrix"], 5)
	require.Len(t, matrices["ccip-matrix"], 3)
}

func TestWriteMultiOutput(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", tempFile)

	matrices := map[string][]Entry{
		"ccip-matrix": {
			{TestName: "Test1", TestID: "0", RunsOn: "runner1"},
		},
		"cre-matrix": {
			{TestName: "Test2", TestID: "1", RunsOn: "runner2"},
		},
	}

	var stdout strings.Builder
	err := WriteMultiOutput(&stdout, matrices, true)
	require.NoError(t, err)
	require.Contains(t, stdout.String(), `"ccip-matrix"`)
	require.Contains(t, stdout.String(), `"cre-matrix"`)

	fileBytes, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	outputStr := string(fileBytes)
	require.Contains(t, outputStr, "ccip-matrix=")
	require.Contains(t, outputStr, "cre-matrix=")
}
