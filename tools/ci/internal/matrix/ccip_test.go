package matrix_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func TestBuildCCIPSystemMatrix(t *testing.T) {
	t.Parallel()

	res, err := matrix.BuildCCIPSystemMatrix(context.Background(), matrix.CCIPSystemOptions{
		RunID:      "112233",
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.NoError(t, err)
	require.Len(t, res, 5)

	assert.Equal(t, "Test_CCIPGasPriceUpdatesWriteFrequency", res[0].TestName)
	assert.Equal(t, 0, res[0].TestID)
	assert.Equal(t, "15m", res[0].Timeout)
	assert.Equal(t, "SIMULATED_1,SIMULATED_2", res[0].SelectedNetwork)
	assert.False(t, res[0].MixedVersion)
	assert.Equal(t, "runs-on=112233-0-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	assert.Equal(t, "TestRMN_GlobalCurseTwoMessagesOnTwoLanes", res[1].TestName)
	assert.Equal(t, 1, res[1].TestID)
	assert.Equal(t, "master-amd6416f5d86", res[1].RMNRageProxyVersion)
	assert.Equal(t, "master-amd64-10b42b2", res[1].RMNAFN2ProxyVersion)
	assert.False(t, res[1].MixedVersion)
	assert.Equal(t, "runs-on=112233-1-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[1].RunsOn)

	assert.Equal(t, "TestDeleteCCIPJobs-TestRevokeJobs", res[2].TestName)
	assert.Equal(t, 2, res[2].TestID)
	assert.Equal(t, 20, res[2].JobTimeout)
	assert.False(t, res[2].MixedVersion)
	assert.Equal(t, "runs-on=112233-2-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[2].RunsOn)

	assert.Equal(t, "Test_CCIPMixedVersionDON", res[3].TestName)
	assert.Equal(t, 3, res[3].TestID)
	assert.Equal(t, "20m", res[3].Timeout)
	assert.Equal(t, "SIMULATED_1,SIMULATED_2", res[3].SelectedNetwork)
	assert.True(t, res[3].MixedVersion)
	assert.Equal(t, "runs-on=112233-3-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[3].RunsOn)

	assert.Equal(t, "Test_CCIPRollingUpgrade", res[4].TestName)
	assert.Equal(t, 4, res[4].TestID)
	assert.Equal(t, "30m", res[4].Timeout)
	assert.Equal(t, "SIMULATED_1,SIMULATED_2", res[4].SelectedNetwork)
	assert.True(t, res[4].MixedVersion)
	assert.Equal(t, "runs-on=112233-4-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[4].RunsOn)
}

func TestBuildCCIPSystemMatrix_MissingRunID(t *testing.T) {
	t.Parallel()

	_, err := matrix.BuildCCIPSystemMatrix(context.Background(), matrix.CCIPSystemOptions{
		RunAttempt: "1",
		SpotFlag:   "spot=co",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "run ID is required")
}

func TestBuildCCIPSystemMatrix_Defaults(t *testing.T) {
	t.Parallel()

	res, err := matrix.BuildCCIPSystemMatrix(context.Background(), matrix.CCIPSystemOptions{
		RunID: "112233",
	})
	require.NoError(t, err)
	require.Len(t, res, 5)

	assert.Equal(t, "runs-on=112233-0-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)
}

func TestBuildCCIPSystemMatrix_MixedVersionOnly(t *testing.T) {
	t.Parallel()

	res, err := matrix.BuildCCIPSystemMatrix(context.Background(), matrix.CCIPSystemOptions{
		RunID:            "112233",
		RunAttempt:       "1",
		SpotFlag:         "spot=co",
		MixedVersionOnly: true,
	})
	require.NoError(t, err)
	require.Len(t, res, 2)

	assert.Equal(t, "Test_CCIPMixedVersionDON", res[0].TestName)
	assert.Equal(t, 3, res[0].TestID)
	assert.Equal(t, "runs-on=112233-3-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	assert.Equal(t, "Test_CCIPRollingUpgrade", res[1].TestName)
	assert.Equal(t, 4, res[1].TestID)
	assert.Equal(t, "runs-on=112233-4-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[1].RunsOn)
}
