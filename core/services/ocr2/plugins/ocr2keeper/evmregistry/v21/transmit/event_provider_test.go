package transmit

import (
	"math/big"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestTransmitEventProvider_Sanity(t *testing.T) {
	ctx := testutils.Context(t)

	lp := new(mocks.LogPoller)

	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)

	provider, err := NewTransmitEventProvider(ctx, logger.TestLogger(t), lp, common.HexToAddress("0x"), client.NewNullClient(big.NewInt(1), logger.TestLogger(t)), 32)
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
						ac.IAutomationV21PlusCommonUpkeepPerformed{}.Topic().Bytes(),
					},
					EventSig: ac.IAutomationV21PlusCommonUpkeepPerformed{}.Topic(),
					Data:     []byte{},
				},
			},
			false,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lp.On("LatestBlock", mock.Anything).Return(logpoller.Block{BlockNumber: tc.latestBlock}, nil)
			lp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.logs, nil)

			res, err := provider.GetLatestEvents(ctx)
			require.Equal(t, tc.errored, err != nil)
			require.Len(t, res, tc.resultsLen)
		})
	}
}

func TestTransmitEventProvider_ProcessLogs(t *testing.T) {
	lp := new(mocks.LogPoller)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	client := clienttest.NewClient(t)
	ctx := testutils.Context(t)

	provider, err := NewTransmitEventProvider(ctx, logger.TestLogger(t), lp, common.HexToAddress("0x"), client, 250)
	require.NoError(t, err)

	id := core.GenUpkeepID(types.LogTrigger, "1111111111111111")

	tests := []struct {
		name            string
		parsedPerformed []transmitEventLog
		latestBlock     int64
		want            []ocr2keepers.TransmitEvent
		errored         bool
	}{
		{
			"happy flow",
			[]transmitEventLog{
				{
					Log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0x0102030405060708010203040506070801020304050607080102030405060708"),
					},
					Performed: &ac.IAutomationV21PlusCommonUpkeepPerformed{
						Id: id.BigInt(),
						Trigger: func() []byte {
							b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111abc0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
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
					CheckBlock: ocr2keepers.BlockNumber(1),
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
		{
			"same log twice", // shouldn't happen in practice as log poller should not return duplicate logs
			[]transmitEventLog{
				{
					Log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0x0102030405060708010203040506070801020304050607080102030405060708"),
					},
					Performed: &ac.IAutomationV21PlusCommonUpkeepPerformed{
						Id: id.BigInt(),
						Trigger: func() []byte {
							b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111abc0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
							return b
						}(),
					},
				},
				{
					Log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0x0102030405060708010203040506070801020304050607080102030405060708"),
					},
					Performed: &ac.IAutomationV21PlusCommonUpkeepPerformed{
						Id: id.BigInt(),
						Trigger: func() []byte {
							b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111abc0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
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
					CheckBlock: ocr2keepers.BlockNumber(1),
				},
				{
					Type:       ocr2keepers.PerformEvent,
					UpkeepID:   id,
					CheckBlock: ocr2keepers.BlockNumber(1),
				},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parseResults := make(map[string]transmitEventLog, len(tc.parsedPerformed))
			performedLogs := make([]logpoller.Log, len(tc.parsedPerformed))
			for i, l := range tc.parsedPerformed {
				performedLogs[i] = l.Log
				if _, ok := parseResults[provider.logKey(l.Log)]; ok {
					continue
				}
				parseResults[provider.logKey(l.Log)] = l
			}
			provider.mu.Lock()
			provider.cache = newTransmitEventCache(provider.cache.cap)
			provider.parseLog = func(registry *ac.IAutomationV21PlusCommon, log logpoller.Log) (transmitEventLog, error) {
				return parseResults[provider.logKey(log)], nil
			}
			provider.mu.Unlock()

			results, err := provider.processLogs(tc.latestBlock, performedLogs...)
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
