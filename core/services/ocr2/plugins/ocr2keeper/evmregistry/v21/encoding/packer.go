package encoding

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

// triggerWrapper is a wrapper for the different trigger types (log and condition triggers).
// NOTE: we use log trigger because it extends condition trigger,
type triggerWrapper = ac.IAutomationV21PlusCommonLogTrigger

type abiPacker struct {
	autoV2CommonABI abi.ABI
	utilsABI        abi.ABI
	streamsABI      abi.ABI
}

var _ Packer = (*abiPacker)(nil)

func NewAbiPacker() *abiPacker {
	return &abiPacker{autoV2CommonABI: core.AutoV2CommonABI, utilsABI: core.CompatibleUtilsABI, streamsABI: core.StreamsCompatibleABI}
}

func (p *abiPacker) UnpackCheckResult(payload ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		// decode failed, not retryable
		return GetIneligibleCheckResultWithoutPerformData(payload, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false),
			fmt.Errorf("upkeepId %s failed to decode checkUpkeep result %s: %s", payload.UpkeepID.String(), raw, err)
	}

	out, err := p.autoV2CommonABI.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		// unpack failed, not retryable
		return GetIneligibleCheckResultWithoutPerformData(payload, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false),
			fmt.Errorf("upkeepId %s failed to unpack checkUpkeep result %s: %s", payload.UpkeepID.String(), raw, err)
	}

	result := ocr2keepers.CheckResult{
		Eligible:            *abi.ConvertType(out[0], new(bool)).(*bool),
		Retryable:           false,
		GasAllocated:        uint64((*abi.ConvertType(out[4], new(*big.Int)).(**big.Int)).Int64()),
		UpkeepID:            payload.UpkeepID,
		Trigger:             payload.Trigger,
		WorkID:              payload.WorkID,
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

func (p *abiPacker) UnpackPerformResult(raw string) (PipelineExecutionState, bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	out, err := p.autoV2CommonABI.Methods["simulatePerformUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	return NoPipelineError, *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
func (p *abiPacker) UnpackLogTriggerConfig(raw []byte) (ac.IAutomationV21PlusCommonLogTriggerConfig, error) {
	var cfg ac.IAutomationV21PlusCommonLogTriggerConfig

	out, err := core.CompatibleUtilsABI.Methods["_logTriggerConfig"].Inputs.UnpackValues(raw)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack _logTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(ac.IAutomationV21PlusCommonLogTriggerConfig)).(*ac.IAutomationV21PlusCommonLogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type during UnpackLogTriggerConfig")
	}
	return *converted, nil
}

// PackReport packs the report with abi definitions from the contract.
func (p *abiPacker) PackReport(report ac.IAutomationV21PlusCommonReport) ([]byte, error) {
	bts, err := p.utilsABI.Methods["_report"].Inputs.Pack(&report)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to pack report", err)
	}
	return bts, nil
}

// UnpackReport unpacks the report from the given raw data.
func (p *abiPacker) UnpackReport(raw []byte) (ac.IAutomationV21PlusCommonReport, error) {
	unpacked, err := p.utilsABI.Methods["_report"].Inputs.Unpack(raw)
	if err != nil {
		return ac.IAutomationV21PlusCommonReport{}, fmt.Errorf("%w: failed to unpack report", err)
	}
	converted, ok := abi.ConvertType(unpacked[0], new(ac.IAutomationV21PlusCommonReport)).(*ac.IAutomationV21PlusCommonReport)
	if !ok {
		return ac.IAutomationV21PlusCommonReport{}, fmt.Errorf("failed to convert type")
	}
	report := ac.IAutomationV21PlusCommonReport{
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

// GetIneligibleCheckResultWithoutPerformData returns an ineligible check result with ineligibility reason and pipeline execution state but without perform data
func GetIneligibleCheckResultWithoutPerformData(p ocr2keepers.UpkeepPayload, reason UpkeepFailureReason, state PipelineExecutionState, retryable bool) ocr2keepers.CheckResult {
	return ocr2keepers.CheckResult{
		IneligibilityReason:    uint8(reason),
		PipelineExecutionState: uint8(state),
		Retryable:              retryable,
		UpkeepID:               p.UpkeepID,
		Trigger:                p.Trigger,
		WorkID:                 p.WorkID,
		FastGasWei:             big.NewInt(0),
		LinkNative:             big.NewInt(0),
	}
}
