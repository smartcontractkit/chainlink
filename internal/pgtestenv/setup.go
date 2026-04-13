// Package pgtestenv starts an ephemeral Postgres (testcontainers) for tests when CL_DATABASE_URL is unset.
// It lives under internal/ (not core/internal/) so the module root test package can import it.
package pgtestenv

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq" // ping CL_DATABASE_URL
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/store"
)

// EnvDisableAutoContainers, when set to "true", disables starting a Postgres testcontainer
// when CL_DATABASE_URL is unset (tests will fail as before if the URL is required).
const EnvDisableAutoContainers = "CHAINLINK_PGTEST_DISABLE_AUTOCONTAINERS"

var (
	autoOnce sync.Once
	autoErr  error
)

type postgresTestResetConfig struct{ u url.URL }

func (c postgresTestResetConfig) URL() url.URL { return c.u }

func (c postgresTestResetConfig) DriverName() string { return pgcommon.DriverPostgres }

func (c postgresTestResetConfig) DefaultIdleInTxSessionTimeout() time.Duration { return time.Hour }

func (c postgresTestResetConfig) DefaultLockTimeout() time.Duration { return 15 * time.Second }

func (c postgresTestResetConfig) MaxOpenConns() int { return 100 }

func (c postgresTestResetConfig) MaxIdleConns() int { return 10 }

var _ store.Config = postgresTestResetConfig{}

// EnsureAutoPostgres starts an isolated Postgres (via testcontainers) when CL_DATABASE_URL is unset,
// runs migrations and [store.PrepareTestDB] fixtures once per process, and sets CL_DATABASE_URL.
// It is a no-op under [testing.Short], when CL_DATABASE_URL is set, or when [EnvDisableAutoContainers] is "true".
//
// If CL_DATABASE_URL is set but unreachable (stale port from an old container, stopped Postgres, etc.),
// it is cleared and a fresh testcontainer is started so `go test` does not fail with connection reset.
//
// Each OS process gets its own container (no cross-process reuse). Parallel `go test` runs many
// processes; sharing one DB and running DROP DATABASE from each caused connection resets.
func EnsureAutoPostgres(t testing.TB) {
	t.Helper()
	if testing.Short() {
		return
	}
	if os.Getenv(EnvDisableAutoContainers) == "true" {
		return
	}
	if raw := string(env.DatabaseURL.Get()); raw != "" {
		if pingPostgres(raw) {
			return
		}
		if u, err := url.Parse(raw); err == nil {
			t.Logf("pgtestenv: CL_DATABASE_URL was set but unreachable (%s); starting ephemeral Postgres instead", u.Redacted())
		} else {
			t.Logf("pgtestenv: CL_DATABASE_URL was set but unreachable; starting ephemeral Postgres instead")
		}
		if err := os.Unsetenv("CL_DATABASE_URL"); err != nil {
			t.Fatalf("pgtestenv: unset stale CL_DATABASE_URL: %v", err)
		}
	}

	autoOnce.Do(func() {
		autoErr = startEmbeddedPostgres(t)
	})
	require.NoError(t, autoErr, "pgtest autosetup: start test postgres")
}

func pingPostgres(dsn string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return false
	}
	defer db.Close()
	return db.PingContext(ctx) == nil
}

func startEmbeddedPostgres(t testing.TB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	lggr := logger.Test(t)

	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithCmdArgs(
			"-c", "synchronous_commit=off",
			"-c", "full_page_writes=off",
			"-c", "shared_buffers=128MB",
			"-c", "work_mem=32MB",
			"-c", "max_connections=500",
			"-c", "random_page_cost=1.1",
		),
	}

	ctr, err := postgres.Run(ctx, "postgres:16-alpine", opts...)
	if err != nil {
		return fmt.Errorf("postgres.Run: %w", err)
	}

	connStr := ctr.MustConnectionString(ctx, "sslmode=disable")
	u, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("parse connection string: %w", err)
	}
	if u.Scheme == "postgres" {
		u.Scheme = "postgresql"
	}
	u.Path = "/chainlink_test"
	if u.Path == "" || u.Path[0] != '/' {
		return errors.New("expected database path for chainlink_test")
	}

	if err := os.Setenv("CL_DATABASE_URL", u.String()); err != nil {
		return fmt.Errorf("set CL_DATABASE_URL: %w", err)
	}

	cfg := postgresTestResetConfig{u: *u}
	if err := store.ResetDatabaseQuick(ctx, lggr, cfg, false); err != nil {
		return fmt.Errorf("ResetDatabaseQuick: %w", err)
	}
	if err := store.PrepareTestDB(lggr, *u, false); err != nil {
		return fmt.Errorf("PrepareTestDB: %w", err)
	}
	return nil
}
