package v2_test

import (
	"bytes"
	"context"
	"fmt"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
)

var ReportCodec = staticReportCodec{}

type ReportCodecEvaluator interface {
	mercury_v2_types.ReportCodec
	testtypes.Evaluator[mercury_v2_types.ReportCodec]
}

type staticReportCodec struct{}

var _ ReportCodecEvaluator = staticReportCodec{}

type StaticReportCodecValues struct {
	Report               ocrtypes.Report
	MaxReportLength      int
	ObservationTimestamp uint32
}

var StaticReportCodecFixtures = StaticReportCodecValues{
	Report:               ocrtypes.Report([]byte("mercury v2 report")),
	MaxReportLength:      20,
	ObservationTimestamp: 23,
}

func (s staticReportCodec) BuildReport(fields mercury_v2_types.ReportFields) (types.Report, error) {
	return StaticReportCodecFixtures.Report, nil
}

// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
// The output of BuildReport must respect this maximum length.
func (s staticReportCodec) MaxReportLength(n int) (int, error) {
	return StaticReportCodecFixtures.MaxReportLength, nil
}

// CurrentBlockNumFromReport returns the median current block number from a report
func (s staticReportCodec) ObservationTimestampFromReport(types.Report) (uint32, error) {
	return StaticReportCodecFixtures.ObservationTimestamp, nil
}

func (s staticReportCodec) Evaluate(ctx context.Context, other mercury_v2_types.ReportCodec) error {
	gotReport, err := other.BuildReport(Fixtures.ReportFields)
	if err != nil {
		return fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, Fixtures.Report) {
		return fmt.Errorf("expected Report %x but got %x", Fixtures.Report, gotReport)
	}
	gotMax, err := other.MaxReportLength(Fixtures.MaxReportLength)
	if err != nil {
		return fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != Fixtures.MaxReportLength {
		return fmt.Errorf("expected MaxReportLength %d but got %d", Fixtures.MaxReportLength, gotMax)
	}
	gotObservedTimestamp, err := other.ObservationTimestampFromReport(gotReport)
	if err != nil {
		return fmt.Errorf("failed to get ObservationTimestampFromReport: %w", err)
	}
	if gotObservedTimestamp != Fixtures.ObservationTimestamp {
		return fmt.Errorf("expected ObservationTimestampFromReport %d but got %d", Fixtures.ObservationTimestamp, gotObservedTimestamp)
	}
	return nil
}
