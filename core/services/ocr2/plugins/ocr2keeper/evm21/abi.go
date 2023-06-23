package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
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

type UpkeepInfo = iregistry21.KeeperRegistryBase21UpkeepInfo

type evmRegistryPackerV2_1 struct {
	abi        abi.ABI
	logDataABI abi.ABI
}

func NewEvmRegistryPackerV2_1(abi abi.ABI, logDataABI abi.ABI) *evmRegistryPackerV2_1 {
	return &evmRegistryPackerV2_1{abi: abi, logDataABI: logDataABI}
}

func (rp *evmRegistryPackerV2_1) PackLogData(log logpoller.Log) ([]byte, error) {
	topics := [][32]byte{}
	for _, topic := range log.GetTopics() {
		topics = append(topics, topic)
	}
	return rp.logDataABI.Pack("checkLog", &i_log_automation.Log{
		Index:       big.NewInt(log.LogIndex),
		TxIndex:     big.NewInt(0), // TODO
		TxHash:      log.TxHash,
		BlockNumber: big.NewInt(log.BlockNumber),
		BlockHash:   log.BlockHash,
		Source:      log.Address,
		Topics:      topics,
		Data:        log.Data,
	})
}

// TODO: remove for 2.1
func (rp *evmRegistryPackerV2_1) UnpackCheckResult(key ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error) {
	var result ocr2keepers.CheckResult

	fmt.Printf("[FeedLookup] payload: %v", key)

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
		GasAllocated: 5000000, // use upkeep gas limit/execute gas
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

	fmt.Printf("[FeedLookup] rawPerformData: %s", hexutil.Encode(rawPerformData))

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

func (rp *evmRegistryPackerV2_1) UnpackTransmitTxInput(raw []byte) ([]ocr2keepers.UpkeepResult, error) {
	var (
		enc     = EVMAutomationEncoder21{}
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

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
func (rp *evmRegistryPackerV2_1) UnpackLogTriggerConfig(raw []byte) (iregistry21.KeeperRegistryBase21LogTriggerConfig, error) {
	var cfg iregistry21.KeeperRegistryBase21LogTriggerConfig

	out, err := rp.abi.Methods["getLogTriggerConfig"].Outputs.UnpackValues(raw)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack getLogTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(iregistry21.KeeperRegistryBase21LogTriggerConfig)).(*iregistry21.KeeperRegistryBase21LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type")
	}
	return *converted, nil
}
