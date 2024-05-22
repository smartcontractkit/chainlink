package logpoller_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type block struct {
	number    int64
	hash      common.Hash
	timestamp int64
}

var lpOpts = logpoller.Opts{
	FinalityDepth:            2,
	BackfillBatchSize:        3,
	RpcBatchSize:             2,
	KeepFinalizedBlocksDepth: 1000,
}

func GenLog(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address) logpoller.Log {
	return GenLogWithTimestamp(chainID, logIndex, blockNum, blockHash, topic1, address, time.Now())
}

func GenLogWithTimestamp(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address, blockTimestamp time.Time) logpoller.Log {
	return logpoller.Log{
		EvmChainId:     ubig.New(chainID),
		LogIndex:       logIndex,
		BlockHash:      common.HexToHash(blockHash),
		BlockNumber:    blockNum,
		EventSig:       common.BytesToHash(topic1),
		Topics:         [][]byte{topic1, topic1},
		Address:        address,
		TxHash:         common.HexToHash("0x1234"),
		Data:           append([]byte("hello "), byte(blockNum)),
		BlockTimestamp: blockTimestamp,
	}
}

func GenLogWithData(chainID *big.Int, address common.Address, eventSig common.Hash, logIndex int64, blockNum int64, data []byte) logpoller.Log {
	return logpoller.Log{
		EvmChainId:     ubig.New(chainID),
		LogIndex:       logIndex,
		BlockHash:      utils.RandomBytes32(),
		BlockNumber:    blockNum,
		EventSig:       eventSig,
		Topics:         [][]byte{},
		Address:        address,
		TxHash:         utils.RandomBytes32(),
		Data:           data,
		BlockTimestamp: time.Now(),
	}
}

func TestLogPoller_Batching(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	th := SetupTH(t, lpOpts)
	var logs []logpoller.Log
	// Inserts are limited to 65535 parameters. A log being 10 parameters this results in
	// a maximum of 6553 log inserts per tx. As inserting more than 6553 would result in
	// an error without batching, this test makes sure batching is enabled.
	for i := 0; i < 15000; i++ {
		logs = append(logs, GenLog(th.ChainID, int64(i+1), 1, "0x3", EmitterABI.Events["Log1"].ID.Bytes(), th.EmitterAddress1))
	}
	require.NoError(t, th.ORM.InsertLogs(ctx, logs))
	lgs, err := th.ORM.SelectLogsByBlockRange(ctx, 1, 1)
	require.NoError(t, err)
	// Make sure all logs are inserted
	require.Equal(t, len(logs), len(lgs))
}

func TestORM_GetBlocks_From_Range(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
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
		require.NoError(t, o1.InsertBlock(ctx, b.hash, b.number, time.Unix(b.timestamp, 0).UTC(), 0))
	}

	var blockNumbers []int64
	for _, b := range blocks {
		blockNumbers = append(blockNumbers, b.number)
	}

	lpBlocks, err := o1.GetBlocksRange(ctx, blockNumbers[0], blockNumbers[len(blockNumbers)-1])
	require.NoError(t, err)
	assert.Len(t, lpBlocks, len(blocks))

	// Ignores non-existent block
	lpBlocks2, err := o1.GetBlocksRange(ctx, blockNumbers[0], 15)
	require.NoError(t, err)
	assert.Len(t, lpBlocks2, len(blocks))

	// Only non-existent blocks
	lpBlocks3, err := o1.GetBlocksRange(ctx, 15, 15)
	require.NoError(t, err)
	assert.Len(t, lpBlocks3, 0)
}

func TestORM_GetBlocks_From_Range_Recent_Blocks(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
	// Insert many blocks and read them back together
	var recentBlocks []block
	for i := 1; i <= 256; i++ {
		recentBlocks = append(recentBlocks, block{number: int64(i), hash: common.HexToHash(fmt.Sprintf("0x%d", i))})
	}
	for _, b := range recentBlocks {
		require.NoError(t, o1.InsertBlock(ctx, b.hash, b.number, time.Now(), 0))
	}

	var blockNumbers []int64
	for _, b := range recentBlocks {
		blockNumbers = append(blockNumbers, b.number)
	}

	lpBlocks, err := o1.GetBlocksRange(ctx, blockNumbers[0], blockNumbers[len(blockNumbers)-1])
	require.NoError(t, err)
	assert.Len(t, lpBlocks, len(recentBlocks))

	// Ignores non-existent block
	lpBlocks2, err := o1.GetBlocksRange(ctx, blockNumbers[0], 257)
	require.NoError(t, err)
	assert.Len(t, lpBlocks2, len(recentBlocks))

	// Only non-existent blocks
	lpBlocks3, err := o1.GetBlocksRange(ctx, 257, 257)
	require.NoError(t, err)
	assert.Len(t, lpBlocks3, 0)
}

func TestORM(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	o2 := th.ORM2
	ctx := testutils.Context(t)
	// Insert and read back a block.
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 10, time.Now(), 0))
	b, err := o1.SelectBlockByHash(ctx, common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b.BlockNumber, int64(10))
	assert.Equal(t, b.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b.EvmChainId.String(), th.ChainID.String())

	// Insert blocks from a different chain
	require.NoError(t, o2.InsertBlock(ctx, common.HexToHash("0x1234"), 11, time.Now(), 0))
	require.NoError(t, o2.InsertBlock(ctx, common.HexToHash("0x1235"), 12, time.Now(), 0))
	b2, err := o2.SelectBlockByHash(ctx, common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b2.BlockNumber, int64(11))
	assert.Equal(t, b2.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b2.EvmChainId.String(), th.ChainID2.String())

	latest, err := o1.SelectLatestBlock(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(10), latest.BlockNumber)

	latest, err = o2.SelectLatestBlock(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(12), latest.BlockNumber)

	// Delete a block (only 10 on chain).
	require.NoError(t, o1.DeleteLogsAndBlocksAfter(ctx, 10))
	_, err = o1.SelectBlockByHash(ctx, common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))

	// Delete blocks from another chain.
	require.NoError(t, o2.DeleteLogsAndBlocksAfter(ctx, 11))
	_, err = o2.SelectBlockByHash(ctx, common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))
	// Delete blocks after should also delete block 12.
	_, err = o2.SelectBlockByHash(ctx, common.HexToHash("0x1235"))
	require.Error(t, err)
	assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))

	// Should be able to insert and read back a log.
	topic := common.HexToHash("0x1599")
	topic2 := common.HexToHash("0x1600")
	require.NoError(t, o1.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       1,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(10),
			EventSig:       topic,
			Topics:         [][]byte{topic[:]},
			Address:        common.HexToAddress("0x1234"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       2,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(11),
			EventSig:       topic,
			Topics:         [][]byte{topic[:]},
			Address:        common.HexToAddress("0x1234"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       3,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(12),
			EventSig:       topic,
			Topics:         [][]byte{topic[:]},
			Address:        common.HexToAddress("0x1235"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       4,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(13),
			EventSig:       topic,
			Topics:         [][]byte{topic[:]},
			Address:        common.HexToAddress("0x1235"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       5,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(14),
			EventSig:       topic2,
			Topics:         [][]byte{topic2[:]},
			Address:        common.HexToAddress("0x1234"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello2"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       6,
			BlockHash:      common.HexToHash("0x1234"),
			BlockNumber:    int64(15),
			EventSig:       topic2,
			Topics:         [][]byte{topic2[:]},
			Address:        common.HexToAddress("0x1235"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello2"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       7,
			BlockHash:      common.HexToHash("0x1237"),
			BlockNumber:    int64(16),
			EventSig:       topic,
			Topics:         [][]byte{topic[:]},
			Address:        common.HexToAddress("0x1236"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello short retention"),
			BlockTimestamp: time.Now(),
		},
		{
			EvmChainId:     ubig.New(th.ChainID),
			LogIndex:       8,
			BlockHash:      common.HexToHash("0x1238"),
			BlockNumber:    int64(17),
			EventSig:       topic2,
			Topics:         [][]byte{topic2[:]},
			Address:        common.HexToAddress("0x1236"),
			TxHash:         common.HexToHash("0x1888"),
			Data:           []byte("hello2 long retention"),
			BlockTimestamp: time.Now(),
		},
	}))

	t.Log(latest.BlockNumber)
	logs, err := o1.SelectLogsByBlockRange(ctx, 1, 17)
	require.NoError(t, err)
	require.Len(t, logs, 8)

	logs, err = o1.SelectLogsByBlockRange(ctx, 10, 10)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
	assert.Equal(t, []byte("hello"), logs[0].Data)

	logs, err = o1.SelectLogs(ctx, 1, 1, common.HexToAddress("0x1234"), topic)
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))
	logs, err = o1.SelectLogs(ctx, 10, 10, common.HexToAddress("0x1234"), topic)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))

	// With no blocks, should be an error
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 0)
	require.Error(t, err)
	require.True(t, pkgerrors.Is(err, sql.ErrNoRows))
	// With block 10, only 0 confs should work
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 10, time.Now(), 0))
	log, err := o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	assert.Equal(t, int64(10), log.BlockNumber)
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 1)
	require.Error(t, err)
	assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))
	// With block 12, anything <=2 should work
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 11, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1235"), 12, time.Now(), 0))
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 0)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 1)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 2)
	require.NoError(t, err)
	_, err = o1.SelectLatestLogByEventSigWithConfs(ctx, topic, common.HexToAddress("0x1234"), 3)
	require.Error(t, err)
	assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))

	// Required for confirmations to work
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 13, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1235"), 14, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1236"), 15, time.Now(), 0))

	// Latest log for topic for addr "0x1234" is @ block 11
	lgs, err := o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic}, 0)
	require.NoError(t, err)

	require.Equal(t, 1, len(lgs))
	require.Equal(t, int64(11), lgs[0].BlockNumber)

	// should return two entries one for each address with the latest update
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic}, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return two entries one for each topic for addr 0x1234
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234")}, []common.Hash{topic, topic2}, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	// should return 4 entries one for each (address,topic) combination
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 0)
	require.NoError(t, err)
	require.Equal(t, 4, len(lgs))

	// should return 3 entries of logs with atleast 1 confirmation
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 1)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))

	// should return 2 entries of logs with atleast 2 confirmation
	lgs, err = o1.SelectLatestLogEventSigsAddrsWithConfs(ctx, 0 /* startBlock */, []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0x1235")}, []common.Hash{topic, topic2}, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))

	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1237"), 16, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1238"), 17, time.Now(), 0))

	filter0 := logpoller.Filter{
		Name:      "permanent retention filter",
		Addresses: []common.Address{common.HexToAddress("0x1234")},
		EventSigs: types.HashArray{topic, topic2},
	}

	filter12 := logpoller.Filter{ // retain both topic1 and topic2 on contract3 for at least 1ms
		Name:      "short retention filter",
		Addresses: []common.Address{common.HexToAddress("0x1236")},
		EventSigs: types.HashArray{topic, topic2},
		Retention: time.Millisecond,
	}
	filter2 := logpoller.Filter{ // retain topic2 on contract3 for at least 1 hour
		Name:      "long retention filter",
		Addresses: []common.Address{common.HexToAddress("0x1236")},
		EventSigs: types.HashArray{topic2},
		Retention: time.Hour,
	}

	// Test inserting filters and reading them back
	require.NoError(t, o1.InsertFilter(ctx, filter0))
	require.NoError(t, o1.InsertFilter(ctx, filter12))
	require.NoError(t, o1.InsertFilter(ctx, filter2))

	filters, err := o1.LoadFilters(ctx)
	require.NoError(t, err)
	require.Len(t, filters, 3)
	assert.Equal(t, filter0, filters["permanent retention filter"])
	assert.Equal(t, filter12, filters["short retention filter"])
	assert.Equal(t, filter2, filters["long retention filter"])

	latest, err = o1.SelectLatestBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(17), latest.BlockNumber)
	logs, err = o1.SelectLogsByBlockRange(ctx, 1, latest.BlockNumber)
	require.NoError(t, err)
	require.Len(t, logs, 8)

	// Delete expired logs
	time.Sleep(2 * time.Millisecond) // just in case we haven't reached the end of the 1ms retention period
	deleted, err := o1.DeleteExpiredLogs(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted)
	logs, err = o1.SelectLogsByBlockRange(ctx, 1, latest.BlockNumber)
	require.NoError(t, err)
	// The only log which should be deleted is the one which matches filter1 (ret=1ms) but not filter12 (ret=1 hour)
	// Importantly, it shouldn't delete any logs matching only filter0 (ret=0 meaning permanent retention).  Anything
	// matching filter12 should be kept regardless of what other filters it matches.
	assert.Len(t, logs, 7)

	// Delete logs after should delete all logs.
	err = o1.DeleteLogsAndBlocksAfter(ctx, 1)
	require.NoError(t, err)
	logs, err = o1.SelectLogsByBlockRange(ctx, 1, latest.BlockNumber)
	require.NoError(t, err)
	require.Zero(t, len(logs))
}

type PgxLogger struct {
	lggr logger.Logger
}

func NewPgxLogger(lggr logger.Logger) PgxLogger {
	return PgxLogger{lggr}
}

func (l PgxLogger) Log(ctx context.Context, log pgx.LogLevel, msg string, data map[string]interface{}) {

}

func TestLogPollerFilters(t *testing.T) {
	lggr := logger.Test(t)
	chainID := testutils.NewRandomEVMChainID()

	dbx := pgtest.NewSqlxDB(t)
	orm := logpoller.NewORM(chainID, dbx, lggr)

	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address := common.HexToAddress("0x1234")
	topicA := common.HexToHash("0x1111")
	topicB := common.HexToHash("0x2222")
	topicC := common.HexToHash("0x3333")
	topicD := common.HexToHash("0x4444")

	ctx := testutils.Context(t)

	filters := []logpoller.Filter{{
		Name:      "filter by topic2",
		EventSigs: types.HashArray{event1, event2},
		Addresses: types.AddressArray{address},
		Topic2:    types.HashArray{topicA, topicB},
	}, {
		Name:      "filter by topic3",
		Addresses: types.AddressArray{address},
		EventSigs: types.HashArray{event1},
		Topic3:    types.HashArray{topicB, topicC, topicD},
	}, {
		Name:      "filter by topic4",
		Addresses: types.AddressArray{address},
		EventSigs: types.HashArray{event1},
		Topic4:    types.HashArray{topicC},
	}, {
		Name:      "filter by topics 2 and 4",
		Addresses: types.AddressArray{address},
		EventSigs: types.HashArray{event2},
		Topic2:    types.HashArray{topicA},
		Topic4:    types.HashArray{topicC, topicD},
	}, {
		Name:         "10 lpb rate limit, 1M max logs",
		Addresses:    types.AddressArray{address},
		EventSigs:    types.HashArray{event1},
		MaxLogsKept:  1000000,
		LogsPerBlock: 10,
	}, { // ensure that the UNIQUE CONSTRAINT isn't too strict (should only error if all fields are identical)
		Name:      "duplicate of filter by topic4",
		Addresses: types.AddressArray{address},
		EventSigs: types.HashArray{event1},
		Topic3:    types.HashArray{topicC},
	}}

	for _, filter := range filters {
		t.Run("Save filter: "+filter.Name, func(t *testing.T) {
			var count int
			err := orm.InsertFilter(ctx, filter)
			require.NoError(t, err)
			err = dbx.Get(&count, `SELECT COUNT(*) FROM evm.log_poller_filters WHERE evm_chain_id = $1 AND name = $2`, ubig.New(chainID), filter.Name)
			require.NoError(t, err)
			expectedCount := len(filter.Addresses) * len(filter.EventSigs)
			if len(filter.Topic2) > 0 {
				expectedCount *= len(filter.Topic2)
			}
			if len(filter.Topic3) > 0 {
				expectedCount *= len(filter.Topic3)
			}
			if len(filter.Topic4) > 0 {
				expectedCount *= len(filter.Topic4)
			}
			assert.Equal(t, count, expectedCount)
		})
	}

	// Make sure they all come back the same when we reload them
	t.Run("Load filters", func(t *testing.T) {
		loadedFilters, err := orm.LoadFilters(ctx)
		require.NoError(t, err)
		for _, filter := range filters {
			loadedFilter, ok := loadedFilters[filter.Name]
			require.True(t, ok, `Failed to reload filter "%s"`, filter.Name)
			assert.Equal(t, filter, loadedFilter)
		}
	})
}

func insertLogsTopicValueRange(t *testing.T, chainID *big.Int, o logpoller.ORM, addr common.Address, blockNumber int, eventSig common.Hash, start, stop int) {
	var lgs []logpoller.Log
	for i := start; i <= stop; i++ {
		lgs = append(lgs, logpoller.Log{
			EvmChainId:  ubig.New(chainID),
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
	require.NoError(t, o.InsertLogs(testutils.Context(t), lgs))
}

func TestORM_IndexedLogs(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1"), 1, time.Now(), 0))
	insertLogsTopicValueRange(t, th.ChainID, o1, addr, 1, eventSig, 1, 3)
	insertLogsTopicValueRange(t, th.ChainID, o1, addr, 2, eventSig, 4, 4) // unconfirmed

	filtersForTopics := func(topicIdx uint64, topicValues []uint64) query.Expression {
		topicFilters := query.BoolExpression{
			Expressions:  make([]query.Expression, len(topicValues)),
			BoolOperator: query.OR,
		}

		for idx, value := range topicValues {
			topicFilters.Expressions[idx] = logpoller.NewEventByTopicFilter(topicIdx, []primitives.ValueComparator{
				{Value: logpoller.EvmWord(value).Hex(), Operator: primitives.Eq},
			})
		}

		return query.Expression{BoolExpression: topicFilters}
	}

	limiter := query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))
	standardFilter := func(topicIdx uint64, topicValues []uint64) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(addr),
				logpoller.NewEventSigFilter(eventSig),
				filtersForTopics(topicIdx, topicValues),
				query.Confirmation(primitives.Unconfirmed),
			},
		}
	}

	lgs, err := o1.SelectIndexedLogs(ctx, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)}, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(1).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.FilteredLogs(ctx, standardFilter(1, []uint64{1}), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(1).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexedLogs(ctx, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1), logpoller.EvmWord(2)}, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, standardFilter(1, []uint64{1, 2}), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	blockRangeFilter := func(start, end uint64, topicIdx uint64, topicValues []uint64) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(addr),
				logpoller.NewEventSigFilter(eventSig),
				filtersForTopics(topicIdx, topicValues),
				query.Block(start, primitives.Gte),
				query.Block(end, primitives.Lte),
			},
		}
	}

	lgs, err = o1.SelectIndexedLogsByBlockRange(ctx, 1, 1, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, blockRangeFilter(1, 1, 1, []uint64{1}), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.SelectIndexedLogsByBlockRange(ctx, 1, 2, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(2)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, blockRangeFilter(1, 2, 1, []uint64{2}), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.SelectIndexedLogsByBlockRange(ctx, 1, 2, addr, eventSig, 1, []common.Hash{logpoller.EvmWord(1)})
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, blockRangeFilter(1, 2, 1, []uint64{1}), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))

	_, err = o1.SelectIndexedLogsByBlockRange(ctx, 1, 2, addr, eventSig, 0, []common.Hash{logpoller.EvmWord(1)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 0")

	_, err = o1.FilteredLogs(ctx, blockRangeFilter(1, 2, 0, []uint64{1}), limiter, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 0")

	_, err = o1.SelectIndexedLogsByBlockRange(ctx, 1, 2, addr, eventSig, 4, []common.Hash{logpoller.EvmWord(1)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 4")

	_, err = o1.FilteredLogs(ctx, blockRangeFilter(1, 2, 4, []uint64{1}), limiter, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index for topic: 4")

	lgs, err = o1.SelectIndexedLogsTopicGreaterThan(ctx, addr, eventSig, 1, logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	filter := query.KeyFilter{
		Expressions: []query.Expression{
			logpoller.NewAddressFilter(addr),
			logpoller.NewEventSigFilter(eventSig),
			logpoller.NewEventByTopicFilter(1, []primitives.ValueComparator{
				{Value: logpoller.EvmWord(2).Hex(), Operator: primitives.Gte},
			}),
			query.Confirmation(primitives.Unconfirmed),
		},
	}

	lgs, err = o1.FilteredLogs(ctx, filter, limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	rangeFilter := func(topicIdx uint64, min, max uint64) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(addr),
				logpoller.NewEventSigFilter(eventSig),
				logpoller.NewEventByTopicFilter(topicIdx, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(min).Hex(), Operator: primitives.Gte},
				}),
				logpoller.NewEventByTopicFilter(topicIdx, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(max).Hex(), Operator: primitives.Lte},
				}),
				query.Confirmation(primitives.Unconfirmed),
			},
		}
	}

	lgs, err = o1.SelectIndexedLogsTopicRange(ctx, addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(3).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.FilteredLogs(ctx, rangeFilter(1, 3, 3), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
	assert.Equal(t, logpoller.EvmWord(3).Bytes(), lgs[0].GetTopics()[1].Bytes())

	lgs, err = o1.SelectIndexedLogsTopicRange(ctx, addr, eventSig, 1, logpoller.EvmWord(1), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, rangeFilter(1, 1, 3), limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 3, len(lgs))

	// Check confirmations work as expected.
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x2"), 2, time.Now(), 0))
	lgs, err = o1.SelectIndexedLogsTopicRange(ctx, addr, eventSig, 1, logpoller.EvmWord(4), logpoller.EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x3"), 3, time.Now(), 0))
	lgs, err = o1.SelectIndexedLogsTopicRange(ctx, addr, eventSig, 1, logpoller.EvmWord(4), logpoller.EvmWord(4), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lgs))
}

func TestORM_SelectIndexedLogsByTxHash(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
	eventSig := common.HexToHash("0x1599")
	txHash := common.HexToHash("0x1888")
	addr := common.HexToAddress("0x1234")

	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1"), 1, time.Now(), 0))
	logs := []logpoller.Log{
		{
			EvmChainId:  ubig.New(th.ChainID),
			LogIndex:    int64(0),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      txHash,
			Data:        logpoller.EvmWord(1).Bytes(),
		},
		{
			EvmChainId:  ubig.New(th.ChainID),
			LogIndex:    int64(1),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      txHash,
			Data:        append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...),
		},
		// Different txHash
		{
			EvmChainId:  ubig.New(th.ChainID),
			LogIndex:    int64(2),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    eventSig,
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      common.HexToHash("0x1889"),
			Data:        append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...),
		},
		// Different eventSig
		{
			EvmChainId:  ubig.New(th.ChainID),
			LogIndex:    int64(3),
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			EventSig:    common.HexToHash("0x1600"),
			Topics:      [][]byte{eventSig[:]},
			Address:     addr,
			TxHash:      txHash,
			Data:        append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...),
		},
	}
	require.NoError(t, o1.InsertLogs(ctx, logs))

	retrievedLogs, err := o1.SelectIndexedLogsByTxHash(ctx, addr, eventSig, txHash)
	require.NoError(t, err)

	require.Equal(t, 2, len(retrievedLogs))
	require.Equal(t, retrievedLogs[0].LogIndex, logs[0].LogIndex)
	require.Equal(t, retrievedLogs[1].LogIndex, logs[1].LogIndex)

	limiter := query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))
	filter := query.KeyFilter{
		Expressions: []query.Expression{
			logpoller.NewAddressFilter(addr),
			logpoller.NewEventSigFilter(eventSig),
			query.TxHash(txHash.Hex()),
		},
	}

	retrievedLogs, err = o1.FilteredLogs(ctx, filter, limiter, "")
	require.NoError(t, err)

	require.Equal(t, 2, len(retrievedLogs))
	require.Equal(t, retrievedLogs[0].LogIndex, logs[0].LogIndex)
	require.Equal(t, retrievedLogs[1].LogIndex, logs[1].LogIndex)
}

func TestORM_DataWords(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
	eventSig := common.HexToHash("0x1599")
	addr := common.HexToAddress("0x1234")
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1"), 1, time.Now(), 0))
	require.NoError(t, o1.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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

	wordFilter := func(wordIdx uint8, word1, word2 uint64) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(addr),
				logpoller.NewEventSigFilter(eventSig),
				logpoller.NewEventByWordFilter(eventSig, wordIdx, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(word1).Hex(), Operator: primitives.Gte},
				}),
				logpoller.NewEventByWordFilter(eventSig, wordIdx, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(word2).Hex(), Operator: primitives.Lte},
				}),
				query.Confirmation(primitives.Unconfirmed),
			},
		}
	}

	limiter := query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))

	// Outside range should fail.
	lgs, err := o1.SelectLogsDataWordRange(ctx, addr, eventSig, 0, logpoller.EvmWord(2), logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, wordFilter(0, 2, 2), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 0, len(lgs))

	// Range including log should succeed
	lgs, err = o1.SelectLogsDataWordRange(ctx, addr, eventSig, 0, logpoller.EvmWord(1), logpoller.EvmWord(2), 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, wordFilter(0, 1, 2), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))

	// Range only covering log should succeed
	lgs, err = o1.SelectLogsDataWordRange(ctx, addr, eventSig, 0, logpoller.EvmWord(1), logpoller.EvmWord(1), 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, wordFilter(0, 1, 1), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))

	// Cannot query for unconfirmed second log.
	lgs, err = o1.SelectLogsDataWordRange(ctx, addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(lgs))

	lgs, err = o1.FilteredLogs(ctx, wordFilter(1, 3, 3), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 0, len(lgs))

	// Confirm it, then can query.
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x2"), 2, time.Now(), 0))
	lgs, err = o1.SelectLogsDataWordRange(ctx, addr, eventSig, 1, logpoller.EvmWord(3), logpoller.EvmWord(3), 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	require.Equal(t, lgs[0].Data, append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...))

	lgs, err = o1.FilteredLogs(ctx, wordFilter(1, 3, 3), limiter, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	require.Equal(t, lgs[0].Data, append(logpoller.EvmWord(2).Bytes(), logpoller.EvmWord(3).Bytes()...))

	// Check greater than 1 yields both logs.
	lgs, err = o1.SelectLogsDataWordGreaterThan(ctx, addr, eventSig, 0, logpoller.EvmWord(1), 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))

	filter := query.KeyFilter{
		Expressions: []query.Expression{
			logpoller.NewAddressFilter(addr),
			logpoller.NewEventSigFilter(eventSig),
			logpoller.NewEventByWordFilter(eventSig, 0, []primitives.ValueComparator{
				{Value: logpoller.EvmWord(1).Hex(), Operator: primitives.Gte},
			}),
			query.Confirmation(primitives.Unconfirmed),
		},
	}

	lgs, err = o1.FilteredLogs(ctx, filter, limiter, "")
	require.NoError(t, err)
	assert.Equal(t, 2, len(lgs))
}

func TestORM_SelectLogsWithSigsByBlockRangeFilter(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)

	// Insert logs on different topics, should be able to read them
	// back using SelectLogsWithSigs and specifying
	// said topics.
	topic := common.HexToHash("0x1599")
	topic2 := common.HexToHash("0x1600")
	sourceAddr := common.HexToAddress("0x12345")
	inputLogs := []logpoller.Log{
		{
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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
			EvmChainId:  ubig.New(th.ChainID),
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
	require.NoError(t, o1.InsertLogs(ctx, inputLogs))

	filter := func(sigs []common.Hash, startBlock, endBlock int64) query.KeyFilter {
		filters := []query.Expression{
			logpoller.NewAddressFilter(sourceAddr),
		}

		if len(sigs) > 0 {
			exp := make([]query.Expression, len(sigs))
			for idx, val := range sigs {
				exp[idx] = logpoller.NewEventSigFilter(val)
			}

			filters = append(filters, query.Expression{
				BoolExpression: query.BoolExpression{
					Expressions:  exp,
					BoolOperator: query.OR,
				},
			})
		}

		filters = append(filters, query.Expression{
			BoolExpression: query.BoolExpression{
				Expressions: []query.Expression{
					query.Block(uint64(startBlock), primitives.Gte),
					query.Block(uint64(endBlock), primitives.Lte),
				},
				BoolOperator: query.AND,
			},
		})

		return query.KeyFilter{
			Expressions: filters,
		}
	}

	limiter := query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortBySequence(query.Asc)},
	}

	assertion := func(t *testing.T, logs []logpoller.Log, err error, startBlock, endBlock int64) {
		require.NoError(t, err)
		assert.Len(t, logs, 4)
		for _, l := range logs {
			assert.Equal(t, sourceAddr, l.Address, "wrong log address")
			assert.True(t, bytes.Equal(topic.Bytes(), l.EventSig.Bytes()) || bytes.Equal(topic2.Bytes(), l.EventSig.Bytes()), "wrong log topic")
			assert.True(t, l.BlockNumber >= startBlock && l.BlockNumber <= endBlock)
		}
	}

	startBlock, endBlock := int64(10), int64(15)
	logs, err := o1.SelectLogsWithSigs(ctx, startBlock, endBlock, sourceAddr, []common.Hash{
		topic,
		topic2,
	})

	assertion(t, logs, err, startBlock, endBlock)

	logs, err = th.ORM.FilteredLogs(ctx, filter([]common.Hash{topic, topic2}, startBlock, endBlock), limiter, "")

	assertion(t, logs, err, startBlock, endBlock)
}

func TestORM_DeleteBlocksBefore(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	ctx := testutils.Context(t)
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 1, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1235"), 2, time.Now(), 0))
	deleted, err := o1.DeleteBlocksBefore(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
	// 1 should be gone.
	_, err = o1.SelectBlockByNumber(ctx, 1)
	require.Equal(t, err, sql.ErrNoRows)
	b, err := o1.SelectBlockByNumber(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(2), b.BlockNumber)
	// Clear multiple
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1236"), 3, time.Now(), 0))
	require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1237"), 4, time.Now(), 0))
	deleted, err = o1.DeleteBlocksBefore(ctx, 3, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), deleted)
	_, err = o1.SelectBlockByNumber(ctx, 2)
	require.Equal(t, err, sql.ErrNoRows)
	_, err = o1.SelectBlockByNumber(ctx, 3)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestLogPoller_Logs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	th := SetupTH(t, lpOpts)
	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1-3
	require.NoError(t, th.ORM.InsertLogs(ctx, []logpoller.Log{
		GenLog(th.ChainID, 1, 1, "0x3", event1[:], address1),
		GenLog(th.ChainID, 2, 1, "0x3", event2[:], address2),
		GenLog(th.ChainID, 1, 2, "0x4", event1[:], address2),
		GenLog(th.ChainID, 2, 2, "0x4", event2[:], address1),
		GenLog(th.ChainID, 1, 3, "0x5", event1[:], address1),
		GenLog(th.ChainID, 2, 3, "0x5", event2[:], address2),
	}))

	// Select for all Addresses
	lgs, err := th.ORM.SelectLogsByBlockRange(ctx, 1, 3)
	require.NoError(t, err)
	require.Equal(t, 6, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[1].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[2].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[3].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[4].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[5].BlockHash.String())

	logFilter := func(start, end uint64, address common.Address) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(address),
				logpoller.NewEventSigFilter(event1),
				query.Block(start, primitives.Gte),
				query.Block(end, primitives.Lte),
			},
		}
	}

	// Filter by Address and topic
	lgs, err = th.ORM.SelectLogs(ctx, 1, 3, address1, event1)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	lgs, err = th.ORM.FilteredLogs(ctx, logFilter(1, 3, address1), query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortBySequence(query.Asc)},
	}, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	// Filter by block
	lgs, err = th.ORM.SelectLogs(ctx, 2, 2, address2, event1)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])

	lgs, err = th.ORM.FilteredLogs(ctx, logFilter(2, 2, address2), query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortBySequence(query.Asc)},
	}, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
}

func BenchmarkLogs(b *testing.B) {
	th := SetupTH(b, lpOpts)
	o := th.ORM
	ctx := testutils.Context(b)
	var lgs []logpoller.Log
	addr := common.HexToAddress("0x1234")
	for i := 0; i < 10_000; i++ {
		lgs = append(lgs, logpoller.Log{
			EvmChainId:  ubig.New(th.ChainID),
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
	require.NoError(b, o.InsertLogs(ctx, lgs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lgs, err := o.SelectLogsDataWordRange(ctx, addr, EmitterABI.Events["Log1"].ID, 0, logpoller.EvmWord(8000), logpoller.EvmWord(8002), 0)
		require.NoError(b, err)
		// TODO: Why is SelectLogsDataWordRange not returning any logs?!
		fmt.Println("len logs:", len(lgs))
	}
}

func TestSelectLogsWithSigsExcluding(t *testing.T) {
	th := SetupTH(t, lpOpts)
	orm := th.ORM
	ctx := testutils.Context(t)
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
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x1"), 1, time.Now(), 0))

	//Get any requestSigA from addressA that do not have a equivalent responseSigA
	logs, err := orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigA, responseSigA, 1, addressA, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-A1"))

	//Get any requestSigB from addressB that do not have a equivalent responseSigB
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Insert a log that mimics response for requestID-A1
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x2"), 2, time.Now(), 0))

	//Should return nothing as requestID-A1 has been fulfilled
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigA, responseSigA, 1, addressA, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//requestID-B1 should still be unfulfilled
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Insert 3 request from addressC (matching will be on topic index 3)
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x3"), 3, time.Now(), 0))

	//Get all unfulfilled requests from addressC, match on topic index 3
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 3)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))
	require.Equal(t, logs[1].Data, []byte("requestID-C2"))
	require.Equal(t, logs[2].Data, []byte("requestID-C3"))

	//Fulfill requestID-C2
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))
	require.Equal(t, logs[1].Data, []byte("requestID-C3"))

	//Fulfill requestID-C3
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 4, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))

	//Should return no logs as the number of confirmations is not satisfied
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 4, 3)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x4"), 4, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x5"), 5, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x6"), 6, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x7"), 7, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x8"), 8, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x9"), 9, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x10"), 10, time.Now(), 0))

	//Fulfill requestID-C3
	require.NoError(t, orm.InsertLogs(ctx, []logpoller.Log{
		{
			EvmChainId:     (*ubig.Big)(th.ChainID),
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
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 10, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//Should return 1 log as it does not satisfy the required number of confirmations
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 10, 3)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-C1"))

	//Insert 3 more blocks so that the requestID-C1 has enough confirmations
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x11"), 11, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x12"), 12, time.Now(), 0))
	require.NoError(t, orm.InsertBlock(ctx, common.HexToHash("0x13"), 13, time.Now(), 0))

	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 3, addressC, 0, 10, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	//AddressB should still have an unfulfilled log (requestID-B1)
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 1, addressB, 0, 3, 0)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-B1"))

	//Should return requestID-A1 as the fulfillment event is out of the block range
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigA, responseSigA, 1, addressA, 0, 1, 10)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0].Data, []byte("requestID-A1"))

	//Should return nothing as requestID-B1 is before the block range
	logs, err = orm.SelectIndexedLogsWithSigsExcluding(ctx, requestSigB, responseSigB, 1, addressB, 2, 13, 0)
	require.NoError(t, err)
	require.Len(t, logs, 0)
}

func TestSelectLatestBlockNumberEventSigsAddrsWithConfs(t *testing.T) {
	ctx := testutils.Context(t)
	th := SetupTH(t, lpOpts)
	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address1 := utils.RandomAddress()
	address2 := utils.RandomAddress()

	require.NoError(t, th.ORM.InsertLogs(ctx, []logpoller.Log{
		GenLog(th.ChainID, 1, 1, utils.RandomAddress().String(), event1[:], address1),
		GenLog(th.ChainID, 2, 1, utils.RandomAddress().String(), event2[:], address2),
		GenLog(th.ChainID, 2, 2, utils.RandomAddress().String(), event2[:], address2),
		GenLog(th.ChainID, 2, 3, utils.RandomAddress().String(), event2[:], address2),
	}))
	require.NoError(t, th.ORM.InsertBlock(ctx, utils.RandomHash(), 3, time.Now(), 1))

	tests := []struct {
		name                string
		events              []common.Hash
		addrs               []common.Address
		confs               evmtypes.Confirmations
		fromBlock           int64
		expectedBlockNumber int64
	}{
		{
			name:                "no matching logs returns 0 block number",
			events:              []common.Hash{event2},
			addrs:               []common.Address{address1},
			confs:               0,
			fromBlock:           0,
			expectedBlockNumber: 0,
		},
		{
			name:                "not enough confirmations block returns 0 block number",
			events:              []common.Hash{event2},
			addrs:               []common.Address{address2},
			confs:               5,
			fromBlock:           0,
			expectedBlockNumber: 0,
		},
		{
			name:                "single matching event and address returns last block",
			events:              []common.Hash{event1},
			addrs:               []common.Address{address1},
			confs:               0,
			fromBlock:           0,
			expectedBlockNumber: 1,
		},
		{
			name:                "only finalized log is picked",
			events:              []common.Hash{event1, event2},
			addrs:               []common.Address{address1, address2},
			confs:               evmtypes.Finalized,
			fromBlock:           0,
			expectedBlockNumber: 1,
		},
		{
			name:                "picks max block from two events",
			events:              []common.Hash{event1, event2},
			addrs:               []common.Address{address1, address2},
			confs:               0,
			fromBlock:           0,
			expectedBlockNumber: 3,
		},
		{
			name:                "picks previous block number for confirmations set to 1",
			events:              []common.Hash{event2},
			addrs:               []common.Address{address2},
			confs:               1,
			fromBlock:           0,
			expectedBlockNumber: 2,
		},
		{
			name:                "returns 0 if from block is not matching",
			events:              []common.Hash{event1, event2},
			addrs:               []common.Address{address1, address2},
			confs:               0,
			fromBlock:           3,
			expectedBlockNumber: 0,
		},
		{
			name:                "picks max block from two events when from block is lower",
			events:              []common.Hash{event1, event2},
			addrs:               []common.Address{address1, address2},
			confs:               0,
			fromBlock:           2,
			expectedBlockNumber: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := th.ORM.SelectLatestBlockByEventSigsAddrsWithConfs(ctx, tt.fromBlock, tt.events, tt.addrs, tt.confs)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBlockNumber, blockNumber)
		})
	}
}

func TestSelectLogsCreatedAfter(t *testing.T) {
	ctx := testutils.Context(t)
	th := SetupTH(t, lpOpts)
	event := EmitterABI.Events["Log1"].ID
	address := utils.RandomAddress()

	block1ts := time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC)
	block2ts := time.Date(2020, 1, 1, 12, 12, 12, 0, time.UTC)
	block3ts := time.Date(2030, 1, 1, 12, 12, 12, 0, time.UTC)

	require.NoError(t, th.ORM.InsertLogs(ctx, []logpoller.Log{
		GenLogWithTimestamp(th.ChainID, 1, 1, utils.RandomAddress().String(), event[:], address, block1ts),
		GenLogWithTimestamp(th.ChainID, 1, 2, utils.RandomAddress().String(), event[:], address, block2ts),
		GenLogWithTimestamp(th.ChainID, 2, 2, utils.RandomAddress().String(), event[:], address, block2ts),
		GenLogWithTimestamp(th.ChainID, 1, 3, utils.RandomAddress().String(), event[:], address, block3ts),
	}))
	require.NoError(t, th.ORM.InsertBlock(ctx, utils.RandomHash(), 1, block1ts, 0))
	require.NoError(t, th.ORM.InsertBlock(ctx, utils.RandomHash(), 2, block2ts, 1))
	require.NoError(t, th.ORM.InsertBlock(ctx, utils.RandomHash(), 3, block3ts, 2))

	type expectedLog struct {
		block int64
		log   int64
	}

	tests := []struct {
		name         string
		confs        evmtypes.Confirmations
		after        time.Time
		expectedLogs []expectedLog
	}{
		{
			name:  "picks logs after block 1",
			confs: 0,
			after: block1ts,
			expectedLogs: []expectedLog{
				{block: 2, log: 1},
				{block: 2, log: 2},
				{block: 3, log: 1},
			},
		},
		{
			name:  "skips blocks with not enough confirmations",
			confs: 1,
			after: block1ts,
			expectedLogs: []expectedLog{
				{block: 2, log: 1},
				{block: 2, log: 2},
			},
		},
		{
			name:  "limits number of blocks by block_timestamp",
			confs: 0,
			after: block2ts,
			expectedLogs: []expectedLog{
				{block: 3, log: 1},
			},
		},
		{
			name:         "returns empty dataset for future timestamp",
			confs:        0,
			after:        block3ts,
			expectedLogs: []expectedLog{},
		},
		{
			name:         "returns empty dataset when too many confirmations are required",
			confs:        3,
			after:        block1ts,
			expectedLogs: []expectedLog{},
		},
		{
			name:  "returns only finalized log",
			confs: evmtypes.Finalized,
			after: block1ts,
			expectedLogs: []expectedLog{
				{block: 2, log: 1},
				{block: 2, log: 2},
			},
		},
	}

	filter := func(timestamp time.Time, confs evmtypes.Confirmations, topicIdx int, topicVals []common.Hash) query.KeyFilter {
		var queryConfs primitives.ConfirmationLevel

		switch confs {
		case evmtypes.Finalized:
			queryConfs = primitives.Finalized
		case evmtypes.Unconfirmed:
			queryConfs = primitives.Unconfirmed
		default:
			fmt.Println("default")
			queryConfs = primitives.ConfirmationLevel(confs)
		}

		filters := []query.Expression{
			logpoller.NewAddressFilter(address),
			logpoller.NewEventSigFilter(event),
		}

		if len(topicVals) > 0 {
			exp := make([]query.Expression, len(topicVals))
			for idx, val := range topicVals {
				exp[idx] = logpoller.NewEventByTopicFilter(uint64(topicIdx), []primitives.ValueComparator{
					{Value: val.String(), Operator: primitives.Eq},
				})
			}

			filters = append(filters, query.Expression{
				BoolExpression: query.BoolExpression{
					Expressions:  exp,
					BoolOperator: query.OR,
				},
			})
		}

		filters = append(filters, []query.Expression{
			query.Timestamp(uint64(timestamp.Unix()), primitives.Gt),
			query.Confirmation(queryConfs),
		}...)

		return query.KeyFilter{
			Expressions: filters,
		}
	}

	limiter := query.LimitAndSort{
		SortBy: []query.SortBy{
			query.NewSortBySequence(query.Asc),
		},
	}

	assertion := func(t *testing.T, logs []logpoller.Log, err error, exp []expectedLog) {
		require.NoError(t, err)
		require.Len(t, logs, len(exp))

		for i, log := range logs {
			assert.Equal(t, exp[i].block, log.BlockNumber)
			assert.Equal(t, exp[i].log, log.LogIndex)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := th.ORM.SelectLogsCreatedAfter(ctx, address, event, tt.after, tt.confs)

			assertion(t, logs, err, tt.expectedLogs)

			logs, err = th.ORM.FilteredLogs(ctx, filter(tt.after, tt.confs, 0, nil), limiter, "")

			assertion(t, logs, err, tt.expectedLogs)
		})
	}

	t.Run("SelectIndexedLogsCreatedAfter", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				logs, err := th.ORM.SelectIndexedLogsCreatedAfter(ctx, address, event, 1, []common.Hash{event}, tt.after, tt.confs)

				assertion(t, logs, err, tt.expectedLogs)

				logs, err = th.ORM.FilteredLogs(ctx, filter(tt.after, tt.confs, 1, []common.Hash{event}), limiter, "")

				assertion(t, logs, err, tt.expectedLogs)
			})
		}
	})
}

func TestNestedLogPollerBlocksQuery(t *testing.T) {
	ctx := testutils.Context(t)
	th := SetupTH(t, lpOpts)
	event := EmitterABI.Events["Log1"].ID
	address := utils.RandomAddress()

	require.NoError(t, th.ORM.InsertLogs(ctx, []logpoller.Log{
		GenLog(th.ChainID, 1, 8, utils.RandomAddress().String(), event[:], address),
	}))

	// Empty logs when block are not persisted
	logs, err := th.ORM.SelectIndexedLogs(ctx, address, event, 1, []common.Hash{event}, evmtypes.Unconfirmed)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	// Persist block
	require.NoError(t, th.ORM.InsertBlock(ctx, utils.RandomHash(), 10, time.Now(), 0))

	// Check if query actually works well with provided dataset
	logs, err = th.ORM.SelectIndexedLogs(ctx, address, event, 1, []common.Hash{event}, evmtypes.Unconfirmed)
	require.NoError(t, err)
	require.Len(t, logs, 1)

	// Empty logs when number of confirmations is too deep
	logs, err = th.ORM.SelectIndexedLogs(ctx, address, event, 1, []common.Hash{event}, evmtypes.Confirmations(4))
	require.NoError(t, err)
	require.Len(t, logs, 0)
}

func TestInsertLogsWithBlock(t *testing.T) {
	chainID := testutils.NewRandomEVMChainID()
	event := utils.RandomBytes32()
	address := utils.RandomAddress()
	ctx := testutils.Context(t)

	// We need full db here, because we want to test transaction rollbacks.
	// Using pgtest.NewSqlxDB(t) will run all tests in TXs which is not desired for this type of test
	// (inner tx rollback will rollback outer tx, blocking rest of execution)
	_, db := heavyweight.FullTestDBV2(t, nil)
	o := logpoller.NewORM(chainID, db, logger.Test(t))

	correctLog := GenLog(chainID, 1, 1, utils.RandomAddress().String(), event[:], address)
	invalidLog := GenLog(chainID, -10, -10, utils.RandomAddress().String(), event[:], address)
	correctBlock := logpoller.NewLogPollerBlock(utils.RandomBytes32(), 20, time.Now(), 10)
	invalidBlock := logpoller.NewLogPollerBlock(utils.RandomBytes32(), -10, time.Now(), -10)

	tests := []struct {
		name           string
		logs           []logpoller.Log
		block          logpoller.LogPollerBlock
		shouldRollback bool
	}{
		{
			name:           "properly persist all data",
			logs:           []logpoller.Log{correctLog},
			block:          correctBlock,
			shouldRollback: false,
		},
		{
			name:           "rollbacks transaction when block is invalid",
			logs:           []logpoller.Log{correctLog},
			block:          invalidBlock,
			shouldRollback: true,
		},
		{
			name:           "rollbacks transaction when log is invalid",
			logs:           []logpoller.Log{invalidLog},
			block:          correctBlock,
			shouldRollback: true,
		},
		{
			name:           "rollback when only some logs are invalid",
			logs:           []logpoller.Log{correctLog, invalidLog},
			block:          correctBlock,
			shouldRollback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clean all logs and blocks between test cases
			defer func() { _ = o.DeleteLogsAndBlocksAfter(ctx, 0) }()
			insertError := o.InsertLogsWithBlock(ctx, tt.logs, tt.block)

			logs, logsErr := o.SelectLogs(ctx, 0, math.MaxInt, address, event)
			block, blockErr := o.SelectLatestBlock(ctx)

			if tt.shouldRollback {
				assert.Error(t, insertError)

				assert.NoError(t, logsErr)
				assert.Len(t, logs, 0)

				assert.Error(t, blockErr)
			} else {
				assert.NoError(t, insertError)

				assert.NoError(t, logsErr)
				assert.Len(t, logs, len(tt.logs))

				assert.NoError(t, blockErr)
				assert.Equal(t, block.BlockNumber, tt.block.BlockNumber)
			}
		})
	}
}

func TestInsertLogsInTx(t *testing.T) {
	chainID := testutils.NewRandomEVMChainID()
	event := utils.RandomBytes32()
	address := utils.RandomAddress()
	maxLogsSize := 9000
	ctx := testutils.Context(t)

	// We need full db here, because we want to test transaction rollbacks.
	_, db := heavyweight.FullTestDBV2(t, nil)
	o := logpoller.NewORM(chainID, db, logger.Test(t))

	logs := make([]logpoller.Log, maxLogsSize, maxLogsSize+1)
	for i := 0; i < maxLogsSize; i++ {
		logs[i] = GenLog(chainID, int64(i+1), int64(i+1), utils.RandomAddress().String(), event[:], address)
	}
	invalidLog := GenLog(chainID, -10, -10, utils.RandomAddress().String(), event[:], address)

	tests := []struct {
		name           string
		logs           []logpoller.Log
		shouldRollback bool
	}{
		{
			name:           "all logs persisted",
			logs:           logs,
			shouldRollback: false,
		},
		{
			name:           "rollback when invalid log is passed",
			logs:           append(logs, invalidLog),
			shouldRollback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clean all logs and blocks between test cases
			defer func() { _, _ = db.Exec("truncate evm.logs") }()

			insertErr := o.InsertLogs(ctx, tt.logs)
			logsFromDb, err := o.SelectLogs(ctx, 0, math.MaxInt, address, event)
			assert.NoError(t, err)

			if tt.shouldRollback {
				assert.Error(t, insertErr)
				assert.Len(t, logsFromDb, 0)
			} else {
				assert.NoError(t, insertErr)
				assert.Len(t, logsFromDb, len(tt.logs))
			}
		})
	}
}

func TestSelectLogsDataWordBetween(t *testing.T) {
	ctx := testutils.Context(t)
	address := utils.RandomAddress()
	eventSig := utils.RandomBytes32()
	th := SetupTH(t, lpOpts)

	firstLogData := make([]byte, 0, 64)
	firstLogData = append(firstLogData, logpoller.EvmWord(1).Bytes()...)
	firstLogData = append(firstLogData, logpoller.EvmWord(10).Bytes()...)

	secondLogData := make([]byte, 0, 64)
	secondLogData = append(secondLogData, logpoller.EvmWord(5).Bytes()...)
	secondLogData = append(secondLogData, logpoller.EvmWord(20).Bytes()...)

	err := th.ORM.InsertLogsWithBlock(ctx,
		[]logpoller.Log{
			GenLogWithData(th.ChainID, address, eventSig, 1, 1, firstLogData),
			GenLogWithData(th.ChainID, address, eventSig, 2, 2, secondLogData),
		},
		logpoller.NewLogPollerBlock(utils.RandomBytes32(), 10, time.Now(), 1),
	)
	require.NoError(t, err)
	limiter := query.LimitAndSort{
		SortBy: []query.SortBy{
			query.NewSortByBlock(query.Asc),
			query.NewSortBySequence(query.Asc),
		},
	}

	tests := []struct {
		name         string
		wordValue    uint64
		expectedLogs []int64
	}{
		{
			name:         "returns only first log",
			wordValue:    2,
			expectedLogs: []int64{1},
		},
		{
			name:         "returns only second log",
			wordValue:    11,
			expectedLogs: []int64{2},
		},
		{
			name:         "returns both logs if word value is between",
			wordValue:    5,
			expectedLogs: []int64{1, 2},
		},
		{
			name:         "returns no logs if word value is outside of the range",
			wordValue:    21,
			expectedLogs: []int64{},
		},
	}

	wordFilter := func(word uint64) query.KeyFilter {
		return query.KeyFilter{
			Expressions: []query.Expression{
				logpoller.NewAddressFilter(address),
				logpoller.NewEventSigFilter(eventSig),
				logpoller.NewEventByWordFilter(eventSig, 0, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(word).Hex(), Operator: primitives.Lte},
				}),
				logpoller.NewEventByWordFilter(eventSig, 1, []primitives.ValueComparator{
					{Value: logpoller.EvmWord(word).Hex(), Operator: primitives.Gte},
				}),
				query.Confirmation(primitives.Unconfirmed),
			},
		}
	}

	assertion := func(t *testing.T, logs []logpoller.Log, err error, expected []int64) {
		require.NoError(t, err)
		assert.Len(t, logs, len(expected))

		for index := range logs {
			assert.Equal(t, expected[index], logs[index].BlockNumber)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := th.ORM.SelectLogsDataWordBetween(ctx, address, eventSig, 0, 1, logpoller.EvmWord(tt.wordValue), evmtypes.Unconfirmed)

			assertion(t, logs, err, tt.expectedLogs)

			logs, err = th.ORM.FilteredLogs(ctx, wordFilter(tt.wordValue), limiter, "")

			assertion(t, logs, err, tt.expectedLogs)
		})
	}
}

func Benchmark_LogsDataWordBetween(b *testing.B) {
	chainId := big.NewInt(137)
	_, db := heavyweight.FullTestDBV2(b, nil)
	o := logpoller.NewORM(chainId, db, logger.Test(b))
	ctx := testutils.Context(b)

	numberOfReports := 100_000
	numberOfMessagesPerReport := 256

	commitStoreAddress := utils.RandomAddress()
	commitReportAccepted := utils.RandomBytes32()

	var dbLogs []logpoller.Log
	for i := 0; i < numberOfReports; i++ {
		data := make([]byte, 64)
		// MinSeqNr
		data = append(data, logpoller.EvmWord(uint64(numberOfMessagesPerReport*i+1)).Bytes()...)
		// MaxSeqNr
		data = append(data, logpoller.EvmWord(uint64(numberOfMessagesPerReport*(i+1))).Bytes()...)

		dbLogs = append(dbLogs, logpoller.Log{
			EvmChainId:     ubig.New(chainId),
			LogIndex:       int64(i + 1),
			BlockHash:      utils.RandomBytes32(),
			BlockNumber:    int64(i + 1),
			BlockTimestamp: time.Now(),
			EventSig:       commitReportAccepted,
			Topics:         [][]byte{},
			Address:        commitStoreAddress,
			TxHash:         utils.RandomHash(),
			Data:           data,
			CreatedAt:      time.Now(),
		})
	}
	require.NoError(b, o.InsertBlock(ctx, utils.RandomHash(), int64(numberOfReports*numberOfMessagesPerReport), time.Now(), int64(numberOfReports*numberOfMessagesPerReport)))
	require.NoError(b, o.InsertLogs(ctx, dbLogs))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logs, err := o.SelectLogsDataWordBetween(ctx,
			commitStoreAddress,
			commitReportAccepted,
			2,
			3,
			logpoller.EvmWord(uint64(numberOfReports*numberOfMessagesPerReport/2)), // Pick the middle report
			evmtypes.Unconfirmed,
		)
		assert.NoError(b, err)
		assert.Len(b, logs, 1)
	}
}

func Benchmark_DeleteExpiredLogs(b *testing.B) {
	chainId := big.NewInt(137)
	_, db := heavyweight.FullTestDBV2(b, nil)
	o := logpoller.NewORM(chainId, db, logger.Test(b))
	ctx := testutils.Context(b)

	numberOfReports := 200_000
	commitStoreAddress := utils.RandomAddress()
	commitReportAccepted := utils.RandomBytes32()

	past := time.Now().Add(-1 * time.Hour)

	err := o.InsertFilter(ctx, logpoller.Filter{
		Name:      "test filter",
		EventSigs: []common.Hash{commitReportAccepted},
		Addresses: []common.Address{commitStoreAddress},
		Retention: 1 * time.Millisecond,
	})
	require.NoError(b, err)

	for j := 0; j < 5; j++ {
		var dbLogs []logpoller.Log
		for i := 0; i < numberOfReports; i++ {
			dbLogs = append(dbLogs, logpoller.Log{
				EvmChainId:     ubig.New(chainId),
				LogIndex:       int64(i + 1),
				BlockHash:      utils.RandomBytes32(),
				BlockNumber:    int64(i + 1),
				BlockTimestamp: past,
				EventSig:       commitReportAccepted,
				Topics:         [][]byte{},
				Address:        commitStoreAddress,
				TxHash:         utils.RandomHash(),
				Data:           []byte{},
				CreatedAt:      past,
			})
		}
		require.NoError(b, o.InsertLogs(ctx, dbLogs))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx, err1 := db.Beginx()
		assert.NoError(b, err1)

		_, err1 = o.DeleteExpiredLogs(ctx, 0)
		assert.NoError(b, err1)

		err1 = tx.Rollback()
		assert.NoError(b, err1)
	}
}

func TestSelectOldestBlock(t *testing.T) {
	th := SetupTH(t, lpOpts)
	o1 := th.ORM
	o2 := th.ORM2
	ctx := testutils.Context(t)
	t.Run("Selects oldest within given chain", func(t *testing.T) {
		// insert blocks
		require.NoError(t, o2.InsertBlock(ctx, common.HexToHash("0x1231"), 11, time.Now(), 0))
		require.NoError(t, o2.InsertBlock(ctx, common.HexToHash("0x1232"), 12, time.Now(), 0))
		// insert newer block from different chain
		require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1233"), 13, time.Now(), 0))
		require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1231"), 14, time.Now(), 0))
		block, err := o1.SelectOldestBlock(ctx, 0)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, block.BlockNumber, int64(13))
		require.Equal(t, block.BlockHash, common.HexToHash("0x1233"))
	})
	t.Run("Does not select blocks older than specified limit", func(t *testing.T) {
		require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1232"), 11, time.Now(), 0))
		require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1233"), 13, time.Now(), 0))
		require.NoError(t, o1.InsertBlock(ctx, common.HexToHash("0x1234"), 15, time.Now(), 0))
		block, err := o1.SelectOldestBlock(ctx, 12)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, block.BlockNumber, int64(13))
		require.Equal(t, block.BlockHash, common.HexToHash("0x1233"))
	})
}
