package ocr3

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/requests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type factory struct {
	store      *requests.Store
	capability *capability
	batchSize  int
	lggr       logger.Logger

	services.StateMachine
}

const (
	defaultMaxPhaseOutputBytes = 100000
	defaultMaxReportCount      = 20
)

func newFactory(s *requests.Store, c *capability, batchSize int, lggr logger.Logger) (*factory, error) {
	return &factory{
		store:      s,
		capability: c,
		batchSize:  batchSize,
		lggr:       logger.Named(lggr, "OCR3ReportingPluginFactory"),
	}, nil
}

func (o *factory) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	rp, err := newReportingPlugin(o.store, o.capability, o.batchSize, config, o.lggr)
	info := ocr3types.ReportingPluginInfo{
		Name: "OCR3 Capability Plugin",
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       defaultMaxPhaseOutputBytes,
			MaxObservationLength: defaultMaxPhaseOutputBytes,
			MaxOutcomeLength:     defaultMaxPhaseOutputBytes,
			MaxReportLength:      defaultMaxPhaseOutputBytes,
			MaxReportCount:       defaultMaxReportCount,
		},
	}
	return rp, info, err
}

func (o *factory) Start(ctx context.Context) error {
	return o.StartOnce("OCR3ReportingPlugin", func() error {
		return nil
	})
}

func (o *factory) Close() error {
	return o.StopOnce("OCR3ReportingPlugin", func() error {
		return nil
	})
}

func (o *factory) Name() string { return o.lggr.Name() }

func (o *factory) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}
