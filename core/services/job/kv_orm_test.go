package job_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

func TestJobKVStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	lggr := logger.TestLogger(t)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobID := int32(1337)
	kvStore := job.NewKVStore(jobID, db, lggr)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, cltest.NewKeyStore(t, db))

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb.ID = jobID
	require.NoError(t, jobORM.CreateJob(testutils.Context(t), &jb))

	var values = [][]byte{
		[]byte("Hello"),
		[]byte("World"),
		[]byte("Go"),
	}

	for i, insertBytes := range values {
		testKey := "test_key_" + fmt.Sprint(i)
		require.NoError(t, kvStore.Store(ctx, testKey, insertBytes))

		var readBytes []byte
		readBytes, err = kvStore.Get(ctx, testKey)
		assert.NoError(t, err)

		require.Equal(t, insertBytes, readBytes)
	}

	key := "test_key_updating"
	td1 := []byte("value1")
	td2 := []byte("value2")

	require.NoError(t, kvStore.Store(ctx, key, td1))
	fetchedBytes, err := kvStore.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, td1, fetchedBytes)

	require.NoError(t, kvStore.Store(ctx, key, td2))
	fetchedBytes, err = kvStore.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, td2, fetchedBytes)

	require.NoError(t, jobORM.DeleteJob(ctx, jobID))
}
