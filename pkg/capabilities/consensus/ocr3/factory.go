package ocr3

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type factory struct {
	store *store
	*capability
	batchSize int
	lggr      logger.Logger
}

const (
	defaultMaxPhaseOutputBytes = 100000
	defaultMaxReportCount      = 20
)

func newFactory(s *store, batchSize int, lggr logger.Logger) (*factory, error) {
	return &factory{
		store:     s,
		batchSize: batchSize,
		lggr:      lggr,
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
