package cre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// orgResolverCacheTable is the durable owner->orgID mapping table backing the
// OrgResolver cache. See migration 0305_org_resolver_cache.sql.
const orgResolverCacheTable = "cre.org_resolver_cache"

// orgResolverStore is a Postgres-backed implementation of orgresolver.Cache.
type orgResolverStore struct {
	ds sqlutil.DataSource
}

// NewOrgResolverStore creates a durable cache store for the OrgResolver.
func NewOrgResolverStore(ds sqlutil.DataSource) *orgResolverStore {
	return &orgResolverStore{ds: ds}
}

// Get returns the cached entry for owner. ok is false if no entry exists.
func (s *orgResolverStore) Get(ctx context.Context, owner string) (orgresolver.CacheEntry, bool, error) {
	const q = `SELECT org_id, updated_at FROM ` + orgResolverCacheTable + ` WHERE workflow_owner = $1`
	var row struct {
		OrgID     string    `db:"org_id"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := s.ds.GetContext(ctx, &row, q, owner); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orgresolver.CacheEntry{}, false, nil
		}
		return orgresolver.CacheEntry{}, false, fmt.Errorf("failed to get cached org for owner %s: %w", owner, err)
	}
	return orgresolver.CacheEntry{OrgID: row.OrgID, RefreshedAt: row.UpdatedAt}, true, nil
}

// Set stores or updates the mapping for owner.
func (s *orgResolverStore) Set(ctx context.Context, owner string, entry orgresolver.CacheEntry) error {
	const q = `
INSERT INTO ` + orgResolverCacheTable + ` (workflow_owner, org_id, updated_at)
VALUES ($1, $2, $3)
ON CONFLICT (workflow_owner) DO UPDATE SET org_id = EXCLUDED.org_id, updated_at = EXCLUDED.updated_at`
	if _, err := s.ds.ExecContext(ctx, q, owner, entry.OrgID, entry.RefreshedAt); err != nil {
		return fmt.Errorf("failed to upsert org for owner %s: %w", owner, err)
	}
	return nil
}

var _ orgresolver.Cache = (*orgResolverStore)(nil)
