package consensus

import (
	"context"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type config struct {
	RequestTimeout *time.Duration
	BatchSize      int

	clock clockwork.Clock
}

type factoryService struct {
	store *store
	*capability
	batchSize int

	services.StateMachine
}

const (
	defaultRequestExpiry time.Duration = 1 * time.Hour
	defaultBatchSize                   = 1000
)

func newFactoryService(config *config) (*factoryService, error) {
	if config.RequestTimeout == nil {
		dre := defaultRequestExpiry
		config.RequestTimeout = &dre
	}

	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}

	if config.clock == nil {
		config.clock = clockwork.NewRealClock()
	}

	s := newStore(*config.RequestTimeout, config.clock)
	cp := newCapability(s, config.clock)
	return &factoryService{
		capability: cp,
		store:      s,
		batchSize:  config.BatchSize,
	}, nil
}

func (o *factoryService) NewReportingPlugin(ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	//return newReportingPlugin(o.store, o.capability, o.batchSize), ocr3types.ReportingPluginInfo{Name: "OCR3 Capability Plugin"}, nil
	return nil, ocr3types.ReportingPluginInfo{}, nil
}

func (o *factoryService) Start(ctx context.Context) error {
	return o.StartOnce("plugin factory service", func() error {
		return o.capability.Start(ctx)
	})
}

func (o *factoryService) Close() error {
	return o.StopOnce("plugin factory service", func() error {
		return o.capability.Close()
	})
}

func (o *factoryService) Name() string { return "ocr3PluginFactoryService" }

func (o *factoryService) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}
