package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

// Handle owns the ephemeral Postgres used for a run. When the user supplied
// CL_DATABASE_URL the container is nil and Reset/Cleanup are no-ops.
type Handle struct {
	container *postgres.PostgresContainer
	conf      *config.App
}

// Ensure starts a Postgres testcontainer when CL_DATABASE_URL is unset, exports
// CL_DATABASE_URL, runs preparetest --force, then snapshots the prepared state
// so Reset can restore it between iterations. When CL_DATABASE_URL is already
// set, it is a no-op and Reset/Cleanup do nothing.
func Ensure(ctx context.Context, conf *config.App) (*Handle, error) {
	start := time.Now()

	if conf.PostgresVersion == "" {
		return &Handle{conf: conf}, errors.New("postgres version is required")
	}

	if conf.DatabaseURL != "" {
		fmt.Printf("Skipping database setup, using provided database URL: %s\n", conf.DatabaseURL)
		return &Handle{conf: conf}, nil
	}
	// We'll do our own cleanup, so disable Ryuk
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return &Handle{conf: conf}, fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", err)
	}

	if !conf.AIOutput {
		fmt.Print("Setting up postgres...")
	}

	c, err := postgres.Run(ctx,
		fmt.Sprintf("docker.io/postgres:%s-alpine", conf.PostgresVersion),
		postgres.WithDatabase("chainlink_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return &Handle{conf: conf}, fmt.Errorf("postgres testcontainer: %w", err)
	}

	h := &Handle{container: c, conf: conf}

	// Build the connection string for CL tests to use
	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return h, errors.Join(fmt.Errorf("connection string: %w", err), h.Cleanup())
	}

	// Set the connection string for CL tests to use
	if err := os.Setenv("CL_DATABASE_URL", connStr); err != nil {
		return h, errors.Join(err, h.Cleanup())
	}

	// Run preparetest --force to set up the database for tests
	prep := exec.CommandContext(ctx, "go", "run", "./core/store/cmd/preparetest", "--force")
	prep.Dir = conf.RepoRoot
	prep.Env = os.Environ()
	if err := prep.Run(); err != nil {
		return h, errors.Join(fmt.Errorf("preparetest --force: %w", err), h.Cleanup())
	}

	// Snapshot the prepared schema so Reset can restore it quickly between iterations.
	if err := c.Snapshot(ctx); err != nil {
		return h, errors.Join(fmt.Errorf("snapshot prepared database: %w", err), h.Cleanup())
	}

	if !conf.AIOutput {
		fmt.Printf(" ✅ (%s)\n", time.Since(start).Round(time.Millisecond))
	}

	return h, nil
}

// Reset restores the database to its freshly-prepared snapshot. No-op when the
// user supplied CL_DATABASE_URL (we don't own the database).
func (h *Handle) Reset(ctx context.Context) error {
	if h == nil || h.container == nil {
		return nil
	}
	start := time.Now()
	if !h.conf.AIOutput {
		fmt.Print("Resetting database...")
	}
	if err := h.container.Restore(ctx); err != nil {
		if !h.conf.AIOutput {
			fmt.Println(" ❌")
		}
		return fmt.Errorf("restore snapshot: %w", err)
	}
	if !h.conf.AIOutput {
		fmt.Printf(" ✅ (%s)\n", time.Since(start).Round(time.Millisecond))
	}
	return nil
}

// Cleanup terminates the Postgres testcontainer. Safe to call on a nil or
// no-container Handle.
func (h *Handle) Cleanup() error {
	if h == nil || h.container == nil {
		return nil
	}
	if !h.conf.AIOutput {
		fmt.Print("Tearing down postgres...")
	}
	termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := h.container.Terminate(termCtx); err != nil {
		if !h.conf.AIOutput {
			fmt.Println(" ❌")
		}
		return fmt.Errorf("error terminating postgres container, you need to terminate it manually: %w", err)
	}
	if !h.conf.AIOutput {
		fmt.Println(" ✅")
	}
	return nil
}
