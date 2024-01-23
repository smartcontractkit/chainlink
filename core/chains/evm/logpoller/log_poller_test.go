package logpoller_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func logRuntime(t testing.TB, start time.Time) {
	t.Log("runtime", time.Since(start))
}

func populateDatabase(t testing.TB, o *logpoller.DbORM, chainID *big.Int) (common.Hash, common.Address, common.Address) {
	event1 := EmitterABI.Events["Log1"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")
	startDate := time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC)

	for j := 1; j < 100; j++ {
		var logs []logpoller.Log
		// Max we can insert per batch
		for i := 0; i < 1000; i++ {
			addr := address1
			if (i+(1000*j))%2 == 0 {
				addr = address2
			}
			blockNumber := int64(i + (1000 * j))
			blockTimestamp := startDate.Add(time.Duration(j*1000) * time.Hour)

			logs = append(logs, logpoller.Log{
				EvmChainId:     ubig.New(chainID),
				LogIndex:       1,
				BlockHash:      common.HexToHash(fmt.Sprintf("0x%d", i+(1000*j))),
				BlockNumber:    blockNumber,
				BlockTimestamp: blockTimestamp,
				EventSig:       event1,
				Topics:         [][]byte{event1[:], logpoller.EvmWord(uint64(i + 1000*j)).Bytes()},
				Address:        addr,
				TxHash:         utils.RandomHash(),
				Data:           logpoller.EvmWord(uint64(i + 1000*j)).Bytes(),
				CreatedAt:      blockTimestamp,
			})

		}
		require.NoError(t, o.InsertLogs(logs))
		require.NoError(t, o.InsertBlock(utils.RandomHash(), int64((j+1)*1000-1), startDate.Add(time.Duration(j*1000)*time.Hour), 0))
	}

	return event1, address1, address2
}

func BenchmarkSelectLogsCreatedAfter(b *testing.B) {
	chainId := big.NewInt(137)
	_, db := heavyweight.FullTestDBV2(b, nil)
	o := logpoller.NewORM(chainId, db, logger.Test(b), pgtest.NewQConfig(false))
	event, address, _ := populateDatabase(b, o, chainId)

	// Setting searchDate to pick around 5k logs
	searchDate := time.Date(2020, 1, 1, 12, 12, 12, 0, time.UTC)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logs, err := o.SelectLogsCreatedAfter(address, event, searchDate, 500)
		require.NotZero(b, len(logs))
		require.NoError(b, err)
	}
}

func TestPopulateLoadedDB(t *testing.T) {
	t.Skip("Only for local load testing and query analysis")
	_, db := heavyweight.FullTestDBV2(t, nil)
	chainID := big.NewInt(137)

	o := logpoller.NewORM(big.NewInt(137), db, logger.Test(t), pgtest.NewQConfig(true))
	event1, address1, address2 := populateDatabase(t, o, chainID)

	func() {
		defer logRuntime(t, time.Now())
		_, err1 := o.SelectLogs(750000, 800000, address1, event1)
		require.NoError(t, err1)
	}()
	func() {
		defer logRuntime(t, time.Now())
		_, err1 := o.SelectLatestLogEventSigsAddrsWithConfs(0, []common.Address{address1}, []common.Hash{event1}, 0)
		require.NoError(t, err1)
	}()

	// Confirm all the logs.
	require.NoError(t, o.InsertBlock(common.HexToHash("0x10"), 1000000, time.Now(), 0))
	func() {
		defer logRuntime(t, time.Now())
		lgs, err1 := o.SelectLogsDataWordRange(address1, event1, 0, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err1)
		// 10 since every other log is for address1
		assert.Equal(t, 10, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err1 := o.SelectIndexedLogs(address2, event1, 1, []common.Hash{logpoller.EvmWord(500000), logpoller.EvmWord(500020)}, 0)
		require.NoError(t, err1)
		assert.Equal(t, 2, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err1 := o.SelectIndexedLogsTopicRange(address1, event1, 1, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err1)
		assert.Equal(t, 10, len(lgs))
	}()
}

func TestLogPoller_Integration(t *testing.T) {
	th := SetupTH(t, false, 2, 3, 2, 1000)
	th.Client.Commit() // Block 2. Ensure we have finality number of blocks

	require.NoError(t, th.LogPoller.RegisterFilter(logpoller.Filter{"Integration test", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1}, 0}))
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Addresses, 1)
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Topics, 1)

	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Addresses, 1)
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Topics, 1)

	// Emit some logs in blocks 3->7.
	for i := 0; i < 5; i++ {
		_, err1 := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter1.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}
	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(testutils.Context(t)))

	// The poller starts on a new chain at latest-finality (5 in this case),
	// Replaying from block 4 should guarantee we have block 4 immediately.  (We will also get
	// block 3 once the backup poller runs, since it always starts 100 blocks behind.)
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 4))

	// We should immediately have at least logs 4-7
	logs, err := th.LogPoller.Logs(4, 7, EmitterABI.Events["Log1"].ID, th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)))
	require.NoError(t, err)
	require.Equal(t, 4, len(logs))

	// Once the backup poller runs we should also have the log from block 3
	testutils.AssertEventually(t, func() bool {
		l, err2 := th.LogPoller.Logs(3, 3, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
		require.NoError(t, err2)
		return len(l) == 1
	})

	// Now let's update the Filter and replay to get Log2 logs.
	err = th.LogPoller.RegisterFilter(logpoller.Filter{
		"Emitter - log2", []common.Hash{EmitterABI.Events["Log2"].ID},
		[]common.Address{th.EmitterAddress1}, 0,
	})
	require.NoError(t, err)
	// Replay an invalid block should error
	assert.Error(t, th.LogPoller.Replay(testutils.Context(t), 0))
	assert.Error(t, th.LogPoller.Replay(testutils.Context(t), 20))

	// Still shouldn't have any Log2 logs yet
	logs, err = th.LogPoller.Logs(2, 7, EmitterABI.Events["Log2"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	require.Len(t, logs, 0)

	// Replay only from block 4, so we should see logs in block 4,5,6,7 (4 logs)
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 4))

	// We should immediately see 4 logs2 logs.
	logs, err = th.LogPoller.Logs(2, 7, EmitterABI.Events["Log2"].ID, th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)))
	require.NoError(t, err)
	assert.Equal(t, 4, len(logs))

	assert.NoError(t, th.LogPoller.Close())

	// Cancelling a replay should return an error synchronously.
	ctx, cancel := context.WithCancel(testutils.Context(t))
	cancel()
	assert.ErrorIs(t, th.LogPoller.Replay(ctx, 4), logpoller.ErrReplayRequestAborted)
}

// Simulate a badly behaving rpc server, where unfinalized blocks can return different logs
// for the same block hash.  We should be able to handle this without missing any logs, as
// long as the logs returned for finalized blocks are consistent.
func Test_BackupLogPoller(t *testing.T) {
	tests := []struct {
		name          string
		finalityDepth int64
		finalityTag   bool
	}{
		{
			name:          "fixed finality depth without finality tag",
			finalityDepth: 2,
			finalityTag:   false,
		},
		{
			name:          "chain finality in use",
			finalityDepth: 0,
			finalityTag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, tt.finalityTag, tt.finalityDepth, 3, 2, 1000)
			// later, we will need at least 32 blocks filled with logs for cache invalidation
			for i := int64(0); i < 32; i++ {
				// to invalidate geth's internal read-cache, a matching log must be found in the bloom Filter
				// for each of the 32 blocks
				tx, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(i + 7)})
				require.NoError(t, err)
				require.NotNil(t, tx)
				th.Client.Commit()
			}

			ctx := testutils.Context(t)

			filter1 := logpoller.Filter{"filter1", []common.Hash{
				EmitterABI.Events["Log1"].ID,
				EmitterABI.Events["Log2"].ID},
				[]common.Address{th.EmitterAddress1},
				0}
			err := th.LogPoller.RegisterFilter(filter1)
			require.NoError(t, err)

			filters, err := th.ORM.LoadFilters(pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			require.Equal(t, 1, len(filters))
			require.Equal(t, filter1, filters["filter1"])

			err = th.LogPoller.RegisterFilter(
				logpoller.Filter{"filter2",
					[]common.Hash{EmitterABI.Events["Log1"].ID},
					[]common.Address{th.EmitterAddress2}, 0})
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, th.LogPoller.UnregisterFilter("filter1"))
			}()
			defer func() {
				assert.NoError(t, th.LogPoller.UnregisterFilter("filter2"))
			}()

			// generate some tx's with logs
			tx1, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
			require.NoError(t, err)
			require.NotNil(t, tx1)

			tx2, err := th.Emitter1.EmitLog2(th.Owner, []*big.Int{big.NewInt(2)})
			require.NoError(t, err)
			require.NotNil(t, tx2)

			tx3, err := th.Emitter2.EmitLog1(th.Owner, []*big.Int{big.NewInt(3)})
			require.NoError(t, err)
			require.NotNil(t, tx3)

			th.Client.Commit() // commit block 34 with 3 tx's included

			h := th.Client.Blockchain().CurrentHeader() // get latest header
			require.Equal(t, uint64(34), h.Number.Uint64())

			// save these 3 receipts for later
			receipts := rawdb.ReadReceipts(th.EthDB, h.Hash(), h.Number.Uint64(), uint64(time.Now().Unix()), params.AllEthashProtocolChanges)
			require.NotZero(t, receipts.Len())

			// Simulate a situation where the rpc server has a block, but no logs available for it yet
			//  this can't happen with geth itself, but can with other clients.
			rawdb.WriteReceipts(th.EthDB, h.Hash(), h.Number.Uint64(), types.Receipts{}) // wipes out all logs for block 34

			body := rawdb.ReadBody(th.EthDB, h.Hash(), h.Number.Uint64())
			require.Equal(t, 3, len(body.Transactions))
			txs := body.Transactions                 // save transactions for later
			body.Transactions = types.Transactions{} // number of tx's must match # of logs for GetLogs() to succeed
			rawdb.WriteBody(th.EthDB, h.Hash(), h.Number.Uint64(), body)

			currentBlockNumber := th.PollAndSaveLogs(ctx, 1)
			assert.Equal(t, int64(35), currentBlockNumber)

			// simulate logs becoming available
			rawdb.WriteReceipts(th.EthDB, h.Hash(), h.Number.Uint64(), receipts)
			require.True(t, rawdb.HasReceipts(th.EthDB, h.Hash(), h.Number.Uint64()))
			body.Transactions = txs
			rawdb.WriteBody(th.EthDB, h.Hash(), h.Number.Uint64(), body)

			// flush out cached block 34 by reading logs from first 32 blocks
			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(2)),
				ToBlock:   big.NewInt(int64(33)),
				Addresses: []common.Address{th.EmitterAddress1},
				Topics:    [][]common.Hash{{EmitterABI.Events["Log1"].ID}},
			}
			fLogs, err := th.Client.FilterLogs(ctx, query)
			require.NoError(t, err)
			require.Equal(t, 32, len(fLogs))

			// logs shouldn't show up yet
			logs, err := th.LogPoller.Logs(34, 34, EmitterABI.Events["Log1"].ID, th.EmitterAddress1,
				pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			assert.Equal(t, 0, len(logs))

			th.Client.Commit()
			th.Client.Commit()
			markBlockAsFinalized(t, th, 34)

			// Run ordinary poller + backup poller at least once
			currentBlock, _ := th.LogPoller.LatestBlock(pg.WithParentCtx(testutils.Context(t)))
			th.LogPoller.PollAndSaveLogs(ctx, currentBlock.BlockNumber+1)
			th.LogPoller.BackupPollAndSaveLogs(ctx, 100)
			currentBlock, _ = th.LogPoller.LatestBlock(pg.WithParentCtx(testutils.Context(t)))

			require.Equal(t, int64(37), currentBlock.BlockNumber+1)

			// logs still shouldn't show up, because we don't want to backfill the last finalized log
			//  to help with reorg detection
			logs, err = th.LogPoller.Logs(34, 34, EmitterABI.Events["Log1"].ID, th.EmitterAddress1,
				pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			assert.Equal(t, 0, len(logs))
			th.Client.Commit()
			markBlockAsFinalized(t, th, 35)

			// Run ordinary poller + backup poller at least once more
			th.LogPoller.PollAndSaveLogs(ctx, currentBlockNumber+1)
			th.LogPoller.BackupPollAndSaveLogs(ctx, 100)
			currentBlock, _ = th.LogPoller.LatestBlock(pg.WithParentCtx(testutils.Context(t)))

			require.Equal(t, int64(38), currentBlock.BlockNumber+1)

			// all 3 logs in block 34 should show up now, thanks to backup logger
			logs, err = th.LogPoller.Logs(30, 37, EmitterABI.Events["Log1"].ID, th.EmitterAddress1,
				pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			assert.Equal(t, 5, len(logs))
			logs, err = th.LogPoller.Logs(34, 34, EmitterABI.Events["Log2"].ID, th.EmitterAddress1,
				pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			assert.Equal(t, 1, len(logs))
			logs, err = th.LogPoller.Logs(32, 36, EmitterABI.Events["Log1"].ID, th.EmitterAddress2,
				pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
			assert.Equal(t, 1, len(logs))
		})
	}
}

func TestLogPoller_BackupPollAndSaveLogsWithPollerNotWorking(t *testing.T) {
	emittedLogs := 30
	// Intentionally use very low backupLogPollerDelay to verify if finality is used properly
	backupLogPollerDelay := int64(0)
	ctx := testutils.Context(t)
	th := SetupTH(t, true, 0, 3, 2, 1000)

	header, err := th.Client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)

	// Emit some logs in blocks
	for i := 0; i < emittedLogs; i++ {
		_, err2 := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err2)
		th.Client.Commit()
	}

	// First PollAndSave, no filters are registered
	// 0 (finalized) -> 1 -> 2 -> ...
	currentBlock := th.PollAndSaveLogs(ctx, 1)
	// currentBlock should be blockChain start + number of emitted logs + 1
	assert.Equal(t, int64(emittedLogs)+header.Number.Int64()+1, currentBlock)

	// LogPoller not working, but chain in the meantime has progressed
	// 0 -> 1 -> 2 -> ... -> currentBlock - 10 (finalized) -> .. -> currentBlock
	markBlockAsFinalized(t, th, currentBlock-10)

	err = th.LogPoller.RegisterFilter(logpoller.Filter{
		Name:      "Test Emitter",
		EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
		Addresses: []common.Address{th.EmitterAddress1},
	})
	require.NoError(t, err)

	// LogPoller should backfill starting from the last finalized block stored in db (genesis block)
	// till the latest finalized block reported by chain.
	th.LogPoller.BackupPollAndSaveLogs(ctx, backupLogPollerDelay)
	require.NoError(t, err)

	logs, err := th.LogPoller.Logs(
		0,
		currentBlock,
		EmitterABI.Events["Log1"].ID,
		th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)),
	)
	require.NoError(t, err)
	require.Len(t, logs, emittedLogs-10)

	// Progressing even more, move blockchain forward by 1 block and mark it as finalized
	th.Client.Commit()
	markBlockAsFinalized(t, th, currentBlock)
	th.LogPoller.BackupPollAndSaveLogs(ctx, backupLogPollerDelay)

	// All emitted logs should be backfilled
	logs, err = th.LogPoller.Logs(
		0,
		currentBlock+1,
		EmitterABI.Events["Log1"].ID,
		th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)),
	)
	require.NoError(t, err)
	require.Len(t, logs, emittedLogs)
}

func TestLogPoller_BackupPollAndSaveLogsWithDeepBlockDelay(t *testing.T) {
	emittedLogs := 30
	ctx := testutils.Context(t)
	th := SetupTH(t, true, 0, 3, 2, 1000)

	// Emit some logs in blocks
	for i := 0; i < emittedLogs; i++ {
		_, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		th.Client.Commit()
	}
	// Emit one more empty block
	th.Client.Commit()

	header, err := th.Client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	// Mark everything as finalized
	markBlockAsFinalized(t, th, header.Number.Int64())

	// First PollAndSave, no filters are registered, but finalization is the same as the latest block
	// 1 -> 2 -> ...
	th.PollAndSaveLogs(ctx, 1)

	// Check that latest block has the same properties as the head
	latestBlock, err := th.LogPoller.LatestBlock()
	require.NoError(t, err)
	assert.Equal(t, latestBlock.BlockNumber, header.Number.Int64())
	assert.Equal(t, latestBlock.FinalizedBlockNumber, header.Number.Int64())
	assert.Equal(t, latestBlock.BlockHash, header.Hash())

	// Register filter
	err = th.LogPoller.RegisterFilter(logpoller.Filter{
		Name:      "Test Emitter",
		EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
		Addresses: []common.Address{th.EmitterAddress1},
	})
	require.NoError(t, err)

	// Should fallback to the backupPollerBlockDelay when finalization was very high in a previous PollAndSave
	th.LogPoller.BackupPollAndSaveLogs(ctx, int64(emittedLogs))
	require.NoError(t, err)

	// All emitted logs should be backfilled
	logs, err := th.LogPoller.Logs(
		0,
		header.Number.Int64()+1,
		EmitterABI.Events["Log1"].ID,
		th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)),
	)
	require.NoError(t, err)
	require.Len(t, logs, emittedLogs)
}

func TestLogPoller_BackupPollAndSaveLogsSkippingLogsThatAreTooOld(t *testing.T) {
	logsBatch := 10
	// Intentionally use very low backupLogPollerDelay to verify if finality is used properly
	ctx := testutils.Context(t)
	th := SetupTH(t, true, 0, 3, 2, 1000)

	//header, err := th.Client.HeaderByNumber(ctx, nil)
	//require.NoError(t, err)

	// Emit some logs in blocks
	for i := 1; i <= logsBatch; i++ {
		_, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		th.Client.Commit()
	}

	// First PollAndSave, no filters are registered, but finalization is the same as the latest block
	// 1 -> 2 -> ... -> firstBatchBlock
	firstBatchBlock := th.PollAndSaveLogs(ctx, 1)
	// Mark current tip of the chain as finalized (after emitting 10 logs)
	markBlockAsFinalized(t, th, firstBatchBlock)

	// Emit 2nd batch of block
	for i := 1; i <= logsBatch; i++ {
		_, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(100 + i))})
		require.NoError(t, err)
		th.Client.Commit()
	}

	// 1 -> 2 -> ... -> firstBatchBlock (finalized) -> .. -> firstBatchBlock + emitted logs
	secondBatchBlock := th.PollAndSaveLogs(ctx, firstBatchBlock)
	// Mark current tip of the block as finalized (after emitting 20 logs)
	markBlockAsFinalized(t, th, secondBatchBlock)

	// Register filter
	err := th.LogPoller.RegisterFilter(logpoller.Filter{
		Name:      "Test Emitter",
		EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
		Addresses: []common.Address{th.EmitterAddress1},
	})
	require.NoError(t, err)

	// Should pick logs starting from one block behind the latest finalized block
	th.LogPoller.BackupPollAndSaveLogs(ctx, 0)
	require.NoError(t, err)

	// Only the 2nd batch + 1 log from a previous batch should be backfilled, because we perform backfill starting
	// from one block behind the latest finalized block
	logs, err := th.LogPoller.Logs(
		0,
		secondBatchBlock,
		EmitterABI.Events["Log1"].ID,
		th.EmitterAddress1,
		pg.WithParentCtx(testutils.Context(t)),
	)
	require.NoError(t, err)
	require.Len(t, logs, logsBatch+1)
	require.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000009`), logs[0].Data)
}

func TestLogPoller_BlockTimestamps(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	th := SetupTH(t, false, 2, 3, 2, 1000)

	addresses := []common.Address{th.EmitterAddress1, th.EmitterAddress2}
	topics := []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"convertLogs", topics, addresses, 0})
	require.NoError(t, err)

	blk, err := th.Client.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), blk.Number())
	start := blk.Time()

	// There is automatically a 10s delay between each block.  To make sure it's including the correct block timestamps,
	// we introduce irregularities by inserting two additional block delays. We can't control the block times for
	// blocks produced by the log emitter, but we can adjust the time on empty blocks in between.  Simulated time
	// sequence:  [ #1 ] ..(10s + delay1).. [ #2 ] ..10s.. [ #3 (LOG1) ] ..(10s + delay2).. [ #4 ] ..10s.. [ #5 (LOG2) ]
	const delay1 = 589
	const delay2 = 643
	time1 := start + 20 + delay1
	time2 := time1 + 20 + delay2

	require.NoError(t, th.Client.AdjustTime(delay1*time.Second))
	hash := th.Client.Commit()

	blk, err = th.Client.BlockByHash(ctx, hash)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(2), blk.Number())
	assert.Equal(t, time1-10, blk.Time())

	_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)
	hash = th.Client.Commit()

	blk, err = th.Client.BlockByHash(ctx, hash)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(3), blk.Number())
	assert.Equal(t, time1, blk.Time())

	require.NoError(t, th.Client.AdjustTime(delay2*time.Second))
	th.Client.Commit()
	_, err = th.Emitter2.EmitLog2(th.Owner, []*big.Int{big.NewInt(2)})
	require.NoError(t, err)
	hash = th.Client.Commit()

	blk, err = th.Client.BlockByHash(ctx, hash)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(5), blk.Number())
	assert.Equal(t, time2, blk.Time())

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2),
		ToBlock:   big.NewInt(5),
		Topics:    [][]common.Hash{topics},
		Addresses: []common.Address{th.EmitterAddress1, th.EmitterAddress2}}

	gethLogs, err := th.Client.FilterLogs(ctx, query)
	require.NoError(t, err)
	require.Len(t, gethLogs, 2)

	lb, _ := th.LogPoller.LatestBlock(pg.WithParentCtx(testutils.Context(t)))
	th.PollAndSaveLogs(ctx, lb.BlockNumber+1)
	lg1, err := th.LogPoller.Logs(0, 20, EmitterABI.Events["Log1"].ID, th.EmitterAddress1,
		pg.WithParentCtx(ctx))
	require.NoError(t, err)
	lg2, err := th.LogPoller.Logs(0, 20, EmitterABI.Events["Log2"].ID, th.EmitterAddress2,
		pg.WithParentCtx(ctx))
	require.NoError(t, err)

	// Logs should have correct timestamps
	b, _ := th.Client.BlockByHash(ctx, lg1[0].BlockHash)
	t.Log(len(lg1), lg1[0].BlockTimestamp)
	assert.Equal(t, int64(b.Time()), lg1[0].BlockTimestamp.UTC().Unix(), time1)
	b2, _ := th.Client.BlockByHash(ctx, lg2[0].BlockHash)
	assert.Equal(t, int64(b2.Time()), lg2[0].BlockTimestamp.UTC().Unix(), time2)
}

func TestLogPoller_SynchronizedWithGeth(t *testing.T) {
	t.Parallel()
	// The log poller's blocks table should remain synchronized
	// with the canonical chain of geth's despite arbitrary mixes of mining and reorgs.
	testParams := gopter.DefaultTestParameters()
	testParams.MinSuccessfulTests = 100
	p := gopter.NewProperties(testParams)
	numChainInserts := 3
	finalityDepth := 5
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)

	owner := testutils.MustNewSimTransactor(t)
	owner.GasPrice = big.NewInt(10e9)
	p.Property("synchronized with geth", prop.ForAll(func(mineOrReorg []uint64) bool {
		// After the set of reorgs, we should have the same canonical blocks that geth does.
		t.Log("Starting test", mineOrReorg)
		chainID := testutils.NewRandomEVMChainID()
		// Set up a test chain with a log emitting contract deployed.
		orm := logpoller.NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
		// Note this property test is run concurrently and the sim is not threadsafe.
		ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
			owner.From: {
				Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
			},
		}, 10e6)
		_, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
		require.NoError(t, err)
		lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, false, int64(finalityDepth), 3, 2, 1000)
		for i := 0; i < finalityDepth; i++ { // Have enough blocks that we could reorg the full finalityDepth-1.
			ec.Commit()
		}
		currentBlockNumber := int64(1)
		lp.PollAndSaveLogs(testutils.Context(t), currentBlockNumber)
		currentBlock, err := lp.LatestBlock(pg.WithParentCtx(testutils.Context(t)))
		require.NoError(t, err)
		matchesGeth := func() bool {
			// Check every block is identical
			latest, err1 := ec.BlockByNumber(testutils.Context(t), nil)
			require.NoError(t, err1)
			for i := 1; i < int(latest.NumberU64()); i++ {
				ourBlock, err1 := lp.BlockByNumber(int64(i))
				require.NoError(t, err1)
				gethBlock, err1 := ec.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
				require.NoError(t, err1)
				if ourBlock.BlockHash != gethBlock.Hash() {
					t.Logf("Initial poll our block differs at height %d got %x want %x\n", i, ourBlock.BlockHash, gethBlock.Hash())
					return false
				}
			}
			return true
		}
		if !matchesGeth() {
			return false
		}
		// Randomly pick to mine or reorg
		for i := 0; i < numChainInserts; i++ {
			if rand.Bool() {
				// Mine blocks
				for j := 0; j < int(mineOrReorg[i]); j++ {
					ec.Commit()
					latest, err1 := ec.BlockByNumber(testutils.Context(t), nil)
					require.NoError(t, err1)
					t.Log("mined block", latest.Hash())
				}
			} else {
				// Reorg blocks
				latest, err1 := ec.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err1)
				reorgedBlock := big.NewInt(0).Sub(latest.Number(), big.NewInt(int64(mineOrReorg[i])))
				reorg, err1 := ec.BlockByNumber(testutils.Context(t), reorgedBlock)
				require.NoError(t, err1)
				require.NoError(t, ec.Fork(testutils.Context(t), reorg.Hash()))
				t.Logf("Reorging from (%v, %x) back to (%v, %x)\n", latest.NumberU64(), latest.Hash(), reorgedBlock.Uint64(), reorg.Hash())
				// Actually need to change the block here to trigger the reorg.
				_, err1 = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(1)})
				require.NoError(t, err1)
				for j := 0; j < int(mineOrReorg[i]+1); j++ { // Need +1 to make it actually longer height so we detect it.
					ec.Commit()
				}
				latest, err1 = ec.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err1)
				t.Logf("New latest (%v, %x), latest parent %x)\n", latest.NumberU64(), latest.Hash(), latest.ParentHash())
			}
			lp.PollAndSaveLogs(testutils.Context(t), currentBlock.BlockNumber)
			currentBlock, err = lp.LatestBlock(pg.WithParentCtx(testutils.Context(t)))
			require.NoError(t, err)
		}
		return matchesGeth()
	}, gen.SliceOfN(numChainInserts, gen.UInt64Range(1, uint64(finalityDepth-1))))) // Max reorg depth is finality depth - 1
	p.TestingRun(t)
}

func TestLogPoller_PollAndSaveLogs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		finalityDepth int64
		finalityTag   bool
	}{
		{
			name:          "fixed finality depth without finality tag",
			finalityDepth: 3,
			finalityTag:   false,
		},
		{
			name:          "chain finality in use",
			finalityDepth: 0,
			finalityTag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, tt.finalityTag, tt.finalityDepth, 3, 2, 1000)

			// Set up a log poller listening for log emitter logs.
			err := th.LogPoller.RegisterFilter(logpoller.Filter{
				"Test Emitter 1 & 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID},
				[]common.Address{th.EmitterAddress1, th.EmitterAddress2}, 0,
			})
			require.NoError(t, err)

			b, err := th.Client.BlockByNumber(testutils.Context(t), nil)
			require.NoError(t, err)
			require.Equal(t, uint64(1), b.NumberU64())
			require.Equal(t, uint64(10), b.Time())

			// Test scenario: single block in chain, no logs.
			// Chain genesis <- 1
			// DB: empty
			newStart := th.PollAndSaveLogs(testutils.Context(t), 1)
			assert.Equal(t, int64(2), newStart)

			// We expect to have saved block 1.
			lpb, err := th.ORM.SelectBlockByNumber(1)
			require.NoError(t, err)
			assert.Equal(t, lpb.BlockHash, b.Hash())
			assert.Equal(t, lpb.BlockNumber, int64(b.NumberU64()))
			assert.Equal(t, int64(1), int64(b.NumberU64()))
			assert.Equal(t, uint64(10), b.Time())

			// No logs.
			lgs, err := th.ORM.SelectLogsByBlockRange(1, 1)
			require.NoError(t, err)
			assert.Equal(t, 0, len(lgs))
			th.assertHaveCanonical(t, 1, 1)

			// Polling again should be a noop, since we are at the latest.
			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(2), newStart)
			latest, err := th.ORM.SelectLatestBlock()
			require.NoError(t, err)
			assert.Equal(t, int64(1), latest.BlockNumber)
			th.assertHaveCanonical(t, 1, 1)

			// Test scenario: one log 2 block chain.
			// Chain gen <- 1 <- 2 (L1)
			// DB: 1
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
			require.NoError(t, err)
			th.Client.Commit()

			// Polling should get us the L1 log.
			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(3), newStart)
			latest, err = th.ORM.SelectLatestBlock()
			require.NoError(t, err)
			assert.Equal(t, int64(2), latest.BlockNumber)
			lgs, err = th.ORM.SelectLogsByBlockRange(1, 3)
			require.NoError(t, err)
			require.Equal(t, 1, len(lgs))
			assert.Equal(t, th.EmitterAddress1, lgs[0].Address)
			assert.Equal(t, latest.BlockHash, lgs[0].BlockHash)
			assert.Equal(t, latest.BlockTimestamp, lgs[0].BlockTimestamp)
			assert.Equal(t, hexutil.Encode(lgs[0].Topics[0]), EmitterABI.Events["Log1"].ID.String())
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`),
				lgs[0].Data)

			// Test scenario: single block reorg with log.
			// Chain gen <- 1 <- 2 (L1_1)
			//                \ 2'(L1_2) <- 3
			// DB: 1, 2
			// - Detect a reorg,
			// - Update the block 2's hash
			// - Save L1'
			// - L1_1 deleted
			reorgedOutBlock, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(2))
			require.NoError(t, err)
			lca, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(1))
			require.NoError(t, err)
			require.NoError(t, th.Client.Fork(testutils.Context(t), lca.Hash()))
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(2)})
			require.NoError(t, err)
			// Create 2'
			th.Client.Commit()
			// Create 3 (we need a new block for us to do any polling and detect the reorg).
			th.Client.Commit()

			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(4), newStart)
			latest, err = th.ORM.SelectLatestBlock()
			require.NoError(t, err)
			assert.Equal(t, int64(3), latest.BlockNumber)
			lgs, err = th.ORM.SelectLogsByBlockRange(1, 3)
			require.NoError(t, err)
			require.Equal(t, 1, len(lgs))
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[0].Data)
			th.assertHaveCanonical(t, 1, 3)

			// Test scenario: reorg back to previous tip.
			// Chain gen <- 1 <- 2 (L1_1) <- 3' (L1_3) <- 4
			//                \ 2'(L1_2) <- 3
			require.NoError(t, th.Client.Fork(testutils.Context(t), reorgedOutBlock.Hash()))
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(3)})
			require.NoError(t, err)
			// Create 3'
			th.Client.Commit()
			// Create 4
			th.Client.Commit()
			// Mark block 1 as finalized
			markBlockAsFinalized(t, th, 1)
			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(5), newStart)
			latest, err = th.ORM.SelectLatestBlock()
			require.NoError(t, err)
			assert.Equal(t, int64(4), latest.BlockNumber)
			lgs, err = th.ORM.SelectLogsByBlockRange(1, 3)
			require.NoError(t, err)
			// We expect ONLY L1_1 and L1_3 since L1_2 is reorg'd out.
			assert.Equal(t, 2, len(lgs))
			assert.Equal(t, int64(2), lgs[0].BlockNumber)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`), lgs[0].Data)
			assert.Equal(t, int64(3), lgs[1].BlockNumber)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000003`), lgs[1].Data)
			th.assertHaveCanonical(t, 1, 1)
			th.assertHaveCanonical(t, 3, 4)
			th.assertDontHave(t, 2, 2) // 2 gets backfilled

			// Test scenario: multiple logs per block for many blocks (also after reorg).
			// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6)
			//                \ 2'(L1_2) <- 3
			// DB: 1, 2', 3'
			// - Should save 4, 5, 6 blocks
			// - Should obtain logs L1_3, L2_5, L1_6
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(4)})
			require.NoError(t, err)
			_, err = th.Emitter2.EmitLog1(th.Owner, []*big.Int{big.NewInt(5)})
			require.NoError(t, err)
			// Create 4
			th.Client.Commit()
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(6)})
			require.NoError(t, err)
			// Create 5
			th.Client.Commit()
			// Mark block 2 as finalized
			markBlockAsFinalized(t, th, 3)

			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(7), newStart)
			lgs, err = th.ORM.SelectLogsByBlockRange(4, 6)
			require.NoError(t, err)
			require.Equal(t, 3, len(lgs))
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000004`), lgs[0].Data)
			assert.Equal(t, th.EmitterAddress1, lgs[0].Address)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000005`), lgs[1].Data)
			assert.Equal(t, th.EmitterAddress2, lgs[1].Address)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000006`), lgs[2].Data)
			assert.Equal(t, th.EmitterAddress1, lgs[2].Address)
			th.assertHaveCanonical(t, 1, 1)
			th.assertDontHave(t, 2, 2) // 2 gets backfilled
			th.assertHaveCanonical(t, 3, 6)

			// Test scenario: node down for exactly finality + 2 blocks
			// Note we only backfill up to finalized - 1 blocks, because we need to save the
			// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9) <- 10 (L1_10)
			//                \ 2'(L1_2) <- 3
			// DB: 1, 2, 3, 4, 5, 6
			// - We expect block 7 to backfilled (treated as finalized)
			// - Then block 8-10 to be handled block by block (treated as unfinalized).
			for i := 7; i < 11; i++ {
				_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
				require.NoError(t, err)
				th.Client.Commit()
			}
			// Mark block 7 as finalized
			markBlockAsFinalized(t, th, 7)

			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(11), newStart)
			lgs, err = th.ORM.SelectLogsByBlockRange(7, 9)
			require.NoError(t, err)
			require.Equal(t, 3, len(lgs))
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000007`), lgs[0].Data)
			assert.Equal(t, int64(7), lgs[0].BlockNumber)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000008`), lgs[1].Data)
			assert.Equal(t, int64(8), lgs[1].BlockNumber)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000009`), lgs[2].Data)
			assert.Equal(t, int64(9), lgs[2].BlockNumber)
			th.assertDontHave(t, 7, 7) // Do not expect to save backfilled blocks.
			th.assertHaveCanonical(t, 8, 10)

			// Test scenario large backfill (multiple batches)
			// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9) <- 10..16
			//                \ 2'(L1_2) <- 3
			// DB: 1, 2, 3, 4, 5, 6, (backfilled 7), 8, 9, 10
			// - 11, 12, 13 backfilled in batch 1
			// - 14 backfilled in batch 2
			// - 15, 16, 17 to be treated as unfinalized
			for i := 11; i < 18; i++ {
				_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
				require.NoError(t, err)
				th.Client.Commit()
			}
			// Mark block 14 as finalized
			markBlockAsFinalized(t, th, 14)

			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(18), newStart)
			lgs, err = th.ORM.SelectLogsByBlockRange(11, 17)
			require.NoError(t, err)
			assert.Equal(t, 7, len(lgs))
			th.assertHaveCanonical(t, 14, 16) // Should have last finalized block plus unfinalized blocks
			th.assertDontHave(t, 11, 13)      // Should not have older finalized blocks

			// Verify that a custom block timestamp will get written to db correctly also
			b, err = th.Client.BlockByNumber(testutils.Context(t), nil)
			require.NoError(t, err)
			require.Equal(t, uint64(17), b.NumberU64())
			require.Equal(t, uint64(170), b.Time())
			require.NoError(t, th.Client.AdjustTime(1*time.Hour))
			th.Client.Commit()

			b, err = th.Client.BlockByNumber(testutils.Context(t), nil)
			require.NoError(t, err)
			require.Equal(t, uint64(180+time.Hour.Seconds()), b.Time())
		})
	}
}

func TestLogPoller_PollAndSaveLogsDeepReorg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		finalityDepth int64
		finalityTag   bool
	}{
		{
			name:          "fixed finality depth without finality tag",
			finalityDepth: 3,
			finalityTag:   false,
		},
		{
			name:          "chain finality in use",
			finalityDepth: 0,
			finalityTag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, tt.finalityTag, tt.finalityDepth, 3, 2, 1000)

			// Set up a log poller listening for log emitter logs.
			err := th.LogPoller.RegisterFilter(logpoller.Filter{
				Name:      "Test Emitter",
				EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
				Addresses: []common.Address{th.EmitterAddress1},
			})
			require.NoError(t, err)

			// Test scenario: one log 2 block chain.
			// Chain gen <- 1 <- 2 (L1_1)
			// DB: 1
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
			require.NoError(t, err)
			th.Client.Commit()
			markBlockAsFinalized(t, th, 1)

			// Polling should get us the L1 log.
			newStart := th.PollAndSaveLogs(testutils.Context(t), 1)
			assert.Equal(t, int64(3), newStart)
			// Check that L1_1 has a proper data payload
			lgs, err := th.ORM.SelectLogsByBlockRange(2, 2)
			require.NoError(t, err)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`), lgs[0].Data)

			// Single block reorg and log poller not working for a while, mine blocks and progress with finalization
			// Chain gen <- 1 <- 2 (L1_1)
			//                \ 2'(L1_2) <- 3 <- 4 <- 5 <- 6 (finalized on chain) <- 7 <- 8 <- 9 <- 10
			lca, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(1))
			require.NoError(t, err)
			require.NoError(t, th.Client.Fork(testutils.Context(t), lca.Hash()))
			// Create 2'
			_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(2)})
			require.NoError(t, err)
			th.Client.Commit()
			// Create 3-10
			for i := 3; i < 10; i++ {
				_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
				require.NoError(t, err)
				th.Client.Commit()
			}
			markBlockAsFinalized(t, th, 6)

			newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
			assert.Equal(t, int64(10), newStart)

			// Expect L1_2 to be properly updated
			lgs, err = th.ORM.SelectLogsByBlockRange(2, 2)
			require.NoError(t, err)
			assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[0].Data)
			th.assertHaveCanonical(t, 1, 1)
			th.assertDontHave(t, 2, 3) // These blocks are backfilled
			th.assertHaveCanonical(t, 5, 10)
		})
	}
}

func TestLogPoller_LoadFilters(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, false, 2, 3, 2, 1000)

	filter1 := logpoller.Filter{"first Filter", []common.Hash{
		EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}, 0}
	filter2 := logpoller.Filter{"second Filter", []common.Hash{
		EmitterABI.Events["Log2"].ID, EmitterABI.Events["Log3"].ID}, []common.Address{th.EmitterAddress2}, 0}
	filter3 := logpoller.Filter{"third Filter", []common.Hash{
		EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}, 0}

	assert.True(t, filter1.Contains(nil))
	assert.False(t, filter1.Contains(&filter2))
	assert.False(t, filter2.Contains(&filter1))
	assert.True(t, filter1.Contains(&filter3))

	err := th.LogPoller.RegisterFilter(filter1)
	require.NoError(t, err)
	err = th.LogPoller.RegisterFilter(filter2)
	require.NoError(t, err)
	err = th.LogPoller.RegisterFilter(filter3)
	require.NoError(t, err)

	filters, err := th.ORM.LoadFilters()
	require.NoError(t, err)
	require.NotNil(t, filters)
	require.Len(t, filters, 3)

	filter, ok := filters["first Filter"]
	require.True(t, ok)
	assert.True(t, filter.Contains(&filter1))
	assert.True(t, filter1.Contains(&filter))

	filter, ok = filters["second Filter"]
	require.True(t, ok)
	assert.True(t, filter.Contains(&filter2))
	assert.True(t, filter2.Contains(&filter))

	filter, ok = filters["third Filter"]
	require.True(t, ok)
	assert.True(t, filter.Contains(&filter3))
	assert.True(t, filter3.Contains(&filter))

	t.Run("HasFilter", func(t *testing.T) {
		assert.True(t, th.LogPoller.HasFilter("first Filter"))
		assert.True(t, th.LogPoller.HasFilter("second Filter"))
		assert.True(t, th.LogPoller.HasFilter("third Filter"))
		assert.False(t, th.LogPoller.HasFilter("fourth Filter"))
	})
}

func TestLogPoller_GetBlocks_Range(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, false, 2, 3, 2, 1000)

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"GetBlocks Test", []common.Hash{
		EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}, 0},
	)
	require.NoError(t, err)

	// LP retrieves 0 blocks
	blockNums := []uint64{}
	blocks, err := th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 0, len(blocks))

	// LP retrieves block 1
	blockNums = []uint64{1}
	blocks, err = th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 1, len(blocks))
	assert.Equal(t, 1, int(blocks[0].BlockNumber))

	// LP fails to retrieve block 2 because it's neither in DB nor returned by RPC
	blockNums = []uint64{2}
	_, err = th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.Error(t, err)
	assert.Equal(t, "blocks were not found in db or RPC call: [2]", err.Error())

	// Emit a log and mine block #2
	_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)
	th.Client.Commit()

	// Assert block 2 is not yet in DB
	_, err = th.ORM.SelectBlockByNumber(2)
	require.Error(t, err)

	// getBlocksRange is able to retrieve block 2 by calling RPC
	rpcBlocks, err := th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 1, len(rpcBlocks))
	assert.Equal(t, 2, int(rpcBlocks[0].BlockNumber))

	// Emit a log and mine block #3
	_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(2)})
	require.NoError(t, err)
	th.Client.Commit()

	// Assert block 3 is not yet in DB
	_, err = th.ORM.SelectBlockByNumber(3)
	require.Error(t, err)

	// getBlocksRange is able to retrieve blocks 1 and 3, without retrieving block 2
	blockNums2 := []uint64{1, 3}
	rpcBlocks2, err := th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(rpcBlocks2))
	assert.Equal(t, 1, int(rpcBlocks2[0].BlockNumber))
	assert.Equal(t, 3, int(rpcBlocks2[1].BlockNumber))

	// after calling PollAndSaveLogs, block 2 & 3 are persisted in DB
	th.LogPoller.PollAndSaveLogs(testutils.Context(t), 1)
	block, err := th.ORM.SelectBlockByNumber(2)
	require.NoError(t, err)
	assert.Equal(t, 2, int(block.BlockNumber))
	block, err = th.ORM.SelectBlockByNumber(3)
	require.NoError(t, err)
	assert.Equal(t, 3, int(block.BlockNumber))

	// getBlocksRange should still be able to return block 2 by fetching from DB
	lpBlocks, err := th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 1, len(lpBlocks))
	assert.Equal(t, rpcBlocks[0].BlockNumber, lpBlocks[0].BlockNumber)
	assert.Equal(t, rpcBlocks[0].BlockHash, lpBlocks[0].BlockHash)

	// getBlocksRange return multiple blocks
	blockNums = []uint64{1, 2}
	blocks, err = th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 1, int(blocks[0].BlockNumber))
	assert.NotEmpty(t, blocks[0].BlockHash)
	assert.Equal(t, 2, int(blocks[1].BlockNumber))
	assert.NotEmpty(t, blocks[1].BlockHash)

	// getBlocksRange return blocks in requested order
	blockNums = []uint64{2, 1}
	reversedBlocks, err := th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, blocks[0].BlockNumber, reversedBlocks[1].BlockNumber)
	assert.Equal(t, blocks[0].BlockHash, reversedBlocks[1].BlockHash)
	assert.Equal(t, blocks[1].BlockNumber, reversedBlocks[0].BlockNumber)
	assert.Equal(t, blocks[1].BlockHash, reversedBlocks[0].BlockHash)

	// test RPC context cancellation
	ctx, cancel := context.WithCancel(testutils.Context(t))
	cancel()
	_, err = th.LogPoller.GetBlocksRange(ctx, blockNums)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	// test still works when qopts is cancelled
	// but context object is not
	ctx, cancel = context.WithCancel(testutils.Context(t))
	qopts := pg.WithParentCtx(ctx)
	cancel()
	_, err = th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums, qopts)
	require.NoError(t, err)
}

func TestGetReplayFromBlock(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, false, 2, 3, 2, 1000)
	// Commit a few blocks
	for i := 0; i < 10; i++ {
		th.Client.Commit()
	}

	// Nothing in the DB yet, should use whatever we specify.
	requested := int64(5)
	fromBlock, err := th.LogPoller.GetReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	assert.Equal(t, requested, fromBlock)

	// Do a poll, then we should have up to block 11 (blocks 0 & 1 are contract deployments, 2-10 logs).
	nextBlock := th.PollAndSaveLogs(testutils.Context(t), 1)
	require.Equal(t, int64(12), nextBlock)

	// Commit a few more so chain is ahead.
	for i := 0; i < 3; i++ {
		th.Client.Commit()
	}
	// Should take min(latest, requested), in this case latest.
	requested = int64(15)
	fromBlock, err = th.LogPoller.GetReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	latest, err := th.LogPoller.LatestBlock(pg.WithParentCtx(testutils.Context(t)))
	require.NoError(t, err)
	assert.Equal(t, latest.BlockNumber, fromBlock)

	// Should take min(latest, requested) in this case requested.
	requested = int64(7)
	fromBlock, err = th.LogPoller.GetReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	assert.Equal(t, requested, fromBlock)
}

func TestLogPoller_DBErrorHandling(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)
	chainID1 := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	o := logpoller.NewORM(chainID1, db, lggr, pgtest.NewQConfig(true))

	owner := testutils.MustNewSimTransactor(t)
	ethDB := rawdb.NewMemoryDatabase()
	ec := backends.NewSimulatedBackendWithDatabase(ethDB, map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	_, _, emitter, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	_, err = emitter.EmitLog1(owner, []*big.Int{big.NewInt(9)})
	require.NoError(t, err)
	_, err = emitter.EmitLog1(owner, []*big.Int{big.NewInt(7)})
	require.NoError(t, err)
	ec.Commit()
	ec.Commit()
	ec.Commit()

	lp := logpoller.NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID2), lggr, 1*time.Hour, false, 2, 3, 2, 1000)

	err = lp.Replay(ctx, 5) // block number too high
	require.ErrorContains(t, err, "Invalid replay block number")

	// Force a db error while loading the filters (tx aborted, already rolled back)
	require.Error(t, commonutils.JustError(db.Exec(`invalid query`)))
	go func() {
		err = lp.Replay(ctx, 2)
		assert.ErrorContains(t, err, "current transaction is aborted")
	}()

	time.Sleep(100 * time.Millisecond)
	require.NoError(t, lp.Start(ctx))
	require.Eventually(t, func() bool {
		return observedLogs.Len() >= 5
	}, 2*time.Second, 20*time.Millisecond)
	lp.Close()

	logMsgs := make(map[string]int)
	for _, obs := range observedLogs.All() {
		_, ok := logMsgs[obs.Entry.Message]
		if ok {
			logMsgs[(obs.Entry.Message)] = 1
		} else {
			logMsgs[(obs.Entry.Message)]++
		}
	}

	assert.Contains(t, logMsgs, "SQL ERROR")
	assert.Contains(t, logMsgs, "Failed loading filters in main logpoller loop, retrying later")
	assert.Contains(t, logMsgs, "Error executing replay, could not get fromBlock")
	assert.Contains(t, logMsgs, "Backup log poller ran before filters loaded, skipping")
}

type getLogErrData struct {
	From  string
	To    string
	Limit int
}

func TestTooManyLogResults(t *testing.T) {
	ctx := testutils.Context(t)
	ec := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr, obs := logger.TestObserved(t, zapcore.DebugLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	o := logpoller.NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	lp := logpoller.NewLogPoller(o, ec, lggr, 1*time.Hour, false, 2, 20, 10, 1000)
	expected := []int64{10, 5, 2, 1}

	clientErr := client.JsonError{
		Code:    -32005,
		Data:    getLogErrData{"0x100E698", "0x100E6D4", 10000},
		Message: "query returned more than 10000 results. Try with this block range [0x100E698, 0x100E6D4].",
	}

	call1 := ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(func(ctx context.Context, blockNumber *big.Int) (*evmtypes.Head, error) {
		if blockNumber == nil {
			return &evmtypes.Head{Number: 300}, nil // Simulate currentBlock = 300
		}
		return &evmtypes.Head{Number: blockNumber.Int64()}, nil
	})

	call2 := ec.On("FilterLogs", mock.Anything, mock.Anything).Return(func(ctx context.Context, fq ethereum.FilterQuery) (logs []types.Log, err error) {
		if fq.BlockHash != nil {
			return []types.Log{}, nil // succeed when single block requested
		}
		from := fq.FromBlock.Uint64()
		to := fq.ToBlock.Uint64()
		if to-from >= 4 {
			return []types.Log{}, &clientErr // return "too many results" error if block range spans 4 or more blocks
		}
		return logs, err
	})

	addr := testutils.NewAddress()
	err := lp.RegisterFilter(logpoller.Filter{"Integration test", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{addr}, 0})
	require.NoError(t, err)
	lp.PollAndSaveLogs(ctx, 5)
	block, err2 := o.SelectLatestBlock()
	require.NoError(t, err2)
	assert.Equal(t, int64(298), block.BlockNumber)

	logs := obs.FilterLevelExact(zapcore.WarnLevel).FilterMessageSnippet("halving block range batch size").FilterFieldKey("newBatchSize").All()
	// Should have tried again 3 times--first reducing batch size to 10, then 5, then 2
	require.Len(t, logs, 3)
	for i, s := range expected[:3] {
		assert.Equal(t, s, logs[i].ContextMap()["newBatchSize"])
	}

	obs.TakeAll()
	call1.Unset()
	call2.Unset()

	// Now jump to block 500, but return error no matter how small the block range gets.
	//  Should exit the loop with a critical error instead of hanging.
	call1.On("HeadByNumber", mock.Anything, mock.Anything).Return(func(ctx context.Context, blockNumber *big.Int) (*evmtypes.Head, error) {
		if blockNumber == nil {
			return &evmtypes.Head{Number: 500}, nil // Simulate currentBlock = 300
		}
		return &evmtypes.Head{Number: blockNumber.Int64()}, nil
	})
	call2.On("FilterLogs", mock.Anything, mock.Anything).Return(func(ctx context.Context, fq ethereum.FilterQuery) (logs []types.Log, err error) {
		if fq.BlockHash != nil {
			return []types.Log{}, nil // succeed when single block requested
		}
		return []types.Log{}, &clientErr // return "too many results" error if block range spans 4 or more blocks
	})

	lp.PollAndSaveLogs(ctx, 298)
	block, err2 = o.SelectLatestBlock()
	require.NoError(t, err2)
	assert.Equal(t, int64(298), block.BlockNumber)
	warns := obs.FilterMessageSnippet("halving block range").FilterLevelExact(zapcore.WarnLevel).All()
	crit := obs.FilterMessageSnippet("failed to retrieve logs").FilterLevelExact(zapcore.DPanicLevel).All()
	require.Len(t, warns, 4)
	for i, s := range expected {
		assert.Equal(t, s, warns[i].ContextMap()["newBatchSize"])
	}

	require.Len(t, crit, 1)
	assert.Contains(t, crit[0].Message, "Too many log results in a single block")
}

func Test_PollAndQueryFinalizedBlocks(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	firstBatchLen := 3
	secondBatchLen := 5

	th := SetupTH(t, true, 2, 3, 2, 1000)

	eventSig := EmitterABI.Events["Log1"].ID
	err := th.LogPoller.RegisterFilter(logpoller.Filter{
		Name:      "GetBlocks Test",
		EventSigs: []common.Hash{eventSig},
		Addresses: []common.Address{th.EmitterAddress1}},
	)
	require.NoError(t, err)

	// Generate block that will be finalized
	for i := 0; i < firstBatchLen; i++ {
		_, err1 := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}

	// Mark current head as finalized
	h := th.Client.Blockchain().CurrentHeader()
	th.Client.Blockchain().SetFinalized(h)

	// Generate next blocks, not marked as finalized
	for i := 0; i < secondBatchLen; i++ {
		_, err1 := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}

	currentBlock := th.PollAndSaveLogs(ctx, 1)
	require.Equal(t, int(currentBlock), firstBatchLen+secondBatchLen+2)

	finalizedLogs, err := th.LogPoller.LogsDataWordGreaterThan(
		eventSig,
		th.EmitterAddress1,
		0,
		common.Hash{},
		logpoller.Finalized,
	)
	require.NoError(t, err)
	require.Len(t, finalizedLogs, firstBatchLen)

	numberOfConfirmations := 1
	logsByConfs, err := th.LogPoller.LogsDataWordGreaterThan(
		eventSig,
		th.EmitterAddress1,
		0,
		common.Hash{},
		logpoller.Confirmations(numberOfConfirmations),
	)
	require.NoError(t, err)
	require.Len(t, logsByConfs, firstBatchLen+secondBatchLen-numberOfConfirmations)
}

func Test_PollAndSavePersistsFinalityInBlocks(t *testing.T) {
	ctx := testutils.Context(t)
	numberOfBlocks := 10

	tests := []struct {
		name                   string
		useFinalityTag         bool
		finalityDepth          int64
		expectedFinalizedBlock int64
	}{
		{
			name:                   "using fixed finality depth",
			useFinalityTag:         false,
			finalityDepth:          2,
			expectedFinalizedBlock: int64(numberOfBlocks - 2),
		},
		{
			name:                   "setting last finalized block number to 0 if finality is too deep",
			useFinalityTag:         false,
			finalityDepth:          20,
			expectedFinalizedBlock: 0,
		},
		{
			name:                   "using finality from chain",
			useFinalityTag:         true,
			finalityDepth:          0,
			expectedFinalizedBlock: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, tt.useFinalityTag, tt.finalityDepth, 3, 2, 1000)
			// Should return error before the first poll and save
			_, err := th.LogPoller.LatestBlock()
			require.Error(t, err)

			// Mark first block as finalized
			h := th.Client.Blockchain().CurrentHeader()
			th.Client.Blockchain().SetFinalized(h)

			// Create a couple of blocks
			for i := 0; i < numberOfBlocks-1; i++ {
				th.Client.Commit()
			}

			th.PollAndSaveLogs(ctx, 1)

			latestBlock, err := th.LogPoller.LatestBlock()
			require.NoError(t, err)
			require.Equal(t, int64(numberOfBlocks), latestBlock.BlockNumber)
			require.Equal(t, tt.expectedFinalizedBlock, latestBlock.FinalizedBlockNumber)
		})
	}
}

func Test_CreatedAfterQueriesWithBackfill(t *testing.T) {
	emittedLogs := 60
	ctx := testutils.Context(t)

	tests := []struct {
		name          string
		finalityDepth int64
		finalityTag   bool
	}{
		{
			name:          "fixed finality depth without finality tag",
			finalityDepth: 10,
			finalityTag:   false,
		},
		{
			name:          "chain finality in use",
			finalityDepth: 0,
			finalityTag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, tt.finalityTag, tt.finalityDepth, 3, 2, 1000)

			header, err := th.Client.HeaderByNumber(ctx, nil)
			require.NoError(t, err)

			genesisBlockTime := time.UnixMilli(int64(header.Time))

			// Emit some logs in blocks
			for i := 0; i < emittedLogs; i++ {
				_, err2 := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
				require.NoError(t, err2)
				th.Client.Commit()
			}

			// First PollAndSave, no filters are registered
			currentBlock := th.PollAndSaveLogs(ctx, 1)

			err = th.LogPoller.RegisterFilter(logpoller.Filter{
				Name:      "Test Emitter",
				EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
				Addresses: []common.Address{th.EmitterAddress1},
			})
			require.NoError(t, err)

			// Emit blocks to cover finality depth, because backup always backfill up to the one block before last finalized
			for i := 0; i < int(tt.finalityDepth)+1; i++ {
				bh := th.Client.Commit()
				markBlockAsFinalizedByHash(t, th, bh)
			}

			// LogPoller should backfill entire history
			th.LogPoller.BackupPollAndSaveLogs(ctx, 100)
			require.NoError(t, err)

			// Make sure that all logs are backfilled
			logs, err := th.LogPoller.Logs(
				0,
				currentBlock,
				EmitterABI.Events["Log1"].ID,
				th.EmitterAddress1,
				pg.WithParentCtx(testutils.Context(t)),
			)
			require.NoError(t, err)
			require.Len(t, logs, emittedLogs)

			// We should get all the logs by the block_timestamp
			logs, err = th.LogPoller.LogsCreatedAfter(
				EmitterABI.Events["Log1"].ID,
				th.EmitterAddress1,
				genesisBlockTime,
				0,
				pg.WithParentCtx(testutils.Context(t)),
			)
			require.NoError(t, err)
			require.Len(t, logs, emittedLogs)
		})
	}
}

func Test_PruneOldBlocks(t *testing.T) {
	ctx := testutils.Context(t)

	tests := []struct {
		name                     string
		keepFinalizedBlocksDepth int64
		blockToCreate            int
		blocksLeft               int
		wantErr                  bool
	}{
		{
			name:                     "returns error if no blocks yet",
			keepFinalizedBlocksDepth: 10,
			blockToCreate:            0,
			wantErr:                  true,
		},
		{
			name:                     "returns if there is not enough blocks in the db",
			keepFinalizedBlocksDepth: 11,
			blockToCreate:            10,
			blocksLeft:               10,
		},
		{
			name:                     "prunes matching blocks",
			keepFinalizedBlocksDepth: 1000,
			blockToCreate:            2000,
			blocksLeft:               1010, // last finalized block is 10 block behind
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := SetupTH(t, true, 0, 3, 2, tt.keepFinalizedBlocksDepth)

			for i := 1; i <= tt.blockToCreate; i++ {
				err := th.ORM.InsertBlock(utils.RandomBytes32(), int64(i+10), time.Now(), int64(i))
				require.NoError(t, err)
			}

			if tt.wantErr {
				require.Error(t, th.LogPoller.PruneOldBlocks(ctx))
				return
			}

			require.NoError(t, th.LogPoller.PruneOldBlocks(ctx))
			blocks, err := th.ORM.GetBlocksRange(0, math.MaxInt64, pg.WithParentCtx(ctx))
			require.NoError(t, err)
			assert.Len(t, blocks, tt.blocksLeft)
		})
	}
}

func markBlockAsFinalized(t *testing.T, th TestHarness, blockNumber int64) {
	b, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(blockNumber))
	require.NoError(t, err)
	th.Client.Blockchain().SetFinalized(b.Header())
}

func markBlockAsFinalizedByHash(t *testing.T, th TestHarness, blockHash common.Hash) {
	b, err := th.Client.BlockByHash(testutils.Context(t), blockHash)
	require.NoError(t, err)
	th.Client.Blockchain().SetFinalized(b.Header())
}
