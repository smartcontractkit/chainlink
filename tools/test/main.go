// Package main launches the tools/test CLI: a thin testrig consumer that
// provides Chainlink's ephemeral/persistent Postgres setup as testrig resources.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/testrig"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/dbdetect"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/dbflags"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
)

func main() {
	testrig.Run(
		testrig.WithRootCommand("make test"),
		// Add --database-url / --postgres-version to the root command.
		testrig.WithRootFlags(dbflags.Register),
		// Provide one prepared Postgres per diagnose worker (or one for run/gotestsum).
		testrig.WithResources(dbProvider),
	)
}

// dbProvider stands up `count` prepared test databases and exposes each as a
// testrig.Resource. With --database-url (or CL_DATABASE_URL) set, a single
// external database is reused and count>1 is rejected; otherwise each worker
// gets an isolated ephemeral container with a snapshot for fast Reset.
func dbProvider(ctx context.Context, count int) ([]testrig.Resource, error) {
	return dbProviderForArgs(ctx, count, os.Args[1:])
}

func dbProviderForArgs(ctx context.Context, count int, args []string) ([]testrig.Resource, error) {
	conf, err := dbflags.AppConfig()
	if err != nil {
		return nil, err
	}

	// Statically check if Postgres is actually needed for the tests being run.
	needsDB, err := dbdetect.NeedsPostgres(conf.RepoRoot, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[dbdetect] Error checking DB need: %v. Defaulting to Postgres needed.\n", err)
		needsDB = true
	}

	if !needsDB {
		return noopResources(count), nil
	}

	conf.ParallelIterations = count
	out := output.NewFromApp(conf)

	pool, err := db.EnsurePool(ctx, conf, out, count)
	if err != nil {
		return nil, err
	}

	handles := pool.Handles()
	resources := make([]testrig.Resource, 0, len(handles))
	for _, h := range handles {
		resources = append(resources, testrig.Resource{
			Env:             h.Env(),
			Reset:           h.Reset,
			DumpDiagnostics: h.DumpDiagnostics,
			Cleanup:         h.Cleanup,
		})
	}
	return resources, nil
}

// noopResources returns empty testrig.Resource values for runs that do not need
// Postgres. testrig treats all Resource fields as optional (see testrig
// hooks.Resource); Reset, DumpDiagnostics, and Cleanup are only invoked when
// non-nil.
func noopResources(count int) []testrig.Resource {
	resources := make([]testrig.Resource, count)
	for i := range count {
		resources[i] = testrig.Resource{}
	}
	return resources
}
