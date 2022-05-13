package vrf

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	eth_client_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	vrf_mocks "github.com/smartcontractkit/chainlink/core/services/vrf/mocks"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func addEthTx(t *testing.T, db *sqlx.DB, from common.Address, state txmgr.EthTxState, maxLink string, subID uint64) {
	_, err := db.Exec(`INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
		VALUES (
		$1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9, $10, $11
		)
		RETURNING "eth_txes".*`,
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		state,
		txmgr.EthTxMeta{
			MaxLink: &maxLink,
			SubID:   &subID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *sqlx.DB, from common.Address, maxLink string, subID, nonce uint64) {
	_, err := db.Exec(`INSERT INTO eth_txes (nonce, broadcast_at, initial_broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
		VALUES (
		$1, NOW(), NOW(), NULL, $2, $3, $4, $5, $6, 'confirmed', NOW(), $7, $8, $9, $10, $11
		)
		RETURNING "eth_txes".*`,
		nonce,          // nonce
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		txmgr.EthTxMeta{
			MaxLink: &maxLink,
			SubID:   &subID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil)
	require.NoError(t, err)
}

// cannot import cltest, because of circular imports
type config struct{}

func (c *config) LogSQL() bool {
	return false
}

type executionRevertedError struct{}

func (executionRevertedError) Error() string {
	return "execution reverted"
}

type networkError struct{}

func (networkError) Error() string {
	return "network Error"
}

func TestMaybeSubtractReservedLink(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	q := pg.NewQ(db, lggr, &config{})
	ks := keystore.New(db, utils.FastScryptParams, lggr, &config{})
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	k, err := ks.Eth().Create(big.NewInt(int64(chainID)))
	require.NoError(t, err)

	subID := uint64(1)

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address.Address(), txmgr.EthTxUnstarted, "10000", subID)
	start, err := MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address.Address(), "10000", subID, 1)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTx(t, db, k.Address.Address(), txmgr.EthTxUnstarted, "10000", subID)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved link should not affect other subscribers prospective balance.
	otherSubID := uint64(2)
	require.NoError(t, err)
	addEthTx(t, db, k.Address.Address(), txmgr.EthTxUnstarted, "10000", otherSubID)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	anotherSubID := uint64(3)
	addEthTx(t, db, k2.Address.Address(), txmgr.EthTxUnstarted, "10000", anotherSubID)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the link reserved across multiple keys,
	// i.e, gas lanes.
	addEthTx(t, db, k2.Address.Address(), txmgr.EthTxUnstarted, "10000", subID)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}

func TestListener_GetConfirmedAt(t *testing.T) {
	j, err := ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
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
	confirmedAt := listener.getConfirmedAt(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}, uint32(nodeMinConfs))
	require.Equal(t, uint64(210), confirmedAt) // log block number + # of confirmations

	// Requester asks for 100 confirmations, we have a delay of 0,
	// so we should wait for max(nodeMinConfs, requestedConfs + requestedConfsDelay) = 100 confirmations
	j, err = ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 0,
	}).Toml())
	require.NoError(t, err)
	listener.job = j
	confirmedAt = listener.getConfirmedAt(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}, uint32(nodeMinConfs))
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

func TestListener_ShouldProcessSub_NotEnoughBalance(t *testing.T) {
	mockAggregator := &vrf_mocks.AggregatorV3Interface{}
	mockAggregator.On("LatestRoundData", mock.Anything).Return(
		aggregator_v3_interface.LatestRoundData{
			Answer: decimal.RequireFromString("9821673525377230000").BigInt(),
		},
		nil,
	)
	defer mockAggregator.AssertExpectations(t)

	cfg := &vrf_mocks.Config{}
	cfg.On("KeySpecificMaxGasPriceWei", mock.Anything).Return(
		assets.GWei(200),
	)
	defer cfg.AssertExpectations(t)

	lsn := &listenerV2{
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{
					ethkey.EIP55Address("0x7Bf4E7069d96eEce4f48F50A9768f8615A8cD6D8"),
				},
			},
		},
		aggregator: mockAggregator,
		l:          logger.TestLogger(t),
		chainID:    big.NewInt(1337),
		cfg:        cfg,
	}
	subID := uint64(1)
	sub := vrf_coordinator_v2.GetSubscription{
		Balance: assets.GWei(100),
	}
	shouldProcess := lsn.shouldProcessSub(subID, sub, []pendingRequest{
		{
			req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
				CallbackGasLimit: 1e6,
				RequestId:        big.NewInt(1),
			},
		},
	})
	assert.False(t, shouldProcess) // estimated fee: 24435754189944131 juels, 100 GJuels not enough
}

func TestListener_ShouldProcessSub_EnoughBalance(t *testing.T) {
	mockAggregator := &vrf_mocks.AggregatorV3Interface{}
	mockAggregator.On("LatestRoundData", mock.Anything).Return(
		aggregator_v3_interface.LatestRoundData{
			Answer: decimal.RequireFromString("9821673525377230000").BigInt(),
		},
		nil,
	)
	defer mockAggregator.AssertExpectations(t)

	cfg := &vrf_mocks.Config{}
	cfg.On("KeySpecificMaxGasPriceWei", mock.Anything).Return(
		assets.GWei(200),
	)
	defer cfg.AssertExpectations(t)

	lsn := &listenerV2{
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{
					ethkey.EIP55Address("0x7Bf4E7069d96eEce4f48F50A9768f8615A8cD6D8"),
				},
			},
		},
		aggregator: mockAggregator,
		l:          logger.TestLogger(t),
		chainID:    big.NewInt(1337),
		cfg:        cfg,
	}
	subID := uint64(1)
	sub := vrf_coordinator_v2.GetSubscription{
		Balance: assets.Ether(1),
	}
	shouldProcess := lsn.shouldProcessSub(subID, sub, []pendingRequest{
		{
			req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
				CallbackGasLimit: 1e6,
				RequestId:        big.NewInt(1),
			},
		},
	})
	assert.True(t, shouldProcess) // estimated fee: 24435754189944131 juels, 1 LINK is enough.
}

func TestListener_ShouldProcessSub_NoLinkEthPrice(t *testing.T) {
	mockAggregator := &vrf_mocks.AggregatorV3Interface{}
	mockAggregator.On("LatestRoundData", mock.Anything).Return(
		aggregator_v3_interface.LatestRoundData{},
		errors.New("aggregator error"),
	)
	defer mockAggregator.AssertExpectations(t)

	cfg := &vrf_mocks.Config{}
	cfg.On("KeySpecificMaxGasPriceWei", mock.Anything).Return(
		assets.GWei(200),
	)
	defer cfg.AssertExpectations(t)

	lsn := &listenerV2{
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{
					ethkey.EIP55Address("0x7Bf4E7069d96eEce4f48F50A9768f8615A8cD6D8"),
				},
			},
		},
		aggregator: mockAggregator,
		l:          logger.TestLogger(t),
		chainID:    big.NewInt(1337),
		cfg:        cfg,
	}
	subID := uint64(1)
	sub := vrf_coordinator_v2.GetSubscription{
		Balance: assets.Ether(1),
	}
	shouldProcess := lsn.shouldProcessSub(subID, sub, []pendingRequest{
		{
			req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
				CallbackGasLimit: 1e6,
				RequestId:        big.NewInt(1),
			},
		},
	})
	assert.True(t, shouldProcess) // no fee available, try to process it.
}

func TestListener_ShouldProcessSub_NoFromAddresses(t *testing.T) {
	mockAggregator := &vrf_mocks.AggregatorV3Interface{}
	defer mockAggregator.AssertExpectations(t)

	cfg := &vrf_mocks.Config{}
	defer cfg.AssertExpectations(t)

	lsn := &listenerV2{
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{},
			},
		},
		aggregator: mockAggregator,
		l:          logger.TestLogger(t),
		chainID:    big.NewInt(1337),
		cfg:        cfg,
	}
	subID := uint64(1)
	sub := vrf_coordinator_v2.GetSubscription{
		Balance: assets.Ether(1),
	}
	shouldProcess := lsn.shouldProcessSub(subID, sub, []pendingRequest{
		{
			req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
				CallbackGasLimit: 1e6,
				RequestId:        big.NewInt(1),
			},
		},
	})
	assert.True(t, shouldProcess) // no addresses, but try to process it.
}

func TestListener_ProcessPendingVRFRequests_SubscriptionNotFound(t *testing.T) {
	// given
	reqs := []pendingRequest{
		{
			confirmedAtBlock: 100,
			req:              &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{},
			lb:               log.NewLogBroadcast(types.Log{}, *big.NewInt(1), nil),
		},
	}
	coordinatorMock := &vrf_mocks.VRFCoordinatorV2Interface{}

	coordinatorMock.On("GetSubscription", mock.Anything, mock.Anything).Return(vrf_coordinator_v2.GetSubscription{}, executionRevertedError{})
	defer coordinatorMock.AssertExpectations(t)

	broadcasterMock := &mocks.Broadcaster{}
	broadcasterMock.On("MarkConsumed", mock.Anything).Return(nil)
	defer broadcasterMock.AssertExpectations(t)

	lsn := &listenerV2{
		coordinator:    coordinatorMock,
		logBroadcaster: broadcasterMock,
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{},
			},
		},
		l:                  logger.NullLogger,
		blockNumberToReqID: pairing.New(),
		reqs:               reqs,
		latestHeadNumber:   100,
	}

	// when
	lsn.processPendingVRFRequests(context.Background())

	// then
	assert.Empty(t, lsn.reqs)
}

func TestListener_ProcessPendingVRFRequests_ProcessedLogsMarkedAsConsumed(t *testing.T) {
	// given
	reqs := []pendingRequest{
		{
			confirmedAtBlock: 100,
			req:              &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{},
			lb:               log.NewLogBroadcast(types.Log{}, *big.NewInt(1), nil),
		},
	}
	coordinatorMock := &vrf_mocks.VRFCoordinatorV2Interface{}

	coordinatorMock.On("GetSubscription", mock.Anything, mock.Anything).Return(vrf_coordinator_v2.GetSubscription{Balance: big.NewInt(1)}, nil)
	defer coordinatorMock.AssertExpectations(t)

	broadcasterMock := &mocks.Broadcaster{}
	broadcasterMock.On("MarkConsumed", mock.Anything).Return(nil)
	defer broadcasterMock.AssertExpectations(t)

	ethClientMock := &eth_client_mocks.Client{}
	ethClientMock.On("ChainID").Return(big.NewInt(1))
	defer ethClientMock.AssertExpectations(t)

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	q := pg.NewQ(db, lggr, &config{})

	lsn := &listenerV2{
		coordinator:    coordinatorMock,
		logBroadcaster: broadcasterMock,
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{},
			},
		},
		l:                  logger.NullLogger,
		blockNumberToReqID: pairing.New(),
		reqs:               reqs,
		q:                  q,
		latestHeadNumber:   100,
		ethClient:          ethClientMock,
	}

	// when
	lsn.processPendingVRFRequests(context.Background())

	// then
	assert.Empty(t, lsn.reqs)
}

func TestListener_ProcessPendingVRFRequests_LogNotMarkedAsConsumed_WhenNetworkError(t *testing.T) {
	// given
	reqs := []pendingRequest{
		{
			confirmedAtBlock: 100,
			req:              &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{},
			lb:               log.NewLogBroadcast(types.Log{}, *big.NewInt(1), nil),
		},
	}
	coordinatorMock := &vrf_mocks.VRFCoordinatorV2Interface{}

	coordinatorMock.On("GetSubscription", mock.Anything, mock.Anything).Return(vrf_coordinator_v2.GetSubscription{}, networkError{})
	defer coordinatorMock.AssertExpectations(t)

	broadcasterMock := &mocks.Broadcaster{}
	broadcasterMock.AssertNotCalled(t, "MarkConsumed", mock.Anything)
	defer broadcasterMock.AssertExpectations(t)

	lsn := &listenerV2{
		coordinator:    coordinatorMock,
		logBroadcaster: broadcasterMock,
		job: job.Job{
			VRFSpec: &job.VRFSpec{
				FromAddresses: []ethkey.EIP55Address{},
			},
		},
		l:                  logger.NullLogger,
		blockNumberToReqID: pairing.New(),
		reqs:               reqs,
		latestHeadNumber:   100,
	}

	// when
	lsn.processPendingVRFRequests(context.Background())

	// then
	assert.True(t, len(lsn.reqs) == 1)
}
