package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

// enum UpkeepFailureReason is defined by https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/dev/automation/2_1/interfaces/AutomationRegistryInterface2_1.sol#L97
// make sure failure reasons are in sync between contract and offchain enum
const (
	UPKEEP_FAILURE_REASON_NONE = iota
	UPKEEP_FAILURE_REASON_UPKEEP_CANCELLED
	UPKEEP_FAILURE_REASON_UPKEEP_PAUSED
	UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED
	UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
	UPKEEP_FAILURE_REASON_PERFORM_DATA_EXCEEDS_LIMIT
	UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE
	UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED
	UPKEEP_FAILURE_REASON_REVERT_DATA_EXCEEDS_LIMIT
	UPKEEP_FAILURE_REASON_REGISTRY_PAUSED

	// Start of offchain failure types. All onchain failure reasons from
	// contract should be put above
	UPKEEP_FAILURE_REASON_MERCURY_ACCESS_NOT_ALLOWED
)

var utilsABI = types.MustGetABI(automation_utils_2_1.AutomationUtilsABI)

type UpkeepInfo = iregistry21.KeeperRegistryBase21UpkeepInfo

// triggerWrapper is a wrapper for the different trigger types (log and condition triggers).
// NOTE: we use log trigger because it extends condition trigger,
type triggerWrapper = automation_utils_2_1.KeeperRegistryBase21LogTrigger

type evmRegistryPackerV2_1 struct {
	abi      abi.ABI
	utilsAbi abi.ABI
}

func NewEvmRegistryPackerV2_1(abi abi.ABI, utilsAbi abi.ABI) *evmRegistryPackerV2_1 {
	return &evmRegistryPackerV2_1{abi: abi, utilsAbi: utilsAbi}
}

func (rp *evmRegistryPackerV2_1) UnpackCheckResult(key ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error) {
	var result ocr2keepers.CheckResult

	b, err := hexutil.Decode(raw)
	if err != nil {
		return result, err
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return result, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, raw)
	}

	result = ocr2keepers.CheckResult{
		Eligible:     *abi.ConvertType(out[0], new(bool)).(*bool),
		Retryable:    false,
		GasAllocated: uint64((*abi.ConvertType(out[4], new(*big.Int)).(**big.Int)).Int64()),
		Payload:      key,
	}
	ext := EVMAutomationResultExtension21{
		FastGasWei:    *abi.ConvertType(out[5], new(*big.Int)).(**big.Int),
		LinkNative:    *abi.ConvertType(out[6], new(*big.Int)).(**big.Int),
		FailureReason: *abi.ConvertType(out[2], new(uint8)).(*uint8),
	}
	result.Extension = ext
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if ext.FailureReason == UPKEEP_FAILURE_REASON_NONE || (ext.FailureReason == UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED && len(rawPerformData) > 0) {
		result.PerformData = rawPerformData
	}

	return result, nil
}

func (rp *evmRegistryPackerV2_1) UnpackCheckCallbackResult(callbackResp []byte) (bool, []byte, uint8, *big.Int, error) {
	out, err := rp.abi.Methods["checkCallback"].Outputs.UnpackValues(callbackResp)
	if err != nil {
		return false, nil, 0, nil, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, hexutil.Encode(callbackResp))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return upkeepNeeded, rawPerformData, failureReason, gasUsed, nil
}

func (rp *evmRegistryPackerV2_1) UnpackPerformResult(raw string) (bool, error) {
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

func (rp *evmRegistryPackerV2_1) UnpackUpkeepInfo(id *big.Int, raw string) (UpkeepInfo, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return UpkeepInfo{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return UpkeepInfo{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	info := *abi.ConvertType(out[0], new(UpkeepInfo)).(*UpkeepInfo)

	return info, nil
}

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
func (rp *evmRegistryPackerV2_1) UnpackLogTriggerConfig(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
	var cfg automation_utils_2_1.LogTriggerConfig

	out, err := utilsABI.Methods["_logTriggerConfig"].Inputs.UnpackValues(raw)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack _logTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(automation_utils_2_1.LogTriggerConfig)).(*automation_utils_2_1.LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type")
	}
	return *converted, nil
}

// PackTrigger packs the trigger into the format expected by the contract,
// according to the upkeep type of the given id.
func (rp *evmRegistryPackerV2_1) PackTrigger(id *big.Int, trig triggerWrapper) ([]byte, error) {
	var trigger []byte
	var err error
	upkeepType := getUpkeepType(id.Bytes())
	switch upkeepType {
	case conditionTrigger:
		trig := automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
		}
		trigger, err = rp.utilsAbi.Pack("_conditionalTrigger", &trig)
	case logTrigger:
		logTrig := automation_utils_2_1.KeeperRegistryBase21LogTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
			LogIndex:  trig.LogIndex,
			TxHash:    trig.TxHash,
		}
		trigger, err = rp.utilsAbi.Pack("_logTrigger", &logTrig)
	default:
		err = fmt.Errorf("unknown trigger type: %d", upkeepType)
	}
	if err != nil {
		return nil, err
	}
	return trigger[4:], nil
}

// UnpackTrigger unpacks the trigger from the given raw data, according to the upkeep type of the given id.
func (rp *evmRegistryPackerV2_1) UnpackTrigger(id *big.Int, raw []byte) (triggerWrapper, error) {
	upkeepType := getUpkeepType(id.Bytes())
	switch upkeepType {
	case conditionTrigger:
		unpacked, err := rp.utilsAbi.Methods["_conditionalTrigger"].Inputs.Unpack(raw)
		if err != nil {
			return triggerWrapper{}, fmt.Errorf("%w: failed to unpack conditional trigger", err)
		}
		converted, ok := abi.ConvertType(unpacked[0], new(automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger)).(*automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger)
		if !ok {
			return triggerWrapper{}, fmt.Errorf("failed to convert type")
		}
		triggerW := triggerWrapper{
			BlockNum: converted.BlockNum,
		}
		copy(triggerW.BlockHash[:], converted.BlockHash[:])
		return triggerW, nil
	case logTrigger:
		unpacked, err := rp.utilsAbi.Methods["_logTrigger"].Inputs.Unpack(raw)
		if err != nil {
			return triggerWrapper{}, fmt.Errorf("%w: failed to unpack log trigger", err)
		}
		converted, ok := abi.ConvertType(unpacked[0], new(automation_utils_2_1.KeeperRegistryBase21LogTrigger)).(*automation_utils_2_1.KeeperRegistryBase21LogTrigger)
		if !ok {
			return triggerWrapper{}, fmt.Errorf("failed to convert type")
		}
		triggerW := triggerWrapper{
			BlockNum: converted.BlockNum,
			LogIndex: converted.LogIndex,
		}
		copy(triggerW.BlockHash[:], converted.BlockHash[:])
		copy(triggerW.TxHash[:], converted.TxHash[:])
		return triggerW, nil
	default:
		return triggerWrapper{}, fmt.Errorf("unknown trigger type: %d", upkeepType)
	}
}

// PackReport packs the report with abi definitions from the contract.
func (rp *evmRegistryPackerV2_1) PackReport(report automation_utils_2_1.KeeperRegistryBase21Report) ([]byte, error) {
	bts, err := rp.utilsAbi.Pack("_report", &report)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to pack report", err)
	}

	return bts[4:], nil
}

// UnpackReport unpacks the report from the given raw data.
func (rp *evmRegistryPackerV2_1) UnpackReport(raw []byte) (automation_utils_2_1.KeeperRegistryBase21Report, error) {
	unpacked, err := rp.utilsAbi.Methods["_report"].Inputs.Unpack(raw)
	if err != nil {
		return automation_utils_2_1.KeeperRegistryBase21Report{}, fmt.Errorf("%w: failed to unpack report", err)
	}
	converted, ok := abi.ConvertType(unpacked[0], new(automation_utils_2_1.KeeperRegistryBase21Report)).(*automation_utils_2_1.KeeperRegistryBase21Report)
	if !ok {
		return automation_utils_2_1.KeeperRegistryBase21Report{}, fmt.Errorf("failed to convert type")
	}
	report := automation_utils_2_1.KeeperRegistryBase21Report{
		FastGasWei:   converted.FastGasWei,
		LinkNative:   converted.LinkNative,
		UpkeepIds:    make([]*big.Int, len(converted.UpkeepIds)),
		GasLimits:    make([]*big.Int, len(converted.GasLimits)),
		Triggers:     make([][]byte, len(converted.Triggers)),
		PerformDatas: make([][]byte, len(converted.PerformDatas)),
	}
	if len(report.UpkeepIds) > 0 {
		copy(report.UpkeepIds, converted.UpkeepIds)
		copy(report.GasLimits, converted.GasLimits)
		copy(report.Triggers, converted.Triggers)
		copy(report.PerformDatas, converted.PerformDatas)
	}

	return report, nil
}
