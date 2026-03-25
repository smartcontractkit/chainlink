package vault

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// KVStoreWrapper provides a migration-aware layer over the underlying KVStore.
// A single instance is created per plugin function call (matching the existing
// one-store-per-call pattern), and orgID/workflowOwner are passed per operation
// since a batch may contain requests from different owners.
type KVStoreWrapper struct {
	adapter *ownerMigrationAdapter
}

// NewKVStoreWrapper creates a wrapper that delegates to an internal
// ownerMigrationAdapter for transitioning secrets from workflow_owner-keyed
// entries to org_id-keyed entries.
func NewKVStoreWrapper(store WriteKVStore, lggr logger.Logger) *KVStoreWrapper {
	return &KVStoreWrapper{
		adapter: newOwnerMigrationAdapter(store, lggr),
	}
}

// GetSecret tries orgID first, falling back to workflowOwner for legacy entries.
func (w *KVStoreWrapper) GetSecret(id *vault.SecretIdentifier, orgID, workflowOwner string) (*vault.StoredSecret, error) {
	return w.adapter.getSecret(id, orgID, workflowOwner)
}

// GetMetadata merges metadata from both orgID and workflowOwner, deduplicating
// by namespace::key and rewriting all Owner fields to orgID.
//
// The merged count cannot exceed the per-owner secret limit: deduplication by
// namespace::key collapses entries that exist under both owners (transient
// mid-migration state) into a single entry, so the result reflects the true
// number of unique secrets the owner has.
func (w *KVStoreWrapper) GetMetadata(orgID, workflowOwner string) (*vault.StoredMetadata, error) {
	return w.adapter.getMetadata(orgID, workflowOwner)
}

// GetSecretIdentifiersCountForOwner returns the count of unique secrets across
// both orgID and workflowOwner after deduplication.
func (w *KVStoreWrapper) GetSecretIdentifiersCountForOwner(orgID, workflowOwner string) (int, error) {
	return w.adapter.getSecretIdentifiersCountForOwner(orgID, workflowOwner)
}

// WriteSecret writes the secret under orgID. If a legacy entry exists under
// workflowOwner with the same namespace/key, it is deleted (lazy migration).
func (w *KVStoreWrapper) WriteSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret, orgID, workflowOwner string) error {
	return w.adapter.writeSecret(id, secret, orgID, workflowOwner)
}

// WriteMetadata writes metadata under orgID.
func (w *KVStoreWrapper) WriteMetadata(orgID string, metadata *vault.StoredMetadata) error {
	return w.adapter.writeMetadata(orgID, metadata)
}

// DeleteSecret deletes the secret from orgID if present, falling back to
// workflowOwner for legacy entries. If the secret exists under both owners
// (transient mid-migration state), both entries are deleted.
func (w *KVStoreWrapper) DeleteSecret(id *vault.SecretIdentifier, orgID, workflowOwner string) error {
	return w.adapter.deleteSecret(id, orgID, workflowOwner)
}

// GetPendingQueue is a pass-through (pending queue is not owner-scoped).
func (w *KVStoreWrapper) GetPendingQueue() ([]*vault.StoredPendingQueueItem, error) {
	return w.adapter.store.GetPendingQueue()
}

// WritePendingQueue is a pass-through (pending queue is not owner-scoped).
func (w *KVStoreWrapper) WritePendingQueue(pending []*vault.StoredPendingQueueItem) error {
	return w.adapter.store.WritePendingQueue(pending)
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

func (a *ownerMigrationAdapter) getSecret(id *vault.SecretIdentifier, orgID, workflowOwner string) (*vault.StoredSecret, error) {
	if id == nil {
		return a.store.GetSecret(id)
	}

	orgIDSid := withOwner(id, orgID)
	secret, err := a.store.GetSecret(orgIDSid)
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
	return a.store.GetSecret(woSid)
}

func (a *ownerMigrationAdapter) getMetadata(orgID, workflowOwner string) (*vault.StoredMetadata, error) {
	orgMd, err := a.store.GetMetadata(orgID)
	if err != nil {
		return nil, err
	}

	if !needsMigration(orgID, workflowOwner) {
		return orgMd, nil
	}

	woMd, err := a.store.GetMetadata(workflowOwner)
	if err != nil {
		return nil, err
	}

	return mergeMetadata(orgMd, woMd, orgID, a.lggr), nil
}

func (a *ownerMigrationAdapter) getSecretIdentifiersCountForOwner(orgID, workflowOwner string) (int, error) {
	md, err := a.getMetadata(orgID, workflowOwner)
	if err != nil {
		return 0, err
	}
	if md == nil {
		return 0, nil
	}
	return len(md.SecretIdentifiers), nil
}

func (a *ownerMigrationAdapter) writeSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret, orgID, workflowOwner string) error {
	if id == nil {
		return a.store.WriteSecret(id, secret)
	}

	orgIDSid := withOwner(id, orgID)
	if err := a.store.WriteSecret(orgIDSid, secret); err != nil {
		return err
	}

	if !needsMigration(orgID, workflowOwner) {
		return nil
	}

	woSid := withOwner(id, workflowOwner)
	legacySecret, err := a.store.GetSecret(woSid)
	if err != nil {
		return fmt.Errorf("failed to check for legacy entry during write: %w", err)
	}
	if legacySecret != nil {
		if err := a.store.DeleteSecret(woSid); err != nil {
			return fmt.Errorf("failed to delete legacy entry during migration: %w", err)
		}
	}

	return nil
}

func (a *ownerMigrationAdapter) writeMetadata(orgID string, metadata *vault.StoredMetadata) error {
	return a.store.WriteMetadata(orgID, metadata)
}

func (a *ownerMigrationAdapter) deleteSecret(id *vault.SecretIdentifier, orgID, workflowOwner string) error {
	if id == nil {
		return a.store.DeleteSecret(id)
	}

	orgIDSid := withOwner(id, orgID)
	orgSecret, err := a.store.GetSecret(orgIDSid)
	if err != nil {
		return fmt.Errorf("failed to check org_id entry for deletion: %w", err)
	}
	if orgSecret != nil {
		if err := a.store.DeleteSecret(orgIDSid); err != nil {
			return fmt.Errorf("failed to delete org_id entry: %w", err)
		}
		if needsMigration(orgID, workflowOwner) {
			woSid := withOwner(id, workflowOwner)
			woSecret, woErr := a.store.GetSecret(woSid)
			if woErr != nil {
				return fmt.Errorf("failed to check legacy entry after org_id deletion: %w", woErr)
			}
			if woSecret != nil {
				if woErr = a.store.DeleteSecret(woSid); woErr != nil {
					return fmt.Errorf("failed to clean up legacy entry after org_id deletion: %w", woErr)
				}
			}
		}
		return nil
	}

	if needsMigration(orgID, workflowOwner) {
		woSid := withOwner(id, workflowOwner)
		woSecret, woErr := a.store.GetSecret(woSid)
		if woErr != nil {
			return fmt.Errorf("failed to check legacy entry for deletion: %w", woErr)
		}
		if woSecret != nil {
			return a.store.DeleteSecret(woSid)
		}
	}

	// Not found under either owner — delegate to inner which will produce
	// the appropriate error from metadata removal.
	return a.store.DeleteSecret(orgIDSid)
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
