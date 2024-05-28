package transmit

import (
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestTransmitEventLog(t *testing.T) {
	uid := core.GenUpkeepID(types.ConditionTrigger, "111")

	tests := []struct {
		name  string
		log   transmitEventLog
		etype ocr2keepers.TransmitEventType
	}{
		{
			"performed",
			transmitEventLog{
				Log: logpoller.Log{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x010203040"),
				},
				Performed: &ac.IAutomationV21PlusCommonUpkeepPerformed{
					Id:      uid.BigInt(),
					Trigger: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			ocr2keepers.PerformEvent,
		},
		{
			"stale",
			transmitEventLog{
				Log: logpoller.Log{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x010203040"),
				},
				Stale: &ac.IAutomationV21PlusCommonStaleUpkeepReport{
					Id:      uid.BigInt(),
					Trigger: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			ocr2keepers.StaleReportEvent,
		},
		{
			"insufficient funds",
			transmitEventLog{
				Log: logpoller.Log{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x010203040"),
				},
				InsufficientFunds: &ac.IAutomationV21PlusCommonInsufficientFundsUpkeepReport{
					Id:      uid.BigInt(),
					Trigger: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			ocr2keepers.InsufficientFundsReportEvent,
		},
		{
			"reorged",
			transmitEventLog{
				Log: logpoller.Log{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x010203040"),
				},
				Reorged: &ac.IAutomationV21PlusCommonReorgedUpkeepReport{
					Id:      uid.BigInt(),
					Trigger: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			ocr2keepers.ReorgReportEvent,
		},
		{
			"empty",
			transmitEventLog{
				Log: logpoller.Log{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x010203040"),
				},
			},
			ocr2keepers.UnknownEvent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.log.Id() != nil {
				require.Equal(t, uid.BigInt().Int64(), tc.log.Id().Int64())
				require.Equal(t, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8}, tc.log.Trigger())
			}
			require.Equal(t, tc.etype, tc.log.TransmitEventType())
		})
	}
}
