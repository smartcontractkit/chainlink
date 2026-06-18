package keystore_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func Test_OCR2KeyStore_Import_RejectsDuplicateWithoutDelete(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	master := keystore.ExposedNewMaster(t, db)
	require.NoError(t, master.Unlock(testutils.Context(t), cltest.Password))
	ks := master.OCR2()
	ctx := testutils.Context(t)

	chain := corekeys.SupportedChainTypes[0]
	key, err := ks.Create(ctx, chain)
	require.NoError(t, err)

	exportedJSON, err := ks.Export(key.ID(), cltest.Password)
	require.NoError(t, err)

	require.ErrorContains(t, ks.Add(ctx, key), "already exists")

	_, err = ks.Import(ctx, exportedJSON, cltest.Password)
	require.ErrorContains(t, err, "already exists")
}
