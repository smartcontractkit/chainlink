package vault

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// OwnerMigrationReadStore is a transparent adapter that wraps a ReadKVStore and
// implements dual-owner lookup: it tries org_id first, falling back to
// workflow_owner for legacy entries. It is used during the migration period
// where secrets may be keyed under either the new org_id or the legacy
// workflow_owner address.
type OwnerMigrationReadStore struct {
	inner         ReadKVStore
	orgID         string
	workflowOwner string
}

var _ ReadKVStore = (*OwnerMigrationReadStore)(nil)

func NewOwnerMigrationReadStore(inner ReadKVStore, orgID, workflowOwner string) *OwnerMigrationReadStore {
	return &OwnerMigrationReadStore{
		inner:         inner,
		orgID:         orgID,
		workflowOwner: workflowOwner,
	}
}

func (s *OwnerMigrationReadStore) GetSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	return getSecretWithFallback(s.inner, id, s.orgID, s.workflowOwner)
}

func (s *OwnerMigrationReadStore) GetMetadata(owner string) (*vault.StoredMetadata, error) {
	return getMetadataWithFallback(s.inner, s.orgID, s.workflowOwner)
}

func (s *OwnerMigrationReadStore) GetSecretIdentifiersCountForOwner(owner string) (int, error) {
	md, err := s.GetMetadata(owner)
	if err != nil {
		return 0, err
	}
	if md == nil {
		return 0, nil
	}
	return len(md.SecretIdentifiers), nil
}

func (s *OwnerMigrationReadStore) GetPendingQueue() ([]*vault.StoredPendingQueueItem, error) {
	return s.inner.GetPendingQueue()
}

// OwnerMigrationWriteStore is a transparent adapter that wraps a WriteKVStore
// and handles the migration of secrets from workflow_owner-keyed entries to
// org_id-keyed entries. Writes always go under org_id. Updates and deletes
// perform lazy migration by cleaning up legacy workflow_owner entries.
type OwnerMigrationWriteStore struct {
	inner         WriteKVStore
	orgID         string
	workflowOwner string
}

var _ WriteKVStore = (*OwnerMigrationWriteStore)(nil)

func NewOwnerMigrationWriteStore(inner WriteKVStore, orgID, workflowOwner string) *OwnerMigrationWriteStore {
	return &OwnerMigrationWriteStore{
		inner:         inner,
		orgID:         orgID,
		workflowOwner: workflowOwner,
	}
}

func (s *OwnerMigrationWriteStore) GetSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	return getSecretWithFallback(s.inner, id, s.orgID, s.workflowOwner)
}

func (s *OwnerMigrationWriteStore) GetMetadata(owner string) (*vault.StoredMetadata, error) {
	return getMetadataWithFallback(s.inner, s.orgID, s.workflowOwner)
}

func (s *OwnerMigrationWriteStore) GetSecretIdentifiersCountForOwner(owner string) (int, error) {
	md, err := s.GetMetadata(owner)
	if err != nil {
		return 0, err
	}
	if md == nil {
		return 0, nil
	}
	return len(md.SecretIdentifiers), nil
}

func (s *OwnerMigrationWriteStore) GetPendingQueue() ([]*vault.StoredPendingQueueItem, error) {
	return s.inner.GetPendingQueue()
}

func (s *OwnerMigrationWriteStore) WriteSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret) error {
	if id == nil {
		return s.inner.WriteSecret(id, secret)
	}

	orgIDSid := withOwner(id, s.orgID)
	if err := s.inner.WriteSecret(orgIDSid, secret); err != nil {
		return err
	}

	if !needsMigration(s.orgID, s.workflowOwner) {
		return nil
	}

	woSid := withOwner(id, s.workflowOwner)
	legacySecret, err := s.inner.GetSecret(woSid)
	if err != nil {
		return fmt.Errorf("failed to check for legacy entry during write: %w", err)
	}
	if legacySecret != nil {
		if err := s.inner.DeleteSecret(woSid); err != nil {
			return fmt.Errorf("failed to delete legacy entry during migration: %w", err)
		}
	}

	return nil
}

func (s *OwnerMigrationWriteStore) WriteMetadata(owner string, metadata *vault.StoredMetadata) error {
	return s.inner.WriteMetadata(s.orgID, metadata)
}

func (s *OwnerMigrationWriteStore) DeleteSecret(id *vault.SecretIdentifier) error {
	if id == nil {
		return s.inner.DeleteSecret(id)
	}

	orgIDSid := withOwner(id, s.orgID)
	orgSecret, err := s.inner.GetSecret(orgIDSid)
	if err != nil {
		return fmt.Errorf("failed to check org_id entry for deletion: %w", err)
	}
	if orgSecret != nil {
		if err := s.inner.DeleteSecret(orgIDSid); err != nil {
			return fmt.Errorf("failed to delete org_id entry: %w", err)
		}
		// Also clean up any legacy entry if it exists under workflow_owner.
		if needsMigration(s.orgID, s.workflowOwner) {
			woSid := withOwner(id, s.workflowOwner)
			woSecret, woErr := s.inner.GetSecret(woSid)
			if woErr != nil {
				return fmt.Errorf("failed to check legacy entry after org_id deletion: %w", woErr)
			}
			if woSecret != nil {
				if woErr = s.inner.DeleteSecret(woSid); woErr != nil {
					return fmt.Errorf("failed to clean up legacy entry after org_id deletion: %w", woErr)
				}
			}
		}
		return nil
	}

	if needsMigration(s.orgID, s.workflowOwner) {
		woSid := withOwner(id, s.workflowOwner)
		woSecret, woErr := s.inner.GetSecret(woSid)
		if woErr != nil {
			return fmt.Errorf("failed to check legacy entry for deletion: %w", woErr)
		}
		if woSecret != nil {
			return s.inner.DeleteSecret(woSid)
		}
	}

	// Not found under either owner — delegate to inner which will produce
	// the appropriate error from metadata removal.
	return s.inner.DeleteSecret(orgIDSid)
}

func (s *OwnerMigrationWriteStore) WritePendingQueue(pending []*vault.StoredPendingQueueItem) error {
	return s.inner.WritePendingQueue(pending)
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

func getSecretWithFallback(inner ReadKVStore, id *vault.SecretIdentifier, orgID, workflowOwner string) (*vault.StoredSecret, error) {
	if id == nil {
		return inner.GetSecret(id)
	}

	orgIDSid := withOwner(id, orgID)
	secret, err := inner.GetSecret(orgIDSid)
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
	return inner.GetSecret(woSid)
}

func getMetadataWithFallback(inner ReadKVStore, orgID, workflowOwner string) (*vault.StoredMetadata, error) {
	orgMd, err := inner.GetMetadata(orgID)
	if err != nil {
		return nil, err
	}

	if !needsMigration(orgID, workflowOwner) {
		return orgMd, nil
	}

	woMd, err := inner.GetMetadata(workflowOwner)
	if err != nil {
		return nil, err
	}

	return mergeMetadata(orgMd, woMd, orgID), nil
}

// mergeMetadata combines metadata from org_id and workflow_owner, deduplicating
// by namespace::key and rewriting all Owner fields to orgID.
func mergeMetadata(orgMd, woMd *vault.StoredMetadata, orgID string) *vault.StoredMetadata {
	if orgMd == nil && woMd == nil {
		return nil
	}

	seen := map[string]bool{}
	var merged []*vault.SecretIdentifier

	addEntries := func(md *vault.StoredMetadata) {
		if md == nil {
			return
		}
		for _, id := range md.SecretIdentifiers {
			dk := deduplicationKey(id)
			if seen[dk] {
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
	addEntries(orgMd)
	addEntries(woMd)

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
