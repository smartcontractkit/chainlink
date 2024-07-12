package pgtest

import (
	"context"
	"database/sql"
	"flag"
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

var testDB *sqlx.DB

func initUnitTestDB(ctx context.Context) (string, error) {
	testing.Init()
	if !flag.Parsed() {
		flag.Parse()
	}
	if testing.Short() {
		// -short tests don't need a DB
		return "", nil
	}
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
		testURL, err := testdb.CreateOrReplace(*testDBURL, "unit_test_base"+uuid[:16], false)
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
	driver := string(dialects.TransactionWrappedPostgres)
	txdb.RegisterTestDB(driver, testDBURL)
	fmt.Println("Test database URL:", testDBURL.String())

	return driver, nil
}

var initOnce sync.Once
var driver string
var err error

func NewSqlxDB(t testing.TB) *sqlx.DB {

	initOnce.Do(func() {
		driver, err = initUnitTestDB(testutils.Context(t))
	})

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
