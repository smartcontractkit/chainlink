package upkeepstate

import (
	"errors"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPerformedEventsScanner(t *testing.T) {
	ctx := testutils.Context(t)
	registryAddr := common.HexToAddress("0x12345")
	lggr := logger.TestLogger(t)

	tests := []struct {
		name           string
		workIDs        []string
		pollerResults  []logpoller.Log
		scannerResults []string
		pollerErr      error
		errored        bool
	}{
		{
			"empty",
			[]string{},
			[]logpoller.Log{},
			[]string{},
			nil,
			false,
		},
		{
			"log poller error",
			[]string{"111"},
			[]logpoller.Log{},
			[]string{},
			errors.New("test-error"),
			true,
		},
		{
			"one result",
			[]string{"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"},
			[]logpoller.Log{
				{
					BlockNumber: 1,
					Address:     registryAddr,
					Topics: convertTopics([]common.Hash{
						ac.IAutomationV21PlusCommonDedupKeyAdded{}.Topic(),
						common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"),
					}),
				},
			},
			[]string{"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"},
			nil,
			false,
		},
		{
			"missing workID",
			[]string{"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"},
			[]logpoller.Log{
				{
					BlockNumber: 1,
					Address:     registryAddr,
					Topics: convertTopics([]common.Hash{
						ac.IAutomationV21PlusCommonDedupKeyAdded{}.Topic(),
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
			mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
			mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
			scanner := NewPerformedEventsScanner(lggr, mp, registryAddr, 100)

			go func() {
				_ = scanner.Start(ctx)
			}()
			defer func() {
				_ = scanner.Close()
			}()

			mp.On("IndexedLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.pollerResults, tc.pollerErr)

			results, err := scanner.ScanWorkIDs(ctx, tc.workIDs...)
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

func TestPerformedEventsScanner_Batch(t *testing.T) {
	ctx := testutils.Context(t)
	registryAddr := common.HexToAddress("0x12345")
	lggr := logger.TestLogger(t)
	lp := new(mocks.LogPoller)
	scanner := NewPerformedEventsScanner(lggr, lp, registryAddr, 100)

	lp.On("IndexedLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 1,
			Address:     registryAddr,
			Topics: convertTopics([]common.Hash{
				ac.IAutomationV21PlusCommonDedupKeyAdded{}.Topic(),
				common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"),
			}),
		},
	}, nil).Times(1)
	lp.On("IndexedLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 3,
			Address:     registryAddr,
			Topics: convertTopics([]common.Hash{
				ac.IAutomationV21PlusCommonDedupKeyAdded{}.Topic(),
				common.HexToHash("0x331decd9548b62a8d603457658386fc84ba6bc95888008f6362f93160ef3b663"),
			}),
		},
	}, nil).Times(1)

	origWorkIDsBatchSize := workIDsBatchSize
	workIDsBatchSize = 8
	defer func() {
		workIDsBatchSize = origWorkIDsBatchSize
	}()

	ids, err := scanner.ScanWorkIDs(ctx,
		"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
		"1111", "2222", "3333", "4444", "5555", "6666", "7777", "8888", "9999",
		"331decd9548b62a8d603457658386fc84ba6bc95888008f6362f93160ef3b663",
	)

	require.NoError(t, err)
	require.Len(t, ids, 2)
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	require.Equal(t, "290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563", ids[0])
	require.Equal(t, "331decd9548b62a8d603457658386fc84ba6bc95888008f6362f93160ef3b663", ids[1])

	lp.AssertExpectations(t)
}

func convertTopics(topics []common.Hash) [][]byte {
	var topicsForDB [][]byte
	for _, t := range topics {
		topicsForDB = append(topicsForDB, t.Bytes())
	}
	return topicsForDB
}
