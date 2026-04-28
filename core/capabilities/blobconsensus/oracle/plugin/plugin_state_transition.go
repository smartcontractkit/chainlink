package plugin

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin/batching"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func (r *reportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, attributedObservations []types.AttributedObservation,
	keyValueStateReadWriter ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	lggr := logger.With(r.lggr, "seqNr", seqNr)

	requestsQuery := &oracletypes.Query{}
	err := proto.Unmarshal(aq.Query, requestsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}

	previousOutcome, err := readPreviousOutcomeBytesFromKV(keyValueStateReadWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to read previous outcome: %w", err)
	}

	outcomeBatch, err := batching.NewOutcomeBatch(ctx, lggr, seqNr, previousOutcome, r.outcomeExpirySeqNrSpan, int(r.config.MaxOutcomeLengthBytes), r.defaultKeyBundleIDForConsensusFailure,
		r.metrics, r.maxRequestOutcomeSize, r.maxNumberOfReports)
	if err != nil {
		return nil, fmt.Errorf("failed to create new outcome batch: %w", err)
	}

	requestIDToObservations := groupAttributedObservationsByRequestID(lggr, attributedObservations)

	for _, requestID := range requestsQuery.RequestIDs {
		observations := requestIDToObservations[requestID]

		// 2f+1 or more observations have been received, calculate the outcome for the request
		if len(observations) >= 2*r.f+1 {
			hasCapacity, err := r.addRequestOutcomeToBatch(ctx, lggr, requestID, observations, outcomeBatch)
			if err != nil {
				return nil, fmt.Errorf("failed to add request outcome to batch for request %s: %w", requestID, err)
			}

			if !hasCapacity {
				lggr.Debugw("batch does not have capacity to add request outcome - skipping in this round", "requestID", requestID)
				break
			}
			lggr.Debugw("added request outcome to batch", "requestID", requestID, "numObservations", len(observations))
		} else {
			lggr.Debugw("not enough observations to calculate outcome for request - skipping in this round", "requestID", requestID, "numObservations", len(observations))
		}
	}

	serialised, err := outcomeBatch.SerialiseOutcomeBatch(ctx)
	if err != nil {
		return nil, err
	}
	if err := keyValueStateReadWriter.Write(prevOutcomeStateKey, serialised); err != nil {
		return nil, fmt.Errorf("failed to persist outcome to key-value state: %w", err)
	}
	return ocr3_1types.ReportsPlusPrecursor(serialised), nil
}

func (r *reportingPlugin) Committed(_ context.Context, _ uint64, _ ocr3_1types.KeyValueStateReader) error {
	return nil
}
