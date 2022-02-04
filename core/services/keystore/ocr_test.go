package keystore_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
)

func Test_OCRKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	keyStore.Unlock(cltest.Password)
	ks := keyStore.OCR()
	reset := func() {
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		keyStore.ResetXXXTestOnly()
		keyStore.Unlock(cltest.Password)
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("errors when getting non-existant ID", func(t *testing.T) {
		defer reset()
		_, err := ks.Get("non-existant-id")
		require.Error(t, err)
	})

	t.Run("creates a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Delete(key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(exportJSON, cltest.Password)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		newKey, err := ocrkey.NewV2()
		require.NoError(t, err)
		err = ks.Add(newKey)
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		_, err = ks.Delete(newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		err := ks.EnsureKey()
		require.NoError(t, err)
		err = ks.EnsureKey()
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
	})

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		exportedKey := `{"id":"7cfd89bbb018e4778a44fd61172e8834dd24b4a2baf61ead795143b117221c61","onChainSigningAddress":"ocrsad_0x2ed5b18b62dacd7a85b6ed19247ea718bdae6114","offChainPublicKey":"ocroff_62a76d04e13dae5870071badea6b113a5123f4ac1a2cbae6b2fb7070dd9dbf2d","configPublicKey":"ocrcfg_75581baab36744671c2b1d75071b07b08b9cb631b3a7155d2f590744983d9c41","crypto":{"cipher":"aes-128-ctr","ciphertext":"60d2e679f08e0b1538cf609e25f2d32c0b7d408f24cab22dd05bffd3b5580c65552097e203f6546e2d792a4f6adb69449fee0fe4dd7f1060970907518e7c33331abd076388af842f03d05c193b03f22f6bf0423d4ae99dbb563c7158b4eac2a31b03c90fb9fd7be217804243151c36c33504469632bc2c89be33e7b9157edf172a52af4d49fa125b8d0358ea63ace90bc181a7164b548e0f12288ec08b919b46afad1b36dbaeda32d8d657a43908f802b6f2354473f538437ba3bd0b0d374d8e836e623484b655c95f4ef11e30baaa47b9075c6dbb53147c4b489f45a4bdcfa6b56ef2e6eaa9e9b88b570517c991de359d7f07226c00259810a8a4196b7d5331e4126529eac9bd80b47b5540940f89ad0e728b3dd50e6da316d9f3cf9b3be9b87ca6b7868daa7e4142fc4a65fc77deea6f4f2b4bce1e38337aa827160d8c50cad92d157309aa251180b894ab1ca9923d709d","cipherparams":{"iv":"a9507e6f2b073c1da1082d40a24864d1"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"267f9450f52af42a918ab5747043c88bd2035fa3d3e0f0cfd2b621981bc9320f"},"mac":"15aeb3fc1903f514bfe70cb2eb5a23820ba904f5edf8aeb1913d447797f74442"}}`
		importedKey, err := ks.Import([]byte(exportedKey), cltest.Password)
		require.NoError(t, err)
		require.Equal(t, "7cfd89bbb018e4778a44fd61172e8834dd24b4a2baf61ead795143b117221c61", importedKey.ID())
	})
}
