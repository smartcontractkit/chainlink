package ccipsolana

import (
	"encoding/binary"
	"fmt"
	"reflect"

	agbinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
)

const svmDestExecDataKey = "destGasAmount"

// DecodeExtraArgsToMap is a helper function for converting Borsh encoded extra args bytes into map[string]any, which will be saved in ocr report.message.ExtraArgsDecoded
func DecodeExtraArgsToMap(extraArgs []byte) (map[string]any, error) {
	outputMap := make(map[string]any)
	var args ccip_router.AnyExtraArgs
	decoder := agbinary.NewBorshDecoder(extraArgs)
	err := args.UnmarshalWithDecoder(decoder)
	if err != nil {
		return outputMap, fmt.Errorf("failed to decode extra args: %w", err)
	}

	val := reflect.ValueOf(args)
	typ := reflect.TypeOf(args)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		outputMap[field.Name] = fieldValue
	}

	return outputMap, nil
}

func DecodeDestExecDataToMap(destExecData []byte) (map[string]any, error) {
	return map[string]interface{}{
		svmDestExecDataKey: bytesToUint32LE(destExecData),
	}, nil
}

func bytesToUint32LE(b []byte) uint32 {
	if len(b) < 4 {
		var padded [4]byte
		copy(padded[:len(b)], b) // Pad from the right for little-endian
		return binary.LittleEndian.Uint32(padded[:])
	}

	return binary.LittleEndian.Uint32(b)
}
