package txmgr_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func Test_SendEveryStrategy(t *testing.T) {
	t.Parallel()

	s := txmgr.SendEveryStrategy{}

	assert.Equal(t, uuid.NullUUID{}, s.Subject())

	n, err := s.PruneQueue(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func Test_DropOldestStrategy_Subject(t *testing.T) {
	t.Parallel()
	cfg := configtest.NewGeneralConfig(t, nil)

	subject := uuid.NewV4()
	s := txmgr.NewDropOldestStrategy(subject, 1, cfg.DatabaseDefaultQueryTimeout())

	assert.True(t, s.Subject().Valid)
	assert.Equal(t, subject, s.Subject().UUID)
}

func Test_DropOldestStrategy_PruneQueue(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTxStore(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	subj1 := uuid.NewV4()
	subj2 := uuid.NewV4()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var n int64

	cltest.MustInsertFatalErrorEthTx(t, txStore, fromAddress)
	cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, n, fromAddress)
	n++
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, n, 42, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, n, fromAddress)
	initialEtxs := []txmgr.EvmTx{
		cltest.MustInsertUnstartedEthTx(t, txStore, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, txStore, fromAddress, subj2),
		cltest.MustInsertUnstartedEthTx(t, txStore, otherAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, txStore, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, txStore, otherAddress, subj1),
	}

	t.Run("with queue size of 2, removes everything except the newest two transactions for the given subject, ignoring fromAddress", func(t *testing.T) {
		s := txmgr.NewDropOldestStrategy(subj1, 2, cfg.DatabaseDefaultQueryTimeout())

		n, err := s.PruneQueue(txStore, pg.WithQueryer(db))
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)

		// Total inserted was 9. Minus the 2 oldest unstarted makes 7
		cltest.AssertCount(t, db, "eth_txes", 7)

		var dbEtxs []txmgr.DbEthTx
		require.NoError(t, db.Select(&dbEtxs, `SELECT * FROM eth_txes WHERE state = 'unstarted' ORDER BY id asc`))

		require.Len(t, dbEtxs, 3)

		assert.Equal(t, initialEtxs[1].ID, dbEtxs[0].ID)
		assert.Equal(t, initialEtxs[3].ID, dbEtxs[1].ID)
		assert.Equal(t, initialEtxs[4].ID, dbEtxs[2].ID)
	})
}
