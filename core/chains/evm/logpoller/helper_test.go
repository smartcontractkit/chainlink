package logpoller

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type TestHarness struct {
	Lggr                             logger.Logger
	ChainID                          *big.Int
	db                               *sqlx.DB
	ORM                              *ORM
	LogPoller                        *logPoller
	Client                           *backends.SimulatedBackend
	Owner                            *bind.TransactOpts
	Emitter1, Emitter2               *log_emitter.LogEmitter
	EmitterAddress1, EmitterAddress2 common.Address
	EthDB                            ethdb.Database
}

func SetupTH(t *testing.T, finalityDepth, backfillBatchSize, rpcBatchSize int64) TestHarness {
	lggr := logger.TestLogger(t)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_filters_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_logs_evm_chain_id_fkey DEFERRED`)))
	o := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	owner := testutils.MustNewSimTransactor(t)
	ethDB := rawdb.NewMemoryDatabase()
	ec := backends.NewSimulatedBackendWithDatabase(ethDB, map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.
	lp := NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID), lggr, 1*time.Hour, finalityDepth, backfillBatchSize, rpcBatchSize, 1000)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()
	return TestHarness{
		Lggr:            lggr,
		ChainID:         chainID,
		db:              db,
		ORM:             o,
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

// returns next unfinalized block number to be fetched and saved to db
func (lp *logPoller) GetCurrentBlock() int64 {
	lastProcessed, _ := lp.orm.SelectLatestBlock()
	return lastProcessed.BlockNumber + 1
}

func (lp *logPoller) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
	lp.pollAndSaveLogs(ctx, currentBlockNumber)
	return lp.GetCurrentBlock()
}

// Similar to lp.Start(), but works after it's already been run once
func (lp *logPoller) Restart(parentCtx context.Context) error {
	lp.StartStopOnce = utils.StartStopOnce{}
	lp.done = make(chan struct{})
	return lp.Start(parentCtx)
}

func (lp *logPoller) Filter() ethereum.FilterQuery {
	return lp.filter(nil, nil, nil)
}

func (o *ORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
	return o.selectLogsByBlockRange(start, end)
}
