package vaultcrypto

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

var label = [32]byte{31: 1}

type recipient struct {
	pub  *[32]byte
	priv *[32]byte
}

func newRecipient(t *testing.T) recipient {
	pub, priv, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	return recipient{pub: pub, priv: priv}
}

func (r recipient) hex() string { return hex.EncodeToString(r.pub[:]) }

func (r recipient) decrypter() decrypter {
	return decryptFunc(func(b []byte) ([]byte, error) {
		out, ok := box.OpenAnonymous(nil, b, r.pub, r.priv)
		if !ok {
			return nil, assert.AnError
		}
		return out, nil
	})
}

type decryptFunc func(encrypted []byte) ([]byte, error)

func (f decryptFunc) Decrypt(encrypted []byte) ([]byte, error) { return f(encrypted) }

func setup(t *testing.T, k, n int, secret string) (*tdh2easy.PublicKey, []*tdh2easy.PrivateShare, []byte) {
	_, pk, shares, err := tdh2easy.GenerateKeys(k, n)
	require.NoError(t, err)
	ct, err := tdh2easy.EncryptWithLabel(pk, []byte(secret), label)
	require.NoError(t, err)
	ctb, err := ct.Marshal()
	require.NoError(t, err)
	return pk, shares, ctb
}

func TestGenerateShareEncryptAndDecrypt(t *testing.T) {
	const k, n = 2, 4
	pk, keyShares, ctb := setup(t, k, n, "my secret")
	r := newRecipient(t)

	var encrypted [][]byte
	for _, ks := range keyShares[:k] {
		share, err := GeneratePlaintextShare(pk, ks, ctb, label)
		require.NoError(t, err)
		enc, err := EncryptShareBinary(share, r.hex())
		require.NoError(t, err)
		encrypted = append(encrypted, enc)
	}

	plaintext, rejected, err := DecryptSecret(ctb, pk, k, encrypted, r.decrypter())
	require.NoError(t, err)
	assert.Empty(t, rejected)
	assert.Equal(t, "my secret", string(plaintext))
}

func TestGeneratePlaintextShare_Errors(t *testing.T) {
	pk, keyShares, ctb := setup(t, 2, 4, "my secret")

	_, err := GeneratePlaintextShare(pk, keyShares[0], ctb, [32]byte{31: 2})
	require.ErrorContains(t, err, "failed to verify label on secret")

	_, err = GeneratePlaintextShare(pk, keyShares[0], []byte("garbage"), label)
	require.ErrorContains(t, err, "failed to unmarshal ciphertext")

	_, err = GeneratePlaintextShare(nil, keyShares[0], ctb, label)
	require.ErrorContains(t, err, "public key is not set")

	_, err = GeneratePlaintextShare(pk, nil, ctb, label)
	require.ErrorContains(t, err, "private key share is not set")
}

func TestEncryptShareBinary_InvalidKey(t *testing.T) {
	_, err := EncryptShareBinary([]byte("share"), "zz")
	require.ErrorIs(t, err, ErrInvalidRecipientKey)

	_, err = EncryptShareBinary([]byte("share"), "abcd")
	require.ErrorIs(t, err, ErrInvalidRecipientKey)
}

func TestInvalidRecipientKeyError_MatchesItself(t *testing.T) {
	_, err := EncryptShareBinary([]byte("share"), "zz")
	require.ErrorIs(t, err, ErrInvalidRecipientKey)

	other := invalidRecipientKeyError{"some other detail"}
	assert.ErrorIs(t, err, other)
	assert.ErrorIs(t, other, err)
}

func TestDecryptSecret_SkipsBadShares(t *testing.T) {
	const k, n = 2, 4
	pk, keyShares, ctb := setup(t, k, n, "my secret")
	r := newRecipient(t)
	other := newRecipient(t)

	encryptFor := func(ks *tdh2easy.PrivateShare, to recipient) []byte {
		share, err := GeneratePlaintextShare(pk, ks, ctb, label)
		require.NoError(t, err)
		enc, err := EncryptShareBinary(share, to.hex())
		require.NoError(t, err)
		return enc
	}

	// A share addressed to someone else and a garbage share are rejected.
	garbage, err := EncryptShareBinary([]byte("not a share"), r.hex())
	require.NoError(t, err)
	bad := [][]byte{encryptFor(keyShares[0], other), garbage}

	_, rejected, err := DecryptSecret(ctb, pk, k, append(bad, encryptFor(keyShares[0], r)), r.decrypter())
	require.ErrorContains(t, err, "not enough decryption shares")
	assert.Len(t, rejected, 2)

	plaintext, rejected, err := DecryptSecret(ctb, pk, k, append(bad, encryptFor(keyShares[0], r), encryptFor(keyShares[1], r)), r.decrypter())
	require.NoError(t, err)
	assert.Len(t, rejected, 2)
	assert.Equal(t, "my secret", string(plaintext))
}

func TestDecryptSecret_DuplicateShares(t *testing.T) {
	const k, n = 2, 4
	pk, keyShares, ctb := setup(t, k, n, "my secret")
	r := newRecipient(t)

	encryptFor := func(ks *tdh2easy.PrivateShare, to recipient) []byte {
		share, err := GeneratePlaintextShare(pk, ks, ctb, label)
		require.NoError(t, err)
		enc, err := EncryptShareBinary(share, to.hex())
		require.NoError(t, err)
		return enc
	}

	// The same node's share encrypted twice: the sealed boxes differ, but both
	// decrypt to a share with the same node index.
	share0, err := GeneratePlaintextShare(pk, keyShares[0], ctb, label)
	require.NoError(t, err)
	share0Enc1, err := EncryptShareBinary(share0, r.hex())
	require.NoError(t, err)
	share0Enc2, err := EncryptShareBinary(share0, r.hex())
	require.NoError(t, err)

	// An exact-duplicate share and a distinct share still decrypt at threshold.
	plaintext, rejected, err := DecryptSecret(ctb, pk, k, [][]byte{share0Enc1, share0Enc2, encryptFor(keyShares[1], r)}, r.decrypter())
	require.NoError(t, err)
	assert.Equal(t, "my secret", string(plaintext))
	require.Len(t, rejected, 1)
	assert.ErrorContains(t, rejected[0], "duplicate share for node index 0")

	// A re-randomized share from the same node (fresh proof scalars, same u_i
	// and index) is also dropped; byte-level dedup would miss this one.
	plaintext, rejected, err = DecryptSecret(ctb, pk, k, [][]byte{share0Enc1, encryptFor(keyShares[0], r), encryptFor(keyShares[1], r)}, r.decrypter())
	require.NoError(t, err)
	assert.Equal(t, "my secret", string(plaintext))
	require.Len(t, rejected, 1)
	assert.ErrorContains(t, rejected[0], "duplicate share for node index 0")

	// Duplicates never inflate the distinct-share count past the threshold.
	_, rejected, err = DecryptSecret(ctb, pk, k, [][]byte{share0Enc1, share0Enc2}, r.decrypter())
	require.ErrorContains(t, err, "not enough decryption shares")
	require.Len(t, rejected, 1)
	assert.ErrorContains(t, rejected[0], "duplicate share for node index 0")
}

func TestSharesForKey(t *testing.T) {
	entries := []*vaultcommon.EncryptedShares{
		{EncryptionKey: "a", BinaryShares: [][]byte{{1}, {2}}},
		{EncryptionKey: "b", BinaryShares: [][]byte{{3}}},
		{EncryptionKey: "a", Shares: []string{"04"}},
		{EncryptionKey: "a", Shares: []string{"05"}, BinaryShares: [][]byte{{6}}},
	}

	got, err := SharesForKey(entries, "a")
	require.NoError(t, err)
	assert.Equal(t, [][]byte{{1}, {2}, {4}, {6}}, got)

	got, err = SharesForKey(entries, "c")
	require.NoError(t, err)
	assert.Empty(t, got)

	_, err = SharesForKey([]*vaultcommon.EncryptedShares{{EncryptionKey: "a", Shares: []string{"zz"}}}, "a")
	require.Error(t, err)
}
