package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func TestEVMKeyring_SignVerify(t *testing.T) {
	kr1, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newEVMKeyring(cryptorand.Reader)
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

func TestEVMKeyring_Marshalling(t *testing.T) {
	kr1, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)

	m, err := kr1.Marshal()
	require.NoError(t, err)

	kr2 := evmKeyring{}
	err = kr2.Unmarshal(m)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(kr1.PublicKey(), kr2.PublicKey()))
	assert.True(t, bytes.Equal(kr1.privateKey.D.Bytes(), kr2.privateKey.D.Bytes()))

	// Invalid seed size should error
	assert.Error(t, kr2.Unmarshal([]byte{0x01}))
}
