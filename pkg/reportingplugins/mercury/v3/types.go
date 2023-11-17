package mercury_v3

import (
	"math/big"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury"
)

type PAO interface {
	mercury.PAO
	GetBid() (*big.Int, bool)
	GetAsk() (*big.Int, bool)
	GetMaxFinalizedTimestamp() (int64, bool)
	GetLinkFee() (*big.Int, bool)
	GetNativeFee() (*big.Int, bool)
}

type ReportFields struct {
	ValidFromTimestamp uint32
	Timestamp          uint32
	NativeFee          *big.Int
	LinkFee            *big.Int
	ExpiresAt          uint32
	BenchmarkPrice     *big.Int
	Bid                *big.Int
	Ask                *big.Int
}

// ReportCodec All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// BuildReport Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport(ReportFields) (ocrtypes.Report, error)

	// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) (int, error)

	ObservationTimestampFromReport(ocrtypes.Report) (uint32, error)
}
