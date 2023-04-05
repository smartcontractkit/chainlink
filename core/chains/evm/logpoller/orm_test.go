package logpoller_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type block struct {
	number    int64
	hash      common.Hash
	timestamp int64
}

func GenLog(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address) logpoller.Log {
	return logpoller.Log{
		EvmChainId:  utils.NewBig(chainID),
		LogIndex:    logIndex,
		BlockHash:   common.HexToHash(blockHash),
		BlockNumber: blockNum,
		EventSig:    common.BytesToHash(topic1),
		Topics:      [][]byte{topic1},
		Address:     address,
		TxHash:      common.HexToHash("0x1234"),
		Data:        append([]byte("hello "), byte(blockNum)),
	}
}

func TestLogPoller_Batching(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)
	var logs []logpoller.Log
	// Inserts are limited to 65535 parameters. A log being 10 parameters this results in
	// a maximum of 6553 log inserts per tx. As inserting more than 6553 would result in
	// an error without batching, this test makes sure batching is enabled.
	for i := 0; i < 15000; i++ {
		logs = append(logs, GenLog(th.ChainID, int64(i+1), 1, "0x3", EmitterABI.Events["Log1"].ID.Bytes(), th.EmitterAddress1))
	}
	require.NoError(t, th.ORM.InsertLogs(logs))
	lgs, err := th.ORM.SelectLogsByBlockRange(1, 1)
	require.NoError(t, err)
	// Make sure all logs are inserted
	require.Equal(t, len(logs), len(lgs))
}

func TestORM_GetBlocks_From_Range(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	// Insert many blocks and read them back together
	blocks := []block{
		{
			number:    10,
			hash:      common.HexToHash("0x111"),
			timestamp: 0,
		},
		{
			number:    11,
			hash:      common.HexToHash("0x112"),
			timestamp: 10,
		},
		{
			number:    12,
			hash:      common.HexToHash("0x113"),
			timestamp: 20,
		},
		{
			number:    13,
			hash:      common.HexToHash("0x114"),
			timestamp: 30,
		},
		{
			number:    14,
			hash:      common.HexToHash("0x115"),
			timestamp: 40,
		},
	}
	for _, b := range blocks {
		require.NoError(t, o1.InsertBlock(b.hash, b.number, time.Unix(b.timestamp, 0).UTC()))
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
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	// Insert many blocks and read them back together
	var recentBlocks []block
	for i := 1; i <= 256; i++ {
		recentBlocks = append(recentBlocks, block{number: int64(i), hash: common.HexToHash(fmt.Sprintf("0x%d", i))})
	}
	for _, b := range recentBlocks {
		require.NoError(t, o1.InsertBlock(b.hash, b.number, time.Now()))
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
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	o2 := th.ORM2
	// Insert and read back a block.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 10, time.Now()))
	b, err := o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b.BlockNumber, int64(10))
	assert.Equal(t, b.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b.EvmChainId.String(), th.ChainID.String())

	// Insert blocks from a different chain
	require.NoError(t, o2.InsertBlock(common.HexToHash("0x1234"), 11, time.Now()))
	require.NoError(t, o2.InsertBlock(common.HexToHash("0x1235"), 12, time.Now()))
	b2, err := o2.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b2.BlockNumber, int64(11))
	assert.Equal(t, b2.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b2.EvmChainId.String(), th.ChainID2.String())

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
	require.NoError(t, o1.InsertLogs([]logpoller.Log{
		{
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
	logs, err := o1.SelectLogsByBlockRange(10, 10)
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
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 10, time.Now()))
	log, err := o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	assert.Equal(t, int64(10), log.BlockNumber)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 1)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	// With block 12, anything <=2 should work
	require.NoError(t, o1.DeleteBlocksAfter(10))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 11, time.Now()))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 12, time.Now()))
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
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 13, time.Now()))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 14, time.Now()))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1236"), 15, time.Now()))
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
	logs, err = o1.SelectLogsByBlockRange(1, latest.BlockNumber)
	require.NoError(t, err)
	require.Equal(t, 0, len(logs))
}

func insertLogsTopicValueRange(t *testing.T, chainID *big.Int, o *logpoller.ORM, addr common.Address, blockNumber int, eventSig common.Hash, start, stop int) {
	var lgs []logpoller.Log
	for i := start; i <= stop; i++ {
		lgs = append(lgs, logpoller.Log{
			EvmChainId:  utils.NewBig(chainID),
			LogIndex:    int64(i),
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(blockNumber),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:], logpoller.EvmWord(uint64(i)).Bytes()},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		})
	}
	require.NoError(t, o.InsertLogs(lgs))
}

func TestORM_IndexedLogs(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1"), 1, time.Now()))
	insertLogsTopicValueRange(t, th.ChainID, o1, addr, 1, eventSig, 1, 3)
	insertLogsTopicValueRange(t, th.ChainID, o1, addr, 2, eventSig, 4, 4) // unconfirmed

	lgs, err := o1.SelectIndexedLogs(addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)}, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(1).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexedLogs(addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1), logpoller.EvmWord(2)}, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexedLogsByBlockRangeFilter(1, 1, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.SelectIndexedLogsByBlockRangeFilter(1, 2, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(2)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.SelectIndexedLogsByBlockRangeFilter(1, 2, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	_, err = o1.SelectIndexedLogsByBlockRangeFilter(1, 2, addr, eventSig, 0, []common.Hash{logpoller.EvmWord(1)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 0")
	_, err = o1.SelectIndexedLogsByBlockRangeFilter(1, 2, addr, eventSig, 4, []common.Hash{logpoller.EvmWord(1)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 4")

	lgs, err = o1.SelectIndexLogsTopicGreaterThan(addr, eventSig, 1, logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(3).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, logpoller.EvmWord(1), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(lgs))

	// Check confirmations work as expected.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x2"), 2, time.Now()))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, logpoller.EvmWord(4), logpoller.EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x3"), 3, time.Now()))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig, 1, logpoller.EvmWord(4), logpoller.EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
}

func TestORM_DataWords(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1"), 1, time.Now()))
	require.NoError(t, o1.InsertLogs([]logpoller.Log{
		{
			EvmChainId:  utils.NewBig(th.ChainID),
			LogIndex:    int64(0),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        logpoller.EvmWord(1).Bytes(),
		},
		{
			// In block 2, unconfirmed to start
			EvmChainId:  utils.NewBig(th.ChainID),
			LogIndex:    int64(1),
			BlockHash:   common.HexToHash("0x2"),
			BlockNumber: int64(2),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...),
		},
	}))
	// Outside range should fail.
	lgs, err := o1.SelectDataWordRange(addr, eventSig, 0, logpoller.EvmWord(2), logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))

	// Range including log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 0, logpoller.EvmWord(1), logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Range only covering log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 0, logpoller.EvmWord(1), logpoller.EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Cannot query for unconfirmed second log.
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	// Confirm it, then can query.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x2"), 2, time.Now()))
	lgs, err = o1.SelectDataWordRange(addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, lgs[0].Data, append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...))

	// Check greater than 1 yields both logs.
	lgs, err = o1.SelectDataWordGreaterThan(addr, eventSig, 0, logpoller.EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))
}

func TestORM_SelectLogsWithSigsByBlockRangeFilter(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM

	// Insert logs on different topics, should be able to read them
	// back using SelectLogsWithSigsByBlockRangeFilter and specifying
	// said topics.
	topic := common.HexToHash("0x1599")
	topic2 := common.HexToHash("0x1600")
	sourceAddr := common.HexToAddress("0x12345")
	inputLogs := []logpoller.Log{
		{
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
			EvmChainId:  utils.NewBig(th.ChainID),
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
	th := SetupTH(t, 2, 3, 2)
	o1 := th.ORM
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 1, time.Now()))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1235"), 2, time.Now()))
	require.NoError(t, o1.DeleteBlocksBefore(1))
	// 1 should be gone.
	_, err := o1.SelectBlockByNumber(1)
	require.Equal(t, err, sql.ErrNoRows)
	b, err := o1.SelectBlockByNumber(2)
	require.NoError(t, err)
	assert.Equal(t, int64(2), b.BlockNumber)
	// Clear multiple
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1236"), 3, time.Now()))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1237"), 4, time.Now()))
	require.NoError(t, o1.DeleteBlocksBefore(3))
	_, err = o1.SelectBlockByNumber(2)
	require.Equal(t, err, sql.ErrNoRows)
	_, err = o1.SelectBlockByNumber(3)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestLogPoller_Logs(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)
	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1-3
	require.NoError(t, th.ORM.InsertLogs([]logpoller.Log{
		GenLog(th.ChainID, 1, 1, "0x3", event1[:], address1),
		GenLog(th.ChainID, 2, 1, "0x3", event2[:], address2),
		GenLog(th.ChainID, 1, 2, "0x4", event1[:], address2),
		GenLog(th.ChainID, 2, 2, "0x4", event2[:], address1),
		GenLog(th.ChainID, 1, 3, "0x5", event1[:], address1),
		GenLog(th.ChainID, 2, 3, "0x5", event2[:], address2),
	}))

	// Select for all Addresses
	lgs, err := th.ORM.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 6, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[1].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[2].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[3].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[4].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[5].BlockHash.String())

	// Filter by Address and topic
	lgs, err = th.ORM.SelectLogsByBlockRangeFilter(1, 3, address1, event1)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	// Filter by block
	lgs, err = th.ORM.SelectLogsByBlockRangeFilter(2, 2, address2, event1)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
}

func BenchmarkLogs(b *testing.B) {
	th := SetupTH(b, 2, 3, 2)
	o := th.ORM
	var lgs []logpoller.Log
	addr := common.HexToAddress("0x1234")
	for i := 0; i < 10_000; i++ {
		lgs = append(lgs, logpoller.Log{
			EvmChainId:  utils.NewBig(th.ChainID),
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
		_, err := o.SelectDataWordRange(addr, EmitterABI.Events["Log1"].ID, 0, logpoller.EvmWord(8000), logpoller.EvmWord(8002), 0)
		require.NoError(b, err)
	}
}

func TestSelectLogsWithSigsExcluding(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	orm := th.ORM
	addressA := common.HexToAddress("0x11111")
	addressB := common.HexToAddress("0x22222")
	addressC := common.HexToAddress("0x33333")

	requestSigA := common.HexToHash("0x01")
	responseSigA := common.HexToHash("0x02")
	requestSigB := common.HexToHash("0x03")
	responseSigB := common.HexToHash("0x04")

	topicA := common.HexToHash("0x000a")
	topicB := common.HexToHash("0x000b")
	topicC := common.HexToHash("0x000c")
	topicD := common.HexToHash("0x000d")

	//Insert two logs that mimics an oracle request from 2 different addresses (matching will be on topic index 1)
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       1,
			BlockHash:      common.HexToHash("0x1"),
			BlockNumber:    1,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{requestSigA.Bytes(), topicA.Bytes(), topicB.Bytes()},
			EventSig:       requestSigA,
			Address:        addressA,
			TxHash:         common.HexToHash("0x0001"),
			Data:           []byte("requestID-A1"),
		},
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       2,
			BlockHash:      common.HexToHash("0x1"),
			BlockNumber:    1,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{requestSigB.Bytes(), topicA.Bytes(), topicB.Bytes()},
			EventSig:       requestSigB,
			Address:        addressB,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("requestID-B1"),
		},
	}))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x1"), 1, time.Now()))

	//Get any requestSigA from addressA that do not have a equivalent responseSigA
	logs, err := orm.SelectIndexedLogsWithSigsExcluding(requestSigA, responseSigA, 1, addressA, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-A1"))

	//Get any requestSigB from addressB that do not have a equivalent responseSigB
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Insert a log that mimics response for requestID-A1
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       3,
			BlockHash:      common.HexToHash("0x2"),
			BlockNumber:    2,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{responseSigA.Bytes(), topicA.Bytes(), topicC.Bytes(), topicD.Bytes()},
			EventSig:       responseSigA,
			Address:        addressA,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("responseID-A1"),
		},
	}))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x2"), 2, time.Now()))

	//Should return nothing as requestID-A1 has been fulfilled
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigA, responseSigA, 1, addressA, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//requestID-B1 should still be unfulfilled
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Insert 3 request from addressC (matching will be on topic index 3)
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       5,
			BlockHash:      common.HexToHash("0x2"),
			BlockNumber:    3,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{requestSigB.Bytes(), topicD.Bytes(), topicB.Bytes(), topicC.Bytes()},
			EventSig:       requestSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("requestID-C1"),
		},
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       6,
			BlockHash:      common.HexToHash("0x2"),
			BlockNumber:    3,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{requestSigB.Bytes(), topicD.Bytes(), topicB.Bytes(), topicA.Bytes()},
			EventSig:       requestSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("requestID-C2"),
		}, {
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       7,
			BlockHash:      common.HexToHash("0x2"),
			BlockNumber:    3,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{requestSigB.Bytes(), topicD.Bytes(), topicB.Bytes(), topicD.Bytes()},
			EventSig:       requestSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("requestID-C3"),
		},
	}))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x3"), 3, time.Now()))

	//Get all unfulfilled requests from addressC, match on topic index 3
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 3)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))
	require.Equal(t, logs[1].Data, []byte("requestID-C2"))
	require.Equal(t, logs[2].Data, []byte("requestID-C3"))

	//Fulfill requestID-C2
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       8,
			BlockHash:      common.HexToHash("0x3"),
			BlockNumber:    3,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{responseSigB.Bytes(), topicC.Bytes(), topicD.Bytes(), topicA.Bytes()},
			EventSig:       responseSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("responseID-C2"),
		},
	}))

	//Verify that requestID-C2 is now fulfilled (not returned)
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))
	require.Equal(t, logs[1].Data, []byte("requestID-C3"))

	//Fulfill requestID-C3
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       9,
			BlockHash:      common.HexToHash("0x3"),
			BlockNumber:    3,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{responseSigB.Bytes(), topicC.Bytes(), topicD.Bytes(), topicD.Bytes()},
			EventSig:       responseSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("responseID-C3"),
		},
	}))

	//Verify that requestID-C3 is now fulfilled (not returned)
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))

	//Should return no logs as the number of confirmations is not satisfied
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 4, 3)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	require.NoError(t, orm.InsertBlock(common.HexToHash("0x4"), 4, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x5"), 5, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x6"), 6, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x7"), 7, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x8"), 8, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x9"), 9, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x10"), 10, time.Now()))

	//Fulfill requestID-C3
	require.NoError(t, orm.InsertLogs([]logpoller.Log{
		{
			EvmChainId:     (*utils.Big)(th.ChainID),
			LogIndex:       10,
			BlockHash:      common.HexToHash("0x2"),
			BlockNumber:    10,
			BlockTimestamp: time.Now(),
			Topics:         [][]byte{responseSigB.Bytes(), topicD.Bytes(), topicB.Bytes(), topicC.Bytes()},
			EventSig:       responseSigB,
			Address:        addressC,
			TxHash:         common.HexToHash("0x0002"),
			Data:           []byte("responseID-C1"),
		},
	}))

	//All logs for addressC should be fulfilled, query should return 0 logs
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 10, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//Should return 1 log as it does not satisfy the required number of confirmations
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 10, 3)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))

	//Insert 3 more blocks so that the requestID-C1 has enough confirmations
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x11"), 11, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x12"), 12, time.Now()))
	require.NoError(t, orm.InsertBlock(common.HexToHash("0x13"), 13, time.Now()))

	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 3, addressC, 0, 10, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//AddressB should still have an unfulfilled log (requestID-B1)
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Should return requestID-A1 as the fulfillment event is out of the block range
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigA, responseSigA, 1, addressA, 0, 1, 10)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-A1"))

	//Should return nothing as requestID-B1 is before the block range
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(requestSigB, responseSigB, 1, addressB, 2, 13, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)
}
