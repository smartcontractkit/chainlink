package queue

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuepb"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = &consensusQueuePluginFactory{}

type consensusQueuePluginFactory struct {
	lggr logger.Logger
}

func NewConsensusQueuePluginFactory(lggr logger.Logger, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	// TODO pass CentralTriggerQueue, CapacityReporter, EngineRegistry
	// TODO construct limiters
	return &consensusQueuePluginFactory{lggr: lggr}, nil
}

// NewReportingPlugin reads the plugin's OCR limits from config.OffchainConfig, which the config digest covers.
// Returns an error, without starting the plugin, if the offchain config is empty or invalid.
func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	pluginCfg, err := queuepb.Decode(config.OffchainConfig)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("decode reporting plugin config: %w", err)
	}
	pluginLimits, err := pluginCfg.OCRLimits()
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, err
	}
	c.lggr.Infow("Creating consensus-queue reporting plugin", "configDigest", config.ConfigDigest.String(), "limits", pluginLimits)

	plugin := newConsensusQueuePlugin(c.lggr, config)
	pluginInfo := ocr3_1types.ReportingPluginInfo1{
		Name:   "ConsensusQueueReportingPlugin",
		Limits: pluginLimits,
	}
	return plugin, pluginInfo, plugin.Start(ctx)
}

var _ ocr3_1types.ReportingPlugin[[]byte] = &consensusQueuePlugin{}

type consensusQueuePlugin struct {
	services.Service
	eng        *services.Engine
	lggr       logger.Logger
	onchainCfg ocr3types.ReportingPluginConfig
}

func newConsensusQueuePlugin(lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig) *consensusQueuePlugin {
	c := &consensusQueuePlugin{lggr: logger.Named(lggr, "ConsensusQueuePlugin"), onchainCfg: onchainCfg}
	c.Service, c.eng = services.Config{
		Name: "ConsensusQueueReportingPlugin",
	}.NewServiceEngine(lggr)

	return c
}

func (c *consensusQueuePlugin) Query(_ context.Context, _ uint64, _ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return nil, nil
}

func (c *consensusQueuePlugin) Observation(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, kvState ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	return nil, nil
}

func (c *consensusQueuePlugin) ValidateObservation(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher,
) error {
	// TODO validate
	return nil
}

func (c *consensusQueuePlugin) ObservationQuorum(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher,
) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, c.onchainCfg.N, c.onchainCfg.F, aos), nil
}

func (c *consensusQueuePlugin) StateTransition(
	ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	kvState ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	return ocr3_1types.ReportsPlusPrecursor{}, nil
}

func (c *consensusQueuePlugin) Committed(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader) error {
	return nil
}

func (c *consensusQueuePlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, nil
}

func (c *consensusQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (c *consensusQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}
