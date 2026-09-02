package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/cmd"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func TestMatrixSystem_CLI(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "smoke_test.go")
	content := `package smoke_test

import "testing"

func Test_CRE_V2_Demo(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"matrix", "system",
		"--suite", "cre-smoke",
		"--dir", tmpDir,
		"--run-id", "123",
		"--run-attempt", "1",
		"--spot-flag", "spot=co",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res []matrix.CRESmokeEntry
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "Test_CRE_V2_Demo", res[0].TestName)
	assert.Equal(t, "workflow-gateway-capabilities", res[0].Topology)
	assert.Equal(t, "runs-on=123-0-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)
}

func TestMatrixInMemory_CLI(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "in-memory-tests.json")
	content := `[
  {"name":"sample_test.go","test":"TestSample","timeout":"10m","parallel":1,"plugins":false,"runs_on":"","free_disk":false,"aptos":"","sui":""}
]`
	require.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"matrix", "in-memory",
		"--file", configFile,
		"--run-id", "456",
		"--run-attempt", "1",
		"--spot-flag", "spot=pco",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res []matrix.InMemoryEntry
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "sample_test.go", res[0].Name)
	assert.Equal(t, 15, res[0].JobTimeout)
	assert.Equal(t, "ubuntu-latest", res[0].RunsOn)
}

func TestMatrixCCIP_CLI(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"matrix", "ccip",
		"--run-id", "789",
		"--run-attempt", "2",
		"--spot-flag", "spot=co",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res []matrix.CCIPSystemEntry
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	require.Len(t, res, 3)
	assert.Equal(t, "Test_CCIPGasPriceUpdatesWriteFrequency", res[0].TestName)
}

func TestMatrixMixedEnv_CLI(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"matrix", "mixed-env",
		"--run-id", "321",
		"--run-attempt", "1",
		"--spot-flag", "spot=pco",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res []matrix.CREMixedEnvEntry
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	require.Len(t, res, 6)
	assert.Equal(t, "Test_CRE_V2_Suite_Bucket_A", res[0].TestName)
}
