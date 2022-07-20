package ocrkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRKey_ConfigPublicKey(t *testing.T) {
	k := MustNewV2XXXTestingOnly(big.NewInt(1))

	t.Run("fails to unmarshal invalid JSON", func(t *testing.T) {
		pk := ConfigPublicKey(k.PublicKeyConfig())

		err := pk.UnmarshalJSON([]byte(""))

		assert.Error(t, err)
	})

	t.Run("returns serialized instance value", func(t *testing.T) {
		pk := ConfigPublicKey(k.PublicKeyConfig())

		v, err := pk.Value()
		require.NoError(t, err)

		assert.NotEmpty(t, v)
	})

	t.Run("updates current instance by scanning another instance", func(t *testing.T) {
		pk := ConfigPublicKey(k.PublicKeyConfig())

		k2 := MustNewV2XXXTestingOnly(big.NewInt(1))
		pk2 := ConfigPublicKey(k2.PublicKeyConfig())

		err := pk.Scan(pk2[:])
		require.NoError(t, err)

		assert.Equal(t, pk2.Raw(), pk.Raw())
	})
}
