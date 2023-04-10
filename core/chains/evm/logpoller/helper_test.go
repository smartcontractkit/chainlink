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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

type TestHarness struct {
	Lggr logger.Logger
	// Chain2/ORM2 is just a dummy second chain, doesn't have a client.
	ChainID, ChainID2                *big.Int
	ORM, ORM2                        *logpoller.ORM
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
	esc := client.NewSimulatedBackendClient(t, ec, chainID)
	lp := logpoller.NewLogPoller(o, esc, lggr, 1*time.Hour, finalityDepth, backfillBatchSize, rpcBatchSize, 1000)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()
	return TestHarness{
		Lggr:            lggr,
		ChainID:         chainID,
		ChainID2:        chainID2,
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

func (th *TestHarness) assertDontHave(t *testing.T, start, end int) {
	for i := start; i < end; i++ {
		_, err := th.ORM.SelectBlockByNumber(int64(i))
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	}
}

func (th *TestHarness) assertHaveCanonical(t *testing.T, start, end int) {
	for i := start; i < end; i++ {
		blk, err := th.ORM.SelectBlockByNumber(int64(i))
		require.NoError(t, err, "block %v", i)
		chainBlk, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
		require.NoError(t, err)
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes(), "block %v", i)
	}
}
