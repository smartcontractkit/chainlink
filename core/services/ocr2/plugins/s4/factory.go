package s4

import (
	s4_orm "github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

const S4ReportingPluginName = "S4Reporting"

type PluginConfigDecoder func([]byte) (*PluginConfig, *types.ReportingPluginLimits, error)

type S4ReportingPluginFactory struct {
	Logger        commontypes.Logger
	ORM           s4_orm.ORM
	ConfigDecoder PluginConfigDecoder
}

var _ types.ReportingPluginFactory = (*S4ReportingPluginFactory)(nil)

// NewReportingPlugin complies with ReportingPluginFactory
func (f S4ReportingPluginFactory) NewReportingPlugin(rpConfig types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	config, limits, err := f.ConfigDecoder(rpConfig.OffchainConfig)
	if err != nil {
		f.Logger.Error("unable to decode reporting plugin config", commontypes.LogFields{
			"digest": rpConfig.ConfigDigest.String(),
		})
		return nil, types.ReportingPluginInfo{}, err
	}
	info := types.ReportingPluginInfo{
		Name:          S4ReportingPluginName,
		UniqueReports: false,
		Limits:        *limits,
	}
	plugin, err := NewReportingPlugin(f.Logger, config, f.ORM)
	if err != nil {
		f.Logger.Error("unable to create S4 reporting plugin", commontypes.LogFields{})
		return nil, types.ReportingPluginInfo{}, err
	}
	return plugin, info, nil
}
