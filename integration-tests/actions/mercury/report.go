package mercury

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

func GetReportTypes() abi.Arguments {
	return []abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
		{Name: "currentBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockHash", Type: mustNewType("bytes32")},
		{Name: "validFromBlockNum", Type: mustNewType("uint64")},
	}
}

func GetPayloadTypes() abi.Arguments {
	return []abi.Argument{
		{Name: "reportContext", Type: mustNewType("bytes32[3]")},
		{Name: "report", Type: mustNewType("bytes")},
		{Name: "rawRs", Type: mustNewType("bytes32[]")},
		{Name: "rawSs", Type: mustNewType("bytes32[]")},
		{Name: "rawVs", Type: mustNewType("bytes32")},
	}
}

var ReportTypes = GetReportTypes()
var PayloadTypes = GetPayloadTypes()

type ReportWithContext struct {
	Report Report
	Round  uint8
	Epoch  uint32
	Digest []byte
}

func (*ReportWithContext) Pack(reportWithContext *ReportWithContext,
	report [32]byte, rs [][32]byte, ss [][32]byte, vs [32]byte) ([]byte, error) {
	return PayloadTypes.Pack(reportWithContext, report, rs, ss, vs)
}

type Report struct {
	FeedId                [32]byte
	ObservationsTimestamp uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	CurrentBlockNum       uint64
	CurrentBlockHash      [32]uint8
	ValidFromBlockNum     uint64
}

// Use core report types
func (r *Report) Pack() ([]byte, error) {
	return reportcodec.ReportTypes.Pack(r.FeedId, r.ObservationsTimestamp, r.BenchmarkPrice,
		r.Bid, r.Ask, r.CurrentBlockNum, r.CurrentBlockHash, r.ValidFromBlockNum)
}

func DecodeReport(r []byte) (*ReportWithContext, error) {
	payloadElements := map[string]interface{}{}
	if err := PayloadTypes.UnpackIntoMap(payloadElements, r); err != nil {
		return nil, errors.Wrapf(err, "error during payload unpack")
	}

	reportInterface, ok := payloadElements["report"]
	if !ok {
		return nil, errors.Errorf("unpacked payload has no 'report'")
	}
	reportBlob, ok := reportInterface.([]byte)
	if !ok {
		return nil, errors.Errorf("cannot cast report to []byte, type is %T", reportBlob)
	}

	reportCtxInterface, ok := payloadElements["reportContext"]
	if !ok {
		return nil, errors.Errorf("unpacked payload has no 'reportContext'")
	}
	reportCtx, ok := reportCtxInterface.([3][32]byte)
	if !ok {
		return nil, errors.Errorf("cannot cast reportContext to [3][32]byte, type is %T", reportCtx)
	}

	report, err := decodeBlobToReport(reportBlob)
	if err != nil {
		return nil, err
	}

	return &ReportWithContext{
		Report: *report,
		Digest: reportCtx[0][:],
		Round:  reportCtx[1][31],
		Epoch:  binary.BigEndian.Uint32(reportCtx[1][32-5 : 32-1]),
	}, nil
}

func decodeBlobToReport(reportBlob []byte) (*Report, error) {
	r := map[string]interface{}{}
	err := ReportTypes.UnpackIntoMap(r, []byte(reportBlob))
	if err != nil {
		return nil, err
	}

	feedIdInterface, ok := r["feedId"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'feedId'")
	}
	feedID, ok := feedIdInterface.([32]byte)
	if !ok {
		return nil, errors.Errorf("cannot cast feedId to [32]byte, type is %T", feedID)
	}
	log.Trace().Str("FeedID", string(feedID[:])).Msg("Feed ID")

	benchmarkPriceInterface, ok := r["benchmarkPrice"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'benchmarkPrice'")
	}
	benchmarkPrice, ok := benchmarkPriceInterface.(*big.Int)
	if !ok {
		return nil, errors.Errorf("cannot cast 'benchmarkPrice' to *big.Int, type is %T", benchmarkPrice)
	}
	log.Trace().Int64("benchmarkPrice", benchmarkPrice.Int64()).Msg("Benchmark price")

	bidInterface, ok := r["bid"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'bid'")
	}
	bidPrice, ok := bidInterface.(*big.Int)
	if !ok {
		return nil, errors.Errorf("cannot cast 'bid' to *big.Int, type is %T", bidPrice)
	}
	log.Trace().Int64("bid", benchmarkPrice.Int64()).Msg("Bid price")

	askInterface, ok := r["ask"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'ask'")
	}
	askPrice, ok := askInterface.(*big.Int)
	if !ok {
		return nil, errors.Errorf("cannot cast 'bid' to *big.Int, type is %T", askPrice)
	}
	log.Trace().Int64("ask", benchmarkPrice.Int64()).Msg("Ask price")

	currentBlockNumInterface, ok := r["currentBlockNum"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'currentBlockNum'")
	}
	currentBlockNumber, ok := currentBlockNumInterface.(uint64)
	if !ok {
		return nil, errors.Errorf("cannot cast 'currentBlockNum' to uint64, type is %T", currentBlockNumber)
	}
	log.Trace().Uint64("currentBlockNumber", currentBlockNumber).Msg("Observation current block number")

	validFromBlockNumInterface, ok := r["validFromBlockNum"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'validFromBlockNum'")
	}
	validFromBlockNum, ok := validFromBlockNumInterface.(uint64)
	if !ok {
		return nil, errors.Errorf("cannot cast 'validFromBlockNum' to uint64, type is %T", validFromBlockNum)
	}
	log.Trace().Uint64("validFromBlockNum", currentBlockNumber).Msg("Valid from block number")

	currentBlockHashInterface, ok := r["currentBlockHash"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'currentBlockHash'")
	}
	currentBlockHash, ok := currentBlockHashInterface.([32]uint8)
	if !ok {
		return nil, errors.Errorf("cannot cast 'currentBlockHash' to uint64, type is %v", currentBlockHash)
	}
	log.Trace().Any("currentBlockHash", currentBlockHash).Msg("currentBlockHash")

	observationsTimestampInterface, ok := r["observationsTimestamp"]
	if !ok {
		return nil, errors.Errorf("unpacked report has no 'observationsTimestamp'")
	}
	observationsTimestamp, ok := observationsTimestampInterface.(uint32)
	if !ok {
		return nil, errors.Errorf("cannot cast observationsTimestamp to uint32, type is %T", observationsTimestamp)
	}
	log.Trace().Uint32("Timestamp", observationsTimestamp).Msg("Observation timestamp")

	report := &Report{
		FeedId:                feedID,
		ObservationsTimestamp: observationsTimestamp,
		BenchmarkPrice:        benchmarkPrice,
		Bid:                   bidPrice,
		Ask:                   askPrice,
		CurrentBlockNum:       currentBlockNumber,
		CurrentBlockHash:      currentBlockHash,
		ValidFromBlockNum:     validFromBlockNum,
	}

	return report, nil
}

func BuildSampleReport(feedId [32]byte) []byte {
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := utils.NewHash()
	validFromBlockNum := uint64(142)

	b, err := reportcodec.ReportTypes.Pack(feedId, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}
