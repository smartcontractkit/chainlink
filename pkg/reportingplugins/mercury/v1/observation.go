package mercury_v1 //nolint:revive

import (
	"math/big"

	"github.com/smartcontractkit/libocr/commontypes"
)

var _ PAO = parsedAttributedObservation{}

type parsedAttributedObservation struct {
	Timestamp uint32
	Observer  commontypes.OracleID

	BenchmarkPrice *big.Int
	Bid            *big.Int
	Ask            *big.Int
	// All three prices must be valid, or none are (they all should come from one API query and hold invariant bid <= bm <= ask)
	PricesValid bool

	// DEPRECATED
	// TODO: Remove this handling after deployment (https://smartcontract-it.atlassian.net/browse/MERC-2272)
	CurrentBlockNum       int64 // inclusive; current block
	CurrentBlockHash      []byte
	CurrentBlockTimestamp uint64
	// All three block observations must be valid, or none are (they all come from the same block)
	CurrentBlockValid bool

	LatestBlocks []Block

	// MaxFinalizedBlockNumber comes from previous report when present and is
	// only observed from mercury server when previous report is nil
	//
	// MaxFinalizedBlockNumber will be -1 if there is none
	MaxFinalizedBlockNumber      int64
	MaxFinalizedBlockNumberValid bool
}

func NewParsedAttributedObservation(ts uint32, observer commontypes.OracleID, bp *big.Int, bid *big.Int, ask *big.Int, pricesValid bool,
	currentBlockNum int64, currentBlockHash []byte, currentBlockTimestamp uint64, currentBlockValid bool,
	maxFinalizedBlockNumber int64, maxFinalizedBlockNumberValid bool) PAO {
	return parsedAttributedObservation{
		Timestamp: ts,
		Observer:  observer,

		BenchmarkPrice: bp,
		Bid:            bid,
		Ask:            ask,
		PricesValid:    pricesValid,

		CurrentBlockNum:       currentBlockNum,
		CurrentBlockHash:      currentBlockHash,
		CurrentBlockTimestamp: currentBlockTimestamp,
		CurrentBlockValid:     currentBlockValid,

		MaxFinalizedBlockNumber:      maxFinalizedBlockNumber,
		MaxFinalizedBlockNumberValid: maxFinalizedBlockNumberValid,
	}
}

func (pao parsedAttributedObservation) GetTimestamp() uint32 {
	return pao.Timestamp
}

func (pao parsedAttributedObservation) GetObserver() commontypes.OracleID {
	return pao.Observer
}

func (pao parsedAttributedObservation) GetBenchmarkPrice() (*big.Int, bool) {
	return pao.BenchmarkPrice, pao.PricesValid
}

func (pao parsedAttributedObservation) GetBid() (*big.Int, bool) {
	return pao.Bid, pao.PricesValid
}

func (pao parsedAttributedObservation) GetAsk() (*big.Int, bool) {
	return pao.Ask, pao.PricesValid
}

func (pao parsedAttributedObservation) GetCurrentBlockNum() (int64, bool) {
	return pao.CurrentBlockNum, pao.CurrentBlockValid
}

func (pao parsedAttributedObservation) GetCurrentBlockHash() ([]byte, bool) {
	return pao.CurrentBlockHash, pao.CurrentBlockValid
}

func (pao parsedAttributedObservation) GetCurrentBlockTimestamp() (uint64, bool) {
	return pao.CurrentBlockTimestamp, pao.CurrentBlockValid
}

func (pao parsedAttributedObservation) GetLatestBlocks() []Block {
	return pao.LatestBlocks
}

func (pao parsedAttributedObservation) GetMaxFinalizedBlockNumber() (int64, bool) {
	return pao.MaxFinalizedBlockNumber, pao.MaxFinalizedBlockNumberValid
}

func (pao parsedAttributedObservation) GetMaxFinalizedTimestamp() (uint32, bool) {
	panic("current observation doesn't contain the field")
}

func (pao parsedAttributedObservation) GetLinkFee() (*big.Int, bool) {
	panic("current observation doesn't contain the field")
}

func (pao parsedAttributedObservation) GetNativeFee() (*big.Int, bool) {
	panic("current observation doesn't contain the field")
}
