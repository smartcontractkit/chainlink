package evm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/shopspring/decimal"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// Extracts nanosecond timestamps as uint32 number of seconds
func ExtractTimestamps(report llo.Report) (validAfterSeconds, observationTimestampSeconds uint32, err error) {
	vas := report.ValidAfterNanoseconds / 1e9
	ots := report.ObservationTimestampNanoseconds / 1e9
	if vas > math.MaxUint32 {
		err = fmt.Errorf("validAfterSeconds too large: %d", vas)
		return
	}
	if ots > math.MaxUint32 {
		err = fmt.Errorf("observationTimestampSeconds too large: %d", ots)
		return
	}
	return uint32(vas), uint32(ots), nil
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

// EncodePadded uses standard ABI encoding to encode the stream value, padding
// result to 32 bytes
func (a ABIEncoder) EncodePadded(ctx context.Context, sv llo.StreamValue) ([]byte, error) {
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

// EncodePacked uses packed ABI encoding to encode the stream value (no padding)
func (a ABIEncoder) EncodePacked(ctx context.Context, sv llo.StreamValue) ([]byte, error) {
	var v *big.Int
	switch sv := sv.(type) {
	case *llo.Decimal:
		if sv == nil {
			return nil, fmt.Errorf("expected non-nil *Decimal; got: %v", sv)
		}
		v = a.applyMultiplier(sv.Decimal())
	default:
		return nil, errors.New("EncodePacked only currently supports StreamValue type of *llo.Decimal")
	}
	return EncodePackedBigInt(v, a.Type)
}

// regex to match Solidity integer types (uint/int) and extract bit width
var typeRegex = regexp.MustCompile(`^(u?int)(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)$`)

// EncodePackedBigInt converts a *big.Int to packed EVM bytes according to the Solidity type.
// For unsigned types ("uintN"), the value must be non-negative and fit in N bits.
// For signed types ("intN"), the function returns the two's complement representation in N bits.
func EncodePackedBigInt(value *big.Int, typeStr string) ([]byte, error) {
	// Validate and extract type and bit width
	matches := typeRegex.FindStringSubmatch(typeStr)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid Solidity type: %s", typeStr)
	}
	typePrefix := matches[1] // "uint" or "int"
	bitWidth, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse bit width: %w", err)
	}
	byteLen := bitWidth / 8

	// For unsigned integers, value must be non-negative and within range.
	if typePrefix == "uint" {
		if value.Sign() < 0 {
			return nil, fmt.Errorf("negative value provided for unsigned type %s", typeStr)
		}
		// Check that value fits in bitWidth bits.
		max := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth)) // max = 2^bitWidth
		if value.Cmp(max) >= 0 {
			return nil, fmt.Errorf("value %s out of range for type %s", value.String(), typeStr)
		}
		// Create a byte slice of fixed length and fill it with the big-endian bytes.
		result := make([]byte, byteLen)
		value.FillBytes(result)
		return result, nil
	}

	// For signed integers, compute the two's complement representation.
	// Solidity uses two's complement in fixed N-bit representation.
	twoPow := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth)) // twoPow = 2^bitWidth
	// Compute the modulo to obtain the two's complement representation.
	modValue := new(big.Int).Mod(value, twoPow)
	result := make([]byte, byteLen)
	modValue.FillBytes(result)
	return result, nil
}
