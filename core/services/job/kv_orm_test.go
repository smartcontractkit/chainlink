package job_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

func TestJobKVStore(t *testing.T) {
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	lggr := logger.TestLogger(t)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())

	jobID := int32(1337)
	kvStore := job.NewKVStore(jobID, db, config.Database(), lggr)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, cltest.NewKeyStore(t, db, config.Database()), config.Database())

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb.ID = jobID
	require.NoError(t, jobORM.CreateJob(&jb))

	key := "test_key"

	type testData struct {
		Test string `json:"test"`
	}

	td1 := testData{Test: "value1"}
	td2 := testData{Test: "value2"}

	var retData testData
	require.NoError(t, kvStore.Store(key, td1))
	require.NoError(t, kvStore.Get(key, &retData))
	require.Equal(t, td1, retData)

	require.NoError(t, kvStore.Store(key, td2))
	require.NoError(t, kvStore.Get(key, &retData))
	require.Equal(t, td2, retData)
}
