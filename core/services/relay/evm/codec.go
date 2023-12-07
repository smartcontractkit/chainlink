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

var evmDecoderHooks = []mapstructure.DecodeHookFunc{decodeAccountHook, codec.BigIntHook, codec.SliceToArrayVerifySizeHook, sizeVerifyBigIntHook}

func NewCodec(conf types.CodecConfig) (commontypes.RemoteCodec, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]*CodecEntry{},
		decoderDefs: map[string]*CodecEntry{},
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

		item := &CodecEntry{Args: args, mod: mod}
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
	var itemTypes map[string]*CodecEntry
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
