package logpoller

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestLogPoller(t *testing.T) {
	_, db := heavyweight.FullTestDB(t, "logs", true, false)
	chainID := big.NewInt(42)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	lggr := logger.TestLogger(t)

	// Set up a test chain with a log emitting contract deployed.
	orm := NewORM(chainID, db, lggr, pgtest.NewPGCfg(true))
	id := cltest.NewSimulatedBackendIdentity(t)
	ec := cltest.NewSimulatedBackend(t, map[common.Address]core.GenesisAccount{
		id.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	emitterAddress, _, emitter, err := log_emitter.DeployLogEmitter(id, ec)
	require.NoError(t, err)
	ec.Commit()

	// Set up a log poller listening for log emitter logs.
	lp := NewLogPoller(orm, client.NewSimulatedBackendClient(t, ec, chainID), lggr)
	emitterABI, err := abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
	require.NoError(t, err)
	lp.addresses = []common.Address{emitterAddress}
	lp.topics = [][]common.Hash{{emitterABI.Events["Log1"].ID, emitterABI.Events["Log2"].ID}}
	t.Log(emitter, lp)

	b, err := ec.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), b.NumberU64())

	// Chain genesis <- 1
	// DB: empty
	newStart := lp.pollAndSaveLogs(context.Background(), 1)
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
	newStart = lp.pollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(2), newStart)
	latest, err := orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(1), latest.BlockNumber)

	// Chain gen <- 1 <- 2 (L1)
	// DB: 1
	_, err = emitter.EmitLog1(id, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)
	ec.Commit()

	// Polling should get us the L1 hello log.
	newStart = lp.pollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(3), newStart)
	latest, err = orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(2), latest.BlockNumber)
	lgs, err = orm.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 1, len(lgs))
	assert.Equal(t, emitterAddress, lgs[0].Address)
	assert.Equal(t, latest.BlockHash, lgs[0].BlockHash)
	assert.Equal(t, hexutil.Encode(lgs[0].Topics[0]), emitterABI.Events["Log1"].ID.String())
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`),
		lgs[0].Data)

	// Chain gen <- 1 <- 2 (L1)
	//                \ 2'(L1') <- 3
	// DB: 1, 2
	// - Detect a reorg,
	// - Update the block 2's hash
	// - Save L1'
	lca, err := ec.BlockByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.NoError(t, ec.Fork(context.Background(), lca.Hash()))
	_, err = emitter.EmitLog1(id, []*big.Int{big.NewInt(2)})
	require.NoError(t, err)
	// Create 2'
	ec.Commit()
	// Create 3 (we need a new block for us to do any polling and detect the reorg).
	ec.Commit()

	newStart = lp.pollAndSaveLogs(context.Background(), newStart)
	assert.Equal(t, int64(4), newStart)
	latest, err = orm.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(3), latest.BlockNumber)
	lgs, err = orm.SelectLogsByBlockRange(1, 3)
	require.NoError(t, err)
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000001`), lgs[0].Data)
	assert.Equal(t, hexutil.MustDecode(`0x0000000000000000000000000000000000000000000000000000000000000002`), lgs[1].Data)
}
