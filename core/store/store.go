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
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	cutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
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

func ResetDatabase(ctx context.Context, lggr logger.Logger, cfg Config, force bool) error {
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
	schema, err := dumpSchema(u)
	if err != nil {
		return err
	}
	lggr.Debugf("Testing rollback and re-migrate for database: %#v", u.String())
	var baseVersionID int64 = 54
	if err := downAndUpDB(ctx, cfg, baseVersionID); err != nil {
		return err
	}
	return checkSchema(u, schema)
}

type Config interface {
	DefaultIdleInTxSessionTimeout() time.Duration
	DefaultLockTimeout() time.Duration
	MaxOpenConns() int
	MaxIdleConns() int
	URL() url.URL
	DriverName() string
}

var errDBURLMissing = errors.New("You must set CL_DATABASE_URL env variable or provide a secrets TOML with Database.URL set. HINT: If you are running this to set up your local test database, try CL_DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable")

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

func dropAndCreateDB(parsed url.URL, force bool) (err error) {
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
			err = multierr.Append(err, cerr)
		}
	}()
	if force {
		// supports pg < 13. https://stackoverflow.com/questions/17449420/postgresql-unable-to-drop-database-because-of-some-auto-connections-to-db
		_, err = db.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s';", dbname))
		if err != nil {
			return fmt.Errorf("unable to terminate connections to postgres database: %w", err)
		}
	}
	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, dbname))
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %w", err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname))
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %w", err)
	}
	return nil
}

func dropAndCreatePristineDB(db *sqlx.DB, template string) (err error) {
	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, testdb.PristineDBName))
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %w", err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s"`, testdb.PristineDBName, template))
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

func dumpSchema(dbURL url.URL) (string, error) {
	args := []string{
		dbURL.String(),
		"--schema-only",
	}
	cmd := exec.Command(
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

func checkSchema(dbURL url.URL, prevSchema string) error {
	newSchema, err := dumpSchema(dbURL)
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
			err = multierr.Append(err, cerr)
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
	_, err = db.Exec(string(fixturesSQL))
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
	for i := 0; i < nWorkers; i++ {
		go func() {
			defer wg.Done()
			for dbname := range ch {
				lggr.Infof("Dropping old, dangling test database: %q", dbname)
				gerr := cutils.JustError(db.Exec(`DROP DATABASE IF EXISTS ` + dbname))
				errCh <- gerr
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
		err = multierr.Append(err, gerr)
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
	seqRows, err := db.Query(`SELECT sequence_schema, sequence_name, minimum_value FROM information_schema.sequences WHERE sequence_schema IN ($1)`, schemas)
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
		randNum, err = crand.Int(crand.Reader, ubig.NewI(10000).ToInt())
		if err != nil {
			return fmt.Errorf("%s: failed to generate random number", failedToRandomizeTestDBSequencesError{})
		}
		randNum.Add(randNum, big.NewInt(minimumSequenceValue))

		if _, err = db.Exec(fmt.Sprintf("ALTER SEQUENCE %s.%s RESTART WITH %d", sequenceSchema, sequenceName, randNum)); err != nil {
			return fmt.Errorf("%s: failed to alter and restart %s sequence: %w", failedToRandomizeTestDBSequencesError{}, sequenceName, err)
		}
	}

	if err = seqRows.Err(); err != nil {
		return fmt.Errorf("%s: failed to iterate through sequences: %w", failedToRandomizeTestDBSequencesError{}, err)
	}

	return nil
}
