package encoding

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

// triggerWrapper is a wrapper for the different trigger types (log and condition triggers).
// NOTE: we use log trigger because it extends condition trigger,
type triggerWrapper = automation_utils_2_1.KeeperRegistryBase21LogTrigger

type abiPacker struct {
	registryABI abi.ABI
	utilsABI    abi.ABI
	streamsABI  abi.ABI
}

var _ Packer = (*abiPacker)(nil)

func NewAbiPacker() *abiPacker {
	return &abiPacker{registryABI: core.RegistryABI, utilsABI: core.UtilsABI, streamsABI: core.StreamsCompatibleABI}
}

func (p *abiPacker) UnpackCheckResult(payload ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		// decode failed, not retryable
		return GetIneligibleCheckResultWithoutPerformData(payload, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false),
			fmt.Errorf("upkeepId %s failed to decode checkUpkeep result %s: %s", payload.UpkeepID.String(), raw, err)
	}

	out, err := p.registryABI.Methods["checkUpkeep"].Outputs.UnpackValues(b)
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

func (p *abiPacker) PackGetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return p.registryABI.Pack("getUpkeepPrivilegeConfig", upkeepId)
}

func (p *abiPacker) UnpackGetUpkeepPrivilegeConfig(resp []byte) ([]byte, error) {
	out, err := p.registryABI.Methods["getUpkeepPrivilegeConfig"].Outputs.UnpackValues(resp)
	if err != nil {
		return nil, fmt.Errorf("%w: unpack getUpkeepPrivilegeConfig return", err)
	}

	bts := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return bts, nil
}

func (p *abiPacker) UnpackCheckCallbackResult(callbackResp []byte) (PipelineExecutionState, bool, []byte, uint8, *big.Int, error) {
	out, err := p.registryABI.Methods["checkCallback"].Outputs.UnpackValues(callbackResp)
	if err != nil {
		return PackUnpackDecodeFailed, false, nil, 0, nil, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, hexutil.Encode(callbackResp))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return NoPipelineError, upkeepNeeded, rawPerformData, failureReason, gasUsed, nil
}

func (p *abiPacker) UnpackPerformResult(raw string) (PipelineExecutionState, bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	out, err := p.registryABI.Methods["simulatePerformUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return PackUnpackDecodeFailed, false, err
	}

	return NoPipelineError, *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
func (p *abiPacker) UnpackLogTriggerConfig(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
	var cfg automation_utils_2_1.LogTriggerConfig

	out, err := core.UtilsABI.Methods["_logTriggerConfig"].Inputs.UnpackValues(raw)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack _logTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(automation_utils_2_1.LogTriggerConfig)).(*automation_utils_2_1.LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type during UnpackLogTriggerConfig")
	}
	return *converted, nil
}

// PackReport packs the report with abi definitions from the contract.
func (p *abiPacker) PackReport(report automation_utils_2_1.KeeperRegistryBase21Report) ([]byte, error) {
	bts, err := p.utilsABI.Methods["_report"].Inputs.Pack(&report)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to pack report", err)
	}
	return bts, nil
}

// UnpackReport unpacks the report from the given raw data.
func (p *abiPacker) UnpackReport(raw []byte) (automation_utils_2_1.KeeperRegistryBase21Report, error) {
	unpacked, err := p.utilsABI.Methods["_report"].Inputs.Unpack(raw)
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

type StreamsLookupError struct {
	FeedParamKey string
	Feeds        []string
	TimeParamKey string
	Time         *big.Int
	ExtraData    []byte
}

// DecodeStreamsLookupRequest decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (p *abiPacker) DecodeStreamsLookupRequest(data []byte) (*StreamsLookupError, error) {
	e := p.streamsABI.Errors["StreamsLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &StreamsLookupError{
		FeedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		Feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		TimeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		Time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		ExtraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
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
