package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_ llo.ReportCodec = ReportCodecEVMABIEncodeUnpacked{}

	zero = big.NewInt(0)
)

type ReportCodecEVMABIEncodeUnpacked struct {
	logger.Logger
	donID uint32
}

func NewReportCodecEVMABIEncodeUnpacked(lggr logger.Logger, donID uint32) ReportCodecEVMABIEncodeUnpacked {
	return ReportCodecEVMABIEncodeUnpacked{logger.Sugared(lggr).Named("ReportCodecEVMABIEncodeUnpacked"), donID}
}

// Opts format remains unchanged
type ReportFormatEVMABIEncodeOpts struct {
	// BaseUSDFee is the cost on-chain of verifying a report
	BaseUSDFee decimal.Decimal `json:"baseUSDFee"`
	// Expiration window is the length of time in seconds the report is valid
	// for, from the observation timestamp
	ExpirationWindow uint32 `json:"expirationWindow"`
	// FeedID is for compatibility with existing on-chain verifiers
	FeedID common.Hash `json:"feedID"`
	// ABI defines the encoding of the payload. Each element maps to exactly
	// one stream (although sub-arrays may be specified for streams that
	// produce a composite data type).
	//
	// EXAMPLE
	//
	// [{"streamID":123,"multiplier":"10000","type":"uint192"}, ...]
	//
	// See definition of ABIEncoder struct for more details.
	//
	// The total number of streams must be 2+n, where n is the number of
	// top-level elements in this ABI array (stream 0 is always the native
	// token price and stream 1 is the link token price).
	ABI []ABIEncoder `json:"abi"`
}

func (r *ReportFormatEVMABIEncodeOpts) Decode(opts []byte) error {
	return json.Unmarshal(opts, r)
}

func (r *ReportFormatEVMABIEncodeOpts) Encode() ([]byte, error) {
	return json.Marshal(r)
}

type BaseReportFields struct {
	FeedID             common.Hash
	ValidFromTimestamp uint32
	Timestamp          uint32
	NativeFee          *big.Int
	LinkFee            *big.Int
	ExpiresAt          uint32
}

func (r ReportCodecEVMABIEncodeUnpacked) Encode(ctx context.Context, report llo.Report, cd llotypes.ChannelDefinition) ([]byte, error) {
	if report.Specimen {
		return nil, errors.New("ReportCodecEVMABIEncodeUnpacked does not support encoding specimen reports")
	}
	if len(report.Values) < 2 {
		return nil, fmt.Errorf("ReportCodecEVMABIEncodeUnpacked requires at least 2 values (NativePrice, LinkPrice, ...); got report.Values: %v", report.Values)
	}
	nativePrice, err := extractPrice(report.Values[0])
	if err != nil {
		return nil, fmt.Errorf("ReportCodecEVMABIEncodeUnpacked failed to extract native price: %w", err)
	}
	linkPrice, err := extractPrice(report.Values[1])
	if err != nil {
		return nil, fmt.Errorf("ReportCodecEVMABIEncodeUnpacked failed to extract link price: %w", err)
	}

	// NOTE: It seems suboptimal to have to parse the opts on every encode but
	// not sure how to avoid it. Should be negligible performance hit as long
	// as Opts is small.
	opts := ReportFormatEVMABIEncodeOpts{}
	if err = (&opts).Decode(cd.Opts); err != nil {
		return nil, fmt.Errorf("failed to decode opts; got: '%s'; %w", cd.Opts, err)
	}

	rf := BaseReportFields{
		FeedID:             opts.FeedID,
		ValidFromTimestamp: report.ValidAfterSeconds + 1,
		Timestamp:          report.ObservationTimestampSeconds,
		NativeFee:          CalculateFee(nativePrice, opts.BaseUSDFee),
		LinkFee:            CalculateFee(linkPrice, opts.BaseUSDFee),
		ExpiresAt:          report.ObservationTimestampSeconds + opts.ExpirationWindow,
	}

	header, err := r.buildHeader(ctx, rf)
	if err != nil {
		return nil, fmt.Errorf("failed to build base report; %w", err)
	}

	payload, err := r.buildPayload(ctx, opts.ABI, report.Values[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to build payload; %w", err)
	}

	return append(header, payload...), nil
}

// BaseSchema represents the fixed base schema that remains unchanged for all
// EVMABIEncodeUnpacked reports.
//
// An arbitrary payload will be appended to this.
var BaseSchema = getBaseSchema()

func getBaseSchema() abi.Arguments {
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
		{Name: "nativeFee", Type: mustNewType("uint192")},
		{Name: "linkFee", Type: mustNewType("uint192")},
		{Name: "expiresAt", Type: mustNewType("uint32")},
	})
}

func (r ReportCodecEVMABIEncodeUnpacked) buildHeader(ctx context.Context, rf BaseReportFields) ([]byte, error) {
	var merr error
	if rf.LinkFee == nil {
		merr = errors.Join(merr, errors.New("linkFee may not be nil"))
	} else if rf.LinkFee.Cmp(zero) < 0 {
		merr = errors.Join(merr, fmt.Errorf("linkFee may not be negative (got: %s)", rf.LinkFee))
	}
	if rf.NativeFee == nil {
		merr = errors.Join(merr, errors.New("nativeFee may not be nil"))
	} else if rf.NativeFee.Cmp(zero) < 0 {
		merr = errors.Join(merr, fmt.Errorf("nativeFee may not be negative (got: %s)", rf.NativeFee))
	}
	if merr != nil {
		return nil, merr
	}
	b, err := BaseSchema.Pack(rf.FeedID, rf.ValidFromTimestamp, rf.Timestamp, rf.NativeFee, rf.LinkFee, rf.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to pack base report blob; %w", err)
	}
	return b, nil
}

func (r ReportCodecEVMABIEncodeUnpacked) buildPayload(ctx context.Context, encoders []ABIEncoder, values []llo.StreamValue) (payload []byte, merr error) {
	if len(encoders) != len(values) {
		return nil, fmt.Errorf("ABI and values length mismatch; ABI: %d, Values: %d", len(encoders), len(values))
	}

	for i, encoder := range encoders {
		b, err := encoder.Encode(ctx, values[i])
		if err != nil {
			var vStr []byte
			if values[i] == nil {
				vStr = []byte("<nil>")
			} else {
				var marshalErr error
				vStr, marshalErr = values[i].MarshalText()
				if marshalErr != nil {
					vStr = []byte(fmt.Sprintf("%v(failed to marshal: %s)", values[i], marshalErr))
				}
			}
			merr = errors.Join(merr, fmt.Errorf("failed to encode stream value %s at index %d with abi %q; %w", string(vStr), i, encoder.Type, err))
			continue
		}
		payload = append(payload, b...)
	}

	return payload, merr
}

// An ABIEncoder encodes exactly one stream value into a byte slice
type ABIEncoder struct {
	// StreamID is the ID of the stream that this encoder is responsible for.
	// MANDATORY
	StreamID llotypes.StreamID `json:"streamID"`
	// Type is the ABI type of the stream value. E.g. "uint192", "int256", "bool", "string" etc.
	// MANDATORY
	Type string `json:"type"`
	// Multiplier, if provided, will be multiplied with the stream value before
	// encoding.
	// OPTIONAL
	Multiplier *ubig.Big `json:"multiplier"`
}

// getNormalizedMultiplier returns the multiplier as a decimal.Decimal, defaulting
// to 1 if the multiplier is nil.
//
// Negative multipliers are ok and will work as expected, flipping the sign of
// the value.
func (a ABIEncoder) getNormalizedMultiplier() (multiplier decimal.Decimal) {
	if a.Multiplier == nil {
		multiplier = decimal.NewFromInt(1)
	} else {
		multiplier = decimal.NewFromBigInt(a.Multiplier.ToInt(), 0)
	}
	return
}

func (a ABIEncoder) applyMultiplier(d decimal.Decimal) *big.Int {
	return d.Mul(a.getNormalizedMultiplier()).BigInt()
}

func (a ABIEncoder) Encode(ctx context.Context, sv llo.StreamValue) ([]byte, error) {
	var encode interface{}
	switch sv := sv.(type) {
	case *llo.Decimal:
		if sv == nil {
			return nil, fmt.Errorf("expected non-nil *Decimal; got: %v", sv)
		}
		encode = a.applyMultiplier(sv.Decimal())
	default:
		return nil, fmt.Errorf("unhandled type; supported types are: *llo.Decimal; got: %T", sv)
	}
	evmEncoderConfig := fmt.Sprintf(`[{"Name":"streamValue","Type":"%s"}]`, a.Type)

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		"evm": {TypeABI: evmEncoderConfig},
	}}
	c, err := codec.NewCodec(codecConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create codec; %w", err)
	}

	result, err := c.Encode(ctx, map[string]any{"streamValue": encode}, "evm")
	if err != nil {
		return nil, fmt.Errorf("failed to encode stream value %v with ABI type %q; %w", sv, a.Type, err)
	}
	return result, nil
}
