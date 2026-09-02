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
	require.Len(t, res, 3)

	assert.Equal(t, "Test_CCIPGasPriceUpdatesWriteFrequency", res[0].TestName)
	assert.Equal(t, 0, res[0].TestID)
	assert.Equal(t, "15m", res[0].Timeout)
	assert.Equal(t, "SIMULATED_1,SIMULATED_2", res[0].SelectedNetwork)
	assert.Equal(t, "runs-on=112233-0-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[0].RunsOn)

	assert.Equal(t, "TestRMN_GlobalCurseTwoMessagesOnTwoLanes", res[1].TestName)
	assert.Equal(t, 1, res[1].TestID)
	assert.Equal(t, "master-amd6416f5d86", res[1].RMNRageProxyVersion)
	assert.Equal(t, "master-amd64-10b42b2", res[1].RMNAFN2ProxyVersion)
	assert.Equal(t, "runs-on=112233-1-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[1].RunsOn)

	assert.Equal(t, "TestDeleteCCIPJobs-TestRevokeJobs", res[2].TestName)
	assert.Equal(t, 2, res[2].TestID)
	assert.Equal(t, 20, res[2].JobTimeout)
	assert.Equal(t, "runs-on=112233-2-1/cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", res[2].RunsOn)
}
