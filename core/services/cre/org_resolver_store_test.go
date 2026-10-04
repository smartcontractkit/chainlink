package cre

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

// countingDataSource wraps a DataSource and counts every query issued through it,
// so tests can assert on cache hits/misses without inspecting store internals.
func countingDataSource(t *testing.T, ds sqlutil.DataSource) (sqlutil.DataSource, *int32) {
	t.Helper()
	var queries int32
	hook := func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
		atomic.AddInt32(&queries, 1)
		return do(ctx)
	}
	return sqlutil.WrapDataSource(ds, logger.Test(t), hook), &queries
}

func Test_OrgResolverStore_GetSet(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	store := NewOrgResolverStore(db)
	ctx := context.Background()

	_, ok, err := store.Get(ctx, "owner-1")
	require.NoError(t, err)
	assert.False(t, ok)

	entry := orgresolver.CacheEntry{OrgID: "org-1", RefreshedAt: time.Now().UTC().Truncate(time.Microsecond)}
	require.NoError(t, store.Set(ctx, "owner-1", entry))

	got, ok, err := store.Get(ctx, "owner-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, entry.OrgID, got.OrgID)
	assert.WithinDuration(t, entry.RefreshedAt, got.RefreshedAt, time.Second)
}

func Test_OrgResolverStore_Set_UpsertsOnOrgIDChange(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	store := NewOrgResolverStore(db)
	ctx := context.Background()

	require.NoError(t, store.Set(ctx, "owner-1", orgresolver.CacheEntry{OrgID: "org-1", RefreshedAt: time.Now()}))
	require.NoError(t, store.Set(ctx, "owner-1", orgresolver.CacheEntry{OrgID: "org-2", RefreshedAt: time.Now()}))

	var orgID string
	require.NoError(t, db.GetContext(ctx, &orgID,
		`SELECT org_id FROM cre.org_resolver_cache WHERE workflow_owner = $1`, "owner-1"))
	assert.Equal(t, "org-2", orgID)
}

func Test_OrgResolverStore_Get_CachesAfterFirstRead(t *testing.T) {
	t.Parallel()
	rawDB := pgtest.NewSqlxDB(t)
	ctx := context.Background()

	entry := orgresolver.CacheEntry{OrgID: "org-1", RefreshedAt: time.Now().UTC().Truncate(time.Microsecond)}
	require.NoError(t, NewOrgResolverStore(rawDB).Set(ctx, "owner-1", entry))

	ds, queries := countingDataSource(t, rawDB)
	store := NewOrgResolverStore(ds)

	got, ok, err := store.Get(ctx, "owner-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, entry.OrgID, got.OrgID)
	assert.Equal(t, int32(1), atomic.LoadInt32(queries), "expected exactly one DB query on cache miss")

	got, ok, err = store.Get(ctx, "owner-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, entry.OrgID, got.OrgID)
	assert.Equal(t, int32(1), atomic.LoadInt32(queries), "expected no additional DB query on cache hit")
}

func Test_OrgResolverStore_Set_PopulatesCacheWithoutExtraRead(t *testing.T) {
	t.Parallel()
	rawDB := pgtest.NewSqlxDB(t)
	ctx := context.Background()

	ds, queries := countingDataSource(t, rawDB)
	store := NewOrgResolverStore(ds)

	entry := orgresolver.CacheEntry{OrgID: "org-1", RefreshedAt: time.Now().UTC().Truncate(time.Microsecond)}
	require.NoError(t, store.Set(ctx, "owner-1", entry))
	assert.Equal(t, int32(1), atomic.LoadInt32(queries), "Set should issue exactly one query")

	got, ok, err := store.Get(ctx, "owner-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, entry.OrgID, got.OrgID)
	assert.Equal(t, int32(1), atomic.LoadInt32(queries), "Get after Set should be served from cache")
}

func Test_OrgResolverStore_Get_CacheIsPerOwner(t *testing.T) {
	t.Parallel()
	rawDB := pgtest.NewSqlxDB(t)
	ctx := context.Background()

	seedStore := NewOrgResolverStore(rawDB)
	require.NoError(t, seedStore.Set(ctx, "owner-1", orgresolver.CacheEntry{OrgID: "org-1", RefreshedAt: time.Now()}))
	require.NoError(t, seedStore.Set(ctx, "owner-2", orgresolver.CacheEntry{OrgID: "org-2", RefreshedAt: time.Now()}))

	ds, queries := countingDataSource(t, rawDB)
	store := NewOrgResolverStore(ds)

	_, _, err := store.Get(ctx, "owner-1")
	require.NoError(t, err)
	_, _, err = store.Get(ctx, "owner-2")
	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(queries))

	_, _, err = store.Get(ctx, "owner-1")
	require.NoError(t, err)
	_, _, err = store.Get(ctx, "owner-2")
	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(queries), "second round of Gets should be fully served from cache")
}

func Test_OrgResolverStore_Get_DoesNotCacheMisses(t *testing.T) {
	t.Parallel()
	rawDB := pgtest.NewSqlxDB(t)
	ctx := context.Background()

	ds, queries := countingDataSource(t, rawDB)
	store := NewOrgResolverStore(ds)

	_, ok, err := store.Get(ctx, "missing-owner")
	require.NoError(t, err)
	assert.False(t, ok)

	_, ok, err = store.Get(ctx, "missing-owner")
	require.NoError(t, err)
	assert.False(t, ok)

	assert.Equal(t, int32(2), atomic.LoadInt32(queries), "a cache miss should not be remembered, each Get should hit the DB")
}
