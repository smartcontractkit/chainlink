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

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	migrations "github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

// FullTestDB creates an DB which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types
// with it.
func FullTestDB(t *testing.T, name string, migrate bool, loadFixtures bool) (*configtest.TestGeneralConfig, *sqlx.DB) {
	if testing.Short() {
		t.Skip("skipping due to use of FullTestDB")
	}
	overrides := configtest.GeneralConfigOverrides{
		SecretGenerator: cltest.MockSecretGenerator{},
	}
	gcfg := configtest.NewTestGeneralConfigWithOverrides(t, overrides)
	gcfg.SetDialect(dialects.Postgres)

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	migrationTestDBURL, err := dropAndCreateThrowawayTestDB(gcfg.DatabaseURL(), name)
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
	if migrate {
		require.NoError(t, migrations.Migrate(db.DB, lggr))
	}
	if loadFixtures {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("could not get runtime.Caller(0)")
		}
		filepath := path.Join(path.Dir(filename), "../../../store/fixtures/fixtures.sql")
		fixturesSQL, err := ioutil.ReadFile(filepath)
		require.NoError(t, err)
		_, err = db.Exec(string(fixturesSQL))
		require.NoError(t, err)
	}

	return gcfg, db
}

func dropAndCreateThrowawayTestDB(parsed url.URL, postfix string) (string, error) {
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
	// `CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		return "", fmt.Errorf("unable to create postgres migrations test database with name '%s': %v", dbname, err)
	}
	parsed.Path = fmt.Sprintf("/%s", dbname)
	return parsed.String(), nil
}
