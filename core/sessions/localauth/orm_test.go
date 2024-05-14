package localauth_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func setupORM(t *testing.T) (*sqlx.DB, sessions.AuthenticationProvider) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := localauth.NewORM(db, time.Minute, logger.TestLogger(t), &audit.AuditLoggerService{})

	return db, orm
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db, orm := setupORM(t)
	user1 := cltest.MustRandomUser(t)
	user2 := cltest.MustRandomUser(t)

	require.NoError(t, orm.CreateUser(ctx, &user1))
	require.NoError(t, orm.CreateUser(ctx, &user2))
	_, err := db.Exec("UPDATE users SET created_at = now() - interval '1 day' WHERE email = $1", user2.Email)
	require.NoError(t, err)

	actual, err := orm.FindUser(ctx, user1.Email)
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
		wantError       string
	}{
		{"authorized", "correctID", cltest.MustParseDuration(t, "3m"), ""},
		{"expired", "correctID", cltest.MustParseDuration(t, "0m"), sessions.ErrUserSessionExpired.Error()},
		{"incorrect", "wrong", cltest.MustParseDuration(t, "3m"), sessions.ErrUserSessionExpired.Error()},
		{"empty", "", cltest.MustParseDuration(t, "3m"), sessions.ErrEmptySessionID.Error()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			db := pgtest.NewSqlxDB(t)
			orm := localauth.NewORM(db, test.sessionDuration, logger.TestLogger(t), &audit.AuditLoggerService{})

			user := cltest.MustRandomUser(t)
			require.NoError(t, orm.CreateUser(ctx, &user))

			prevSession := cltest.NewSession("correctID")
			prevSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "2m"))
			_, err := db.Exec("INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, $3, now())", prevSession.ID, user.Email, prevSession.LastUsed)
			require.NoError(t, err)

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := orm.AuthorizedUserWithSession(ctx, test.sessionID)
			if test.wantError != "" {
				require.EqualError(t, err, test.wantError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, user.Email, actual.Email)
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
	ctx := testutils.Context(t)
	_, orm := setupORM(t)

	u := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &u))

	err := orm.DeleteUser(ctx, u.Email)
	require.NoError(t, err)

	_, err = orm.FindUser(ctx, u.Email)
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db, orm := setupORM(t)

	u := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &u))

	session := sessions.NewSession()
	_, err := db.Exec("INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, u.Email)
	require.NoError(t, err)

	err = orm.DeleteUserSession(ctx, session.ID)
	require.NoError(t, err)

	_, err = orm.FindUser(ctx, u.Email)
	require.NoError(t, err)

	sessions, err := orm.Sessions(ctx, 0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_DeleteUserCascade(t *testing.T) {
	ctx := testutils.Context(t)
	db, orm := setupORM(t)

	u := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &u))

	session := sessions.NewSession()
	_, err := db.Exec("INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, u.Email)
	require.NoError(t, err)

	err = orm.DeleteUser(ctx, u.Email)
	require.NoError(t, err)

	_, err = orm.FindUser(ctx, u.Email)
	require.Error(t, err)

	sessions, err := orm.Sessions(ctx, 0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(testutils.Context(t), &initial))

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

			sessionID, err := orm.CreateSession(testutils.Context(t), sessionRequest)
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
	ctx := testutils.Context(t)

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &initial))

	was, err := orm.GetUserWebAuthn(ctx, initial.Email)
	require.NoError(t, err)
	assert.Len(t, was, 0)

	cred := webauthn.Credential{
		ID:              []byte("test-id"),
		PublicKey:       []byte("test-key"),
		AttestationType: "test-attestation",
	}
	require.NoError(t, sessions.AddCredentialToUser(ctx, orm, initial.Email, &cred))

	was, err = orm.GetUserWebAuthn(ctx, initial.Email)
	require.NoError(t, err)
	require.NotEmpty(t, was)

	_, err = orm.CreateSession(ctx, sessions.SessionRequest{
		Email:    initial.Email,
		Password: cltest.Password,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "MFA Error")

	ss := sessions.NewWebAuthnSessionStore()
	_, err = orm.CreateSession(ctx, sessions.SessionRequest{
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

	_, err = orm.CreateSession(ctx, sessions.SessionRequest{
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
	require.NoError(t, err)
	_, err = orm.CreateSession(ctx, sessions.SessionRequest{
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
	ctx := testutils.Context(t)

	_, orm := setupORM(t)

	initial := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &initial))

	token, err := orm.CreateAndSetAuthToken(ctx, &initial)
	require.NoError(t, err)

	dbUser, err := orm.FindUser(ctx, initial.Email)
	require.NoError(t, err)

	hashedSecret, err := auth.HashedSecret(token, dbUser.TokenSalt.String)
	require.NoError(t, err)

	assert.NotNil(t, token)
	assert.NotNil(t, token.Secret)
	assert.NotEmpty(t, token.AccessKey)
	assert.Equal(t, dbUser.TokenKey.String, token.AccessKey)
	assert.Equal(t, dbUser.TokenHashedSecret.String, hashedSecret)

	require.NoError(t, orm.DeleteAuthToken(ctx, &initial))
	dbUser, err = orm.FindUser(ctx, initial.Email)
	require.NoError(t, err)
	assert.Empty(t, dbUser.TokenKey.ValueOrZero())
	assert.Empty(t, dbUser.TokenSalt.ValueOrZero())
	assert.Empty(t, dbUser.TokenHashedSecret.ValueOrZero())
}
