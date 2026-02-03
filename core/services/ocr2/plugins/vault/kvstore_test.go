package vault

import (
	"context"
	"errors"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

type response struct {
	data []byte
	err  error
}

type kv struct {
	m map[string]response
}

func (k *kv) Read(key []byte) ([]byte, error) {
	d := k.m[string(key)]
	return d.data, d.err
}

func (k *kv) Delete(key []byte) error {
	delete(k.m, string(key))
	return nil
}

func (k *kv) Write(key []byte, data []byte) error {
	k.m[string(key)] = response{
		data: data,
	}
	return nil
}

type blobber struct {
	blobs [][]byte
	cnt   int
}

func (b *blobber) BroadcastBlob(_ context.Context, data []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	b.blobs = append(b.blobs, data)
	return ocr3_1types.BlobHandle{}, nil
}

func (b *blobber) FetchBlob(_ context.Context, _ ocr3_1types.BlobHandle) ([]byte, error) {
	blob := b.blobs[b.cnt]
	b.cnt++
	return blob, nil
}

var _ (ocr3_1types.BlobBroadcastFetcher) = (*blobber)(nil)

var _ (ocr3_1types.KeyValueReadWriter) = (*kv)(nil)

func TestKVStore_Secrets(t *testing.T) {
	kv := &kv{
		m: make(map[string]response),
	}
	kv.m["Metadata::owner"] = response{
		err: errors.New("not found"),
	}
	store := NewWriteStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}

	_, err := store.GetSecret(id)
	require.ErrorContains(t, err, "not found")

	d, err := proto.Marshal(&vault.StoredSecret{
		EncryptedSecret: []byte("encrypted data"),
	})
	require.NoError(t, err)
	kv.m["Key::owner::main::secret1"] = response{
		data: d,
	}
	d, err = proto.Marshal(&vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{id},
	})
	require.NoError(t, err)
	kv.m["Metadata::owner"] = response{
		data: d,
	}
	s, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Equal(t, s.EncryptedSecret, []byte("encrypted data"))

	delete(kv.m, "Metadata::owner")
	s, err = store.GetSecret(id)
	assert.Nil(t, s)
	require.NoError(t, err)

	newData := []byte("new encrypted data 2")
	ss := &vault.StoredSecret{
		EncryptedSecret: newData,
	}
	err = store.WriteSecret(id, ss)
	require.NoError(t, err)

	s, err = store.GetSecret(id)
	require.NoError(t, err)
	assert.Equal(t, newData, s.EncryptedSecret)
}

func TestKVStore_DeleteSecrets(t *testing.T) {
	kv := &kv{
		m: make(map[string]response),
	}
	store := NewWriteStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}
	err := store.WriteSecret(id, &vault.StoredSecret{
		EncryptedSecret: []byte("encrypted data"),
	})
	require.NoError(t, err)

	err = store.DeleteSecret(id)
	require.NoError(t, err)

	md, err := store.GetMetadata("owner")
	require.NoError(t, err)

	assert.Empty(t, md.SecretIdentifiers)
}

func TestKVStore_Metadata(t *testing.T) {
	owner := "owner"
	kv := &kv{
		m: make(map[string]response),
	}
	kv.m["Metadata::"+owner] = response{
		err: errors.New("not found"),
	}
	store := NewWriteStore(kv)

	_, err := store.GetMetadata(owner)
	require.ErrorContains(t, err, "not found")

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}
	d, err := proto.Marshal(&vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{id},
	})
	require.NoError(t, err)
	kv.m["Metadata::owner"] = response{
		data: d,
	}
	m, err := store.GetMetadata(owner)
	require.NoError(t, err)
	assert.Len(t, m.SecretIdentifiers, 1)
	assert.True(t, proto.Equal(m.SecretIdentifiers[0], id))

	delete(kv.m, "Metadata::"+owner)
	m, err = store.GetMetadata(owner)
	assert.Nil(t, m)
	require.NoError(t, err)

	m = &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{
				Owner:     "owner",
				Namespace: "main",
				Key:       "secret2",
			},
			{
				Owner:     "owner",
				Namespace: "main",
				Key:       "secret3",
			},
		},
	}
	err = store.WriteMetadata(owner, m)
	require.NoError(t, err)

	gotM, err := store.GetMetadata(owner)
	require.NoError(t, err)
	assert.True(t, proto.Equal(m, gotM))

	newKey := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret3",
	}
	err = store.addIDToMetadata(newKey)
	require.NoError(t, err)

	gotM, err = store.GetMetadata(owner)
	require.NoError(t, err)
	assert.Len(t, gotM.SecretIdentifiers, 2)
}

func TestKVStore_Metadata_Delete(t *testing.T) {
	owner := "owner"
	kv := &kv{
		m: make(map[string]response),
	}
	store := NewWriteStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}
	d, err := proto.Marshal(&vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{id},
	})
	require.NoError(t, err)
	kv.m["Metadata::owner"] = response{
		data: d,
	}

	err = store.removeIDFromMetadata(id)
	require.NoError(t, err)

	m, err := store.GetMetadata(owner)
	require.NoError(t, err)

	assert.Empty(t, m.SecretIdentifiers)

	err = store.removeIDFromMetadata(id)
	require.ErrorContains(t, err, "not found in metadata for owner owner")

	delete(kv.m, "Metadata::owner")

	err = store.removeIDFromMetadata(id)
	require.ErrorContains(t, err, "no metadata found for owner owner")
}

func TestKVStore_InconsistentWrites(t *testing.T) {
	kv := &kv{
		m: make(map[string]response),
	}
	store := NewWriteStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}

	d, err := proto.Marshal(&vault.StoredSecret{
		EncryptedSecret: []byte("encrypted data"),
	})
	require.NoError(t, err)
	kv.m["Key::owner::main::secret1"] = response{
		data: d,
	}
	d, err = proto.Marshal(&vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{id},
	})
	require.NoError(t, err)
	kv.m["Metadata::owner"] = response{
		data: d,
	}
	s, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Equal(t, s.EncryptedSecret, []byte("encrypted data"))

	// Simulate a delete which was inconsistent;
	// we deleted the metadata record but not the secret iteslf.
	delete(kv.m, "Metadata::owner")

	// Now fetching the secret should fail
	s, err = store.GetSecret(id)
	assert.Nil(t, s)
	require.NoError(t, err)

	// We can recreate it without an already exists error.
	err = store.WriteSecret(id, &vault.StoredSecret{
		EncryptedSecret: []byte("encrypted data 2"),
	})
	require.NoError(t, err)

	md, err := store.GetMetadata("owner")
	require.NoError(t, err)
	assert.Len(t, md.SecretIdentifiers, 1)

	s, err = store.GetSecret(id)
	assert.NotNil(t, s)
	require.NoError(t, err)

	assert.Equal(t, []byte("encrypted data 2"), s.EncryptedSecret)
}

func TestKVStore_GetPendingRequests(t *testing.T) {
	// Simulating an in-memory kv store.
	kv := &kv{
		m: make(map[string]response),
	}
	store := NewWriteStore(kv)

	// Expect no pending requests on empty store.
	requests, err := store.GetPendingQueue()
	require.NoError(t, err)
	assert.Empty(t, requests)

	// Add mock pending requests.
	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	item := &vault.StoredPendingQueueItem{
		Id:   "test-request-id-123",
		Item: empty,
	}
	d, err := proto.Marshal(item)
	require.NoError(t, err)
	kv.m[pendingQueueItemPrefix+"0"] = response{data: d}

	item2 := &vault.StoredPendingQueueItem{
		Id:   "test-request-id-456",
		Item: empty,
	}
	d, err = proto.Marshal(item2)
	require.NoError(t, err)
	kv.m[pendingQueueItemPrefix+"1"] = response{data: d}

	index := &vault.StoredPendingQueueIndex{Length: 2}
	indexBytes, err := proto.Marshal(index)
	require.NoError(t, err)
	kv.m[pendingQueueIndex] = response{data: indexBytes}

	// Validate retrieval of one pending request.
	requests, err = store.GetPendingQueue()
	require.NoError(t, err)
	assert.Len(t, requests, 2)
	assert.Equal(t, "test-request-id-123", requests[0].Id)
	assert.Equal(t, "test-request-id-456", requests[1].Id)

	// Validate behaviour when the index item is missing.
	delete(kv.m, pendingQueueIndex)
	requests, err = store.GetPendingQueue()
	require.NoError(t, err)
	assert.Empty(t, requests)

	// Validate behaviour when one of the queue items is missing
	index = &vault.StoredPendingQueueIndex{Length: 3}
	indexBytes, err = proto.Marshal(index)
	require.NoError(t, err)
	kv.m[pendingQueueIndex] = response{data: indexBytes}

	requests, err = store.GetPendingQueue()
	require.ErrorContains(t, err, "pending queue item at index 2 not found")
	assert.Empty(t, requests)
}

func TestKVStore_WritePendingRequests(t *testing.T) {
	// Simulating an in-memory kv store.
	kv := &kv{
		m: make(map[string]response),
	}
	store := NewWriteStore(kv)

	// Writing mock pending requests.
	empty, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	item := &vault.StoredPendingQueueItem{
		Id:   "test-request-id-1",
		Item: empty,
	}
	item2 := &vault.StoredPendingQueueItem{
		Id:   "test-request-id-2",
		Item: empty,
	}
	item3 := &vault.StoredPendingQueueItem{
		Id:   "test-request-id-3",
		Item: empty,
	}
	err = store.WritePendingQueue([]*vault.StoredPendingQueueItem{item, item2, item3})
	require.NoError(t, err)

	// Ensure index is correctly written.
	indexBytes, exists := kv.m[pendingQueueIndex]
	assert.True(t, exists)
	index := &vault.StoredPendingQueueIndex{
		Length: 3,
	}
	require.NoError(t, proto.Unmarshal(indexBytes.data, index))
	assert.Equal(t, int64(3), index.Length)

	// Ensure queue items are correctly written.
	itemBytes, exists := kv.m[pendingQueueItemPrefix+"0"]
	assert.True(t, exists)
	item = &vault.StoredPendingQueueItem{}
	require.NoError(t, proto.Unmarshal(itemBytes.data, item))
	assert.Equal(t, "test-request-id-1", item.Id)

	itemBytes, exists = kv.m[pendingQueueItemPrefix+"1"]
	assert.True(t, exists)
	item2 = &vault.StoredPendingQueueItem{}
	require.NoError(t, proto.Unmarshal(itemBytes.data, item2))
	assert.Equal(t, "test-request-id-2", item2.Id)

	itemBytes, exists = kv.m[pendingQueueItemPrefix+"2"]
	assert.True(t, exists)
	item2 = &vault.StoredPendingQueueItem{}
	require.NoError(t, proto.Unmarshal(itemBytes.data, item2))
	assert.Equal(t, "test-request-id-3", item2.Id)

	// Writing a shorter list deletes the old one.
	err = store.WritePendingQueue([]*vault.StoredPendingQueueItem{item, item2})
	require.NoError(t, err)

	_, exists = kv.m[pendingQueueItemPrefix+"3"]
	assert.False(t, exists)
}
