package ocr2key

import (
	cryptorand "crypto/rand"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStarknetKeyring_Sign_Verify(t *testing.T) {
	kr1, err := newStarknetKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newStarknetKeyring(cryptorand.Reader)
	require.NoError(t, err)
	ctx := ocrtypes.ReportContext{}
	report := ocrtypes.Report{
		97, 91, 43, 83, // observations_timestamp
		0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // observers
		2,                                                   // len
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 73, 150, 2, 210, // observation 1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 73, 150, 2, 211, // observation 2
		0, 0, 0, 0, 0, 0, 0, 0, 13, 224, 182, 179, 167, 100, 0, 0, // juels per luna (1 with 18 decimal places)
	}

	t.Run("can verify", func(t *testing.T) {
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		t.Log(len(sig))
		result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
		require.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		result := kr2.Verify(kr1.PublicKey(), ctx, report, []byte{0x01})
		require.False(t, result)

		longSig := [100]byte{}
		result = kr2.Verify(kr1.PublicKey(), ctx, report, longSig[:])
		require.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		result := kr2.Verify([]byte{0x01}, ctx, report, sig)
		require.False(t, result)
	})
}

func TestStarknetKeyring_Marshal(t *testing.T) {
	kr1, err := newStarknetKeyring(cryptorand.Reader)
	require.NoError(t, err)
	m, err := kr1.marshal()
	require.NoError(t, err)
	kr2 := starknetKeyring{}
	err = kr2.unmarshal(m)
	require.NoError(t, err)
	assert.True(t, kr1.privateKey.Equal(&kr2.privateKey))

	// Invalid seed size should error
	require.Error(t, kr2.unmarshal([]byte{0x01}))
}
