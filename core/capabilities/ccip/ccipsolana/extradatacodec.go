package ccipsolana

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	agbinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/fee_quoter"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

const (
	svmDestExecDataKey = "destGasAmount"
	evmGasLimitKey     = "GasLimit"
)

var (
	// tag definition https://github.com/smartcontractkit/chainlink-ccip/blob/1b2ee24da54bddef8f3943dc84102686f2890f87/chains/solana/contracts/programs/ccip-router/src/extra_args.rs#L8C21-L11C45
	// this should be moved to msghasher.go once merged

	// bytes4(keccak256("CCIP SVMExtraArgsV1"));
	svmExtraArgsV1Tag = hexutil.MustDecode("0x1f3b3aba")

	// bytes4(keccak256("CCIP EVMExtraArgsV2"));
	evmExtraArgsV2Tag = hexutil.MustDecode("0x181dcf10")
)

// ExtraDataDecoder is a helper struct for decoding extra data
type ExtraDataDecoder struct{}

// DecodeExtraArgsToMap is a helper function for converting Borsh encoded extra args bytes into map[string]any
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs ccipocr3.Bytes) (map[string]any, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var val reflect.Value
	var typ reflect.Type
	outputMap := make(map[string]any)
	switch string(extraArgs[:4]) {
	case string(evmExtraArgsV2Tag):
		var args fee_quoter.GenericExtraArgsV2
		decoder := agbinary.NewBorshDecoder(extraArgs[4:])
		err := args.UnmarshalWithDecoder(decoder)
		if err != nil {
			return nil, fmt.Errorf("failed to decode extra args: %w", err)
		}
		val = reflect.ValueOf(args)
		typ = reflect.TypeFor[fee_quoter.GenericExtraArgsV2]()
	case string(svmExtraArgsV1Tag):
		var args fee_quoter.SVMExtraArgsV1
		decoder := agbinary.NewBorshDecoder(extraArgs[4:])
		err := args.UnmarshalWithDecoder(decoder)
		if err != nil {
			return nil, fmt.Errorf("failed to decode extra args: %w", err)
		}
		val = reflect.ValueOf(args)
		typ = reflect.TypeFor[fee_quoter.SVMExtraArgsV1]()
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs[:4])
	}

	for i := range val.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		if field.Name == evmGasLimitKey {
			// convert SVM Borsh specific type uint128 to *big.Int for EVM gas limit
			gl, ok := fieldValue.(agbinary.Uint128)
			if !ok {
				return nil, fmt.Errorf("expected field %s to be of type agbinary.Uint128, got %T", field.Name, fieldValue)
			}

			fieldValue = gl.BigInt()
		}
		outputMap[field.Name] = fieldValue
	}

	return outputMap, nil
}

// DecodeDestExecDataToMap is a helper function for converting dest exec data bytes into map[string]any
func (d ExtraDataDecoder) DecodeDestExecDataToMap(destExecData ccipocr3.Bytes) (map[string]any, error) {
	return map[string]any{
		svmDestExecDataKey: binary.BigEndian.Uint32(destExecData),
	}, nil
}

// Ensure ExtraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &ExtraDataDecoder{}
