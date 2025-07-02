package ccipton

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const (
	tvmDestExecDataKey = "destGasAmount"
)

var (
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
		var args bindings.GenericExtraArgsV2
		c, err := cell.FromBOC(extraArgs[4:])
		if err != nil {
			return outputMap, fmt.Errorf("decode BOC: %w", err)
		}

		if err = tlb.LoadFromCell(&args, c.BeginParse()); err != nil {
			return nil, fmt.Errorf("failed to tlb load extra args from cell: %w", err)
		}

		val = reflect.ValueOf(args)
		typ = reflect.TypeOf(args)
	case string(svmExtraArgsV1Tag):
		var tlbArgs bindings.TLBSVMExtraArgsV1
		c, err := cell.FromBOC(extraArgs[4:])
		if err != nil {
			return outputMap, fmt.Errorf("decode BOC: %w", err)
		}

		if err = tlb.LoadFromCell(&tlbArgs, c.BeginParse()); err != nil {
			return nil, fmt.Errorf("failed to tlb load extra args from cell: %w", err)
		}

		args, err := tlbArgs.ExportSVMExtraArgsV1()
		if err != nil {
			return nil, fmt.Errorf("failed to export SVM extra args with tlb format: %w", err)
		}

		val = reflect.ValueOf(args)
		typ = reflect.TypeOf(args)
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs[:4])
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
		tvmDestExecDataKey: binary.BigEndian.Uint32(destExecData),
	}, nil
}

// Ensure ExtraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &ExtraDataDecoder{}
