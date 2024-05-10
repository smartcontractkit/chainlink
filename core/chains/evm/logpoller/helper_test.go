package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	pkgerrors "github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
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
	Client                           *client.SimulatedBackendClient
	Backend                          *simulated.Backend
	Owner                            *bind.TransactOpts
	Emitter1, Emitter2               *log_emitter.LogEmitter
	EmitterAddress1, EmitterAddress2 common.Address
}

func SetupTH(t testing.TB, opts logpoller.Opts) TestHarness {
	lggr := logger.Test(t)
	chainID := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	//dataDir, err := os.MkdirTemp("", "simgethdata")
	//require.NoError(t, err)

	o := logpoller.NewORM(chainID, db, lggr)
	o2 := logpoller.NewORM(chainID2, db, lggr)
	owner := testutils.MustNewSimTransactor(t)

	backend := simulated.NewBackend(types.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))

	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.

	esc := client.NewSimulatedBackendClient(t, backend, chainID)

	if opts.PollPeriod == 0 {
		opts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(o, esc, lggr, opts)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, backend.Client())
	require.NoError(t, err)
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	return TestHarness{
		Lggr:            lggr,
		ChainID:         chainID,
		ChainID2:        chainID2,
		ORM:             o,
		ORM2:            o2,
		LogPoller:       lp,
		Client:          esc,
		Backend:         backend,
		Owner:           owner,
		Emitter1:        emitter1,
		Emitter2:        emitter2,
		EmitterAddress1: emitterAddress1,
		EmitterAddress2: emitterAddress2,
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

// Simulates an RPC failover event to an alternate rpc server. This can also be used to
// simulate switching back to the primary rpc after it recovers.
func (th *TestHarness) SetActiveClient(backend *simulated.Backend, optimismMode bool) {
	th.Backend = backend
	th.Client.SetBackend(backend, optimismMode)
}
