package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	lp := logpoller.NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 15*time.Second, 4, 3)
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
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`), lgs[0].Data)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[1].Data)

	// Test scenario: multiple logs per block for many blocks (also after reorg).
	// Chain gen <- 1 <- 2 (L1_1)
	//                \ 2'(L1_2) <- 3 <- 4 (L1_3, L2_4) <- 5 (L1_5)
	// DB: 1, 2, 3
	// - Should save 4, 5, 6 blocks
	// - Should obtain logs L1_3, L2_4, L1_5, L2_6
	_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(3)})
	require.NoError(t, err)
	_, err = emitter2.EmitLog1(id, []*big.Int{big.NewInt(4)})
	require.NoError(t, err)
	// Create 4
	ec.Commit()
	_, err = emitter2.EmitLog1(id, []*big.Int{big.NewInt(5)})
	require.NoError(t, err)
	// Create 5
	ec.Commit()

	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(6), newStart)
	lgs, err = orm.SelectLogsByBlockRange(4, 5)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000003`), lgs[0].Data)
	assert.Equal(t, emitterAddress1, lgs[0].Address)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000004`), lgs[1].Data)
	assert.Equal(t, emitterAddress2, lgs[1].Address)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000005`), lgs[2].Data)
	assert.Equal(t, emitterAddress2, lgs[2].Address)

	// Test scenario: node down for exactly finality + 1 block
	// Chain gen <- 1 <- 2 (L1_1)
	//                \ 2'(L1_2) <- 3 <- 4 (L1_3, L2_4) <- 5 (L1_5) <- 6 (L2_6) <- 7 (L1_7) <- 8 <- 9 <- 10 <- 11 (L1_8)
	// DB: 1, 2, 3, 4, 5
	// - We expect block 6 to backfilled
	// - Then block 7-11 to be handled block by block (treated as unfinalized).
	_, err = emitter2.EmitLog1(id, []*big.Int{big.NewInt(6)})
	require.NoError(t, err)
	ec.Commit()
	for i := 0; i < 5; i++ {
		if i == 0 {
			_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(7)})
			require.NoError(t, err)
		}
		if i == 4 {
			_, err = emitter1.EmitLog1(id, []*big.Int{big.NewInt(8)})
			require.NoError(t, err)
		}
		ec.Commit()
	}
	newStart = lp.PollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(12), newStart)
	lgs, err = orm.SelectLogsByBlockRange(6, 11)
	require.NoError(t, err)
	require.Equal(t, 3, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000006`), lgs[0].Data)
	assert.Equal(t, int64(6), lgs[0].BlockNumber)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000007`), lgs[1].Data)
	assert.Equal(t, int64(7), lgs[1].BlockNumber)
	_, err = orm.SelectBlockByNumber(6)
	require.True(t, errors.Is(err, sql.ErrNoRows)) // Do not expect to save backfilled blocks.
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

func TestCanonicalQuery(t *testing.T) {
	lggr := logger.TestLogger(t)
	_, db := heavyweight.FullTestDB(t, "logs", true, false)
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	topic1 := emitterABI.Events["Log1"].ID
	topic2 := emitterABI.Events["Log2"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// Block 1 and 2
	require.NoError(t, o.InsertLogs([]logpoller.Log{
		genLog(chainID, 1, 1, "0x1", topic1[:], address1),
		genLog(chainID, 2, 1, "0x1", topic2[:], address2),
		genLog(chainID, 1, 2, "0x2", topic1[:], address2),
		genLog(chainID, 2, 2, "0x2", topic2[:], address1),
	}))

	// Block 1' and 2' and 3
	require.NoError(t, o.InsertLogs([]logpoller.Log{
		genLog(chainID, 1, 1, "0x3", topic1[:], address1),
		genLog(chainID, 2, 1, "0x3", topic2[:], address2),
		genLog(chainID, 1, 2, "0x4", topic1[:], address2),
		genLog(chainID, 2, 2, "0x4", topic2[:], address1),
		genLog(chainID, 1, 3, "0x5", topic1[:], address1),
		genLog(chainID, 2, 3, "0x5", topic2[:], address2),
	}))

	lgs, err := o.SelectCanonicalLogsByBlockRange(1, 3)
	require.NoError(t, err)
	// We expect only logs from block hash 0x3 and 0x4 as they are more recent for the same block height.
	require.Equal(t, 6, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[1].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[2].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[3].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[4].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[5].BlockHash.String())

	// Filter by address and topic
	lgs, err = o.SelectCanonicalLogsByBlockRangeTopicAddress(1, 3, address1, [][]byte{topic1[:]})
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, address1, lgs[0].Address)
	assert.Equal(t, topic1.Bytes(), lgs[0].Topics[0])
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", lgs[1].BlockHash.String())
	assert.Equal(t, address1, lgs[1].Address)

	// Filter by block
	lgs, err = o.SelectCanonicalLogsByBlockRangeTopicAddress(2, 2, address2, [][]byte{topic1[:]})
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[0].BlockHash.String())
	assert.Equal(t, int64(1), lgs[0].LogIndex)
	assert.Equal(t, address2, lgs[0].Address)
	assert.Equal(t, topic1.Bytes(), lgs[0].Topics[0])

	// TODO: Benchmarking
}
