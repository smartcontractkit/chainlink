package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/encoding"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
)

var (
	ErrEmptyResults = fmt.Errorf("empty results; cannot encode")
)

type EVMAutomationEncoder21 struct {
	encoding.BasicEncoder
	packer *evmRegistryPackerV2_1
}

type EVMAutomationUpkeepResult21 struct {
	// Block is the block number used to build an UpkeepKey for this result
	Block uint32
	// ID is the unique identifier for the upkeep
	ID            *big.Int
	Eligible      bool
	FailureReason uint8
	GasUsed       *big.Int
	PerformData   []byte
	FastGasWei    *big.Int
	LinkNative    *big.Int
	// CheckBlockNumber is the block number that the contract indicates the
	// upkeep was checked on
	CheckBlockNumber uint32
	CheckBlockHash   [32]byte
	ExecuteGas       uint32
	Retryable        bool
}

type EVMAutomationResultExtension21 struct {
	FastGasWei    *big.Int
	LinkNative    *big.Int
	FailureReason uint8 // this is not encoded, only pass along for the purpose of pipeline run
}

func (enc EVMAutomationEncoder21) Encode(results ...ocr2keepers.CheckResult) ([]byte, error) {
	if len(results) == 0 {
		return nil, ErrEmptyResults
	}

	report := automation_utils_2_1.KeeperRegistryBase21Report{
		UpkeepIds:    make([]*big.Int, len(results)),
		GasLimits:    make([]*big.Int, len(results)),
		Triggers:     make([][]byte, len(results)),
		PerformDatas: make([][]byte, len(results)),
	}

	for i, result := range results {
		ext, ok := result.Extension.(EVMAutomationResultExtension21)
		if !ok {
			return nil, fmt.Errorf("unexpected check result extension struct")
		}

		// only take these values from the first result
		// TODO: find a new way to get these values
		if i == 0 {
			report.FastGasWei = ext.FastGasWei
			report.LinkNative = ext.LinkNative
		}

		id, ok := big.NewInt(0).SetString(string(result.Payload.Upkeep.ID), 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse big int from upkeep id: %s", string(result.Payload.Upkeep.ID))
		}
		report.UpkeepIds[i] = id
		report.GasLimits[i] = big.NewInt(0).SetUint64(result.GasAllocated)

		triggerW := triggerWrapper{
			BlockNum:  uint32(result.Payload.Trigger.BlockNumber),
			BlockHash: common.HexToHash(result.Payload.Trigger.BlockHash),
		}
		switch getUpkeepType(id.Bytes()) {
		case logTrigger:
			trExt, ok := result.Payload.Trigger.Extension.(logprovider.LogTriggerExtension)
			if !ok {
				return nil, fmt.Errorf("unrecognized trigger extension data")
			}
			hex, err := common.ParseHexOrString(trExt.TxHash)
			if err != nil {
				return nil, fmt.Errorf("tx hash parse error: %w", err)
			}
			triggerW.TxHash = common.BytesToHash(hex[:])
			triggerW.LogIndex = uint32(trExt.LogIndex)
		default:
		}
		trigger, err := enc.packer.PackTrigger(id, triggerW)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to pack trigger", err)
		}
		report.Triggers[i] = trigger
		report.PerformDatas[i] = result.PerformData
	}

	return enc.packer.PackReport(report)
}

// Extract the plugin will call this function to accept/transmit reports
func (enc EVMAutomationEncoder21) Extract(raw []byte) ([]ocr2keepers.ReportedUpkeep, error) {
	report, err := enc.packer.UnpackReport(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unpack report", err)
	}
	reportedUpkeeps := make([]ocr2keepers.ReportedUpkeep, len(report.UpkeepIds))
	for i, upkeepId := range report.UpkeepIds {
		triggerW, err := enc.packer.UnpackTrigger(upkeepId, report.Triggers[i])
		if err != nil {
			// TODO: log error and continue instead?
			return nil, fmt.Errorf("%w: failed to unpack trigger", err)
		}
		logExt := logprovider.LogTriggerExtension{}

		switch getUpkeepType(upkeepId.Bytes()) {
		case logTrigger:
			logExt.TxHash = common.BytesToHash(triggerW.TxHash[:]).Hex()
			logExt.LogIndex = int64(triggerW.LogIndex)
		default:
		}
		trigger := ocr2keepers.NewTrigger(
			int64(triggerW.BlockNum),
			common.BytesToHash(triggerW.BlockHash[:]).Hex(),
			logExt,
		)
		triggerID, _ := UpkeepTriggerID(upkeepId, report.Triggers[i])
		reportedUpkeeps[i] = ocr2keepers.ReportedUpkeep{
			ID:          triggerID,
			UpkeepID:    ocr2keepers.UpkeepIdentifier(upkeepId.String()),
			Trigger:     trigger,
			PerformData: report.PerformDatas[i],
		}
	}

	return reportedUpkeeps, nil
}

type BlockKeyHelper[T uint32 | int64] struct {
}

func (kh BlockKeyHelper[T]) MakeBlockKey(b T) ocr2keepers.BlockKey {
	return ocr2keepers.BlockKey(fmt.Sprintf("%d", b))
}

type UpkeepKeyHelper[T uint32 | int64] struct {
}

func (kh UpkeepKeyHelper[T]) MakeUpkeepKey(b T, id *big.Int) ocr2keepers.UpkeepKey {
	return ocr2keepers.UpkeepKey(fmt.Sprintf("%d%s%s", b, separator, id))
}
