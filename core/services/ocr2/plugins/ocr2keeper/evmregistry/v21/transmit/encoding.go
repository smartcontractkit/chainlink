package transmit

import (
	"fmt"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	autov2common "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v2_common"
)

// defaultLogParser parses logs from the registry contract
func defaultLogParser(registry *autov2common.IAutomationV2Common, log logpoller.Log) (transmitEventLog, error) {
	rawLog := log.ToGethLog()
	abilog, err := registry.ParseLog(rawLog)
	if err != nil {
		return transmitEventLog{}, fmt.Errorf("%w: failed to parse log", err)
	}

	switch l := abilog.(type) {
	case *autov2common.IAutomationV2CommonUpkeepPerformed:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:       log,
			Performed: l,
		}, nil
	case *autov2common.IAutomationV2CommonReorgedUpkeepReport:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:     log,
			Reorged: l,
		}, nil
	case *autov2common.IAutomationV2CommonStaleUpkeepReport:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:   log,
			Stale: l,
		}, nil
	case *autov2common.IAutomationV2CommonInsufficientFundsUpkeepReport:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:               log,
			InsufficientFunds: l,
		}, nil
	default:
		return transmitEventLog{}, fmt.Errorf("unknown log type: %v", l)
	}
	return transmitEventLog{}, fmt.Errorf("log with bad structure")
}

// transmitEventLog is a wrapper around logpoller.Log and the parsed log
type transmitEventLog struct {
	logpoller.Log
	Performed         *autov2common.IAutomationV2CommonUpkeepPerformed
	Stale             *autov2common.IAutomationV2CommonStaleUpkeepReport
	Reorged           *autov2common.IAutomationV2CommonReorgedUpkeepReport
	InsufficientFunds *autov2common.IAutomationV2CommonInsufficientFundsUpkeepReport
}

func (l transmitEventLog) Id() *big.Int {
	switch {
	case l.Performed != nil:
		return l.Performed.Id
	case l.Stale != nil:
		return l.Stale.Id
	case l.Reorged != nil:
		return l.Reorged.Id
	case l.InsufficientFunds != nil:
		return l.InsufficientFunds.Id
	default:
		return nil
	}
}

func (l transmitEventLog) Trigger() []byte {
	switch {
	case l.Performed != nil:
		return l.Performed.Trigger
	case l.Stale != nil:
		return l.Stale.Trigger
	case l.Reorged != nil:
		return l.Reorged.Trigger
	case l.InsufficientFunds != nil:
		return l.InsufficientFunds.Trigger
	default:
		return []byte{}
	}
}

func (l transmitEventLog) TransmitEventType() ocr2keepers.TransmitEventType {
	switch {
	case l.Performed != nil:
		return ocr2keepers.PerformEvent
	case l.Stale != nil:
		return ocr2keepers.StaleReportEvent
	case l.Reorged != nil:
		return ocr2keepers.ReorgReportEvent
	case l.InsufficientFunds != nil:
		return ocr2keepers.InsufficientFundsReportEvent
	default:
		return ocr2keepers.UnknownEvent
	}
}
