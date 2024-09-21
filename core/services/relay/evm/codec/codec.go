package codec

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-viper/mapstructure/v2"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// DecoderHooks
//
// decodeAccountAndAllowArraySliceHook allows:
//
//	strings to be converted to [32]byte allowing config to represent them as 0x...
//	slices or arrays to be converted to a pointer to that type
//
// BigIntHook allows *big.Int to be represented as any integer type or a string and to go back to them.
// Useful for config, or if when a model may use a go type that isn't a *big.Int when Pack expects one.
// Eg: int32 in a go struct from a plugin could require a *big.Int in Pack for int24, if it fits, we shouldn't care.
// SliceToArrayVerifySizeHook verifies that slices have the correct size when converting to an array
// sizeVerifyBigIntHook allows our custom types that verify the number fits in the on-chain type to be converted as-if
// it was a *big.Int
var DecoderHooks = []mapstructure.DecodeHookFunc{
	decodeAccountAndAllowArraySliceHook,
	commoncodec.BigIntHook,
	commoncodec.SliceToArrayVerifySizeHook,
	sizeVerifyBigIntHook,
	commoncodec.NumberHook,
	addressStringDecodeHook,
}

// NewCodec creates a new [commontypes.RemoteCodec] for EVM.
// Note that names in the ABI are converted to Go names using [abi.ToCamelCase],
// this is per convention in [abi.MakeTopics], [abi.Arguments.Pack] etc.
// This allows names on-chain to be in go convention when generated.
// It means that if you need to use a [codec.Modifier] to reference a field
// you need to use the Go name instead of the name on-chain.
// eg: rename FooBar -> Bar, not foo_bar_ to Bar if the name on-chain is foo_bar_
func NewCodec(conf types.CodecConfig) (commontypes.RemoteCodec, error) {
	parsed := &ParsedTypes{
		EncoderDefs: map[string]types.CodecEntry{},
		DecoderDefs: map[string]types.CodecEntry{},
	}

	for k, v := range conf.Configs {
		args := abi.Arguments{}
		if err := json.Unmarshal(([]byte)(v.TypeABI), &args); err != nil {
			return nil, err
		}

		mod, err := v.ModifierConfigs.ToModifier(DecoderHooks...)
		if err != nil {
			return nil, err
		}

		item := types.NewCodecEntry(args, nil, mod)
		if err = item.Init(); err != nil {
			return nil, err
		}

		parsed.EncoderDefs[k] = item
		parsed.DecoderDefs[k] = item
	}

	return parsed.ToCodec()
}

type evmCodec struct {
	*encoder
	*decoder
	*ParsedTypes
}

func (c *evmCodec) CreateType(itemType string, forEncoding bool) (any, error) {
	var itemTypes map[string]types.CodecEntry
	if forEncoding {
		itemTypes = c.EncoderDefs
	} else {
		itemTypes = c.DecoderDefs
	}

	def, ok := itemTypes[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find type name %s", commontypes.ErrInvalidType, itemType)
	}

	return reflect.New(def.CheckedType()).Interface(), nil
}

func WrapItemType(contractName, itemType string, isParams bool) string {
	if isParams {
		return fmt.Sprintf("params.%s.%s", contractName, itemType)
	}

	return fmt.Sprintf("return.%s.%s", contractName, itemType)
}

var bigIntType = reflect.TypeOf((*big.Int)(nil))

func sizeVerifyBigIntHook(from, to reflect.Type, data any) (any, error) {
	if from.Implements(types.SizedBigIntType()) &&
		!to.Implements(types.SizedBigIntType()) &&
		!reflect.PointerTo(to).Implements(types.SizedBigIntType()) {
		return commoncodec.BigIntHook(from, bigIntType, reflect.ValueOf(data).Convert(bigIntType).Interface())
	}

	if !to.Implements(types.SizedBigIntType()) {
		return data, nil
	}

	var err error
	data, err = commoncodec.BigIntHook(from, bigIntType, data)
	if err != nil {
		return nil, err
	}

	bi, ok := data.(*big.Int)
	if !ok {
		return data, nil
	}

	converted := reflect.ValueOf(bi).Convert(to).Interface().(types.SizedBigInt)
	return converted, converted.Verify()
}

func decodeAccountAndAllowArraySliceHook(from, to reflect.Type, data any) (any, error) {
	if from.Kind() == reflect.String &&
		(to == reflect.TypeOf(common.Address{}) || to == reflect.TypeOf(&common.Address{})) {
		return decodeAddress(data)
	}

	if from.Kind() == reflect.Pointer && to.Kind() != reflect.Pointer && from != nil &&
		(from.Elem().Kind() == reflect.Slice || from.Elem().Kind() == reflect.Array) {
		return reflect.ValueOf(data).Elem().Interface(), nil
	}

	return data, nil
}

// decodeAddress attempts to decode a hex-encoded Ethereum address from a string.
// Returns an error in the following cases:
//   - If the input is not a valid hex string
//   - If the decoded length does not match the expected length
func decodeAddress(data any) (any, error) {
	decoded, err := hexutil.Decode(data.(string))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	} else if len(decoded) != common.AddressLength {
		return nil, fmt.Errorf(
			"%w: wrong number size for address expected %v got %v",
			commontypes.ErrSliceWrongLen,
			common.AddressLength, len(decoded))
	}

	return common.Address(decoded), nil
}

// addressStringDecodeHook is a decode hook that converts between `from` and `to` types involving string and common.Address types.
// It handles the following conversions:
// 1. `from` string or *string -> `to` common.Address or *common.Address
// 2. `from` common.Address or *common.Address -> `to` string or *string
//
// The function handles invalid `from` values and `nil` pointers:
//   - If `from` is a string or *string and is invalid (e.g., an empty string or a non-hex string),
//     it returns an appropriate error (types.ErrInvalidType).
//   - If `from` is an empty common.Address{} or *common.Address, the function returns an error
//     (types.ErrInvalidType) instead of treating it as the zero address ("0x0000000000000000000000000000000000000000").
//   - If `from` is a nil *string or nil *common.Address, the function returns an error (types.ErrInvalidType).
//
// For unsupported `from` and `to` conversions, the function returns the original value unchanged.
func addressStringDecodeHook(from reflect.Type, to reflect.Type, value interface{}) (interface{}, error) {
	// Helper variables for types to improve readability
	stringType, stringPtrType := reflect.TypeOf(""), reflect.PointerTo(reflect.TypeOf(""))
	addressType, addressPtrType := reflect.TypeOf(common.Address{}), reflect.TypeOf(&common.Address{})

	if (from == stringType || from == stringPtrType) &&
		(to == addressType || to == addressPtrType) {
		strValue, err := extractStringValue(from, value, stringPtrType)
		if err != nil {
			return nil, err
		}

		address, err := decodeAddress(strValue)
		if err != nil {
			return nil, err
		}

		return returnAddressValue(to, address, addressPtrType)
	}

	if (from == addressType || from == addressPtrType) &&
		(to == stringType || to == stringPtrType) {
		if from == addressPtrType && (value == nil || reflect.ValueOf(value).IsNil()) {
			return nil, fmt.Errorf("%w: nil *common.Address value", commontypes.ErrInvalidType)
		}

		addressStr, err := extractAddressToHex(from, value, addressPtrType)
		if err != nil {
			return nil, err
		}

		return returnStringValue(to, addressStr, stringPtrType)
	}

	// Return the original value unchanged for unsupported conversions
	return value, nil
}

// Helper function to extract string or *string values
func extractStringValue(from reflect.Type, value interface{}, stringPtrType reflect.Type) (string, error) {
	if from == stringPtrType {
		if value == nil || reflect.ValueOf(value).IsNil() {
			return "", fmt.Errorf("%w: nil *string value", commontypes.ErrInvalidType)
		}
		return *value.(*string), nil
	}

	return value.(string), nil
}

// Helper function to return *common.Address or common.Address depending on 'to' type
func returnAddressValue(to reflect.Type, address interface{}, addressPtrType reflect.Type) (interface{}, error) {
	if to == addressPtrType {
		addr := address.(common.Address)
		return &addr, nil
	}
	return address, nil
}

// Helper function to extract hex string from common.Address or *common.Address
func extractAddressToHex(from reflect.Type, value interface{}, addressPtrType reflect.Type) (string, error) {
	if from == addressPtrType {
		if (*value.(*common.Address) == common.Address{}) {
			return "", fmt.Errorf("%w: empty address", commontypes.ErrInvalidType)
		}
		return value.(*common.Address).Hex(), nil
	}

	if (value.(common.Address) == common.Address{}) {
		return "", fmt.Errorf("%w: empty address", commontypes.ErrInvalidType)
	}

	return value.(common.Address).Hex(), nil
}

// Helper function to return string as *string or string depending on 'to' type
func returnStringValue(to reflect.Type, addressStr string, stringPtrType reflect.Type) (interface{}, error) {
	if to == stringPtrType {
		return &addressStr, nil
	}

	return addressStr, nil
}
