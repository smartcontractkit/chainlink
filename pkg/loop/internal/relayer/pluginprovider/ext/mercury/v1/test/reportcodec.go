package v1_test

import (
	"bytes"
	"context"
	"fmt"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
)

var ReportCodec = staticReportCodec{}

type ReportCodecEvaluator interface {
	mercury_v1_types.ReportCodec
	testtypes.Evaluator[mercury_v1_types.ReportCodec]
}

type staticReportCodec struct{}

var _ mercury_v1_types.ReportCodec = staticReportCodec{}

func (s staticReportCodec) BuildReport(ctx context.Context, fields mercury_v1_types.ReportFields) (ocrtypes.Report, error) {
	return Fixtures.Report, nil
}

// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
// The output of BuildReport must respect this maximum length.
func (s staticReportCodec) MaxReportLength(ctx context.Context, n int) (int, error) {
	return Fixtures.MaxReportLength, nil
}

// CurrentBlockNumFromReport returns the median current block number from a report
func (s staticReportCodec) CurrentBlockNumFromReport(context.Context, ocrtypes.Report) (int64, error) {
	return Fixtures.CurrentBlockNum, nil
}

func (s staticReportCodec) Evaluate(ctx context.Context, other mercury_v1_types.ReportCodec) error {
	gotReport, err := other.BuildReport(ctx, Fixtures.ReportFields)
	if err != nil {
		return fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, Fixtures.Report) {
		return fmt.Errorf("expected Report %x but got %x", Fixtures.Report, gotReport)
	}
	gotMax, err := other.MaxReportLength(ctx, Fixtures.MaxReportLength)
	if err != nil {
		return fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != Fixtures.MaxReportLength {
		return fmt.Errorf("expected MaxReportLength %d but got %d", Fixtures.MaxReportLength, gotMax)
	}
	gotCurrentBlockNum, err := other.CurrentBlockNumFromReport(ctx, gotReport)
	if err != nil {
		return fmt.Errorf("failed to get ObservationTimestampFromReport: %w", err)
	}
	if gotCurrentBlockNum != Fixtures.CurrentBlockNum {
		return fmt.Errorf("expected ObservationTimestampFromReport %d but got %d", Fixtures.CurrentBlockNum, gotCurrentBlockNum)
	}
	return nil
}
