package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Ensure starts a Postgres testcontainer when CL_DATABASE_URL is unset, exports CL_DATABASE_URL,
// and runs preparetest --force from repoRoot. When CL_DATABASE_URL is already set, it is a no-op.
func Ensure(ctx context.Context, databaseURL, postgresVersion, repoRoot string) (cleanup func(), err error) {
	if postgresVersion == "" {
		return func() {}, fmt.Errorf("postgres version is required")
	}

	if databaseURL != "" {
		fmt.Printf("skipping database setup, using %s", databaseURL)
		return func() {}, nil
	}

	err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	if err != nil {
		return nil, fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", err)
	}

	c, err := postgres.Run(ctx,
		fmt.Sprintf("docker.io/postgres:%s-alpine", postgresVersion),
		postgres.WithDatabase("chainlink_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("postgres testcontainer: %w", err)
	}

	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.Terminate(termCtx)
		if err != nil {
			fmt.Println("error terminating postgres container, you need to terminate it manually:", err)
		}
		cancel()
		return nil, fmt.Errorf("connection string: %w", err)
	}

	if err := os.Setenv("CL_DATABASE_URL", connStr); err != nil {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.Terminate(termCtx)
		if err != nil {
			fmt.Println("error terminating postgres container, you need to terminate it manually:", err)
		}
		cancel()
		return nil, err
	}

	cleanup = func() {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := c.Terminate(termCtx)
		if err != nil {
			fmt.Println("error terminating postgres container, you need to terminate it manually:", err)
		}
	}

	prep := exec.CommandContext(ctx, "go", "run", "./core/store/cmd/preparetest", "--force")
	prep.Dir = repoRoot
	prep.Env = os.Environ()
	if err := prep.Run(); err != nil {
		cleanup()
		return nil, fmt.Errorf("preparetest --force: %w", err)
	}

	return cleanup, nil
}
