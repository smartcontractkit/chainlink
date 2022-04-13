package heavyweight

// The heavyweight package contains cltest items that are costly and you should
// think **real carefully** before using in your tests.

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
)

// FullTestDB creates a pristine DB which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types with it.
func FullTestDB(t *testing.T, name string) (*configtest.TestGeneralConfig, *sqlx.DB) {
	return prepareFullTestDB(t, name, false, true)
}

// FullTestDBNoFixtures is the same as FullTestDB, but it does not load fixtures.
func FullTestDBNoFixtures(t *testing.T, name string) (*configtest.TestGeneralConfig, *sqlx.DB) {
	return prepareFullTestDB(t, name, false, false)
}

// FullTestDBEmpty creates an empty DB (without migrations).
func FullTestDBEmpty(t *testing.T, name string) (*configtest.TestGeneralConfig, *sqlx.DB) {
	return prepareFullTestDB(t, name, true, false)
}

func prepareFullTestDB(t *testing.T, name string, empty bool, loadFixtures bool) (*configtest.TestGeneralConfig, *sqlx.DB) {
	testutils.SkipShort(t, "FullTestDB")

	if empty && loadFixtures {
		t.Fatal("could not load fixtures into an empty DB")
	}

	overrides := configtest.GeneralConfigOverrides{
		SecretGenerator: cltest.MockSecretGenerator{},
	}
	gcfg := configtest.NewTestGeneralConfigWithOverrides(t, overrides)
	gcfg.Overrides.Dialect = dialects.Postgres

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	migrationTestDBURL, err := dropAndCreateThrowawayTestDB(gcfg.DatabaseURL(), name, empty)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	db, err := pg.NewConnection(migrationTestDBURL, string(dialects.Postgres), pg.Config{
		Logger:       lggr,
		MaxOpenConns: gcfg.ORMMaxOpenConns(),
		MaxIdleConns: gcfg.ORMMaxIdleConns(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		os.RemoveAll(gcfg.RootDir())
	})
	gcfg.Overrides.DatabaseURL = null.StringFrom(migrationTestDBURL)

	if loadFixtures {
		_, filename, _, ok := runtime.Caller(1)
		if !ok {
			t.Fatal("could not get runtime.Caller(1)")
		}
		filepath := path.Join(path.Dir(filename), "../../../store/fixtures/fixtures.sql")
		fixturesSQL, err := ioutil.ReadFile(filepath)
		require.NoError(t, err)
		_, err = db.Exec(string(fixturesSQL))
		require.NoError(t, err)
	}

	return gcfg, db
}

func dropAndCreateThrowawayTestDB(parsed url.URL, postfix string, empty bool) (string, error) {
	if parsed.Path == "" {
		return "", errors.New("path missing from database URL")
	}

	dbname := fmt.Sprintf("%s_%s", parsed.Path[1:], postfix)
	if len(dbname) > 62 {
		return "", fmt.Errorf("dbname %v too long, max is 63 bytes. Try a shorter postfix", dbname)
	}
	// Cannot drop test database if we are connected to it, so we must connect
	// to a different one. 'postgres' should be present on all postgres installations
	parsed.Path = "/postgres"
	db, err := sql.Open(string(dialects.Postgres), parsed.String())
	if err != nil {
		return "", fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		return "", fmt.Errorf("unable to drop postgres migrations test database: %v", err)
	}
	if empty {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	} else {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", dbname, cmd.PristineDBName))
	}
	if err != nil {
		return "", fmt.Errorf("unable to create postgres test database with name '%s': %v", dbname, err)
	}
	parsed.Path = fmt.Sprintf("/%s", dbname)
	return parsed.String(), nil
}
