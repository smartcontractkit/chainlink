package v2

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	EmitterABI, _    = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
	VRFEmitterABI, _ = abi.JSON(strings.NewReader(vrf_log_emitter.VRFLogEmitterABI))
)

type TestHarness struct {
	Lggr              logger.Logger
	ChainID           *big.Int
	ORM               *logpoller.DbORM
	LogPoller         logpoller.LogPollerTest
	Client            *backends.SimulatedBackend
	Emitter           *log_emitter.LogEmitter
	EmitterAddress    common.Address
	VRFLogEmitter     *vrf_log_emitter.VRFLogEmitter
	VRFEmitterAddress common.Address
	Owner             *bind.TransactOpts
	EthDB             ethdb.Database
	Db                *sqlx.DB
	Listener          *listenerV2
	Ctx               context.Context
}

func SetupTH(t *testing.T,
	useFinalityTag bool,
	finalityDepth, backfillBatchSize,
	rpcBatchSize, keepFinalizedBlocksDepth int64,
	mockChainUpdateFn func(*evmmocks.Chain, *TestHarness)) *TestHarness {

	lggr := logger.TestLogger(t)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	o := logpoller.NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	owner := testutils.MustNewSimTransactor(t)
	ethDB := rawdb.NewMemoryDatabase()
	ec := backends.NewSimulatedBackendWithDatabase(ethDB, map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	// VRF Listener relies on block timestamps, but SimulatedBackend uses by default clock starting from 1970-01-01
	// This trick is used to move the clock closer to the current time. We set first block to be X hours ago.
	// FirstBlockAge is used to compute first block's timestamp in SimulatedBackend (time.Now() - FirstBlockAge)
	const FirstBlockAge = 24 * time.Hour
	blockTime := time.UnixMilli(int64(ec.Blockchain().CurrentHeader().Time))
	err := ec.AdjustTime(time.Since(blockTime) - FirstBlockAge)
	require.NoError(t, err)
	ec.Commit()

	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.
	esc := client.NewSimulatedBackendClient(t, ec, chainID)
	// Mark genesis block as finalized to avoid any nulls in the tests
	head := esc.Backend().Blockchain().CurrentHeader()
	esc.Backend().Blockchain().SetFinalized(head)

	lp := logpoller.NewLogPoller(o, esc, lggr, 1*time.Hour, useFinalityTag, finalityDepth, backfillBatchSize, rpcBatchSize, keepFinalizedBlocksDepth)

	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	vrfLogEmitterAddress, _, vrfLogEmitter, err := vrf_log_emitter.DeployVRFLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()

	// Log Poller Listener
	cfg := pgtest.NewQConfig(false)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr, cfg)
	require.NoError(t, ks.Unlock("blah"))
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
		EVMChainID:          chainID.String(),
	}).Toml())
	require.NoError(t, err)

	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(vrfLogEmitter.Address(), ec)
	require.Nil(t, err)
	coordinator := NewCoordinatorV2(coordinatorV2)

	chain := evmmocks.NewChain(t)
	listener := &listenerV2{
		respCount:   map[string]uint64{},
		job:         j,
		chain:       chain,
		l:           logger.Sugared(lggr),
		coordinator: coordinator,
	}
	ctx := testutils.Context(t)

	// Filter registration is idempotent, so we can just call it every time
	// and retry on errors using the ticker.
	err = lp.RegisterFilter(logpoller.Filter{
		Name: fmt.Sprintf("vrf_%s_keyhash_%s_job_%d", "v2", listener.job.VRFSpec.PublicKey.MustHash().String(), listener.job.ID),
		EventSigs: evmtypes.HashArray{
			vrf_log_emitter.VRFLogEmitterRandomWordsRequested{}.Topic(),
			vrf_log_emitter.VRFLogEmitterRandomWordsFulfilled{}.Topic(),
		},
		// EventSigs: []common.Hash{VRFEmitterABI.Events["RandomWordsRequested"].ID},
		Addresses: evmtypes.AddressArray{
			vrfLogEmitter.Address(),
			// listener.job.VRFSpec.CoordinatorAddress.Address(),
		},
	})
	require.Nil(t, err)
	require.NoError(t, lp.RegisterFilter(logpoller.Filter{
		Name:      "Integration test",
		EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID},
		Addresses: []common.Address{emitterAddress1},
		Retention: 0}))
	require.Nil(t, err)
	require.Len(t, lp.Filter(nil, nil, nil).Addresses, 2)
	require.Len(t, lp.Filter(nil, nil, nil).Topics, 1)

	th := &TestHarness{
		Lggr:              lggr,
		ChainID:           chainID,
		ORM:               o,
		LogPoller:         lp,
		Emitter:           emitter1,
		EmitterAddress:    emitterAddress1,
		VRFLogEmitter:     vrfLogEmitter,
		VRFEmitterAddress: vrfLogEmitterAddress,
		Client:            ec,
		Owner:             owner,
		EthDB:             ethDB,
		Db:                db,
		Listener:          listener,
		Ctx:               ctx,
	}
	mockChainUpdateFn(chain, th)
	return th
}

func TestInitProcessedBlock_NoVRFReqs(t *testing.T) {
	t.Parallel()

	finalityDepth := int64(3)
	th := SetupTH(t, false, finalityDepth, 3, 2, 1000, func(mockChain *evmmocks.Chain, th *TestHarness) {
		// txstore := txmgr.NewTxStore(db, lggr, cfg)
		// txm := makeTestTxm(t, txstore, ks)
		// chain.On("TxManager").Return(txm)
		mockChain.On("ID").Return(th.ChainID)
		mockChain.On("LogPoller").Return(th.LogPoller)
	})

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(finalityDepth); i++ {
		th.Client.Commit()
	}

	// Emit some logs in blocks finalityDepth+1 to finalityDepth+4.
	numBlocks := 5
	for i := 0; i < numBlocks; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(testutils.Context(t)))

	// The poller starts on a new chain at latest-finality (finalityDepth + 5 in this case),
	// Replaying from block 4 should guarantee we have block 4 immediately.  (We will also get
	// block 3 once the backup poller runs, since it always starts 100 blocks behind.)
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 4))

	// We should immediately have at least logs 4-7
	logs, err := th.LogPoller.Logs(4, 7, EmitterABI.Events["Log1"].ID, th.EmitterAddress,
		pg.WithParentCtx(testutils.Context(t)))
	require.NoError(t, err)
	require.Equal(t, 3, len(logs))

	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(th.Ctx)
	require.Nil(t, err)
	require.Equal(t, int64(6), lastProcessedBlock)
}

func TestInitProcessedBlock_NoUnfulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	finalityDepth := int64(3)
	th := SetupTH(t, false, finalityDepth, 3, 2, 1000, func(mockChain *evmmocks.Chain, curTH *TestHarness) {
		// txstore := txmgr.NewTxStore(db, lggr, cfg)
		// txm := makeTestTxm(t, txstore, ks)
		// chain.On("TxManager").Return(txm)
		mockChain.On("ID").Return(curTH.ChainID)
		mockChain.On("LogPoller").Return(curTH.LogPoller)
	})

	// Block 2 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(finalityDepth); i++ {
		th.Client.Commit()
	}

	// Make VRF requests and fulfill them
	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
		big.NewInt(1), big.NewInt(105), 1, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Client.Commit()
	_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, big.NewInt(1), big.NewInt(105), big.NewInt(10), true)
	require.NoError(t, err2)

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	numBlocks := 5
	th.Client.Commit()
	for i := 0; i < numBlocks; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(th.Ctx))

	latestBlock := finalityDepth + int64(numBlocks) + 1
	require.NoError(t, th.LogPoller.Replay(th.Ctx, latestBlock))

	// initializeLastProcessedBlock must return the finalizedBlockNumber instead of
	// unfulfilled VRF request block number, since all VRF requests are fulfilled
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(th.Ctx)
	require.Nil(t, err)
	require.Equal(t, int64(8), lastProcessedBlock)
}

func TestInitProcessedBlock_OneUnfulfilledVRFReq(t *testing.T) {
	t.Parallel()

	finalityDepth := int64(3)
	th := SetupTH(t, false, finalityDepth, 3, 2, 1000, func(mockChain *evmmocks.Chain, curTH *TestHarness) {
		// txstore := txmgr.NewTxStore(db, lggr, cfg)
		// txm := makeTestTxm(t, txstore, ks)
		// chain.On("TxManager").Return(txm)
		mockChain.On("ID").Return(curTH.ChainID)
		mockChain.On("LogPoller").Return(curTH.LogPoller)
	})

	// Block 2 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(finalityDepth); i++ {
		th.Client.Commit()
	}

	// Make VRF requests and fulfill them
	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
		big.NewInt(1), big.NewInt(105), 1, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Client.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	numBlocks := 5
	th.Client.Commit()
	for i := 0; i < numBlocks; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(th.Ctx))

	latestBlock := finalityDepth + int64(numBlocks) + 1
	require.NoError(t, th.LogPoller.Replay(th.Ctx, latestBlock))

	// initializeLastProcessedBlock must return the unfulfilled VRF request block number instead of
	// finalizedBlockNumber
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(th.Ctx)
	require.Nil(t, err)
	require.Equal(t, int64(5), lastProcessedBlock)
}

func TestInitProcessedBlock_SomeUnfulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	finalityDepth := int64(3)
	th := SetupTH(t, false, finalityDepth, 3, 2, 1000, func(mockChain *evmmocks.Chain, curTH *TestHarness) {
		// txstore := txmgr.NewTxStore(db, lggr, cfg)
		// txm := makeTestTxm(t, txstore, ks)
		// chain.On("TxManager").Return(txm)
		mockChain.On("ID").Return(curTH.ChainID)
		mockChain.On("LogPoller").Return(curTH.LogPoller)
	})

	// Block 2 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(finalityDepth); i++ {
		th.Client.Commit()
	}

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	numBlocks := 5
	th.Client.Commit()
	for i := 0; i < numBlocks; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
		// Make VRF requests and fulfill them
		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
			big.NewInt(int64(2*i)), big.NewInt(105), 1, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Client.Commit()
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
			big.NewInt(int64(2*i+1)), big.NewInt(106), 1, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Client.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(th.Ctx))

	latestBlock := finalityDepth + int64(numBlocks) + 1
	require.NoError(t, th.LogPoller.Replay(th.Ctx, latestBlock))

	// initializeLastProcessedBlock must return the earliest unfulfilled VRF request block
	// number instead of finalizedBlockNumber
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(th.Ctx)
	require.Nil(t, err)
	require.Equal(t, int64(7), lastProcessedBlock)
}

func TestInitProcessedBlock_UnfulfilledNFulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	finalityDepth := int64(3)
	th := SetupTH(t, false, finalityDepth, 3, 2, 1000, func(mockChain *evmmocks.Chain, curTH *TestHarness) {
		// txstore := txmgr.NewTxStore(db, lggr, cfg)
		// txm := makeTestTxm(t, txstore, ks)
		// chain.On("TxManager").Return(txm)
		mockChain.On("ID").Return(curTH.ChainID)
		mockChain.On("LogPoller").Return(curTH.LogPoller)
	})

	// Block 2 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(finalityDepth); i++ {
		th.Client.Commit()
	}

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	numBlocks := 5
	th.Client.Commit()
	for i := 0; i < numBlocks; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Client.Commit()
		// Make VRF requests and fulfill them
		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
			big.NewInt(int64(2*i)), big.NewInt(105), 1, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Client.Commit()
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			[32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes()),
			big.NewInt(int64(2*i+1)), big.NewInt(106), 1, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, big.NewInt(int64(2*i)), big.NewInt(105), big.NewInt(10), true)
		require.NoError(t, err2)
		th.Client.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(th.Ctx))

	latestBlock := finalityDepth + int64(numBlocks) + 1
	require.NoError(t, th.LogPoller.Replay(th.Ctx, latestBlock))

	// initializeLastProcessedBlock must return the earliest unfulfilled VRF request block
	// number instead of finalizedBlockNumber
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(th.Ctx)
	require.Nil(t, err)
	require.Equal(t, int64(8), lastProcessedBlock)
}
