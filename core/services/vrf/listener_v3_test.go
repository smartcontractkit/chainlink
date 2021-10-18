package vrf

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func addEthTx(t *testing.T, db *gorm.DB, from common.Address, state bulletprooftxmanager.EthTxState, maxLink string) {
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
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil,
		false).Error
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *gorm.DB, from common.Address, maxLink string) {
	err := db.Exec(`INSERT INTO eth_txes (nonce, broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		10, NOW(), NULL, ?,?,?,?,?,'confirmed',NOW(),?,?,?,?,?,?
		)
		RETURNING "eth_txes".*`,
		from,           // from
		from,           // to
		[]byte(`blah`), // payload
		0,              // value
		0,              // limit
		bulletprooftxmanager.EthTxMeta{
			MaxLink: maxLink,
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
	ks := keystore.New(db, utils.FastScryptParams, logger.Default)
	require.NoError(t, ks.Unlock("blah"))
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000")
	start, err := MaybeSubtractReservedLink(logger.Default, db, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address.Address(), "10000")
	start, err = MaybeSubtractReservedLink(logger.Default, db, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// Another unstarted should
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000")
	start, err = MaybeSubtractReservedLink(logger.Default, db, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())
}
