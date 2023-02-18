package logpoller

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type block struct {
	number int64
	hash   common.Hash
}

// Setup creates two orms representing logs from different chains.
func setup(t testing.TB) (*ORM, *ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_logs_evm_chain_id_fkey DEFERRED`)))
	o1 := NewORM(big.NewInt(137), db, lggr, pgtest.NewQConfig(true))
	o2 := NewORM(big.NewInt(138), db, lggr, pgtest.NewQConfig(true))
	return o1, o2
}

func TestORM_GetBlocks_From_Range(t *testing.T) {

	o1, _ := setup(t)
	// Insert many blocks and read them back together
	blocks := []block{
		{
			number: 10,
			hash:   common.HexToHash("0x111"),
		},
		{
			number: 11,
			hash:   common.HexToHash("0x112"),
		},
		{
			number: 12,
			hash:   common.HexToHash("0x113"),
		},
		{
			number: 13,
			hash:   common.HexToHash("0x114"),
		},
		{
			number: 14,
			hash:   common.HexToHash("0x115"),
		},
	}
	for _, b := range blocks {
		require.NoError(t, o1.InsertBlock(b.hash, b.number))
	}

	var blockNumbers []uint64
	for _, b := range blocks {
		blockNumbers = append(blockNumbers, uint64(b.number))
	}

	lpBlocks, err := o1.GetBlocksRange(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
	require.NoError(t, err)
	assert.Len(t, lpBlocks, len(blocks))

	// Ignores non-existent block
	lpBlocks2, err := o1.GetBlocksRange(blockNumbers[0], 15)
	require.NoError(t, err)
	assert.Len(t, lpBlocks2, len(blocks))

	// Only non-existent blocks
	lpBlocks3, err := o1.GetBlocksRange(15, 15)
	require.NoError(t, err)
	assert.Len(t, lpBlocks3, 0)
}

func TestORM_GetBlocks_From_Range_Recent_Blocks(t *testing.T) {

	o1, _ := setup(t)
	// Insert many blocks and read them back together
	var recentBlocks []block
	for i := 1; i <= 256; i++ {
		recentBlocks = append(recentBlocks, block{number: int64(i), hash: common.HexToHash(fmt.Sprintf("0x%d", i))})
	}
	for _, b := range recentBlocks {
		require.NoError(t, o1.InsertBlock(b.hash, b.number))
	}

	var blockNumbers []uint64
	for _, b := range recentBlocks {
		blockNumbers = append(blockNumbers, uint64(b.number))
	}

	lpBlocks, err := o1.GetBlocksRange(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
	require.NoError(t, err)
	assert.Len(t, lpBlocks, len(recentBlocks))

	// Ignores non-existent block
	lpBlocks2, err := o1.GetBlocksRange(blockNumbers[0], 257)
	require.NoError(t, err)
	assert.Len(t, lpBlocks2, len(recentBlocks))

	// Only non-existent blocks
	lpBlocks3, err := o1.GetBlocksRange(257, 257)
	require.NoError(t, err)
	assert.Len(t, lpBlocks3, 0)
}

func TestORM(t *testing.T) {
	o1, o2 := setup(t)
	// Insert and read back a block.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 10))
	b, err := o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b.BlockNumber, int64(10))
	assert.Equal(t, b.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b.EvmChainId.String(), "137")

	// Insert blocks from a different chain
	require.NoError(t, o2.InsertBlock(common.HexToHash("0x1234"), 11))
	require.NoError(t, o2.InsertBlock(common.HexToHash("0x1235"), 12))
	b2, err := o2.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b2.BlockNumber, int64(11))
	assert.Equal(t, b2.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b2.EvmChainId.String(), "138")

	latest, err := o1.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(10), latest.BlockNumber)

	latest, err = o2.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(12), latest.BlockNumber)

	// Delete a block (only 10 on chain).
	require.NoError(t, o1.DeleteBlocksAfter(10))
	_, err = o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Delete blocks from another chain.
	require.NoError(t, o2.DeleteBlocksAfter(11))
	_, err = o2.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	// Delete blocks after should also delete block 12.
	_, err = o2.SelectBlockByHash(common.HexToHash("0x1235"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Should be able to insert and read back a log.
	topic := common.HexToHash("0x1599")
	topic2 := common.HexToHash("0x1600")
	require.NoError(t, o1.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(10),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    2,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(11),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    3,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(12),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1235"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    4,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(13),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1235"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    5,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(14),
			EventSig:    topic2,
			Topics:      [][]byte{topic2[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello2"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    6,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(15),
			EventSig:    topic2,
			Topics:      [][]byte{topic2[:]},
			Address:     common.HexToAddress("0x1235"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello2"),
		},
	}))
	logs, err := o1.selectLogsByBlockRange(10, 10)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
	assert.Equal(t, []byte("hello"), logs[0].Data)

	logs, err = o1.SelectLogsByBlockRangeFilter(1, 1, common.HexToAddress("0x1234"), topic)
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))
	logs, err = o1.SelectLogsByBlockRangeFilter(10, 10, common.HexToAddress("0x1234"), topic)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))

	// With no blocks, should be an error
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 0)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	// With block 10, only 0 confs should work
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 10))
	log, err := o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	assert.Equal(t, int64(10), log.BlockNumber)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 1)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	// With block 12, anything <=2 should work
	require.NoError(t, o1.DeleteBlocksAfter(10))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 11))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 12))
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 1)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 2)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 3)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Required for confirmations to work
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 13))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 14))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1236"), 15))
	// Latest log for topic for addr "0x1234" is @ block 11
	lgs, err := o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic}, 0)
	require.NoError(t, err)

	require.Equal(t, 1, len(lgs))
	require.Equal(t, int64(11), lgs[0].BlockNumber)

	// should return two entries one for each address with the latest update
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic}, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return two entries one for each topic for addr 0x1234
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic, topic2}, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return 4 entries one for each (address,topic) combination
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 0)
	require.NoError(t, err)
	require.Equal(t, 4, len(lgs))

	// should return 3 entries of logs with atleast 1 confirmation
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 1)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))

	// should return 2 entries of logs with atleast 2 confirmation
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// Delete logs after should delete all logs.
	err = o1.DeleteLogsAfter(1)
	require.NoError(t, err)
	latest, err = o1.SelectLatestBlock()
	require.NoError(t, err)
	t.Log(latest.BlockNumber)
	logs, err = o1.selectLogsByBlockRange(1, latest.BlockNumber)
	require.NoError(t, err)
	require.Equal(t, 0, len(logs))
}

func insertLogsTopicValueRange(t *testing.T, o *ORM, addr common.Address, blockNumber int, eventSig common.Hash, start, stop int) {
	var lgs []Log
	for i := start; i <= stop; i++ {
		lgs = append(lgs, Log{
			EvmChainId:  utils.NewBig(o.chainID),
			LogIndex:    int64(i),
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(blockNumber),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:], EvmWord(uint64(i)).Bytes()},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		})
	}
	require.NoError(t, o.InsertLogs(lgs))
}

func TestORM_IndexedLogs(t *testing.T) {
	o1, _ := setup(t)
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1"), 1))
	insertLogsTopicValueRange(t, o1, addr, 1, eventSig, 1, 3)
	insertLogsTopicValueRange(t, o1, addr, 2, eventSig, 4, 4) // unconfirmed

	lgs, err := o1.SelectIndexedLogs(addr, eventSig, 1, []common.Hash{EvmWord(1)}, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, EvmWord(1).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexedLogs(addr, eventSig, 1, []common.Hash{EvmWord(1), EvmWord(2)}, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexLogsTopicGreaterThan(addr, eventSig, 1, EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, EvmWord(3).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, EvmWord(1), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(lgs))

	// Check confirmations work as expected.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x2"), 2))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, EvmWord(4), EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x3"), 3))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, EvmWord(4), EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
}

func TestORM_DataWords(t *testing.T) {
	o1, _ := setup(t)
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1"), 1))
	require.NoError(t, o1.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBig(o1.chainID),
			LogIndex:    int64(0),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        EvmWord(1).Bytes(),
		},
		{
			// In block 2, unconfirmed to start
			EvmChainId:  utils.NewBig(o1.chainID),
			LogIndex:    int64(1),
			BlockHash:   common.HexToHash("0x2"),
			BlockNumber: int64(2),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        append(EvmWord(2).Bytes(), EvmWord(3).Bytes()...),
		},
	}))
	// Outside range should fail.
	lgs, err := o1.SelectDataWordRange(addr, eventSig, 0, EvmWord(2), EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))

	// Range including log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 0, EvmWord(1), EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Range only covering log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 0, EvmWord(1), EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Cannot query for unconfirmed second log.
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	// Confirm it, then can query.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x2"), 2))
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, lgs[0].Data, append(EvmWord(2).Bytes(), EvmWord(3).Bytes()...))

	// Check greater than 1 yields both logs.
	lgs, err = o1.SelectDataWordGreaterThan(addr, eventSig, 0, EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))
}

func TestORM_SelectLogsWithSigsByBlockRangeFilter(t *testing.T) {
	o1, _ := setup(t)

	// Insert logs on different topics, should be able to read them
	// back using SelectLogsWithSigsByBlockRangeFilter and specifying
	// said topics.
	topic := common.HexToHash("0x1599")
	topic2 := common.HexToHash("0x1600")
	sourceAddr := common.HexToAddress("0x12345")
	inputLogs := []Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(10),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     sourceAddr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello1"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    2,
			BlockHash:   common.HexToHash("0x1235"),
			BlockNumber: int64(11),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     sourceAddr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello2"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    3,
			BlockHash:   common.HexToHash("0x1236"),
			BlockNumber: int64(12),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1235"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello3"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    4,
			BlockHash:   common.HexToHash("0x1237"),
			BlockNumber: int64(13),
			EventSig:    topic,
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1235"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello4"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    5,
			BlockHash:   common.HexToHash("0x1238"),
			BlockNumber: int64(14),
			EventSig:    topic2,
			Topics:      [][]byte{topic2[:]},
			Address:     sourceAddr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello5"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    6,
			BlockHash:   common.HexToHash("0x1239"),
			BlockNumber: int64(15),
			EventSig:    topic2,
			Topics:      [][]byte{topic2[:]},
			Address:     sourceAddr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello6"),
		},
	}
	require.NoError(t, o1.InsertLogs(inputLogs))

	startBlock, endBlock := int64(10), int64(15)
	logs, err := o1.SelectLogsWithSigsByBlockRangeFilter(startBlock, endBlock, sourceAddr, []common.Hash{
		topic,
		topic2,
	})
	require.NoError(t, err)
	assert.Len(t, logs, 4)
	for _, l := range logs {
		assert.Equal(t, sourceAddr, l.Address, "wrong log address")
		assert.True(t, bytes.Equal(topic.Bytes(), l.EventSig.Bytes()) || bytes.Equal(topic2.Bytes(), l.EventSig.Bytes()), "wrong log topic")
		assert.True(t, l.BlockNumber >= startBlock && l.BlockNumber <= endBlock)
	}
}

func TestORM_DeleteBlocksBefore(t *testing.T) {
	o1, _ := setup(t)
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 1))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 2))
	require.NoError(t, o1.DeleteBlocksBefore(1))
	// 1 should be gone.
	_, err := o1.SelectBlockByNumber(1)
	require.Equal(t, err, sql.ErrNoRows)
	b, err := o1.SelectBlockByNumber(2)
	require.NoError(t, err)
	assert.Equal(t, int64(2), b.BlockNumber)
	// Clear multiple
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1236"), 3))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1237"), 4))
	require.NoError(t, o1.DeleteBlocksBefore(3))
	_, err = o1.SelectBlockByNumber(2)
	require.Equal(t, err, sql.ErrNoRows)
	_, err = o1.SelectBlockByNumber(3)
	require.Equal(t, err, sql.ErrNoRows)
}

func BenchmarkLogs(b *testing.B) {
	o, _ := setup(b)
	var lgs []Log
	addr := common.HexToAddress("0x1234")
	for i := 0; i < 10_000; i++ {
		lgs = append(lgs, Log{
			EvmChainId:  utils.NewBig(o.chainID),
			LogIndex:    int64(i),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: 1,
			EventSig:    EmitterABI.Events["Log1"].ID,
			Topics:      [][]byte{},
			Address:     addr,
			TxHash:      common.HexToHash("0x1234"),
			Data:        common.HexToHash(fmt.Sprintf("0x%d", i)).Bytes(),
		})
	}
	require.NoError(b, o.InsertLogs(lgs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := o.SelectDataWordRange(addr, EmitterABI.Events["Log1"].ID, 0, EvmWord(8000), EvmWord(8002), 0)
		require.NoError(b, err)
	}
}
