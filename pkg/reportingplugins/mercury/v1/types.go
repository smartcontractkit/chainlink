package mercury_v1

import (
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
)

type PAO interface {
	mercury.PAO

	GetBid() (*big.Int, bool)
	GetAsk() (*big.Int, bool)
	GetCurrentBlockNum() (int64, bool)
	GetCurrentBlockHash() ([]byte, bool)
	GetCurrentBlockTimestamp() (uint64, bool)
	GetMaxFinalizedBlockNumber() (int64, bool)
}

type ReportFields struct {
	Timestamp             uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	CurrentBlockNum       int64
	CurrentBlockHash      []byte
	ValidFromBlockNum     int64
	CurrentBlockTimestamp uint64
}

// ReportCodec All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// BuildReport Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport(fields ReportFields) (ocrtypes.Report, error)

	// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) (int, error)

	// CurrentBlockNumFromReport returns the median current block number from a report
	CurrentBlockNumFromReport(types.Report) (int64, error)
}
