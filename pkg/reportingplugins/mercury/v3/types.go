package mercury_v3

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
)

type ParsedAttributedObservation interface {
	mercury.ParsedAttributedObservation
}

func Convert(pao []ParsedAttributedObservation) []mercury.ParsedAttributedObservation {
	var ret []mercury.ParsedAttributedObservation
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}

// ReportCodec All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// BuildReport Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport(paos []ParsedAttributedObservation, f int, validFromTimestamp, expiresAt uint32) (ocrtypes.Report, error)

	// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) (int, error)

	ObservationTimestampFromReport(ocrtypes.Report) (uint32, error)
}
