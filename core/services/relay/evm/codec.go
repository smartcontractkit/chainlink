package evm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

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
var evmDecoderHooks = []mapstructure.DecodeHookFunc{decodeAccountAndAllowArraySliceHook, codec.BigIntHook, codec.SliceToArrayVerifySizeHook, sizeVerifyBigIntHook}

// NewCodec creates a new [commontypes.RemoteCodec] for EVM.
// Note that names in the ABI are converted to Go names using [abi.ToCamelCase],
// this is per convention in [abi.MakeTopics], [abi.Arguments.Pack] etc.
// This allows names on-chain to be in go convention when generated.
// It means that if you need to use a [codec.Modifier] to reference a field
// you need to use the Go name instead of the name on-chain.
// eg: rename FooBar -> Bar, not foo_bar_ to Bar if the name on-chain is foo_bar_
func NewCodec(conf types.CodecConfig) (commontypes.RemoteCodec, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]types.CodecEntry{},
		decoderDefs: map[string]types.CodecEntry{},
	}

	for k, v := range conf.Configs {
		args := abi.Arguments{}
		if err := json.Unmarshal(([]byte)(v.TypeABI), &args); err != nil {
			return nil, err
		}

		mod, err := v.ModifierConfigs.ToModifier(evmDecoderHooks...)
		if err != nil {
			return nil, err
		}

		item := types.NewCodecEntry(args, nil, mod)
		if err = item.Init(); err != nil {
			return nil, err
		}

		parsed.encoderDefs[k] = item
		parsed.decoderDefs[k] = item
	}

	return parsed.toCodec()
}

type evmCodec struct {
	*encoder
	*decoder
	*parsedTypes
}

func (c *evmCodec) CreateType(itemType string, forEncoding bool) (any, error) {
	var itemTypes map[string]types.CodecEntry
	if forEncoding {
		itemTypes = c.encoderDefs
	} else {
		itemTypes = c.decoderDefs
	}

	def, ok := itemTypes[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find type name %s", commontypes.ErrInvalidType, itemType)
	}

	return reflect.New(def.CheckedType()).Interface(), nil
}

var bigIntType = reflect.TypeOf((*big.Int)(nil))

func sizeVerifyBigIntHook(from, to reflect.Type, data any) (any, error) {
	if from.Implements(types.SizedBigIntType()) &&
		!to.Implements(types.SizedBigIntType()) &&
		!reflect.PointerTo(to).Implements(types.SizedBigIntType()) {
		return codec.BigIntHook(from, bigIntType, reflect.ValueOf(data).Convert(bigIntType).Interface())
	}

	if !to.Implements(types.SizedBigIntType()) {
		return data, nil
	}

	var err error
	data, err = codec.BigIntHook(from, bigIntType, data)
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
