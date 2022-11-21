package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
)

func (rs *RegistrySynchronizer) JobID() int32 {
	return rs.job.ID
}

func (rs *RegistrySynchronizer) HandleLog(broadcast log.Broadcast) {
	eventLog := broadcast.DecodedLog()
	if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
		rs.logger.Panicf("HandleLog: ignoring nil value, type: %T", broadcast)
		return
	}

	svcLogger := rs.logger.With(
		"logType", reflect.TypeOf(eventLog),
		"txHash", broadcast.RawLog().TxHash.Hex(),
	)

	svcLogger.Debug("received log, waiting for confirmations")

	wasOverCapacity := rs.mbLogs.Deliver(broadcast)
	if wasOverCapacity {
		svcLogger.Errorf("mailbox is over capacity - dropped the oldest unprocessed item")
	}
}
