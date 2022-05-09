package logpoller

import (
	"database/sql"
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

// Setup creates two orms representing logs from different chains.
func setup(t *testing.T) (*ORM, *ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	o1 := NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	o2 := NewORM(big.NewInt(138), db, lggr, pgtest.NewPGCfg(true))
	return o1, o2
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

	// Delete a block
	require.NoError(t, o1.DeleteRangeBlocks(10, 10))
	_, err = o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Delete block from another chain.
	require.NoError(t, o2.DeleteRangeBlocks(11, 11))
	_, err = o2.SelectBlockByHash(common.HexToHash("0x1234"))
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
			EventSig:    topic[:],
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
			EventSig:    topic[:],
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
			EventSig:    topic[:],
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
			EventSig:    topic[:],
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
			EventSig:    topic2[:],
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
			EventSig:    topic2[:],
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

	logs, err = o1.SelectLogsByBlockRangeFilter(10, 10, common.HexToAddress("0x1234"), topic[:])
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
	require.NoError(t, o1.DeleteRangeBlocks(10, 10))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 11))
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 12))
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 1)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 2)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogEventSigWithConfs(topic, common.HexToAddress("0x1234"), 3)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Latest log for topic for addr "0x1234" is @ block 11
	lgs, err := o1.LatestLogEventSigsAddrs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic})
	require.NoError(t, err)

	require.Equal(t, 1, len(lgs))
	require.Equal(t, int64(11), lgs[0].BlockNumber)

	// should return two entries one for each address with the latest update
	lgs, err = o1.LatestLogEventSigsAddrs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic})
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return two entries one for each topic for addr 0x1234
	lgs, err = o1.LatestLogEventSigsAddrs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic, topic2})
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return 4 entries one for each (address,topic) combination
	lgs, err = o1.LatestLogEventSigsAddrs(0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2})
	require.NoError(t, err)
	require.Equal(t, 4, len(lgs))
}

func insertLogsTopicValueRange(t *testing.T, o *ORM, addr common.Address, blockNumber int, eventSig []byte, start, stop int) {
	var lgs []Log
	for i := start; i <= stop; i++ {
		lgs = append(lgs, Log{
			EvmChainId:  utils.NewBig(o.chainID),
			LogIndex:    int64(i),
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(blockNumber),
			EventSig:    eventSig[:],
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
	bh := common.HexToHash("0x1234")
	require.NoError(t, o1.InsertBlock(bh, 1))
	insertLogsTopicValueRange(t, o1, addr, 1, eventSig.Bytes(), 1, 3)
	insertLogsTopicValueRange(t, o1, addr, 2, eventSig.Bytes(), 4, 4) // unconfirmed

	lgs, err := o1.SelectIndexedLogs(addr, eventSig[:], 1, []common.Hash{EvmWord(1)}, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, EvmWord(1).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexedLogs(addr, eventSig[:], 1, []common.Hash{EvmWord(1), EvmWord(2)}, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexLogsTopicGreaterThan(addr, eventSig[:], 1, EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig[:], 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, EvmWord(3).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig[:], 1, EvmWord(1), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(lgs))

	// Check confirmations work as expected.
	require.NoError(t, o1.InsertBlock(bh, 2))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig[:], 1, EvmWord(4), EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	require.NoError(t, o1.InsertBlock(bh, 3))
	lgs, err = o1.SelectIndexLogsTopicRange(addr, eventSig[:], 1, EvmWord(4), EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
}

func TestORM_DataWords(t *testing.T) {
	o1, _ := setup(t)
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	bh := common.HexToHash("0x1234")
	require.NoError(t, o1.InsertBlock(bh, 1))
	require.NoError(t, o1.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBig(o1.chainID),
			LogIndex:    int64(0),
			BlockHash:   bh,
			BlockNumber: int64(1),
			EventSig:    eventSig[:],
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        EvmWord(1).Bytes(),
		},
		{
			// In block 2, unconfirmed to start
			EvmChainId:  utils.NewBig(o1.chainID),
			LogIndex:    int64(1),
			BlockHash:   bh,
			BlockNumber: int64(2),
			EventSig:    eventSig[:],
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1888"),
			Data:        append(EvmWord(2).Bytes(), EvmWord(3).Bytes()...),
		},
	}))
	// Outside range should fail.
	lgs, err := o1.SelectDataWordRange(addr, eventSig[:], 0, EvmWord(2), EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))

	// Range including log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig[:], 0, EvmWord(1), EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Range only covering log should succeed
	lgs, err = o1.SelectDataWordRange(addr, eventSig[:], 0, EvmWord(1), EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	// Cannot query for unconfirmed second log.
	lgs, err = o1.SelectDataWordRange(addr, eventSig[:], 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	// Confirm it, then can query.
	require.NoError(t, o1.InsertBlock(bh, 2))
	lgs, err = o1.SelectDataWordRange(addr, eventSig[:], 1, EvmWord(3), EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, lgs[0].Data, append(EvmWord(2).Bytes(), EvmWord(3).Bytes()...))

	// Check greater than 1 yields both logs.
	lgs, err = o1.SelectDataWordGreaterThan(addr, eventSig[:], 0, EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))
}
