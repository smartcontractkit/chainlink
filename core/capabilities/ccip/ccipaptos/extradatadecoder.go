package ccipaptos

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_aptos_utils"
)

// ExtraDataDecoder is a concrete implementation of ccipcommon.ExtraDataDecoder
// Compatible with ccip::fee_quoter version 1.6.0
type ExtraDataDecoder struct{}

const (
	aptosDestExecDataKey = "destGasAmount"
)

var aptosUtilsABI = abihelpers.MustParseABI(ccip_aptos_utils.AptosUtilsABI)

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
	destGasAmount, err := abiDecodeUint32(destExecData)
	if err != nil {
		return nil, fmt.Errorf("decode dest gas amount: %w", err)
	}

	return map[string]any{
		aptosDestExecDataKey: destGasAmount,
	}, nil
}

// DecodeExtraArgsToMap reformats bytes into a chain agnostic map[string]any representation for extra args
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var method string
	var argName string
	switch string(extraArgs[:4]) {
	case string(evmExtraArgsV1Tag):
		method = "exposeEVMExtraArgsV1"
		argName = "evmExtraArgsV1"
	case string(genericExtraArgsV2Tag):
		method = "exposeGenericExtraArgsV2"
		argName = "evmExtraArgsV2"
	case string(svmExtraArgsV1Tag):
		method = "exposeSVMExtraArgsV1"
		argName = "svmExtraArgsV1"
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}

	output := make(map[string]any)
	args := make(map[string]any)
	err := aptosUtilsABI.Methods[method].Inputs.UnpackIntoMap(args, extraArgs[4:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args %v: %w", method, err)
	}

	argValue, exists := args[argName]
	if !exists {
		return nil, fmt.Errorf("failed to get arg value for %s", argName)
	}

	val := reflect.ValueOf(argValue)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct for %s, got %T", argName, argValue)
	}

	switch argName {
	case "evmExtraArgsV1":
		gasLimitField := val.FieldByName("GasLimit")
		if gasLimitField.IsValid() {
			output["gasLimit"] = gasLimitField.Interface()
		} else {
			output["gasLimit"] = nil
		}
	case "evmExtraArgsV2":
		gasLimitField := val.FieldByName("GasLimit")
		if gasLimitField.IsValid() {
			output["gasLimit"] = gasLimitField.Interface()
		} else {
			output["gasLimit"] = nil
		}

		allowOutOfOrderField := val.FieldByName("AllowOutOfOrderExecution")
		if allowOutOfOrderField.IsValid() {
			output["allowOutOfOrderExecution"] = allowOutOfOrderField.Interface()
		} else {
			output["allowOutOfOrderExecution"] = false
		}
	case "svmExtraArgsV1":
		computeUnitsField := val.FieldByName("ComputeUnits")
		if computeUnitsField.IsValid() {
			output["computeUnits"] = computeUnitsField.Interface()
		} else {
			output["computeUnits"] = nil
		}

		bitmapField := val.FieldByName("AccountIsWritableBitmap")
		if bitmapField.IsValid() {
			output["accountIsWritableBitmap"] = bitmapField.Interface()
		} else {
			output["accountIsWritableBitmap"] = nil
		}

		allowOutOfOrderField := val.FieldByName("AllowOutOfOrderExecution")
		if allowOutOfOrderField.IsValid() {
			output["allowOutOfOrderExecution"] = allowOutOfOrderField.Interface()
		} else {
			output["allowOutOfOrderExecution"] = false
		}

		tokenReceiverField := val.FieldByName("TokenReceiver")
		if tokenReceiverField.IsValid() {
			output["tokenReceiver"] = tokenReceiverField.Interface()
		} else {
			output["tokenReceiver"] = nil
		}

		accountsField := val.FieldByName("Accounts")
		if accountsField.IsValid() {
			output["accounts"] = accountsField.Interface()
		} else {
			output["accounts"] = nil
		}
	}
	return output, nil
}

var _ ccipcommon.ExtraDataDecoder = (*ExtraDataDecoder)(nil)
