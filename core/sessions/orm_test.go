package sessions_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setupORM(t *testing.T) (*sqlx.DB, sessions.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := sessions.NewORM(db, time.Minute, logger.TestLogger(t))

	return db, orm
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()

	db, orm := setupORM(t)
	user1 := cltest.MustNewUser(t, "test1@email1.net", "longlonglongpassword1")
	user2 := cltest.MustNewUser(t, "test2@email2.net", "longlonglongpassword2")

	require.NoError(t, orm.CreateUser(&user1))
	require.NoError(t, orm.CreateUser(&user2))
	_, err := db.Exec("UPDATE users SET created_at = now() - interval '1 day' WHERE email = $1", user2.Email)
	require.NoError(t, err)

	actual, err := orm.FindUser()
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
			db := pgtest.NewSqlxDB(t)
			orm := sessions.NewORM(db, test.sessionDuration, logger.TestLogger(t))

			user := cltest.MustNewUser(t, "have@email", testutils.Password)
			require.NoError(t, orm.CreateUser(&user))

			prevSession := cltest.NewSession("correctID")
			prevSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "2m"))
			_, err := db.Exec("INSERT INTO sessions (id, last_used, created_at) VALUES ($1, $2, now())", prevSession.ID, prevSession.LastUsed)
			require.NoError(t, err)

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := orm.AuthorizedUserWithSession(test.sessionID)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.wantEmail, actual.Email)
				var bumpedSession sessions.Session
				err = db.Get(&bumpedSession, "SELECT * FROM sessions WHERE ID = $1", prevSession.ID)
				require.NoError(t, err)
				assert.Equal(t, expectedTime[0:13], utils.ISO8601UTC(bumpedSession.LastUsed)[0:13]) // only compare up to the hour
			}
		})
	}
}

func TestORM_DeleteUser(t *testing.T) {
	t.Parallel()
	_, orm := setupORM(t)

	_, err := orm.FindUser()
	require.NoError(t, err)

	err = orm.DeleteUser()
	require.NoError(t, err)

	_, err = orm.FindUser()
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()

	db, orm := setupORM(t)

	session := sessions.NewSession()
	_, err := db.Exec("INSERT INTO sessions (id, last_used, created_at) VALUES ($1, now(), now())", session.ID)
	require.NoError(t, err)

	err = orm.DeleteUserSession(session.ID)
	require.NoError(t, err)

	_, err = orm.FindUser()
	require.NoError(t, err)

	sessions, err := orm.Sessions(0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(&initial))

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
			sessionRequest := sessions.SessionRequest{
				Email:    test.email,
				Password: test.password,
			}

			sessionID, err := orm.CreateSession(sessionRequest)
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

func TestORM_WebAuthn(t *testing.T) {
	t.Parallel()

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(&initial))

	was, err := orm.GetUserWebAuthn(initial.Email)
	require.NoError(t, err)
	assert.Len(t, was, 0)

	cred := webauthn.Credential{
		ID:              []byte("test-id"),
		PublicKey:       []byte("test-key"),
		AttestationType: "test-attestation",
	}
	require.NoError(t, sessions.AddCredentialToUser(orm, initial.Email, &cred))

	was, err = orm.GetUserWebAuthn(initial.Email)
	require.NoError(t, err)
	require.NotEmpty(t, was)

	_, err = orm.CreateSession(sessions.SessionRequest{
		Email:    initial.Email,
		Password: cltest.Password,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "MFA Error")

	ss := sessions.NewWebAuthnSessionStore()
	_, err = orm.CreateSession(sessions.SessionRequest{
		Email:    initial.Email,
		Password: cltest.Password,
		WebAuthnConfig: sessions.WebAuthnConfiguration{
			RPID:     "test-rpid",
			RPOrigin: "test-rporigin",
		},
		SessionStore: ss,
	})
	require.Error(t, err)
	var ca protocol.CredentialAssertion
	require.NoError(t, json.Unmarshal([]byte(err.Error()), &ca))
	require.Equal(t, "test-rpid", ca.Response.RelyingPartyID)

	_, err = orm.CreateSession(sessions.SessionRequest{
		Email:    initial.Email,
		Password: cltest.Password,
		WebAuthnConfig: sessions.WebAuthnConfiguration{
			RPID:     "test-rpid",
			RPOrigin: "test-rporigin",
		},
		SessionStore: ss,
		WebAuthnData: "invalid-format",
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "MFA Error")

	challengeResp, err := json.Marshal(protocol.CredentialAssertionResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				ID:   "test-id",
				Type: "test-type",
			},
		},
	})
	_, err = orm.CreateSession(sessions.SessionRequest{
		Email:    initial.Email,
		Password: cltest.Password,
		WebAuthnConfig: sessions.WebAuthnConfiguration{
			RPID:     "test-rpid",
			RPOrigin: "test-rporigin",
		},
		WebAuthnData: string(challengeResp),
		SessionStore: ss,
	})
	require.Error(t, err)
}

func TestOrm_GenerateAuthToken(t *testing.T) {
	t.Parallel()

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(&initial))

	token, err := orm.CreateAndSetAuthToken(&initial)
	require.NoError(t, err)

	dbUser, err := orm.FindUser()
	require.NoError(t, err)

	hashedSecret, err := auth.HashedSecret(token, dbUser.TokenSalt.String)
	require.NoError(t, err)

	assert.NotNil(t, token)
	assert.NotNil(t, token.Secret)
	assert.NotEmpty(t, token.AccessKey)
	assert.Equal(t, dbUser.TokenKey.String, token.AccessKey)
	assert.Equal(t, dbUser.TokenHashedSecret.String, hashedSecret)

	require.NoError(t, orm.DeleteAuthToken(&initial))
	dbUser, err = orm.FindUser()
	require.NoError(t, err)
	assert.Empty(t, dbUser.TokenKey.ValueOrZero())
	assert.Empty(t, dbUser.TokenSalt.ValueOrZero())
	assert.Empty(t, dbUser.TokenHashedSecret.ValueOrZero())
}
