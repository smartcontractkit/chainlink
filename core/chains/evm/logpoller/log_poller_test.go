package logpoller_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

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
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func logRuntime(t *testing.T, start time.Time) {
	t.Log("runtime", time.Since(start))
}

func TestPopulateLoadedDB(t *testing.T) {
	t.Skip("only for local load testing and query analysis")
	lggr := logger.TestLogger(t)
	_, db := heavyweight.FullTestDBV2(t, "logs_scale", nil)
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewQConfig(true))
	event1 := EmitterABI.Events["Log1"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// We start at 1 just so block number > 0
	for j := 1; j < 1000; j++ {
		var logs []logpoller.Log
		// Max we can insert per batch
		for i := 0; i < 1000; i++ {
			addr := address1
			if (i+(1000*j))%2 == 0 {
				addr = address2
			}
			logs = append(logs, logpoller.Log{
				EvmChainId:  utils.NewBig(chainID),
				LogIndex:    1,
				BlockHash:   common.HexToHash(fmt.Sprintf("0x%d", i+(1000*j))),
				BlockNumber: int64(i + (1000 * j)),
				EventSig:    event1,
				Topics:      [][]byte{event1[:], logpoller.EvmWord(uint64(i + 1000*j)).Bytes()},
				Address:     addr,
				TxHash:      common.HexToHash("0x1234"),
				Data:        logpoller.EvmWord(uint64(i + 1000*j)).Bytes(),
			})
		}
		require.NoError(t, o.InsertLogs(logs))
	}
	func() {
		defer logRuntime(t, time.Now())
		_, err := o.SelectLogsByBlockRangeFilter(750000, 800000, address1, event1)
		require.NoError(t, err)
	}()
	func() {
		defer logRuntime(t, time.Now())
		_, err = o.SelectLatestLogEventSigsAddrsWithConfs(0, []common.Address{address1}, []common.Hash{event1}, 0)
		require.NoError(t, err)
	}()

	// Confirm all the logs.
	require.NoError(t, o.InsertBlock(common.HexToHash("0x10"), 1000000, time.Now()))
	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectDataWordRange(address1, event1, 0, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err)
		// 10 since every other log is for address1
		assert.Equal(t, 10, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectIndexedLogs(address2, event1, 1, []common.Hash{logpoller.EvmWord(500000), logpoller.EvmWord(500020)}, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectIndexLogsTopicRange(address1, event1, 1, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err)
		assert.Equal(t, 10, len(lgs))
	}()
}

func TestLogPoller_Integration(t *testing.T) {
	th := SetupTH(t, 2, 3, 2)
	th.Client.Commit() // Block 2. Ensure we have finality number of blocks

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"Integration test", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1}})
	require.NoError(t, err)
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Addresses, 1)
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Topics, 1)

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(testutils.Context(t)))
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Addresses, 1)
	require.Len(t, th.LogPoller.Filter(nil, nil, nil).Topics, 1)

	// Emit some logs in blocks 3->7.
	for i := 0; i < 5; i++ {
		_, err := th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		_, err = th.Emitter1.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		th.Client.Commit()
	}
	// The poller starts on a new chain at latest-finality (5 in this case),
	// replay to ensure we get all the logs.
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 1))

	// We should immediately have all those Log1 logs.
	logs, err := th.LogPoller.Logs(2, 7, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 5, len(logs))
	// Now let's update the Filter and replay to get Log2 logs.
	err = th.LogPoller.RegisterFilter(logpoller.Filter{
		"Emitter - log2", []common.Hash{EmitterABI.Events["Log2"].ID},
		[]common.Address{th.EmitterAddress1},
	})
	require.NoError(t, err)
	// Replay an invalid block should error
	assert.Error(t, th.LogPoller.Replay(testutils.Context(t), 0))
	assert.Error(t, th.LogPoller.Replay(testutils.Context(t), 20))
	// Replay only from block 4, so we should see logs in block 4,5,6,7 (4 logs)
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 4))

	// We should immediately see 4 logs2 logs.
	logs, err = th.LogPoller.Logs(2, 7, EmitterABI.Events["Log2"].ID, th.EmitterAddress1)
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
	th := SetupTH(t, 2, 3, 2)
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
		[]common.Address{th.EmitterAddress1}}
	err := th.LogPoller.RegisterFilter(filter1)
	require.NoError(t, err)

	filters, err := th.ORM.LoadFilters(pg.WithParentCtx(testutils.Context(t)))
	require.NoError(t, err)
	require.Equal(t, 1, len(filters))
	require.Equal(t, filter1, filters["filter1"])

	err = th.LogPoller.RegisterFilter(
		logpoller.Filter{"filter2",
			[]common.Hash{EmitterABI.Events["Log1"].ID},
			[]common.Address{th.EmitterAddress2}})
	require.NoError(t, err)

	defer th.LogPoller.UnregisterFilter("filter1", nil)
	defer th.LogPoller.UnregisterFilter("filter2", nil)

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
	receipts := rawdb.ReadReceipts(th.EthDB, h.Hash(), h.Number.Uint64(), params.AllEthashProtocolChanges)
	require.NotZero(t, receipts.Len())

	// Simulate a situation where the rpc server has a block, but no logs available for it yet
	//  this can't happen with geth itself, but can with other clients.
	rawdb.WriteReceipts(th.EthDB, h.Hash(), h.Number.Uint64(), types.Receipts{}) // wipes out all logs for block 34

	body := rawdb.ReadBody(th.EthDB, h.Hash(), h.Number.Uint64())
	require.Equal(t, 3, len(body.Transactions))
	txs := body.Transactions                 // save transactions for later
	body.Transactions = types.Transactions{} // number of tx's must match # of logs for GetLogs() to succeed
	rawdb.WriteBody(th.EthDB, h.Hash(), h.Number.Uint64(), body)

	currentBlock := th.PollAndSaveLogs(ctx, 1)
	assert.Equal(t, int64(35), currentBlock)

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
	logs, err := th.LogPoller.Logs(34, 34, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	th.Client.Commit()
	th.Client.Commit()

	// Run ordinary poller + backup poller at least once
	currentBlock, _ = th.LogPoller.LatestBlock()
	th.LogPoller.PollAndSaveLogs(ctx, currentBlock+1)
	th.LogPoller.BackupPollAndSaveLogs(ctx, 100)
	currentBlock, _ = th.LogPoller.LatestBlock()

	require.Equal(t, int64(37), currentBlock+1)

	// logs still shouldn't show up, because we don't want to backfill the last finalized log
	//  to help with reorg detection
	logs, err = th.LogPoller.Logs(34, 34, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	th.Client.Commit()

	// Run ordinary poller + backup poller at least once more
	th.LogPoller.PollAndSaveLogs(ctx, currentBlock+1)
	th.LogPoller.BackupPollAndSaveLogs(ctx, 100)
	currentBlock, _ = th.LogPoller.LatestBlock()

	require.Equal(t, int64(38), currentBlock+1)

	// all 3 logs in block 34 should show up now, thanks to backup logger
	logs, err = th.LogPoller.Logs(30, 37, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 5, len(logs))
	logs, err = th.LogPoller.Logs(34, 34, EmitterABI.Events["Log2"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(logs))
	logs, err = th.LogPoller.Logs(32, 36, EmitterABI.Events["Log1"].ID, th.EmitterAddress2)
	require.NoError(t, err)
	assert.Equal(t, 1, len(logs))
}

func TestLogPoller_BlockTimestamps(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	th := SetupTH(t, 2, 3, 2)

	addresses := []common.Address{th.EmitterAddress1, th.EmitterAddress2}
	topics := []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"convertLogs", topics, addresses})
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

	lb, _ := th.LogPoller.LatestBlock()
	th.PollAndSaveLogs(context.Background(), lb+1)
	lg1, err := th.LogPoller.Logs(0, 20, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	lg2, err := th.LogPoller.Logs(0, 20, EmitterABI.Events["Log2"].ID, th.EmitterAddress2)
	require.NoError(t, err)

	// Logs should have correct timestamps
	b, _ := th.Client.BlockByHash(context.Background(), lg1[0].BlockHash)
	t.Log(len(lg1), lg1[0].BlockTimestamp)
	assert.Equal(t, int64(b.Time()), lg1[0].BlockTimestamp.UTC().Unix(), time1)
	b2, _ := th.Client.BlockByHash(context.Background(), lg2[0].BlockHash)
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
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_logs_evm_chain_id_fkey DEFERRED`)))
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
		lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, int64(finalityDepth), 3, 2, 1000)
		for i := 0; i < finalityDepth; i++ { // Have enough blocks that we could reorg the full finalityDepth-1.
			ec.Commit()
		}
		currentBlock := int64(1)
		lp.PollAndSaveLogs(testutils.Context(t), currentBlock)
		currentBlock, err = lp.LatestBlock()
		require.NoError(t, err)
		matchesGeth := func() bool {
			// Check every block is identical
			latest, err := ec.BlockByNumber(testutils.Context(t), nil)
			require.NoError(t, err)
			for i := 1; i < int(latest.NumberU64()); i++ {
				ourBlock, err := lp.BlockByNumber(int64(i))
				require.NoError(t, err)
				gethBlock, err := ec.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
				require.NoError(t, err)
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
					latest, err := ec.BlockByNumber(testutils.Context(t), nil)
					require.NoError(t, err)
					t.Log("mined block", latest.Hash())
				}
			} else {
				// Reorg blocks
				latest, err := ec.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)
				reorgedBlock := big.NewInt(0).Sub(latest.Number(), big.NewInt(int64(mineOrReorg[i])))
				reorg, err := ec.BlockByNumber(testutils.Context(t), reorgedBlock)
				require.NoError(t, err)
				require.NoError(t, ec.Fork(testutils.Context(t), reorg.Hash()))
				t.Logf("Reorging from (%v, %x) back to (%v, %x)\n", latest.NumberU64(), latest.Hash(), reorgedBlock.Uint64(), reorg.Hash())
				// Actually need to change the block here to trigger the reorg.
				_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(1)})
				require.NoError(t, err)
				for j := 0; j < int(mineOrReorg[i]+1); j++ { // Need +1 to make it actually longer height so we detect it.
					ec.Commit()
				}
				latest, err = ec.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)
				t.Logf("New latest (%v, %x), latest parent %x)\n", latest.NumberU64(), latest.Hash(), latest.ParentHash())
			}
			lp.PollAndSaveLogs(testutils.Context(t), currentBlock)
			currentBlock, err = lp.LatestBlock()
			require.NoError(t, err)
		}
		return matchesGeth()
	}, gen.SliceOfN(numChainInserts, gen.UInt64Range(1, uint64(finalityDepth-1))))) // Max reorg depth is finality depth - 1
	p.TestingRun(t)
}

func TestLogPoller_PollAndSaveLogs(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	// Set up a log poller listening for log emitter logs.
	err := th.LogPoller.RegisterFilter(logpoller.Filter{
		"Test Emitter 1 & 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID},
		[]common.Address{th.EmitterAddress1, th.EmitterAddress2},
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
	newStart = th.PollAndSaveLogs(testutils.Context(t), newStart)
	assert.Equal(t, int64(18), newStart)
	lgs, err = th.ORM.SelectLogsByBlockRange(11, 17)
	require.NoError(t, err)
	assert.Equal(t, 7, len(lgs))
	th.assertHaveCanonical(t, 15, 16)
	th.assertDontHave(t, 11, 14) // Do not expect to save backfilled blocks.

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
}

func TestLogPoller_LoadFilters(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	filter1 := logpoller.Filter{"first Filter", []common.Hash{
		EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}}
	filter2 := logpoller.Filter{"second Filter", []common.Hash{
		EmitterABI.Events["Log2"].ID, EmitterABI.Events["Log3"].ID}, []common.Address{th.EmitterAddress2}}
	filter3 := logpoller.Filter{"third Filter", []common.Hash{
		EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}}

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
}

func TestLogPoller_GetBlocks_Range(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"GetBlocks Test", []common.Hash{
		EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}},
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
	th := SetupTH(t, 2, 3, 2)
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
	latest, err := th.LogPoller.LatestBlock()
	require.NoError(t, err)
	assert.Equal(t, latest, fromBlock)

	// Should take min(latest, requested) in this case requested.
	requested = int64(7)
	fromBlock, err = th.LogPoller.GetReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	assert.Equal(t, requested, fromBlock)
}

func TestLogPoller_DBErrorHandling(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
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

	lp := logpoller.NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID2), lggr, 1*time.Hour, 2, 3, 2, 1000)

	err = lp.Replay(ctx, 5) // block number too high
	require.ErrorContains(t, err, "Invalid replay block number")

	// Force a db error while loading the filters (tx aborted, already rolled back)
	require.Error(t, utils.JustError(db.Exec(`invalid query`)))
	go func() {
		err = lp.Replay(ctx, 2)
		assert.ErrorContains(t, err, "current transaction is aborted")
	}()

	time.Sleep(100 * time.Millisecond)
	lp.Start(ctx)
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
	assert.Contains(t, logMsgs, "backup log poller ran before filters loaded, skipping")
}
