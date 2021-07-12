package bulletprooftxmanager_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SendEveryStrategy(t *testing.T) {
	t.Parallel()

	s := bulletprooftxmanager.SendEveryStrategy{}

	assert.Equal(t, uuid.NullUUID{}, s.Subject())

	n, err := s.PruneQueue(nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func Test_DropOldestStrategy_Subject(t *testing.T) {
	t.Parallel()

	subject := uuid.NewV4()
	s := bulletprooftxmanager.NewDropOldestStrategy(subject, 1)

	assert.True(t, s.Subject().Valid)
	assert.Equal(t, subject, s.Subject().UUID)
}

func Test_DropOldestStrategy_PruneQueue(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	subj1 := uuid.NewV4()
	subj2 := uuid.NewV4()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var n int64 = 0

	cltest.MustInsertFatalErrorEthTx(t, store, fromAddress)
	cltest.MustInsertInProgressEthTxWithAttempt(t, store, n, fromAddress)
	n++
	cltest.MustInsertConfirmedEthTxWithAttempt(t, store, n, 42, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, n, fromAddress)
	n++
	initialEtxs := []bulletprooftxmanager.EthTx{
		cltest.MustInsertUnstartedEthTx(t, store, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, store, fromAddress, subj2),
		cltest.MustInsertUnstartedEthTx(t, store, otherAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, store, fromAddress, subj1),
		cltest.MustInsertUnstartedEthTx(t, store, otherAddress, subj1),
	}

	db := store.DB

	t.Run("with queue size of 2, removes everything except the newest two transactions for the given subject, ignoring fromAddress", func(t *testing.T) {
		s := bulletprooftxmanager.NewDropOldestStrategy(subj1, 2)

		n, err := s.PruneQueue(db)
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)

		// Total inserted was 9. Minus the 2 oldest unstarted makes 7
		cltest.AssertCount(t, store, &bulletprooftxmanager.EthTx{}, 7)

		var etxs []bulletprooftxmanager.EthTx
		require.NoError(t, db.Raw(`SELECT * FROM eth_txes WHERE state = 'unstarted' ORDER BY id asc`).Scan(&etxs).Error)

		require.Len(t, etxs, 3)

		assert.Equal(t, initialEtxs[1].ID, etxs[0].ID)
		assert.Equal(t, initialEtxs[3].ID, etxs[1].ID)
		assert.Equal(t, initialEtxs[4].ID, etxs[2].ID)
	})
}
