package encoding

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

var (
	ErrEmptyResults = errors.New("empty results; cannot encode")
)

type reportEncoder struct {
	packer Packer
}

var _ ocr2keepers.Encoder = (*reportEncoder)(nil)

func NewReportEncoder(p Packer) ocr2keepers.Encoder {
	return &reportEncoder{
		packer: p,
	}
}

func (e reportEncoder) Encode(results ...ocr2keepers.CheckResult) ([]byte, error) {
	if len(results) == 0 {
		return nil, ErrEmptyResults
	}

	report := ac.IAutomationV21PlusCommonReport{
		FastGasWei:   big.NewInt(0),
		LinkNative:   big.NewInt(0),
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
			if result.FastGasWei != nil {
				report.FastGasWei = result.FastGasWei
			}
			if result.LinkNative != nil {
				report.LinkNative = result.LinkNative
			}
		}

		id := result.UpkeepID.BigInt()
		report.UpkeepIds[i] = id
		report.GasLimits[i] = big.NewInt(0).SetUint64(result.GasAllocated)

		triggerW := triggerWrapper{
			BlockNum:  uint32(result.Trigger.BlockNumber),
			BlockHash: result.Trigger.BlockHash,
		}
		switch core.GetUpkeepType(result.UpkeepID) {
		case types.LogTrigger:
			triggerW.TxHash = result.Trigger.LogTriggerExtension.TxHash
			triggerW.LogIndex = result.Trigger.LogTriggerExtension.Index
			triggerW.LogBlockHash = result.Trigger.LogTriggerExtension.BlockHash
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

	return e.packer.PackReport(report)
}

// Extract extracts a slice of reported upkeeps (upkeep id, trigger, and work id) from raw bytes. the plugin will call this function to accept/transmit reports.
func (e reportEncoder) Extract(raw []byte) ([]ocr2keepers.ReportedUpkeep, error) {
	report, err := e.packer.UnpackReport(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unpack report", err)
	}
	reportedUpkeeps := make([]ocr2keepers.ReportedUpkeep, len(report.UpkeepIds))
	for i, upkeepId := range report.UpkeepIds {
		triggerW, err := core.UnpackTrigger(upkeepId, report.Triggers[i])
		if err != nil {
			return nil, fmt.Errorf("%w: failed to unpack trigger", err)
		}
		id := &ocr2keepers.UpkeepIdentifier{}
		id.FromBigInt(upkeepId)

		trigger := ocr2keepers.NewTrigger(
			ocr2keepers.BlockNumber(triggerW.BlockNum),
			triggerW.BlockHash,
		)
		switch core.GetUpkeepType(*id) {
		case types.LogTrigger:
			trigger.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{}
			trigger.LogTriggerExtension.TxHash = triggerW.TxHash
			trigger.LogTriggerExtension.Index = triggerW.LogIndex
			trigger.LogTriggerExtension.BlockHash = triggerW.LogBlockHash
		default:
		}
		workID := core.UpkeepWorkID(*id, trigger)
		reportedUpkeeps[i] = ocr2keepers.ReportedUpkeep{
			WorkID:   workID,
			UpkeepID: *id,
			Trigger:  trigger,
		}
	}

	return reportedUpkeeps, nil
}
