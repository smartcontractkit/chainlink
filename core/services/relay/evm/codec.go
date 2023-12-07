package evm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// decodeAccountHook allows strings to be converted to [32]byte allowing config to represent them as 0x...
// BigIntHook allows *big.Int to be represented as any integer type or a string and to go back to them.
// Useful for config, or if when a model may use a go type that isn't a *big.Int when Pack expects one.
// Eg: int32 in a go struct from a plugin could require a *big.Int in Pack for int24, if it fits, we shouldn't care.
// SliceToArrayVerifySizeHook verifies that slices have the correct size when converting to an array
// sizeVerifyBigIntHook allows our custom types that verify the number fits in the on-chain type to be converted as-if
// it was a *big.Int
var evmDecoderHooks = []mapstructure.DecodeHookFunc{decodeAccountHook, codec.BigIntHook, codec.SliceToArrayVerifySizeHook, sizeVerifyBigIntHook}

func NewCodec(conf types.CodecConfig) (commontypes.CodecTypeProvider, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]*codecEntry{},
		decoderDefs: map[string]*codecEntry{},
	}

	for k, v := range conf.ChainCodecConfigs {
		args := abi.Arguments{}
		if err := json.Unmarshal(([]byte)(v.TypeAbi), &args); err != nil {
			return nil, err
		}

		mod, err := v.ModifierConfigs.ToModifier(evmDecoderHooks...)
		if err != nil {
			return nil, err
		}

		item := &codecEntry{Args: args, mod: mod}
		if err := item.Init(); err != nil {
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
	var itemTypes map[string]*codecEntry
	if forEncoding {
		itemTypes = c.encoderDefs
	} else {
		itemTypes = c.decoderDefs
	}

	def, ok := itemTypes[itemType]
	if !ok {
		return nil, commontypes.ErrInvalidType
	}

	return reflect.New(def.checkedType).Interface(), nil
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

func decodeAccountHook(from, to reflect.Type, data any) (any, error) {
	b32, _ := types.GetType("bytes32")
	if from.Kind() == reflect.String && to == b32.Checked {
		decoded, err := hexutil.Decode(data.(string))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
		}
		return [32]byte(decoded), nil
	}
	return data, nil
}
