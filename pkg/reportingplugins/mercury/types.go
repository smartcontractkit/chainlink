package mercury

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type PAO interface {
	// These fields are common to all observations
	GetTimestamp() uint32
	GetObserver() commontypes.OracleID
	GetBenchmarkPrice() (*big.Int, bool)
}

type ObsResult[T any] struct {
	Val T
	Err error
}

type OnchainConfigCodec interface {
	Encode(OnchainConfig) ([]byte, error)
	Decode([]byte) (OnchainConfig, error)
}

type MercuryServerFetcher interface { //nolint:revive
	// FetchInitialMaxFinalizedBlockNumber should fetch the initial max finalized block number
	FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error)
	LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error)
	LatestTimestamp(context.Context) (int64, error)
}

type Transmitter interface {
	MercuryServerFetcher
	// NOTE: Mercury doesn't actually transmit on-chain, so there is no
	// "contract" involved with the transmitter.
	// - Transmit should be implemented and send to Mercury server
	// - LatestConfigDigestAndEpoch is a stub method, does not need to do anything
	// - FromAccount() should return CSA public key
	ocrtypes.ContractTransmitter
}

type ChainReader interface {
	// LatestHeads returns an ordered list of the latest specified number of heads
	LatestHeads(context.Context, int) ([]Head, error)
}

type Head struct {
	Number    uint64
	Hash      []byte
	Timestamp uint64
}
