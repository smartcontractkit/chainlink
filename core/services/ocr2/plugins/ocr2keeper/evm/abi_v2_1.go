package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
)

type evmRegistryPackerV2_1 struct {
	abi abi.ABI
}

func NewEvmRegistryPackerV2_1(abi abi.ABI) *evmRegistryPackerV2_1 {
	return &evmRegistryPackerV2_1{abi: abi}
}

// TODO: adjust to 2.1 if needed
func (rp *evmRegistryPackerV2_1) UnpackCheckResult(key types.UpkeepKey, raw string) (types.UpkeepResult, error) {
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

// TODO: adjust to 2.1 if needed
func (rp *evmRegistryPackerV2_1) UnpackMercuryLookupResult(callbackResp []byte) (bool, []byte, error) {
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

func (rp *evmRegistryPackerV2_1) UnpackPerformResult(raw string) (bool, error) {
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

func (rp *evmRegistryPackerV2_1) UnpackUpkeepConfig(raw string) (keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig, error) {
	var cfg keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig
	b, err := hexutil.Decode(raw)
	if err != nil {
		return cfg, err
	}

	out, err := rp.abi.Methods["getLogTriggerConfig"].Outputs.UnpackValues(b)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig)).(*keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type")
	}
	return *converted, nil
}

func (rp *evmRegistryPackerV2_1) UnpackUpkeepInfo(id *big.Int, raw string) (upkeepInfoEntry, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return upkeepInfoEntry{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return upkeepInfoEntry{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	temp := *abi.ConvertType(out[0], new(keeper_registry_logic_b_wrapper_2_1.UpkeepInfo)).(*keeper_registry_logic_b_wrapper_2_1.UpkeepInfo)

	u := upkeepInfoEntry{
		id:              id,
		target:          temp.Target,
		performGasLimit: temp.ExecuteGas,
		offchainConfig:  temp.OffchainConfig,
	}
	if temp.Paused {
		u.state = stateInactive
	}

	return u, nil
}

// TODO: adjust to 2.1 if needed
func (rp *evmRegistryPackerV2_1) UnpackTransmitTxInput(raw []byte) ([]types.UpkeepResult, error) {
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
