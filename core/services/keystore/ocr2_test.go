package keystore_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kslib "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/ocr2offchain"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

func Test_OCR2KeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.OCR2()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
		_, err := db.Exec("DELETE FROM encrypted_key_rings")
		require.NoError(t, err)
		keyStore.ResetXXXTestOnly()
		err = keyStore.Unlock(ctx, cltest.Password)
		require.NoError(t, err)
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Empty(t, keys)
	})

	t.Run("errors when getting non-existent ID", func(t *testing.T) {
		defer reset()
		_, err := ks.Get("non-existent-id")
		require.Error(t, err)
	})

	t.Run("creates a key with valid type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		// lopp through different chain types
		for _, chain := range chaintype.SupportedChainTypes {
			key, err := ks.Create(ctx, chain)
			require.NoError(t, err)
			retrievedKey, err := ks.Get(key.ID())
			require.NoError(t, err)
			require.Equal(t, key, retrievedKey)
		}
	})

	t.Run("gets keys by type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		created := map[chaintype.ChainType]bool{}
		for _, chain := range chaintype.SupportedChainTypes {
			// validate no keys exist for chain
			keys, err := ks.GetAllOfType(chain)
			require.NoError(t, err)
			require.Empty(t, keys)

			_, err = ks.Create(ctx, chain)
			require.NoError(t, err)
			created[chain] = true

			// validate that only 1 of each exists after creation
			for _, c := range chaintype.SupportedChainTypes {
				keys, err := ks.GetAllOfType(c)
				require.NoError(t, err)
				if created[c] {
					require.Len(t, keys, 1)
					continue
				}
				require.Empty(t, keys)
			}
		}
	})

	t.Run("errors when creating a key with an invalid type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		_, err := ks.Create(ctx, "foobar")
		require.Error(t, err)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		for _, chain := range chaintype.SupportedChainTypes {
			key, err := ks.Create(ctx, chain)
			require.NoError(t, err)
			exportJSON, err := ks.Export(key.ID(), cltest.Password)
			require.NoError(t, err)
			_, err = ks.Export("non-existent", cltest.Password)
			assert.Error(t, err)
			err = ks.Delete(ctx, key.ID())
			require.NoError(t, err)
			_, err = ks.Get(key.ID())
			require.Error(t, err)
			importedKey, err := ks.Import(ctx, exportJSON, cltest.Password)
			require.NoError(t, err)
			_, err = ks.Import(ctx, []byte(""), cltest.Password)
			assert.Error(t, err)
			require.Equal(t, key.ID(), importedKey.ID())
			retrievedKey, err := ks.Get(key.ID())
			require.NoError(t, err)
			require.Equal(t, importedKey, retrievedKey)
			require.Equal(t, importedKey.ChainType(), retrievedKey.ChainType())
		}
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		for _, chain := range chaintype.SupportedChainTypes {
			newKey, err := ocr2key.New(chain)
			require.NoError(t, err)
			err = ks.Add(ctx, newKey)
			require.NoError(t, err)
			err = ks.Add(ctx, newKey)
			assert.Error(t, err)
			keys, err := ks.GetAll()
			require.NoError(t, err)
			require.Len(t, keys, 1)
			err = ks.Delete(ctx, newKey.ID())
			require.NoError(t, err)
			err = ks.Delete(ctx, newKey.ID())
			assert.Error(t, err)
			keys, err = ks.GetAll()
			require.NoError(t, err)
			require.Empty(t, keys)
			_, err = ks.Get(newKey.ID())
			require.Error(t, err)
		}
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKeys(ctx, chaintype.SupportedChainTypes...)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		require.Len(t, keys, len(chaintype.SupportedChainTypes))

		err = ks.EnsureKeys(ctx, chaintype.SupportedChainTypes...)
		assert.NoError(t, err)

		// loop through different supported chain types
		for _, chain := range chaintype.SupportedChainTypes {
			keys, err := ks.GetAllOfType(chain)
			assert.NoError(t, err)
			require.Len(t, keys, 1)
		}
	})

	t.Run("ensures key only for enabled chains", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKeys(ctx, chaintype.EVM)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, chaintype.EVM, keys[0].ChainType())

		err = ks.EnsureKeys(ctx, chaintype.Cosmos)
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)
		require.Len(t, keys, 2)

		cosmosKeys, err := ks.GetAllOfType(chaintype.Cosmos)
		assert.NoError(t, err)
		require.Len(t, cosmosKeys, 1)
		require.Equal(t, chaintype.Cosmos, cosmosKeys[0].ChainType())

		err = ks.EnsureKeys(ctx, chaintype.StarkNet)
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)
		require.Len(t, keys, 3)

		starknetKeys, err := ks.GetAllOfType(chaintype.StarkNet)
		require.NoError(t, err)
		require.Len(t, starknetKeys, 1)
		require.Equal(t, chaintype.StarkNet, starknetKeys[0].ChainType())

		err = ks.EnsureKeys(ctx, chaintype.Tron)
		require.NoError(t, err)

		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 4)

		tronKeys, err := ks.GetAllOfType(chaintype.Tron)
		require.NoError(t, err)
		require.Len(t, tronKeys, 1)
		require.Equal(t, chaintype.Tron, tronKeys[0].ChainType())
	})
}

func Test_OCR2KeyStore_KeystoreLibCompatibility(t *testing.T) {
	ctx := testutils.Context(t)
	password := "my-password"
	chainType := corekeys.ChainType("evm")

	ks, err := kslib.LoadKeystore(ctx, kslib.NewMemoryStorage(), password)
	require.NoError(t, err)

	corekeysKs := corekeys.NewKeystore(ks)

	encryptedBundle, err := corekeysKs.GenerateEncryptedOCRKeyBundle(ctx, chainType, password)
	require.NoError(t, err)

	signingKeyPath := kslib.NewKeyPath(ocr2offchain.PrefixOCR2Offchain, "default", ocr2offchain.OCR2OffchainSigning)
	encryptionKeyPath := kslib.NewKeyPath(ocr2offchain.PrefixOCR2Offchain, "default", ocr2offchain.OCR2OffchainEncryption)
	onchainKeyPath := kslib.NewKeyPath(corekeys.PrefixOCR2Onchain, "default", string(chainType))

	getKeysResp, err := ks.GetKeys(ctx, kslib.GetKeysRequest{
		KeyNames: []string{
			signingKeyPath.String(),
			encryptionKeyPath.String(),
			onchainKeyPath.String(),
		},
	})
	require.NoError(t, err)
	require.Len(t, getKeysResp.Keys, 3)

	var expectedSigningPubKey, expectedEncryptionPubKey, expectedOnchainPubKey []byte
	for _, key := range getKeysResp.Keys {
		switch key.KeyInfo.Name {
		case signingKeyPath.String():
			expectedSigningPubKey = key.KeyInfo.PublicKey
		case encryptionKeyPath.String():
			expectedEncryptionPubKey = key.KeyInfo.PublicKey
		case onchainKeyPath.String():
			expectedOnchainPubKey = key.KeyInfo.PublicKey
		}
	}

	require.NotEmpty(t, expectedSigningPubKey)
	require.NotEmpty(t, expectedEncryptionPubKey)
	require.NotEmpty(t, expectedOnchainPubKey)

	db := pgtest.NewSqlxDB(t)
	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(ctx, cltest.Password))

	ocr2KeyStore := keyStore.OCR2()

	importedKey, err := ocr2KeyStore.Import(ctx, encryptedBundle, password)
	require.NoError(t, err)

	assert.Equal(t, chaintype.EVM, importedKey.ChainType())

	importedSigningPubKey := importedKey.OffchainPublicKey()
	assert.Equal(t, expectedSigningPubKey, importedSigningPubKey[:])

	importedConfigPubKey := importedKey.ConfigEncryptionPublicKey()
	assert.Equal(t, expectedEncryptionPubKey, importedConfigPubKey[:])

	bundle, err := corekeys.FromEncryptedOCRKeyBundle(encryptedBundle, password)
	require.NoError(t, err)

	onchainPrivKey, err := crypto.ToECDSA(bundle.OnchainSigningKey)
	require.NoError(t, err)
	derivedOnchainPubKey := crypto.FromECDSAPub(&onchainPrivKey.PublicKey)
	require.Equal(t, expectedOnchainPubKey, derivedOnchainPubKey)

	expectedOnchainAddress := crypto.PubkeyToAddress(onchainPrivKey.PublicKey)

	importedOnchainPubKey := importedKey.PublicKey()
	assert.Equal(t, expectedOnchainAddress[:], []byte(importedOnchainPubKey))
}
