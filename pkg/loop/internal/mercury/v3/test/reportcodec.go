package v3_test

import (
	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

type StaticReportCodec struct{}

var _ mercury_v3_types.ReportCodec = StaticReportCodec{}

func (StaticReportCodec) BuildReport(fields mercury_v3_types.ReportFields) (ocr2plus_types.Report, error) {
	return Fixtures.Report, nil
}

func (StaticReportCodec) MaxReportLength(n int) (int, error) {
	return Fixtures.MaxReportLength, nil
}

func (StaticReportCodec) ObservationTimestampFromReport(report ocr2plus_types.Report) (uint32, error) {
	return Fixtures.ObservationTimestamp, nil
}
