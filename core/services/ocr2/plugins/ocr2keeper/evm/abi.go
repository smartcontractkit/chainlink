package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

type RegistryVersion string

const (
	KeeperRegistryV20 RegistryVersion = "KeeperRegistry 2.0"
	KeeperRegistryV21 RegistryVersion = "KeeperRegistry 2.1"
)

type evmRegistryPacker struct {
	version RegistryVersion
	abi     abi.ABI
}

// enum UpkeepFailureReason
// https://github.com/smartcontractkit/chainlink/blob/d9dee8ea6af26bc82463510cb8786b951fa98585/contracts/src/v0.8/interfaces/AutomationRegistryInterface2_0.sol#L94
const (
	UPKEEP_FAILURE_REASON_NONE = iota
	UPKEEP_FAILURE_REASON_UPKEEP_CANCELLED
	UPKEEP_FAILURE_REASON_UPKEEP_PAUSED
	UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED
	UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
	UPKEEP_FAILURE_REASON_PERFORM_DATA_EXCEEDS_LIMIT
	UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE
	UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED
	UPKEEP_FAILURE_REASON_MERCURY_ACCESS_NOT_ALLOWED
)

func NewEvmRegistryPacker(version RegistryVersion, abi abi.ABI) *evmRegistryPacker {
	return &evmRegistryPacker{version: version, abi: abi}
}

func (rp *evmRegistryPacker) UnpackCheckResult(key ocr2keepers.UpkeepKey, raw string) (EVMAutomationUpkeepResult20, error) {
	var result EVMAutomationUpkeepResult20

	b, err := hexutil.Decode(raw)
	if err != nil {
		return result, err
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return result, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, raw)
	}

	block, id, err := splitKey(key)
	if err != nil {
		return result, err
	}

	result = EVMAutomationUpkeepResult20{
		Block:    uint32(block.Uint64()),
		ID:       id,
		Eligible: true,
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	result.FailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	result.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	result.FastGasWei = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	result.LinkNative = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	if !upkeepNeeded {
		result.Eligible = false
	}
	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result.FailureReason == UPKEEP_FAILURE_REASON_NONE || (result.FailureReason == UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED && len(rawPerformData) > 0) {
		var ret0 = new(performDataWrapper)
		err = pdataABI.UnpackIntoInterface(ret0, "check", rawPerformData)
		if err != nil {
			return result, err
		}

		result.CheckBlockNumber = ret0.Result.CheckBlockNumber
		result.CheckBlockHash = ret0.Result.CheckBlockhash
		result.PerformData = ret0.Result.PerformData
	}

	// This is a default placeholder which is used since we do not get the execute gas
	// from checkUpkeep result. This field is overwritten later from the execute gas
	// we have for an upkeep in memory. TODO (AUTO-1482): Refactor this
	result.ExecuteGas = 5_000_000

	return result, nil
}

// UnpackMercuryLookupResult should only be called on v2_1 registry
func (rp *evmRegistryPacker) UnpackMercuryLookupResult(callbackResp []byte) (bool, []byte, uint8, *big.Int, error) {
	if rp.version != KeeperRegistryV21 {
		return false, nil, 0, nil, fmt.Errorf("registry version %s does not support mercury lookup", rp.version)
	}

	out, err := rp.abi.Methods["mercuryCallback"].Outputs.UnpackValues(callbackResp)
	if err != nil {
		return false, nil, 0, nil, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, hexutil.Encode(callbackResp))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return upkeepNeeded, rawPerformData, failureReason, gasUsed, nil
}

func (rp *evmRegistryPacker) UnpackSimulatePerformResult(raw string) (bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return false, err
	}

	out, err := rp.abi.Methods["simulatePerformUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return false, fmt.Errorf("%w: unpack simulatePerformUpkeep return: %s", err, raw)
	}

	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// UnpackUpkeepResult will include a Forwarder field in upkeepInfo struct
func (rp *evmRegistryPacker) UnpackUpkeepResult(id *big.Int, raw string) (activeUpkeep, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return activeUpkeep{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return activeUpkeep{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	au := activeUpkeep{ID: id}
	if rp.version == KeeperRegistryV20 {
		temp := *abi.ConvertType(out[0], new(keeper_registry_wrapper2_0.UpkeepInfo)).(*keeper_registry_wrapper2_0.UpkeepInfo)
		au.PerformGasLimit = temp.ExecuteGas
		au.CheckData = temp.CheckData
	} else if rp.version == KeeperRegistryV21 {
		temp := *abi.ConvertType(out[0], new(i_keeper_registry_master_wrapper_2_1.UpkeepInfo)).(*i_keeper_registry_master_wrapper_2_1.UpkeepInfo)
		au.PerformGasLimit = temp.ExecuteGas
		au.CheckData = temp.CheckData
	} else {
		return activeUpkeep{}, fmt.Errorf("failed to unpack upkeep result for registry version %s", rp.version)
	}

	return au, nil
}

func (rp *evmRegistryPacker) UnpackTransmitTxInput(raw []byte) ([]ocr2keepers.UpkeepResult, error) {
	var (
		enc     = EVMAutomationEncoder20{}
		decoded []ocr2keepers.UpkeepResult
		out     []interface{}
		err     error
		b       []byte
		ok      bool
	)

	if out, err = rp.abi.Methods["transmit"].Inputs.UnpackValues(raw); err != nil {
		return nil, fmt.Errorf("%w: unpack TransmitTxInput return: %s", err, raw)
	}

	if len(out) < 2 {
		return nil, fmt.Errorf("invalid unpacking of TransmitTxInput in %s", raw)
	}

	if b, ok = out[1].([]byte); !ok {
		return nil, fmt.Errorf("unexpected value type in transaction")
	}

	if decoded, err = enc.DecodeReport(b); err != nil {
		return nil, fmt.Errorf("error during decoding report while unpacking TransmitTxInput: %w", err)
	}

	return decoded, nil
}

var (
	// rawPerformData is abi encoded tuple(uint32, bytes32, bytes). We create an ABI with dummy
	// function which returns this tuple in order to decode the bytes
	pdataABI, _ = abi.JSON(strings.NewReader(`[{
		"name":"check",
		"type":"function",
		"outputs":[{
			"name":"ret",
			"type":"tuple",
			"components":[
				{"type":"uint32","name":"checkBlockNumber"},
				{"type":"bytes32","name":"checkBlockhash"},
				{"type":"bytes","name":"performData"}
				]
			}]
		}]`,
	))
)

type performDataWrapper struct {
	Result performDataStruct
}
type performDataStruct struct {
	CheckBlockNumber uint32   `abi:"checkBlockNumber"`
	CheckBlockhash   [32]byte `abi:"checkBlockhash"`
	PerformData      []byte   `abi:"performData"`
}
