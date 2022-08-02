package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

func GenLog(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address) logpoller.Log {
	return logpoller.Log{
		EvmChainId:  utils.NewBig(chainID),
		LogIndex:    logIndex,
		BlockHash:   common.HexToHash(blockHash),
		BlockNumber: blockNum,
		EventSig:    topic1,
		Topics:      [][]byte{topic1},
		Address:     address,
		TxHash:      common.HexToHash("0x1234"),
		Data:        append([]byte("hello "), byte(blockNum)),
	}
}

func assertDontHave(t *testing.T, start, end int, orm *logpoller.ORM) {
	for i := start; i < end; i++ {
		_, err := orm.SelectBlockByNumber(int64(i))
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	}
}

func assertHaveCanonical(t *testing.T, start, end int, ec *backends.SimulatedBackend, orm *logpoller.ORM) {
	for i := start; i < end; i++ {
		blk, err := orm.SelectBlockByNumber(int64(i))
		require.NoError(t, err, "block %v", i)
		chainBlk, err := ec.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes(), "block %v", i)
	}
}

func TestLogPoller_SynchronizedWithGeth(t *testing.T) {
	// The log poller's blocks table should remain synchronized
	// with the canonical chain of geth's despite arbitrary mixes of mining and reorgs.
	testParams := gopter.DefaultTestParameters()
	testParams.MinSuccessfulTests = 100
	p := gopter.NewProperties(testParams)
	numChainInserts := 3
	finalityDepth := 5
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	owner := testutils.MustNewSimTransactor(t)
	owner.GasPrice = big.NewInt(10e9)
	p.Property("synchronized with geth", prop.ForAll(func(mineOrReorg []uint64) bool {
		// After the set of reorgs, we should have the same canonical blocks that geth does.
		t.Log("Starting test", mineOrReorg)
		chainID := testutils.NewRandomEVMChainID()
		// Set up a test chain with a log emitting contract deployed.
		orm := logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true))
		// Note this property test is run concurrently and the sim is not threadsafe.
		ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
			owner.From: {
				Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
			},
		}, 10e6)
		_, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
		require.NoError(t, err)
		lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, int64(finalityDepth), 3)
		for i := 0; i < finalityDepth; i++ { // Have enough blocks that we could reorg the full finalityDepth-1.
			ec.Commit()
		}
		currentBlock := int64(1)
		currentBlock = lp.PollAndSaveLogs(context.Background(), currentBlock)
		matchesGeth := func() bool {
			// Check every block is identical
			latest, err := ec.BlockByNumber(context.Background(), nil)
			require.NoError(t, err)
			for i := 1; i < int(latest.NumberU64()); i++ {
				ourBlock, err := lp.BlockByNumber(int64(i))
				require.NoError(t, err)
				gethBlock, err := ec.BlockByNumber(context.Background(), big.NewInt(int64(i)))
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
					latest, err := ec.BlockByNumber(context.Background(), nil)
					require.NoError(t, err)
					t.Log("mined block", latest.Hash())
				}
			} else {
				// Reorg blocks
				// Get the hash of reorg block
				latest, err := ec.BlockByNumber(context.Background(), nil)
				require.NoError(t, err)
				reorgedBlock := big.NewInt(0).Sub(latest.Number(), big.NewInt(int64(mineOrReorg[i])))
				reorg, err := ec.BlockByNumber(context.Background(), reorgedBlock)
				require.NoError(t, err)
				require.NoError(t, ec.Fork(context.Background(), reorg.Hash()))
				t.Logf("Reorging from (%v, %x) back to (%v, %x)\n", latest.NumberU64(), latest.Hash(), reorgedBlock.Uint64(), reorg.Hash())
				// Actually need to change the block here to trigger the reorg.
				_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(1)})
				require.NoError(t, err)
				for j := 0; j < int(mineOrReorg[i]+1); j++ { // Need +1 to make it actually longer height so we detect it.
					ec.Commit()
				}
				latest, err = ec.BlockByNumber(context.Background(), nil)
				require.NoError(t, err)
				t.Logf("New latest (%v, %x), latest parent %x)\n", latest.NumberU64(), latest.Hash(), latest.ParentHash())
			}
			currentBlock = lp.PollAndSaveLogs(context.Background(), currentBlock)
		}
		return matchesGeth()
	}, gen.SliceOfN(numChainInserts, gen.UInt64Range(1, uint64(finalityDepth-1))))) // Max reorg depth is finality depth - 1
	p.TestingRun(t)
}

func TestLogPoller_PollAndSaveLogs(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	chainID := testutils.NewRandomEVMChainID()
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))

	// Set up a test chain with a log emitting contract deployed.
	orm := logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true))
	owner := testutils.MustNewSimTransactor(t)
	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()

	// Set up a log poller listening for log emitter logs.
	lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, 2, 3)
	lp.MergeFilter([]common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{emitterAddress1, emitterAddress2})

	b, err := ec.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), b.NumberU64())

	// Test scenario: single block in chain, no logs.
	// Chain genesis <- 1
	// DB: empty
	newStart := lp.PollAndSaveLogs(context.Background(), 1)
	assert.Equal(t, int64(2), newStart)

	// We expect to have saved block 1.
	lpb, err := orm.SelectBlockByNumber(1)
	require.NoError(t, err)
	assert.Equal(t, lpb.BlockHash, b.Hash())
	assert.Equal(t, lpb.BlockNumber, int64(b.NumberU64()))
	assert.Equal(t, int64(1), int64(b.NumberU64()))
	// No logs.
	lgs, err := orm.SelectLogsByBlockRange(1, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(lgs))
	assertHaveCanonical(t, 1, 1, ec, orm)

	// Polling again should be a noop, since we are at the latest.
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(2), newStart)
	latest, err := orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(1), latest.BlockNumber)
	assertHaveCanonical(t, 1, 1, ec, orm)

	// Test scenario: one log 2 block chain.
	// Chain gen <- 1 <- 2 (L1)
	// DB: 1
	_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)
	ec.Commit()

	// Polling should get us the L1 log.
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(3), newStart)
	latest, err = orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(2), latest.BlockNumber)
	lgs, err = orm.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, emitterAddress1, lgs[0].Address)
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
	reorgedOutBlock, err := ec.BlockByNumber(context.Background(), big.NewInt(2))
	require.NoError(t, err)
	lca, err := ec.BlockByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.NoError(t, ec.Fork(context.Background(), lca.Hash()))
	_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(2)})
	require.NoError(t, err)
	// Create 2'
	ec.Commit()
	// Create 3 (we need a new block for us to do any polling and detect the reorg).
	ec.Commit()

	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(4), newStart)
	latest, err = orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(3), latest.BlockNumber)
	lgs, err = orm.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[0].Data)
	assertHaveCanonical(t, 1, 3, ec, orm)

	// Test scenario: reorg back to previous tip.
	// Chain gen <- 1 <- 2 (L1_1) <- 3' (L1_3) <- 4
	//                \ 2'(L1_2) <- 3
	require.NoError(t, ec.Fork(context.Background(), reorgedOutBlock.Hash()))
	_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(3)})
	require.NoError(t, err)
	// Create 3'
	ec.Commit()
	// Create 4
	ec.Commit()
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(5), newStart)
	latest, err = orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(4), latest.BlockNumber)
	lgs, err = orm.SelectLogsByBlockRange(1, 3)
	// We expect ONLY L1_1 and L1_3 since L1_2 is reorg'd out.
	assert.Equal(t, 2, len(lgs))
	assert.Equal(t, int64(2), lgs[0].BlockNumber)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`), lgs[0].Data)
	assert.Equal(t, int64(3), lgs[1].BlockNumber)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000003`), lgs[1].Data)
	assertHaveCanonical(t, 1, 1, ec, orm)
	assertHaveCanonical(t, 3, 4, ec, orm)
	assertDontHave(t, 2, 2, orm) // 2 gets backfilled

	// Test scenario: multiple logs per block for many blocks (also after reorg).
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6)
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2', 3'
	// - Should save 4, 5, 6 blocks
	// - Should obtain logs L1_3, L2_5, L1_6
	_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(4)})
	require.NoError(t, err)
	_, err = emitter2.EmitLog1(owner, []*big.Int{big.NewInt(5)})
	require.NoError(t, err)
	// Create 4
	ec.Commit()
	_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(6)})
	require.NoError(t, err)
	// Create 5
	ec.Commit()

	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(7), newStart)
	lgs, err = orm.SelectLogsByBlockRange(4, 6)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000004`), lgs[0].Data)
	assert.Equal(t, emitterAddress1, lgs[0].Address)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000005`), lgs[1].Data)
	assert.Equal(t, emitterAddress2, lgs[1].Address)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000006`), lgs[2].Data)
	assert.Equal(t, emitterAddress1, lgs[2].Address)
	assertHaveCanonical(t, 1, 1, ec, orm)
	assertDontHave(t, 2, 2, orm) // 2 gets backfilled
	assertHaveCanonical(t, 3, 6, ec, orm)

	// Test scenario: node down for exactly finality + 2 blocks
	// Note we only backfill up to finalized - 1 blocks, because we need to save the
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9) <- 10 (L1_10)
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2, 3, 4, 5, 6
	// - We expect block 7 to backfilled (treated as finalized)
	// - Then block 8-10 to be handled block by block (treated as unfinalized).
	for i := 7; i < 11; i++ {
		_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		ec.Commit()
	}
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(11), newStart)
	lgs, err = orm.SelectLogsByBlockRange(7, 9)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000007`), lgs[0].Data)
	assert.Equal(t, int64(7), lgs[0].BlockNumber)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000008`), lgs[1].Data)
	assert.Equal(t, int64(8), lgs[1].BlockNumber)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000009`), lgs[2].Data)
	assert.Equal(t, int64(9), lgs[2].BlockNumber)
	assertDontHave(t, 7, 7, orm) // Do not expect to save backfilled blocks.
	assertHaveCanonical(t, 8, 10, ec, orm)

	// Test scenario large backfill (multiple batches)
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9) <- 10..16
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2, 3, 4, 5, 6, (backfilled 7), 8, 9, 10
	// - 11, 12, 13 backfilled in batch 1
	// - 14 backfilled in batch 2
	// - 15, 16, 17 to be treated as unfinalized
	for i := 11; i < 18; i++ {
		_, err = emitter1.EmitLog1(owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		ec.Commit()
	}
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(18), newStart)
	lgs, err = orm.SelectLogsByBlockRange(11, 17)
	require.NoError(t, err)
	assert.Equal(t, 7, len(lgs))
	assertHaveCanonical(t, 15, 16, ec, orm)
	assertDontHave(t, 11, 14, orm) // Do not expect to save backfilled blocks.
}

func TestLogPoller_Logs(t *testing.T) {
	lggr := logger.TestLogger(t)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	o := logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true))
	event1 := EmitterABI.Events["Log1"].ID
	event2 := EmitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1-3
	require.NoError(t, o.InsertLogs([]logpoller.Log{
		GenLog(chainID, 1, 1, "0x3", event1[:], address1),
		GenLog(chainID, 2, 1, "0x3", event2[:], address2),
		GenLog(chainID, 1, 2, "0x4", event1[:], address2),
		GenLog(chainID, 2, 2, "0x4", event2[:], address1),
		GenLog(chainID, 1, 3, "0x5", event1[:], address1),
		GenLog(chainID, 2, 3, "0x5", event2[:], address2),
	}))

	// Select for all addresses
	lgs, err := o.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 6, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[1].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[2].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[3].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[4].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[5].BlockHash.String())

	// Filter by address and topic
	lgs, err = o.SelectLogsByBlockRangeFilter(1, 3, address1, event1[:])
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	// Filter by block
	lgs, err = o.SelectLogsByBlockRangeFilter(2, 2, address2, event1[:])
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, event1.Bytes(), lgs[0].Topics[0])
}

func TestLogPoller_MergeFilter(t *testing.T) {
	lp := logpoller.NewLogPoller(nil, nil, nil, 15*time.Second, 1, 1)
	a1 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	a2 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")
	lp.MergeFilter([]common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1})
	assert.Equal(t, []common.Address{a1}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, lp.Filter().Topics)

	// Should de-dupe eventSigs
	lp.MergeFilter([]common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2})
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter().Topics)

	// Should de-dupe addresses
	lp.MergeFilter([]common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2})
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter().Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter().Topics)
}
