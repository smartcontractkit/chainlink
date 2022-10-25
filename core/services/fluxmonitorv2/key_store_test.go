package fluxmonitorv2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
)

func TestKeyStore_EnabledKeysForChain(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(true)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	ks := fluxmonitorv2.NewKeyStore(ethKeyStore)

	key, err := ethKeyStore.Create(testutils.FixtureChainID)
	require.NoError(t, err)
	key2, err := ethKeyStore.Create(testutils.SimulatedChainID)
	require.NoError(t, err)

	keys, err := ks.EnabledKeysForChain(testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key, keys[0])

	keys, err = ks.EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key2, keys[0])
}

func TestKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(true)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, k0Address := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	ks := fluxmonitorv2.NewKeyStore(ethKeyStore)

	// Gets the only address in the keystore
	addr, err := ks.GetRoundRobinAddress(testutils.FixtureChainID)
	require.NoError(t, err)
	require.Equal(t, k0Address, addr)
}
