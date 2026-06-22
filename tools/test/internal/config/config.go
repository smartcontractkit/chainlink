// Package config holds the database/runtime configuration consumed by the
// tools/test database helpers.
package config

import "fmt"

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
	// DiagnoseMode is true when the harness is running testrig diagnose.
	DiagnoseMode bool
	// WorkerIndex is the 1-based diagnose worker slot when DiagnoseMode is set.
	WorkerIndex int
	// PackageSlug is a short, docker-safe token derived from the package patterns
	// under test (e.g. core_services).
	PackageSlug string
}

// PostgresContainerName returns the docker container name for an ephemeral
// Postgres instance. Non-diagnose runs use test_<slug>; diagnose workers use
// iteration_<n>_<slug>.
func (c *App) PostgresContainerName() string {
	if c == nil {
		return "test_pkgs"
	}
	slug := c.PackageSlug
	if slug == "" {
		slug = "pkgs"
	}
	if c.DiagnoseMode {
		return fmt.Sprintf("iteration_%d_%s", max(c.WorkerIndex, 1), slug)
	}
	return "test_" + slug
}
