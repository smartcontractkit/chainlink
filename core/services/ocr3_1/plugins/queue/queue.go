package queue

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = &consensusQueuePluginFactory{}

type consensusQueuePluginFactory struct {
	lggr logger.Logger
	cfg  *PluginConfig
}

func NewConsensusQueuePluginFactory(lggr logger.Logger, cfg *PluginConfig, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	// TODO pass CentralTriggerQueue, CapacityReporter, EngineRegistry
	// TODO construct limiters
	return &consensusQueuePluginFactory{lggr: lggr, cfg: cfg}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr)
	pluginInfo := ocr3_1types.ReportingPluginInfo1{
		Name:   "ConsensusQueueReportingPlugin",
		Limits: c.cfg.Limits,
	}
	return plugin, pluginInfo, nil
}

var _ ocr3_1types.ReportingPlugin[[]byte] = &consensusQueuePlugin{}

type consensusQueuePlugin struct {
	lggr logger.Logger
}

func newConsensusQueuePlugin(lggr logger.Logger) ocr3_1types.ReportingPlugin[[]byte] {
	return &consensusQueuePlugin{lggr: logger.Named(lggr, "ConsensusQueuePlugin")}
}

func (c *consensusQueuePlugin) Query(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return nil, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) Observation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, keyValueStateReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	var _ []triggers.CoordinatedEvent
	// TODO observe queue via c.CentralTriggerQueue.TakeForObservation()
	//TODO serialize observed events
	return nil, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation, keyValueStateReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return errors.ErrUnsupported
}

func (c *consensusQueuePlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueStateReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return false, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) StateTransition(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueStateReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	var key, value []byte
	key = []byte(strconv.Itoa(rand.IntN(256)))
	value = []byte(cryptorand.Text())
	err := keyValueStateReadWriter.Write(key, value)
	if err != nil {
		return nil, fmt.Errorf("failed to write random data to KV store: %w", err)
	}
	//TODO precursor?
	return nil, nil
}

func (c *consensusQueuePlugin) Committed(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader) error {
	return nil // no-op
}

func (c *consensusQueuePlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errors.ErrUnsupported
}

func (c *consensusQueuePlugin) Close() error {
	return nil
}
