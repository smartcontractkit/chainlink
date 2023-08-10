package evm

import (
	"fmt"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

var (
	ErrEmptyResults = fmt.Errorf("empty results; cannot encode")
)

type EVMAutomationEncoder21 struct {
	packer *evmRegistryPackerV2_1
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

	encoded := 0
	highestCheckBlock := big.NewInt(0)

	for i, result := range results {
		checkBlock := big.NewInt(int64(result.Trigger.BlockNumber))

		if checkBlock.Cmp(highestCheckBlock) == 1 {
			highestCheckBlock = checkBlock
			report.FastGasWei = result.FastGasWei
			report.LinkNative = result.LinkNative
		}

		id := result.UpkeepID.BigInt()
		report.UpkeepIds[i] = id
		report.GasLimits[i] = big.NewInt(0).SetUint64(result.GasAllocated)

		triggerW := triggerWrapper{
			BlockNum:  uint32(result.Trigger.BlockNumber),
			BlockHash: result.Trigger.BlockHash,
		}
		switch core.GetUpkeepType(result.UpkeepID) {
		case ocr2keepers.LogTrigger:
			triggerW.TxHash = result.Trigger.LogTriggerExtension.TxHash
			triggerW.LogIndex = result.Trigger.LogTriggerExtension.Index
		default:
			// no special handling here for conditional triggers
		}

		trigger, err := core.PackTrigger(id, triggerW)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to pack trigger", err)
		}

		report.Triggers[i] = trigger
		report.PerformDatas[i] = result.PerformData

		encoded++
	}

	fmt.Printf("[automation-ocr3|EvmRegistry|Encoder] encoded %d out of %d results\n", encoded, len(results))

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
		triggerW, err := core.UnpackTrigger(upkeepId, report.Triggers[i])
		if err != nil {
			// TODO: log error and continue instead?
			return nil, fmt.Errorf("%w: failed to unpack trigger", err)
		}
		id := &ocr2keepers.UpkeepIdentifier{}
		id.FromBigInt(upkeepId)

		trigger := ocr2keepers.NewTrigger(
			ocr2keepers.BlockNumber(triggerW.BlockNum),
			triggerW.BlockHash,
		)
		switch core.GetUpkeepType(*id) {
		case ocr2keepers.LogTrigger:
			trigger.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{}
			trigger.LogTriggerExtension.TxHash = triggerW.TxHash
			trigger.LogTriggerExtension.Index = triggerW.LogIndex
		default:
		}
		workID, _ := core.UpkeepWorkID(upkeepId, trigger)
		reportedUpkeeps[i] = ocr2keepers.ReportedUpkeep{
			WorkID:   workID,
			UpkeepID: *id,
			Trigger:  trigger,
		}
	}

	return reportedUpkeeps, nil
}
