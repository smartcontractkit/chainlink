package sessions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	sqlxTypes "github.com/smartcontractkit/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
)

func TestWebAuthnSessionStore(t *testing.T) {
	const key = "test-key"
	data := webauthn.SessionData{
		Challenge: "challenge-string",
		UserID:    []byte("test-user-id"),
		AllowedCredentialIDs: [][]byte{
			[]byte("test"),
			[]byte("foo"),
			[]byte("bar"),
		},
		UserVerification: protocol.UserVerificationRequirement("test-user-verification"),
	}
	s := NewWebAuthnSessionStore()

	val, ok := s.take(key)
	assert.Equal(t, "", val)
	require.False(t, ok)

	require.NoError(t, s.SaveWebauthnSession(key, &data))

	got, err := s.GetWebauthnSession(key)
	require.NoError(t, err)
	assert.Equal(t, data, got)

	val, ok = s.take(key)
	assert.Equal(t, "", val)
	require.False(t, ok)

	_, err = s.GetWebauthnSession(key)
	assert.ErrorContains(t, err, "assertion not in challenge store")

	user := mustRandomUser(t)
	cred := webauthn.Credential{
		ID:              []byte("test-id"),
		PublicKey:       []byte("test-key"),
		AttestationType: "test-attestation",
	}
	credj, err := json.Marshal(cred)
	require.NoError(t, err)

	token := WebAuthn{
		Email:         user.Email,
		PublicKeyData: sqlxTypes.JSONText(credj),
	}
	uwas := []WebAuthn{token}
	wcfg := WebAuthnConfiguration{RPID: "test-rpid", RPOrigin: "test-rporigin"}
	cc, err := s.BeginWebAuthnRegistration(user, uwas, wcfg)
	require.NoError(t, err)
	require.Equal(t, "Chainlink Operator", cc.Response.RelyingParty.CredentialEntity.Name)
	require.Equal(t, "test-rpid", cc.Response.RelyingParty.ID)
	require.Equal(t, user.Email, cc.Response.User.Name)
	require.Equal(t, user.Email, cc.Response.User.DisplayName)

	_, err = s.FinishWebAuthnRegistration(user, uwas, nil, wcfg)
	require.Error(t, err)
}

func mustRandomUser(t testing.TB) User {
	email := fmt.Sprintf("user-%v@chainlink.test", testutils.NewRandomPositiveInt64())
	r, err := NewUser(email, testutils.Password, UserRoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
