package transmit

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
	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestTransmitEventProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lp := new(mocks.LogPoller)

	lp.On("RegisterFilter", mock.Anything).Return(nil)

	provider, err := NewTransmitEventProvider(logger.TestLogger(t), lp, common.HexToAddress("0x"), client.NewNullClient(big.NewInt(1), logger.TestLogger(t)), 32)
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
			lp.On("LatestBlock", mock.Anything).Return(tc.latestBlock, nil)
			lp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.logs, nil)

			res, err := provider.GetLatestEvents(ctx)
			require.Equal(t, tc.errored, err != nil)
			require.Len(t, res, tc.resultsLen)
		})
	}
}

// {"EvmChainId":"1337","LogIndex":11,"BlockHash":"0x0f220443e2dc988851c69569c8929e2125ebfc81338c030380fb713658efcc52","BlockNumber":536,"BlockTimestamp":"1970-01-01T03:29:20+02:00","Topics":["rYzJV5sh3+LC9uo1uhW2VuRrT1sMtCT1Jzm4zlysnFs=","GTkpnAAAAAAAAAAAAAAAAQ+7hJM0Mbulrh8+1CxZ/nI=","AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE="],"EventSig":"0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b","Address":"0x0ef848e980980c633864139e5cf4004ef0a52324","TxHash":"0xe2d4ca9ebf8321d37deb1466965e9f70d8e9e67f9726a6fd1ed8537ae963cb53","Data":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAU/hhVa6AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABV1QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAes1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgEBqTDkJvytuYCUcOW0LhtlXd22ArWxnmXyu7kdz8FHxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACF0H2cZJJyfkjLQd6BuXkGbidAlLftLwn6eJ/Ew82eea7","CreatedAt":"2023-08-29T18:01:45.829126+03:00"}
func TestTransmitEventProvider_ProcessLogs(t *testing.T) {
	lp := new(mocks.LogPoller)
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	client := evmClientMocks.NewClient(t)

	provider, err := NewTransmitEventProvider(logger.TestLogger(t), lp, common.HexToAddress("0x"), client, 250)
	require.NoError(t, err)

	id := core.GenUpkeepID(ocr2keepers.LogTrigger, "1111111111111111")

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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parseResults := make(map[string]transmitEventLog, len(tc.parsedPerformed))
			performedLogs := make([]logpoller.Log, len(tc.parsedPerformed))
			for i, l := range tc.parsedPerformed {
				performedLogs[i] = l.Log
				parseResults[logKey(l.Log)] = l
			}
			provider.mu.Lock()
			provider.parseLog = func(registry *iregistry21.IKeeperRegistryMaster, log logpoller.Log) (transmitEventLog, error) {
				return parseResults[logKey(log)], nil
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
