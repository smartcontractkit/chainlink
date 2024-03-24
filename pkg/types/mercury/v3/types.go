package v3

import (
	"context"
	"math/big"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

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
	BuildReport(context.Context, ReportFields) (ocrtypes.Report, error)

	// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(ctx context.Context, n int) (int, error)

	ObservationTimestampFromReport(context.Context, ocrtypes.Report) (uint32, error)
}

// DataSource implementations must be thread-safe. Observe may be called by many
// different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error. Once the
	// context is expires, Observe may still do cheap computations and return a
	// result, but should return as quickly as possible.
	//
	// More details: In the current implementation, the context passed to
	// Observe will time out after MaxDurationObservation. However, Observe
	// should *not* make any assumptions about context timeout behavior. Once
	// the context times out, Observe should prioritize returning as quickly as
	// possible, but may still perform fast computations to return a result
	// rather than error. For example, if Observe medianizes a number of data
	// sources, some of which already returned a result to Observe prior to the
	// context's expiry, Observe might still compute their median, and return it
	// instead of an error.
	//
	// Important: Observe should not perform any potentially time-consuming
	// actions like database access, once the context passed has expired.
	//
	// TODO/WARNING: the type of repts is ocrtypes.ReportTimestamp, but the type of the
	// argument to Observe is types.ReportTimestamp. the underlying struct is the same, but
	// updated here to be self-consistent.
	Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (Observation, error)
}

type Observation struct {
	BenchmarkPrice mercury.ObsResult[*big.Int]
	Bid            mercury.ObsResult[*big.Int]
	Ask            mercury.ObsResult[*big.Int]

	MaxFinalizedTimestamp mercury.ObsResult[int64]

	LinkPrice   mercury.ObsResult[*big.Int]
	NativePrice mercury.ObsResult[*big.Int]
}
