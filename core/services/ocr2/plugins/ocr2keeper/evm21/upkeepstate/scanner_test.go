package upkeepstate

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPerformedEventsScanner(t *testing.T) {
	ctx := testutils.Context(t)
	registryAddr := common.HexToAddress("0x12345")
	lggr := logger.TestLogger(t)

	tests := []struct {
		name           string
		pollerResults  []logpoller.Log
		scannerResults []string
		pollerErr      error
		errored        bool
	}{
		{
			"no logs",
			[]logpoller.Log{},
			[]string{},
			nil,
			false,
		},
		{
			"log poller error",
			[]logpoller.Log{},
			[]string{},
			fmt.Errorf("test-error"),
			true,
		},
		{
			"one result",
			[]logpoller.Log{
				{
					BlockNumber: 1,
					Address:     registryAddr,
					Topics: convertTopics([]common.Hash{
						iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
						common.HexToHash("0x1111"),
					}),
				},
			},
			[]string{common.HexToHash("0x1111").Hex()},
			nil,
			false,
		},
		{
			"missing workID",
			[]logpoller.Log{
				{
					BlockNumber: 1,
					Address:     registryAddr,
					Topics: convertTopics([]common.Hash{
						iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
					}),
				},
			},
			[]string{},
			nil,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mp := new(mocks.LogPoller)
			mp.On("RegisterFilter", mock.Anything).Return(nil)
			mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
			scanner := NewPerformedEventsScanner(lggr, mp, registryAddr)

			go func() {
				_ = scanner.Start(ctx)
			}()
			defer func() {
				_ = scanner.Close()
			}()

			mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.pollerResults, tc.pollerErr)

			results, err := scanner.WorkIDsInRange(ctx, 0, 100)
			if tc.errored {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, len(tc.scannerResults), len(results))

			for _, result := range results {
				require.Contains(t, tc.scannerResults, result)
			}
		})
	}
}

func TestPerformedEventsScanner_LogPollerErrors(t *testing.T) {
	ctx := testutils.Context(t)
	registryAddr := common.HexToAddress("0x12345")
	lggr := logger.TestLogger(t)

	mp := new(mocks.LogPoller)
	scanner := NewPerformedEventsScanner(lggr, mp, registryAddr)

	mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error"))

	workIDs, err := scanner.WorkIDsInRange(ctx, 0, 100)
	require.Error(t, err)
	require.Nil(t, workIDs)
}

func convertTopics(topics []common.Hash) [][]byte {
	var topicsForDB [][]byte
	for _, t := range topics {
		topicsForDB = append(topicsForDB, t.Bytes())
	}
	return topicsForDB
}
