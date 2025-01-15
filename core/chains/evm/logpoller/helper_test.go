package logpoller_test

import (
	"context"
	"database/sql"
	"math/big"
	"strings"
	"testing"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

type Backend struct {
	*simulated.Backend
	t             testing.TB
	expectPending bool
}

// SetExpectPending sets whether the backend should expect txes to be pending
// after a Fork. We do this to avoid breaking the existing evmtypes.Backend interface (by
// for example passing in a pending bool to Fork).
func (b *Backend) SetExpectPending(pending bool) {
	b.expectPending = pending
}

// Fork as an override exists to maintain the same behaviour as the old
// simulated backend. Description of the changed behaviour
// here https://github.com/ethereum/go-ethereum/pull/30465#issuecomment-2362967508
// Basically the new simulated backend (post 1.14) will automatically
// put forked txes back in the mempool whereas the old one didn't
// so they would just remain on the fork.
func (b *Backend) Fork(parentHash common.Hash) error {
	if err := b.Backend.Fork(parentHash); err != nil {
		return err
	}
	// TODO: Fairly sure we need to upstream a tx pool sync like this:
	// func (c *SimulatedBeacon) Rollback() {
	//	// Flush all transactions from the transaction pools
	//	+       c.eth.TxPool().Sync()
	//	maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(common.Big1, 256), common.Big1)
	// Otherwise its possible the fork adds the txes to the pool
	// _after_ we Rollback so the rollback is ineffective.
	// In the meantime we can just wait for the txes to be pending as workaround.
	require.Eventually(b.t, func() bool {
		p, err := b.Backend.Client().PendingTransactionCount(context.Background())
		if err != nil {
			return false
		}
		b.t.Logf("waiting for forked txes to be pending, have %v, want %v\n", p, b.expectPending)
		return p > 0 == b.expectPending
	}, testutils.DefaultWaitTimeout, 500*time.Millisecond)
	b.Rollback()
	return nil
}

type TestHarness struct {
	Lggr logger.Logger
	// Chain2/ORM2 is just a dummy second chain, doesn't have a client.
	ChainID, ChainID2                *big.Int
	ORM, ORM2                        logpoller.ORM
	LogPoller                        logpoller.LogPollerTest
	Client                           *client.SimulatedBackendClient
	Backend                          evmtypes.Backend
	Owner                            *bind.TransactOpts
	Emitter1, Emitter2               *log_emitter.LogEmitter
	EmitterAddress1, EmitterAddress2 common.Address
}

func SetupTH(t testing.TB, opts logpoller.Opts) TestHarness {
	lggr := logger.Test(t)
	chainID := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	o := logpoller.NewORM(chainID, db, lggr)
	o2 := logpoller.NewORM(chainID2, db, lggr)
	owner := testutils.MustNewSimTransactor(t)
	// Needed for the new sim if you are using Rollback
	owner.GasTipCap = big.NewInt(1000000000)

	backend := simulated.NewBackend(types.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))

	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.

	esc := client.NewSimulatedBackendClient(t, backend, chainID)

	headTracker := headtracker.NewSimulatedHeadTracker(esc, opts.UseFinalityTag, opts.FinalityDepth)
	if opts.PollPeriod == 0 {
		opts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(o, esc, lggr, headTracker, opts)

	pendingNonce, err := backend.Client().PendingNonceAt(testutils.Context(t), owner.From)
	require.NoError(t, err)

	owner.Nonce = big.NewInt(0).SetUint64(pendingNonce)
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, backend.Client())
	require.NoError(t, err)

	owner.Nonce.Add(owner.Nonce, big.NewInt(1)) // Avoid race where DeployLogEmitter returns before PendingNonce has been incremented
	emitterAddress2, _, emitter2, err := log_emitter.DeployLogEmitter(owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	owner.Nonce = nil // Just use pending nonce after this

	return TestHarness{
		Lggr:            lggr,
		ChainID:         chainID,
		ChainID2:        chainID2,
		ORM:             o,
		ORM2:            o2,
		LogPoller:       lp,
		Client:          esc,
		Backend:         &Backend{t: t, Backend: backend, expectPending: true},
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
func (th *TestHarness) SetActiveClient(backend evmtypes.Backend, chainType chaintype.ChainType) {
	th.Backend = backend
	th.Client.SetBackend(backend, chainType)
}

func (th *TestHarness) finalizeThroughBlock(t *testing.T, blockNumber int64) {
	client.FinalizeThroughBlock(t, th.Backend, th.Client, blockNumber)
}
