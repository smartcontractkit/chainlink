package vault

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	store := newWriteStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}

	_, err := store.getSecret(id)
	require.ErrorContains(t, err, "not found")

	d, err := proto.Marshal(&vault.StoredSecret{
		EncryptedSecret: []byte("encrypted data"),
	})
	require.NoError(t, err)
	kv.m["Key::owner::main::secret1"] = response{
		data: d,
	}
	s, err := store.getSecret(id)
	require.NoError(t, err)
	assert.Equal(t, s.EncryptedSecret, []byte("encrypted data"))

	delete(kv.m, "Key::owner::main::secret1")
	s, err = store.getSecret(id)
	assert.Nil(t, s)
	assert.NoError(t, err)

	newData := []byte("new encrypted data 2")
	ss := &vault.StoredSecret{
		EncryptedSecret: newData,
	}
	err = store.writeSecret(id, ss)
	assert.NoError(t, err)

	s, err = store.getSecret(id)
	assert.NoError(t, err)
	assert.Equal(t, newData, s.EncryptedSecret)
}

func TestKVStore_Metadata(t *testing.T) {
	owner := "owner"
	kv := &kv{
		m: make(map[string]response),
	}
	kv.m["Metadata::"+owner] = response{
		err: errors.New("not found"),
	}
	store := newWriteStore(kv)

	_, err := store.getMetadata(owner)
	require.ErrorContains(t, err, "not found")

	key := "Key::owner::main::secret1"
	d, err := proto.Marshal(&vault.StoredMetadata{
		Keys: []string{key},
	})
	require.NoError(t, err)
	kv.m["Metadata::owner"] = response{
		data: d,
	}
	m, err := store.getMetadata(owner)
	require.NoError(t, err)
	assert.Len(t, m.Keys, 1)
	assert.Equal(t, m.Keys[0], key)

	delete(kv.m, "Metadata::"+owner)
	m, err = store.getMetadata(owner)
	assert.Nil(t, m)
	assert.NoError(t, err)

	m = &vault.StoredMetadata{
		Keys: []string{"Key::owner::main::secret1", "Key::owner::main::secret2"},
	}
	err = store.writeMetadata(owner, m)
	assert.NoError(t, err)

	gotM, err := store.getMetadata(owner)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(m, gotM))

	newKey := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret3",
	}
	err = store.addKeyToMetadata(newKey)
	assert.NoError(t, err)

	gotM, err = store.getMetadata(owner)
	assert.NoError(t, err)
	assert.Len(t, gotM.Keys, 3)
}
