package batching

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type ObservationBatch struct {
	oracletypes.Observation

	lggr                       logger.Logger
	currentSerialisedBatchSize int
	metrics                    metrics
	maxObservationLengthBytes  int
}

func NewObservationBatch(ctx context.Context, lggr logger.Logger, maxObservationLengthBytes int, metrics metrics) *ObservationBatch {
	metrics.IncBatchRequestsTotal(ctx, "observation")
	observations := make(map[string]*oracletypes.RequestObservation)
	obs := &oracletypes.Observation{Observations: observations}
	messageSize := calculateMessageSize(obs)

	return &ObservationBatch{
		Observation:                oracletypes.Observation{Observations: observations},
		lggr:                       lggr,
		currentSerialisedBatchSize: messageSize,
		maxObservationLengthBytes:  maxObservationLengthBytes,
		metrics:                    metrics,
	}
}

func (ob *ObservationBatch) AddObservation(ctx context.Context, reqObs *oracletypes.RequestObservation) bool {
	// Adding an entry to a map is the same as adding a key-value pair to a slice for size calculation purposes
	mapEntry := oracletypes.ObservationMapEntry{
		Key:   reqObs.Metadata.RequestId,
		Value: reqObs,
	}

	mapEntrySize := proto.Size(&mapEntry)
	ok, newSize := batchHasCapacityForSliceBytes(ob.currentSerialisedBatchSize, mapEntrySize, ob.maxObservationLengthBytes)

	if !ok {
		ob.metrics.IncBatchCapacityExceeded(ctx, "observation")
		// Observation doesn't fit in current batch - will retry next round when batch is empty
		return false
	}

	ob.currentSerialisedBatchSize = newSize
	ob.Observations[reqObs.Metadata.RequestId] = reqObs

	return true
}

func (ob *ObservationBatch) CurrentSerialisedBatchSize() int {
	return ob.currentSerialisedBatchSize
}

func (ob *ObservationBatch) SerialiseObservationBatch(ctx context.Context) ([]byte, error) {
	serialisedBatch, err := proto.MarshalOptions{Deterministic: true}.Marshal(&ob.Observation)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise batch of observations: %w", err)
	}

	allIDs := make([]string, 0, len(ob.Observations))
	for id := range ob.Observations {
		allIDs = append(allIDs, id)
	}

	batchSize := len(serialisedBatch)
	ob.metrics.RecordObservationBatchSize(ctx, float64(batchSize))
	ob.lggr.Debugw("serialised observation batch", "numObservations", len(ob.Observations), "actualSizeBytes", batchSize,
		"calculatedSizeBytes", ob.currentSerialisedBatchSize,
		"maxObservationLengthBytes", ob.maxObservationLengthBytes,
		"executionIDs", allIDs,
	)

	return serialisedBatch, nil
}

func (ob *ObservationBatch) NumObservationsInBatch() int {
	return len(ob.Observations)
}
