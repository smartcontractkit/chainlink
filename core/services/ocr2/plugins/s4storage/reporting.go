package s4storage

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type S4ReportingPluginFactory struct {
	Logger commontypes.Logger
	JobID  uuid.UUID
}

var _ types.ReportingPluginFactory = (*S4ReportingPluginFactory)(nil)

type s4Reporting struct {
	logger        commontypes.Logger
	jobID         uuid.UUID
	genericConfig *types.ReportingPluginConfig
}

var _ types.ReportingPlugin = &s4Reporting{}

// NewReportingPlugin complies with ReportingPluginFactory
func (f S4ReportingPluginFactory) NewReportingPlugin(rpConfig types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	info := types.ReportingPluginInfo{
		Name:          "s4Reporting",
		UniqueReports: true, // Enforces (N+F+1)/2 signatures. Must match setting in OCR2Base.sol.
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       100000,
			MaxObservationLength: 100000,
			MaxReportLength:      100000,
		},
	}
	plugin := s4Reporting{
		logger:        f.Logger,
		jobID:         f.JobID,
		genericConfig: &rpConfig,
	}
	return &plugin, info, nil
}

// Query() complies with ReportingPlugin
func (r *s4Reporting) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	return nil, nil
}

// Observation() complies with ReportingPlugin
func (r *s4Reporting) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	return []byte{}, nil
}

// Report() complies with ReportingPlugin
func (r *s4Reporting) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	return true, []byte{}, nil
}

// ShouldAcceptFinalizedReport() complies with ReportingPlugin
func (r *s4Reporting) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	return false, nil
}

// ShouldTransmitAcceptedReport() complies with ReportingPlugin
func (r *s4Reporting) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	return false, nil
}

// Close() complies with ReportingPlugin
func (r *s4Reporting) Close() error {
	return nil
}
