package evm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

func TestTransmitEventProvider_ConvertToTransmitEvents(t *testing.T) {
	provider := &TransmitEventProvider{}
	id := genUpkeepID(logTrigger, "111")
	tests := []struct {
		name        string
		performed   []transmitEventLog
		latestBlock int64
		want        []ocr2keepers.TransmitEvent
		errored     bool
	}{
		{
			"happy flow",
			[]transmitEventLog{
				{
					Log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0x0102030405060708010203040506070801020304050607080102030405060708"),
					},
					Performed: &iregistry21.IKeeperRegistryMasterUpkeepPerformed{
						Id: id,
					},
				},
			},
			1,
			[]ocr2keepers.TransmitEvent{
				{
					Type:       ocr2keepers.PerformEvent,
					UpkeepID:   ocr2keepers.UpkeepIdentifier(id.Bytes()),
					CheckBlock: ocr2keepers.BlockKey(""), // empty for log triggers
				},
			},
			false,
		},
		{
			"empty performed",
			[]transmitEventLog{},
			1,
			[]ocr2keepers.TransmitEvent{},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := provider.convertToTransmitEvents(tc.performed, tc.latestBlock)
			require.Equal(t, tc.errored, err != nil)
			require.Len(t, results, len(tc.want))
			for i, res := range results {
				require.Equal(t, tc.want[i].Type, res.Type)
				require.Equal(t, tc.want[i].UpkeepID, res.UpkeepID)
				require.Equal(t, tc.want[i].CheckBlock, res.CheckBlock)
			}
		})
	}
}
