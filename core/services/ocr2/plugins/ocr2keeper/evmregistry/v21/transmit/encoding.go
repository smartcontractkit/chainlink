package transmit

import (
	"fmt"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

// defaultLogParser parses logs from the registry contract
func defaultLogParser(registry *iregistry21.IKeeperRegistryMaster, log logpoller.Log) (transmitEventLog, error) {
	rawLog := log.ToGethLog()
	abilog, err := registry.ParseLog(rawLog)
	if err != nil {
		return transmitEventLog{}, fmt.Errorf("%w: failed to parse log", err)
	}

	switch l := abilog.(type) {
	case *iregistry21.IKeeperRegistryMasterUpkeepPerformed:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:       log,
			Performed: l,
		}, nil
	case *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:     log,
			Reorged: l,
		}, nil
	case *iregistry21.IKeeperRegistryMasterStaleUpkeepReport:
		if l == nil {
			break
		}
		return transmitEventLog{
			Log:   log,
			Stale: l,
		}, nil
	case *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport:
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
	Performed         *iregistry21.IKeeperRegistryMasterUpkeepPerformed
	Stale             *iregistry21.IKeeperRegistryMasterStaleUpkeepReport
	Reorged           *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport
	InsufficientFunds *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport
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
