package ocr2key

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
)

type (
	XXXOldEVMKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		EVMKeyring      []byte
	}
	XXXOldSolanaKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		SolanaKeyring   []byte
	}
	XXXOldTerraKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		TerraKeyring    []byte
	}
)

func TestGenericKeyBundle_Migrate_UnmarshalMarshal(t *testing.T) {
	// offchain key
	offKey, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)
	offBytes, err := offKey.marshal()
	require.NoError(t, err)

	t.Run("EVM", func(t *testing.T) {
		// onchain key
		onKey, err := newEVMKeyring(cryptorand.Reader)
		require.NoError(t, err)
		onBytes, err := onKey.Marshal()
		require.NoError(t, err)

		// marshal old key format
		oldKey := XXXOldEVMKeyBundleRawData{
			ChainType:       chaintype.EVM,
			OffchainKeyring: offBytes,
			EVMKeyring:      onBytes,
		}
		bundleBytes, err := json.Marshal(oldKey)
		require.NoError(t, err)

		// test Unmarshal with old raw bundle
		bundle := newKeyBundle(&evmKeyring{})
		require.NoError(t, bundle.Unmarshal(bundleBytes))
		newBundleBytes, err := bundle.Marshal()
		require.NoError(t, err)

		// new bundle == old bundle (only difference is <chain>Keyring == Keyring)
		var newBundle keyBundleRawData
		require.NoError(t, json.Unmarshal(newBundleBytes, &newBundle))
		assert.Equal(t, oldKey.ChainType, newBundle.ChainType)
		assert.Equal(t, oldKey.OffchainKeyring, newBundle.OffchainKeyring)
		assert.Equal(t, oldKey.EVMKeyring, newBundle.Keyring)
	})

	t.Run("Solana", func(t *testing.T) {
		// onchain key
		onKey, err := newSolanaKeyring(cryptorand.Reader)
		require.NoError(t, err)
		onBytes, err := onKey.Marshal()
		require.NoError(t, err)

		// marshal old key format
		oldKey := XXXOldSolanaKeyBundleRawData{
			ChainType:       chaintype.Solana,
			OffchainKeyring: offBytes,
			SolanaKeyring:   onBytes,
		}
		bundleBytes, err := json.Marshal(oldKey)
		require.NoError(t, err)

		// test Unmarshal with old raw bundle
		bundle := newKeyBundle(&solanaKeyring{})
		require.NoError(t, bundle.Unmarshal(bundleBytes))
		newBundleBytes, err := bundle.Marshal()
		require.NoError(t, err)

		// new bundle == old bundle (only difference is <chain>Keyring == Keyring)
		var newBundle keyBundleRawData
		require.NoError(t, json.Unmarshal(newBundleBytes, &newBundle))
		assert.Equal(t, oldKey.ChainType, newBundle.ChainType)
		assert.Equal(t, oldKey.OffchainKeyring, newBundle.OffchainKeyring)
		assert.Equal(t, oldKey.SolanaKeyring, newBundle.Keyring)
	})

	t.Run("Terra", func(t *testing.T) {
		// onchain key
		onKey, err := newTerraKeyring(cryptorand.Reader)
		require.NoError(t, err)
		onBytes, err := onKey.Marshal()
		require.NoError(t, err)

		// marshal old key format
		oldKey := XXXOldTerraKeyBundleRawData{
			ChainType:       chaintype.Terra,
			OffchainKeyring: offBytes,
			TerraKeyring:    onBytes,
		}
		bundleBytes, err := json.Marshal(oldKey)
		require.NoError(t, err)

		// test Unmarshal with old raw bundle
		bundle := newKeyBundle(&terraKeyring{})
		require.NoError(t, bundle.Unmarshal(bundleBytes))
		newBundleBytes, err := bundle.Marshal()
		require.NoError(t, err)

		// new bundle == old bundle (only difference is <chain>Keyring == Keyring)
		var newBundle keyBundleRawData
		require.NoError(t, json.Unmarshal(newBundleBytes, &newBundle))
		assert.Equal(t, oldKey.ChainType, newBundle.ChainType)
		assert.Equal(t, oldKey.OffchainKeyring, newBundle.OffchainKeyring)
		assert.Equal(t, oldKey.TerraKeyring, newBundle.Keyring)
	})
}
