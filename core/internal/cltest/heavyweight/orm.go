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

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	migrations "github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

// FullTestORM creates an ORM which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types
// with it.
func FullTestORM(t *testing.T, name string, migrate bool, loadFixtures ...bool) (*configtest.TestEVMConfig, *orm.ORM, func()) {
	overrides := configtest.GeneralConfigOverrides{
		SecretGenerator: cltest.MockSecretGenerator{},
	}
	gcfg := configtest.NewTestGeneralConfigWithOverrides(t, overrides)
	config := configtest.NewTestEVMConfig(t, gcfg)
	config.SetDialect(dialects.PostgresWithoutLock)

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	migrationTestDBURL, err := dropAndCreateThrowawayTestDB(config.DatabaseURL(), name)
	require.NoError(t, err)
	orm, err := orm.NewORM(migrationTestDBURL, config.DatabaseTimeout(), gracefulpanic.NewSignal(), dialects.PostgresWithoutLock, 0, config.GlobalLockRetryInterval().Duration(), config.ORMMaxOpenConns(), config.ORMMaxIdleConns())
	require.NoError(t, err)
	orm.SetLogging(config.LogSQLMigrations())
	config.GeneralConfig.Overrides.DatabaseURL = null.StringFrom(migrationTestDBURL)
	if migrate {
		require.NoError(t, migrations.Migrate(postgres.UnwrapGormDB(orm.DB).DB))
	}
	if len(loadFixtures) > 0 && loadFixtures[0] {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("could not get runtime.Caller(0)")
		}
		filepath := path.Join(path.Dir(filename), "../../../store/fixtures/fixtures.sql")
		fixturesSQL, err := ioutil.ReadFile(filepath)
		require.NoError(t, err)
		err = orm.DB.Exec(string(fixturesSQL)).Error
		require.NoError(t, err)
	}
	orm.SetLogging(config.LogSQLStatements())

	return config, orm, func() {
		assert.NoError(t, orm.Close())
		os.RemoveAll(config.RootDir())
	}
}

func dropAndCreateThrowawayTestDB(parsed url.URL, postfix string) (string, error) {
	if parsed.Path == "" {
		return "", errors.New("path missing from database URL")
	}

	dbname := fmt.Sprintf("%s_%s", parsed.Path[1:], postfix)
	if len(dbname) > 62 {
		return "", errors.New("dbname too long, max is 63 bytes. Try a shorter postfix")
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
