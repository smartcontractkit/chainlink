package plugin

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin/batching"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

func (r *reportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	_ = seqNr
	allRequests, err := r.store.All()
	if err != nil {
		return nil, fmt.Errorf("failed to get all requests: %w", err)
	}

	previousOutcome, err := readPreviousOutcomeBytesFromKV(keyValueStateReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read previous outcome from key-value state: %w", err)
	}

	// Get only those requests that are pending consensus to prevent completed requests being included in the new query
	pendingRequests, err := r.getPendingRequestsFromPreviousOutcome(previousOutcome, allRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to remove completed requests: %w", err)
	}

	queryBatch := batching.NewQueryBatch(ctx, r.lggr, int(r.config.MaxQueryLengthBytes), r.metrics)

	for _, rq := range pendingRequests {
		hasCapacity := queryBatch.AddRequestID(ctx, rq.RequestID)
		if !hasCapacity {
			break
		}
	}

	r.lggr.Debugw("consensus plugin query complete", "seqNr", seqNr, "number of request ids", queryBatch.NumberOfRequestIDs())
	return queryBatch.SerialiseQueryBatch()
}

// Removes any requests that have already been completed (successfully/failed/errored) from the batch
// leaving only those requests that are pending consensus.
func (r *reportingPlugin) getPendingRequestsFromPreviousOutcome(previousOutcome []byte, allRequests []*oracle.ConsensusRequest) ([]*oracle.ConsensusRequest, error) {
	var pendingRequests []*oracle.ConsensusRequest
	if len(previousOutcome) == 0 {
		return allRequests, nil
	}

	prevOutcome := &oracletypes.Outcome{}
	err := proto.Unmarshal(previousOutcome, prevOutcome)
	if err != nil {
		r.lggr.Errorw("could not unmarshal previous outcome", "error", err)
		return nil, err
	}

	// Remove any requests from the batch that already have a historical outcome to prevent duplicate outcome generation
	for _, rq := range allRequests {
		_, exists := prevOutcome.HistoricalOutcomes[rq.ID()]
		if exists {
			continue
		}

		pendingRequests = append(pendingRequests, rq)
	}

	return pendingRequests, nil
}
