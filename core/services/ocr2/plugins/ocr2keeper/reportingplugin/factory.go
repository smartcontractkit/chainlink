package reportingplugin

import (
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// factory implements types.ReportingPluginFactory interface and creates keepers reporting plugin.
type factory struct {
	logger logger.Logger
	hb     httypes.HeadBroadcaster
}

// NewFactory is the constructor of factory
func NewFactory(logger logger.Logger, hb httypes.HeadBroadcaster) types.ReportingPluginFactory {
	return &factory{
		logger: logger,
		hb:     hb,
	}
}

func (f *factory) NewReportingPlugin(rpc types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	p := NewPlugin(f.logger, f.hb)
	pi := types.ReportingPluginInfo{
		Name:          "OCR2Keeper",
		UniqueReports: false,
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       1000,
			MaxObservationLength: 1000,
			MaxReportLength:      1000,
		},
	}
	return p, pi, nil
}
