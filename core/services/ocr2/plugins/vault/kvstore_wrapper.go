package vault

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// KVStoreWrapper provides a migration-aware layer over the underlying KVStore.
// A single instance is created per plugin function call (matching the existing
// one-store-per-call pattern), and orgID/workflowOwner are passed per operation
// since a batch may contain requests from different owners.
//
// When migrationEnabled is false the wrapper is a pure pass-through: every
// call goes directly to the inner store and the orgID/workflowOwner parameters
// are ignored. Gate this with cresettings.Default.VaultOrgIdAsSecretOwnerEnabled.
type KVStoreWrapper struct {
	store            WriteKVStore
	adapter          *ownerMigrationAdapter
	migrationEnabled bool
}

// requestScopedKVStore binds a wrapper to the orgID/workflowOwner for a single
// top-level plugin request while preserving the existing ReadKVStore /
// WriteKVStore interfaces used throughout plugin.go.
type requestScopedKVStore struct {
	wrapper       *KVStoreWrapper
	orgID         string
	workflowOwner string
}

// NewKVStoreWrapper creates a wrapper around the given store.
// When migrationEnabled is true, an internal ownerMigrationAdapter handles
// the transition from workflow_owner-keyed entries to org_id-keyed entries.
// When false, all operations pass through directly to the inner store.
func NewKVStoreWrapper(store WriteKVStore, migrationEnabled bool, lggr logger.Logger) *KVStoreWrapper {
	w := &KVStoreWrapper{
		store:            store,
		migrationEnabled: migrationEnabled,
	}
	if migrationEnabled {
		w.adapter = newOwnerMigrationAdapter(store, lggr)
	}
	return w
}

// WithRequest returns a store view bound to the top-level request's orgID and
// workflowOwner. When migration is enabled, owner-scoped operations are routed
// through the migration adapter using the bound orgID/workflowOwner while
// preserving the plugin's existing store interface usage.
func (w *KVStoreWrapper) WithRequest(orgID, workflowOwner string) WriteKVStore {
	return &requestScopedKVStore{
		wrapper:       w,
		orgID:         orgID,
		workflowOwner: workflowOwner,
	}
}

// GetSecret tries orgID first, falling back to workflowOwner for legacy entries.
// When migration is disabled, delegates directly to the inner store using id as-is.
func (w *KVStoreWrapper) GetSecret(ctx context.Context, id *vault.SecretIdentifier, orgID, workflowOwner string) (*vault.StoredSecret, error) {
	if !w.migrationEnabled {
		return w.store.GetSecret(ctx, id)
	}
	return w.adapter.getSecret(ctx, id, orgID, workflowOwner)
}

// GetMetadata merges metadata from both orgID and workflowOwner, deduplicating
// by namespace::key and rewriting all Owner fields to orgID.
// When migration is disabled, delegates directly to the inner store using orgID.
//
// The merged count cannot exceed the per-owner secret limit: deduplication by
// namespace::key collapses entries that exist under both owners (transient
// mid-migration state) into a single entry, so the result reflects the true
// number of unique secrets the owner has.
func (w *KVStoreWrapper) GetMetadata(ctx context.Context, orgID, workflowOwner string) (*vault.StoredMetadata, error) {
	if !w.migrationEnabled {
		return w.store.GetMetadata(ctx, orgID)
	}
	return w.adapter.getMetadata(ctx, orgID, workflowOwner)
}

// GetSecretIdentifiersCountForOwner returns the count of unique secrets across
// both orgID and workflowOwner after deduplication.
// When migration is disabled, delegates directly to the inner store using orgID.
func (w *KVStoreWrapper) GetSecretIdentifiersCountForOwner(ctx context.Context, orgID, workflowOwner string) (int, error) {
	if !w.migrationEnabled {
		return w.store.GetSecretIdentifiersCountForOwner(ctx, orgID)
	}
	return w.adapter.getSecretIdentifiersCountForOwner(ctx, orgID, workflowOwner)
}

// WriteSecret writes the secret under orgID. If a legacy entry exists under
// workflowOwner with the same namespace/key, it is deleted (lazy migration).
// When migration is disabled, delegates directly to the inner store.
func (w *KVStoreWrapper) WriteSecret(ctx context.Context, id *vault.SecretIdentifier, secret *vault.StoredSecret, orgID, workflowOwner string) error {
	if !w.migrationEnabled {
		return w.store.WriteSecret(ctx, id, secret)
	}
	return w.adapter.writeSecret(ctx, id, secret, orgID, workflowOwner)
}

// WriteMetadata writes metadata under orgID.
// When migration is disabled, delegates directly to the inner store.
func (w *KVStoreWrapper) WriteMetadata(ctx context.Context, orgID string, metadata *vault.StoredMetadata) error {
	if !w.migrationEnabled {
		return w.store.WriteMetadata(ctx, orgID, metadata)
	}
	return w.adapter.writeMetadata(ctx, orgID, metadata)
}

// DeleteSecret deletes the secret from orgID if present, falling back to
// workflowOwner for legacy entries. If the secret exists under both owners
// (transient mid-migration state), both entries are deleted.
// When migration is disabled, delegates directly to the inner store.
func (w *KVStoreWrapper) DeleteSecret(ctx context.Context, id *vault.SecretIdentifier, orgID, workflowOwner string) error {
	if !w.migrationEnabled {
		return w.store.DeleteSecret(ctx, id)
	}
	return w.adapter.deleteSecret(ctx, id, orgID, workflowOwner)
}

// GetPendingQueue is always a pass-through (pending queue is not owner-scoped).
func (w *KVStoreWrapper) GetPendingQueue(ctx context.Context) ([]*vault.StoredPendingQueueItem, error) {
	return w.store.GetPendingQueue(ctx)
}

// WritePendingQueue is always a pass-through (pending queue is not owner-scoped).
func (w *KVStoreWrapper) WritePendingQueue(ctx context.Context, pending []*vault.StoredPendingQueueItem) error {
	return w.store.WritePendingQueue(ctx, pending)
}

func (s *requestScopedKVStore) effectiveOwner(owner string) string {
	if s.wrapper.migrationEnabled && s.orgID != "" {
		return s.orgID
	}
	return owner
}

func (s *requestScopedKVStore) GetSecret(ctx context.Context, id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	orgID := ""
	if id != nil {
		orgID = s.effectiveOwner(id.Owner)
	}
	return s.wrapper.GetSecret(ctx, id, orgID, s.workflowOwner)
}

func (s *requestScopedKVStore) GetMetadata(ctx context.Context, owner string) (*vault.StoredMetadata, error) {
	return s.wrapper.GetMetadata(ctx, s.effectiveOwner(owner), s.workflowOwner)
}

func (s *requestScopedKVStore) GetSecretIdentifiersCountForOwner(ctx context.Context, owner string) (int, error) {
	return s.wrapper.GetSecretIdentifiersCountForOwner(ctx, s.effectiveOwner(owner), s.workflowOwner)
}

func (s *requestScopedKVStore) WriteSecret(ctx context.Context, id *vault.SecretIdentifier, secret *vault.StoredSecret) error {
	orgID := ""
	if id != nil {
		orgID = s.effectiveOwner(id.Owner)
	}
	return s.wrapper.WriteSecret(ctx, id, secret, orgID, s.workflowOwner)
}

func (s *requestScopedKVStore) WriteMetadata(ctx context.Context, owner string, metadata *vault.StoredMetadata) error {
	return s.wrapper.WriteMetadata(ctx, s.effectiveOwner(owner), metadata)
}

func (s *requestScopedKVStore) DeleteSecret(ctx context.Context, id *vault.SecretIdentifier) error {
	orgID := ""
	if id != nil {
		orgID = s.effectiveOwner(id.Owner)
	}
	return s.wrapper.DeleteSecret(ctx, id, orgID, s.workflowOwner)
}

func (s *requestScopedKVStore) GetPendingQueue(ctx context.Context) ([]*vault.StoredPendingQueueItem, error) {
	return s.wrapper.GetPendingQueue(ctx)
}

func (s *requestScopedKVStore) WritePendingQueue(ctx context.Context, pending []*vault.StoredPendingQueueItem) error {
	return s.wrapper.WritePendingQueue(ctx, pending)
}

// ownerMigrationAdapter handles the migration of secrets from workflow_owner-keyed
// entries to org_id-keyed entries. It performs dual-lookup reads, org_id-based
// writes, lazy migration on update, metadata merge for list, and dual-owner
// deletion.
type ownerMigrationAdapter struct {
	store WriteKVStore
	lggr  logger.Logger
}

func newOwnerMigrationAdapter(store WriteKVStore, lggr logger.Logger) *ownerMigrationAdapter {
	return &ownerMigrationAdapter{store: store, lggr: lggr}
}

func (a *ownerMigrationAdapter) getSecret(ctx context.Context, id *vault.SecretIdentifier, orgID, workflowOwner string) (*vault.StoredSecret, error) {
	if id == nil {
		return a.store.GetSecret(ctx, id)
	}

	orgIDSid := withOwner(id, orgID)
	secret, err := a.store.GetSecret(ctx, orgIDSid)
	if err != nil {
		return nil, err
	}
	if secret != nil {
		return secret, nil
	}

	if !needsMigration(orgID, workflowOwner) {
		return nil, nil
	}

	woSid := withOwner(id, workflowOwner)
	return a.store.GetSecret(ctx, woSid)
}

func (a *ownerMigrationAdapter) getMetadata(ctx context.Context, orgID, workflowOwner string) (*vault.StoredMetadata, error) {
	orgMd, err := a.store.GetMetadata(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if !needsMigration(orgID, workflowOwner) {
		return orgMd, nil
	}

	woMd, err := a.store.GetMetadata(ctx, workflowOwner)
	if err != nil {
		return nil, err
	}

	return mergeMetadata(orgMd, woMd, orgID, a.lggr), nil
}

func (a *ownerMigrationAdapter) getSecretIdentifiersCountForOwner(ctx context.Context, orgID, workflowOwner string) (int, error) {
	md, err := a.getMetadata(ctx, orgID, workflowOwner)
	if err != nil {
		return 0, err
	}
	if md == nil {
		return 0, nil
	}
	return len(md.SecretIdentifiers), nil
}

func (a *ownerMigrationAdapter) writeSecret(ctx context.Context, id *vault.SecretIdentifier, secret *vault.StoredSecret, orgID, workflowOwner string) error {
	if id == nil {
		return a.store.WriteSecret(ctx, id, secret)
	}

	orgIDSid := withOwner(id, orgID)
	if err := a.store.WriteSecret(ctx, orgIDSid, secret); err != nil {
		return err
	}

	if !needsMigration(orgID, workflowOwner) {
		return nil
	}

	woSid := withOwner(id, workflowOwner)
	legacySecret, err := a.store.GetSecret(ctx, woSid)
	if err != nil {
		return fmt.Errorf("failed to check for legacy entry during write: %w", err)
	}
	if legacySecret != nil {
		if err := a.store.DeleteSecret(ctx, woSid); err != nil {
			return fmt.Errorf("failed to delete legacy entry during migration: %w", err)
		}
	}

	return nil
}

func (a *ownerMigrationAdapter) writeMetadata(ctx context.Context, orgID string, metadata *vault.StoredMetadata) error {
	return a.store.WriteMetadata(ctx, orgID, metadata)
}

func (a *ownerMigrationAdapter) deleteSecret(ctx context.Context, id *vault.SecretIdentifier, orgID, workflowOwner string) error {
	if id == nil {
		return a.store.DeleteSecret(ctx, id)
	}

	orgIDSid := withOwner(id, orgID)
	orgSecret, err := a.store.GetSecret(ctx, orgIDSid)
	if err != nil {
		return fmt.Errorf("failed to check org_id entry for deletion: %w", err)
	}
	if orgSecret != nil {
		if err := a.store.DeleteSecret(ctx, orgIDSid); err != nil {
			return fmt.Errorf("failed to delete org_id entry: %w", err)
		}
		if needsMigration(orgID, workflowOwner) {
			woSid := withOwner(id, workflowOwner)
			woSecret, woErr := a.store.GetSecret(ctx, woSid)
			if woErr != nil {
				return fmt.Errorf("failed to check legacy entry after org_id deletion: %w", woErr)
			}
			if woSecret != nil {
				if woErr = a.store.DeleteSecret(ctx, woSid); woErr != nil {
					return fmt.Errorf("failed to clean up legacy entry after org_id deletion: %w", woErr)
				}
			}
		}
		return nil
	}

	if needsMigration(orgID, workflowOwner) {
		woSid := withOwner(id, workflowOwner)
		woSecret, woErr := a.store.GetSecret(ctx, woSid)
		if woErr != nil {
			return fmt.Errorf("failed to check legacy entry for deletion: %w", woErr)
		}
		if woSecret != nil {
			return a.store.DeleteSecret(ctx, woSid)
		}
	}

	// Not found under either owner — delegate to inner which will produce
	// the appropriate error from metadata removal.
	return a.store.DeleteSecret(ctx, orgIDSid)
}

// --- shared helpers ---

func withOwner(id *vault.SecretIdentifier, owner string) *vault.SecretIdentifier {
	return &vault.SecretIdentifier{
		Key:       id.Key,
		Namespace: id.Namespace,
		Owner:     owner,
	}
}

func needsMigration(orgID, workflowOwner string) bool {
	return workflowOwner != "" && workflowOwner != orgID
}

// mergeMetadata combines metadata from org_id and workflow_owner, deduplicating
// by namespace::key and rewriting all Owner fields to orgID.
func mergeMetadata(orgMd, woMd *vault.StoredMetadata, orgID string, lggr logger.Logger) *vault.StoredMetadata {
	if orgMd == nil && woMd == nil {
		return nil
	}

	seen := map[string]bool{}
	var merged []*vault.SecretIdentifier

	addEntries := func(md *vault.StoredMetadata, source string) {
		if md == nil {
			return
		}
		for _, id := range md.SecretIdentifiers {
			dk := deduplicationKey(id)
			if seen[dk] {
				lggr.Criticalw(
					"duplicate secret identifier found during owner migration metadata merge",
					"orgID", orgID,
					"duplicateKey", dk,
					"namespace", id.Namespace,
					"key", id.Key,
					"owner", id.Owner,
					"source", source,
				)
				continue
			}
			seen[dk] = true
			merged = append(merged, &vault.SecretIdentifier{
				Key:       id.Key,
				Namespace: id.Namespace,
				Owner:     orgID,
			})
		}
	}

	// org_id entries take priority in deduplication.
	addEntries(orgMd, "org_id")
	addEntries(woMd, "workflow_owner")

	return &vault.StoredMetadata{
		SecretIdentifiers: merged,
	}
}

func deduplicationKey(id *vault.SecretIdentifier) string {
	namespace := id.Namespace
	if namespace == "" {
		namespace = vaulttypes.DefaultNamespace
	}
	return namespace + "::" + id.Key
}
