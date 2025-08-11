package job_test

import (
	"context"
	"strconv"
	"testing"
	"time"

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

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobID := int32(1337)
	kvStore := job.NewKVStore(jobID, db)
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
		testKey := "test_key_" + strconv.Itoa(i)
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

	require.NoError(t, jobORM.DeleteJob(ctx, jobID, jb.Type))
}

func TestJobKVStore_PruneExpiredEntries(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, cltest.NewKeyStore(t, db))

	// Create two test jobs
	jobID1 := int32(1337)
	jobID2 := int32(1338)
	kvStore1 := job.NewKVStore(jobID1, db)
	kvStore2 := job.NewKVStore(jobID2, db)

	jb1, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb1.ID = jobID1
	require.NoError(t, jobORM.CreateJob(testutils.Context(t), &jb1))

	jb2, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb2.ID = jobID2
	require.NoError(t, jobORM.CreateJob(testutils.Context(t), &jb2))

	testData := []struct {
		key   string
		value []byte
	}{
		{"old_key_1", []byte("old_value_1")},
		{"old_key_2", []byte("old_value_2")},
		{"new_key_1", []byte("new_value_1")},
		{"new_key_2", []byte("new_value_2")},
	}

	for i := 0; i < 2; i++ {
		require.NoError(t, kvStore1.Store(ctx, testData[i].key, testData[i].value))
	}

	for i := 0; i < 2; i++ {
		require.NoError(t, kvStore2.Store(ctx, testData[i].key, testData[i].value))
	}

	// Simulate old entries by updating updated_at to be in the past
	pastTime := time.Now().Add(-2 * time.Hour)
	_, err = db.ExecContext(ctx,
		"UPDATE job_kv_store SET updated_at = $1 WHERE key LIKE 'old_key_%'",
		pastTime)
	require.NoError(t, err)

	// Store recent entries for both jobs
	for i := 2; i < 4; i++ {
		require.NoError(t, kvStore1.Store(ctx, testData[i].key, testData[i].value))
		require.NoError(t, kvStore2.Store(ctx, testData[i].key, testData[i].value))
	}

	t.Run("PruneExpiredEntries for specific job", func(t *testing.T) {
		maxAge := 1 * time.Hour
		deletedCount, err := kvStore1.PruneExpiredEntries(ctx, maxAge)
		require.NoError(t, err)
		require.Equal(t, int64(2), deletedCount, "Should delete 2 old entries for job 1")

		_, err = kvStore1.Get(ctx, "old_key_1")
		require.Error(t, err, "old_key_1 should be deleted")
		_, err = kvStore1.Get(ctx, "old_key_2")
		require.Error(t, err, "old_key_2 should be deleted")

		val, err := kvStore1.Get(ctx, "new_key_1")
		require.NoError(t, err)
		require.Equal(t, []byte("new_value_1"), val)
		val, err = kvStore2.Get(ctx, "old_key_1")
		require.NoError(t, err)
		require.Equal(t, []byte("old_value_1"), val)
	})

	t.Run("PruneExpiredEntries for job 2", func(t *testing.T) {
		maxAge := 1 * time.Hour
		deletedCount, err := kvStore2.PruneExpiredEntries(ctx, maxAge)
		require.NoError(t, err)
		require.Equal(t, int64(2), deletedCount, "Should delete 2 old entries for job 2")

		_, err = kvStore2.Get(ctx, "old_key_1")
		require.Error(t, err, "old_key_1 should be deleted")
		_, err = kvStore2.Get(ctx, "old_key_2")
		require.Error(t, err, "old_key_2 should be deleted")
		val, err := kvStore2.Get(ctx, "new_key_1")
		require.NoError(t, err)
		require.Equal(t, []byte("new_value_1"), val)
	})

	t.Run("PruneExpiredEntries with no expired entries", func(t *testing.T) {
		maxAge := 5 * time.Hour
		deletedCount, err := kvStore1.PruneExpiredEntries(ctx, maxAge)
		require.NoError(t, err)
		require.Equal(t, int64(0), deletedCount, "Should delete no entries")
	})

	// Cleanup
	require.NoError(t, jobORM.DeleteJob(ctx, jobID1, jb1.Type))
	require.NoError(t, jobORM.DeleteJob(ctx, jobID2, jb2.Type))
}
