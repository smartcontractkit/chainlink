package pgtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func requireBenchDB(b *testing.B) {
	b.Helper()
	testutils.SkipShortDB(b)
	if string(env.DatabaseURL.Get()) == "" {
		b.Skip("CL_DATABASE_URL required")
	}
}

// BenchmarkNewSqlxDB measures txdb connection provisioning (open transaction + session
// timeouts). Each iteration closes the connection explicitly so the cost reflects a
// single test's DB setup, not accumulated transactions.
//
// Run: CL_DATABASE_URL='postgres://...' go test -bench=BenchmarkNewSqlxDB -benchmem ./core/internal/testutils/pgtest/
func BenchmarkNewSqlxDB(b *testing.B) {
	requireBenchDB(b)
	b.ReportAllocs()

	for b.Loop() {
		db := NewSqlxDB(b)
		_, err := db.Exec("SELECT 1")
		require.NoError(b, err)
		require.NoError(b, db.Close())
	}
}
