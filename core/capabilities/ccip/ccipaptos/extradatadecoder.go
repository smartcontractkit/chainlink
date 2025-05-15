package ccipaptos

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

	uint32Type     = mustNewType("uint32")
	uint64Type     = mustNewType("uint64")
	uint256Type    = mustNewType("uint256")
	boolType       = mustNewType("bool")
	bytes32Type    = mustNewType("bytes32")
	bytes32ArrType = mustNewType("bytes32[]")

	// Arguments for decoding destGasAmount
	destGasAmountArguments = abi.Arguments{
		{Name: aptosDestExecDataKey, Type: uint32Type},
	}

	// Arguments matching the fields of EVMExtraArgsV1 struct
	evmExtraArgsV1Fields = abi.Arguments{
		{Name: "gasLimit", Type: uint256Type},
	}

	// Arguments matching the fields of GenericExtraArgsV2 struct
	genericExtraArgsV2Fields = abi.Arguments{
		{Name: "gasLimit", Type: uint256Type},
		{Name: "allowOutOfOrderExecution", Type: boolType},
	}

	// Arguments matching the fields of SVMExtraArgsV1 struct
	svmExtraArgsV1Fields = abi.Arguments{
		{Name: "computeUnits", Type: uint32Type},
		{Name: "accountIsWritableBitmap", Type: uint64Type},
		{Name: "allowOutOfOrderExecution", Type: boolType},
		{Name: "tokenReceiver", Type: bytes32Type},
		{Name: "accounts", Type: bytes32ArrType},
	}
)

func mustNewType(typeStr string) abi.Type {
	t, err := abi.NewType(typeStr, "", nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create ABI type %s: %v", typeStr, err))
	}
	return t
}

// DecodeDestExecDataToMap reformats bytes into a chain agnostic map[string]interface{} representation for dest exec data
func (d ExtraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	args := make(map[string]any)
	err := destGasAmountArguments.UnpackIntoMap(args, destExecData)
	if err != nil {
		if len(destExecData) != 32 {
			return nil, fmt.Errorf("decode dest gas amount: expected 32 bytes for uint32, got %d: %w", len(destExecData), err)
		}

		var val big.Int
		val.SetBytes(destExecData)
		if val.Cmp(big.NewInt(0xFFFFFFFF)) > 0 {
			return nil, fmt.Errorf("decode dest gas amount: value %s exceeds uint32 max: %w", val.String(), err)
		}

		return nil, fmt.Errorf("decode dest gas amount: %w", err)
	}

	if _, ok := args[aptosDestExecDataKey]; !ok {
		return nil, fmt.Errorf("failed to unpack key '%s' for dest gas amount", aptosDestExecDataKey)
	}

	return args, nil
}

// DecodeExtraArgsToMap reformats bytes into a chain agnostic map[string]any representation for extra args
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var decoderArgs abi.Arguments
	tag := string(extraArgs[:4])
	argsData := extraArgs[4:]

	switch tag {
	case string(evmExtraArgsV1Tag):
		decoderArgs = evmExtraArgsV1Fields
	case string(genericExtraArgsV2Tag):
		decoderArgs = genericExtraArgsV2Fields
	case string(svmExtraArgsV1Tag):
		decoderArgs = svmExtraArgsV1Fields
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs[:4])
	}

	unpackedArgs := make(map[string]any)
	err := decoderArgs.UnpackIntoMap(unpackedArgs, argsData)
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args (tag %x): %w", extraArgs[:4], err)
	}

	output := make(map[string]any)

	switch tag {
	case string(evmExtraArgsV1Tag):
		if gasLimit, ok := unpackedArgs["gasLimit"].(*big.Int); ok {
			output["gasLimit"] = gasLimit
		} else if unpackedArgs["gasLimit"] != nil {
			return nil, fmt.Errorf("unexpected type for gasLimit in EVMExtraArgsV1: %T", unpackedArgs["gasLimit"])
		} else {
			output["gasLimit"] = nil // Field not present or nil
		}

	case string(genericExtraArgsV2Tag):
		if gasLimit, ok := unpackedArgs["gasLimit"].(*big.Int); ok {
			output["gasLimit"] = gasLimit
		} else if unpackedArgs["gasLimit"] != nil {
			return nil, fmt.Errorf("unexpected type for gasLimit in GenericExtraArgsV2: %T", unpackedArgs["gasLimit"])
		} else {
			output["gasLimit"] = nil
		}

		if allow, ok := unpackedArgs["allowOutOfOrderExecution"].(bool); ok {
			output["allowOutOfOrderExecution"] = allow
		} else if unpackedArgs["allowOutOfOrderExecution"] != nil {
			return nil, fmt.Errorf("unexpected type for allowOutOfOrderExecution in GenericExtraArgsV2: %T", unpackedArgs["allowOutOfOrderExecution"])
		} else {
			// Default to false if not present, consistent with original code.
			// Note: ABI decoding of bool usually doesn't result in nil, but checking doesn't hurt.
			output["allowOutOfOrderExecution"] = false
		}

	case string(svmExtraArgsV1Tag):
		if computeUnits, ok := unpackedArgs["computeUnits"].(uint32); ok {
			output["computeUnits"] = computeUnits
		} else if unpackedArgs["computeUnits"] != nil {
			return nil, fmt.Errorf("unexpected type for computeUnits in SVMExtraArgsV1: %T", unpackedArgs["computeUnits"])
		} else {
			output["computeUnits"] = nil
		}

		if bitmap, ok := unpackedArgs["accountIsWritableBitmap"].(uint64); ok {
			output["accountIsWritableBitmap"] = bitmap
		} else if unpackedArgs["accountIsWritableBitmap"] != nil {
			return nil, fmt.Errorf("unexpected type for accountIsWritableBitmap in SVMExtraArgsV1: %T", unpackedArgs["accountIsWritableBitmap"])
		} else {
			output["accountIsWritableBitmap"] = nil
		}

		if allow, ok := unpackedArgs["allowOutOfOrderExecution"].(bool); ok {
			output["allowOutOfOrderExecution"] = allow
		} else if unpackedArgs["allowOutOfOrderExecution"] != nil {
			return nil, fmt.Errorf("unexpected type for allowOutOfOrderExecution in SVMExtraArgsV1: %T", unpackedArgs["allowOutOfOrderExecution"])
		} else {
			output["allowOutOfOrderExecution"] = false
		}

		if tokenReceiver, ok := unpackedArgs["tokenReceiver"].([32]byte); ok {
			output["tokenReceiver"] = tokenReceiver
		} else if unpackedArgs["tokenReceiver"] != nil {
			return nil, fmt.Errorf("unexpected type for tokenReceiver in SVMExtraArgsV1: %T", unpackedArgs["tokenReceiver"])
		} else {
			output["tokenReceiver"] = nil
		}

		if accounts, ok := unpackedArgs["accounts"].([][32]byte); ok {
			output["accounts"] = accounts
		} else if unpackedArgs["accounts"] != nil {
			return nil, fmt.Errorf("unexpected type for accounts in SVMExtraArgsV1: %T", unpackedArgs["accounts"])
		} else {
			output["accounts"] = nil
		}
	}

	return output, nil
}
