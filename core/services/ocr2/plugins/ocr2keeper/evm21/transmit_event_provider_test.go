package evm

import (
	"context"
	"math/big"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
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

	for provider.Ready() != nil {
		runtime.Gosched()
	}
	errs := provider.HealthReport()
	require.Len(t, errs, 1)
	require.NoError(t, errs[provider.Name()])

	tests := []struct {
		name        string
		latestBlock int64
		logs        []logpoller.Log
		errored     bool
		resultsLen  int
	}{
		// TODO: add more tests
		{
			"empty logs",
			100,
			[]logpoller.Log{},
			false,
			0,
		},
		{
			"invalid log",
			101,
			[]logpoller.Log{
				{
					BlockNumber: 101,
					BlockHash:   common.HexToHash("0x1"),
					TxHash:      common.HexToHash("0x1"),
					LogIndex:    1,
					Address:     common.HexToAddress("0x1"),
					Topics: [][]byte{
						iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic().Bytes(),
					},
					EventSig: iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
					Data:     []byte{},
				},
			},
			false,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mp.On("LatestBlock", mock.Anything).Return(tc.latestBlock, nil)
			mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.logs, nil)

			res, err := provider.GetLatestEvents(ctx)
			require.Equal(t, tc.errored, err != nil)
			require.Len(t, res, tc.resultsLen)
		})
	}
}

func TestTransmitEventProvider_ConvertToTransmitEvents(t *testing.T) {
	provider := &TransmitEventProvider{}
	id := genUpkeepID(ocr2keepers.LogTrigger, "1111111111111111")
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
						Id: id.BigInt(),
						Trigger: func() []byte {
							b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
							return b
						}(),
					},
				},
			},
			1,
			[]ocr2keepers.TransmitEvent{
				{
					Type:       ocr2keepers.PerformEvent,
					UpkeepID:   id,
					CheckBlock: ocr2keepers.BlockNumber(1), // empty for log triggers
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

func TestTransmitEventLog(t *testing.T) {
	uid := genUpkeepID(ocr2keepers.ConditionTrigger, "111")

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
