package terratest

import (
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
)

func MustInsertChain(t testing.TB, db *sqlx.DB, chain *types.Chain) {
	query, args, e := db.BindNamed(`
INSERT INTO terra_chains (id, cfg, enabled, created_at, updated_at) VALUES (:id, :cfg, :enabled, NOW(), NOW()) RETURNING *;`, chain)
	require.NoError(t, e)
	err := db.Get(chain, query, args...)
	require.NoError(t, err)
}
