package v2

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/vrf_log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	emitterABI, _    = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
	vrfEmitterABI, _ = abi.JSON(strings.NewReader(vrf_log_emitter.VRFLogEmitterABI))
)

type vrfLogPollerListenerTH struct {
	FinalityDepth     int64
	Lggr              logger.Logger
	ChainID           *big.Int
	ORM               logpoller.ORM
	LogPoller         logpoller.LogPollerTest
	Backend           *simulated.Backend
	Emitter           *log_emitter.LogEmitter
	EmitterAddress    common.Address
	VRFLogEmitter     *vrf_log_emitter.VRFLogEmitter
	VRFEmitterAddress common.Address
	Owner             *bind.TransactOpts
	Db                *sqlx.DB
	Listener          *listenerV2
}

func setupVRFLogPollerListenerTH(t *testing.T) *vrfLogPollerListenerTH {
	const (
		useFinalityTag           = false
		finalityDepth            = 3
		backfillBatchSize        = 3
		rpcBatchSize             = 2
		keepFinalizedBlocksDepth = 1000
	)

	ctx := testutils.Context(t)

	lggr := logger.TestLogger(t)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	o := logpoller.NewORM(chainID, db, lggr)
	owner := testutils.MustNewSimTransactor(t)
	backend := simulated.NewBackend(ethtypes.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	ec := backend.Client()

	h, err := ec.HeaderByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	require.LessOrEqual(t, h.Time, uint64(math.MaxInt64))
	blockTime := time.Unix(int64(h.Time), 0) //nolint:gosec // G115 false positive
	// VRF Listener relies on block timestamps, but SimulatedBackend uses by default clock starting from 1970-01-01
	// This trick is used to move the clock closer to the current time. We set first block to be 24 hours ago.
	err = backend.AdjustTime(time.Since(blockTime) - 24*time.Hour)
	require.NoError(t, err)
	backend.Commit()

	esc := client.NewSimulatedBackendClient(t, backend, chainID)
	// Mark genesis block as finalized to avoid any nulls in the tests
	client.FinalizeLatest(t, esc.Backend())

	// Poll period doesn't matter, we intend to call poll and save logs directly in the test.
	// Set it to some insanely high value to not interfere with any tests.

	lpOpts := logpoller.Opts{
		PollPeriod:               time.Hour,
		UseFinalityTag:           useFinalityTag,
		FinalityDepth:            finalityDepth,
		BackfillBatchSize:        backfillBatchSize,
		RpcBatchSize:             rpcBatchSize,
		KeepFinalizedBlocksDepth: keepFinalizedBlocksDepth,
	}
	ht := headtracker.NewSimulatedHeadTracker(esc, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(o, esc, lggr, ht, lpOpts)

	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	vrfLogEmitterAddress, _, vrfLogEmitter, err := vrf_log_emitter.DeployVRFLogEmitter(owner, ec)
	require.NoError(t, err)
	backend.Commit()

	// Log Poller Listener
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
		EVMChainID:          chainID.String(),
	}).Toml())
	require.NoError(t, err)

	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(vrfLogEmitter.Address(), ec)
	require.Nil(t, err)
	coordinator := NewCoordinatorV2(coordinatorV2)

	chain := evmmocks.NewChain(t)
	chain.On("ID").Maybe().Return(chainID)
	chain.On("LogPoller").Maybe().Return(lp)

	listener := &listenerV2{
		respCount:     map[string]uint64{},
		job:           j,
		chain:         chain,
		l:             logger.Sugared(lggr),
		coordinator:   coordinator,
		inflightCache: vrfcommon.NewInflightCache(10),
		chStop:        make(chan struct{}),
	}

	// Filter registration is idempotent, so we can just call it every time
	// and retry on errors using the ticker.
	err = lp.RegisterFilter(ctx, logpoller.Filter{
		Name: fmt.Sprintf("vrf_%s_keyhash_%s_job_%d", "v2", listener.job.VRFSpec.PublicKey.MustHash().String(), listener.job.ID),
		EventSigs: evmtypes.HashArray{
			vrf_log_emitter.VRFLogEmitterRandomWordsRequested{}.Topic(),
			vrf_log_emitter.VRFLogEmitterRandomWordsFulfilled{}.Topic(),
		},
		Addresses: evmtypes.AddressArray{
			vrfLogEmitter.Address(),
			// listener.job.VRFSpec.CoordinatorAddress.Address(),
		},
	})
	require.Nil(t, err)
	require.NoError(t, lp.RegisterFilter(ctx, logpoller.Filter{
		Name:      "Integration test",
		EventSigs: []common.Hash{emitterABI.Events["Log1"].ID},
		Addresses: []common.Address{emitterAddress1},
		Retention: 0}))
	require.Nil(t, err)
	require.Len(t, lp.Filter(nil, nil, nil).Addresses, 2)
	require.Len(t, lp.Filter(nil, nil, nil).Topics, 1)
	require.Len(t, lp.Filter(nil, nil, nil).Topics[0], 3)

	th := &vrfLogPollerListenerTH{
		FinalityDepth:     finalityDepth,
		Lggr:              lggr,
		ChainID:           chainID,
		ORM:               o,
		LogPoller:         lp,
		Emitter:           emitter1,
		EmitterAddress:    emitterAddress1,
		VRFLogEmitter:     vrfLogEmitter,
		VRFEmitterAddress: vrfLogEmitterAddress,
		Backend:           backend,
		Owner:             owner,
		Db:                db,
		Listener:          listener,
	}
	return th
}

/* Tests for initializeLastProcessedBlock: BEGIN
 * TestInitProcessedBlock_NoVRFReqs
 * TestInitProcessedBlock_NoUnfulfilledVRFReqs
 * TestInitProcessedBlock_OneUnfulfilledVRFReq
 * TestInitProcessedBlock_SomeUnfulfilledVRFReqs
 * TestInitProcessedBlock_UnfulfilledNFulfilledVRFReqs
 */

func TestInitProcessedBlock_NoVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Emit some logs from block 5 to 9 (Inclusive)
	for i := 0; i < 5; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 5 (EmitLog blocks) = 9

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	servicetest.Run(t, th.LogPoller)

	// The poller starts on a new chain at latest-finality (finalityDepth + 5 in this case),
	// Replaying from block 4 should guarantee we have block 4 immediately.  (We will also get
	// block 3 once the backup poller runs, since it always starts 100 blocks behind.)
	require.NoError(t, th.LogPoller.Replay(testutils.Context(t), 4))

	// Should return logs from block 5 to 7 (inclusive)
	logs, err := th.LogPoller.Logs(testutils.Context(t), 4, 7, emitterABI.Events["Log1"].ID, th.EmitterAddress)
	require.NoError(t, err)
	require.Equal(t, 3, len(logs))

	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(6), lastProcessedBlock)
}

func TestLogPollerFilterRegistered(t *testing.T) {
	t.Parallel()
	// Instantiate listener.
	th := setupVRFLogPollerListenerTH(t)

	func() {
		// Run the log listener. This should register the log poller filter.
		go th.Listener.runLogListener(time.Second, 1)
		// Close the listener to avoid an orphaned goroutine.
		defer close(th.Listener.chStop)

		// Wait for the log poller filter to be registered.
		filterName := th.Listener.getLogPollerFilterName()
		require.Eventually(t, func() bool {
			return th.Listener.chain.LogPoller().HasFilter(filterName)
		}, testutils.WaitTimeout(t), time.Second)

		// Once registered, expect the filter to stay registered.
		gomega.NewWithT(t).Consistently(func() bool {
			return th.Listener.chain.LogPoller().HasFilter(filterName)
		}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
	}()

	// Assert channel is closed.
	_, ok := <-th.Listener.chStop
	assert.False(t, ok)
}

func TestInitProcessedBlock_NoUnfulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Create VRF request block and a fulfillment block
	keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
	preSeed := big.NewInt(105)
	subID := uint64(1)
	reqID := big.NewInt(1)
	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()
	_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, reqID, preSeed, big.NewInt(10), true)
	require.NoError(t, err2)
	th.Backend.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 2 (VRF req/resp block) + 5 (EmitLog blocks) = 11
	latestBlock := int64(2 + 2 + 2 + 5)

	// A replay is needed so that log poller has a latest block
	// Replay from block 11 (latest) onwards, so that log poller has a latest block
	// Then test if log poller is able to replay from finalizedBlockNumber (8 --> onwards)
	// since there are no pending VRF requests
	// Blocks: 1 2 3 4 [5;Request] [6;Fulfilment] 7 8 9 10 11
	require.NoError(t, th.LogPoller.Replay(ctx, latestBlock))

	// initializeLastProcessedBlock must return the finalizedBlockNumber (8) instead of
	// VRF request block number (5), since all VRF requests are fulfilled
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(8), lastProcessedBlock)
}

func TestInitProcessedBlock_OneUnfulfilledVRFReq(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Make a VRF request without fulfilling it
	keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
	preSeed := big.NewInt(105)
	subID := uint64(1)
	reqID := big.NewInt(1)
	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	th.Backend.Commit()
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 1 (VRF req block) + 5 (EmitLog blocks) = 10
	latestBlock := int64(2 + 2 + 1 + 5)

	// A replay is needed so that log poller has a latest block
	// Replay from block 10 (latest) onwards, so that log poller has a latest block
	// Then test if log poller is able to replay from earliestUnprocessedBlock (5 --> onwards)
	// Blocks: 1 2 3 4 [5;Request] 6 7 8 9 10
	require.NoError(t, th.LogPoller.Replay(ctx, latestBlock))

	// initializeLastProcessedBlock must return the unfulfilled VRF
	// request block number (5) instead of finalizedBlockNumber (8)
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(5), lastProcessedBlock)
}

func TestInitProcessedBlock_SomeUnfulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Emit some logs in blocks with VRF reqs interspersed
	// No fulfillment for any VRF requests
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()

		// Create 2 blocks with VRF requests in each iteration
		keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
		preSeed := big.NewInt(105)
		subID := uint64(1)
		reqID1 := big.NewInt(int64(2 * i))
		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()

		reqID2 := big.NewInt(int64(2*i + 1))
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID2, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 3*5 (EmitLog + VRF req/resp blocks) = 19
	latestBlock := int64(2 + 2 + 3*5)

	// A replay is needed so that log poller has a latest block
	// Replay from block 19 (latest) onwards, so that log poller has a latest block
	// Then test if log poller is able to replay from earliestUnprocessedBlock (6 --> onwards)
	// Blocks: 1 2 3 4 5 [6;Request] [7;Request] 8 [9;Request] [10;Request]
	// 11 [12;Request] [13;Request] 14 [15;Request] [16;Request]
	// 17 [18;Request] [19;Request]
	require.NoError(t, th.LogPoller.Replay(ctx, latestBlock))

	// initializeLastProcessedBlock must return the earliest unfulfilled VRF request block
	// number instead of finalizedBlockNumber
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(6), lastProcessedBlock)
}

func TestInitProcessedBlock_UnfulfilledNFulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Emit some logs in blocks with VRF reqs interspersed
	// One VRF request in each iteration is fulfilled to imitate mixed workload
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()

		// Create 2 blocks with VRF requests in each iteration and fulfill one
		// of them. This creates a mixed workload of fulfilled and unfulfilled
		// VRF requests for testing the VRF listener
		keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
		preSeed := big.NewInt(105)
		subID := uint64(1)
		reqID1 := big.NewInt(int64(2 * i))
		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()

		reqID2 := big.NewInt(int64(2*i + 1))
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID2, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)

		_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, reqID1, preSeed, big.NewInt(10), true)
		require.NoError(t, err2)
		th.Backend.Commit()
	}

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 3*5 (EmitLog + VRF req/resp blocks) = 19
	latestBlock := int64(2 + 2 + 3*5)
	// A replay is needed so that log poller has a latest block
	// Replay from block 19 (latest) onwards, so that log poller has a latest block
	// Then test if log poller is able to replay from earliestUnprocessedBlock (7 --> onwards)
	// Blocks: 1 2 3 4 5 [6;Request] [7;Request;6-Fulfilment] 8 [9;Request] [10;Request;9-Fulfilment]
	// 11 [12;Request] [13;Request;12-Fulfilment] 14 [15;Request] [16;Request;15-Fulfilment]
	// 17 [18;Request] [19;Request;18-Fulfilment]
	require.NoError(t, th.LogPoller.Replay(ctx, latestBlock))

	// initializeLastProcessedBlock must return the earliest unfulfilled VRF request block
	// number instead of finalizedBlockNumber
	lastProcessedBlock, err := th.Listener.initializeLastProcessedBlock(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(7), lastProcessedBlock)
}

/* Tests for initializeLastProcessedBlock: END */

/* Tests for updateLastProcessedBlock: BEGIN
 * TestUpdateLastProcessedBlock_NoVRFReqs
 * TestUpdateLastProcessedBlock_NoUnfulfilledVRFReqs
 * TestUpdateLastProcessedBlock_OneUnfulfilledVRFReq
 * TestUpdateLastProcessedBlock_SomeUnfulfilledVRFReqs
 * TestUpdateLastProcessedBlock_UnfulfilledNFulfilledVRFReqs
 */

func TestUpdateLastProcessedBlock_NoVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Create VRF request logs
	keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
	preSeed := big.NewInt(105)
	subID := uint64(1)
	reqID1 := big.NewInt(int64(1))

	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()

	reqID2 := big.NewInt(int64(2))
	_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID2, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 2 (VRF req blocks) + 5 (EmitLog blocks) = 11

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// We've to replay from before VRF request log, since updateLastProcessedBlock
	// does not internally call LogPoller.Replay
	require.NoError(t, th.LogPoller.Replay(ctx, 4))

	// updateLastProcessedBlock must return the finalizedBlockNumber as there are
	// no VRF requests, after currLastProcessedBlock (block 6). The VRF requests
	// made above are before the currLastProcessedBlock (7) passed in below
	lastProcessedBlock, err := th.Listener.updateLastProcessedBlock(ctx, 7)
	require.Nil(t, err)
	require.Equal(t, int64(8), lastProcessedBlock)
}

func TestUpdateLastProcessedBlock_NoUnfulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Create VRF request log block with a fulfillment log block
	keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
	preSeed := big.NewInt(105)
	subID := uint64(1)
	reqID1 := big.NewInt(int64(1))

	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()

	_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, reqID1, preSeed, big.NewInt(10), true)
	require.NoError(t, err2)
	th.Backend.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 2 (VRF req/resp blocks) + 5 (EmitLog blocks) = 11

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// We've to replay from before VRF request log, since updateLastProcessedBlock
	// does not internally call LogPoller.Replay
	require.NoError(t, th.LogPoller.Replay(ctx, 4))

	// updateLastProcessedBlock must return the finalizedBlockNumber (8) though we have
	// a VRF req at block (5) after currLastProcessedBlock (4) passed below, because
	// the VRF request is fulfilled
	lastProcessedBlock, err := th.Listener.updateLastProcessedBlock(ctx, 4)
	require.Nil(t, err)
	require.Equal(t, int64(8), lastProcessedBlock)
}

func TestUpdateLastProcessedBlock_OneUnfulfilledVRFReq(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Create VRF request logs without a fulfillment log block
	keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
	preSeed := big.NewInt(105)
	subID := uint64(1)
	reqID1 := big.NewInt(int64(1))

	_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
		keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
	require.NoError(t, err2)
	th.Backend.Commit()

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 1 (VRF req block) + 5 (EmitLog blocks) = 10

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// We've to replay from before VRF request log, since updateLastProcessedBlock
	// does not internally call LogPoller.Replay
	require.NoError(t, th.LogPoller.Replay(ctx, 4))

	// updateLastProcessedBlock must return the VRF req at block (5) instead of
	// finalizedBlockNumber (8) after currLastProcessedBlock (4) passed below,
	// because the VRF request is unfulfilled
	lastProcessedBlock, err := th.Listener.updateLastProcessedBlock(ctx, 4)
	require.Nil(t, err)
	require.Equal(t, int64(5), lastProcessedBlock)
}

func TestUpdateLastProcessedBlock_SomeUnfulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()

		// Create 2 blocks with VRF requests in each iteration
		keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
		preSeed := big.NewInt(105)
		subID := uint64(1)
		reqID1 := big.NewInt(int64(2 * i))

		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()

		reqID2 := big.NewInt(int64(2*i + 1))
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID2, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 3*5 (EmitLog + VRF req blocks) = 19

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// We've to replay from before VRF request log, since updateLastProcessedBlock
	// does not internally call LogPoller.Replay
	require.NoError(t, th.LogPoller.Replay(ctx, 4))

	// updateLastProcessedBlock must return the VRF req at block (6) instead of
	// finalizedBlockNumber (16) after currLastProcessedBlock (4) passed below,
	// as block 6 contains the earliest unfulfilled VRF request
	lastProcessedBlock, err := th.Listener.updateLastProcessedBlock(ctx, 4)
	require.Nil(t, err)
	require.Equal(t, int64(6), lastProcessedBlock)
}

func TestUpdateLastProcessedBlock_UnfulfilledNFulfilledVRFReqs(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	t.Parallel()
	ctx := tests.Context(t)

	th := setupVRFLogPollerListenerTH(t)

	// Block 3 to finalityDepth. Ensure we have finality number of blocks
	for i := 1; i < int(th.FinalityDepth); i++ {
		th.Backend.Commit()
	}

	// Emit some logs in blocks to make the VRF req and fulfillment older than finalityDepth from latestBlock
	n := 5
	for i := 0; i < n; i++ {
		_, err1 := th.Emitter.EmitLog1(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		_, err1 = th.Emitter.EmitLog2(th.Owner, []*big.Int{big.NewInt(int64(i))})
		require.NoError(t, err1)
		th.Backend.Commit()

		// Create 2 blocks with VRF requests in each iteration and fulfill one
		// of them. This creates a mixed workload of fulfilled and unfulfilled
		// VRF requests for testing the VRF listener
		keyHash := [32]byte(th.Listener.job.VRFSpec.PublicKey.MustHash().Bytes())
		preSeed := big.NewInt(105)
		subID := uint64(1)
		reqID1 := big.NewInt(int64(2 * i))

		_, err2 := th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID1, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		th.Backend.Commit()

		reqID2 := big.NewInt(int64(2*i + 1))
		_, err2 = th.VRFLogEmitter.EmitRandomWordsRequested(th.Owner,
			keyHash, reqID2, preSeed, subID, 10, 10000, 2, th.Owner.From)
		require.NoError(t, err2)
		_, err2 = th.VRFLogEmitter.EmitRandomWordsFulfilled(th.Owner, reqID1, preSeed, big.NewInt(10), true)
		require.NoError(t, err2)
		th.Backend.Commit()
	}

	// Blocks till now: 2 (in SetupTH) + 2 (empty blocks) + 3*5 (EmitLog + VRF req blocks) = 19

	// Calling Start() after RegisterFilter() simulates a node restart after job creation, should reload Filter from db.
	require.NoError(t, th.LogPoller.Start(ctx))

	// We've to replay from before VRF request log, since updateLastProcessedBlock
	// does not internally call LogPoller.Replay
	require.NoError(t, th.LogPoller.Replay(ctx, 4))

	// updateLastProcessedBlock must return the VRF req at block (7) instead of
	// finalizedBlockNumber (16) after currLastProcessedBlock (4) passed below,
	// as block 7 contains the earliest unfulfilled VRF request. VRF request
	// in block 6 has been fulfilled in block 7.
	lastProcessedBlock, err := th.Listener.updateLastProcessedBlock(ctx, 4)
	require.Nil(t, err)
	require.Equal(t, int64(7), lastProcessedBlock)
}

/* Tests for updateLastProcessedBlock: END */

/* Tests for getUnfulfilled: BEGIN
 * TestGetUnfulfilled_NoVRFReqs
 * TestGetUnfulfilled_NoUnfulfilledVRFReqs
 * TestGetUnfulfilled_OneUnfulfilledVRFReq
 * TestGetUnfulfilled_SomeUnfulfilledVRFReqs
 * TestGetUnfulfilled_UnfulfilledNFulfilledVRFReqs
 */

func SetupGetUnfulfilledTH(t *testing.T) (*listenerV2, *ubig.Big) {
	chainID := ubig.New(big.NewInt(12345))
	lggr := logger.TestLogger(t)
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
		EVMChainID:          chainID.String(),
	}).Toml())
	require.NoError(t, err)
	chain := evmmocks.NewChain(t)

	// Construct CoordinatorV2_X object for VRF listener
	owner := testutils.MustNewSimTransactor(t)
	b := simulated.NewBackend(ethtypes.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	_, _, vrfLogEmitter, err := vrf_log_emitter.DeployVRFLogEmitter(owner, b.Client())
	require.NoError(t, err)
	b.Commit()
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(vrfLogEmitter.Address(), b.Client())
	require.Nil(t, err)
	coordinator := NewCoordinatorV2(coordinatorV2)

	listener := &listenerV2{
		respCount:   map[string]uint64{},
		job:         j,
		chain:       chain,
		l:           logger.Sugared(lggr),
		coordinator: coordinator,
	}
	return listener, chainID
}

func TestGetUnfulfilled_NoVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := 0; i < 10; i++ {
		logs = append(logs, logpoller.Log{
			EvmChainId:     chainID,
			LogIndex:       0,
			BlockHash:      common.BigToHash(big.NewInt(int64(i))),
			BlockNumber:    int64(i),
			BlockTimestamp: time.Now(),
			Topics: [][]byte{
				[]byte("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
			},
			EventSig:  emitterABI.Events["Log1"].ID,
			Address:   common.Address{},
			TxHash:    common.BigToHash(big.NewInt(int64(i))),
			Data:      nil,
			CreatedAt: time.Now(),
		})
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Empty(t, unfulfilled)
	require.Empty(t, fulfilled)
}

func TestGetUnfulfilled_NoUnfulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := 0; i < 10; i++ {
		eventSig := emitterABI.Events["Log1"].ID
		topics := [][]byte{
			common.FromHex("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
		}
		if i%2 == 0 {
			eventSig = vrfEmitterABI.Events["RandomWordsRequested"].ID
			topics = [][]byte{
				common.FromHex("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772"),
				common.FromHex("0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf"),
				common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"),
				common.FromHex("0x0000000000000000000000005ee3b50502b5c4c9184dcb281471a0614d4b2ef9"),
			}
		}
		logs = append(logs, logpoller.Log{
			EvmChainId:     chainID,
			LogIndex:       0,
			BlockHash:      common.BigToHash(big.NewInt(int64(2 * i))),
			BlockNumber:    int64(2 * i),
			BlockTimestamp: time.Now(),
			Topics:         topics,
			EventSig:       eventSig,
			Address:        common.Address{},
			TxHash:         common.BigToHash(big.NewInt(int64(2 * i))),
			Data:           common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i) + "000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000027100000000000000000000000000000000000000000000000000000000000000002"),
			CreatedAt:      time.Now(),
		})
		if i%2 == 0 {
			logs = append(logs, logpoller.Log{
				EvmChainId:     chainID,
				LogIndex:       0,
				BlockHash:      common.BigToHash(big.NewInt(int64(2*i + 1))),
				BlockNumber:    int64(2*i + 1),
				BlockTimestamp: time.Now(),
				Topics: [][]byte{
					common.FromHex("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4"),
					common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i)),
				},
				EventSig:  vrfEmitterABI.Events["RandomWordsFulfilled"].ID,
				Address:   common.Address{},
				TxHash:    common.BigToHash(big.NewInt(int64(2*i + 1))),
				Data:      common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000069000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001"),
				CreatedAt: time.Now(),
			})
		}
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Empty(t, unfulfilled)
	require.Len(t, fulfilled, 5)
}

func TestGetUnfulfilled_OneUnfulfilledVRFReq(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := 0; i < 10; i++ {
		eventSig := emitterABI.Events["Log1"].ID
		topics := [][]byte{
			common.FromHex("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
		}
		if i == 4 {
			eventSig = vrfEmitterABI.Events["RandomWordsRequested"].ID
			topics = [][]byte{
				common.FromHex("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772"),
				common.FromHex("0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf"),
				common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"),
				common.FromHex("0x0000000000000000000000005ee3b50502b5c4c9184dcb281471a0614d4b2ef9"),
			}
		}
		logs = append(logs, logpoller.Log{
			EvmChainId:     chainID,
			LogIndex:       0,
			BlockHash:      common.BigToHash(big.NewInt(int64(2 * i))),
			BlockNumber:    int64(2 * i),
			BlockTimestamp: time.Now(),
			Topics:         topics,
			EventSig:       eventSig,
			Address:        common.Address{},
			TxHash:         common.BigToHash(big.NewInt(int64(2 * i))),
			Data:           common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i) + "000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000027100000000000000000000000000000000000000000000000000000000000000002"),
			CreatedAt:      time.Now(),
		})
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Equal(t, unfulfilled[0].RequestID().Int64(), big.NewInt(4).Int64())
	require.Len(t, unfulfilled, 1)
	require.Empty(t, fulfilled)
}

func TestGetUnfulfilled_SomeUnfulfilledVRFReq(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := 0; i < 10; i++ {
		eventSig := emitterABI.Events["Log1"].ID
		topics := [][]byte{
			common.FromHex("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
		}
		if i%2 == 0 {
			eventSig = vrfEmitterABI.Events["RandomWordsRequested"].ID
			topics = [][]byte{
				common.FromHex("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772"),
				common.FromHex("0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf"),
				common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"),
				common.FromHex("0x0000000000000000000000005ee3b50502b5c4c9184dcb281471a0614d4b2ef9"),
			}
		}
		logs = append(logs, logpoller.Log{
			EvmChainId:     chainID,
			LogIndex:       0,
			BlockHash:      common.BigToHash(big.NewInt(int64(2 * i))),
			BlockNumber:    int64(2 * i),
			BlockTimestamp: time.Now(),
			Topics:         topics,
			EventSig:       eventSig,
			Address:        common.Address{},
			TxHash:         common.BigToHash(big.NewInt(int64(2 * i))),
			Data:           common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i) + "000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000027100000000000000000000000000000000000000000000000000000000000000002"),
			CreatedAt:      time.Now(),
		})
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Len(t, unfulfilled, 5)
	require.Len(t, fulfilled, 0)
	expected := map[int64]bool{0: true, 2: true, 4: true, 6: true, 8: true}
	for _, u := range unfulfilled {
		v, ok := expected[u.RequestID().Int64()]
		require.Equal(t, ok, true)
		require.Equal(t, v, true)
	}
	require.Equal(t, len(expected), len(unfulfilled))
}

func TestGetUnfulfilled_UnfulfilledNFulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := 0; i < 10; i++ {
		eventSig := emitterABI.Events["Log1"].ID
		topics := [][]byte{
			common.FromHex("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
		}
		if i%2 == 0 {
			eventSig = vrfEmitterABI.Events["RandomWordsRequested"].ID
			topics = [][]byte{
				common.FromHex("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772"),
				common.FromHex("0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf"),
				common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"),
				common.FromHex("0x0000000000000000000000005ee3b50502b5c4c9184dcb281471a0614d4b2ef9"),
			}
		}
		logs = append(logs, logpoller.Log{
			EvmChainId:     chainID,
			LogIndex:       0,
			BlockHash:      common.BigToHash(big.NewInt(int64(2 * i))),
			BlockNumber:    int64(2 * i),
			BlockTimestamp: time.Now(),
			Topics:         topics,
			EventSig:       eventSig,
			Address:        common.Address{},
			TxHash:         common.BigToHash(big.NewInt(int64(2 * i))),
			Data:           common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i) + "000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000027100000000000000000000000000000000000000000000000000000000000000002"),
			CreatedAt:      time.Now(),
		})
		if i%2 == 0 && i < 6 {
			logs = append(logs, logpoller.Log{
				EvmChainId:     chainID,
				LogIndex:       0,
				BlockHash:      common.BigToHash(big.NewInt(int64(2*i + 1))),
				BlockNumber:    int64(2*i + 1),
				BlockTimestamp: time.Now(),
				Topics: [][]byte{
					common.FromHex("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4"),
					common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%d", i)),
				},
				EventSig:  vrfEmitterABI.Events["RandomWordsFulfilled"].ID,
				Address:   common.Address{},
				TxHash:    common.BigToHash(big.NewInt(int64(2*i + 1))),
				Data:      common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000069000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001"),
				CreatedAt: time.Now(),
			})
		}
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Len(t, unfulfilled, 2)
	require.Len(t, fulfilled, 3)
	expected := map[int64]bool{6: true, 8: true}
	for _, u := range unfulfilled {
		v, ok := expected[u.RequestID().Int64()]
		require.Equal(t, ok, true)
		require.Equal(t, v, true)
	}
	require.Equal(t, len(expected), len(unfulfilled))
}

/* Tests for getUnfulfilled: END */
