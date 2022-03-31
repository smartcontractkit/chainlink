package logpoller_test

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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	emitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

func assertDontHave(t *testing.T, start, end int, orm *logpoller.ORM) {
	for i := start; i < end; i++ {
		_, err := orm.SelectBlockByNumber(int64(i))
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	}
}

func assertHaveCanonical(t *testing.T, start, end int, ec *backends.SimulatedBackend, orm *logpoller.ORM) {
	for i := start; i < end; i++ {
		blk, err := orm.SelectBlockByNumber(int64(i))
		require.NoError(t, err)
		chainBlk, err := ec.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes())
	}
}

func TestLogPoller(t *testing.T) {
	_, db := heavyweight.FullTestDB(t, "logs", true, false)
	chainID := big.NewInt(42)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	lggr := logger.TestLogger(t)

	// Set up a test chain with a log emitting contract deployed.
	orm := logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true))
	id := cltest.NewSimulatedBackendIdentity(t)
	ec := cltest.NewSimulatedBackend(t, map[common.Address]core.GenesisAccount{
		id.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(id, ec)
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(id, ec)
	require.NoError(t, err)
	ec.Commit()

	// Set up a log poller listening for log emitter logs.
	lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, 2, 3)
	lp.MergeFilter([][]common.Hash{{emitterABI.Events["Log1"].ID, emitterABI.Events["Log2"].ID}}, []common.Address{emitterAddress1})
	lp.MergeFilter([][]common.Hash{{emitterABI.Events["Log1"].ID, emitterABI.Events["Log2"].ID}}, []common.Address{emitterAddress2})

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

	// Polling again should be a noop, since we are at the latest.
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(2), newStart)
	latest, err := orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(1), latest.BlockNumber)

	// Test scenario: one log 2 block chain.
	// Chain gen <- 1 <- 2 (L1)
	// DB: 1
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(1)})
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
	assert.Equal(t, hexutil.Encode(lgs[0].Topics[0]), emitterABI.Events["Log1"].ID.String())
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
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(2)})
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
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(3)})
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
	assertHaveCanonical(t, 1, 4, ec, orm)

	// Test scenario: multiple logs per block for many blocks (also after reorg).
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6)
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2', 3'
	// - Should save 4, 5, 6 blocks
	// - Should obtain logs L1_3, L2_5, L1_6
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(4)})
	require.NoError(t, err)
	_, err = emitter2.EmitLog1(id, []*big.Int{big.NewInt(5)})
	require.NoError(t, err)
	// Create 4
	ec.Commit()
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(6)})
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
	assertHaveCanonical(t, 1, 6, ec, orm)

	// Test scenario: node down for exactly finality + 1 block
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9)
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2, 3, 4, 5, 6
	// - We expect block 7 to backfilled (treated as finalized)
	// - Then block 8-9 to be handled block by block (treated as unfinalized).
	for i := 7; i < 10; i++ {
		_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		ec.Commit()
	}
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(10), newStart)
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
	assertHaveCanonical(t, 8, 9, ec, orm)

	// Test scenario large backfill (multiple batches)
	// Chain gen <- 1 <- 2 (L1_1) <- 3' L1_3 <- 4 <- 5 (L1_4, L2_5) <- 6 (L1_6) <- 7 (L1_7) <- 8 (L1_8) <- 9 (L1_9) <- 10..15
	//                \ 2'(L1_2) <- 3
	// DB: 1, 2, 3, 4, 5, 6, (backfilled 7), 8, 9
	// - 10, 11, 12 backfilled in batch 1
	// - 13 backfilled in batch 2
	// - 14, 15 to be treated as unfinalized
	for i := 10; i < 16; i++ {
		_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err)
		ec.Commit()
	}
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(16), newStart)
	lgs, err = orm.SelectLogsByBlockRange(10, 15)
	require.NoError(t, err)
	assert.Equal(t, 6, len(lgs))
	assertHaveCanonical(t, 14, 15, ec, orm)
	assertDontHave(t, 10, 13, orm) // Do not expect to save backfilled blocks.
}

func genLog(chainID *big.Int, logIndex int64, blockNum int64, blockHash string, topic1 []byte, address common.Address) logpoller.Log {
	return logpoller.Log{
		EvmChainId:  utils.NewBig(chainID),
		LogIndex:    logIndex,
		BlockHash:   common.HexToHash(blockHash),
		BlockNumber: blockNum,
		Topics:      [][]byte{topic1},
		Address:     address,
		TxHash:      common.HexToHash("0x1234"),
		Data:        []byte("hello"),
	}
}

func TestLogsQuerying(t *testing.T) {
	lggr := logger.TestLogger(t)
	chainID := big.NewInt(137)
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	topic1 := emitterABI.Events["Log1"].ID
	topic2 := emitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1-3
	require.NoError(t, o.InsertLogs([]logpoller.Log{
		genLog(chainID, 1, 1, "0x3", topic1[:], address1),
		genLog(chainID, 2, 1, "0x3", topic2[:], address2),
		genLog(chainID, 1, 2, "0x4", topic1[:], address2),
		genLog(chainID, 2, 2, "0x4", topic2[:], address1),
		genLog(chainID, 1, 3, "0x5", topic1[:], address1),
		genLog(chainID, 2, 3, "0x5", topic2[:], address2),
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
	lgs, err = o.SelectLogsByBlockRangeTopicAddress(1, 3, address1, [][]byte{topic1[:]})
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, topic1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	// Filter by block
	lgs, err = o.SelectLogsByBlockRangeTopicAddress(2, 2, address2, [][]byte{topic1[:]})
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, topic1.Bytes(), lgs[0].Topics[0])
}

func TestPopulateLoadedDB(t *testing.T) {
	t.Skip()
	lggr := logger.TestLogger(t)
	_, db := heavyweight.FullTestDB(t, "logs_scale", true, false)
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	topic1 := emitterABI.Events["Log1"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	var logs []logpoller.Log
	for i := 0; i < 1000; i++ {
		logs = append(logs, genLog(chainID, 1, int64(i), fmt.Sprintf("0x%d", i), topic1[:], address1))
	}
	require.NoError(t, o.InsertLogs(logs))
	var logs2 []logpoller.Log
	for i := 1001; i < 2000; i++ {
		logs2 = append(logs2, genLog(chainID, 1, int64(i), fmt.Sprintf("0x%d", i), topic1[:], address2))
	}
	require.NoError(t, o.InsertLogs(logs2))
}
