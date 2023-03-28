package heavyweight

// The heavyweight package contains cltest items that are costly and you should
// think **real carefully** before using in your tests.

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// FullTestDBV2 creates a pristine DB which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types with it.
func FullTestDBV2(t *testing.T, name string, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return prepareFullTestDBV2(t, name, false, true, overrideFn)
}

// FullTestDBNoFixturesV2 is the same as FullTestDB, but it does not load fixtures.
func FullTestDBNoFixturesV2(t *testing.T, name string, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return prepareFullTestDBV2(t, name, false, false, overrideFn)
}

// FullTestDBEmptyV2 creates an empty DB (without migrations).
func FullTestDBEmptyV2(t *testing.T, name string, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return prepareFullTestDBV2(t, name, true, false, overrideFn)
}

func prepareFullTestDBV2(t *testing.T, name string, empty bool, loadFixtures bool, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	testutils.SkipShort(t, "FullTestDB")

	if empty && loadFixtures {
		t.Fatal("could not load fixtures into an empty DB")
	}

	gcfg := configtest2.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Dialect = dialects.Postgres
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	name = fmt.Sprintf("%s_%x", name, rand.Intn(0xFFF)) // to avoid name collisions
	migrationTestDBURL, err := dropAndCreateThrowawayTestDB(gcfg.DatabaseURL(), name, empty)
	require.NoError(t, err)
	db, err := pg.NewConnection(migrationTestDBURL, dialects.Postgres, gcfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		os.RemoveAll(gcfg.RootDir())
	})

	gcfg = configtest2.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Dialect = dialects.Postgres
		s.Database.URL = models.MustSecretURL(migrationTestDBURL)
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	if loadFixtures {
		_, filename, _, ok := runtime.Caller(1)
		if !ok {
			t.Fatal("could not get runtime.Caller(1)")
		}
		filepath := path.Join(path.Dir(filename), "../../../store/fixtures/fixtures.sql")
		fixturesSQL, err := os.ReadFile(filepath)
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

	// Match the naming schema that our dangling DB cleanup methods expect
	dbname := cmd.TestDBNamePrefix + postfix
	if l := len(dbname); l > 63 {
		return "", fmt.Errorf("dbname %v too long (%d), max is 63 bytes. Try a shorter postfix", dbname, l)
	}
	// Cannot drop test database if we are connected to it, so we must connect
	// to a different one. 'postgres' should be present on all postgres installations
	parsed.Path = "/postgres"
	db, err := sql.Open(string(dialects.Postgres), parsed.String())
	if err != nil {
		return "", fmt.Errorf("In order to drop the test database, we need to connect to a separate database"+
			" called 'postgres'. But we are unable to open 'postgres' database: %+v\n", err)
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
