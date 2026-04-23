package batching_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin/batching"
)

func TestQueryBatchCapacityCalculation(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	testMetrics := newTestMetrics(t, "query")
	queryBatch := batching.NewQueryBatch(ctx, testLogger, 100000000, testMetrics)

	for i := 0; i < 1000; i++ {
		added := queryBatch.AddRequestID(ctx, uuid.NewString())

		require.True(t, added)

		serialisedBatch, err := queryBatch.SerialiseQueryBatch()
		require.NoError(t, err)

		require.Equal(t, queryBatch.CurrentSerialisedBatchSize(), len(serialisedBatch))
	}

	require.Equal(t, 1, testMetrics.batchRequestsTotal)
	require.Equal(t, 0, testMetrics.batchCapacityExceeded)
}

func TestQueryBatchCapacityExceeded(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	testMetrics := newTestMetrics(t, "query")

	queryBatch := batching.NewQueryBatch(ctx, testLogger, 100, testMetrics)

	for i := 0; i < 1000; i++ {
		added := queryBatch.AddRequestID(ctx, uuid.NewString())

		if !added {
			require.Equal(t, 1, testMetrics.batchRequestsTotal)
			require.Equal(t, 1, testMetrics.batchCapacityExceeded)
			return
		}

		require.True(t, added)

		serialisedBatch, err := queryBatch.SerialiseQueryBatch()
		require.NoError(t, err)

		require.Equal(t, queryBatch.CurrentSerialisedBatchSize(), len(serialisedBatch))
	}

	t.Fatal("expected batch capacity to be exceeded")
}

type testMetrics struct {
	t                     *testing.T
	stepName              string
	batchRequestsTotal    int
	batchCapacityExceeded int
}

func newTestMetrics(t *testing.T, stepName string) *testMetrics {
	return &testMetrics{
		t:        t,
		stepName: stepName,
	}
}

func (tm *testMetrics) IncBatchRequestsTotal(_ context.Context, stepName string) {
	require.Equal(tm.t, tm.stepName, stepName)
	tm.batchRequestsTotal++
}

func (tm *testMetrics) IncBatchCapacityExceeded(_ context.Context, stepName string) {
	require.Equal(tm.t, tm.stepName, stepName)
	tm.batchCapacityExceeded++
}

func (tm *testMetrics) RecordObservationBatchSize(_ context.Context, size float64) {}

func (tm *testMetrics) RecordOutcomeBatchSize(_ context.Context, size float64) {}
