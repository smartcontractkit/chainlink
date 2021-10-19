package orm_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestORM_CreateExternalInitiator(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	exi2, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.Equal(t, `ERROR: duplicate key value violates unique constraint "external_initiators_name_key" (SQLSTATE 23505)`, store.CreateExternalInitiator(exi2).Error())
}

func TestORM_DeleteExternalInitiator(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	_, err = store.FindExternalInitiator(token)
	require.NoError(t, err)

	err = store.DeleteExternalInitiator(exi.Name)
	require.NoError(t, err)

	_, err = store.FindExternalInitiator(token)
	require.Error(t, err)

	require.NoError(t, store.CreateExternalInitiator(exi))
}

func TestORM_FindBridge(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, store.CreateBridgeType(&bt))

	cases := []struct {
		description string
		name        models.TaskType
		want        models.BridgeType
		errored     bool
	}{
		{"actual external adapter", bt.Name, bt, false},
		{"core adapter", "ethtx", models.BridgeType{}, true},
		{"non-existent adapter", "nonExistent", models.BridgeType{}, true},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			tt, err := store.FindBridge(test.name)
			tt.CreatedAt = test.want.CreatedAt
			tt.UpdatedAt = test.want.UpdatedAt
			assert.Equal(t, test.want, tt)
			assert.Equal(t, test.errored, err != nil)
		})
	}
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	user1 := cltest.MustNewUser(t, "test1@email1.net", "password1")
	user2 := cltest.MustNewUser(t, "test2@email2.net", "password2")
	user2.CreatedAt = time.Now().Add(-24 * time.Hour)

	require.NoError(t, store.SaveUser(&user1))
	require.NoError(t, store.SaveUser(&user2))

	actual, err := store.FindUser()
	require.NoError(t, err)
	assert.Equal(t, user1.Email, actual.Email)
	assert.Equal(t, user1.HashedPassword, actual.HashedPassword)
}

func TestORM_AuthorizedUserWithSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		sessionID       string
		sessionDuration time.Duration
		wantError       bool
		wantEmail       string
	}{
		{"authorized", "correctID", cltest.MustParseDuration(t, "3m"), false, "have@email"},
		{"expired", "correctID", cltest.MustParseDuration(t, "0m"), true, ""},
		{"incorrect", "wrong", cltest.MustParseDuration(t, "3m"), true, ""},
		{"empty", "", cltest.MustParseDuration(t, "3m"), true, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			user := cltest.MustNewUser(t, "have@email", "password")
			require.NoError(t, store.SaveUser(&user))

			prevSession := cltest.NewSession("correctID")
			prevSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "2m"))
			require.NoError(t, store.DB.Save(&prevSession).Error)

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := store.ORM.AuthorizedUserWithSession(test.sessionID, test.sessionDuration)
			assert.Equal(t, test.wantEmail, actual.Email)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var bumpedSession models.Session
				err = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
					return db.First(&bumpedSession, "ID = ?", prevSession.ID).Error
				})
				require.NoError(t, err)
				assert.Equal(t, expectedTime[0:13], utils.ISO8601UTC(bumpedSession.LastUsed)[0:13]) // only compare up to the hour
			}
		})
	}
}

func TestORM_DeleteUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, err := store.FindUser()
	require.NoError(t, err)

	err = store.DeleteUser()
	require.NoError(t, err)

	_, err = store.FindUser()
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	session := models.NewSession()
	require.NoError(t, store.DB.Save(&session).Error)

	err := store.DeleteUserSession(session.ID)
	require.NoError(t, err)

	_, err = store.FindUser()
	require.NoError(t, err)

	sessions, err := postgres.Sessions(store.DB, 0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	initial := cltest.MustRandomUser()
	require.NoError(t, store.SaveUser(&initial))

	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"correct", initial.Email, cltest.Password, true},
		{"incorrect email", "bogus@town.org", cltest.Password, false},
		{"incorrect pwd", initial.Email, "jamaicandundada", false},
		{"incorrect both", "dudus@coke.ja", "jamaicandundada", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionRequest := models.SessionRequest{
				Email:    test.email,
				Password: test.password,
			}

			sessionID, err := store.CreateSession(sessionRequest)
			if test.wantSession {
				require.NoError(t, err)
				assert.NotEmpty(t, sessionID)
			} else {
				require.Error(t, err)
				assert.Empty(t, sessionID)
			}
		})
	}
}

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewEthTxAttempt(t, tx2.ID)
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, store.DB.Create(&attempt).Error)

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = bulletprooftxmanager.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, store.DB.Save(&tx3).Error)

	count, err := store.CountOf(bulletprooftxmanager.EthTx{})
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := store.EthTransactionsWithAttempts(0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
	assert.Equal(t, int64(0), *txs[1].Nonce, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].EthTxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].EthTxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].EthTxAttempts[0].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].EthTxAttempts[1].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")

	txs, count, err = store.EthTransactionsWithAttempts(0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
}

func TestORM_UpdateBridgeType(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	firstBridge := &models.BridgeType{
		Name: "UniqueName",
		URL:  cltest.WebURL(t, "http:/oneurl.com"),
	}

	require.NoError(t, store.CreateBridgeType(firstBridge))

	updateBridge := &models.BridgeTypeRequest{
		URL: cltest.WebURL(t, "http:/updatedurl.com"),
	}

	require.NoError(t, store.UpdateBridgeType(firstBridge, updateBridge))

	foundbridge, err := store.FindBridge("UniqueName")
	require.NoError(t, err)
	require.Equal(t, updateBridge.URL, foundbridge.URL)
}
