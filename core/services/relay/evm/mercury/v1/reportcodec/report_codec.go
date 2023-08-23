package reportcodec

import (
	"math"

	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	reportcodec "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/types"
)

// NOTE:
// This report codec is based on the original median evmreportcodec
// here:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/evmreportcodec/reportcodec.go
var ReportTypes = reporttypes.GetSchema()
var maxReportLength = 32 * len(ReportTypes) // each arg is 256 bit EVM word

var _ reportcodec.ReportCodec = &ReportCodec{}

type ReportCodec struct {
	logger logger.Logger
	feedID utils.FeedID
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
func (r *ReportCodec) Decode(report ocrtypes.Report) (*reporttypes.Report, error) {
	return reporttypes.Decode(report)
}
