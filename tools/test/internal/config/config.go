// Package config holds the database/runtime configuration consumed by the
// tools/test database helpers.
package config

// DefaultPostgresVersion is the Postgres major version used for ephemeral and
// persistent test databases when none is specified.
const DefaultPostgresVersion = "16"

// App is the configuration the database helpers need to stand up Postgres for a
// test run. It is populated from the testrig root flags (see internal/dbflags).
type App struct {
	// DatabaseURL, when set, points at an existing database to use instead of an
	// ephemeral container.
	DatabaseURL string
	// PostgresVersion is the Postgres major version for ephemeral containers.
	PostgresVersion string
	// RepoRoot is the repository root, used to run preparetest.
	RepoRoot string
	// AIOutput selects sparse output for agent tooling.
	AIOutput bool
	// ParallelIterations records the requested diagnose worker count (used only
	// to reject external databases with parallel runs).
	ParallelIterations int
}
