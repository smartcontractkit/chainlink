package store

import (
	"context"
	crand "crypto/rand"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/diff"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	cutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

//go:embed fixtures/fixtures.sql
var fixturesSQL string

func FixturesSQL() string { return fixturesSQL }

func PrepareTestDB(lggr logger.Logger, dbURL url.URL, userOnly bool) error {
	db, err := sqlx.Open(pgcommon.DriverPostgres, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()
	templateDB := strings.Trim(dbURL.Path, "/")
	if err = dropAndCreatePristineDB(db, templateDB); err != nil {
		return err
	}

	fixturePath := "../store/fixtures/fixtures.sql"
	if userOnly {
		fixturePath = "../store/fixtures/users_only_fixture.sql"
	}
	if err = insertFixtures(dbURL, fixturePath); err != nil {
		return err
	}
	if err = dropDanglingTestDBs(lggr, db); err != nil {
		return err
	}
	return randomizeTestDBSequences(db)
}

func ResetDatabase(ctx context.Context, lggr logger.Logger, cfg Config, force, deterministic bool) error {
	u := cfg.URL()
	lggr.Infof("Resetting database: %#v", u.String())
	lggr.Debugf("Dropping and recreating database: %#v", u.String())
	if err := dropAndCreateDB(u, force); err != nil {
		return err
	}
	lggr.Debugf("Migrating database: %#v", u.String())
	if err := migrateDB(ctx, cfg); err != nil {
		return err
	}
	var restrictKey string
	if deterministic {
		restrictKey = "chainlinktestrestrictkey"
	}
	schema, err := dumpSchema(u, restrictKey)
	if err != nil {
		return err
	}
	lggr.Debugf("Testing rollback and re-migrate for database: %#v", u.String())
	var baseVersionID int64 = 54
	if err := downAndUpDB(ctx, cfg, baseVersionID); err != nil {
		return err
	}
	return checkSchema(u, schema, restrictKey)
}

type Config interface {
	DefaultIdleInTxSessionTimeout() time.Duration
	DefaultLockTimeout() time.Duration
	MaxOpenConns() int
	MaxIdleConns() int
	URL() url.URL
	DriverName() string
}

var errDBURLMissing = errors.New("you must set CL_DATABASE_URL env variable or provide a secrets TOML with Database.URL set; if you are running this to set up your local test database, try CL_DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable")

func NewConnection(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	parsed := cfg.URL()
	if parsed.String() == "" {
		return nil, errDBURLMissing
	}
	return pg.NewConnection(ctx, parsed.String(), cfg.DriverName(), cfg)
}

func migrateDB(ctx context.Context, config Config) error {
	db, err := NewConnection(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	if err = migrate.Migrate(ctx, db.DB); err != nil {
		return fmt.Errorf("migrateDB failed: %w", err)
	}
	return db.Close()
}

func dropAndCreateDB(parsed url.URL, _ bool) (err error) {
	// Cannot drop the database if we are connected to it, so we must connect
	// to a different one. template1 should be present on all postgres installations
	dbname := parsed.Path[1:]
	parsed.Path = "/template1"
	db, err := sql.Open(pgcommon.DriverPostgres, parsed.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()
	// DROP ... WITH (FORCE) requires PostgreSQL 13+; replaces pg_terminate_backend + DROP for older versions.
	// Second parameter kept for ResetDatabase API compatibility (preparetest --force).
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// PostgreSQL does not support bound parameters for database names; pq.QuoteIdentifier is the supported escape.
	_, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+pq.QuoteIdentifier(dbname)+" WITH (FORCE)") //nolint:gosec // G701 false positive: identifier from pq.QuoteIdentifier only
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %w", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "CREATE DATABASE "+pq.QuoteIdentifier(dbname)) //nolint:gosec // G701 false positive: identifier from pq.QuoteIdentifier only
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %w", err)
	}
	return nil
}

func dropAndCreatePristineDB(db *sqlx.DB, template string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+pq.QuoteIdentifier(testdb.PristineDBName)+" WITH (FORCE)")
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %w", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "CREATE DATABASE "+pq.QuoteIdentifier(testdb.PristineDBName)+" WITH TEMPLATE "+pq.QuoteIdentifier(template)) //nolint:gosec // G701 false positive: identifiers from pq.QuoteIdentifier only
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %w", err)
	}
	return nil
}

func downAndUpDB(ctx context.Context, cfg Config, baseVersionID int64) error {
	db, err := NewConnection(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}
	if err = migrate.Rollback(ctx, db.DB, null.IntFrom(baseVersionID)); err != nil {
		return fmt.Errorf("test rollback failed: %w", err)
	}
	if err = migrate.Migrate(ctx, db.DB); err != nil {
		return fmt.Errorf("second migrateDB failed: %w", err)
	}
	return db.Close()
}

func dumpSchema(dbURL url.URL, restrictKey string) (string, error) {
	args := []string{
		dbURL.String(),
		"--schema-only",
	}

	// Only add restrict-key if it's supported (PostgreSQL v17+).
	// This is used for deterministic schema dumps which CI runs to compare
	// previous and new schemas.
	if restrictKey != "" {
		// Test if pg_dump supports --restrict-key
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		testCmd := exec.CommandContext(ctx, "pg_dump", "--help")
		helpOutput, err := testCmd.Output()
		if err == nil && strings.Contains(string(helpOutput), "--restrict-key") {
			args = append(args, "--restrict-key="+restrictKey)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx,
		"pg_dump", args...,
	)

	schema, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return "", fmt.Errorf("failed to dump schema: %w\n%s", err, string(ee.Stderr))
		}
		return "", fmt.Errorf("failed to dump schema: %w", err)
	}
	return string(schema), nil
}

func checkSchema(dbURL url.URL, prevSchema string, restrictKey string) error {
	newSchema, err := dumpSchema(dbURL, restrictKey)
	if err != nil {
		return err
	}
	df := diff.Diff(prevSchema, newSchema)
	if len(df) > 0 {
		fmt.Println(df)
		return errors.New("schema pre- and post- rollback does not match (ctrl+f for '+' or '-' to find the changed lines)")
	}
	return nil
}
func insertFixtures(dbURL url.URL, pathToFixtures string) (err error) {
	db, err := sql.Open(pgcommon.DriverPostgres, dbURL.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("could not get runtime.Caller(1)")
	}
	filepath := path.Join(path.Dir(filename), pathToFixtures)
	fixturesSQL, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, string(fixturesSQL))
	return err
}

func dropDanglingTestDBs(lggr logger.Logger, db *sqlx.DB) (err error) {
	// Drop all old dangling databases
	var dbs []string
	if err = db.Select(&dbs, `SELECT datname FROM pg_database WHERE datistemplate = false;`); err != nil {
		return err
	}

	// dropping database is very slow in postgres so we parallelise it here
	nWorkers := 25
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	errCh := make(chan error, len(dbs))
	for range nWorkers {
		go func() {
			defer wg.Done()
			for dbname := range ch {
				lggr.Infof("Dropping old, dangling test database: %q", dbname)
				errCh <- func() error {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					return cutils.JustError(db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+pq.QuoteIdentifier(dbname)+" WITH (FORCE)"))
				}()
			}
		}()
	}
	for _, dbname := range dbs {
		if strings.HasPrefix(dbname, testdb.TestDBNamePrefix) && !strings.HasSuffix(dbname, "_pristine") {
			ch <- dbname
		}
	}
	close(ch)
	wg.Wait()
	close(errCh)
	for gerr := range errCh {
		err = errors.Join(err, gerr)
	}
	return
}

type failedToRandomizeTestDBSequencesError struct{}

func (m *failedToRandomizeTestDBSequencesError) Error() string {
	return "failed to randomize test db sequences"
}

// randomizeTestDBSequences randomizes sequenced table columns sequence
// This is necessary as to avoid false positives in some test cases.
func randomizeTestDBSequences(db *sqlx.DB) error {
	// not ideal to hard code this, but also not safe to do it programmatically :(
	schemas := pq.Array([]string{"public", "evm"})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	seqRows, err := db.QueryContext(ctx, `SELECT sequence_schema, sequence_name, minimum_value FROM information_schema.sequences WHERE sequence_schema IN ($1)`, schemas)
	if err != nil {
		return fmt.Errorf("%s: error fetching sequences: %w", failedToRandomizeTestDBSequencesError{}, err)
	}

	defer seqRows.Close()
	for seqRows.Next() {
		var sequenceSchema, sequenceName string
		var minimumSequenceValue int64
		if err = seqRows.Scan(&sequenceSchema, &sequenceName, &minimumSequenceValue); err != nil {
			return fmt.Errorf("%s: failed scanning sequence rows: %w", failedToRandomizeTestDBSequencesError{}, err)
		}

		if sequenceName == "goose_migrations_id_seq" || sequenceName == "configurations_id_seq" {
			continue
		}

		var randNum *big.Int
		randNum, err = crand.Int(crand.Reader, sqlutil.NewI(10000).ToInt())
		if err != nil {
			return fmt.Errorf("%s: failed to generate random number", failedToRandomizeTestDBSequencesError{})
		}
		randNum.Add(randNum, big.NewInt(minimumSequenceValue))

		if err = func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return cutils.JustError(db.ExecContext(ctx, fmt.Sprintf("ALTER SEQUENCE %s.%s RESTART WITH %d", sequenceSchema, sequenceName, randNum)))
		}(); err != nil {
			return fmt.Errorf("%s: failed to alter and restart %s sequence: %w", failedToRandomizeTestDBSequencesError{}, sequenceName, err)
		}
	}

	if err = seqRows.Err(); err != nil {
		return fmt.Errorf("%s: failed to iterate through sequences: %w", failedToRandomizeTestDBSequencesError{}, err)
	}

	return nil
}
