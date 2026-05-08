package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

// Handle owns the ephemeral Postgres used for a run. When the user supplied
// CL_DATABASE_URL the container is nil and Reset/Cleanup are no-ops.
type Handle struct {
	container *postgres.PostgresContainer
	conf      *config.App
	out       *output.Printer
	connStr   string
}

// Resource is the runner-facing view of one prepared test database.
type Resource struct {
	Env             []string
	Reset           func(context.Context) error
	DumpDiagnostics func(context.Context, string, int) error
}

// Pool owns the database handles used by a command.
type Pool struct {
	handles []*Handle
}

// Ensure configures CL_DATABASE_URL for child test processes. If --database-url
// is set, that value is exported as CL_DATABASE_URL (failing if CL_DATABASE_URL is
// already set to something else). Otherwise it starts an ephemeral Postgres
// container, sets CL_DATABASE_URL to its connection string, runs preparetest
// --force, and snapshots the prepared state so Reset can restore it between
// diagnose iterations.
func Ensure(ctx context.Context, conf *config.App, out *output.Printer) (h *Handle, err error) {
	return ensure(ctx, conf, out, true)
}

func ensure(ctx context.Context, conf *config.App, out *output.Printer, setGlobalDatabaseURL bool) (h *Handle, err error) {
	if out == nil {
		out = output.New(conf.AIOutput, io.Discard, io.Discard, output.SkipFD)
	}
	start := time.Now()

	if conf.PostgresVersion == "" {
		return &Handle{conf: conf, out: out}, errors.New("postgres version is required")
	}

	if conf.DatabaseURL != "" {
		if existing := os.Getenv("CL_DATABASE_URL"); existing != "" && existing != conf.DatabaseURL {
			return &Handle{conf: conf, out: out}, errors.New("CL_DATABASE_URL is already set to a different value than --database-url (refusing to override)")
		}
		if setGlobalDatabaseURL {
			if err = os.Setenv("CL_DATABASE_URL", conf.DatabaseURL); err != nil {
				return &Handle{conf: conf, out: out}, fmt.Errorf("set CL_DATABASE_URL: %w", err)
			}
		}
		out.IfHuman(func() {
			out.HumanStdout(
				termstyle.Muted.Render("Skipping database setup, using provided database URL: ") +
					termstyle.Label.Render(conf.DatabaseURL))
		})
		return &Handle{conf: conf, out: out, connStr: conf.DatabaseURL}, nil
	}
	// Intentional: Ryuk is disabled because this harness always tears down via
	// Handle.Cleanup(); Ryuk can conflict with that lifecycle in some setups.
	if err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return &Handle{conf: conf, out: out}, fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", err)
	}

	// Progress on stderr, same escape and TTY rules as diagnoseIteration /
	// renderDiagnoseProgressLine (runner).
	setupPartial := false
	out.IfHuman(func() {
		out.HumanFprint(termstyle.Label.Render("Setting up Postgres..."))
		setupPartial = true
	})
	abortSetupPartial := func() {
		if !setupPartial {
			return
		}
		_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K\n")
		setupPartial = false
	}
	defer func() {
		if err != nil {
			abortSetupPartial()
		}
	}()

	c, err := postgres.Run(ctx,
		fmt.Sprintf("docker.io/postgres:%s-alpine", conf.PostgresVersion),
		postgres.WithDatabase("chainlink_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithCmdArgs("-c", "max_connections=1000"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return &Handle{conf: conf, out: out}, fmt.Errorf("postgres testcontainer: %w", err)
	}

	h = &Handle{container: c, conf: conf, out: out}

	// Build the connection string for CL tests to use
	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return h, errors.Join(fmt.Errorf("connection string: %w", err), h.Cleanup())
	}
	h.connStr = connStr

	if setGlobalDatabaseURL {
		// Set the connection string for CL tests to use.
		if err := os.Setenv("CL_DATABASE_URL", connStr); err != nil {
			return h, errors.Join(err, h.Cleanup())
		}
	}

	// Run preparetest --force to set up the database for tests
	prepareOutput := bytes.NewBuffer(nil)
	prep := exec.CommandContext(ctx, "go", "run", "./core/store/cmd/preparetest", "--force")
	prep.Dir = conf.RepoRoot
	prep.Env = append(os.Environ(), h.Env()...)
	prep.Stdout = prepareOutput
	prep.Stderr = prepareOutput
	if err := prep.Run(); err != nil {
		return h, errors.Join(fmt.Errorf("preparetest --force: %w\n%s", err, prepareOutput.String()), h.Cleanup())
	}

	// Snapshot the prepared schema so Reset can restore it quickly between iterations.
	if err := c.Snapshot(ctx); err != nil {
		return h, errors.Join(fmt.Errorf("snapshot prepared database: %w", err), h.Cleanup())
	}

	out.IfHuman(func() {
		_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K")
		out.HumanStderr(
			termstyle.Label.Render("Setup Postgres") + " " +
				termstyle.OK.Render("✅") + " " +
				termstyle.Muted.Render(fmt.Sprintf("(%s)", time.Since(start).Round(time.Millisecond))))
		setupPartial = false
	})

	return h, nil
}

// EnsurePool creates the prepared databases needed by diagnose. Parallel
// diagnose uses one ephemeral Postgres per worker; external database URLs are
// only allowed for serial runs because the harness cannot isolate or reset them.
func EnsurePool(ctx context.Context, conf *config.App, out *output.Printer, size int) (*Pool, error) {
	if size < 1 {
		return nil, errors.New("database pool size must be >= 1")
	}
	if size > 1 && conf.DatabaseURL != "" {
		return nil, errors.New("--parallel-iterations > 1 cannot be used with --database-url")
	}

	start := time.Now()
	setupPartial := false
	if size > 1 && out != nil {
		out.IfHuman(func() {
			out.HumanFprint(termstyle.Label.Render(fmt.Sprintf("Setting up %d Postgres...", size)))
			setupPartial = true
		})
	}

	poolCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	handles := make([]*Handle, size)
	errs := make(chan error, size)
	var wg sync.WaitGroup

	for i := range size {
		wg.Go(func() {
			workerOut := out
			if size > 1 {
				workerOut = output.New(conf.AIOutput, io.Discard, io.Discard, output.SkipFD)
			}
			h, err := ensure(poolCtx, conf, workerOut, size == 1)
			if err != nil {
				cancel()
				errs <- err
				return
			}
			handles[i] = h
		})
	}
	wg.Wait()
	close(errs)

	var err error
	for e := range errs {
		err = errors.Join(err, e)
	}

	pool := &Pool{handles: handles}
	if err != nil {
		if setupPartial && out != nil {
			_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K\n")
		}
		return pool, errors.Join(err, pool.Cleanup())
	}

	for _, h := range handles {
		if h == nil {
			if setupPartial && out != nil {
				_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K\n")
			}
			return pool, errors.Join(errors.New("database pool factory returned nil handle"), pool.Cleanup())
		}
	}

	if setupPartial && out != nil {
		out.IfHuman(func() {
			_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K")
			out.HumanStderr(
				termstyle.Label.Render(fmt.Sprintf("Setup %d Postgres", size)) + " " +
					termstyle.OK.Render("✅") + " " +
					termstyle.Muted.Render(fmt.Sprintf("(%s)", time.Since(start).Round(time.Millisecond))))
		})
	}

	return pool, nil
}

// Handles returns the underlying handles for tests and low-level callers.
func (p *Pool) Handles() []*Handle {
	if p == nil {
		return nil
	}
	return p.handles
}

// Resources returns the runner-facing DB resources.
func (p *Pool) Resources() []Resource {
	if p == nil {
		return nil
	}
	resources := make([]Resource, 0, len(p.handles))
	for _, h := range p.handles {
		resources = append(resources, Resource{
			Env:             h.Env(),
			Reset:           h.Reset,
			DumpDiagnostics: h.DumpDiagnostics,
		})
	}
	return resources
}

// Cleanup tears down every database handle in the pool.
func (p *Pool) Cleanup() error {
	if p == nil || len(p.handles) == 0 {
		return nil
	}
	errs := make(chan error, len(p.handles))
	var wg sync.WaitGroup
	for _, h := range p.handles {
		wg.Go(func() {
			errs <- h.Cleanup()
		})
	}
	wg.Wait()
	close(errs)
	var err error
	for e := range errs {
		err = errors.Join(err, e)
	}
	return err
}

// Env returns environment overrides for child test processes using this handle.
func (h *Handle) Env() []string {
	if h == nil || h.connStr == "" {
		return nil
	}
	return []string{"CL_DATABASE_URL=" + h.connStr}
}

// Reset restores the database to its freshly-prepared snapshot. No-op when the
// user supplied CL_DATABASE_URL (we don't own the database).
func (h *Handle) Reset(ctx context.Context) error {
	if h == nil || h.container == nil {
		return nil
	}
	if err := h.container.Restore(ctx); err != nil {
		return fmt.Errorf("restore snapshot: %w", err)
	}
	return nil
}

// DumpDiagnostics writes postgres-state-<iteration>.md to dir with the container log
// and key system-view snapshots for that iteration. No-op when the user
// supplied CL_DATABASE_URL (we don't own that database).
func (h *Handle) DumpDiagnostics(ctx context.Context, dir string, iteration int) error {
	if h == nil || h.container == nil {
		return nil
	}

	name := fmt.Sprintf("postgres-state-%d.md", iteration)
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("create %s: %w", name, err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# Postgres State\n\nCaptured: %s\n\n", time.Now().UTC().Format(time.RFC3339))

	// Container (server) log.
	fmt.Fprint(f, "## Server Log\n\n```\n")
	logs, logErr := h.container.Logs(ctx)
	if logErr != nil {
		return fmt.Errorf("fetch logs: %w", logErr)
	}
	defer func() {
		if logCloseErr := logs.Close(); logCloseErr != nil {
			// Ignore errors closing logs, we're dumping diagnostics to a file.
			fmt.Fprintf(f, "error closing logs: %v\n", logCloseErr)
		}
	}()
	_, err = io.Copy(f, logs)
	if err != nil {
		return fmt.Errorf("copy logs: %w", err)
	}
	fmt.Fprint(f, "```\n\n")

	type query struct {
		heading string
		sql     string
	}
	queries := []query{
		{
			"Active Connections (pg_stat_activity)",
			`SELECT pid, state, wait_event_type, wait_event, query_start, left(query,120) AS query ` +
				`FROM pg_stat_activity WHERE datname='chainlink_test' ORDER BY query_start;`,
		},
		{
			"Locks (pg_locks + pg_stat_activity)",
			`SELECT l.pid, l.locktype, l.relation::regclass, l.mode, l.granted, left(a.query,80) AS query ` +
				`FROM pg_locks l LEFT JOIN pg_stat_activity a ON a.pid=l.pid ` +
				`WHERE l.relation IS NOT NULL ORDER BY l.granted, l.pid;`,
		},
		{
			"Table Statistics (pg_stat_user_tables)",
			`SELECT relname, seq_scan, idx_scan, n_tup_ins, n_tup_upd, n_tup_del, n_live_tup, n_dead_tup ` +
				`FROM pg_stat_user_tables ORDER BY n_live_tup DESC LIMIT 30;`,
		},
		{
			"Database Size",
			`SELECT pg_size_pretty(pg_database_size('chainlink_test')) AS db_size;`,
		},
	}

	for _, q := range queries {
		fmt.Fprintf(f, "## %s\n\n```\n", q.heading)
		exitCode, out, execErr := h.container.Exec(ctx,
			[]string{"psql", "-U", "postgres", "-d", "chainlink_test", "-P", "pager=off", "-c", q.sql},
		)
		switch {
		case execErr != nil:
			fmt.Fprintf(f, "error: %v\n", execErr)
		case exitCode != 0:
			fmt.Fprintf(f, "psql exit %d\n", exitCode)
		}
		_, err = io.Copy(f, out)
		if err != nil {
			return fmt.Errorf("copy output: %w", err)
		}
		fmt.Fprint(f, "```\n\n")
	}

	return nil
}

// Cleanup terminates the Postgres testcontainer. Safe to call on a nil or
// no-container Handle.
func (h *Handle) Cleanup() error {
	if h == nil || h.container == nil {
		return nil
	}
	if h.out != nil {
		h.out.IfHuman(func() {
			h.out.HumanFprint(termstyle.Label.Render("Tearing down postgres..."))
		})
	}
	termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := h.container.Terminate(termCtx); err != nil {
		if h.out != nil {
			h.out.IfHuman(func() {
				h.out.HumanStderr(" " + termstyle.Bad.Render("❌"))
			})
		}
		return fmt.Errorf("error terminating postgres container, you need to terminate it manually: %w", err)
	}
	if h.out != nil {
		h.out.IfHuman(func() {
			h.out.HumanStderr(" " + termstyle.OK.Render("✅"))
		})
	}
	return nil
}
