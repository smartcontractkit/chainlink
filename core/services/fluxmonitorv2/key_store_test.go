package fluxmonitorv2_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/stretchr/testify/require"
)

func TestKeyStore_SendingKeys(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	ks := fluxmonitorv2.NewKeyStore(s.KeyStore)

	s.KeyStore.Unlock(cltest.Password)
	key, err := s.KeyStore.CreateNewKey()
	require.NoError(t, err)

	keys, err := ks.SendingKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key, keys[0])
}

func TestKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	cltest.MustAddRandomKeyToKeystore(t, s, 0, true)
	_, k0Address := cltest.MustAddRandomKeyToKeystore(t, s, 0)

	ks := fluxmonitorv2.NewKeyStore(s.KeyStore)

	// Gets the only address in the keystore
	addr, err := ks.GetRoundRobinAddress()
	require.NoError(t, err)
	require.Equal(t, k0Address, addr)
}
