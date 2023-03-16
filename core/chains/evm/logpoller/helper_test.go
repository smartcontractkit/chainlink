package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//	func (lp *logPoller) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
//		lp.PollAndSaveLogs(ctx, currentBlockNumber)
//		lastProcessed, _ := lp.orm.SelectLatestBlock()
//		return lastProcessed.BlockNumber + 1
//	}
//
// // Similar to lp.Start(), but works after it's already been run once
//
//	func (lp *logPoller) Restart(parentCtx context.Context) error {
//		lp.StartStopOnce = utils.StartStopOnce{}
//		lp.done = make(chan struct{})
//		return lp.Start(parentCtx)
//	}
//
//	func (lp *logPoller) Filter() ethereum.FilterQuery {
//		return lp.Filter(nil, nil, nil)
//	}
//
//	func (lp *logPoller) ConvertLogs(gethLogs []types.Log, blocks []LogPollerBlock) []Log {
//		return convertLogs(gethLogs, blocks, lp.lggr, lp.ec.ChainID())
//	}
//
//	func (lp *logPoller) BlocksFromLogs(ctx context.Context, logs []types.Log) (blocks []LogPollerBlock, err error) {
//		return lp.blocksFromLogs(ctx, logs)
//	}
var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

type TestHarness struct {
	Lggr                             logger.Logger
	ChainID                          *big.Int
	ORM                              *logpoller.ORM
	ORM2                             *logpoller.ORM // Dummy second chain
	LogPoller                        logpoller.LogPollerTest
	Client                           *backends.SimulatedBackend
	Owner                            *bind.TransactOpts
	Emitter1, Emitter2               *log_emitter.LogEmitter
	EmitterAddress1, EmitterAddress2 common.Address
	EthDB                            ethdb.Database
}

func SetupTH(t testing.TB, finalityDepth, backfillBatchSize, rpcBatchSize int64) TestHarness {
	lggr := logger.TestLogger(t)
	chainID := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_filters_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_logs_evm_chain_id_fkey DEFERRED`)))
	o := logpoller.NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	o2 := logpoller.NewORM(chainID2, db, lggr, pgtest.NewQConfig(true))
	owner := testutils.MustNewSimTransactor(t)
	ethDB := rawdb.NewMemoryDatabase()
	ec := backends.NewSimulatedBackendWithDatabase(ethDB, map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.
	lp := logpoller.NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 1*time.Hour, finalityDepth, backfillBatchSize, rpcBatchSize, 1000)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()
	return TestHarness{
		Lggr:            lggr,
		ChainID:         chainID,
		ORM:             o,
		ORM2:            o2,
		LogPoller:       lp,
		Client:          ec,
		Owner:           owner,
		Emitter1:        emitter1,
		Emitter2:        emitter2,
		EmitterAddress1: emitterAddress1,
		EmitterAddress2: emitterAddress2,
		EthDB:           ethDB,
	}
}

func (th *TestHarness) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
	th.LogPoller.PollAndSaveLogs(ctx, currentBlockNumber)
	latest, _ := th.LogPoller.LatestBlock()
	return latest + 1
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
		chainBlk, err := ec.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
		require.NoError(t, err)
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes(), "block %v", i)
	}
}
