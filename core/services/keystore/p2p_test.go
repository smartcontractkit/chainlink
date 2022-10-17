package keystore_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_P2PKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	ks := keyStore.P2P()
	reset := func() {
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(cltest.Password))
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("errors when getting non-existent ID", func(t *testing.T) {
		defer reset()
		_, err := ks.Get("non-existent-id")
		require.Error(t, err)
	})

	t.Run("creates a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.PeerID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.PeerID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Export("non-existent", cltest.Password)
		assert.Error(t, err)
		_, err = ks.Delete(key.PeerID())
		require.NoError(t, err)
		_, err = ks.Get(key.PeerID())
		require.Error(t, err)
		importedKey, err := ks.Import(exportJSON, cltest.Password)
		require.NoError(t, err)
		_, err = ks.Import(exportJSON, cltest.Password)
		assert.Error(t, err)
		_, err = ks.Import([]byte(""), cltest.Password)
		assert.Error(t, err)
		require.Equal(t, key.PeerID(), importedKey.PeerID())
		retrievedKey, err := ks.Get(key.PeerID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		newKey, err := p2pkey.NewV2()
		require.NoError(t, err)
		err = ks.Add(newKey)
		require.NoError(t, err)
		err = ks.Add(newKey)
		assert.Error(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		_, err = ks.Delete(newKey.PeerID())
		require.NoError(t, err)
		_, err = ks.Delete(newKey.PeerID())
		assert.Error(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.PeerID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		err := ks.EnsureKey()
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, 1, len(keys))

		err = ks.EnsureKey()
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, 1, len(keys))
	})

	t.Run("GetOrFirst", func(t *testing.T) {
		defer reset()
		_, err := ks.GetOrFirst("")
		require.Contains(t, err.Error(), "no p2p keys exist")
		id := p2pkey.PeerID("a0")
		_, err = ks.GetOrFirst(id)
		require.Contains(t, err.Error(), fmt.Sprintf("unable to find P2P key with id %s", id))
		k1, err := ks.Create()
		require.NoError(t, err)
		k2, err := ks.GetOrFirst("")
		require.NoError(t, err)
		require.Equal(t, k1, k2)
		k3, err := ks.GetOrFirst(k1.PeerID())
		require.NoError(t, err)
		require.Equal(t, k1, k3)
		_, err = ks.Create()
		require.NoError(t, err)
		_, err = ks.GetOrFirst("")
		require.Contains(t, err.Error(), "multiple p2p keys found")
		k4, err := ks.GetOrFirst(k1.PeerID())
		require.NoError(t, err)
		require.Equal(t, k1, k4)
	})

	t.Run("clears p2p_peers on delete", func(t *testing.T) {
		key, err := ks.Create()
		require.NoError(t, err)
		p2pPeer1 := ocrcommon.P2PPeer{
			ID:     cltest.NewPeerID().String(),
			Addr:   testutils.NewAddress().Hex(),
			PeerID: cltest.DefaultPeerID, // different p2p key
		}
		p2pPeer2 := ocrcommon.P2PPeer{
			ID:     cltest.NewPeerID().String(),
			Addr:   testutils.NewAddress().Hex(),
			PeerID: key.PeerID().Raw(),
		}
		const p2pTableName = "p2p_peers"
		sql := fmt.Sprintf(`INSERT INTO %s (id, addr, peer_id, created_at, updated_at)
		VALUES (:id, :addr, :peer_id, now(), now())
		RETURNING *;`, p2pTableName)
		stmt, err := db.PrepareNamed(sql)
		require.NoError(t, err)
		require.NoError(t, stmt.Get(&p2pPeer1, &p2pPeer1))
		require.NoError(t, stmt.Get(&p2pPeer2, &p2pPeer2))
		cltest.AssertCount(t, db, p2pTableName, 2)
		_, err = ks.Delete(key.PeerID())
		require.NoError(t, err)
		cltest.AssertCount(t, db, p2pTableName, 1)
	})

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		exportedKey := `{"publicKey":"fcc1fdebde28322dde17233fe7bd6dcde447d60d5cc1de518962deed102eea35","peerID":"p2p_12D3KooWSq2UZgSXvhGLG5uuAAmz1JNjxHMJViJB39aorvbbYo8p","crypto":{"cipher":"aes-128-ctr","ciphertext":"adb2dff72148a8cd467f6f06a03869e7cedf180cf2a4decdb86875b2e1cf3e58c4bd2b721ecdaa88a0825fa9abfc309bf32dbb35a5c0b6cb01ac89a956d78e0550eff351","cipherparams":{"iv":"6cc4381766a4efc39f762b2b8d09dfba"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"ff5055ae4cdcdc2d0404307d578262e2caeb0210f82db3a0ecbdba727c6f5259"},"mac":"d37e4f1dea98d85960ef3205099fc71741715ae56a3b1a8f9215a78de9b95595"}}`
		importedKey, err := ks.Import([]byte(exportedKey), "p4SsW0rD1!@#_")
		require.NoError(t, err)
		require.Equal(t, "12D3KooWSq2UZgSXvhGLG5uuAAmz1JNjxHMJViJB39aorvbbYo8p", importedKey.ID())
	})

	t.Run("returns V1 keys as V2", func(t *testing.T) {
		defer reset()
		defer require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_p2p_keys")))

		p1 := cltest.MustRandomP2PPeerID(t)
		err := utils.JustError(db.Exec(`INSERT INTO encrypted_p2p_keys (peer_id, pub_key, encrypted_priv_key, created_at, updated_at, deleted_at) VALUES ($1, $2, '{"cipher":"aes-128-ctr","ciphertext":"adb2dff72148a8cd467f6f06a03869e7cedf180cf2a4decdb86875b2e1cf3e58c4bd2b721ecdaa88a0825fa9abfc309bf32dbb35a5c0b6cb01ac89a956d78e0550eff351","cipherparams":{"iv":"6cc4381766a4efc39f762b2b8d09dfba"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"ff5055ae4cdcdc2d0404307d578262e2caeb0210f82db3a0ecbdba727c6f5259"},"mac":"d37e4f1dea98d85960ef3205099fc71741715ae56a3b1a8f9215a78de9b95595"}', NOW(), NOW(), NULL)`, p1.Pretty(), utils.NewHash()))
		require.NoError(t, err)

		keyStore.SetPassword("p4SsW0rD1!@#_")

		keys, err := ks.GetV1KeysAsV2()
		require.NoError(t, err)

		assert.Len(t, keys, 1)
		assert.Equal(t, fmt.Sprintf("P2PKeyV2{PrivateKey: <redacted>, PeerID: %s}", keys[0].PeerID().Raw()), keys[0].GoString())
	})
}
