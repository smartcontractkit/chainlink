package vault

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	defaultBatchSize = 20
)

func NewReportingPluginFactory(store *requests.Store[*Request]) *ReportingPluginFactory {
	return &ReportingPluginFactory{
		store:     store,
		batchSize: defaultBatchSize, // TODO fetch from onchain config
	}
}

type ReportingPluginFactory struct {
	store     *requests.Store[*Request]
	batchSize int
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	return &ReportingPlugin{store: r.store, batchSize: r.batchSize}, ocr3_1types.ReportingPluginInfo{}, nil
}

type ReportingPlugin struct {
	store     *requests.Store[*Request]
	batchSize int
}

func (r *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	batch, err := r.store.FirstN(r.batchSize)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
	}

	for _, req := range batch {
	}

	return types.Observation{}, nil
}

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return nil
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return false, nil
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	return ocr3_1types.ReportsPlusPrecursor{}, nil
}

func (r *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader) error {
	// Not currently used by the protocol, so we noop here.
	return nil
}

func (r *ReportingPlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, nil
}

func (r *ReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func (r *ReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func (r *ReportingPlugin) Close() error {
	return nil
}
