package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func TestSolanaKeyring_Sign_Verify(t *testing.T) {
	kr1, err := newSolanaKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newSolanaKeyring(cryptorand.Reader)
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

func TestSolanaKeyring_Marshalling(t *testing.T) {
	kr1, err := newSolanaKeyring(cryptorand.Reader)
	require.NoError(t, err)
	m, err := kr1.Marshal()
	require.NoError(t, err)
	kr2 := solanaKeyring{}
	err = kr2.Unmarshal(m)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(kr1.PublicKey(), kr2.PublicKey()))
	assert.True(t, bytes.Equal(kr1.privateKey().D.Bytes(), kr2.privateKey().D.Bytes()))
}
