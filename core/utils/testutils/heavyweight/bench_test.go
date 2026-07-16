package heavyweight

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

func requireBenchDB(b *testing.B) {
	b.Helper()
	if testing.Short() {
		b.Skip("skipping DB benchmark with -short")
	}
	if string(env.DatabaseURL.Get()) == "" {
		b.Skip("CL_DATABASE_URL required; run via make test or set CL_DATABASE_URL")
	}
}

// BenchmarkFullTestDBNoFixturesV2 measures isolated DB provisioning via
// CreateOrReplace (DROP + CREATE DATABASE … WITH TEMPLATE chainlink_test_pristine).
// Database drops are deferred to benchmark teardown via t.Cleanup; use a modest
// benchtime to avoid accumulating hundreds of databases locally.
//
// Run: CL_DATABASE_URL='postgres://...' go test -bench=BenchmarkFullTestDB -benchtime=5x -benchmem ./core/utils/testutils/heavyweight/
func BenchmarkFullTestDBNoFixturesV2(b *testing.B) {
	requireBenchDB(b)

	// Warm up template; cost excluded from timer.
	_, db := FullTestDBNoFixturesV2(b, nil)
	_, err := db.ExecContext(b.Context(), "SELECT 1")
	require.NoError(b, err)
	require.NoError(b, db.Close())

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, db := FullTestDBNoFixturesV2(b, nil)
		_, err := db.ExecContext(b.Context(), "SELECT 1")
		require.NoError(b, err)
		require.NoError(b, db.Close())
	}
}

// BenchmarkFullTestDBEmptyV2 measures empty DB provisioning (no migrations template).
func BenchmarkFullTestDBEmptyV2(b *testing.B) {
	requireBenchDB(b)
	b.ReportAllocs()

	// Warm up template; cost excluded from timer.
	_, db := FullTestDBEmptyV2(b, nil)
	_, err := db.ExecContext(b.Context(), "SELECT 1")
	require.NoError(b, err)
	require.NoError(b, db.Close())

	for b.Loop() {
		_, db := FullTestDBEmptyV2(b, nil)
		_, err := db.ExecContext(b.Context(), "SELECT 1")
		require.NoError(b, err)
		require.NoError(b, db.Close())
	}
}
