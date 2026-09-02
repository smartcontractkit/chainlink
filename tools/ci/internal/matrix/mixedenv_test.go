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
	require.Len(t, res, 6)

	// Bucket A with default configs
	assert.Equal(t, "Test_CRE_V2_Suite_Bucket_A", res[0].TestName)
	assert.Equal(t, 0, res[0].TestID)
	assert.Equal(t, "configs/mixed-env-don.toml", res[0].Configs)
	assert.Equal(t, "runs-on=445566-0-1/cpu=16/ram=64/family=m7i+m8i/spot=pco/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	// ConfidentialWorkflows with override configs
	assert.Equal(t, "Test_CRE_V2_ConfidentialWorkflows_Relay", res[5].TestName)
	assert.Equal(t, 5, res[5].TestID)
	assert.Equal(t, "configs/mixed-env-confidential-workflows.toml", res[5].Configs)
	assert.Equal(t, "runs-on=445566-5-1/cpu=16/ram=64/family=m7i+m8i/spot=pco/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[5].RunsOn)
}
