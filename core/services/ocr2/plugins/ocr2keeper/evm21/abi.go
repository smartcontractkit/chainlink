package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

type UpkeepFailureReason uint8
type PipelineExecutionState uint8

const (
	// upkeep failure onchain reasons
	UpkeepFailureReasonNone                    UpkeepFailureReason = 0
	UpkeepFailureReasonUpkeepCancelled         UpkeepFailureReason = 1
	UpkeepFailureReasonUpkeepPaused            UpkeepFailureReason = 2
	UpkeepFailureReasonTargetCheckReverted     UpkeepFailureReason = 3
	UpkeepFailureReasonUpkeepNotNeeded         UpkeepFailureReason = 4
	UpkeepFailureReasonPerformDataExceedsLimit UpkeepFailureReason = 5
	UpkeepFailureReasonInsufficientBalance     UpkeepFailureReason = 6
	UpkeepFailureReasonMercuryCallbackReverted UpkeepFailureReason = 7
	UpkeepFailureReasonRevertDataExceedsLimit  UpkeepFailureReason = 8
	UpkeepFailureReasonRegistryPaused          UpkeepFailureReason = 9
	// leaving a gap here for more onchain failure reasons in the future
	// upkeep failure offchain reasons
	UpkeepFailureReasonMercuryAccessNotAllowed UpkeepFailureReason = 32
	UpkeepFailureReasonLogBlockNoLongerExists  UpkeepFailureReason = 31
	UpkeepFailureReasonLogBlockInvalid         UpkeepFailureReason = 32
	UpkeepFailureReasonTxHashNoLongerExists    UpkeepFailureReason = 33

	// pipeline execution error
	NoPipelineError             PipelineExecutionState = 0
	CheckBlockTooOld            PipelineExecutionState = 1
	CheckBlockInvalid           PipelineExecutionState = 2
	RpcFlakyFailure             PipelineExecutionState = 3
	MercuryFlakyFailure         PipelineExecutionState = 4
	PackUnpackDecodeFailed      PipelineExecutionState = 5
	MercuryUnmarshalError       PipelineExecutionState = 6
	InvalidMercuryRequest       PipelineExecutionState = 7
	FailedToReadMercuryResponse PipelineExecutionState = 8
	InvalidRevertDataInput      PipelineExecutionState = 9
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

func (rp *evmRegistryPackerV2_1) UnpackCheckResult(p ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		// decode failed, not retryable
		return getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false), fmt.Errorf("upkeepId %s failed to decode checkUpkeep result %s: %s", p.UpkeepID.String(), raw, err)
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		// unpack failed, not retryable
		return getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false), fmt.Errorf("upkeepId %s failed to unpack checkUpkeep result %s: %s", p.UpkeepID.String(), raw, err)
	}

	result := ocr2keepers.CheckResult{
		Eligible:            *abi.ConvertType(out[0], new(bool)).(*bool),
		Retryable:           false,
		GasAllocated:        uint64((*abi.ConvertType(out[4], new(*big.Int)).(**big.Int)).Int64()),
		UpkeepID:            p.UpkeepID,
		Trigger:             p.Trigger,
		WorkID:              p.WorkID,
		FastGasWei:          *abi.ConvertType(out[5], new(*big.Int)).(**big.Int),
		LinkNative:          *abi.ConvertType(out[6], new(*big.Int)).(**big.Int),
		IneligibilityReason: *abi.ConvertType(out[2], new(uint8)).(*uint8),
	}

	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result.IneligibilityReason == uint8(UpkeepFailureReasonNone) || (result.IneligibilityReason == uint8(UpkeepFailureReasonTargetCheckReverted) && len(rawPerformData) > 0) {
		result.PerformData = rawPerformData
	}

	return result, nil
}

func (rp *evmRegistryPackerV2_1) UnpackCheckCallbackResult(callbackResp []byte) (PipelineExecutionState, bool, []byte, uint8, *big.Int, error) {
	out, err := rp.abi.Methods["checkCallback"].Outputs.UnpackValues(callbackResp)
	if err != nil {
		return PackUnpackDecodeFailed, false, nil, 0, nil, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, hexutil.Encode(callbackResp))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return NoPipelineError, upkeepNeeded, rawPerformData, failureReason, gasUsed, nil
}

func (rp *evmRegistryPackerV2_1) UnpackPerformResult(raw string) (PipelineExecutionState, bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	out, err := rp.abi.Methods["simulatePerformUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	return NoPipelineError, *abi.ConvertType(out[0], new(bool)).(*bool), nil
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
