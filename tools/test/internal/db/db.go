package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"charm.land/huh/v2/spinner"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Ensure starts a Postgres testcontainer when CL_DATABASE_URL is unset, exports CL_DATABASE_URL,
// and runs preparetest --force from repoRoot. When CL_DATABASE_URL is already set, it is a no-op.
func Ensure(ctx context.Context, conf *config.App) (cleanupDB func() error, err error) {
	if conf.PostgresVersion == "" {
		return func() error { return nil }, fmt.Errorf("postgres version is required")
	}

	if conf.DatabaseURL != "" {
		fmt.Printf("Skipping database setup, using provided database URL: %s\n", conf.DatabaseURL)
		return func() error { return nil }, nil
	}
	// We'll do our own cleanup, so disable Ryuk
	err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	if err != nil {
		return func() error { return nil }, fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", err)
	}

	aiOutput := ctx.Value("aiOutput").(bool)

	var c *postgres.PostgresContainer
	setupDB := func(ctx context.Context) error {
		c, err = postgres.Run(ctx,
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
			return fmt.Errorf("postgres testcontainer: %w", err)
		}
		return nil
	}

	if !aiOutput {
		err = spinner.New().Title("Setting up postgres...").Context(ctx).ActionWithErr(setupDB).Run()
		if err != nil {
			return func() error { return nil }, fmt.Errorf("setting up postgres: %w", err)
		}
	}

	// Build the connection string for CL tests to use
	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.Terminate(termCtx)
		if err != nil {
			fmt.Println("error terminating postgres container, you need to terminate it manually:", err)
		}
		cancel()
		return func() error { return nil }, fmt.Errorf("connection string: %w", err)
	}

	// Set the connection string for CL tests to use
	if err := os.Setenv("CL_DATABASE_URL", connStr); err != nil {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.Terminate(termCtx)
		if err != nil {
			fmt.Println("error terminating postgres container, you need to terminate it manually:", err)
		}
		cancel()
		return func() error { return err }, nil
	}

	// Setup the cleanup function to terminate the postgres container
	cleanupDB = func() error {
		termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := c.Terminate(termCtx)
		if err != nil {
			return fmt.Errorf("error terminating postgres container, you need to terminate it manually: %w", err)
		}
		return nil
	}

	// Run preparetest --force to set up the database for tests
	prep := exec.CommandContext(ctx, "go", "run", "./core/store/cmd/preparetest", "--force")
	prep.Dir = conf.RepoRoot
	prep.Env = os.Environ()
	if err := prep.Run(); err != nil { // If it fails, cleanup the database and return the error
		prepareErr := fmt.Errorf("preparetest --force: %w", err)
		cleanupErr := cleanupDB()
		if cleanupErr != nil {
			prepareErr = errors.Join(prepareErr, cleanupErr)
		}
		return func() error { return prepareErr }, prepareErr
	}

	return cleanupDB, nil
}
