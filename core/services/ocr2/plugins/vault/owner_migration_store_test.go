package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

const (
	testOrgID         = "org_abc123"
	testWorkflowOwner = "0xABCDEF1234567890"
)

func newTestWriteStore(t *testing.T) (*kv, WriteKVStore) {
	t.Helper()
	kvStore := &kv{m: make(map[string]response)}
	return kvStore, NewWriteStore(kvStore)
}

func writeTestSecret(t *testing.T, store WriteKVStore, owner, namespace, key string, data []byte) {
	t.Helper()
	id := &vault.SecretIdentifier{Owner: owner, Namespace: namespace, Key: key}
	require.NoError(t, store.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: data}))
}

// --- GetSecret tests ---

func TestOwnerMigrationReadStore_GetSecret_FoundUnderOrgID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "secret1", []byte("org-data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("org-data"), secret.EncryptedSecret)
}

func TestOwnerMigrationReadStore_GetSecret_FallbackToWorkflowOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("legacy-data"), secret.EncryptedSecret)
}

func TestOwnerMigrationReadStore_GetSecret_NotFoundUnderEither(t *testing.T) {
	_, inner := newTestWriteStore(t)

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "nonexistent"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationReadStore_GetSecret_PrefersOrgIDOverWorkflowOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "secret1", []byte("org-data"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("org-data"), secret.EncryptedSecret)
}

func TestOwnerMigrationReadStore_GetSecret_NoFallbackWhenWorkflowOwnerEmpty(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, "")
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationReadStore_GetSecret_NoFallbackWhenSameOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)

	store := NewOwnerMigrationReadStore(inner, testOrgID, testOrgID)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationReadStore_GetSecret_NilID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)

	_, err := store.GetSecret(nil)
	require.Error(t, err)
}

// --- GetMetadata tests ---

func TestOwnerMigrationReadStore_GetMetadata_OnlyOrgID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "secret1", []byte("data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, md.SecretIdentifiers[0].Owner)
	assert.Equal(t, "secret1", md.SecretIdentifiers[0].Key)
}

func TestOwnerMigrationReadStore_GetMetadata_OnlyWorkflowOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testWorkflowOwner, "main", "legacy1", []byte("data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, md.SecretIdentifiers[0].Owner)
	assert.Equal(t, "legacy1", md.SecretIdentifiers[0].Key)
}

func TestOwnerMigrationReadStore_GetMetadata_MergeAndDedup(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "shared_key", []byte("org-data"))
	writeTestSecret(t, inner, testOrgID, "main", "org_only", []byte("data1"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "shared_key", []byte("legacy-data"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "legacy_only", []byte("data2"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)

	// shared_key should be deduped (org_id wins), so 3 total: shared_key, org_only, legacy_only
	assert.Len(t, md.SecretIdentifiers, 3)
	for _, sid := range md.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner, "all owners should be rewritten to orgID")
	}

	keys := map[string]bool{}
	for _, sid := range md.SecretIdentifiers {
		keys[sid.Key] = true
	}
	assert.True(t, keys["shared_key"])
	assert.True(t, keys["org_only"])
	assert.True(t, keys["legacy_only"])
}

func TestOwnerMigrationReadStore_GetMetadata_BothEmpty(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)

	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	assert.Nil(t, md)
}

func TestOwnerMigrationReadStore_GetMetadata_NoMergeWhenSameOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "secret1", []byte("data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testOrgID)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
}

func TestOwnerMigrationReadStore_GetMetadata_CrossNamespaceDedup(t *testing.T) {
	_, inner := newTestWriteStore(t)
	// Same key name, different namespaces — should NOT be deduped.
	writeTestSecret(t, inner, testOrgID, "ns1", "secret1", []byte("data-ns1"))
	writeTestSecret(t, inner, testWorkflowOwner, "ns2", "secret1", []byte("data-ns2"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
}

// --- GetSecretIdentifiersCountForOwner tests ---

func TestOwnerMigrationReadStore_GetSecretIdentifiersCountForOwner_Merged(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "s1", []byte("data"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "s2", []byte("data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	count, err := store.GetSecretIdentifiersCountForOwner(testOrgID)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestOwnerMigrationReadStore_GetSecretIdentifiersCountForOwner_Deduped(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "shared", []byte("data"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "shared", []byte("data"))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	count, err := store.GetSecretIdentifiersCountForOwner(testOrgID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestOwnerMigrationReadStore_GetSecretIdentifiersCountForOwner_Empty(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)

	count, err := store.GetSecretIdentifiersCountForOwner(testOrgID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// --- GetPendingQueue pass-through tests ---

func TestOwnerMigrationReadStore_GetPendingQueue_Passthrough(t *testing.T) {
	_, inner := newTestWriteStore(t)

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{
		{Id: "req-1", Item: empty},
		{Id: "req-2", Item: empty},
	}
	require.NoError(t, inner.WritePendingQueue(items))

	store := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	got, err := store.GetPendingQueue()
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "req-1", got[0].Id)
	assert.Equal(t, "req-2", got[1].Id)
}

// --- WriteKVStore: WriteSecret tests ---

func TestOwnerMigrationWriteStore_WriteSecret_WritesUnderOrgID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "new_secret"}
	require.NoError(t, store.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("data")}))

	orgID := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "new_secret"}
	secret, err := inner.GetSecret(orgID)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data"), secret.EncryptedSecret)

	woID := &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "new_secret"}
	secret, err = inner.GetSecret(woID)
	require.NoError(t, err)
	assert.Nil(t, secret, "should not exist under workflow_owner")
}

func TestOwnerMigrationWriteStore_WriteSecret_LazyMigration(t *testing.T) {
	_, inner := newTestWriteStore(t)

	// Pre-populate a legacy entry under workflow_owner.
	writeTestSecret(t, inner, testWorkflowOwner, "main", "legacy_secret", []byte("old-data"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_secret"}
	require.NoError(t, store.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("new-data")}))

	// Verify new entry under org_id.
	orgIDSid := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_secret"}
	secret, err := inner.GetSecret(orgIDSid)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("new-data"), secret.EncryptedSecret)

	// Verify legacy entry deleted.
	woSid := &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "legacy_secret"}
	secret, err = inner.GetSecret(woSid)
	require.NoError(t, err)
	assert.Nil(t, secret, "legacy entry should be deleted after migration")

	// Verify legacy metadata cleaned up.
	woMd, err := inner.GetMetadata(testWorkflowOwner)
	require.NoError(t, err)
	if woMd != nil {
		assert.Empty(t, woMd.SecretIdentifiers)
	}
}

func TestOwnerMigrationWriteStore_WriteSecret_NoMigrationWhenNoLegacy(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "brand_new"}
	require.NoError(t, store.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("data")}))

	secret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "brand_new"})
	require.NoError(t, err)
	require.NotNil(t, secret)
}

func TestOwnerMigrationWriteStore_WriteSecret_NoMigrationWhenSameOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testOrgID)

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("data")}))
}

// --- WriteKVStore: DeleteSecret tests ---

func TestOwnerMigrationWriteStore_DeleteSecret_DeletesFromOrgID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "to_delete", []byte("data"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"}
	require.NoError(t, store.DeleteSecret(id))

	secret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationWriteStore_DeleteSecret_FallsBackToWorkflowOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testWorkflowOwner, "main", "legacy_del", []byte("data"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_del"}
	require.NoError(t, store.DeleteSecret(id))

	secret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "legacy_del"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationWriteStore_DeleteSecret_CleansBothOwners(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "both_owners", []byte("org-data"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "both_owners", []byte("legacy-data"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "both_owners"}
	require.NoError(t, store.DeleteSecret(id))

	orgSecret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "both_owners"})
	require.NoError(t, err)
	assert.Nil(t, orgSecret)

	woSecret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "both_owners"})
	require.NoError(t, err)
	assert.Nil(t, woSecret)
}

func TestOwnerMigrationWriteStore_DeleteSecret_NotFoundAnywhere(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "nonexistent"}
	err := store.DeleteSecret(id)
	require.Error(t, err, "deleting a non-existent secret should error")
}

// --- WriteKVStore: WritePendingQueue pass-through ---

func TestOwnerMigrationWriteStore_WritePendingQueue_Passthrough(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{
		{Id: "pq-1", Item: empty},
	}
	require.NoError(t, store.WritePendingQueue(items))

	got, err := inner.GetPendingQueue()
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "pq-1", got[0].Id)
}

// --- WriteKVStore: WriteMetadata ---

func TestOwnerMigrationWriteStore_WriteMetadata_WritesUnderOrgID(t *testing.T) {
	_, inner := newTestWriteStore(t)
	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	md := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	require.NoError(t, store.WriteMetadata("anything", md))

	got, err := inner.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Len(t, got.SecretIdentifiers, 1)
}

// --- End-to-end migration scenarios ---

func TestOwnerMigrationStore_CreateOldFlow_ReadNewFlow(t *testing.T) {
	_, inner := newTestWriteStore(t)

	// Simulate old flow: secret created under workflow_owner.
	writeTestSecret(t, inner, testWorkflowOwner, "main", "migrating_secret", []byte("old-data"))

	// New flow: read through migration adapter.
	readStore := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating_secret"}
	secret, err := readStore.GetSecret(id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("old-data"), secret.EncryptedSecret)
}

func TestOwnerMigrationStore_CreateOldFlow_UpdateNewFlow_LazyMigration(t *testing.T) {
	_, inner := newTestWriteStore(t)

	// Old flow: write under workflow_owner.
	writeTestSecret(t, inner, testWorkflowOwner, "main", "migrating", []byte("old"))

	writeStore := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	// New flow: update (which triggers lazy migration).
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating"}
	require.NoError(t, writeStore.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("new")}))

	// Verify: exists under org_id, not under workflow_owner.
	orgSecret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating"})
	require.NoError(t, err)
	require.NotNil(t, orgSecret)
	assert.Equal(t, []byte("new"), orgSecret.EncryptedSecret)

	woSecret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "migrating"})
	require.NoError(t, err)
	assert.Nil(t, woSecret)
}

func TestOwnerMigrationStore_CreateOldFlow_ListNewFlow(t *testing.T) {
	_, inner := newTestWriteStore(t)

	writeTestSecret(t, inner, testWorkflowOwner, "main", "legacy1", []byte("d1"))
	writeTestSecret(t, inner, testWorkflowOwner, "alt", "legacy2", []byte("d2"))

	readStore := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := readStore.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
	for _, sid := range md.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner)
	}
}

func TestOwnerMigrationStore_CreateOldFlow_DeleteNewFlow(t *testing.T) {
	_, inner := newTestWriteStore(t)

	writeTestSecret(t, inner, testWorkflowOwner, "main", "to_delete", []byte("data"))

	writeStore := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"}
	require.NoError(t, writeStore.DeleteSecret(id))

	// Verify deleted from workflow_owner.
	secret, err := inner.GetSecret(&vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "to_delete"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestOwnerMigrationStore_UpdateMigration_ThenListShowsNoDuplicates(t *testing.T) {
	_, inner := newTestWriteStore(t)

	// Create two secrets under legacy owner.
	writeTestSecret(t, inner, testWorkflowOwner, "main", "s1", []byte("old1"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "s2", []byte("old2"))

	writeStore := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)

	// Update (migrate) s1.
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, writeStore.WriteSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("new1")}))

	// List should show 2 unique secrets, all with org_id owner.
	readStore := NewOwnerMigrationReadStore(inner, testOrgID, testWorkflowOwner)
	md, err := readStore.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
	for _, sid := range md.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner)
	}
}

// --- mergeMetadata unit tests ---

func TestMergeMetadata_BothNil(t *testing.T) {
	result := mergeMetadata(nil, nil, testOrgID)
	assert.Nil(t, result)
}

func TestMergeMetadata_OrgOnly(t *testing.T) {
	orgMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	result := mergeMetadata(orgMd, nil, testOrgID)
	require.NotNil(t, result)
	assert.Len(t, result.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, result.SecretIdentifiers[0].Owner)
}

func TestMergeMetadata_WorkflowOwnerOnly(t *testing.T) {
	woMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testWorkflowOwner, Namespace: "main", Key: "s1"},
		},
	}
	result := mergeMetadata(nil, woMd, testOrgID)
	require.NotNil(t, result)
	assert.Len(t, result.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, result.SecretIdentifiers[0].Owner, "owner should be rewritten to orgID")
}

func TestMergeMetadata_Deduplication(t *testing.T) {
	orgMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "shared"},
			{Owner: testOrgID, Namespace: "main", Key: "org_only"},
		},
	}
	woMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testWorkflowOwner, Namespace: "main", Key: "shared"},
			{Owner: testWorkflowOwner, Namespace: "main", Key: "wo_only"},
		},
	}
	result := mergeMetadata(orgMd, woMd, testOrgID)
	require.NotNil(t, result)
	assert.Len(t, result.SecretIdentifiers, 3)

	keys := map[string]bool{}
	for _, sid := range result.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner)
		keys[sid.Key] = true
	}
	assert.True(t, keys["shared"])
	assert.True(t, keys["org_only"])
	assert.True(t, keys["wo_only"])
}

func TestMergeMetadata_DefaultNamespaceNormalization(t *testing.T) {
	orgMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	// Empty namespace normalizes to "main" for dedup purposes.
	woMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testWorkflowOwner, Namespace: "", Key: "s1"},
		},
	}
	result := mergeMetadata(orgMd, woMd, testOrgID)
	require.NotNil(t, result)
	assert.Len(t, result.SecretIdentifiers, 1, "empty namespace should dedup against 'main'")
}

// --- deduplicationKey tests ---

func TestDeduplicationKey(t *testing.T) {
	tests := []struct {
		name      string
		id        *vault.SecretIdentifier
		expected  string
	}{
		{
			name:     "normal",
			id:       &vault.SecretIdentifier{Owner: "any", Namespace: "ns1", Key: "k1"},
			expected: "ns1::k1",
		},
		{
			name:     "empty namespace defaults to main",
			id:       &vault.SecretIdentifier{Owner: "any", Namespace: "", Key: "k1"},
			expected: "main::k1",
		},
		{
			name:     "main namespace explicit",
			id:       &vault.SecretIdentifier{Owner: "any", Namespace: "main", Key: "k1"},
			expected: "main::k1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, deduplicationKey(tt.id))
		})
	}
}

// --- needsMigration tests ---

func TestNeedsMigration(t *testing.T) {
	assert.True(t, needsMigration("org1", "wo1"))
	assert.False(t, needsMigration("org1", "org1"))
	assert.False(t, needsMigration("org1", ""))
}

// --- WriteKVStore GetSecret (exercising the write store's read path) ---

func TestOwnerMigrationWriteStore_GetSecret_FallbackToWorkflowOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("legacy"), secret.EncryptedSecret)
}

func TestOwnerMigrationWriteStore_GetMetadata_Merged(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "s1", []byte("d1"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "s2", []byte("d2"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	md, err := store.GetMetadata(testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
}

func TestOwnerMigrationWriteStore_GetSecretIdentifiersCountForOwner(t *testing.T) {
	_, inner := newTestWriteStore(t)
	writeTestSecret(t, inner, testOrgID, "main", "s1", []byte("d1"))
	writeTestSecret(t, inner, testWorkflowOwner, "main", "s2", []byte("d2"))

	store := NewOwnerMigrationWriteStore(inner, testOrgID, testWorkflowOwner)
	count, err := store.GetSecretIdentifiersCountForOwner(testOrgID)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

// --- Interface compliance (compile-time checks happen at top of production file) ---

func TestOwnerMigrationReadStore_ImplementsReadKVStore(t *testing.T) {
	var _ ReadKVStore = (*OwnerMigrationReadStore)(nil)
}

func TestOwnerMigrationWriteStore_ImplementsWriteKVStore(t *testing.T) {
	var _ WriteKVStore = (*OwnerMigrationWriteStore)(nil)
}

// --- Error propagation tests ---

func TestOwnerMigrationReadStore_GetSecret_PropagatesInnerError(t *testing.T) {
	inner := &kv{m: map[string]response{}}
	inner.m["Metadata::"+testOrgID] = response{err: assert.AnError}
	store := NewOwnerMigrationReadStore(NewReadStore(inner), testOrgID, testWorkflowOwner)

	_, err := store.GetSecret(&vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"})
	require.Error(t, err)
}

func TestOwnerMigrationReadStore_GetMetadata_PropagatesOrgIDError(t *testing.T) {
	inner := &kv{m: map[string]response{}}
	inner.m["Metadata::"+testOrgID] = response{err: assert.AnError}
	store := NewOwnerMigrationReadStore(NewReadStore(inner), testOrgID, testWorkflowOwner)

	_, err := store.GetMetadata(testOrgID)
	require.Error(t, err)
}

func TestOwnerMigrationReadStore_GetMetadata_PropagatesWorkflowOwnerError(t *testing.T) {
	inner := &kv{m: map[string]response{}}
	orgMdBytes, _ := proto.Marshal(&vault.StoredMetadata{SecretIdentifiers: []*vault.SecretIdentifier{}})
	inner.m["Metadata::"+testOrgID] = response{data: orgMdBytes}
	inner.m["Metadata::"+testWorkflowOwner] = response{err: assert.AnError}
	store := NewOwnerMigrationReadStore(NewReadStore(inner), testOrgID, testWorkflowOwner)

	_, err := store.GetMetadata(testOrgID)
	require.Error(t, err)
}
