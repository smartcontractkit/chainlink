package keystore_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
)

func Test_CSAKeyStore(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ks := cltest.NewKeyStore(t, store.DB).CSA()
	ks.Unlock(cltest.Password)

	t.Run("it can count keys", func(tt *testing.T) {
		count, err := ks.CountCSAKeys()
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
	})

	t.Run("it can list keys", func(tt *testing.T) {
		keys, err := ks.ListCSAKeys()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("it can create new keys", func(tt *testing.T) {
		_, err := ks.CreateCSAKey()
		require.NoError(t, err)

		count, err := ks.CountCSAKeys()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
	})

	t.Run("it won't allow more than one key", func(tt *testing.T) {
		_, err := ks.CreateCSAKey()
		require.Error(t, err)

		count, err := ks.CountCSAKeys()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
	})
}
