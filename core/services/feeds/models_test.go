package feeds_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PublicKey_String(t *testing.T) {
	t.Parallel()

	pubKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	pk := feeds.PublicKey(pubKey)
	expected := hex.EncodeToString(pubKey)

	assert.Equal(t, expected, pk.String())
}

func Test_PublicKey_MarshalJSON(t *testing.T) {
	t.Parallel()

	pubKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	hexKey := hex.EncodeToString(pubKey)

	pk := feeds.PublicKey(pubKey)
	actual, err := pk.MarshalJSON()
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf(`"%s"`, hexKey), string(actual))
}

func Test_PublicKey_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	pubKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	hexKey := hex.EncodeToString(pubKey)

	actual := &feeds.PublicKey{}
	err = actual.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, hexKey)))
	require.NoError(t, err)

	assert.Equal(t, feeds.PublicKey(pubKey), *actual)
}

func Test_PublicKey_Scan(t *testing.T) {
	pubKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	actual := &feeds.PublicKey{}

	// Error if not bytes
	err = actual.Scan("not bytes")
	assert.Error(t, err)

	// Nil
	err = actual.Scan(nil)
	require.NoError(t, err)
	nilPk := feeds.PublicKey(nil)
	assert.Equal(t, &nilPk, actual)

	// Bytes
	err = actual.Scan([]byte(pubKey))
	require.NoError(t, err)
	assert.Equal(t, feeds.PublicKey(pubKey), *actual)
}

func Test_PublicKey_Value(t *testing.T) {
	pubKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	pk := feeds.PublicKey(pubKey)
	dv, err := pk.Value()
	require.NoError(t, err)
	assert.Equal(t, []byte(pubKey), dv)
}
