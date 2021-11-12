package vrf

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var nonce int8

// Non-threadsafe nonce just for testing.
func getNonce() int8 {
	nonce++
	return nonce
}

func addEthTx(t *testing.T, db *gorm.DB, from common.Address, state bulletprooftxmanager.EthTxState, maxLink string, subID uint64) {
	err := db.Exec(`INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		?,?,?,?,?,?,NOW(),?,?,?,?,?,?
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
		false).Error
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *gorm.DB, from common.Address, maxLink string, subID uint64) {
	err := db.Exec(`INSERT INTO eth_txes (nonce, broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		?, NOW(), NULL, ?,?,?,?,?,'confirmed',NOW(),?,?,?,?,?,?
		)
		RETURNING "eth_txes".*`,
		getNonce(),     // nonce
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
		false).Error
	require.NoError(t, err)
}

func TestMaybeSubtractReservedLink(t *testing.T) {
	db := pgtest.NewGormDB(t)
	lggr := logger.TestLogger(t)
	ks := keystore.New(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock("blah"))
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	subID := uint64(1)

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", subID)
	start, err := MaybeSubtractReservedLink(lggr, db, k.Address.Address(), big.NewInt(100_000), subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address.Address(), "10000", subID)
	start, err = MaybeSubtractReservedLink(lggr, db, k.Address.Address(), big.NewInt(100_000), subID)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// Another unstarted should
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000", subID)
	start, err = MaybeSubtractReservedLink(lggr, db, k.Address.Address(), big.NewInt(100_000), subID)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's balance should not affect other subscribers balances
	otherSubID := uint64(2)
	require.NoError(t, err)
	addConfirmedEthTx(t, db, k.Address.Address(), "20000", otherSubID)
	start, err = MaybeSubtractReservedLink(lggr, db, k.Address.Address(), big.NewInt(100_000), subID)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())
}
