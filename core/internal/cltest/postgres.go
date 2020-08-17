package cltest

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

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dropAndCreateThrowawayTestDB(databaseURL string, postfix string) (string, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}

	if parsed.Path == "" {
		return "", errors.New("path missing from database URL")
	}

	dbname := fmt.Sprintf("%s_%s", parsed.Path[1:], postfix)
	if len(dbname) > 62 {
		return "", errors.New("dbname too long, max is 63 bytes. Try a shorter postfix")
	}
	// Cannot drop test database if we are connected to it, so we must connect
	// to a different one. template1 should be present on all postgres installations
	parsed.Path = "/template1"
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
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
		return "", fmt.Errorf("unable to create postgres migrations test database: %v", err)
	}
	parsed.Path = fmt.Sprintf("/%s", dbname)
	return parsed.String(), nil
}

// BootstrapThrowawayORM creates an ORM which runs in a separate database
// than the normal unit tests, so it you can do things like use other
// Postgres connection types with it.
func BootstrapThrowawayORM(t *testing.T, name string, migrate bool, loadFixtures ...bool) (*TestConfig, *orm.ORM, func()) {
	tc, cleanup := NewConfig(t)
	config := tc.Config

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	migrationTestDBURL, err := dropAndCreateThrowawayTestDB(tc.DatabaseURL(), name)
	require.NoError(t, err)
	orm, err := orm.NewORM(migrationTestDBURL, config.DatabaseTimeout(), gracefulpanic.NewSignal(), orm.DialectPostgres, config.GetAdvisoryLockIDConfiguredOrDefault())
	require.NoError(t, err)
	orm.SetLogging(true)
	tc.Config.Set("DATABASE_URL", migrationTestDBURL)
	if migrate {
		require.NoError(t, orm.RawDB(func(db *gorm.DB) error { return migrations.Migrate(db) }))
	}
	if len(loadFixtures) > 0 && loadFixtures[0] {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("could not get runtime.Caller(0)")
		}
		filepath := path.Join(path.Dir(filename), "../../store/testdata/fixtures.sql")
		fixturesSQL, err := ioutil.ReadFile(filepath)
		require.NoError(t, err)
		err = orm.RawDB(func(db *gorm.DB) error {
			return db.Exec(string(fixturesSQL)).Error
		})
		require.NoError(t, err)
	}

	return tc, orm, func() {
		assert.NoError(t, orm.Close())
		cleanup()
		os.RemoveAll(config.RootDir())
	}
}
