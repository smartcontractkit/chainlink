package transmit

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestTransmitEventLog(t *testing.T) {
	uid := core.GenUpkeepID(ocr2keepers.ConditionTrigger, "111")

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
				Performed: &iregistry21.IKeeperRegistryMasterUpkeepPerformed{
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
				Stale: &iregistry21.IKeeperRegistryMasterStaleUpkeepReport{
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
				InsufficientFunds: &iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{
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
				Reorged: &iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{
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
