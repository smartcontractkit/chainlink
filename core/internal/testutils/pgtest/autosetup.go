package pgtest

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

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

// EnvTerminateContainer, when set to "true", terminates the reuse-named Postgres container
// when the first test that triggered autosetup finishes (default: leave it for reuse).
const EnvTerminateContainer = "CHAINLINK_PGTEST_TERMINATE_CONTAINER"

var (
	autoOnce sync.Once
	autoErr  error
	autoCtr  *postgres.PostgresContainer
)

// postgresTestResetConfig implements store.Config for prepare/migrate helpers.
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
// Reuses a named container across packages when Docker + testcontainers reuse are available.
func EnsureAutoPostgres(t testing.TB) {
	t.Helper()
	if testing.Short() {
		return
	}
	if os.Getenv(EnvDisableAutoContainers) == "true" {
		return
	}
	if string(env.DatabaseURL.Get()) != "" {
		return
	}

	autoOnce.Do(func() {
		autoErr = startEmbeddedPostgres(t)
		if autoErr != nil {
			return
		}
		if os.Getenv(EnvTerminateContainer) == "true" && autoCtr != nil {
			c := autoCtr
			t.Cleanup(func() {
				_ = c.Terminate(context.Background()) //nolint:contextcheck // best-effort cleanup
			})
		}
	})
	require.NoError(t, autoErr, "pgtest autosetup: start test postgres")
}

func startEmbeddedPostgres(t testing.TB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	lggr := logger.Test(t)

	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithReuseByName("chainlink-core-gotest-postgres-16"),
		// Supplement default `fsync=off` from the postgres module.
		testcontainers.WithCmdArgs(
			"-c", "synchronous_commit=off",
			"-c", "full_page_writes=off",
			"-c", "shared_buffers=256MB",
			"-c", "work_mem=64MB",
			"-c", "max_connections=500",
			"-c", "random_page_cost=1.1",
		),
	}

	ctr, err := postgres.Run(ctx, "postgres:16-alpine", opts...)
	if err != nil {
		return fmt.Errorf("postgres.Run: %w", err)
	}
	autoCtr = ctr

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
