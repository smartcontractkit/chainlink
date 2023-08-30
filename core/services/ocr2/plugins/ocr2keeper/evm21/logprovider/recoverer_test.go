package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestLogRecoverer_GetRecoverables(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewLogRecoverer(logger.TestLogger(t), nil, nil, nil, nil, nil, NewOptions(200))

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
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
			},
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
			},
			false,
		},
		{
			"rate limiting",
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "3", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "4", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "5", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "6", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
			},
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "3", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "4", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "5", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
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

func TestLogRecoverer_Clean(t *testing.T) {
	tests := []struct {
		name        string
		pending     []ocr2keepers.UpkeepPayload
		visited     map[string]visitedRecord
		latestBlock int64
		states      []ocr2keepers.UpkeepState
		wantPending []ocr2keepers.UpkeepPayload
		wantVisited []string
	}{
		{
			"empty",
			[]ocr2keepers.UpkeepPayload{},
			map[string]visitedRecord{},
			0,
			[]ocr2keepers.UpkeepState{},
			[]ocr2keepers.UpkeepPayload{},
			[]string{},
		},
		{
			"clean expired",
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
			},
			map[string]visitedRecord{
				"1": visitedRecord{time.Now(), ocr2keepers.UpkeepPayload{
					WorkID: "1",
					Trigger: ocr2keepers.Trigger{
						LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
							BlockNumber: 50,
						},
					},
				}},
				"2": visitedRecord{time.Now(), ocr2keepers.UpkeepPayload{
					WorkID: "2",
					Trigger: ocr2keepers.Trigger{
						LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
							BlockNumber: 50,
						},
					},
				}},
				"3": visitedRecord{time.Now().Add(-time.Hour), ocr2keepers.UpkeepPayload{
					WorkID: "3",
					Trigger: ocr2keepers.Trigger{
						LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
							BlockNumber: 50,
						},
					},
				}},
				"4": visitedRecord{time.Now().Add(-time.Hour), ocr2keepers.UpkeepPayload{
					WorkID: "4",
					Trigger: ocr2keepers.Trigger{
						LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
							BlockNumber: 50,
						},
					},
				}},
			},
			200,
			[]ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.UnknownState,
			},
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1")},
				{WorkID: "2", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2")},
				{WorkID: "4", UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "4")},
			},
			[]string{"1", "2", "4"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(testutils.Context(t))
			defer cancel()

			lookbackBlocks := int64(100)
			r, _, lp, statesReader := setupTestRecoverer(t, time.Millisecond*50, lookbackBlocks)

			lp.On("LatestBlock", mock.Anything).Return(tc.latestBlock, nil)
			statesReader.On("SelectByWorkIDs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.states, nil)

			r.lock.Lock()
			r.pending = tc.pending
			r.visited = tc.visited
			r.lock.Unlock()

			r.clean(ctx)

			r.lock.RLock()
			defer r.lock.RUnlock()

			pending := r.pending
			require.Equal(t, len(tc.wantPending), len(pending))
			sort.Slice(pending, func(i, j int) bool {
				return pending[i].WorkID < pending[j].WorkID
			})
			for i := range pending {
				require.Equal(t, tc.wantPending[i].WorkID, pending[i].WorkID)
			}
			require.Equal(t, len(tc.wantVisited), len(r.visited))
			for _, id := range tc.wantVisited {
				_, ok := r.visited[id]
				require.True(t, ok)
			}
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
		statesErr        error
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
			nil,
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
			nil,
			[]logpoller.Log{},
			nil,
			fmt.Errorf("test error"),
			[]string{},
		},
		{
			"states error",
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
				},
			},
			nil,
			fmt.Errorf("test error"),
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
				},
			},
			[]ocr2keepers.UpkeepState{},
			nil,
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
				},
				{
					upkeepID: big.NewInt(2),
					addr:     common.HexToAddress("0x2").Bytes(),
					topics: []common.Hash{
						common.HexToHash("0x2"),
					},
					configUpdateBlock: 150, // should be filtered out
				},
			},
			[]ocr2keepers.UpkeepState{ocr2keepers.UnknownState},
			nil,
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
			[]string{"84c83c79c2be2c3eabd8d35986a2a798d9187564d7f4f8f96c5a0f40f50bed3f"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookbackBlocks := int64(100)
			recoverer, filterStore, lp, statesReader := setupTestRecoverer(t, time.Millisecond*50, lookbackBlocks)

			filterStore.AddActiveUpkeeps(tc.active...)
			lp.On("LatestBlock", mock.Anything).Return(tc.latestBlock, tc.latestBlockErr)
			lp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.logs, tc.logsErr)
			statesReader.On("SelectByWorkIDs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.states, tc.statesErr)

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
	n := recoveryBatchSize*2 + 2
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

func TestLogRecoverer_getFilterBatch(t *testing.T) {
	tests := []struct {
		name        string
		offsetBlock int64
		filters     []upkeepFilter
		want        int
	}{
		{
			"empty",
			2,
			[]upkeepFilter{},
			0,
		},
		{
			"filter out of range",
			100,
			[]upkeepFilter{
				{
					upkeepID:        big.NewInt(1),
					addr:            common.HexToAddress("0x1").Bytes(),
					lastRePollBlock: 50,
				},
				{
					upkeepID:          big.NewInt(2),
					addr:              common.HexToAddress("0x2").Bytes(),
					lastRePollBlock:   50,
					configUpdateBlock: 101, // out of range
				},
				{
					upkeepID:          big.NewInt(3),
					addr:              common.HexToAddress("0x3").Bytes(),
					configUpdateBlock: 99,
				},
			},
			2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recoverer, filterStore, _, _ := setupTestRecoverer(t, time.Millisecond*50, int64(100))
			filterStore.AddActiveUpkeeps(tc.filters...)
			batch := recoverer.getFilterBatch(tc.offsetBlock)
			require.Equal(t, tc.want, len(batch))
		})
	}
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

func TestLogRecoverer_GetProposalData(t *testing.T) {
	for _, tc := range []struct {
		name        string
		proposal    ocr2keepers.CoordinatedBlockProposal
		skipFilter  bool
		filterStore UpkeepFilterStore
		logPoller   logpoller.LogPoller
		client      client.Client
		stateReader core.UpkeepStateReader
		wantBytes   []byte
		expectErr   bool
		wantErr     error
	}{
		{
			name:      "passing an empty proposal with an empty upkeep ID returns an error",
			proposal:  ocr2keepers.CoordinatedBlockProposal{},
			expectErr: true,
			wantErr:   errors.New("not a log trigger upkeep ID"),
		},
		{
			name: "if a filter is not found for the upkeep ID, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
			},
			skipFilter: true,
			expectErr:  true,
			wantErr:    errors.New("filter not found for upkeep 452312848583266388373324160190187140457511065560374322131410487042692349952"),
		},
		{
			name: "if an error is encountered fetching the latest block, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 0,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 0, errors.New("latest block boom")
				},
			},
			expectErr: true,
			wantErr:   errors.New("latest block boom"),
		},
		{
			name: "if an error is encountered fetching the tx receipt, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 0,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			client: &mockClient{
				CallContextFn: func(ctx context.Context, receipt *types.Receipt, method string, args ...interface{}) error {
					return errors.New("tx receipt boom")
				},
			},
			expectErr: true,
			wantErr:   errors.New("tx receipt boom"),
		},
		{
			name: "if the tx block is nil, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 0,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			client: &mockClient{
				CallContextFn: func(ctx context.Context, receipt *types.Receipt, method string, args ...interface{}) error {
					return nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("failed to get tx block"),
		},
		{
			name: "if a log trigger extension block number is 0, and the block number on the tx receipt is not recoverable, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 0,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			client: &mockClient{
				CallContextFn: func(ctx context.Context, receipt *types.Receipt, method string, args ...interface{}) error {
					receipt.Status = 1
					receipt.BlockNumber = big.NewInt(200)
					return nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("log block is not recoverable"),
		},
		{
			name: "if a log block is not recoverable, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 200,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("log block is not recoverable"),
		},
		{
			name: "if a log block is recoverable, when the upkeep state reader errors, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 80,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return nil, errors.New("upkeep state boom")
				},
			},
			expectErr: true,
			wantErr:   errors.New("upkeep state boom"),
		},
		{
			name: "if a log block is recoverable, when the upkeep state reader returns a non recoverable state, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 80,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.Ineligible,
					}, nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("upkeep state is not recoverable"),
		},
		{
			name: "if a log block is recoverable, when the filter address is empty, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 80,
					},
				},
			},
			filterStore: &mockFilterStore{
				HasFn: func(id *big.Int) bool {
					return true
				},
				RangeFiltersByIDsFn: func(iterator func(int, upkeepFilter), ids ...*big.Int) {

				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.UnknownState,
					}, nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("invalid filter found for upkeepID 452312848583266388373324160190187140457511065560374322131410487042692349952"),
		},
		{
			name: "if a log block is recoverable, when the log poller returns an error fetching logs, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 80,
					},
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
				LogsWithSigsFn: func(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return nil, errors.New("logs with sigs boom")
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.UnknownState,
					}, nil
				},
			},
			expectErr: true,
			wantErr:   errors.New("could not read logs: logs with sigs boom"),
		},
		{
			name: "if a log block is recoverable, when logs cannot be found for an upkeep ID, an error is returned",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: ocr2keepers.Trigger{
					LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
						BlockNumber: 80,
					},
				},
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
				LogsWithSigsFn: func(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 80,
						},
					}, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.UnknownState,
					}, nil
				},
			},
			expectErr: true,
			wantErr:   errors.New(`no log found for upkeepID 452312848583266388373324160190187140457511065560374322131410487042692349952 and trigger {"BlockNumber":0,"BlockHash":"0000000000000000000000000000000000000000000000000000000000000000","LogTriggerExtension":{"BlockHash":"0000000000000000000000000000000000000000000000000000000000000000","BlockNumber":80,"Index":0,"TxHash":"0000000000000000000000000000000000000000000000000000000000000000"}}`),
		},
		{
			name: "happy path with empty check data",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: func() ocr2keepers.Trigger {
					t := ocr2keepers.NewTrigger(
						ocr2keepers.BlockNumber(80),
						[32]byte{1},
					)
					t.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{
						TxHash:      [32]byte{2},
						Index:       uint32(3),
						BlockHash:   [32]byte{1},
						BlockNumber: ocr2keepers.BlockNumber(80),
					}
					return t
				}(),
				WorkID: "d91c6f090b8477f434cf775182e4ff12c90618ba4da5b8ec06aa719768b7743a",
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
				LogsWithSigsFn: func(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 80,
							BlockHash:   [32]byte{1},
							TxHash:      [32]byte{2},
							LogIndex:    3,
						},
					}, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.UnknownState,
					}, nil
				},
			},
			wantBytes: []byte(nil),
		},
		{
			name: "happy path with check data",
			proposal: ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123"),
				Trigger: func() ocr2keepers.Trigger {
					t := ocr2keepers.NewTrigger(
						ocr2keepers.BlockNumber(80),
						[32]byte{1},
					)
					t.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{
						TxHash:      [32]byte{2},
						Index:       uint32(3),
						BlockHash:   [32]byte{1},
						BlockNumber: ocr2keepers.BlockNumber(80),
					}
					return t
				}(),
				WorkID: "d91c6f090b8477f434cf775182e4ff12c90618ba4da5b8ec06aa719768b7743a",
			},
			logPoller: &mockLogPoller{
				LatestBlockFn: func(qopts ...pg.QOpt) (int64, error) {
					return 100, nil
				},
				LogsWithSigsFn: func(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							EvmChainId:     utils.NewBig(big.NewInt(1)),
							LogIndex:       3,
							BlockHash:      [32]byte{1},
							BlockNumber:    80,
							BlockTimestamp: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
							EventSig:       common.HexToHash("abc"),
							TxHash:         [32]byte{2},
							Data:           []byte{1, 2, 3},
							CreatedAt:      time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						},
					}, nil
				},
			},
			stateReader: &mockStateReader{
				SelectByWorkIDsFn: func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
					return []ocr2keepers.UpkeepState{
						ocr2keepers.UnknownState,
					}, nil
				},
			},
			wantBytes: []byte{1, 2, 3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recoverer, filterStore, _, _ := setupTestRecoverer(t, time.Second, 10)

			if !tc.skipFilter {
				filterStore.AddActiveUpkeeps(upkeepFilter{
					addr:     []byte("test"),
					topics:   []common.Hash{common.HexToHash("0x1"), common.HexToHash("0x2"), common.HexToHash("0x3"), common.HexToHash("0x4")},
					upkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "123").BigInt(),
				})
			}

			if tc.filterStore != nil {
				recoverer.filterStore = tc.filterStore
			}
			if tc.logPoller != nil {
				recoverer.poller = tc.logPoller
			}
			if tc.client != nil {
				recoverer.client = tc.client
			}
			if tc.stateReader != nil {
				recoverer.states = tc.stateReader
			}

			b, err := recoverer.GetProposalData(context.Background(), tc.proposal)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantBytes, b)
			}
		})
	}
}

type mockFilterStore struct {
	UpkeepFilterStore
	HasFn               func(id *big.Int) bool
	RangeFiltersByIDsFn func(iterator func(int, upkeepFilter), ids ...*big.Int)
}

func (s *mockFilterStore) RangeFiltersByIDs(iterator func(int, upkeepFilter), ids ...*big.Int) {
	s.RangeFiltersByIDsFn(iterator, ids...)
}

func (s *mockFilterStore) Has(id *big.Int) bool {
	return s.HasFn(id)
}

type mockLogPoller struct {
	logpoller.LogPoller
	LatestBlockFn  func(qopts ...pg.QOpt) (int64, error)
	LogsWithSigsFn func(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error)
}

func (p *mockLogPoller) LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	return p.LogsWithSigsFn(start, end, eventSigs, address, qopts...)
}
func (p *mockLogPoller) LatestBlock(qopts ...pg.QOpt) (int64, error) {
	return p.LatestBlockFn(qopts...)
}

type mockClient struct {
	client.Client
	CallContextFn func(ctx context.Context, receipt *types.Receipt, method string, args ...interface{}) error
}

func (c *mockClient) CallContext(ctx context.Context, r interface{}, method string, args ...interface{}) error {
	receipt := r.(*types.Receipt)
	return c.CallContextFn(ctx, receipt, method, args)
}

type mockStateReader struct {
	SelectByWorkIDsFn func(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error)
}

func (r *mockStateReader) SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
	return r.SelectByWorkIDsFn(ctx, workIDs...)
}

func setupTestRecoverer(t *testing.T, interval time.Duration, lookbackBlocks int64) (*logRecoverer, UpkeepFilterStore, *lpmocks.LogPoller, *mocks.UpkeepStateReader) {
	lp := new(lpmocks.LogPoller)
	statesReader := new(mocks.UpkeepStateReader)
	filterStore := NewUpkeepFilterStore()
	opts := NewOptions(lookbackBlocks)
	opts.ReadInterval = interval / 5
	opts.LookbackBlocks = lookbackBlocks
	recoverer := NewLogRecoverer(logger.TestLogger(t), lp, nil, statesReader, &mockedPacker{}, filterStore, opts)
	return recoverer, filterStore, lp, statesReader
}
