package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core/mocks"
)

func TestLogRecoverer_GetRecoverables(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewLogRecoverer(logger.TestLogger(t), nil, nil, nil, nil, time.Millisecond*10, 0)

	tests := []struct {
		name    string
		pending []ocr2keepers.UpkeepPayload
		want    []ocr2keepers.UpkeepPayload
		wantErr bool
	}{
		{
			"empty",
			[]ocr2keepers.UpkeepPayload{},
			[]ocr2keepers.UpkeepPayload{},
			false,
		},
		{
			"happy flow",
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1"}, {WorkID: "2"},
			},
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1"}, {WorkID: "2"},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r.lock.Lock()
			r.pending = tc.pending
			r.lock.Unlock()

			got, err := r.GetRecoveryProposals(ctx)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Len(t, got, len(tc.want))
		})
	}
}

func TestLogRecoverer_Recover(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name             string
		lookbackBlocks   int64
		latestBlock      int64
		latestBlockErr   error
		active           []upkeepFilter
		states           []ocr2keepers.UpkeepState
		logs             []logpoller.Log
		logsErr          error
		recoverErr       error
		proposalsWorkIDs []string
	}{
		{
			"no filters",
			200,
			300,
			nil,
			[]upkeepFilter{},
			[]ocr2keepers.UpkeepState{},
			[]logpoller.Log{},
			nil,
			nil,
			[]string{},
		},
		{
			"latest block error",
			200,
			0,
			fmt.Errorf("test error"),
			[]upkeepFilter{},
			[]ocr2keepers.UpkeepState{},
			[]logpoller.Log{},
			nil,
			fmt.Errorf("test error"),
			[]string{},
		},
		{
			"get logs error",
			200,
			300,
			nil,
			[]upkeepFilter{
				{
					upkeepID: big.NewInt(1),
					addr:     common.HexToAddress("0x1").Bytes(),
					topics: []common.Hash{
						common.HexToHash("0x1"),
					},
					lastPollBlock: 0,
				},
			},
			[]ocr2keepers.UpkeepState{},
			[]logpoller.Log{},
			fmt.Errorf("test error"),
			nil,
			[]string{},
		},
		{
			"happy flow",
			100,
			200,
			nil,
			[]upkeepFilter{
				{
					upkeepID: big.NewInt(1),
					addr:     common.HexToAddress("0x1").Bytes(),
					topics: []common.Hash{
						common.HexToHash("0x1"),
					},
					lastPollBlock: 0,
				},
			},
			[]ocr2keepers.UpkeepState{ocr2keepers.UnknownState},
			[]logpoller.Log{
				{
					BlockNumber: 2,
					TxHash:      common.HexToHash("0x111"),
					LogIndex:    1,
					BlockHash:   common.HexToHash("0x2"),
				},
			},
			nil,
			nil,
			[]string{"ef1833ce5eb633189430bf52332ef8a1263ae20b9243f77be2809acb58966616"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookbackBlocks := int64(100)
			recoverer, filterStore, lp, statesReader := setupTestRecoverer(t, time.Millisecond*50, lookbackBlocks)

			filterStore.AddActiveUpkeeps(tc.active...)
			lp.On("LatestBlock", mock.Anything).Return(tc.latestBlock, tc.latestBlockErr)
			lp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.logs, tc.logsErr)
			statesReader.On("SelectByWorkIDsInRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.states, nil)

			err := recoverer.recover(ctx)
			if tc.recoverErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			proposals, err := recoverer.GetRecoveryProposals(ctx)
			require.NoError(t, err)
			require.Equal(t, len(tc.proposalsWorkIDs), len(proposals))
			if len(proposals) > 0 {
				sort.Slice(proposals, func(i, j int) bool {
					return proposals[i].WorkID < proposals[j].WorkID
				})
			}
			for i := range proposals {
				require.Equal(t, tc.proposalsWorkIDs[i], proposals[i].WorkID)
			}
		})
	}
}

func TestLogRecoverer_SelectFilterBatch(t *testing.T) {
	n := (recoveryBatchSize*2 + 2)
	filters := []upkeepFilter{}
	for i := 0; i < n; i++ {
		filters = append(filters, upkeepFilter{
			upkeepID: big.NewInt(int64(i)),
		})
	}
	recoverer, _, _, _ := setupTestRecoverer(t, time.Millisecond*50, int64(100))

	batch := recoverer.selectFilterBatch(filters)
	require.Equal(t, recoveryBatchSize, len(batch))

	batch = recoverer.selectFilterBatch(filters[:recoveryBatchSize/2])
	require.Equal(t, recoveryBatchSize/2, len(batch))
}

func TestLogRecoverer_FilterFinalizedStates(t *testing.T) {
	tests := []struct {
		name   string
		logs   []logpoller.Log
		states []ocr2keepers.UpkeepState
		want   []logpoller.Log
	}{
		{
			"empty",
			[]logpoller.Log{},
			[]ocr2keepers.UpkeepState{},
			[]logpoller.Log{},
		},
		{
			"happy flow",
			[]logpoller.Log{
				{LogIndex: 0}, {LogIndex: 2}, {LogIndex: 2},
			},
			[]ocr2keepers.UpkeepState{
				ocr2keepers.UnknownState,
				ocr2keepers.Performed,
				ocr2keepers.Ineligible,
			},
			[]logpoller.Log{
				{LogIndex: 0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recoverer, _, _, _ := setupTestRecoverer(t, time.Millisecond*50, int64(100))
			state := recoverer.filterFinalizedStates(upkeepFilter{}, tc.logs, tc.states)
			require.Equal(t, len(tc.want), len(state))
			for i := range state {
				require.Equal(t, tc.want[i].LogIndex, state[i].LogIndex)
			}
		})
	}
}

func setupTestRecoverer(t *testing.T, interval time.Duration, lookbackBlocks int64) (*logRecoverer, UpkeepFilterStore, *lpmocks.LogPoller, *mocks.UpkeepStateReader) {
	lp := new(lpmocks.LogPoller)
	statesReader := new(mocks.UpkeepStateReader)
	filterStore := NewUpkeepFilterStore()
	recoverer := NewLogRecoverer(logger.TestLogger(t), lp, statesReader, &mockedPacker{}, filterStore, interval, lookbackBlocks)
	return recoverer, filterStore, lp, statesReader
}
