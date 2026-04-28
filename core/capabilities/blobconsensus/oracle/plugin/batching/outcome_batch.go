package batching

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const FailureMessageTruncated = ":TRUNCATED"

type OutcomeBatch struct {
	oracletypes.Outcome

	lggr                  logger.Logger
	metrics               metrics
	maxRequestOutcomeSize int

	seqNr                          uint64
	maxOutcomeLengthBytes          int
	currentSerialisedBatchSize     int
	initialBatchOverheadSize       int
	keybundleIDForConsensusFailure string
	maxNumberOfReports             int
}

func NewOutcomeBatch(ctx context.Context, lggr logger.Logger, seqNr uint64, previousOutcome []byte, outcomeExpirySeqNrSpan uint64, maxOutcomeLengthBytes int,
	keybundleIDForConsensusFailure string, metrics metrics, maxRequestOutcomeSize int, maxNumberOfReports int) (*OutcomeBatch, error) {
	metrics.IncBatchRequestsTotal(ctx, "outcome")
	historicalOutcomes, err := getNonExpiredHistoricalRequestOutcomes(lggr, seqNr, previousOutcome, outcomeExpirySeqNrSpan)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous outcomes: %w", err)
	}

	initialOverhead := calculateMessageSize(&oracletypes.Outcome{HistoricalOutcomes: historicalOutcomes})

	return &OutcomeBatch{
		Outcome: oracletypes.Outcome{
			HistoricalOutcomes: historicalOutcomes,
		},
		lggr:                           lggr,
		seqNr:                          seqNr,
		currentSerialisedBatchSize:     initialOverhead,
		initialBatchOverheadSize:       initialOverhead,
		keybundleIDForConsensusFailure: keybundleIDForConsensusFailure,
		metrics:                        metrics,
		maxOutcomeLengthBytes:          maxOutcomeLengthBytes,
		maxRequestOutcomeSize:          maxRequestOutcomeSize,
		maxNumberOfReports:             maxNumberOfReports,
	}, nil
}

func (o *OutcomeBatch) CurrentSerialisedBatchSize() int {
	return o.currentSerialisedBatchSize
}

// AddSuccessfulConsensusRequestOutcomeToBatch adds a successful consensus request outcome to the outcome batch.
// Returns (true, nil) if the outcome was added successfully.
// Returns (false, nil) if batch does not have capacity to add the outcome (will be retried next round).
// Returns (false, <error>) if there is a system error.
// If the outcome is too large to ever fit in any batch, a failure message (user error) will be added to the batch instead.
func (o *OutcomeBatch) AddSuccessfulConsensusRequestOutcomeToBatch(ctx context.Context, metadata *oracletypes.RequestMetaData, value *valuespb.Value, timestamp *timestamppb.Timestamp) (bool, error) {
	requestID := metadata.RequestId

	serialisedValue, err := proto.MarshalOptions{Deterministic: true}.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal successful consensus outcome value when adding request %s to batch: %w", requestID, err)
	}

	requestOutcome := &oracletypes.ConsensusOutcome{
		Outcome: &oracletypes.ConsensusOutcome_Success{
			Success: &oracletypes.ConsensusSuccessOutcome{
				Metadata:  metadata,
				Outcome:   serialisedValue,
				Timestamp: timestamp,
			},
		},
	}

	hasCapacity := o.checkOutcomeBatchHasCapacity(requestID, requestOutcome, o.seqNr)
	if !hasCapacity {
		o.metrics.IncBatchCapacityExceeded(ctx, "outcome")

		fitsInEmptyBatch, _ := o.outcomeWouldFit(o.initialBatchOverheadSize, requestID, requestOutcome, o.seqNr)
		if !fitsInEmptyBatch {
			failureMsg := "outcome too large: the consensus result for this request exceeds the maximum allowed size"
			o.lggr.Errorw(failureMsg, "requestID", requestID,
				"outcomeSize", proto.Size(requestOutcome), "maxOutcomeLengthBytes", o.maxOutcomeLengthBytes,
				"initialBatchOverheadSize", o.initialBatchOverheadSize)
			return o.AddFailedConsensusRequestOutcomeToBatch(ctx, requestID, failureMsg,
				oracletypes.ConsensusFailureCode_OUTCOME_TOO_LARGE)
		}
		return false, nil
	}

	o.Outcomes = append(o.Outcomes, requestOutcome)
	o.HistoricalOutcomes[requestID] = o.seqNr

	return true, nil
}

// AddFailedConsensusRequestOutcomeToBatch adds a failed consensus request outcome to the outcome batch. Returns false if batch does not have capacity to add the outcome.
func (o *OutcomeBatch) AddFailedConsensusRequestOutcomeToBatch(ctx context.Context, requestID, failureMessage string,
	failureCode oracletypes.ConsensusFailureCode) (bool, error) {
	requestOutcome := &oracletypes.ConsensusOutcome{
		Outcome: &oracletypes.ConsensusOutcome_Failure{
			Failure: &oracletypes.ConsensusFailedOutcome{
				RequestID:      requestID,
				KeyBundleId:    o.keybundleIDForConsensusFailure,
				FailureMessage: failureMessage,
				Code:           failureCode,
			},
		},
	}

	if proto.Size(requestOutcome) > o.maxRequestOutcomeSize {
		requestOutcome = o.truncateFailedRequestOutcome(requestOutcome.GetFailure())
	}

	hasCapacity := o.checkOutcomeBatchHasCapacity(requestID, requestOutcome, o.seqNr)
	if !hasCapacity {
		o.metrics.IncBatchCapacityExceeded(ctx, "outcome")
		return false, nil
	}

	o.Outcomes = append(o.Outcomes, requestOutcome)
	o.HistoricalOutcomes[requestID] = o.seqNr

	return true, nil
}

func (o *OutcomeBatch) truncateFailedRequestOutcome(failedRequestOutcome *oracletypes.ConsensusFailedOutcome) *oracletypes.ConsensusOutcome {
	// Truncate the failure message so the request outcome fits within the request outcome size limit
	requestOutcomeWithoutFailureMessage := &oracletypes.ConsensusOutcome{Outcome: &oracletypes.ConsensusOutcome_Failure{
		Failure: &oracletypes.ConsensusFailedOutcome{
			RequestID:   failedRequestOutcome.RequestID,
			KeyBundleId: failedRequestOutcome.KeyBundleId,
			Code:        failedRequestOutcome.Code,
		},
	},
	}
	sizeOfRequestOutcomeWithoutFailureMessage := proto.Size(requestOutcomeWithoutFailureMessage)

	allowedFailureMessageSize := o.maxRequestOutcomeSize - sizeOfRequestOutcomeWithoutFailureMessage - len(FailureMessageTruncated)

	if allowedFailureMessageSize <= 0 {
		o.lggr.Warn("unable to fit any part of failure message within request outcome size limit", "requestID", failedRequestOutcome.RequestID,
			"maxRequestOutcomeSize", o.maxRequestOutcomeSize, "failure message", failedRequestOutcome.FailureMessage)
		return requestOutcomeWithoutFailureMessage
	}

	truncatedFailureMessage := failedRequestOutcome.FailureMessage[:allowedFailureMessageSize] + FailureMessageTruncated
	o.lggr.Warnw("truncated failure message to fit within request outcome size limit", "requestID", failedRequestOutcome.RequestID,
		"originalSize", len(failedRequestOutcome.FailureMessage), "truncatedSize", len(truncatedFailureMessage),
		"maxRequestOutcomeSize", o.maxRequestOutcomeSize, "failure message", failedRequestOutcome.FailureMessage)

	return &oracletypes.ConsensusOutcome{
		Outcome: &oracletypes.ConsensusOutcome_Failure{
			Failure: &oracletypes.ConsensusFailedOutcome{
				RequestID:      failedRequestOutcome.RequestID,
				KeyBundleId:    failedRequestOutcome.KeyBundleId,
				FailureMessage: truncatedFailureMessage,
				Code:           failedRequestOutcome.Code,
			},
		},
	}
}

// FailConsensusWithDefaultCheck handles a consensus failure by checking if a default value is available to use.
// If a default value is available, it adds a successful consensus outcome with the default value to the batch.
// If no default value is available, it adds a failed consensus outcome to the batch.
func (o *OutcomeBatch) FailConsensusWithDefaultCheck(ctx context.Context, lggr logger.Logger, requestID string, fullConsensusFailedMsg string, logMsg string,
	code oracletypes.ConsensusFailureCode, consensusMDD *oracletypes.RequestObservation, timestamp *timestamppb.Timestamp) (bool, error) {
	lggr.Debugw("consensus failed with a user error", "reason", logMsg, "requestID", requestID)

	defaultVal, err := values.FromProto(consensusMDD.Input.Default)
	if err != nil {
		errMsg := fmt.Sprintf("could not convert default value from proto for request %s: %v", requestID, err)
		lggr.Error(errMsg)
		return o.AddFailedConsensusRequestOutcomeToBatch(ctx, requestID, errMsg, code)
	}

	if defaultVal != nil {
		lggr.Debugw("using default value for request", "requestID", requestID, "defaultValue", defaultVal)
		return o.AddSuccessfulConsensusRequestOutcomeToBatch(ctx, consensusMDD.Metadata, consensusMDD.Input.Default, timestamp)
	}

	return o.AddFailedConsensusRequestOutcomeToBatch(ctx, requestID, fullConsensusFailedMsg, code)
}

func (o *OutcomeBatch) SerialiseOutcomeBatch(ctx context.Context) ([]byte, error) {
	serialisedBatch, err := proto.MarshalOptions{Deterministic: true}.Marshal(&o.Outcome)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise batch of outcomes: %w", err)
	}

	batchSize := len(serialisedBatch)
	o.metrics.RecordOutcomeBatchSize(ctx, float64(batchSize))
	o.lggr.Debugw("serialised outcome batch", "numOutcomes", len(o.Outcomes),
		"actualSizeBytes", batchSize, "calculatedSizeBytes", o.currentSerialisedBatchSize,
		"maxOutcomeLengthBytes", o.maxOutcomeLengthBytes)

	return serialisedBatch, nil
}

// outcomeWouldFit checks if an outcome (message + historical map entry) would fit given a starting batch size.
// Returns (fits, newSize) where newSize is the resulting batch size if it fits.
func (o *OutcomeBatch) outcomeWouldFit(startingSize int, requestID string, requestOutcome proto.Message, historicalSeqNr uint64) (bool, int) {
	ok, newSize := batchHasCapacityForMessageOnSlice(startingSize, requestOutcome, o.maxOutcomeLengthBytes)
	if !ok {
		return false, startingSize
	}

	// Adding an entry to a map is the same as adding a key-value pair to a slice for size calculation purposes
	mapEntry := oracletypes.HistoricalOutcomeMapEntry{
		Key:   requestID,
		Value: historicalSeqNr,
	}
	mapEntrySize := proto.Size(&mapEntry)

	ok, newSize = batchHasCapacityForSliceBytes(newSize, mapEntrySize, o.maxOutcomeLengthBytes)
	if !ok {
		return false, startingSize
	}

	return true, newSize
}

func (o *OutcomeBatch) checkOutcomeBatchHasCapacity(requestID string, requestOutcome proto.Message,
	historicalSeqNr uint64) bool {

	// The number of outcomes should never exceed the max number of reports as this would cause outcomes to be lost and need recalculating on subsequent rounds
	if len(o.Outcomes) == o.maxNumberOfReports {
		o.lggr.Debugw("outcome batch has reached maximum number of reports", "maxNumberOfReports", o.maxNumberOfReports)
		return false
	}

	fits, newSize := o.outcomeWouldFit(o.currentSerialisedBatchSize, requestID, requestOutcome, historicalSeqNr)
	if fits {
		o.currentSerialisedBatchSize = newSize
	}
	return fits
}

func getNonExpiredHistoricalRequestOutcomes(lggr logger.Logger, seqNr uint64, previousOutcome []byte, outcomeExpirySeqNrSpan uint64) (map[string]uint64, error) {
	nonExpiredHistoricalOutcomes := map[string]uint64{}
	if len(previousOutcome) > 0 {
		prevOutcome := &oracletypes.Outcome{}
		err := proto.Unmarshal(previousOutcome, prevOutcome)
		if err != nil {
			lggr.Errorw("could not unmarshal previous outcome", "error", err)
			return nil, err
		}

		for requestID, outcomeSeqNr := range prevOutcome.HistoricalOutcomes {
			if seqNr-outcomeSeqNr <= outcomeExpirySeqNrSpan {
				nonExpiredHistoricalOutcomes[requestID] = outcomeSeqNr
			}
		}
	}

	return nonExpiredHistoricalOutcomes, nil
}
