package ring

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

const (
	defaultMaxPhaseOutputBytes = 1000000 // 1 MB
	defaultMaxReportCount      = 1
	defaultBatchSize           = 100
)

var _ core.OCR3ReportingPluginFactory = &Factory{}

type Factory struct {
	ringStore              *Store
	shardOrchestratorStore *shardorchestrator.Store
	arbiterScaler          ringpb.ArbiterScalerClient
	config                 *ConsensusConfig
	lggr                   logger.Logger

	services.StateMachine
}

func NewFactory(s *Store, shardOrchestratorStore *shardorchestrator.Store, arbiterScaler ringpb.ArbiterScalerClient, lggr logger.Logger, cfg *ConsensusConfig) (*Factory, error) {
	if arbiterScaler == nil {
		return nil, errors.New("arbiterScaler is required")
	}
	if cfg == nil {
		cfg = &ConsensusConfig{
			BatchSize: defaultBatchSize,
		}
	}
	return &Factory{
		ringStore:              s,
		shardOrchestratorStore: shardOrchestratorStore,
		arbiterScaler:          arbiterScaler,
		config:                 cfg,
		lggr:                   logger.Named(lggr, "RingPluginFactory"),
	}, nil
}

func (o *Factory) NewReportingPlugin(_ context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	plugin, err := NewPlugin(o.ringStore, o.arbiterScaler, config, o.lggr, o.config)
	pluginInfo := ocr3types.ReportingPluginInfo{
		Name: "RingPlugin",
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       defaultMaxPhaseOutputBytes,
			MaxObservationLength: defaultMaxPhaseOutputBytes,
			MaxOutcomeLength:     defaultMaxPhaseOutputBytes,
			MaxReportLength:      defaultMaxPhaseOutputBytes,
			MaxReportCount:       defaultMaxReportCount,
		},
	}
	return plugin, pluginInfo, err
}

func (o *Factory) Start(ctx context.Context) error {
	return o.StartOnce("RingPlugin", func() error {
		return nil
	})
}

func (o *Factory) Close() error {
	return o.StopOnce("RingPlugin", func() error {
		return nil
	})
}

func (o *Factory) Name() string { return o.lggr.Name() }

func (o *Factory) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}
