package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

type TestHarness struct {
	Lggr logger.Logger
	// Chain2/ORM2 is just a dummy second chain, doesn't have a client.
	ChainID, ChainID2                *big.Int
	ORM, ORM2                        logpoller.ORM
	LogPoller                        logpoller.LogPollerTest
	Client                           *backends.SimulatedBackend
	Owner                            *bind.TransactOpts
	Emitter1, Emitter2               *log_emitter.LogEmitter
	EmitterAddress1, EmitterAddress2 common.Address
	EthDB                            ethdb.Database
}

func SetupTH(t testing.TB, opts logpoller.Opts) TestHarness {
	lggr := logger.Test(t)
	chainID := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	o := logpoller.NewORM(chainID, db, lggr)
	o2 := logpoller.NewORM(chainID2, db, lggr)
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
	// Mark genesis block as finalized to avoid any nulls in the tests
	head := esc.Backend().Blockchain().CurrentHeader()
	esc.Backend().Blockchain().SetFinalized(head)

	headTracker := headtracker.NewSimulatedHeadTracker(esc, opts.UseFinalityTag, opts.FinalityDepth)
	if opts.PollPeriod == 0 {
		opts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(o, esc, lggr, headTracker, opts)
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
	latest, _ := th.LogPoller.LatestBlock(ctx)
	return latest.BlockNumber + 1
}

func (th *TestHarness) assertDontHave(t *testing.T, start, end int) {
	for i := start; i < end; i++ {
		_, err := th.ORM.SelectBlockByNumber(testutils.Context(t), int64(i))
		assert.True(t, pkgerrors.Is(err, sql.ErrNoRows))
	}
}

func (th *TestHarness) assertHaveCanonical(t *testing.T, start, end int) {
	for i := start; i < end; i++ {
		blk, err := th.ORM.SelectBlockByNumber(testutils.Context(t), int64(i))
		require.NoError(t, err, "block %v", i)
		chainBlk, err := th.Client.BlockByNumber(testutils.Context(t), big.NewInt(int64(i)))
		require.NoError(t, err)
		assert.Equal(t, chainBlk.Hash().Bytes(), blk.BlockHash.Bytes(), "block %v", i)
	}
}
