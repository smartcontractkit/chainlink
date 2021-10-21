package versioning

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/stretchr/testify/require"
)

func TestORM_NodeVersion_UpsertNodeVersion(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db)

	err := orm.UpsertNodeVersion(NewNodeVersion("9.9.8"))
	require.NoError(t, err)

	ver, err := orm.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Equal(t, "9.9.8", ver.Version)
	require.NotZero(t, ver.CreatedAt)

	// Testing Upsert
	require.NoError(t, orm.UpsertNodeVersion(NewNodeVersion("9.9.8")))
	require.NoError(t, orm.UpsertNodeVersion(NewNodeVersion("9.9.7")))
	require.NoError(t, orm.UpsertNodeVersion(NewNodeVersion("9.9.9")))

	ver, err = orm.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Equal(t, "9.9.9", ver.Version)
}

func TestORM_NodeVersion_FindLatestNodeVersion(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db)

	// Not Found
	_, err := orm.FindLatestNodeVersion()
	require.Error(t, err)

	err = orm.UpsertNodeVersion(NewNodeVersion("9.9.8"))
	require.NoError(t, err)

	ver, err := orm.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Equal(t, "9.9.8", ver.Version)
	require.NotZero(t, ver.CreatedAt)
}
