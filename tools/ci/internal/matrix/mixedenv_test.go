package matrix_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func TestBuildCREMixedEnvMatrix(t *testing.T) {
	t.Parallel()

	res, err := matrix.BuildCREMixedEnvMatrix(context.Background(), matrix.CREMixedEnvOptions{
		RunID:      "445566",
		RunAttempt: "1",
		SpotFlag:   "spot=pco",
	})
	require.NoError(t, err)
	require.Len(t, res, 5)

	// Bucket A with default configs
	assert.Equal(t, "Test_CRE_V2_Suite_Bucket_A", res[0].TestName)
	assert.Equal(t, 0, res[0].TestID)
	assert.Equal(t, "configs/mixed-env-don.toml", res[0].Configs)
	assert.Equal(t, "runs-on=445566-0-1/cpu=16/ram=64/family=m7i+m8i/spot=pco/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	// TxArtifacts with default configs
	assert.Equal(t, "Test_CRE_V2_EVM_Read_TxArtifacts", res[4].TestName)
	assert.Equal(t, 4, res[4].TestID)
	assert.Equal(t, "configs/mixed-env-don.toml", res[4].Configs)
	assert.Equal(t, "runs-on=445566-4-1/cpu=16/ram=64/family=m7i+m8i/spot=pco/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[4].RunsOn)
}

func TestBuildCREMixedEnvMatrix_MissingRunID(t *testing.T) {
	t.Parallel()

	_, err := matrix.BuildCREMixedEnvMatrix(context.Background(), matrix.CREMixedEnvOptions{
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "run ID is required")
}
