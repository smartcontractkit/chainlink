package log

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func TestUnit_AddLog(t *testing.T) {
	t.Parallel()
	var p iLogPool = newLogPool()

	blockHash := common.BigToHash(big.NewInt(1))
	l1 := types.Log{
		BlockHash:   blockHash,
		Index:       42,
		BlockNumber: 1,
	}
	// 1st log added should be the minimum
	assert.True(t, p.addLog(l1), "AddLog should have returned true for first log added")
	require.Equal(t, 1, p.testOnly_getNumLogsForBlock(blockHash))

	// Reattempting to add same log should work, but shouldn't be the minimum
	assert.False(t, p.addLog(l1), "AddLog should have returned false for a 2nd reattempt")
	require.Equal(t, 1, p.testOnly_getNumLogsForBlock(blockHash))

	// 2nd log with same loghash should add a new log, which shouldn't be minimum
	l2 := l1
	l2.Index = 43
	assert.False(t, p.addLog(l2), "AddLog should have returned false for same log added")
	require.Equal(t, 2, p.testOnly_getNumLogsForBlock(blockHash))

	// New log with different larger BlockNumber/loghash should add a new log, not as minimum
	l3 := l1
	l3.BlockNumber = 3
	l3.BlockHash = common.BigToHash(big.NewInt(3))
	assert.False(t, p.addLog(l3), "AddLog should have returned false for same log added")
	assert.Equal(t, 2, p.testOnly_getNumLogsForBlock(blockHash))
	require.Equal(t, 1, p.testOnly_getNumLogsForBlock(l3.BlockHash))

	// New log with different smaller BlockNumber/loghash should add a new log, as minimum
	l4 := l1
	l4.BlockNumber = 0 // New minimum block number
	l4.BlockHash = common.BigToHash(big.NewInt(0))
	assert.True(t, p.addLog(l4), "AddLog should have returned true for smallest BlockNumber")
	assert.Equal(t, 2, p.testOnly_getNumLogsForBlock(blockHash))
	assert.Equal(t, 1, p.testOnly_getNumLogsForBlock(l3.BlockHash))
	require.Equal(t, 1, p.testOnly_getNumLogsForBlock(l4.BlockHash))
}

func TestUnit_GetAndDeleteAll(t *testing.T) {
	t.Parallel()
	var p iLogPool = newLogPool()
	p.addLog(L1)
	p.addLog(L1) // duplicate an add
	p.addLog(L21)
	p.addLog(L22)
	p.addLog(L3)

	logsOnBlock, lowest, highest := p.getAndDeleteAll()

	assert.Equal(t, int64(1), lowest)
	assert.Equal(t, int64(3), highest)
	assert.Len(t, logsOnBlock, 3)
	for _, logs := range logsOnBlock {
		switch logs.BlockNumber {
		case 1:
			l1s := [1]types.Log{L1}
			assert.ElementsMatch(t, l1s, logs.Logs)
		case 2:
			l2s := [2]types.Log{L21, L22}
			assert.ElementsMatch(t, l2s, logs.Logs)
		case 3:
			l3s := [1]types.Log{L3}
			assert.ElementsMatch(t, l3s, logs.Logs)
		default:
			t.Errorf("Received unexpected BlockNumber in results: %d", logs.BlockNumber)
		}
	}
	assert.Equal(t, 0, p.testOnly_getNumLogsForBlock(L1.BlockHash))
	assert.Equal(t, 0, p.testOnly_getNumLogsForBlock(L21.BlockHash))
	assert.Equal(t, 0, p.testOnly_getNumLogsForBlock(L3.BlockHash))
}

func TestUnit_GetLogsToSendWhenEmptyPool(t *testing.T) {
	t.Parallel()
	var p iLogPool = newLogPool()
	logsOnBlocks, minBlockNumToSend := p.getLogsToSend(1)
	assert.Equal(t, int64(0), minBlockNumToSend)
	assert.ElementsMatch(t, []logsOnBlock{}, logsOnBlocks)
}

func TestUnit_GetLogsToSend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		latestBlockNumber         int64
		expectedMinBlockNumToSend int64
		expectedLogs              []logsOnBlock
	}{
		{
			name:                      "NoLogsToSend",
			latestBlockNumber:         0,
			expectedMinBlockNumToSend: 1,
			expectedLogs:              []logsOnBlock{},
		},
		{
			name:                      "PartialLogsToSend",
			latestBlockNumber:         2,
			expectedMinBlockNumToSend: 1,
			expectedLogs: []logsOnBlock{
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
			expectedLogs: []logsOnBlock{
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

	var p iLogPool = newLogPool()
	p.addLog(L1)
	p.addLog(L21)
	p.addLog(L3)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logsOnBlocks, minBlockNumToSend := p.getLogsToSend(test.latestBlockNumber)
			assert.Equal(t, test.expectedMinBlockNumToSend, minBlockNumToSend)
			assert.ElementsMatch(t, test.expectedLogs, logsOnBlocks)
		})
	}
}

func TestUnit_DeleteOlderLogsWhenEmptyPool(t *testing.T) {
	t.Parallel()
	var p iLogPool = newLogPool()
	keptDepth := p.deleteOlderLogs(1)
	var expectedKeptDepth *int64 = nil
	require.Equal(t, expectedKeptDepth, keptDepth)
}

func TestUnit_DeleteOlderLogs(t *testing.T) {
	t.Parallel()
	keptDepth3 := int64(3)
	keptDepth1 := int64(1)
	tests := []struct {
		name                string
		keptDepth           int64
		expectedOldestBlock *int64
		expectedKeptLogs    []logsOnBlock
	}{
		{
			name:                "AllLogsDeleted",
			keptDepth:           4,
			expectedOldestBlock: nil,
			expectedKeptLogs:    []logsOnBlock{},
		},
		{
			name:                "PartialLogsDeleted",
			keptDepth:           3,
			expectedOldestBlock: &keptDepth3,
			expectedKeptLogs: []logsOnBlock{
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
			expectedKeptLogs: []logsOnBlock{
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
			var p iLogPool = newLogPool()
			p.addLog(L1)
			p.addLog(L21)
			p.addLog(L3)

			oldestKeptBlock := p.deleteOlderLogs(test.keptDepth)

			assert.Equal(t, test.expectedOldestBlock, oldestKeptBlock)
			keptLogs, _ := p.getLogsToSend(4)
			assert.ElementsMatch(t, test.expectedKeptLogs, keptLogs)
		})
	}
}

func TestUnit_RemoveBlockWhenEmptyPool(t *testing.T) {
	t.Parallel()
	var p iLogPool = newLogPool()
	p.removeBlock(L1.BlockHash, L1.BlockNumber)
}

func TestUnit_RemoveBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		blockHash             common.Hash
		blockNumber           uint64
		expectedRemainingLogs []logsOnBlock
	}{
		{
			name:                  "BlockNotFound",
			blockHash:             L1.BlockHash,
			blockNumber:           L1.BlockNumber,
			expectedRemainingLogs: []logsOnBlock{},
		},
		{
			name:                  "BlockNumberWasUnique",
			blockHash:             L3.BlockHash,
			blockNumber:           L3.BlockNumber,
			expectedRemainingLogs: []logsOnBlock{},
		},
		{
			name:        "MultipleBlocksWithSameBlockNumber",
			blockHash:   L21.BlockHash,
			blockNumber: L21.BlockNumber,
			expectedRemainingLogs: []logsOnBlock{
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
			var p iLogPool = newLogPool()
			p.addLog(L21)
			p.addLog(L22)
			p.addLog(L23)
			p.addLog(L3)

			p.removeBlock(test.blockHash, test.blockNumber)

			assert.Equal(t, 0, p.testOnly_getNumLogsForBlock(test.blockHash))
			p.deleteOlderLogs(int64(test.blockNumber)) // Pruning logs for easier testing next line
			logsOnBlock, _ := p.getLogsToSend(int64(test.blockNumber))
			assert.ElementsMatch(t, test.expectedRemainingLogs, logsOnBlock)
		})
	}
}
