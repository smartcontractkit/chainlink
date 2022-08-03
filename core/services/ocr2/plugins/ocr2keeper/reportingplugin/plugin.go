package reportingplugin

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// plugin implements types.ReportingPlugin interface with the keepers-specific logic.
type plugin struct {
	logger logger.Logger
}

// NewPlugin is the constructor of plugin
func NewPlugin(logger logger.Logger) types.ReportingPlugin {
	return &plugin{
		logger: logger,
	}
}

func (p *plugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	p.logger.Info("Query()", nil)
	return []byte("Query()"), nil
}

func (p *plugin) Observation(_ context.Context, _ types.ReportTimestamp, q types.Query) (types.Observation, error) {
	p.logger.Info("Observation()", string(q))
	return []byte("Observation()"), nil
}

func (p *plugin) Report(_ context.Context, _ types.ReportTimestamp, q types.Query, _ []types.AttributedObservation) (bool, types.Report, error) {
	p.logger.Info("Report()", string(q))
	return true, []byte("Report()"), nil
}

func (p *plugin) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldAcceptFinalizedReport()", string(r))
	return true, nil
}

func (p *plugin) ShouldTransmitAcceptedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldTransmitAcceptedReport()", string(r))
	return true, nil
}

func (p *plugin) Close() error {
	p.logger.Info("Close()", nil)
	return nil
}
