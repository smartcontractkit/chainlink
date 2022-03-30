package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

func TestTerraKeyRing_Sign_Verify(t *testing.T) {
	kr1, err := newTerraKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newTerraKeyring(cryptorand.Reader)
	require.NoError(t, err)
	ctx := ocrtypes.ReportContext{}

	t.Run("can verify", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		t.Log(len(sig))
		result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
		require.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		report := ocrtypes.Report{}
		result := kr2.Verify(kr1.PublicKey(), ctx, report, []byte{0x01})
		require.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		result := kr2.Verify([]byte{0x01}, ctx, report, sig)
		require.False(t, result)
	})
}

func TestTerraKeyRing_Marshalling(t *testing.T) {
	kr1, err := newTerraKeyring(cryptorand.Reader)
	require.NoError(t, err)
	m, err := kr1.marshal()
	require.NoError(t, err)
	kr2 := terraKeyring{}
	err = kr2.unmarshal(m)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(kr1.pubKey, kr2.pubKey))
	assert.True(t, bytes.Equal(kr1.privKey, kr2.privKey))

	// Invalid seed size should error
	require.Error(t, kr2.unmarshal([]byte{0x01}))
}
