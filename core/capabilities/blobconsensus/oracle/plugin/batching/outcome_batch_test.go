package batching_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin/batching"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"
)

func TestOutcomeBatchCapacityCalculation(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	prevOutcome := &oracletypes.Outcome{
		HistoricalOutcomes: map[string]uint64{
			"req-1": 10,
			"req-2": 20,
			"req-3": 30,
		},
	}

	serialisedPrevOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(prevOutcome)
	require.NoError(t, err)

	testMetrics := newTestMetrics(t, "outcome")
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: serialisedPrevOutcome,
		SeqNr:           1000,
	}, 1000,
		100_000_000, "evm", testMetrics, 1000, 10000)

	require.NoError(t, err)

	for range 1000 {
		added, err := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
			RequestId:           uuid.NewString(),
			WorkflowExecutionId: generateRandomStringBetweenBounds(1, 10000),
		}, values.Proto(values.NewString("test-outcome-data-1")), &timestamppb.Timestamp{})

		require.True(t, added)
		require.NoError(t, err)

		added, err = outcome.AddFailedConsensusRequestOutcomeToBatch(ctx, uuid.NewString(), generateRandomStringBetweenBounds(1, 10000),
			oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED)

		require.True(t, added)
		require.NoError(t, err)

		serialisedBatch, err := outcome.SerialiseOutcomeBatch(t.Context())
		require.NoError(t, err)

		require.Equal(t, outcome.CurrentSerialisedBatchSize(), len(serialisedBatch))
	}

	require.Equal(t, 1, testMetrics.batchRequestsTotal)
	require.Equal(t, 0, testMetrics.batchCapacityExceeded)
}

func TestOutcomeBatchCapacityExceeded(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	prevOutcome := &oracletypes.Outcome{
		HistoricalOutcomes: map[string]uint64{
			"req-1": 10,
			"req-2": 20,
			"req-3": 30,
		},
	}

	serialisedPrevOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(prevOutcome)
	require.NoError(t, err)

	testMetrics := newTestMetrics(t, "outcome")
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: serialisedPrevOutcome,
		SeqNr:           1000,
	}, 1000,
		1000, "evm", testMetrics, 1000, 10000)

	require.NoError(t, err)

	for range 10 {
		added, err := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
			RequestId:           uuid.NewString(),
			WorkflowExecutionId: "exec-1",
		}, values.Proto(values.NewString("test-outcome-data-1")), &timestamppb.Timestamp{})

		if !added {
			require.Equal(t, 1, testMetrics.batchRequestsTotal)
			require.Equal(t, 1, testMetrics.batchCapacityExceeded)
			return
		}
		require.NoError(t, err)

		added, err = outcome.AddFailedConsensusRequestOutcomeToBatch(ctx, uuid.NewString(), "failed",
			oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED)

		if !added {
			return
		}
		require.NoError(t, err)

		serialisedBatch, err := outcome.SerialiseOutcomeBatch(t.Context())
		require.NoError(t, err)

		require.Equal(t, outcome.CurrentSerialisedBatchSize(), len(serialisedBatch))
	}

	t.Fatal("expected batch capacity to be exceeded")
}

func TestOutcomeBatchMaxNumberOfReportsExceeded(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	prevOutcome := &oracletypes.Outcome{
		HistoricalOutcomes: map[string]uint64{
			"req-1": 10,
			"req-2": 20,
			"req-3": 30,
		},
	}

	serialisedPrevOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(prevOutcome)
	require.NoError(t, err)

	testMetrics := newTestMetrics(t, "outcome")
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: serialisedPrevOutcome,
		SeqNr:           1000,
	}, 1000,
		100000, "evm", testMetrics, 100000, 3)

	require.NoError(t, err)

	for range 10 {
		added, err := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
			RequestId:           uuid.NewString(),
			WorkflowExecutionId: "exec-1",
		}, values.Proto(values.NewString("test-outcome-data-1")), &timestamppb.Timestamp{})

		if !added {
			require.Equal(t, 1, testMetrics.batchRequestsTotal)
			require.Equal(t, 1, testMetrics.batchCapacityExceeded)
			return
		}
		require.NoError(t, err)

		added, err = outcome.AddFailedConsensusRequestOutcomeToBatch(ctx, uuid.NewString(), "failed",
			oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED)

		if !added {
			return
		}
		require.NoError(t, err)

		serialisedBatch, err := outcome.SerialiseOutcomeBatch(t.Context())
		require.NoError(t, err)

		require.Equal(t, outcome.CurrentSerialisedBatchSize(), len(serialisedBatch))
	}

	t.Fatal("expected max number of reports to be exceeded")
}

func TestOutcomeTooLargeToEverFit(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	testMetrics := newTestMetrics(t, "outcome")
	// Create a batch with a very small max size (500 bytes)
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: nil,
		SeqNr:           1,
	}, 1000,
		500, "evm", testMetrics, 1000, 10000)

	require.NoError(t, err)

	// Create a value that is definitely too large to ever fit
	largeData := strings.Repeat("x", 1000)
	added, _ := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
		RequestId:           "req-too-large",
		WorkflowExecutionId: "exec-1",
	}, values.Proto(values.NewString(largeData)), &timestamppb.Timestamp{})

	require.True(t, added) // still added but a user erro, not the large outcome
	require.Equal(t, oracletypes.ConsensusFailureCode_OUTCOME_TOO_LARGE, outcome.Outcomes[len(outcome.Outcomes)-1].GetFailure().GetCode())
	require.Equal(t, 1, testMetrics.batchCapacityExceeded)
}

func TestOutcomeDoesNotFitNowButWouldFitInEmptyBatch(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	testMetrics := newTestMetrics(t, "outcome")
	// Create a batch with just enough space for one outcome but not two.
	// Based on testing, an outcome with small data takes ~32 bytes.
	// With max size 50, one fits (32 < 50) but two don't (64 > 50).
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: nil,
		SeqNr:           1,
	}, 1000,
		50, "evm", testMetrics, 1000, 10000)

	require.NoError(t, err)

	// Add a first outcome that should fit
	added, err := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
		RequestId:           "r1",
		WorkflowExecutionId: "e1",
	}, values.Proto(values.NewString("d")), &timestamppb.Timestamp{})

	require.True(t, added, "first outcome should fit")
	require.NoError(t, err)
	require.Equal(t, 0, testMetrics.batchCapacityExceeded)

	// Add a second outcome - should not fit now, but would fit in an empty batch
	// This should return (false, nil) - not an error
	added, err = outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
		RequestId:           "r2",
		WorkflowExecutionId: "e2",
	}, values.Proto(values.NewString("d")), &timestamppb.Timestamp{})

	require.False(t, added, "second outcome should not fit in current batch")
	require.NoError(t, err, "should not return error when outcome would fit in empty batch")
	require.Equal(t, 1, testMetrics.batchCapacityExceeded)
}

func TestOutcomeTooLargeWithExistingHistoricalOutcomes(t *testing.T) {
	testLogger := logger.Test(t)
	ctx := t.Context()

	// Create previous outcome with some historical data
	prevOutcome := &oracletypes.Outcome{
		HistoricalOutcomes: map[string]uint64{
			"old-req-1": 10,
			"old-req-2": 20,
		},
	}
	serialisedPrevOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(prevOutcome)
	require.NoError(t, err)

	testMetrics := newTestMetrics(t, "outcome")
	// Create a batch with a small max size
	outcome, err := batching.NewOutcomeBatch(ctx, testLogger, ocr3types.OutcomeContext{
		PreviousOutcome: serialisedPrevOutcome,
		SeqNr:           100,
	}, 1000,
		500, "evm", testMetrics, 1000, 10000)

	require.NoError(t, err)

	// Try to add an outcome that is too large even for an empty batch
	largeData := strings.Repeat("y", 1000)
	added, _ := outcome.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, &oracletypes.RequestMetaData{
		RequestId:           "req-too-large",
		WorkflowExecutionId: "exec-1",
	}, values.Proto(values.NewString(largeData)), &timestamppb.Timestamp{})

	require.True(t, added) // still added but a user erro, not the large outcome
	require.Equal(t, oracletypes.ConsensusFailureCode_OUTCOME_TOO_LARGE, outcome.Outcomes[len(outcome.Outcomes)-1].GetFailure().GetCode())
	require.Equal(t, 1, testMetrics.batchCapacityExceeded)
}
