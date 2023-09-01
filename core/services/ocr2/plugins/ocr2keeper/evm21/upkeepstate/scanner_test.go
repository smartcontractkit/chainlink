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
			fmt.Errorf("test-error"),
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
						iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
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

func convertTopics(topics []common.Hash) [][]byte {
	var topicsForDB [][]byte
	for _, t := range topics {
		topicsForDB = append(topicsForDB, t.Bytes())
	}
	return topicsForDB
}
