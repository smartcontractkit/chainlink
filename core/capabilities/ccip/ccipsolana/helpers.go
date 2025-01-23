package ccipsolana

import (
	"fmt"
	"reflect"

	agbinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
)

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
