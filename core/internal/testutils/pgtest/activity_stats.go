package pgtest

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/lib/pq" // postgres driver for diagnostics
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

// LogPostgresActivitySummary prints counts from pg_stat_activity for databases matching `chainlink_test%`.
// Use when investigating lock contention or idle-in-transaction buildup during test runs.
func LogPostgresActivitySummary(t testing.TB) {
	t.Helper()
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		t.Skip("CL_DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err)
	defer db.Close()

	q := `SELECT state, wait_event_type, wait_event, count(*)::int
FROM pg_stat_activity
WHERE datname LIKE 'chainlink_test%'
GROUP BY state, wait_event_type, wait_event
ORDER BY 4 DESC, 1, 2, 3`

	rows, err := db.Query(q)
	require.NoError(t, err)
	defer rows.Close()

	var b strings.Builder
	_, _ = b.WriteString("pg_stat_activity (chainlink_test%):\n")
	for rows.Next() {
		var state, waitType, waitEvent sql.NullString
		var n int
		require.NoError(t, rows.Scan(&state, &waitType, &waitEvent, &n))
		_, _ = fmt.Fprintf(&b, "  state=%s wait_event_type=%s wait_event=%s count=%d\n",
			nullStr(state), nullStr(waitType), nullStr(waitEvent), n)
	}
	require.NoError(t, rows.Err())
	t.Log(b.String())
}

func nullStr(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}
