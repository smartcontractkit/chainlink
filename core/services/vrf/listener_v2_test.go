package vrf

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

func addEthTx(t *testing.T, db *sqlx.DB, from common.Address, state bulletprooftxmanager.EthTxState, maxLink string, subID uint64) {
	_, err := db.Exec(`INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		$1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9, $10, $11, $12
		)
		RETURNING "eth_txes".*`,
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		state,
		bulletprooftxmanager.EthTxMeta{
			MaxLink: maxLink,
			SubID:   subID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil,
		false)
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *sqlx.DB, from common.Address, maxLink string, subID, nonce uint64) {
	_, err := db.Exec(`INSERT INTO eth_txes (nonce, broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		$1, NOW(), NULL, $2, $3, $4, $5, $6, 'confirmed', NOW(), $7, $8, $9, $10, $11, $12
		)
		RETURNING "eth_txes".*`,
		nonce,          // nonce
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		bulletprooftxmanager.EthTxMeta{
			MaxLink: maxLink,
			SubID:   subID,
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil,
		false)
	require.NoError(t, err)
}

// cannot import cltest, because of circular imports
type config struct{}

func (c *config) LogSQL() bool {
	return false
}

func TestMaybeSubtractReservedLink(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	q := pg.NewNewQ(db, lggr, &config{})
	ks := keystore.New(db, utils.FastScryptParams, lggr, &config{})
	require.NoError(t, ks.Unlock("blah"))
	chainID := uint64(1337)
	k, err := ks.Eth().Create(big.NewInt(int64(chainID)))
	require.NoError(t, err)

	subID := uint64(1)

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", subID)
	start, err := MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address.Address(), "10000", subID, 1)
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", subID)
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved link should not affect other subscribers prospective balance.
	otherSubID := uint64(2)
	require.NoError(t, err)
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", otherSubID)
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	anotherSubID := uint64(3)
	addEthTx(t, db, k2.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", anotherSubID)
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the link reserved across multiple keys,
	// i.e, gas lanes.
	addEthTx(t, db, k2.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", subID)
	start, err = MaybeSubtractReservedLink(lggr, q, k2.Address.Address(), big.NewInt(100_000), chainID, subID)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}
