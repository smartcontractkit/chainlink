package vrf

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func addEthTx(t *testing.T, db *sqlx.DB, from common.Address, state txmgr.EthTxState, maxLink string, subID uint64, reqTxHash common.Hash) {
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
			MaxLink:       &maxLink,
			SubID:         &subID,
			RequestTxHash: &reqTxHash,
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

func TestMaybeSubtractReservedLink(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := pgtest.NewQConfig(false)
	q := pg.NewQ(db, lggr, cfg)
	ks := keystore.New(db, utils.FastScryptParams, lggr, cfg)
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	k, err := ks.Eth().Create(big.NewInt(int64(chainID)))
	require.NoError(t, err)

	subID := uint64(1)
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address, txmgr.EthTxUnstarted, "10000", subID, reqTxHash)
	start, err := MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address, "10000", subID, 1)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTx(t, db, k.Address, txmgr.EthTxUnstarted, "10000", subID, reqTxHash)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved link should not affect other subscribers prospective balance.
	otherSubID := uint64(2)
	require.NoError(t, err)
	addEthTx(t, db, k.Address, txmgr.EthTxUnstarted, "10000", otherSubID, reqTxHash)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	anotherSubID := uint64(3)
	addEthTx(t, db, k2.Address, txmgr.EthTxUnstarted, "10000", anotherSubID, reqTxHash)
	start, err = MaybeSubtractReservedLink(q, big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the link reserved across multiple keys,
	// i.e, gas lanes.
	addEthTx(t, db, k2.Address, txmgr.EthTxUnstarted, "10000", subID, reqTxHash)
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
