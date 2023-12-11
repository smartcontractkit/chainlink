package p2pkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestP2PKeys_KeyStruct(t *testing.T) {
	_, pk, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	k := Key{PrivKey: pk}

	t.Run("converts into V2 key", func(t *testing.T) {
		k2 := k.ToV2()

		assert.Equal(t, k.PrivKey, k2.PrivKey)
		assert.Equal(t, k.PeerID(), k2.peerID)
	})

	t.Run("returns PeerID", func(t *testing.T) {
		pid, err := k.GetPeerID()
		require.NoError(t, err)
		pid2 := k.PeerID()

		assert.Equal(t, pid, pid2)
	})
}

func TestP2PKeys_PublicKeyBytes(t *testing.T) {
	pk, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pkb := PublicKeyBytes(pk)
	assert.Equal(t, hex.EncodeToString(pkb), pkb.String())

	b, err := pkb.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, b)

	err = pkb.UnmarshalJSON(b)
	assert.NoError(t, err)

	err = pkb.UnmarshalJSON([]byte(""))
	assert.Error(t, err)

	err = pkb.Scan([]byte(pk))
	assert.NoError(t, err)

	err = pkb.Scan("invalid-type")
	assert.Error(t, err)

	sv, err := pkb.Value()
	assert.NoError(t, err)
	assert.NotEmpty(t, sv)
}

func TestP2PKeys_EncryptedP2PKey(t *testing.T) {
	_, privk, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	k := Key{PrivKey: privk}

	pubkr := k.PrivKey.Public().(ed25519.PublicKey)

	var marshalledPrivK []byte
	marshalledPrivK, err = MarshalPrivateKey(k.PrivKey)
	require.NoError(t, err)
	cryptoJSON, err := keystore.EncryptDataV3(marshalledPrivK, []byte(adulteratedPassword("password")), utils.FastScryptParams.N, utils.FastScryptParams.P)
	require.NoError(t, err)
	encryptedPrivKey, err := json.Marshal(&cryptoJSON)
	require.NoError(t, err)

	p2pk := EncryptedP2PKey{
		ID:               1,
		PeerID:           k.PeerID(),
		PubKey:           []byte(pubkr),
		EncryptedPrivKey: encryptedPrivKey,
	}

	t.Run("sets a different ID", func(t *testing.T) {
		err := p2pk.SetID("12")
		require.NoError(t, err)

		assert.Equal(t, int32(12), p2pk.ID)

		err = p2pk.SetID("invalid")
		assert.Error(t, err)
	})

	t.Run("decrypts key", func(t *testing.T) {
		k, err := p2pk.Decrypt("invalid-pass")
		assert.Empty(t, k)
		assert.Error(t, err)

		k, err = p2pk.Decrypt("password")
		require.NoError(t, err)
		assert.NotEmpty(t, k)
	})
}
