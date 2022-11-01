package terratest

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
)

func MustEnsureChain(t testing.TB, db *sqlx.DB, id string) {
	_, err := db.Exec("INSERT INTO terra_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW()) ON CONFLICT DO NOTHING;", id)
	require.NoError(t, err)
}

// RandomChainID returns a random chain id for testing. Use this instead of a constant to prevent DB collisions.
func RandomChainID() string {
	return fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
}
