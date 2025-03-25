package evm

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
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
	encoders []singleABIEncoder
}

func (a ABIEncoder) MarshalJSON() ([]byte, error) {
	if len(a.encoders) == 1 {
		return json.Marshal(a.encoders[0])
	}
	return json.Marshal(a.encoders)
}

func (a *ABIEncoder) UnmarshalJSON(data []byte) error {
	// Try to unmarshal into a slice
	var encoders []singleABIEncoder
	if err := json.Unmarshal(data, &encoders); err == nil {
		a.encoders = encoders
		return nil
	}

	// Try to unmarshal into a single encoder
	var encoder singleABIEncoder
	if err := json.Unmarshal(data, &encoder); err != nil {
		return fmt.Errorf("failed to unmarshal ABIEncoder: %w", err)
	}
	a.encoders = []singleABIEncoder{encoder}
	return nil
}

func (a ABIEncoder) EncodePacked(sv llo.StreamValue) ([]byte, error) {
	switch v := sv.(type) {
	case *llo.Decimal:
		if len(a.encoders) != 1 {
			return nil, fmt.Errorf("expected exactly one encoder; got: %d", len(a.encoders))
		}
		b, err := a.encoders[0].encodePacked(sv)
		if err != nil {
			return nil, fmt.Errorf("failed to encode stream value; %w", err)
		}
		return b, nil
	case *llo.TimestampedStreamValue:
		if len(a.encoders) != 2 {
			return nil, fmt.Errorf("expected exactly two encoders for *llo.TimestampedStreamValue; got: %d", len(a.encoders))
		}
		// encode as two packed values
		// <type0> observed_at_nanoseconds
		// <type1> stream_value
		encodedTimestamp, err := a.encoders[0].encodeUint64Packed(v.ObservedAtNanoseconds)
		if err != nil {
			return nil, fmt.Errorf("failed to encode timestamped stream value; %w", err)
		}
		encodedDecimal, err := a.encoders[1].encodePacked(v.StreamValue)
		if err != nil {
			return nil, fmt.Errorf("failed to encode nested stream value; %w", err)
		}
		return append(encodedTimestamp, encodedDecimal...), nil
	default:
		return nil, fmt.Errorf("unhandled type; supported types are: *llo.Decimal or *llo.TimestampedStreamValue; got: %T", sv)
	}
}

func (a ABIEncoder) EncodePadded(sv llo.StreamValue) ([]byte, error) {
	switch v := sv.(type) {
	case *llo.Decimal:
		if len(a.encoders) != 1 {
			return nil, fmt.Errorf("expected exactly one encoder; got: %d", len(a.encoders))
		}
		return a.encoders[0].encodeDecimalStreamValuePadded(v)
	case *llo.TimestampedStreamValue:
		if len(a.encoders) != 2 {
			return nil, fmt.Errorf("expected exactly two encoders for *llo.TimestampedStreamValue; got: %d", len(a.encoders))
		}
		// encode as two zero-padded 32 byte evm words:
		// uint256 observed_at_nanoseconds
		// uint256 stream_value
		encodedTimestamp, err := a.encoders[0].encodeUint64Padded(v.ObservedAtNanoseconds)
		if err != nil {
			return nil, fmt.Errorf("failed to encode timestamped stream value; %w", err)
		}
		var encodedDecimal []byte
		switch d := v.StreamValue.(type) {
		case *llo.Decimal:
			var err error
			encodedDecimal, err = a.encoders[1].encodeDecimalStreamValuePadded(d)
			if err != nil {
				return nil, fmt.Errorf("failed to encode nested stream value; %w", err)
			}
		default:
			return nil, fmt.Errorf("unhandled type; supported nested types for *llo.TimestampedStreamValue are: *llo.Decimal; got: %T", d)
		}
		return append(encodedTimestamp, encodedDecimal...), nil
	default:
		return nil, fmt.Errorf("unhandled type; supported types are: *llo.Decimal or *llo.TimestampedStreamValue; got: %T", sv)
	}
}

type singleABIEncoder struct {
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
func (a singleABIEncoder) getNormalizedMultiplier() (multiplier decimal.Decimal) {
	if a.Multiplier == nil {
		multiplier = decimal.NewFromInt(1)
	} else {
		multiplier = decimal.NewFromBigInt(a.Multiplier.ToInt(), 0)
	}
	return
}

func (a singleABIEncoder) applyMultiplier(d decimal.Decimal) *big.Int {
	return d.Mul(a.getNormalizedMultiplier()).BigInt()
}

// EncodePadded uses standard ABI encoding to encode the stream value, padding
// result to a multiple of 32 bytes
func (a singleABIEncoder) EncodePadded(sv llo.StreamValue) (b []byte, err error) {
	switch sv := sv.(type) {
	case *llo.Decimal:
		b, err = a.encodeDecimalStreamValuePadded(sv)
	default:
		return nil, fmt.Errorf("unhandled type; supported types are: *llo.Decimal; got: %T", sv)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode stream value %v (%T) with ABI type %q; %w", sv, sv, a.Type, err)
	}
	return b, nil
}

func (a singleABIEncoder) encodeUint64Padded(v uint64) (b []byte, err error) {
	return EncodePaddedBigInt(a.applyMultiplier(decimal.NewFromBigInt(new(big.Int).SetUint64(v), 0)), a.Type)
}

func (a singleABIEncoder) encodeDecimalStreamValuePadded(sv *llo.Decimal) (b []byte, err error) {
	if sv == nil {
		return nil, fmt.Errorf("expected non-nil *Decimal; got: %v", sv)
	}
	return EncodePaddedBigInt(a.applyMultiplier(sv.Decimal()), a.Type)
}

func EncodePaddedBigInt(v *big.Int, t string) (b []byte, err error) {
	b, err = EncodePackedBigInt(v, t)
	if err != nil {
		return nil, fmt.Errorf("failed to encode stream value %v with ABI type %q; %w", v, t, err)
	}
	if len(b) > 32 {
		return nil, fmt.Errorf("encoded stream value %v with ABI type %q is too large; got: %d bytes, max: 32 bytes", v, t, len(b))
	}
	if v.Sign() < 0 {
		// If the value is negative, we need to sign-extend it to 32 bytes.
		// This is done by padding with 0xff bytes.
		return padWithOnes(b), nil
	}
	return padWithZeroes(b), nil
}

func padWithZeroes(b []byte) []byte {
	return append(make([]byte, 32-len(b)), b...)
}

func padWithOnes(b []byte) []byte {
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	for i := 0; i < 32-len(b); i++ {
		padded[i] = 0xff
	}
	return padded
}

const ZeroBytesSentinel = "bytes0"

func (a singleABIEncoder) encodeUint64Packed(v uint64) ([]byte, error) {
	if a.Type == ZeroBytesSentinel {
		return nil, nil
	}
	return EncodePackedBigInt(a.applyMultiplier(decimal.NewFromBigInt(new(big.Int).SetUint64(v), 0)), a.Type)
}

// EncodePacked uses packed ABI encoding to encode the stream value (no padding)
func (a singleABIEncoder) encodePacked(sv llo.StreamValue) ([]byte, error) {
	if a.Type == ZeroBytesSentinel {
		return nil, nil
	}
	var v *big.Int
	switch sv := sv.(type) {
	case *llo.Decimal:
		if sv == nil {
			return nil, fmt.Errorf("expected non-nil *Decimal; got: %v", sv)
		}
		v = a.applyMultiplier(sv.Decimal())
	default:
		return nil, fmt.Errorf("EncodePacked only currently supports StreamValue type of *llo.Decimal, got: %T", sv)
	}
	return EncodePackedBigInt(v, a.Type)
}

// regex to match Solidity integer types (uint/int) and extract bit width
var typeRegex = regexp.MustCompile(`^(u?int)(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)$`)

// EncodePackedBigInt converts a *big.Int to packed EVM bytes according to the Solidity type.
// For unsigned types ("uintN"), the value must be non-negative and fit in N bits.
// For signed types ("intN"), the value must be between -2^(N-1) and 2^(N-1)-1; then the function returns the two's complement representation in N bits.
func EncodePackedBigInt(value *big.Int, typeStr string) ([]byte, error) {
	// Validate and extract type and bit width.
	matches := typeRegex.FindStringSubmatch(typeStr)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid Solidity type: %s", typeStr)
	}
	typePrefix := matches[1] // "uint" or "int"
	bitWidthI, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse bit width: %w", err)
	}
	bitWidth := uint(bitWidthI) //nolint:gosec // G115 // bit width is never above 256, so no overflow issues
	byteLen := bitWidth / 8

	// For unsigned integers, value must be non-negative and within range.
	if typePrefix == "uint" {
		if value.Sign() < 0 {
			return nil, fmt.Errorf("negative value provided for unsigned type %s", typeStr)
		}
		// Check that value fits in bitWidth bits.
		// maxSize = 2^bitWidth
		maxSize := new(big.Int).Lsh(big.NewInt(1), bitWidth)
		if value.Cmp(maxSize) >= 0 {
			return nil, fmt.Errorf("value %s out of range for type %s", value.String(), typeStr)
		}
		// Create a byte slice of fixed length and fill it with the big-endian bytes.
		result := make([]byte, byteLen)
		value.FillBytes(result)
		return result, nil
	}

	// For signed integers, first ensure that value fits in the representable range.
	// The valid range is from -2^(bitWidth-1) to 2^(bitWidth-1)-1.
	// 2^(bitWidth-1)
	minInt := new(big.Int).Lsh(big.NewInt(1), bitWidth-1)
	// -2^(bitWidth-1)
	minInt.Neg(minInt)
	// 2^(bitWidth-1) - 1
	maxInt := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bitWidth-1), big.NewInt(1))
	if value.Cmp(minInt) < 0 || value.Cmp(maxInt) > 0 {
		return nil, fmt.Errorf("value %s out of range for type %s", value.String(), typeStr)
	}

	// Solidity uses two's complement in fixed N-bit representation.
	// twoPow = 2^bitWidth
	twoPow := new(big.Int).Lsh(big.NewInt(1), bitWidth)
	// Compute the modulo to obtain the two's complement representation.
	modValue := new(big.Int).Mod(value, twoPow)
	result := make([]byte, byteLen)
	modValue.FillBytes(result)
	return result, nil
}
