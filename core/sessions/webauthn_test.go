package sessions

import (
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	got, err = s.GetWebauthnSession(key)
	assert.ErrorContains(t, err, "assertion not in challenge store")
}
