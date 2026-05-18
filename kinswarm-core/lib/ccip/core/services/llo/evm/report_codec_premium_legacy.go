package evm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

var (
	_ llo.ReportCodec = ReportCodecPremiumLegacy{}
)

type ReportCodecPremiumLegacy struct{ logger.Logger }

func NewReportCodecPremiumLegacy(lggr logger.Logger) llo.ReportCodec {
	return ReportCodecPremiumLegacy{lggr.Named("ReportCodecPremiumLegacy")}
}

type ReportFormatEVMPremiumLegacyOpts struct {
	// BaseUSDFee is the cost on-chain of verifying a report
	BaseUSDFee decimal.Decimal `json:"baseUSDFee"`
	// Expiration window is the length of time in seconds the report is valid
	// for, from the observation timestamp
	ExpirationWindow uint32 `json:"expirationWindow"`
	// FeedID is for compatibility with existing on-chain verifiers
	FeedID common.Hash `json:"feedID"`
	// Multiplier is used to scale the bid, benchmark and ask values in the
	// report. If not specified, or zero is used, a multiplier of 1 is assumed.
	Multiplier *ubig.Big `json:"multiplier"`
}

func (r *ReportFormatEVMPremiumLegacyOpts) Decode(opts []byte) error {
	if len(opts) == 0 {
		// special case if opts are unspecified, just use the zero options rather than erroring
		return nil
	}
	return json.Unmarshal(opts, r)
}

func (r ReportCodecPremiumLegacy) Encode(report llo.Report, cd llotypes.ChannelDefinition) ([]byte, error) {
	if report.Specimen {
		return nil, errors.New("ReportCodecPremiumLegacy does not support encoding specimen reports")
	}
	nativePrice, linkPrice, quote, err := ExtractReportValues(report)
	if err != nil {
		return nil, fmt.Errorf("ReportCodecPremiumLegacy cannot encode; got unusable report; %w", err)
	}

	// NOTE: It seems suboptimal to have to parse the opts on every encode but
	// not sure how to avoid it. Should be negligible performance hit as long
	// as Opts is small.
	opts := ReportFormatEVMPremiumLegacyOpts{}
	if err := (&opts).Decode(cd.Opts); err != nil {
		return nil, fmt.Errorf("failed to decode opts; got: '%s'; %w", cd.Opts, err)
	}
	var multiplier decimal.Decimal
	if opts.Multiplier == nil {
		multiplier = decimal.NewFromInt(1)
	} else if opts.Multiplier.IsZero() {
		return nil, errors.New("multiplier, if specified in channel opts, must be non-zero")
	} else {
		multiplier = decimal.NewFromBigInt(opts.Multiplier.ToInt(), 0)
	}

	codec := reportcodecv3.NewReportCodec(opts.FeedID, r.Logger)

	rf := v3.ReportFields{
		ValidFromTimestamp: report.ValidAfterSeconds + 1,
		Timestamp:          report.ObservationTimestampSeconds,
		NativeFee:          CalculateFee(nativePrice.Decimal(), opts.BaseUSDFee),
		LinkFee:            CalculateFee(linkPrice.Decimal(), opts.BaseUSDFee),
		ExpiresAt:          report.ObservationTimestampSeconds + opts.ExpirationWindow,
		BenchmarkPrice:     quote.Benchmark.Mul(multiplier).BigInt(),
		Bid:                quote.Bid.Mul(multiplier).BigInt(),
		Ask:                quote.Ask.Mul(multiplier).BigInt(),
	}
	return codec.BuildReport(rf)
}

func (r ReportCodecPremiumLegacy) Decode(b []byte) (*reporttypes.Report, error) {
	codec := reportcodecv3.NewReportCodec([32]byte{}, r.Logger)
	return codec.Decode(b)
}

// Pack assembles the report values into a payload for verifying on-chain
func (r ReportCodecPremiumLegacy) Pack(digest types.ConfigDigest, seqNr uint64, report ocr2types.Report, sigs []types.AttributedOnchainSignature) ([]byte, error) {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sigs {
		rr, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			return nil, fmt.Errorf("eventTransmit(ev): error in SplitSignature: %w", err)
		}
		rs = append(rs, rr)
		ss = append(ss, s)
		vs[i] = v
	}
	reportCtx := LegacyReportContext(digest, seqNr)
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	payload, err := mercury.PayloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack failed; %w", err)
	}
	return payload, nil
}

// TODO: Test this
// MERC-3524
func ExtractReportValues(report llo.Report) (nativePrice, linkPrice *llo.Decimal, quote *llo.Quote, err error) {
	if len(report.Values) != 3 {
		return nil, nil, nil, fmt.Errorf("ReportCodecPremiumLegacy requires exactly 3 values (NativePrice, LinkPrice, Quote{Bid, Mid, Ask}); got report.Values: %#v", report.Values)
	}
	var is bool
	nativePrice, is = report.Values[0].(*llo.Decimal)
	if nativePrice == nil {
		// Missing price median will cause a zero fee
		nativePrice = llo.ToDecimal(decimal.Zero)
	} else if !is {
		return nil, nil, nil, fmt.Errorf("ReportCodecPremiumLegacy expects first value to be of type *Decimal; got: %T", report.Values[0])
	}
	linkPrice, is = report.Values[1].(*llo.Decimal)
	if linkPrice == nil {
		// Missing price median will cause a zero fee
		linkPrice = llo.ToDecimal(decimal.Zero)
	} else if !is {
		return nil, nil, nil, fmt.Errorf("ReportCodecPremiumLegacy expects second value to be of type *Decimal; got: %T", report.Values[1])
	}
	quote, is = report.Values[2].(*llo.Quote)
	if !is {
		return nil, nil, nil, fmt.Errorf("ReportCodecPremiumLegacy expects third value to be of type *Quote; got: %T", report.Values[2])
	}
	return nativePrice, linkPrice, quote, nil
}

// TODO: Consider embedding the DON ID here?
// MERC-3524
var LLOExtraHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")

func SeqNrToEpochAndRound(seqNr uint64) (epoch uint32, round uint8) {
	// Simulate 256 rounds/epoch
	epoch = uint32(seqNr / 256) // nolint
	round = uint8(seqNr % 256)  // nolint
	return
}

func LegacyReportContext(cd ocr2types.ConfigDigest, seqNr uint64) ocr2types.ReportContext {
	epoch, round := SeqNrToEpochAndRound(seqNr)
	return ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: cd,
			Epoch:        epoch,
			Round:        round,
		},
		ExtraHash: LLOExtraHash, // ExtraHash is always zero for mercury, we use LLOExtraHash here to differentiate from the legacy plugin
	}
}
