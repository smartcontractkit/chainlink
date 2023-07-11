package mercury_v2

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	reportcodec "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var ReportTypes = getReportTypes()

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
		{Name: "validFromTimestamp", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
		{Name: "expiresAt", Type: mustNewType("uint32")},
		{Name: "linkFee", Type: mustNewType("int192")},
		{Name: "nativeFee", Type: mustNewType("int192")},
	})
}

var _ reportcodec.ReportCodec = &ReportCodec{}

type ReportCodec struct {
	logger logger.Logger
	feedID [32]byte
}

func NewReportCodec(feedID [32]byte, lggr logger.Logger) *ReportCodec {
	return &ReportCodec{lggr, feedID}
}

func (r *ReportCodec) BuildReport(paos []relaymercury.ParsedObservation, f int, validFromTimestamp int64) (ocrtypes.Report, error) {
	if len(paos) == 0 {
		return nil, errors.Errorf("cannot build report from empty attributed observations")
	}

	// copy so we can safely sort in place
	paos = append([]relaymercury.ParsedObservation{}, paos...)

	timestamp := relaymercury.GetConsensusTimestamp(paos)
	expiresAt := relaymercury.GetConsensusExpiresAt(paos)

	// todo: add checks for validFromTimestamp

	benchmarkPrice, err := relaymercury.GetConsensusBenchmarkPrice(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusBenchmarkPrice failed")
	}
	bid, err := relaymercury.GetConsensusBid(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusBid failed")
	}
	ask, err := relaymercury.GetConsensusAsk(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusAsk failed")
	}
	linkFee, err := relaymercury.GetConsensusLinkFee(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusLinkFee failed")
	}
	nativeFee, err := relaymercury.GetConsensusNativeFee(paos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusNativeFee failed")
	}

	reportBytes, err := ReportTypes.Pack(r.feedID, timestamp, validFromTimestamp, benchmarkPrice, bid, ask, expiresAt, linkFee, nativeFee)
	return ocrtypes.Report(reportBytes), errors.Wrap(err, "failed to pack report blob")
}

func (r *ReportCodec) MaxReportLength(n int) (int, error) {
	return 8*32 + // feed ID
			32 + // timestamp
			64 + // validFromTimestamp
			192 + // benchmarkPrice
			192 + // bid
			192, // ask
		nil
}

func (r *ReportCodec) CurrentTimestampFromReport(report ocrtypes.Report) (int64, error) {
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

func (r *ReportCodec) ValidFromTimestampFromReport(report ocrtypes.Report) (int64, error) {
	reportElems := map[string]interface{}{}
	if err := ReportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return 0, errors.Errorf("error during unpack: %v", err)
	}

	timestampIface, ok := reportElems["validFromTimestamp"]
	if !ok {
		return 0, errors.Errorf("unpacked report has no 'validFromTimestamp' field")
	}

	timestamp, ok := timestampIface.(uint64)
	if !ok {
		return 0, errors.Errorf("cannot cast blockNum to int64, type is %T", timestampIface)
	}

	if timestamp > math.MaxInt64 {
		return 0, errors.Errorf("timestamp overflows max int64, got: %d", timestamp)
	}

	return int64(timestamp), nil
}
