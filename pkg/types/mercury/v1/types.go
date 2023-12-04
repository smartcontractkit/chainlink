package v1

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

type Block struct {
	Num  int64
	Hash string // Hash is stringified to allow use of block as hash key. It is NOT hex and can be cast directly to []byte
	Ts   uint64
}

func NewBlock(num int64, hash []byte, ts uint64) Block {
	return Block{
		Num:  num,
		Hash: string(hash),
		Ts:   ts,
	}
}

// Less returns true if b1 is less than b2 by comparing in order:
//   - smaller block number
//   - smaller timestamp
//   - largest hash
func (b Block) Less(b2 Block) bool {
	if b.Num == b2.Num && b.Ts == b2.Ts {
		// tie-break on hash, all else being equal
		return b.Hash > b2.Hash
	} else if b.Num == b2.Num {
		// if block number is equal and timestamps differ, take the oldest timestamp
		return b.Ts < b2.Ts
	}
	// if block number is different, take the lower block number
	return b.Num < b2.Num
}

func (b Block) String() string {
	return fmt.Sprintf("%d-0x%x-%d", b.Num, []byte(b.Hash), b.Ts)
}

func (b Block) HashBytes() []byte {
	return []byte(b.Hash)
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

type Observation struct {
	BenchmarkPrice mercury.ObsResult[*big.Int]
	Bid            mercury.ObsResult[*big.Int]
	Ask            mercury.ObsResult[*big.Int]

	CurrentBlockNum       mercury.ObsResult[int64]
	CurrentBlockHash      mercury.ObsResult[[]byte]
	CurrentBlockTimestamp mercury.ObsResult[uint64]

	LatestBlocks []Block

	// MaxFinalizedBlockNumber comes from previous report when present and is
	// only observed from mercury server when previous report is nil
	MaxFinalizedBlockNumber mercury.ObsResult[int64]
}
