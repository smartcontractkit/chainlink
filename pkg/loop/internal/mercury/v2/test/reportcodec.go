package v2_test

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
)

type StaticReportCodec struct{}

var _ mercury_v2_types.ReportCodec = StaticReportCodec{}

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

func (s StaticReportCodec) BuildReport(fields mercury_v2_types.ReportFields) (ocrtypes.Report, error) {
	return StaticReportCodecFixtures.Report, nil
}

// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
// The output of BuildReport must respect this maximum length.
func (s StaticReportCodec) MaxReportLength(n int) (int, error) {
	return StaticReportCodecFixtures.MaxReportLength, nil
}

// CurrentBlockNumFromReport returns the median current block number from a report
func (s StaticReportCodec) ObservationTimestampFromReport(ocrtypes.Report) (uint32, error) {
	return StaticReportCodecFixtures.ObservationTimestamp, nil
}
