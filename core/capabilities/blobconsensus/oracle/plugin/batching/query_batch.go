package batching

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"
)

type QueryBatch struct {
	oracletypes.Query

	lggr                       logger.Logger
	currentSerialisedBatchSize int
	metrics                    metrics
	maxQueryLengthBytes        int
}

func NewQueryBatch(ctx context.Context, lggr logger.Logger, maxQueryLengthBytes int, metrics metrics) *QueryBatch {
	metrics.IncBatchRequestsTotal(ctx, "query")
	obs := &oracletypes.Query{RequestIDs: nil}
	messageSize := calculateMessageSize(obs)

	return &QueryBatch{
		Query:                      oracletypes.Query{RequestIDs: nil},
		lggr:                       lggr,
		currentSerialisedBatchSize: messageSize,
		maxQueryLengthBytes:        maxQueryLengthBytes,
		metrics:                    metrics,
	}
}

func (qb *QueryBatch) AddRequestID(ctx context.Context, requestID string) bool {
	hasCapacity, newSize := batchHasCapacityForStringOnSlice(qb.currentSerialisedBatchSize, requestID, qb.maxQueryLengthBytes)

	if !hasCapacity {
		qb.metrics.IncBatchCapacityExceeded(ctx, "query")
		return false
	}

	qb.RequestIDs = append(qb.RequestIDs, requestID)
	qb.currentSerialisedBatchSize = newSize

	return true
}

func (qb *QueryBatch) CurrentSerialisedBatchSize() int {
	return qb.currentSerialisedBatchSize
}

func (qb *QueryBatch) SerialiseQueryBatch() ([]byte, error) {
	serialisedBatch, err := proto.MarshalOptions{Deterministic: true}.Marshal(&qb.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise batch of request ids: %w", err)
	}

	qb.lggr.Debugw("serialised batch of request ids", "RequestIDs", qb.RequestIDs,
		"actualSizeBytes", len(serialisedBatch), "calculatedSizeBytes", qb.currentSerialisedBatchSize,
		"maxQueryLengthBytes", qb.maxQueryLengthBytes)

	return serialisedBatch, nil
}

func (qb *QueryBatch) NumberOfRequestIDs() int {
	return len(qb.RequestIDs)
}
