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

func addEthTx(t *testing.T, db *sqlx.DB, from common.Address, state bulletprooftxmanager.EthTxState, maxLink string) {
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
		},
		uuid.NullUUID{},
		1337,
		0, // confs
		nil,
		false)
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, db *sqlx.DB, from common.Address, maxLink string) {
	_, err := db.Exec(`INSERT INTO eth_txes (nonce, broadcast_at, error, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, simulate)
		VALUES (
		10, NOW(), NULL, $1, $2, $3, $4, $5, 'confirmed', NOW(), $6, $7, $8, $9, $10, $11
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
	ks := keystore.New(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock("blah"))
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000")
	start, err := MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, db, k.Address.Address(), "10000")
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// Another unstarted should
	addEthTx(t, db, k.Address.Address(), bulletprooftxmanager.EthTxUnstarted, "10000")
	start, err = MaybeSubtractReservedLink(lggr, q, k.Address.Address(), big.NewInt(100000))
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())
}
