package logpoller

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/jackc/pgconn"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

func GenLog(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address) Log {
	return Log{
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

func assertDontHave(t *testing.T, start, end int, orm *ORM) {
	for i := start; i < end; i++ {
		_, err := orm.SelectBlockByNumber(int64(i))
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	}
}

func assertHaveCanonical(t *testing.T, start, end int, ec *backends.SimulatedBackend, orm *ORM) {
	for i := start; i < end; i++ {
		blk, err := orm.SelectBlockByNumber(int64(i))
		require.NoError(t, err, "block %v", i)
		chainBlk, err := ec.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
		require.NoError(t, err)
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes(), "block %v", i)
	}
}

func TestLogPoller_Batching(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)
	var logs []Log
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
		orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
		// Note this property test is run concurrently and the sim is not threadsafe.
		ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
			owner.From: {
				Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
			},
		}, 10e6)
		_, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
		require.NoError(t, err)
		lp := NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, int64(finalityDepth), 3, 2, 1000)
		for i := 0; i < finalityDepth; i++ { // Have enough blocks that we could reorg the full finalityDepth-1.
			ec.Commit()
		}
		currentBlock := int64(1)
		currentBlock = lp.PollAndSaveLogs(testutils.Context(t), currentBlock)
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
			currentBlock = lp.PollAndSaveLogs(testutils.Context(t), currentBlock)
		}
		return matchesGeth()
	}, gen.SliceOfN(numChainInserts, gen.UInt64Range(1, uint64(finalityDepth-1))))) // Max reorg depth is finality depth - 1
	p.TestingRun(t)
}

func TestLogPoller_PollAndSaveLogs(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	// Set up a log poller listening for log emitter logs.
	err := th.LogPoller.RegisterFilter(Filter{
		"Test Emitter 1 & 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID},
		[]common.Address{th.EmitterAddress1, th.EmitterAddress2},
	})
	require.NoError(t, err)

	b, err := th.Client.BlockByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), b.NumberU64())

	// Test scenario: single block in chain, no logs.
	// Chain genesis <- 1
	// DB: empty
	newStart := th.LogPoller.PollAndSaveLogs(testutils.Context(t), 1)
	assert.Equal(t, int64(2), newStart)

	// We expect to have saved block 1.
	lpb, err := th.ORM.SelectBlockByNumber(1)
	require.NoError(t, err)
	assert.Equal(t, lpb.BlockHash, b.Hash())
	assert.Equal(t, lpb.BlockNumber, int64(b.NumberU64()))
	assert.Equal(t, int64(1), int64(b.NumberU64()))

	// No logs.
	lgs, err := th.ORM.SelectLogsByBlockRange(1, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	assertHaveCanonical(t, 1, 1, th.Client, th.ORM)

	// Polling again should be a noop, since we are at the latest.
	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
	assert.Equal(t, int64(2), newStart)
	latest, err := th.ORM.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(1), latest.BlockNumber)
	assertHaveCanonical(t, 1, 1, th.Client, th.ORM)

	// Test scenario: one log 2 block chain.
	// Chain gen <- 1 <- 2 (L1)
	// DB: 1
	_, err = th.Emitter1.EmitLog1(th.Owner, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)
	th.Client.Commit()

	// Polling should get us the L1 log.
	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
	assert.Equal(t, int64(3), newStart)
	latest, err = th.ORM.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(2), latest.BlockNumber)
	lgs, err = th.ORM.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, th.EmitterAddress1, lgs[0].Address)
	assert.Equal(t, latest.BlockHash, lgs[0].BlockHash)
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

	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
	assert.Equal(t, int64(4), newStart)
	latest, err = th.ORM.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(3), latest.BlockNumber)
	lgs, err = th.ORM.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[0].Data)
	assertHaveCanonical(t, 1, 3, th.Client, th.ORM)

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
	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
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
	assertHaveCanonical(t, 1, 1, th.Client, th.ORM)
	assertHaveCanonical(t, 3, 4, th.Client, th.ORM)
	assertDontHave(t, 2, 2, th.ORM) // 2 gets backfilled

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

	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
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
	assertHaveCanonical(t, 1, 1, th.Client, th.ORM)
	assertDontHave(t, 2, 2, th.ORM) // 2 gets backfilled
	assertHaveCanonical(t, 3, 6, th.Client, th.ORM)

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
	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
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
	assertDontHave(t, 7, 7, th.ORM) // Do not expect to save backfilled blocks.
	assertHaveCanonical(t, 8, 10, th.Client, th.ORM)

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
	newStart = th.LogPoller.PollAndSaveLogs(testutils.Context(t), newStart)
	assert.Equal(t, int64(18), newStart)
	lgs, err = th.ORM.SelectLogsByBlockRange(11, 17)
	require.NoError(t, err)
	assert.Equal(t, 7, len(lgs))
	assertHaveCanonical(t, 15, 16, th.Client, th.ORM)
	assertDontHave(t, 11, 14, th.ORM) // Do not expect to save backfilled blocks.
}

func TestLogPoller_Logs(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)
	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1-3
	require.NoError(t, th.ORM.InsertLogs([]Log{
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

// Validate that filters stored in log_filters_table match the filters stored in memory
func validateFiltersTable(t *testing.T, lp *logPoller, orm *ORM) {
	filters, err := orm.LoadFilters()
	require.NoError(t, err)
	require.Equal(t, len(filters), len(lp.filters))
	for name, dbFilter := range filters {
		dbFilter := dbFilter
		memFilter, ok := lp.filters[name]
		require.True(t, ok)
		assert.True(t, memFilter.contains(&dbFilter),
			fmt.Sprintf("in-memory filter %s is missing some addresses or events from db filter table", name))
		assert.True(t, dbFilter.contains(&memFilter),
			fmt.Sprintf("db filter table %s is missing some addresses or events from in-memory filter", name))
	}
}

func TestLogPoller_RegisterFilter(t *testing.T) {
	t.Parallel()
	a1 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	a2 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	lp := NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	filter := Filter{"test filter", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}}
	err := lp.RegisterFilter(filter)
	require.Error(t, err, "RegisterFilter failed to save filter to db")
	require.Equal(t, 1, observedLogs.Len())
	assertForeignConstraintError(t, observedLogs.All()[0], "evm_log_poller_filters", "evm_log_poller_filters_evm_chain_id_fkey")

	db.Close()
	db = pgtest.NewSqlxDB(t)
	lggr = logger.TestLogger(t)
	orm = NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	// disable check that chain id exists for rest of tests
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_filters_evm_chain_id_fkey DEFERRED`)))
	// Set up a test chain with a log emitting contract deployed.

	lp = NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	// We expect a zero filter if nothing registered yet.
	f := lp.filter(nil, nil, nil)
	require.Equal(t, 1, len(f.Addresses))
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), f.Addresses[0])

	err = lp.RegisterFilter(Filter{"Emitter Log 1", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, lp.Filter().Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe EventSigs
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter().Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe Addresses
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2 dupe", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter().Topics)
	validateFiltersTable(t, lp, orm)

	// Address required.
	err = lp.RegisterFilter(Filter{"no address", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{}})
	require.Error(t, err)
	// Event required
	err = lp.RegisterFilter(Filter{"No event", []common.Hash{}, []common.Address{a1}})
	require.Error(t, err)
	validateFiltersTable(t, lp, orm)

	// Removing non-existence filter should error.
	err = lp.UnregisterFilter("filter doesn't exist")
	require.Error(t, err)

	// Check that all filters are still there
	_, ok := lp.filters["Emitter Log 1"]
	require.True(t, ok, "'Emitter Log 1 filter' missing")
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.True(t, ok, "'Emitter Log 1 + 2' filter missing")
	_, ok = lp.filters["Emitter Log 1 + 2 dupe"]
	require.True(t, ok, "'Emitter Log 1 + 2 dupe' filter missing")

	// Removing an existing filter should remove it from both memory and db
	err = lp.UnregisterFilter("Emitter Log 1 + 2")
	require.NoError(t, err)
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.False(t, ok, "'Emitter Log 1 filter' should have been removed by UnregisterFilter()")
	require.Len(t, lp.filters, 2)
	validateFiltersTable(t, lp, orm)

	err = lp.UnregisterFilter("Emitter Log 1 + 2 dupe")
	require.NoError(t, err)
	err = lp.UnregisterFilter("Emitter Log 1")
	require.NoError(t, err)
	assert.Len(t, lp.filters, 0)
	filters, err := lp.orm.LoadFilters()
	require.NoError(t, err)
	assert.Len(t, filters, 0)

	// Make sure cache was invalidated
	assert.Len(t, lp.Filter().Addresses, 1)
	assert.Equal(t, lp.Filter().Addresses[0], common.HexToAddress("0x0000000000000000000000000000000000000000"))
	assert.Len(t, lp.Filter().Topics, 1)
	assert.Len(t, lp.Filter().Topics[0], 0)
}

func assertForeignConstraintError(t *testing.T, observedLog observer.LoggedEntry,
	table string, constraint string) {

	assert.Equal(t, "SQL ERROR", observedLog.Entry.Message)

	field := observedLog.Context[0]
	require.Equal(t, zapcore.ErrorType, field.Type)
	err, ok := field.Interface.(error)
	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	require.True(t, ok)
	assert.Equal(t, "23503", pgErr.SQLState()) // foreign key constraint violation code
	assert.Equal(t, table, pgErr.TableName)
	assert.Equal(t, constraint, pgErr.ConstraintName)
}

func TestLogPoller_DBErrorHandling(t *testing.T) {
	t.Parallel()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
	chainID1 := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	o := NewORM(chainID1, db, lggr, pgtest.NewQConfig(true))

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

	lp := NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID2), lggr, 1*time.Hour, 2, 3, 2, 1000)
	ctx, cancelReplay := context.WithCancel(testutils.Context(t))
	lp.ctx, lp.cancel = context.WithCancel(testutils.Context(t))
	defer cancelReplay()
	defer lp.cancel()

	err = lp.Replay(ctx, 5) // block number too high
	require.ErrorContains(t, err, "Invalid replay block number")

	// Force a db error while loading the filters (tx aborted, already rolled back)
	require.Error(t, utils.JustError(db.Exec(`invalid query`)))
	go func() {
		err = lp.Replay(ctx, 2)
		assert.Error(t, err, ErrReplayAbortedByClient)
	}()

	time.Sleep(100 * time.Millisecond)
	go lp.run()
	require.Eventually(t, func() bool {
		return observedLogs.Len() >= 5
	}, 2*time.Second, 20*time.Millisecond)
	lp.cancel()
	lp.Close()
	<-lp.done

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

func TestLogPoller_LoadFilters(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	filter1 := Filter{"first filter", []common.Hash{
		EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}}
	filter2 := Filter{"second filter", []common.Hash{
		EmitterABI.Events["Log2"].ID, EmitterABI.Events["Log3"].ID}, []common.Address{th.EmitterAddress2}}
	filter3 := Filter{"third filter", []common.Hash{
		EmitterABI.Events["Log1"].ID}, []common.Address{th.EmitterAddress1, th.EmitterAddress2}}

	assert.True(t, filter1.contains(nil))
	assert.False(t, filter1.contains(&filter2))
	assert.False(t, filter2.contains(&filter1))
	assert.True(t, filter1.contains(&filter3))

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

	filter, ok := filters["first filter"]
	require.True(t, ok)
	assert.True(t, filter.contains(&filter1))
	assert.True(t, filter1.contains(&filter))

	filter, ok = filters["second filter"]
	require.True(t, ok)
	assert.True(t, filter.contains(&filter2))
	assert.True(t, filter2.contains(&filter))

	filter, ok = filters["third filter"]
	require.True(t, ok)
	assert.True(t, filter.contains(&filter3))
	assert.True(t, filter3.contains(&filter))
}

func TestLogPoller_GetBlocks_Range(t *testing.T) {
	t.Parallel()
	th := SetupTH(t, 2, 3, 2)

	err := th.LogPoller.RegisterFilter(Filter{"GetBlocks Test", []common.Hash{
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

	// getBlocksRange returns blocks with a nil client
	th.LogPoller.ec = nil
	blockNums = []uint64{1, 2}
	blocks, err = th.LogPoller.GetBlocksRange(testutils.Context(t), blockNums)
	require.NoError(t, err)
	assert.Equal(t, 1, int(blocks[0].BlockNumber))
	assert.NotEmpty(t, blocks[0].BlockHash)
	assert.Equal(t, 2, int(blocks[1].BlockNumber))
	assert.NotEmpty(t, blocks[1].BlockHash)
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
	fromBlock, err := th.LogPoller.getReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	assert.Equal(t, requested, fromBlock)

	// Do a poll, then we should have up to block 11 (blocks 0 & 1 are contract deployments, 2-10 logs).
	nextBlock := th.LogPoller.PollAndSaveLogs(testutils.Context(t), 1)
	require.Equal(t, int64(12), nextBlock)

	// Commit a few more so chain is ahead.
	for i := 0; i < 3; i++ {
		th.Client.Commit()
	}
	// Should take min(latest, requested), in this case latest.
	requested = int64(15)
	fromBlock, err = th.LogPoller.getReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	latest, err := th.LogPoller.LatestBlock()
	require.NoError(t, err)
	assert.Equal(t, latest, fromBlock)

	// Should take min(latest, requested) in this case requested.
	requested = int64(7)
	fromBlock, err = th.LogPoller.getReplayFromBlock(testutils.Context(t), requested)
	require.NoError(t, err)
	assert.Equal(t, requested, fromBlock)
}

func TestFilterName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "a - b:c:d", FilterName("a", "b", "c", "d"))
	assert.Equal(t, "empty args test", FilterName("empty args test"))
}

func benchmarkFilter(b *testing.B, nFilters, nAddresses, nEvents int) {
	lggr := logger.TestLogger(b)
	lp := NewLogPoller(nil, nil, lggr, 1*time.Hour, 2, 3, 2, 1000)
	for i := 0; i < nFilters; i++ {
		var addresses []common.Address
		var events []common.Hash
		for j := 0; j < nAddresses; j++ {
			addresses = append(addresses, common.BigToAddress(big.NewInt(int64(j+1))))
		}
		for j := 0; j < nEvents; j++ {
			events = append(events, common.BigToHash(big.NewInt(int64(j+1))))
		}
		err := lp.RegisterFilter(Filter{Name: "my filter", EventSigs: events, Addresses: addresses})
		require.NoError(b, err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lp.filter(nil, nil, nil)
	}
}

func BenchmarkFilter10_1(b *testing.B) {
	benchmarkFilter(b, 10, 1, 1)
}
func BenchmarkFilter100_10(b *testing.B) {
	benchmarkFilter(b, 100, 10, 10)
}
func BenchmarkFilter1000_100(b *testing.B) {
	benchmarkFilter(b, 1000, 100, 100)
}
