package mercury

import "github.com/smartcontractkit/chainlink-common/pkg/values"

type MercuryReportSet struct {
	// feedID -> report
	Reports map[string]MercuryReport
}

type MercuryReport struct {
	Info       MercuryReportInfo // minimal data extracted from the report for convenience
	FullReport []byte            // full report, acceptable by the verifier contract
}

type MercuryReportInfo struct {
	Timestamp uint32
	Price     float64
}

type MercuryCodec interface {
	// validate each report and convert to MercuryReportSet struct
	Unwrap(raw values.Value) (MercuryReportSet, error)

	// validate each report and convert to Value
	Wrap(reportSet MercuryReportSet) (values.Value, error)
}

// TODO implement a codec
