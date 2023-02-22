package logpoller_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
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
	require.NoError(t, o.InsertBlock(common.HexToHash("0x10"), 1000000))
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
	th := logpoller.SetupTH(t, 2, 3, 2)
	th.Client.Commit() // Block 2. Ensure we have finality number of blocks

	err := th.LogPoller.RegisterFilter(logpoller.Filter{"Integration test", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1}})
	require.NoError(t, err)
	require.Len(t, th.LogPoller.Filter().Addresses, 1)
	require.Len(t, th.LogPoller.Filter().Topics, 1)

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload filter from db.
	require.NoError(t, th.LogPoller.Start(testutils.Context(t)))
	require.Len(t, th.LogPoller.Filter().Addresses, 1)
	require.Len(t, th.LogPoller.Filter().Topics, 1)

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
	// Now let's update the filter and replay to get Log2 logs.
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

	// Cancelling a replay should return an error synchronously.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.True(t, errors.Is(th.LogPoller.Replay(ctx, 4), logpoller.ErrReplayAbortedByClient))

	require.NoError(t, th.LogPoller.Close())
}

// Simulate a badly behaving rpc server, where unfinalized blocks can return different logs
// for the same block hash.  We should be able to handle this without missing any logs, as
// long as the logs returned for finalized blocks are consistent.
func Test_BackupLogPoller(t *testing.T) {
	th := logpoller.SetupTH(t, 2, 3, 2)
	// later, we will need at least 32 blocks filled with logs for cache invalidation
	for i := int64(0); i < 32; i++ {
		// to invalidate geth's internal read-cache, a matching log must be found in the bloom filter
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

	defer th.LogPoller.UnregisterFilter("filter1")
	defer th.LogPoller.UnregisterFilter("filter2")

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

	currentBlock := th.LogPoller.PollAndSaveLogs(ctx, 1)
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
	require.NoError(t, th.LogPoller.Restart(testutils.Context(t)))
	time.Sleep(500 * time.Millisecond)
	require.NoError(t, th.LogPoller.Close())
	currentBlock = th.LogPoller.GetCurrentBlock()

	require.Equal(t, int64(37), currentBlock)

	// logs still shouldn't show up, because we don't want to backfill the last finalized log
	//  to help with reorg detection
	logs, err = th.LogPoller.Logs(34, 34, EmitterABI.Events["Log1"].ID, th.EmitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	th.Client.Commit()

	// Run ordinary poller + backup poller at least once more
	require.NoError(t, th.LogPoller.Restart(testutils.Context(t)))
	time.Sleep(500 * time.Millisecond)
	require.NoError(t, th.LogPoller.Close())
	currentBlock = th.LogPoller.GetCurrentBlock()

	require.Equal(t, int64(38), currentBlock)

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
