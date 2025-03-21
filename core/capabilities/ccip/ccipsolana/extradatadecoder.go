package ccipsolana

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	agbinary "github.com/gagliardetto/binary"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
)

const (
	svmDestExecDataKey = "destGasAmount"
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
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
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
			return outputMap, fmt.Errorf("failed to decode extra args: %w", err)
		}
		val = reflect.ValueOf(args)
		typ = reflect.TypeOf(args)
	case string(svmExtraArgsV1Tag):
		var args fee_quoter.SVMExtraArgsV1
		decoder := agbinary.NewBorshDecoder(extraArgs[4:])
		err := args.UnmarshalWithDecoder(decoder)
		if err != nil {
			return outputMap, fmt.Errorf("failed to decode extra args: %w", err)
		}
		val = reflect.ValueOf(args)
		typ = reflect.TypeOf(args)
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		outputMap[field.Name] = fieldValue
	}

	return outputMap, nil
}

// DecodeDestExecDataToMap is a helper function for converting dest exec data bytes into map[string]any
func (d ExtraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	return map[string]any{
		svmDestExecDataKey: binary.BigEndian.Uint32(destExecData),
	}, nil
}
