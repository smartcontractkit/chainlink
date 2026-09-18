package matrix_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func TestDiscoverGoTestNames(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "sample_test.go")
	content := `package sample_test

import "testing"

func TestMain(m *testing.M) {}
func Test_Smoke_HappyPath(t *testing.T) {}
func Test_Upgrade_Something(t *testing.T) {}
func Example_Something() {}
func Test_Smoke_ComplexFlow(t *testing.T) {}
func helperFunction() {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	names, err := matrix.DiscoverGoTestNames(tmpDir, matrix.DiscoverOptions{
		IgnoredPatterns: []string{"TestMain", "Test_Upgrade.*"},
	})
	require.NoError(t, err)

	assert.Equal(t, []string{"Example_Something", "Test_Smoke_ComplexFlow", "Test_Smoke_HappyPath"}, names)
}

func TestBuildCRESmokeMatrix(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "smoke_test.go")
	content := `package smoke_test

import "testing"

func Test_CRE_V2_Basic(t *testing.T) {}
func Test_CRE_V2_Sharding(t *testing.T) {}
func Test_CRE_V2_Suite_Bucket_B(t *testing.T) {}
func Test_CRE_V2_Solana_Write(t *testing.T) {}
func Test_CRE_V2_Solana_Read_Tx(t *testing.T) {}
func Test_CRE_V2_Stellar_Suite(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	res, err := matrix.BuildCRESmokeMatrix(context.Background(), matrix.CRESmokeOptions{
		Dir:        tmpDir,
		RunID:      "123456",
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.NoError(t, err)
	require.Len(t, res, 7)

	// Test_CRE_V2_Basic has default topology
	basic := res[0]
	assert.Equal(t, "Test_CRE_V2_Basic", basic.TestName)
	assert.Equal(t, 0, basic.TestID)
	assert.Equal(t, "workflow-gateway-capabilities", basic.Topology)
	assert.Equal(t, "configs/workflow-gateway-capabilities-don.toml", basic.Configs)
	assert.Equal(t, "runs-on=123456-0-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", basic.RunsOn)

	// Test_CRE_V2_Sharding has per-test override
	sharding := res[1]
	assert.Equal(t, "Test_CRE_V2_Sharding", sharding.TestName)
	assert.Equal(t, 1, sharding.TestID)
	assert.Equal(t, "workflow-gateway-sharded", sharding.Topology)
	assert.Equal(t, "configs/workflow-gateway-sharded-don.toml", sharding.Configs)
	assert.Equal(t, "runs-on=123456-1-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", sharding.RunsOn)

	// Solana tests use the workflow topology with solana configs
	solanaReadTx := res[2]
	assert.Equal(t, "Test_CRE_V2_Solana_Read_Tx", solanaReadTx.TestName)
	assert.Equal(t, 2, solanaReadTx.TestID)
	assert.Equal(t, "workflow", solanaReadTx.Topology)
	assert.Equal(t, "configs/workflow-don-solana.toml", solanaReadTx.Configs)

	solanaWrite := res[3]
	assert.Equal(t, "Test_CRE_V2_Solana_Write", solanaWrite.TestName)
	assert.Equal(t, 3, solanaWrite.TestID)
	assert.Equal(t, "workflow", solanaWrite.Topology)
	assert.Equal(t, "configs/workflow-don-solana.toml", solanaWrite.Configs)

	// Test_CRE_V2_Stellar_Suite has per-test override
	stellar := res[4]
	assert.Equal(t, "Test_CRE_V2_Stellar_Suite", stellar.TestName)
	assert.Equal(t, 4, stellar.TestID)
	assert.Equal(t, "workflow-gateway-stellar", stellar.Topology)
	assert.Equal(t, "configs/workflow-gateway-don-stellar.toml", stellar.Configs)

	// Test_CRE_V2_Suite_Bucket_B has two topology entries
	bucketB := res[5]
	assert.Equal(t, "Test_CRE_V2_Suite_Bucket_B", bucketB.TestName)
	assert.Equal(t, 5, bucketB.TestID)
	assert.Equal(t, "workflow-gateway-capabilities", bucketB.Topology)
	assert.Equal(t, "configs/workflow-gateway-capabilities-don.toml", bucketB.Configs)

	bucketBVault := res[6]
	assert.Equal(t, "Test_CRE_V2_Suite_Bucket_B", bucketBVault.TestName)
	assert.Equal(t, 6, bucketBVault.TestID)
	assert.Equal(t, "workflow-gateway-capabilities-vault-stall-purge", bucketBVault.Topology)
	assert.Equal(t, "configs/workflow-gateway-capabilities-don-vault-stall-purge.toml", bucketBVault.Configs)
}

func TestBuildCRESmokeMatrix_MissingRunID(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "smoke_test.go")
	content := `package smoke_test

import "testing"

func Test_CRE_V2_Basic(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	_, err := matrix.BuildCRESmokeMatrix(context.Background(), matrix.CRESmokeOptions{
		Dir:        tmpDir,
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "run ID is required")
}

func TestBuildCRERegressionMatrix(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "regression_test.go")
	content := `package regression_test

import "testing"

func Test_CRE_V2_Standard_Regression(t *testing.T) {}
func Test_CRE_V2_Stellar_Regression(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	res, err := matrix.BuildCRERegressionMatrix(context.Background(), matrix.CRERegressionOptions{
		Dir:        tmpDir,
		RunID:      "789012",
		RunAttempt: "2",
		SpotFlag:   "spot=false",
	})
	require.NoError(t, err)
	require.Len(t, res, 2)

	assert.Equal(t, "Test_CRE_V2_Standard_Regression", res[0].TestName)
	assert.Equal(t, 0, res[0].TestID)
	assert.Equal(t, "configs/workflow-gateway-capabilities-don.toml", res[0].Configs)
	assert.Equal(t, "runs-on=789012-0-2/cpu=16/ram=64/family=m7i+m8i/spot=false/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	assert.Equal(t, "Test_CRE_V2_Stellar_Regression", res[1].TestName)
	assert.Equal(t, 1, res[1].TestID)
	assert.Equal(t, "configs/workflow-gateway-don-stellar.toml", res[1].Configs)
	assert.Equal(t, "runs-on=789012-1-2/cpu=16/ram=64/family=m7i+m8i/spot=false/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[1].RunsOn)
}

func TestBuildCRERegressionMatrix_MissingRunID(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "regression_test.go")
	content := `package regression_test

import "testing"

func Test_CRE_V2_Standard_Regression(t *testing.T) {}
`
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o600))

	_, err := matrix.BuildCRERegressionMatrix(context.Background(), matrix.CRERegressionOptions{
		Dir:        tmpDir,
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "run ID is required")
}
