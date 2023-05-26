package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

type evmRegistryPackerV2_0 struct {
	abi abi.ABI
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
)

func NewEvmRegistryPackerV2_0(abi abi.ABI) *evmRegistryPackerV2_0 {
	return &evmRegistryPackerV2_0{abi: abi}
}

func (rp *evmRegistryPackerV2_0) UnpackCheckResult(key types.UpkeepKey, raw string) (types.UpkeepResult, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return types.UpkeepResult{}, err
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return types.UpkeepResult{}, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, raw)
	}

	result := types.UpkeepResult{
		Key:   key,
		State: types.Eligible,
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	result.FailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	result.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	result.FastGasWei = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	result.LinkNative = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	if !upkeepNeeded {
		result.State = types.NotEligible
	}
	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result.FailureReason == UPKEEP_FAILURE_REASON_NONE || (result.FailureReason == UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED && len(rawPerformData) > 0) {
		var ret0 = new(performDataWrapper)
		err = pdataABI.UnpackIntoInterface(ret0, "check", rawPerformData)
		if err != nil {
			return types.UpkeepResult{}, err
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

func (rp *evmRegistryPackerV2_0) UnpackMercuryLookupResult(callbackResp []byte) (bool, []byte, error) {
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return false, nil, fmt.Errorf("abi new bytes type error: %w", err)
	}
	boolTyp, err := abi.NewType("bool", "", nil)
	if err != nil {
		return false, nil, fmt.Errorf("abi new bool type error: %w", err)
	}
	callbackOutput := abi.Arguments{
		{Name: "upkeepNeeded", Type: boolTyp},
		{Name: "performData", Type: typBytes},
	}
	unpack, err := callbackOutput.Unpack(callbackResp)
	if err != nil {
		return false, nil, fmt.Errorf("callback output unpack error: %w", err)
	}

	upkeepNeeded := *abi.ConvertType(unpack[0], new(bool)).(*bool)
	if !upkeepNeeded {
		return false, nil, nil
	}
	performData := *abi.ConvertType(unpack[1], new([]byte)).(*[]byte)
	return true, performData, nil
}

func (rp *evmRegistryPackerV2_0) UnpackPerformResult(raw string) (bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return false, err
	}

	out, err := rp.abi.Methods["simulatePerformUpkeep"].
		Outputs.UnpackValues(b)
	if err != nil {
		return false, fmt.Errorf("%w: unpack simulatePerformUpkeep return: %s", err, raw)
	}

	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

func (rp *evmRegistryPackerV2_0) UnpackUpkeepResult(id *big.Int, raw string) (activeUpkeep, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return activeUpkeep{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return activeUpkeep{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	type upkeepInfo struct {
		Target                 common.Address
		ExecuteGas             uint32
		CheckData              []byte
		Balance                *big.Int
		Admin                  common.Address
		MaxValidBlocknumber    uint64
		LastPerformBlockNumber uint32
		AmountSpent            *big.Int
		Paused                 bool
		OffchainConfig         []byte
	}
	temp := *abi.ConvertType(out[0], new(upkeepInfo)).(*upkeepInfo)

	au := activeUpkeep{
		ID:              id,
		PerformGasLimit: temp.ExecuteGas,
		CheckData:       temp.CheckData,
	}

	return au, nil
}

func (rp *evmRegistryPackerV2_0) UnpackTransmitTxInput(raw []byte) ([]types.UpkeepResult, error) {
	out, err := rp.abi.Methods["transmit"].Inputs.UnpackValues(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: unpack TransmitTxInput return: %s", err, raw)
	}

	if len(out) < 2 {
		return nil, fmt.Errorf("invalid unpacking of TransmitTxInput in %s", raw)
	}
	decodedReport, err := chain.NewEVMReportEncoder().DecodeReport(out[1].([]byte))
	if err != nil {
		return nil, fmt.Errorf("error during decoding report while unpacking TransmitTxInput: %w", err)
	}
	return decodedReport, nil
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
