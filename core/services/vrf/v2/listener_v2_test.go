package v2

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	addEthTxQuery = `INSERT INTO evm.txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
		VALUES (
		$1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9, $10, $11
		)
		RETURNING "txes".*`

	addConfirmedEthTxQuery = `INSERT INTO evm.txes (nonce, broadcast_at, initial_broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
		VALUES (
		$1, NOW(), NOW(), NULL, $2, $3, $4, $5, $6, 'confirmed', NOW(), $7, $8, $9, $10, $11
		)
		RETURNING "txes".*`
)

func txMetaSubIDs(t *testing.T, vrfVersion vrfcommon.Version, subID *big.Int) (*uint64, *string) {
	var (
		txMetaSubID       *uint64
		txMetaGlobalSubID *string
	)
	if vrfVersion == vrfcommon.V2Plus {
		txMetaGlobalSubID = ptr(subID.String())
	} else if vrfVersion == vrfcommon.V2 {
		txMetaSubID = ptr(subID.Uint64())
	} else {
		t.Errorf("unsupported vrf version: %s", vrfVersion)
	}
	return txMetaSubID, txMetaGlobalSubID
}

func addEthTx(t *testing.T, db *sqlx.DB, from common.Address, state txmgrtypes.TxState, maxLink string, subID *big.Int, reqTxHash common.Hash, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	_, err := db.Exec(addEthTxQuery,
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		state,
		txmgr.TxMeta{
			MaxLink:       &maxLink,
			SubID:         txMetaSubID,
			GlobalSubID:   txMetaGlobalSubID,
			RequestTxHash: &reqTxHash,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *sqlx.DB, from common.Address, maxLink string, subID *big.Int, nonce uint64, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	_, err := db.Exec(addConfirmedEthTxQuery,
		nonce,          // nonce
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		txmgr.TxMeta{
			MaxLink:     &maxLink,
			SubID:       txMetaSubID,
			GlobalSubID: txMetaGlobalSubID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

func addEthTxNativePayment(t *testing.T, db *sqlx.DB, from common.Address, state txmgrtypes.TxState, maxNative string, subID *big.Int, reqTxHash common.Hash, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	_, err := db.Exec(addEthTxQuery,
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		state,
		txmgr.TxMeta{
			MaxEth:        &maxNative,
			SubID:         txMetaSubID,
			GlobalSubID:   txMetaGlobalSubID,
			RequestTxHash: &reqTxHash,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

func addConfirmedEthTxNativePayment(t *testing.T, db *sqlx.DB, from common.Address, maxNative string, subID *big.Int, nonce uint64, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	_, err := db.Exec(addConfirmedEthTxQuery,
		nonce,          // nonce
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		txmgr.TxMeta{
			MaxEth:      &maxNative,
			SubID:       txMetaSubID,
			GlobalSubID: txMetaGlobalSubID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

func testMaybeSubtractReservedLink(t *testing.T, vrfVersion vrfcommon.Version) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := pgtest.NewQConfig(false)
	q := pg.NewQ(db, lggr, cfg)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr, cfg)
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	k, err := ks.Eth().Create(big.NewInt(int64(chainID)))
	require.NoError(t, err)

	subID := new(big.Int).SetUint64(1)
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err := MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)

	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address, "10000", subID, 1, vrfVersion)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTx(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved link should not affect other subscribers prospective balance.
	otherSubID := new(big.Int).SetUint64(2)
	require.NoError(t, err)
	addEthTx(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", otherSubID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	anotherSubID := new(big.Int).SetUint64(3)
	addEthTx(t, db, k2.Address, txmgrcommon.TxUnstarted, "10000", anotherSubID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the link reserved across multiple keys,
	// i.e, gas lanes.
	addEthTx(t, db, k2.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}

func TestMaybeSubtractReservedLinkV2(t *testing.T) {
	testMaybeSubtractReservedLink(t, vrfcommon.V2)
}

func TestMaybeSubtractReservedLinkV2Plus(t *testing.T) {
	testMaybeSubtractReservedLink(t, vrfcommon.V2Plus)
}

func testMaybeSubtractReservedNative(t *testing.T, vrfVersion vrfcommon.Version) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := pgtest.NewQConfig(false)
	q := pg.NewQ(db, lggr, cfg)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr, cfg)
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	k, err := ks.Eth().Create(big.NewInt(int64(chainID)))
	require.NoError(t, err)

	subID := new(big.Int).SetUint64(1)
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")

	// Insert an unstarted eth tx with native metadata
	addEthTxNativePayment(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err := MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)

	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTxNativePayment(t, db, k.Address, "10000", subID, 1, vrfVersion)
	start, err = MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTxNativePayment(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved native should not affect other subscribers prospective balance.
	otherSubID := new(big.Int).SetUint64(2)
	require.NoError(t, err)
	addEthTxNativePayment(t, db, k.Address, txmgrcommon.TxUnstarted, "10000", otherSubID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	anotherSubID := new(big.Int).SetUint64(3)
	addEthTxNativePayment(t, db, k2.Address, txmgrcommon.TxUnstarted, "10000", anotherSubID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the native reserved across multiple keys,
	// i.e, gas lanes.
	addEthTxNativePayment(t, db, k2.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}

func TestMaybeSubtractReservedNativeV2Plus(t *testing.T) {
	testMaybeSubtractReservedNative(t, vrfcommon.V2Plus)
}

func TestMaybeSubtractReservedNativeV2(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := pgtest.NewQConfig(false)
	q := pg.NewQ(db, lggr, cfg)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr, cfg)
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	subID := new(big.Int).SetUint64(1)
	// returns error because native payment is not supported for V2
	start, err := MaybeSubtractReservedEth(q, big.NewInt(100_000), chainID, subID, vrfcommon.V2)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), start)
}

func TestListener_GetConfirmedAt(t *testing.T) {
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)

	listener := &listenerV2{
		respCount: map[string]uint64{},
		job:       j,
	}

	// Requester asks for 100 confirmations, we have a delay of 10,
	// so we should wait for max(nodeMinConfs, requestedConfs + requestedConfsDelay) = 110 confirmations
	nodeMinConfs := 10
	confirmedAt := listener.getConfirmedAt(NewV2RandomWordsRequested(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}), uint32(nodeMinConfs))
	require.Equal(t, uint64(210), confirmedAt) // log block number + # of confirmations

	// Requester asks for 100 confirmations, we have a delay of 0,
	// so we should wait for max(nodeMinConfs, requestedConfs + requestedConfsDelay) = 100 confirmations
	j, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 0,
	}).Toml())
	require.NoError(t, err)
	listener.job = j
	confirmedAt = listener.getConfirmedAt(NewV2RandomWordsRequested(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}), uint32(nodeMinConfs))
	require.Equal(t, uint64(200), confirmedAt) // log block number + # of confirmations
}

func TestListener_Backoff(t *testing.T) {
	var tests = []struct {
		name     string
		initial  time.Duration
		max      time.Duration
		last     time.Duration
		retries  int
		expected bool
	}{
		{
			name:     "Backoff disabled, ready",
			expected: true,
		},
		{
			name:     "First try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     0,
			retries:  0,
			expected: true,
		},
		{
			name:     "Second try, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     59 * time.Second,
			retries:  1,
			expected: false,
		},
		{
			name:     "Second try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     61 * time.Second, // Last try was over a minute ago
			retries:  1,
			expected: true,
		},
		{
			name:     "Third try, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     77 * time.Second, // Slightly less than backoffFactor * initial
			retries:  2,
			expected: false,
		},
		{
			name:     "Third try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     79 * time.Second, // Slightly more than backoffFactor * initial
			retries:  2,
			expected: true,
		},
		{
			name:     "Max, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     59 * time.Minute, // Slightly less than max
			retries:  900,
			expected: false,
		},
		{
			name:     "Max, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     61 * time.Minute, // Slightly more than max
			retries:  900,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lsn := &listenerV2{job: job.Job{
				VRFSpec: &job.VRFSpec{
					BackoffInitialDelay: test.initial,
					BackoffMaxDelay:     test.max,
				},
			}}

			req := pendingRequest{
				confirmedAtBlock: 5,
				attempts:         test.retries,
				lastTry:          time.Now().Add(-test.last),
			}

			require.Equal(t, test.expected, lsn.ready(req, 10))
		})
	}
}

func TestListener_handleLog(tt *testing.T) {
	lb := mocks.NewBroadcaster(tt)
	chainID := int64(2)
	minConfs := uint32(3)
	blockNumber := uint64(5)
	requestID := int64(6)
	tt.Run("v2", func(t *testing.T) {
		j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			VRFVersion:          vrfcommon.V2,
			RequestedConfsDelay: 10,
			FromAddresses:       []string{"0xF2982b7Ef6E3D8BB738f8Ea20502229781f6Ad97"},
		}).Toml())
		require.NoError(t, err)
		fulfilledLog := vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{
			RequestId: big.NewInt(requestID),
			Raw:       types.Log{BlockNumber: blockNumber},
		}
		log := log.NewLogBroadcast(types.Log{}, *big.NewInt(chainID), &fulfilledLog)
		lb.On("WasAlreadyConsumed", log).Return(false, nil).Once()
		lb.On("MarkConsumed", log).Return(nil).Once()
		defer lb.AssertExpectations(t)
		listener := &listenerV2{
			respCount:          map[string]uint64{},
			job:                j,
			blockNumberToReqID: pairing.New(),
			latestHeadMu:       sync.RWMutex{},
			logBroadcaster:     lb,
			l:                  logger.TestLogger(t),
		}
		listener.handleLog(log, minConfs)
		require.Equal(t, listener.respCount[fulfilledLog.RequestId.String()], uint64(1))
		req, ok := listener.blockNumberToReqID.FindMin().(fulfilledReqV2)
		require.True(t, ok)
		require.Equal(t, req.blockNumber, blockNumber)
		require.Equal(t, req.reqID, "6")
	})

	tt.Run("v2 plus", func(t *testing.T) {
		j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			VRFVersion:          vrfcommon.V2Plus,
			RequestedConfsDelay: 10,
			FromAddresses:       []string{"0xF2982b7Ef6E3D8BB738f8Ea20502229781f6Ad97"},
		}).Toml())
		require.NoError(t, err)
		fulfilledLog := vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRandomWordsFulfilled{
			RequestId: big.NewInt(requestID),
			Raw:       types.Log{BlockNumber: blockNumber},
		}
		log := log.NewLogBroadcast(types.Log{}, *big.NewInt(chainID), &fulfilledLog)
		lb.On("WasAlreadyConsumed", log).Return(false, nil).Once()
		lb.On("MarkConsumed", log).Return(nil).Once()
		defer lb.AssertExpectations(t)
		listener := &listenerV2{
			respCount:          map[string]uint64{},
			job:                j,
			blockNumberToReqID: pairing.New(),
			latestHeadMu:       sync.RWMutex{},
			logBroadcaster:     lb,
			l:                  logger.TestLogger(t),
		}
		listener.handleLog(log, minConfs)
		require.Equal(t, listener.respCount[fulfilledLog.RequestId.String()], uint64(1))
		req, ok := listener.blockNumberToReqID.FindMin().(fulfilledReqV2)
		require.True(t, ok)
		require.Equal(t, req.blockNumber, blockNumber)
		require.Equal(t, req.reqID, "6")
	})

}
