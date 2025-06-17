package ccipaptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common/hexutil"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// ExtraDataDecoder is a concrete implementation of ccipcommon.ExtraDataDecoder
// Compatible with ccip::fee_quoter version 1.6.0
type ExtraDataDecoder struct{}

var _ ccipcommon.SourceChainExtraDataCodec = ExtraDataDecoder{}

const (
	aptosDestExecDataKey = "destGasAmount"
)

var (
	// bytes4 public constant EVM_EXTRA_ARGS_V1_TAG = 0x97a657c9;
	evmExtraArgsV1Tag = hexutil.MustDecode("0x97a657c9")

	// bytes4 public constant GENERIC_EXTRA_ARGS_V2_TAG = 0x181dcf10;
	genericExtraArgsV2Tag = hexutil.MustDecode("0x181dcf10")

	// bytes4 public constant SVM_EXTRA_EXTRA_ARGS_V1_TAG = 0x1f3b3aba
	svmExtraArgsV1Tag = hexutil.MustDecode("0x1f3b3aba")
)

// DecodeDestExecDataToMap reformats bytes into a chain agnostic map[string]interface{} representation for dest exec data
func (d ExtraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	des := bcs.NewDeserializer(destExecData)
	if des.Remaining() != 4 {
		return nil, fmt.Errorf("dest exec data invalid length: %d, should be 4 bytes", des.Remaining())
	}

	destGasAmount := des.U32()
	if des.Error() != nil {
		return nil, fmt.Errorf("decode dest gas amount: %w", des.Error())
	}

	return map[string]any{aptosDestExecDataKey: destGasAmount}, nil
}

// DecodeExtraArgsToMap reformats bytes into a chain agnostic map[string]any representation for extra args
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	des := bcs.NewDeserializer(extraArgs)
	tag := des.ReadFixedBytes(4)

	switch string(tag) {
	case string(evmExtraArgsV1Tag):
		return d.decodeEvmExtraArgsV1(des)
	case string(genericExtraArgsV2Tag):
		return d.decodeGenericExtraArgsV2(des)
	case string(svmExtraArgsV1Tag):
		return d.decodeSvmExtraArgsV1(des)
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", tag)
	}
}

func (d ExtraDataDecoder) decodeEvmExtraArgsV1(des *bcs.Deserializer) (map[string]any, error) {
	extraArgs := make(map[string]any)

	gasLimit := des.U256()
	if des.Error() != nil {
		return nil, fmt.Errorf("error whilst decoding evm extra args v1: %w", des.Error())
	}

	extraArgs["gasLimit"] = &gasLimit
	return extraArgs, nil
}

func (d ExtraDataDecoder) decodeGenericExtraArgsV2(des *bcs.Deserializer) (map[string]any, error) {
	extraArgs := make(map[string]any)

	gasLimit := des.U256()
	if des.Error() != nil {
		return nil, fmt.Errorf("error whilst decoding generic extra args v2: %w", des.Error())
	}

	extraArgs["gasLimit"] = &gasLimit

	allowOutOfOrderExecution := des.Bool()
	if des.Error() != nil {
		// Default to false if not present, consistent with original code.
		extraArgs["allowOutOfOrderExecution"] = false
	} else {
		extraArgs["allowOutOfOrderExecution"] = allowOutOfOrderExecution
	}

	return extraArgs, nil
}

func (d ExtraDataDecoder) decodeSvmExtraArgsV1(des *bcs.Deserializer) (map[string]any, error) {
	extraArgs := make(map[string]any)

	computeUnits := des.U32()
	if des.Error() != nil {
		return nil, fmt.Errorf("error whilst decoding svm extra args v1: %w", des.Error())
	}
	extraArgs["computeUnits"] = computeUnits

	accountIsWritableBitmap := des.U64()
	if des.Error() != nil {
		extraArgs["accountIsWritableBitmap"] = nil
	} else {
		extraArgs["accountIsWritableBitmap"] = accountIsWritableBitmap
	}

	allowOutOfOrderExecution := des.Bool()
	if des.Error() != nil {
		extraArgs["allowOutOfOrderExecution"] = false
	} else {
		extraArgs["allowOutOfOrderExecution"] = allowOutOfOrderExecution
	}

	tokenReceiver := des.ReadBytes()
	if des.Error() != nil {
		extraArgs["tokenReceiver"] = nil
	} else {
		extraArgs["tokenReceiver"] = tokenReceiver
	}

	// accounts is an array of bytes32 values
	accounts := bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *[]byte) {
		*item = des.ReadBytes()
	})

	if des.Error() != nil {
		extraArgs["accounts"] = nil
	} else {
		extraArgs["accounts"] = accounts
	}

	return extraArgs, nil
}
