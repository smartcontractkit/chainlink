package shim

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// LimitCheckReportingPlugin wraps another ReportingPlugin and checks that
// its outputs respect limits. We use it to surface violations to authors of
// ReportingPlugins as early as possible.
//
// It does not check inputs since those are checked by the SerializingEndpoint.
type LimitCheckReportingPlugin struct {
	Plugin types.ReportingPlugin
	Limits types.ReportingPluginLimits
}

var _ types.ReportingPlugin = LimitCheckReportingPlugin{}

func (rp LimitCheckReportingPlugin) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	query, err := rp.Plugin.Query(ctx, ts)
	if err != nil {
		return nil, err
	}
	if !(len(query) <= rp.Limits.MaxQueryLength) {
		return nil, fmt.Errorf("LimitCheckReportingPlugin: underlying ReportingPlugin returned oversize query (%v vs %v)", len(query), rp.Limits.MaxQueryLength)
	}
	return query, nil
}

func (rp LimitCheckReportingPlugin) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	observation, err := rp.Plugin.Observation(ctx, ts, query)
	if err != nil {
		return nil, err
	}
	if !(len(observation) <= rp.Limits.MaxObservationLength) {
		return nil, fmt.Errorf("LimitCheckReportingPlugin: underlying ReportingPlugin returned oversize observation (%v vs %v)", len(observation), rp.Limits.MaxObservationLength)
	}
	return observation, nil
}

func (rp LimitCheckReportingPlugin) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	shouldReport, report, err := rp.Plugin.Report(ctx, ts, query, aos)
	if err != nil {
		return false, nil, err
	}
	if !(len(report) <= rp.Limits.MaxReportLength) {
		return false, nil, fmt.Errorf("LimitCheckReportingPlugin: underlying ReportingPlugin returned oversize report (%v vs %v)", len(report), rp.Limits.MaxReportLength)
	}
	return shouldReport, report, nil
}

func (rp LimitCheckReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	return rp.Plugin.ShouldAcceptFinalizedReport(ctx, ts, report)
}

func (rp LimitCheckReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	return rp.Plugin.ShouldTransmitAcceptedReport(ctx, ts, report)
}

func (rp LimitCheckReportingPlugin) Close() error {
	return rp.Plugin.Close()
}
