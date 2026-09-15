package cre

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

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
