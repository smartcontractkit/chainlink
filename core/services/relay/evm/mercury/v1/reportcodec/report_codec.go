package reportcodec

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	reportcodec "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
)

// NOTE:
// This report codec is based on the original median evmreportcodec
// here:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/evmreportcodec/reportcodec.go

var ReportTypes = getReportTypes()
var maxReportLength = 32 * len(ReportTypes) // each arg is 256 bit EVM word

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
		{Name: "currentBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockHash", Type: mustNewType("bytes32")},
		{Name: "validFromBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockTimestamp", Type: mustNewType("uint64")},
	})
}

type Report struct {
	FeedId                [32]byte
	ObservationsTimestamp uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	CurrentBlockNum       uint64
	CurrentBlockHash      [32]byte
	ValidFromBlockNum     uint64
	CurrentBlockTimestamp uint64
}

var _ reportcodec.ReportCodec = &ReportCodec{}

type ReportCodec struct {
	logger logger.Logger
	feedID types.FeedID
}

func NewReportCodec(feedID [32]byte, lggr logger.Logger) *ReportCodec {
	return &ReportCodec{lggr, feedID}
}

func (r *ReportCodec) BuildReport(paos []reportcodec.ParsedAttributedObservation, f int, validFromBlockNum int64) (ocrtypes.Report, error) {
	if len(paos) == 0 {
		return nil, errors.Errorf("cannot build report from empty attributed observations")
	}

	mPaos := reportcodec.Convert(paos)

	timestamp := relaymercury.GetConsensusTimestamp(mPaos)
	benchmarkPrice, err := relaymercury.GetConsensusBenchmarkPrice(mPaos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusBenchmarkPrice failed")
	}
	bid, err := relaymercury.GetConsensusBid(mPaos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusBid failed")
	}
	ask, err := relaymercury.GetConsensusAsk(mPaos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusAsk failed")
	}

	currentBlockHash, currentBlockNum, currentBlockTimestamp, err := reportcodec.GetConsensusCurrentBlock(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusCurrentBlock failed")
	}

	if validFromBlockNum > currentBlockNum {
		return nil, errors.Errorf("validFromBlockNum=%d may not be greater than currentBlockNum=%d", validFromBlockNum, currentBlockNum)
	}

	if len(currentBlockHash) != 32 {
		return nil, errors.Errorf("invalid length for currentBlockHash, expected: 32, got: %d", len(currentBlockHash))
	}
	currentBlockHashArray := [32]byte{}
	copy(currentBlockHashArray[:], currentBlockHash)

	reportBytes, err := ReportTypes.Pack(r.feedID, timestamp, benchmarkPrice, bid, ask, uint64(currentBlockNum), currentBlockHashArray, uint64(validFromBlockNum), currentBlockTimestamp)
	return ocrtypes.Report(reportBytes), errors.Wrap(err, "failed to pack report blob")
}

// Maximum length in bytes of Report returned by BuildReport. Used for
// defending against spam attacks.
func (r *ReportCodec) MaxReportLength(n int) (int, error) {
	return maxReportLength, nil
}

func (r *ReportCodec) CurrentBlockNumFromReport(report ocrtypes.Report) (int64, error) {
	reportElems := map[string]interface{}{}
	if err := ReportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return 0, errors.Errorf("error during unpack: %v", err)
	}

	blockNumIface, ok := reportElems["currentBlockNum"]
	if !ok {
		return 0, errors.Errorf("unpacked report has no 'currentBlockNum' field")
	}

	blockNum, ok := blockNumIface.(uint64)
	if !ok {
		return 0, errors.Errorf("cannot cast blockNum to int64, type is %T", blockNumIface)
	}

	if blockNum > math.MaxInt64 {
		return 0, errors.Errorf("blockNum overflows max int64, got: %d", blockNum)
	}

	return int64(blockNum), nil
}

func (r *ReportCodec) ValidFromBlockNumFromReport(report ocrtypes.Report) (int64, error) {
	reportElems := map[string]interface{}{}
	if err := ReportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return 0, errors.Errorf("error during unpack: %v", err)
	}

	blockNumIface, ok := reportElems["validFromBlockNum"]
	if !ok {
		return 0, errors.Errorf("unpacked report has no 'validFromBlockNum' field")
	}

	blockNum, ok := blockNumIface.(uint64)
	if !ok {
		return 0, errors.Errorf("cannot cast blockNum to int64, type is %T", blockNumIface)
	}

	if blockNum > math.MaxInt64 {
		return 0, errors.Errorf("blockNum overflows max int64, got: %d", blockNum)
	}

	return int64(blockNum), nil
}

// Decode is made available to external users (i.e. mercury server)
func (r *ReportCodec) Decode(report ocrtypes.Report) (*Report, error) {
	values, err := ReportTypes.Unpack(report)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report: %w", err)
	}
	decoded := new(Report)
	if err = ReportTypes.Copy(decoded, values); err != nil {
		return nil, fmt.Errorf("failed to copy report values to struct: %w", err)
	}
	return decoded, nil
}
