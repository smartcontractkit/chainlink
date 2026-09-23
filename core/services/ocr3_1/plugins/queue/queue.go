package queue

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = &centralQueuePluginFactory{}

type centralQueuePluginFactory struct {
	lggr  logger.Logger
	meter metric.Meter
}

func NewCentralQueuePluginFactory(lggr logger.Logger, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	//TODO pass CentralTriggerQueue, CapacityReporter, EngineRegistry
	//TODO construct limiters
	return &centralQueuePluginFactory{lggr: lggr}, nil
}

func (c *centralQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	return nil, nil, errors.ErrUnsupported
}

var _ ocr3_1types.ReportingPlugin[[]byte] = &centralQueuePlugin{}

type centralQueuePlugin struct{}

func NewCentralQueuePlugin(lggr logger.Logger) ocr3_1types.ReportingPlugin[[]byte] {
	return &centralQueuePlugin{}
}

func (c *centralQueuePlugin) Query(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return nil, errors.ErrUnsupported
}

func (c *centralQueuePlugin) Observation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, keyValueStateReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	//TODO observe queue via c.CentralTriggerQueue.TakeForObservation()
	return nil, errors.ErrUnsupported
}

func (c *centralQueuePlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation, keyValueStateReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return errors.ErrUnsupported
}

func (c *centralQueuePlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueStateReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return false, errors.ErrUnsupported
}

func (c *centralQueuePlugin) StateTransition(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueStateReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	//TODO write to KV store
	return nil, errors.ErrUnsupported
}

func (c *centralQueuePlugin) Committed(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader) error {
	return nil // no-op
}

func (c *centralQueuePlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, errors.ErrUnsupported
}

func (c *centralQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errors.ErrUnsupported
}

func (c *centralQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errors.ErrUnsupported
}

func (c *centralQueuePlugin) Close() error {
	return nil
}
