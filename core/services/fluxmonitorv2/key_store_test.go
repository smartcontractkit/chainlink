package fluxmonitorv2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
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

	chainID1 := testutils.NextEVMChainID()
	key, err := ethKeyStore.Create(ctx, chainID1)
	require.NoError(t, err)
	chainID2 := testutils.NextEVMChainID()
	key2, err := ethKeyStore.Create(ctx, chainID2)
	require.NoError(t, err)

	keys, err := ks.EnabledKeysForChain(ctx, chainID1)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key.ID(), keys[0].ID())
	require.Equal(t, key.Raw(), keys[0].Raw())
	require.EqualExportedValues(t, key, keys[0])

	keys, err = ks.EnabledKeysForChain(ctx, chainID2)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, key2.ID(), keys[0].ID())
	require.Equal(t, key2.Raw(), keys[0].Raw())
	require.EqualExportedValues(t, key2, keys[0])
}

func TestKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	chainID := testutils.NextEVMChainID()
	sqlChainID := sqlutil.New(chainID)
	_, k0Address := cltest.MustInsertRandomKey(t, ethKeyStore, *sqlChainID)

	ks := fluxmonitorv2.NewKeyStore(ethKeyStore)

	// Gets the only address in the keystore
	addr, err := ks.GetRoundRobinAddress(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, k0Address, addr)
}
