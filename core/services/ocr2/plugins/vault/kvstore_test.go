package vault

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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

var _ (ocr3_1types.KeyValueReadWriter) = (*kv)(nil)

func TestKVStore_Secrets(t *testing.T) {
	kv := &kv{
		m: make(map[string]response),
	}
	kv.m["Key::owner::main::secret1"] = response{
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
	s, err := store.GetSecret(id)
	require.NoError(t, err)
	assert.Equal(t, s.EncryptedSecret, []byte("encrypted data"))

	delete(kv.m, "Key::owner::main::secret1")
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
	err = store.AddIDToMetadata(newKey)
	require.NoError(t, err)

	gotM, err = store.GetMetadata(owner)
	require.NoError(t, err)
	assert.Len(t, gotM.SecretIdentifiers, 3)
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

	err = store.RemoveIDFromMetadata(id)
	require.NoError(t, err)

	m, err := store.GetMetadata(owner)
	require.NoError(t, err)

	assert.Empty(t, m.SecretIdentifiers)

	err = store.RemoveIDFromMetadata(id)
	require.ErrorContains(t, err, "not found in metadata for owner owner")

	delete(kv.m, "Metadata::owner")

	err = store.RemoveIDFromMetadata(id)
	require.ErrorContains(t, err, "no metadata found for owner owner")
}
