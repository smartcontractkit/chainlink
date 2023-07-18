package evm

import (
	"context"
	"math/big"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTransmitEventProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mp := new(mocks.LogPoller)

	mp.On("RegisterFilter", mock.Anything).Return(nil)

	provider, err := NewTransmitEventProvider(logger.TestLogger(t), mp, common.HexToAddress("0x"), client.NewNullClient(big.NewInt(1), logger.TestLogger(t)), 32)
	require.NoError(t, err)
	require.NotNil(t, provider)

	go func() {
		require.Error(t, provider.Ready())
		errs := provider.HealthReport()
		require.Len(t, errs, 1)
		require.Error(t, errs[provider.Name()])
		err := provider.Start(ctx)
		if ctx.Err() == nil {
			require.NoError(t, err)
		}
	}()

	go func() {
		for provider.Ready() != nil {
			runtime.Gosched()
		}
		errs := provider.HealthReport()
		require.Len(t, errs, 1)
		require.NoError(t, errs[provider.Name()])

		mp.On("LatestBlock", mock.Anything).Return(int64(1), nil)
		mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
			// TODO: return values
		}, nil)

		res, err := provider.Events(ctx)
		require.NoError(t, err)
		require.Len(t, res, 0)
		// TODO: check values returned from mock
		// require.Len(t, res, 1)

		cancel()
	}()

	<-ctx.Done()
}

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
			"empty events",
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
