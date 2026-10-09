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

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmmocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/extraargs"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

var emitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))

// v2PlusEventABI returns the parsed VRF V2Plus coordinator ABI used to pack
// synthetic RandomWordsRequested/RandomWordsFulfilled event data.
func v2PlusEventABI(t *testing.T) abi.ABI {
	parsed, err := vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalMetaData.GetAbi()
	require.NoError(t, err)
	return *parsed
}

// packV2PlusRequestedData packs the non-indexed arguments of the V2Plus
// RandomWordsRequested event, matching how an on-chain coordinator would
// encode it.
func packV2PlusRequestedData(t *testing.T, requestID, preSeed *big.Int, minConfs uint16, callbackGasLimit, numWords uint32, extraArgs []byte) []byte {
	data, err := v2PlusEventABI(t).Events["RandomWordsRequested"].Inputs.NonIndexed().Pack(
		requestID, preSeed, minConfs, callbackGasLimit, numWords, extraArgs)
	require.NoError(t, err)
	return data
}

// packV2PlusFulfilledData packs the non-indexed arguments of the V2Plus
// RandomWordsFulfilled event as emitted by the VRFCoordinatorV25 contract.
// Note the coordinator contract's event carries an extra nativePayment
// argument compared to the IVRFCoordinatorV2Plus interface declaration; the
// listener parses fulfilled logs with the coordinator contract ABI, while it
// registers log poller filters by the interface ABI (see
// coordinatorV2_5.RandomWordsFulfilledTopic). The subId argument is indexed
// (a topic), not part of the data.
func packV2PlusFulfilledData(t *testing.T, outputSeed, payment *big.Int, success, onlyPremium bool) []byte {
	parsed, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
	require.NoError(t, err)
	data, err := parsed.Events["RandomWordsFulfilled"].Inputs.NonIndexed().Pack(
		outputSeed, payment, false, success, onlyPremium)
	require.NoError(t, err)
	return data
}

type vrfLogPollerListenerTH struct {
	FinalityDepth  int64
	Lggr           logger.Logger
	ChainID        *big.Int
	ORM            logpoller.ORM
	LogPoller      logpoller.LogPollerTest
	Backend        *simulated.Backend
	Emitter        *log_emitter.LogEmitter
	EmitterAddress common.Address
	Owner          *bind.TransactOpts
	Db             *sqlx.DB
	Listener       *listenerV2
}

func setupVRFLogPollerListenerTH(t *testing.T) *vrfLogPollerListenerTH {
	const (
		useFinalityTag           = false
		finalityDepth            = 3
		backfillBatchSize        = 3
		rpcBatchSize             = 2
		keepFinalizedBlocksDepth = 1000
	)

	ctx := t.Context()

	lggr := logger.Test(t)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	o := logpoller.NewORM(chainID, db, lggr)
	owner := evmtestutils.MustNewSimTransactor(t)
	backend := simulated.NewBackend(ethtypes.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	ec := backend.Client()

	h, err := ec.HeaderByNumber(t.Context(), nil)
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
		RPCBatchSize:             rpcBatchSize,
		KeepFinalizedBlocksDepth: keepFinalizedBlocksDepth,
	}
	ht := headstest.NewSimulatedHeadTracker(esc, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(o, esc, lggr, ht, lpOpts)

	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	emitterAddress2, _, _, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	backend.Commit()

	// Log Poller Listener
	ks := keystore.NewInMemory(db, commonkeystore.FastScryptParams, lggr.Infof)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
		EVMChainID:          chainID.String(),
	}).Toml())
	require.NoError(t, err)

	// Bind the V2Plus coordinator ABI at the second emitter's address, the
	// same trick the V2-era version of this test used with the VRF log
	// emitter: the listener only parses logs through this binding.
	coordinatorV2_5, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(emitterAddress2, ec)
	require.NoError(t, err)
	coordinator := NewCoordinatorV2_5(coordinatorV2_5)

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
		Name: fmt.Sprintf("vrf_%s_keyhash_%s_job_%d", "v2plus", listener.job.VRFSpec.PublicKey.MustHash().String(), listener.job.ID),
		EventSigs: evmtypes.HashArray{
			coordinator.RandomWordsRequestedTopic(),
			coordinator.RandomWordsFulfilledTopic(),
		},
		Addresses: evmtypes.AddressArray{
			coordinator.Address(),
		},
	})
	require.NoError(t, err)
	require.NoError(t, lp.RegisterFilter(ctx, logpoller.Filter{
		Name:      "Integration test",
		EventSigs: []common.Hash{emitterABI.Events["Log1"].ID},
		Addresses: []common.Address{emitterAddress1},
		Retention: 0}))
	require.NoError(t, err)
	require.Len(t, lp.Filter(nil, nil, nil).Addresses, 2)
	require.Len(t, lp.Filter(nil, nil, nil).Topics, 1)
	require.Len(t, lp.Filter(nil, nil, nil).Topics[0], 3)

	th := &vrfLogPollerListenerTH{
		FinalityDepth:  finalityDepth,
		Lggr:           lggr,
		ChainID:        chainID,
		ORM:            o,
		LogPoller:      lp,
		Emitter:        emitter1,
		EmitterAddress: emitterAddress1,
		Backend:        backend,
		Owner:          owner,
		Db:             db,
		Listener:       listener,
	}
	return th
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

// SetupGetUnfulfilledTH returns a listener wired to a V2Plus coordinator ABI
// binding, for exercising getUnfulfilled with synthetic logpoller logs.
func SetupGetUnfulfilledTH(t *testing.T) (*listenerV2, *big.Int) {
	lggr := logger.Test(t)

	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)
	chain := evmmocks.NewChain(t)

	// Construct CoordinatorV2_X object for VRF listener
	owner := evmtestutils.MustNewSimTransactor(t)
	b := simulated.NewBackend(ethtypes.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	backend := b.Client()
	b.Commit()
	// The binding does not require a deployed coordinator: the listener only
	// parses logs through it.
	coordinatorV2_5, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(owner.From, backend)
	require.NoError(t, err)
	coordinator := NewCoordinatorV2_5(coordinatorV2_5)

	chainID := testutils.NewRandomEVMChainID()
	chain.On("ID").Maybe().Return(chainID)

	listener := &listenerV2{
		respCount:   map[string]uint64{},
		job:         j,
		chain:       chain,
		l:           logger.Sugared(lggr),
		coordinator: coordinator,
	}
	return listener, chainID
}

// syntheticV2PlusRequestedLog builds a logpoller.Log carrying a V2Plus
// RandomWordsRequested event for the given request id, with the topics and
// data encoded exactly as an on-chain coordinator would emit them.
func syntheticV2PlusRequestedLog(t *testing.T, chainID *big.Int, listener *listenerV2, requestID *big.Int, blockNumber int64) logpoller.Log {
	keyHash := listener.job.VRFSpec.PublicKey.MustHash()
	subID := big.NewInt(1)
	extraArgs, err := extraargs.EncodeV1(false)
	require.NoError(t, err)

	subIDTopic := common.BytesToHash(subID.Bytes())
	senderTopic := common.BytesToHash(common.HexToAddress("0x5ee3b50502b5c4c9184dcb281471a0614d4b2ef9").Bytes())

	return logpoller.Log{
		EVMChainID:     (*sqlutil.Big)(chainID),
		LogIndex:       0,
		BlockHash:      common.BigToHash(big.NewInt(blockNumber)),
		BlockNumber:    blockNumber,
		BlockTimestamp: time.Now(),
		Topics: [][]byte{
			listener.coordinator.RandomWordsRequestedTopic().Bytes(),
			keyHash.Bytes(),
			subIDTopic.Bytes(),
			senderTopic.Bytes(),
		},
		EventSig:  listener.coordinator.RandomWordsRequestedTopic(),
		Address:   common.Address{},
		TxHash:    common.BigToHash(big.NewInt(blockNumber)),
		Data:      packV2PlusRequestedData(t, requestID, big.NewInt(106), 10, 10000, 2, extraArgs),
		CreatedAt: time.Now(),
	}
}

// syntheticV2PlusFulfilledLog builds a logpoller.Log carrying a V2Plus
// RandomWordsFulfilled event for the given request id.
func syntheticV2PlusFulfilledLog(t *testing.T, chainID *big.Int, listener *listenerV2, requestID *big.Int, blockNumber int64) logpoller.Log {
	subID := big.NewInt(1)

	subIDTopic := common.BytesToHash(subID.Bytes())

	return logpoller.Log{
		EVMChainID:     (*sqlutil.Big)(chainID),
		LogIndex:       0,
		BlockHash:      common.BigToHash(big.NewInt(blockNumber)),
		BlockNumber:    blockNumber,
		BlockTimestamp: time.Now(),
		Topics: [][]byte{
			// The listener parses fulfilled logs with the coordinator
			// contract ABI, which requires the contract's event topic.
			vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled{}.Topic().Bytes(),
			common.BytesToHash(requestID.Bytes()).Bytes(),
			subIDTopic.Bytes(),
		},
		// getUnfulfilled dispatches on the adapter's (interface) event topic.
		EventSig: listener.coordinator.RandomWordsFulfilledTopic(),
		Address:  common.Address{},
		TxHash:   common.BigToHash(big.NewInt(blockNumber)),
		Data:     packV2PlusFulfilledData(t, big.NewInt(105), big.NewInt(10), true, false),
	}
}

// syntheticIrrelevantLog builds a logpoller.Log for an unrelated event, to
// prove getUnfulfilled ignores non-VRF logs.
func syntheticIrrelevantLog(chainID *big.Int, blockNumber int64) logpoller.Log {
	return logpoller.Log{
		EVMChainID:     (*sqlutil.Big)(chainID),
		LogIndex:       0,
		BlockHash:      common.BigToHash(big.NewInt(blockNumber)),
		BlockNumber:    blockNumber,
		BlockTimestamp: time.Now(),
		Topics: [][]byte{
			common.FromHex("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8"),
		},
		EventSig:  emitterABI.Events["Log1"].ID,
		Address:   common.Address{},
		TxHash:    common.BigToHash(big.NewInt(blockNumber)),
		CreatedAt: time.Now(),
	}
}

func TestGetUnfulfilled_NoVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := make([]logpoller.Log, 0, 10)
	for i := range 10 {
		logs = append(logs, syntheticIrrelevantLog(chainID, int64(i)))
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Empty(t, unfulfilled)
	require.Empty(t, fulfilled)
}

func TestGetUnfulfilled_NoUnfulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := range 10 {
		logs = append(logs, syntheticIrrelevantLog(chainID, int64(2*i)))
		if i%2 == 0 {
			logs = append(logs, syntheticV2PlusRequestedLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i)))
			logs = append(logs, syntheticV2PlusFulfilledLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i+1)))
		}
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Empty(t, unfulfilled)
	require.Len(t, fulfilled, 5)
}

func TestGetUnfulfilled_OneUnfulfilledVRFReq(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := make([]logpoller.Log, 0, 10)
	for i := range 10 {
		if i == 4 {
			logs = append(logs, syntheticV2PlusRequestedLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i)))
			continue
		}
		logs = append(logs, syntheticIrrelevantLog(chainID, int64(2*i)))
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Len(t, unfulfilled, 1)
	require.Equal(t, unfulfilled[0].RequestID().Int64(), big.NewInt(4).Int64())
	require.Empty(t, fulfilled)
}

func TestGetUnfulfilled_SomeUnfulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := make([]logpoller.Log, 0, 10)
	for i := range 10 {
		if i%2 == 0 {
			logs = append(logs, syntheticV2PlusRequestedLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i)))
			continue
		}
		logs = append(logs, syntheticIrrelevantLog(chainID, int64(2*i)))
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Len(t, unfulfilled, 5)
	require.Empty(t, fulfilled)
	expected := map[int64]bool{0: true, 2: true, 4: true, 6: true, 8: true}
	for _, u := range unfulfilled {
		v, ok := expected[u.RequestID().Int64()]
		require.True(t, ok)
		require.True(t, v)
	}
	require.Len(t, unfulfilled, len(expected))
}

func TestGetUnfulfilled_UnfulfilledNFulfilledVRFReqs(t *testing.T) {
	t.Parallel()

	listener, chainID := SetupGetUnfulfilledTH(t)

	logs := []logpoller.Log{}
	for i := range 10 {
		logs = append(logs, syntheticIrrelevantLog(chainID, int64(2*i)))
		if i%2 == 0 {
			logs = append(logs, syntheticV2PlusRequestedLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i)))
		}
		if i%2 == 0 && i < 6 {
			logs = append(logs, syntheticV2PlusFulfilledLog(t, chainID, listener, big.NewInt(int64(i)), int64(2*i+1)))
		}
	}

	unfulfilled, _, fulfilled := listener.getUnfulfilled(logs, listener.l)
	require.Len(t, unfulfilled, 2)
	require.Len(t, fulfilled, 3)
	expected := map[int64]bool{6: true, 8: true}
	for _, u := range unfulfilled {
		v, ok := expected[u.RequestID().Int64()]
		require.True(t, ok)
		require.True(t, v)
	}
	require.Len(t, unfulfilled, len(expected))
}
