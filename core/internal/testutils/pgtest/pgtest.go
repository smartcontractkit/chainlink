package pgtest

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/txdb"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

func initUnitTestDB(ctx context.Context) (string, error) {
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		return "", fmt.Errorf("you must provide a CL_DATABASE_URL environment variable")
	}

	testDBURL, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}

	// If the CL_USE_UNIT_TEST_DB env var is set, create a new test database and migrate it
	if os.Getenv("CL_USE_UNIT_TEST_DB") != "" {
		uuid := strings.ReplaceAll(uuid.New().String(), "-", "")
		testURL, err := testdb.CreateOrReplace(*testDBURL, uuid[:16]+"_unit_test", false)
		if err != nil {
			return "", err
		}
		testDBURL, err = url.Parse(testURL)
		if err != nil {
			return "", err
		}
		// migrate the test database
		testDB, err := sql.Open(string(dialects.Postgres), testDBURL.String())
		if err != nil {
			return "", err
		}
		err = migrate.Migrate(ctx, testDB)
		if err != nil {
			return "", err
		}
	}
	driver := string(dialects.TransactionWrappedPostgres) + "unit_test"
	err = txdb.RegisterTestDB(driver, testDBURL)
	if err != nil {
		return "", err
	}
	fmt.Println("Test database URL:", testDBURL.String())

	return driver, nil
}

var initOnce sync.Once
var driver string
var err error

func NewSqlxDB(t testing.TB) *sqlx.DB {
	testutils.SkipShortDB(t)
	initOnce.Do(func() {
		driver, err = initUnitTestDB(testutils.Context(t))
	})
	fmt.Println("driver:", driver)
	require.NoError(t, err)
	require.NotEmpty(t, driver)
	return txdb.New(t, driver)
}

func MustExec(t *testing.T, ds sqlutil.DataSource, stmt string, args ...interface{}) {
	ctx := testutils.Context(t)
	require.NoError(t, utils.JustError(ds.ExecContext(ctx, stmt, args...)))
}

func MustCount(t *testing.T, db *sqlx.DB, stmt string, args ...interface{}) (cnt int) {
	require.NoError(t, db.Get(&cnt, stmt, args...))
	return
}

func TestMain(m *testing.M) {
	initOnce.Do(func() {
		driver, err = initUnitTestDB(context.Background())
	})
	m.Run()

	//panic("wtf")

	//	exitCode := m.Run()

}
