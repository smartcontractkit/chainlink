package log_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/stretchr/testify/require"
)

var (
	L1 = types.Log{
		BlockHash:   common.HexToHash("1"),
		Index:       1,
		BlockNumber: 1,
	}

	L21 = types.Log{
		BlockHash:   common.HexToHash("2"),
		Index:       21,
		BlockNumber: 2,
	}

	// L21 and L22 differ only in index
	L22 = types.Log{
		BlockHash:   common.HexToHash("2"),
		Index:       22,
		BlockNumber: 2,
	}

	// L23 is a different BlockHash than L21 and L22
	L23 = types.Log{
		BlockHash:   common.HexToHash("23"),
		Index:       21,
		BlockNumber: 2,
	}

	L3 = types.Log{
		BlockHash:   common.HexToHash("3"),
		Index:       3,
		BlockNumber: 3,
	}
)

func TestPool_AddLog(t *testing.T) {
	t.Parallel()
	p := log.NewLogPool()

	blockHash := common.BigToHash(big.NewInt(1))
	l1 := types.Log{
		BlockHash:   blockHash,
		Index:       42,
		BlockNumber: 1,
	}
	// 1st log added should be the minimum
	require.True(t, p.AddLog(l1), "AddLog should have returned true for first log added")
	require.Equal(t, 1, p.TestOnly_getNumLogsForBlock(blockHash))

	// Reattempting to add same log should work, but shouldn't be the minimum
	require.False(t, p.AddLog(l1), "AddLog should have returned false for a 2nd reattempt")
	require.Equal(t, 1, p.TestOnly_getNumLogsForBlock(blockHash))

	// 2nd log with same loghash should add a new log, which shouldn't be minimum
	l2 := l1
	l2.Index = 43
	require.False(t, p.AddLog(l2), "AddLog should have returned false for same log added")
	require.Equal(t, 2, p.TestOnly_getNumLogsForBlock(blockHash))

	// New log with different larger BlockNumber/loghash should add a new log, not as minimum
	l3 := l1
	l3.BlockNumber = 3
	l3.BlockHash = common.BigToHash(big.NewInt(3))
	require.False(t, p.AddLog(l3), "AddLog should have returned false for same log added")
	require.Equal(t, 2, p.TestOnly_getNumLogsForBlock(blockHash))
	require.Equal(t, 1, p.TestOnly_getNumLogsForBlock(l3.BlockHash))

	// New log with different smaller BlockNumber/loghash should add a new log, as minimum
	l4 := l1
	l4.BlockNumber = 0 // New minimum block number
	l4.BlockHash = common.BigToHash(big.NewInt(0))
	require.True(t, p.AddLog(l4), "AddLog should have returned true for smallest BlockNumber")
	require.Equal(t, 2, p.TestOnly_getNumLogsForBlock(blockHash))
	require.Equal(t, 1, p.TestOnly_getNumLogsForBlock(l3.BlockHash))
	require.Equal(t, 1, p.TestOnly_getNumLogsForBlock(l4.BlockHash))
}

func TestPool_GetAndDeleteAll(t *testing.T) {
	t.Parallel()
	p := log.NewLogPool()
	p.AddLog(L1)
	p.AddLog(L1) // duplicate an add
	p.AddLog(L21)
	p.AddLog(L22)
	p.AddLog(L3)

	logsOnBlock, lowest, highest := p.GetAndDeleteAll()

	require.Equal(t, int64(1), lowest)
	require.Equal(t, int64(3), highest)
	require.Len(t, logsOnBlock, 3)
	for _, logs := range logsOnBlock {
		switch logs.BlockNumber {
		case 1:
			l1s := [1]types.Log{L1}
			require.ElementsMatch(t, l1s, logs.Logs)
		case 2:
			l2s := [2]types.Log{L21, L22}
			require.ElementsMatch(t, l2s, logs.Logs)
		case 3:
			l3s := [1]types.Log{L3}
			require.ElementsMatch(t, l3s, logs.Logs)
		default:
			t.Errorf("Received unexpected BlockNumber in results: %d", logs.BlockNumber)
		}
	}
	require.Equal(t, 0, p.TestOnly_getNumLogsForBlock(L1.BlockHash))
	require.Equal(t, 0, p.TestOnly_getNumLogsForBlock(L21.BlockHash))
	require.Equal(t, 0, p.TestOnly_getNumLogsForBlock(L3.BlockHash))
}

func TestPool_GetLogsToSendWhenEmptyPool(t *testing.T) {
	t.Parallel()
	p := log.NewLogPool()
	logsOnBlocks, minBlockNumToSend := p.GetLogsToSend(1)
	require.Equal(t, int64(0), minBlockNumToSend)
	require.ElementsMatch(t, []log.LogsOnBlock{}, logsOnBlocks)
}

func TestPool_GetLogsToSend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		latestBlockNumber         int64
		expectedMinBlockNumToSend int64
		expectedLogs              []log.LogsOnBlock
	}{
		{
			name:                      "NoLogsToSend",
			latestBlockNumber:         0,
			expectedMinBlockNumToSend: 1,
			expectedLogs:              []log.LogsOnBlock{},
		},
		{
			name:                      "PartialLogsToSend",
			latestBlockNumber:         2,
			expectedMinBlockNumToSend: 1,
			expectedLogs: []log.LogsOnBlock{
				{
					BlockNumber: 1,
					Logs: []types.Log{
						L1,
					},
				},
				{
					BlockNumber: 2,
					Logs: []types.Log{
						L21,
					},
				},
			},
		},
		{
			name:                      "AllLogsToSend",
			latestBlockNumber:         4,
			expectedMinBlockNumToSend: 1,
			expectedLogs: []log.LogsOnBlock{
				{
					BlockNumber: 1,
					Logs: []types.Log{
						L1,
					},
				},
				{
					BlockNumber: 2,
					Logs: []types.Log{
						L21,
					},
				},
				{
					BlockNumber: 3,
					Logs: []types.Log{
						L3,
					},
				},
			},
		},
	}

	p := log.NewLogPool()
	p.AddLog(L1)
	p.AddLog(L21)
	p.AddLog(L3)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logsOnBlocks, minBlockNumToSend := p.GetLogsToSend(test.latestBlockNumber)
			require.Equal(t, test.expectedMinBlockNumToSend, minBlockNumToSend)
			require.ElementsMatch(t, test.expectedLogs, logsOnBlocks)
		})
	}
}

func TestPool_DeleteOlderLogsWhenEmptyPool(t *testing.T) {
	t.Parallel()
	p := log.NewLogPool()
	keptDepth := p.DeleteOlderLogs(1)
	var expectedKeptDepth *int64 = nil
	require.Equal(t, expectedKeptDepth, keptDepth)
}

func TestPool_DeleteOlderLogs(t *testing.T) {
	t.Parallel()
	keptDepth3 := int64(3)
	keptDepth1 := int64(1)
	tests := []struct {
		name                string
		keptDepth           int64
		expectedOldestBlock *int64
		expectedKeptLogs    []log.LogsOnBlock
	}{
		{
			name:                "AllLogsDeleted",
			keptDepth:           4,
			expectedOldestBlock: nil,
			expectedKeptLogs:    []log.LogsOnBlock{},
		},
		{
			name:                "PartialLogsDeleted",
			keptDepth:           3,
			expectedOldestBlock: &keptDepth3,
			expectedKeptLogs: []log.LogsOnBlock{
				{
					BlockNumber: 3,
					Logs: []types.Log{
						L3,
					},
				},
			},
		},
		{
			name:                "NoLogsDeleted",
			keptDepth:           0,
			expectedOldestBlock: &keptDepth1,
			expectedKeptLogs: []log.LogsOnBlock{
				{
					BlockNumber: 3,
					Logs: []types.Log{
						L3,
					},
				},
				{
					BlockNumber: 2,
					Logs: []types.Log{
						L21,
					},
				},
				{
					BlockNumber: 1,
					Logs: []types.Log{
						L1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := log.NewLogPool()
			p.AddLog(L1)
			p.AddLog(L21)
			p.AddLog(L3)

			oldestKeptBlock := p.DeleteOlderLogs(test.keptDepth)

			require.Equal(t, test.expectedOldestBlock, oldestKeptBlock)
			keptLogs, _ := p.GetLogsToSend(4)
			require.ElementsMatch(t, test.expectedKeptLogs, keptLogs)
		})
	}
}

func TestPool_RemoveBlockWhenEmptyPool(t *testing.T) {
	t.Parallel()
	p := log.NewLogPool()
	p.RemoveBlock(L1.BlockHash, L1.BlockNumber)
}

func TestPool_RemoveBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		blockHash             common.Hash
		blockNumber           uint64
		expectedRemainingLogs []log.LogsOnBlock
	}{
		{
			name:                  "BlockNotFound",
			blockHash:             L1.BlockHash,
			blockNumber:           L1.BlockNumber,
			expectedRemainingLogs: []log.LogsOnBlock{},
		},
		{
			name:                  "BlockNumberWasUnique",
			blockHash:             L3.BlockHash,
			blockNumber:           L3.BlockNumber,
			expectedRemainingLogs: []log.LogsOnBlock{},
		},
		{
			name:        "MultipleBlocksWithSameBlockNumber",
			blockHash:   L21.BlockHash,
			blockNumber: L21.BlockNumber,
			expectedRemainingLogs: []log.LogsOnBlock{
				{
					BlockNumber: L23.BlockNumber,
					Logs: []types.Log{
						L23,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := log.NewLogPool()
			p.AddLog(L21)
			p.AddLog(L22)
			p.AddLog(L23)
			p.AddLog(L3)

			p.RemoveBlock(test.blockHash, test.blockNumber)

			require.Equal(t, 0, p.TestOnly_getNumLogsForBlock(test.blockHash))
			p.DeleteOlderLogs(int64(test.blockNumber)) // Pruning logs for easier testing next line
			logsOnBlock, _ := p.GetLogsToSend(int64(test.blockNumber))
			require.ElementsMatch(t, test.expectedRemainingLogs, logsOnBlock)
		})
	}
}
