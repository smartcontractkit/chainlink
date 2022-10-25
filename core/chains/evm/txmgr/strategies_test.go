package txmgr_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func Test_SendEveryStrategy(t *testing.T) {
	t.Parallel()

	s := txmgr.SendEveryStrategy{}

	assert.Equal(t, uuid.NullUUID{}, s.Subject())

	n, err := s.PruneQueue(nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func Test_DropOldestStrategy_Subject(t *testing.T) {
	t.Parallel()

	subject := uuid.NewV4()
	s := txmgr.NewDropOldestStrategy(subject, 1, pg.DefaultQueryTimeout)

	assert.True(t, s.Subject().Valid)
	assert.Equal(t, subject, s.Subject().UUID)
}

func Test_DropOldestStrategy_PruneQueue(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	subj1 := uuid.NewV4()
	subj2 := uuid.NewV4()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var n int64

	cltest.MustInsertFatalErrorEthTx(t, borm, fromAddress)
	cltest.MustInsertInProgressEthTxWithAttempt(t, borm, n, fromAddress)
	n++
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, n, 42, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, n, fromAddress)
	initialEtxs := []txmgr.EthTx{
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress, subj2),
		cltest.MustInsertUnstartedEthTx(t, borm, otherAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, borm, otherAddress, subj1),
	}

	t.Run("with queue size of 2, removes everything except the newest two transactions for the given subject, ignoring fromAddress", func(t *testing.T) {
		s := txmgr.NewDropOldestStrategy(subj1, 2, pg.DefaultQueryTimeout)

		n, err := s.PruneQueue(db)
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)

		// Total inserted was 9. Minus the 2 oldest unstarted makes 7
		cltest.AssertCount(t, db, "eth_txes", 7)

		var etxs []txmgr.EthTx
		require.NoError(t, db.Select(&etxs, `SELECT * FROM eth_txes WHERE state = 'unstarted' ORDER BY id asc`))

		require.Len(t, etxs, 3)

		assert.Equal(t, initialEtxs[1].ID, etxs[0].ID)
		assert.Equal(t, initialEtxs[3].ID, etxs[1].ID)
		assert.Equal(t, initialEtxs[4].ID, etxs[2].ID)
	})
}
