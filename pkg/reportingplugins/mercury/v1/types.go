package mercury_v1

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
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

// b1 is less than b2 if it is:
// smaller block number
// smaller timestamp
// largest hash
// evaluated in that order
func (b Block) less(b2 Block) bool {
	if b.Num == b2.Num && b.Ts == b2.Ts {
		// tie-break on hash, all else being equal
		return b.Hash > b2.Hash
	} else if b.Num == b2.Num {
		// if block number is equal and timestamps differ, take the oldest timestamp
		return b.Ts < b2.Ts
	} else {
		// if block number is different, take the lower block number
		return b.Num < b2.Num
	}
}

func (b Block) String() string {
	return fmt.Sprintf("%d-0x%x-%d", b.Num, []byte(b.Hash), b.Ts)
}

func (b Block) HashBytes() []byte {
	return []byte(b.Hash)
}

type PAO interface {
	mercury.PAO

	GetBid() (*big.Int, bool)
	GetAsk() (*big.Int, bool)

	// DEPRECATED
	// TODO: Remove this handling after deployment (https://smartcontract-it.atlassian.net/browse/MERC-2272)
	GetCurrentBlockNum() (int64, bool)
	GetCurrentBlockHash() ([]byte, bool)
	GetCurrentBlockTimestamp() (uint64, bool)

	GetLatestBlocks() []Block
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
