package fluxmonitorv2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
)

func TestKeyStore_EnabledKeysForChain(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	ks := fluxmonitorv2.NewKeyStore(ethKeyStore)

	key, err := ethKeyStore.Create(ctx, testutils.FixtureChainID)
	require.NoError(t, err)
	key2, err := ethKeyStore.Create(ctx, testutils.SimulatedChainID)
	require.NoError(t, err)

	keys, err := ks.EnabledKeysForChain(ctx, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key, keys[0])

	keys, err = ks.EnabledKeysForChain(ctx, testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key2, keys[0])
}

func TestKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, k0Address := cltest.MustInsertRandomKey(t, ethKeyStore)

	ks := fluxmonitorv2.NewKeyStore(ethKeyStore)

	// Gets the only address in the keystore
	addr, err := ks.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Equal(t, k0Address, addr)
}
