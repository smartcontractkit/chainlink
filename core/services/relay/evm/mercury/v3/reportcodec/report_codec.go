package reportcodec

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	reportcodec "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
)

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
		{Name: "validFromTimestamp", Type: mustNewType("uint32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "nativeFee", Type: mustNewType("int192")},
		{Name: "linkFee", Type: mustNewType("int192")},
		{Name: "expiresAt", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
	})
}

type Report struct {
	FeedId                [32]byte
	ObservationsTimestamp uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	ValidFromTimestamp    uint32
	ExpiresAt             uint32
	LinkFee               *big.Int
	NativeFee             *big.Int
}

var _ reportcodec.ReportCodec = &ReportCodec{}

type ReportCodec struct {
	logger logger.Logger
	feedID types.FeedID
}

func NewReportCodec(feedID [32]byte, lggr logger.Logger) *ReportCodec {
	return &ReportCodec{lggr, feedID}
}

func (r *ReportCodec) BuildReport(paos []reportcodec.ParsedAttributedObservation, f int, validFromTimestamp, expiresAt uint32) (ocrtypes.Report, error) {
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

	linkFee, err := relaymercury.GetConsensusLinkFee(mPaos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusLinkFee failed")
	}
	nativeFee, err := relaymercury.GetConsensusNativeFee(mPaos, f)
	if err != nil {
		return nil, errors.Wrap(err, "GetConsensusNativeFee failed")
	}

	reportBytes, err := ReportTypes.Pack(r.feedID, validFromTimestamp, timestamp, nativeFee, linkFee, expiresAt, benchmarkPrice, bid, ask)
	return ocrtypes.Report(reportBytes), errors.Wrap(err, "failed to pack report blob")
}

func (r *ReportCodec) MaxReportLength(n int) (int, error) {
	return maxReportLength, nil
}

func (r *ReportCodec) ObservationTimestampFromReport(report ocrtypes.Report) (uint32, error) {
	reportElems := map[string]interface{}{}
	if err := ReportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return 0, errors.Errorf("error during unpack: %v", err)
	}

	timestampIface, ok := reportElems["observationsTimestamp"]
	if !ok {
		return 0, errors.Errorf("unpacked report has no 'timestamp' field")
	}

	timestamp, ok := timestampIface.(uint32)
	if !ok {
		return 0, errors.Errorf("cannot cast timestamp to uint32, type is %T", timestampIface)
	}

	if timestamp > math.MaxInt32 {
		return 0, errors.Errorf("timestamp overflows max uint32, got: %d", timestamp)
	}

	return timestamp, nil
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
