package shim

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// LimitCheckOCR3ReportingPlugin wraps another plugin and checks that its outputs respect
// limits. We use it to surface violations to authors of plugins as early as
// possible.
//
// It does not check inputs since those are checked by the SerializingEndpoint.
type LimitCheckOCR3ReportingPlugin[RI any] struct {
	Plugin ocr3types.ReportingPlugin[RI]
	Limits ocr3types.ReportingPluginLimits
}

var _ ocr3types.ReportingPlugin[struct{}] = LimitCheckOCR3ReportingPlugin[struct{}]{}

func (rp LimitCheckOCR3ReportingPlugin[RI]) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	query, err := rp.Plugin.Query(ctx, outctx)
	if err != nil {
		return nil, err
	}
	if !(len(query) <= rp.Limits.MaxQueryLength) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize query (%v vs %v)", len(query), rp.Limits.MaxQueryLength)
	}
	return query, nil
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	return rp.Plugin.ObservationQuorum(outctx, query)
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	observation, err := rp.Plugin.Observation(ctx, outctx, query)
	if err != nil {
		return nil, err
	}
	if !(len(observation) <= rp.Limits.MaxObservationLength) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize observation (%v vs %v)", len(observation), rp.Limits.MaxObservationLength)
	}
	return observation, nil
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) ValidateObservation(outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	return rp.Plugin.ValidateObservation(outctx, query, ao)
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	outcome, err := rp.Plugin.Outcome(outctx, query, aos)
	if err != nil {
		return nil, err
	}
	if !(len(outcome) <= rp.Limits.MaxOutcomeLength) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize outcome (%v vs %v)", len(outcome), rp.Limits.MaxOutcomeLength)
	}
	return outcome, nil
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[RI], error) {
	reports, err := rp.Plugin.Reports(seqNr, outcome)
	if err != nil {
		return nil, err
	}
	if !(len(reports) <= rp.Limits.MaxReportCount) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned too many reports (%v vs %v)", len(reports), rp.Limits.MaxReportCount)
	}
	for i, reportWithInfo := range reports {
		if !(len(reportWithInfo.Report) <= rp.Limits.MaxReportLength) {
			return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize report at index %v (%v vs %v)", i, len(reportWithInfo.Report), rp.Limits.MaxReportLength)
		}
	}
	return reports, nil
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, report ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldAcceptAttestedReport(ctx, seqNr, report)
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, report ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldTransmitAcceptedReport(ctx, seqNr, report)
}

func (rp LimitCheckOCR3ReportingPlugin[RI]) Close() error {
	return rp.Plugin.Close()
}
