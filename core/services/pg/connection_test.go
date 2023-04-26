package pg

import (
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

func Test_disallowReplica(t *testing.T) {

	testutils.SkipShortDB(t)
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), uuid.NewV4().String())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err = db.Exec("SET session_replication_role= 'origin'")
	require.NoError(t, err)
	err = disallowReplica(db)
	require.NoError(t, err)

	_, err = db.Exec("SET session_replication_role= 'replica'")
	require.NoError(t, err)
	err = disallowReplica(db)
	require.Error(t, err, "replica role should be disallowed")

	_, err = db.Exec("SET session_replication_role= 'not_valid_role'")
	require.Error(t, err)

}
