package vault

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	testOrgID         = "org_abc123"
	testWorkflowOwner = "0xABCDEF1234567890"
)

func newMigrationTestStore(t *testing.T) (*kv, WriteKVStore) {
	t.Helper()
	kvStore := &kv{m: make(map[string]response)}
	return kvStore, NewWriteStore(kvStore, newTestMetrics(t))
}

func writeTestSecret(ctx context.Context, t *testing.T, store WriteKVStore, owner, namespace, key string, data []byte) {
	t.Helper()
	id := &vault.SecretIdentifier{Owner: owner, Namespace: namespace, Key: key}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: data}))
}

// --- GetSecret tests ---

func TestKVStoreWrapper_GetSecret_FoundUnderOrgID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "secret1", []byte("org-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("org-data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_GetSecret_FallbackToWorkflowOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("legacy-data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_GetSecret_NotFoundUnderEither(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "nonexistent"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_GetSecret_PrefersOrgIDOverWorkflowOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "secret1", []byte("org-data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("org-data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_GetSecret_NoFallbackWhenWorkflowOwnerEmpty(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "secret1", []byte("legacy-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, "")
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_GetSecret_NoFallbackWhenSameOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "secret1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testOrgID)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_GetSecret_NilID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	_, err := store.GetSecret(ctx, nil, testOrgID, testWorkflowOwner)
	require.Error(t, err)
}

func TestKVStoreWrapper_GetSecret_DifferentOwnersPerCall(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, "org_A", "main", "s1", []byte("data-A"))
	writeTestSecret(ctx, t, inner, "wo_B", "main", "s2", []byte("data-B"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id1 := &vault.SecretIdentifier{Owner: "org_A", Namespace: "main", Key: "s1"}
	secret, err := store.GetSecret(ctx, id1, "org_A", "wo_A")
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data-A"), secret.EncryptedSecret)

	id2 := &vault.SecretIdentifier{Owner: "org_B", Namespace: "main", Key: "s2"}
	secret, err = store.GetSecret(ctx, id2, "org_B", "wo_B")
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data-B"), secret.EncryptedSecret)
}

// --- GetMetadata tests ---

func TestKVStoreWrapper_GetMetadata_OnlyOrgID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "secret1", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, md.SecretIdentifiers[0].Owner)
	assert.Equal(t, "secret1", md.SecretIdentifiers[0].Key)
}

func TestKVStoreWrapper_GetMetadata_OnlyWorkflowOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy1", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, md.SecretIdentifiers[0].Owner)
	assert.Equal(t, "legacy1", md.SecretIdentifiers[0].Key)
}

func TestKVStoreWrapper_GetMetadata_MergeAndDedup(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "shared_key", []byte("org-data"))
	writeTestSecret(ctx, t, inner, testOrgID, "main", "org_only", []byte("data1"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "shared_key", []byte("legacy-data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy_only", []byte("data2"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)

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

func TestKVStoreWrapper_GetMetadata_BothEmpty(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Nil(t, md)
}

func TestKVStoreWrapper_GetMetadata_NoMergeWhenSameOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "secret1", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testOrgID)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
}

func TestKVStoreWrapper_GetMetadata_CrossNamespaceDedup(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "ns1", "secret1", []byte("data-ns1"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "ns2", "secret1", []byte("data-ns2"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2, "same key in different namespaces should NOT be deduped")
}

// --- GetSecretIdentifiersCountForOwner tests ---

func TestKVStoreWrapper_GetSecretIdentifiersCountForOwner_Merged(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s2", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	count, err := store.GetSecretIdentifiersCountForOwner(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestKVStoreWrapper_GetSecretIdentifiersCountForOwner_Deduped(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "shared", []byte("data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "shared", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	count, err := store.GetSecretIdentifiersCountForOwner(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestKVStoreWrapper_GetSecretIdentifiersCountForOwner_Empty(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	count, err := store.GetSecretIdentifiersCountForOwner(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// --- GetPendingQueue / WritePendingQueue pass-through tests ---

func TestKVStoreWrapper_GetPendingQueue_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{
		{Id: "req-1", Item: empty},
		{Id: "req-2", Item: empty},
	}
	require.NoError(t, inner.WritePendingQueue(ctx, items))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	got, err := store.GetPendingQueue(ctx)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "req-1", got[0].Id)
	assert.Equal(t, "req-2", got[1].Id)
}

func TestKVStoreWrapper_WritePendingQueue_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{
		{Id: "pq-1", Item: empty},
	}
	require.NoError(t, store.WritePendingQueue(ctx, items))

	got, err := inner.GetPendingQueue(ctx)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "pq-1", got[0].Id)
}

// --- WriteSecret tests ---

func TestKVStoreWrapper_WriteSecret_WritesUnderOrgID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "new_secret"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("data")}, testOrgID, testWorkflowOwner))

	orgID := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "new_secret"}
	secret, err := inner.GetSecret(ctx, orgID)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data"), secret.EncryptedSecret)

	woID := &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "new_secret"}
	secret, err = inner.GetSecret(ctx, woID)
	require.NoError(t, err)
	assert.Nil(t, secret, "should not exist under workflow_owner")
}

func TestKVStoreWrapper_WriteSecret_LazyMigration(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy_secret", []byte("old-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_secret"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("new-data")}, testOrgID, testWorkflowOwner))

	orgIDSid := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_secret"}
	secret, err := inner.GetSecret(ctx, orgIDSid)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("new-data"), secret.EncryptedSecret)

	woSid := &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "legacy_secret"}
	secret, err = inner.GetSecret(ctx, woSid)
	require.NoError(t, err)
	assert.Nil(t, secret, "legacy entry should be deleted after migration")

	woMd, err := inner.GetMetadata(ctx, testWorkflowOwner)
	require.NoError(t, err)
	if woMd != nil {
		assert.Empty(t, woMd.SecretIdentifiers)
	}
}

func TestKVStoreWrapper_WriteSecret_NoMigrationWhenNoLegacy(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "brand_new"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("data")}, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "brand_new"})
	require.NoError(t, err)
	require.NotNil(t, secret)
}

func TestKVStoreWrapper_WriteSecret_NoMigrationWhenSameOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("data")}, testOrgID, testOrgID))
}

// --- WriteMetadata test ---

func TestKVStoreWrapper_WriteMetadata_WritesUnderOrgID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	md := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	require.NoError(t, store.WriteMetadata(ctx, testOrgID, md))

	got, err := inner.GetMetadata(ctx, testOrgID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Len(t, got.SecretIdentifiers, 1)
}

// --- DeleteSecret tests ---

func TestKVStoreWrapper_DeleteSecret_DeletesFromOrgID(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "to_delete", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_DeleteSecret_FallsBackToWorkflowOwner(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy_del", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "legacy_del"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "legacy_del"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_DeleteSecret_CleansBothOwners(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "both_owners", []byte("org-data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "both_owners", []byte("legacy-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "both_owners"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	orgSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "both_owners"})
	require.NoError(t, err)
	assert.Nil(t, orgSecret)

	woSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "both_owners"})
	require.NoError(t, err)
	assert.Nil(t, woSecret)
}

func TestKVStoreWrapper_DeleteSecret_NotFoundAnywhere(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "nonexistent"}
	err := store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.Error(t, err, "deleting a non-existent secret should error")
}

// --- End-to-end migration scenarios ---

func TestKVStoreWrapper_CreateOldFlow_ReadNewFlow(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "migrating_secret", []byte("old-data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating_secret"}
	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("old-data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_CreateOldFlow_UpdateNewFlow_LazyMigration(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "migrating", []byte("old"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("new")}, testOrgID, testWorkflowOwner))

	orgSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "migrating"})
	require.NoError(t, err)
	require.NotNil(t, orgSecret)
	assert.Equal(t, []byte("new"), orgSecret.EncryptedSecret)

	woSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "migrating"})
	require.NoError(t, err)
	assert.Nil(t, woSecret)
}

func TestKVStoreWrapper_CreateOldFlow_ListNewFlow(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy1", []byte("d1"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "alt", "legacy2", []byte("d2"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
	for _, sid := range md.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner)
	}
}

func TestKVStoreWrapper_CreateOldFlow_DeleteNewFlow(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "to_delete", []byte("data"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "to_delete"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "to_delete"})
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_UpdateMigration_ThenListShowsNoDuplicates(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s1", []byte("old1"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s2", []byte("old2"))

	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("new1")}, testOrgID, testWorkflowOwner))

	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 2)
	for _, sid := range md.SecretIdentifiers {
		assert.Equal(t, testOrgID, sid.Owner)
	}
}

// --- mergeMetadata unit tests ---

func TestMergeMetadata_BothNil(t *testing.T) {
	result := mergeMetadata(nil, nil, testOrgID, logger.TestLogger(t))
	assert.Nil(t, result)
}

func TestMergeMetadata_OrgOnly(t *testing.T) {
	orgMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	result := mergeMetadata(orgMd, nil, testOrgID, logger.TestLogger(t))
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
	result := mergeMetadata(nil, woMd, testOrgID, logger.TestLogger(t))
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
	result := mergeMetadata(orgMd, woMd, testOrgID, logger.TestLogger(t))
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
	woMd := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testWorkflowOwner, Namespace: "", Key: "s1"},
		},
	}
	result := mergeMetadata(orgMd, woMd, testOrgID, logger.TestLogger(t))
	require.NotNil(t, result)
	assert.Len(t, result.SecretIdentifiers, 1, "empty namespace should dedup against 'main'")
}

// --- deduplicationKey tests ---

func TestDeduplicationKey(t *testing.T) {
	tests := []struct {
		name     string
		id       *vault.SecretIdentifier
		expected string
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

// --- Error propagation tests ---

func TestKVStoreWrapper_GetSecret_PropagatesInnerError(t *testing.T) {
	ctx := t.Context()
	inner := &kv{m: map[string]response{}}
	inner.m["Metadata::"+testOrgID] = response{err: assert.AnError}
	store := NewKVStoreWrapper(NewWriteStore(inner, newTestMetrics(t)), true, logger.TestLogger(t))

	_, err := store.GetSecret(ctx, &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}, testOrgID, testWorkflowOwner)
	require.Error(t, err)
}

func TestKVStoreWrapper_GetMetadata_PropagatesOrgIDError(t *testing.T) {
	ctx := t.Context()
	inner := &kv{m: map[string]response{}}
	inner.m["Metadata::"+testOrgID] = response{err: assert.AnError}
	store := NewKVStoreWrapper(NewWriteStore(inner, newTestMetrics(t)), true, logger.TestLogger(t))

	_, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.Error(t, err)
}

func TestKVStoreWrapper_GetMetadata_PropagatesWorkflowOwnerError(t *testing.T) {
	ctx := t.Context()
	inner := &kv{m: map[string]response{}}
	orgMdBytes, _ := proto.Marshal(&vault.StoredMetadata{SecretIdentifiers: []*vault.SecretIdentifier{}})
	inner.m["Metadata::"+testOrgID] = response{data: orgMdBytes}
	inner.m["Metadata::"+testWorkflowOwner] = response{err: assert.AnError}
	store := NewKVStoreWrapper(NewWriteStore(inner, newTestMetrics(t)), true, logger.TestLogger(t))

	_, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.Error(t, err)
}

// --- Migration disabled (passthrough) tests ---

func TestKVStoreWrapper_Disabled_GetSecret_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_Disabled_GetSecret_NoFallback(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s1", []byte("legacy"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}

	secret, err := store.GetSecret(ctx, id, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Nil(t, secret, "should NOT fall back to workflow_owner when migration is disabled")
}

func TestKVStoreWrapper_Disabled_GetMetadata_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1)
	assert.Equal(t, testOrgID, md.SecretIdentifiers[0].Owner)
}

func TestKVStoreWrapper_Disabled_GetMetadata_NoMerge(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "org_secret", []byte("data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "legacy_secret", []byte("data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	md, err := store.GetMetadata(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	require.NotNil(t, md)
	assert.Len(t, md.SecretIdentifiers, 1, "should only return orgID metadata, not merge with workflow_owner")
	assert.Equal(t, "org_secret", md.SecretIdentifiers[0].Key)
}

func TestKVStoreWrapper_Disabled_GetSecretIdentifiersCountForOwner_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s2", []byte("data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	count, err := store.GetSecretIdentifiersCountForOwner(ctx, testOrgID, testWorkflowOwner)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "should only count orgID secrets, not merged")
}

func TestKVStoreWrapper_Disabled_WriteSecret_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))

	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("data")}, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, secret)
	assert.Equal(t, []byte("data"), secret.EncryptedSecret)
}

func TestKVStoreWrapper_Disabled_WriteSecret_NoLazyMigration(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s1", []byte("legacy"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.WriteSecret(ctx, id, &vault.StoredSecret{EncryptedSecret: []byte("new")}, testOrgID, testWorkflowOwner))

	woSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "s1"})
	require.NoError(t, err)
	assert.NotNil(t, woSecret, "legacy entry should NOT be cleaned up when migration is disabled")
}

func TestKVStoreWrapper_Disabled_WriteMetadata_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))

	md := &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{Owner: testOrgID, Namespace: "main", Key: "s1"},
		},
	}
	require.NoError(t, store.WriteMetadata(ctx, testOrgID, md))

	got, err := inner.GetMetadata(ctx, testOrgID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Len(t, got.SecretIdentifiers, 1)
}

func TestKVStoreWrapper_Disabled_DeleteSecret_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	secret, err := inner.GetSecret(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestKVStoreWrapper_Disabled_DeleteSecret_NoDualOwnerCleanup(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	writeTestSecret(ctx, t, inner, testOrgID, "main", "s1", []byte("org-data"))
	writeTestSecret(ctx, t, inner, testWorkflowOwner, "main", "s1", []byte("legacy-data"))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	id := &vault.SecretIdentifier{Owner: testOrgID, Namespace: "main", Key: "s1"}
	require.NoError(t, store.DeleteSecret(ctx, id, testOrgID, testWorkflowOwner))

	woSecret, err := inner.GetSecret(ctx, &vault.SecretIdentifier{Owner: testWorkflowOwner, Namespace: "main", Key: "s1"})
	require.NoError(t, err)
	assert.NotNil(t, woSecret, "legacy entry should NOT be cleaned up when migration is disabled")
}

func TestKVStoreWrapper_Disabled_GetPendingQueue_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{{Id: "req-1", Item: empty}}
	require.NoError(t, inner.WritePendingQueue(ctx, items))

	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	got, err := store.GetPendingQueue(ctx)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "req-1", got[0].Id)
}

func TestKVStoreWrapper_Disabled_WritePendingQueue_Passthrough(t *testing.T) {
	ctx := t.Context()
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))

	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	items := []*vault.StoredPendingQueueItem{{Id: "pq-1", Item: empty}}
	require.NoError(t, store.WritePendingQueue(ctx, items))

	got, err := inner.GetPendingQueue(ctx)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "pq-1", got[0].Id)
}

func TestKVStoreWrapper_Disabled_AdapterIsNil(t *testing.T) {
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, false, logger.TestLogger(t))
	assert.Nil(t, store.adapter, "adapter should be nil when migration is disabled")
}

func TestKVStoreWrapper_Enabled_AdapterIsSet(t *testing.T) {
	_, inner := newMigrationTestStore(t)
	store := NewKVStoreWrapper(inner, true, logger.TestLogger(t))
	assert.NotNil(t, store.adapter, "adapter should be set when migration is enabled")
}
