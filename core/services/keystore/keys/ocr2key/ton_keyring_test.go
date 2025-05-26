package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func TestTONKeyring_Sign_Verify(t *testing.T) {
	kr1, err := newTONKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newTONKeyring(cryptorand.Reader)
	require.NoError(t, err)
	ctx := ocrtypes.ReportContext{}

	t.Run("can verify", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		t.Log(len(sig))
		result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
		assert.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		report := ocrtypes.Report{}
		result := kr2.Verify(kr1.PublicKey(), ctx, report, []byte{0x01})
		assert.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		result := kr2.Verify([]byte{0x01}, ctx, report, sig)
		assert.False(t, result)
	})
}

func TestTONKeyring_Sign3_Verify3(t *testing.T) {
	kr1, err := newTONKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newTONKeyring(cryptorand.Reader)
	require.NoError(t, err)
	digest := ocrtypes.ConfigDigest{}

	t.Run("can verify", func(t *testing.T) {
		report := ocrtypes.Report{}
		seqNr := uint64(1)
		sig, err := kr1.Sign3(digest, 1, report)
		require.NoError(t, err)
		t.Log(len(sig))

		result := kr2.Verify3(kr1.PublicKey(), digest, seqNr, report, sig)
		assert.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		report := ocrtypes.Report{}
		seqNr := uint64(1)
		result := kr2.Verify3(kr1.PublicKey(), digest, seqNr, report, []byte{0x01})
		assert.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		report := ocrtypes.Report{}
		seqNr := uint64(1)
		sig, err := kr1.Sign3(digest, 1, report)
		require.NoError(t, err)

		result := kr2.Verify3([]byte{0x01}, digest, seqNr, report, sig)
		assert.False(t, result)
	})
}

func TestTONKeyring_Marshalling(t *testing.T) {
	kr1, err := newAptosKeyring(cryptorand.Reader)
	require.NoError(t, err)
	m, err := kr1.Marshal()
	require.NoError(t, err)
	kr2 := aptosKeyring{}
	err = kr2.Unmarshal(m)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(kr1.pubKey, kr2.pubKey))
	assert.True(t, bytes.Equal(kr1.privKey(), kr2.privKey()))

	// Invalid seed size should error
	require.Error(t, kr2.Unmarshal([]byte{0x01}))
}
