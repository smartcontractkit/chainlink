package ocr2key

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
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
	XXXOldV1GenericKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		Keyring         []byte
		// missing ID
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
		newBundleBytes, err := bundle.Marshal() // marshalling migrates to a generic struct
		require.NoError(t, err)

		// new bundle == old bundle (only difference is <chain>Keyring == Keyring)
		var newRawBundle keyBundleRawData
		require.NoError(t, json.Unmarshal(newBundleBytes, &newRawBundle))
		assert.Equal(t, oldKey.ChainType, newRawBundle.ChainType)
		assert.Equal(t, oldKey.OffchainKeyring, newRawBundle.OffchainKeyring)
		assert.Equal(t, oldKey.EVMKeyring, newRawBundle.Keyring)

		// test unmarshalling again to ensure ID has not changed
		// the underlying bytes have changed, but ID should be preserved
		newBundle := newKeyBundle(&evmKeyring{})
		require.NoError(t, newBundle.Unmarshal(newBundleBytes))
		assert.Equal(t, bundle.ID(), newBundle.ID())
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
		var newRawBundle keyBundleRawData
		require.NoError(t, json.Unmarshal(newBundleBytes, &newRawBundle))
		assert.Equal(t, oldKey.ChainType, newRawBundle.ChainType)
		assert.Equal(t, oldKey.OffchainKeyring, newRawBundle.OffchainKeyring)
		assert.Equal(t, oldKey.SolanaKeyring, newRawBundle.Keyring)

		// test unmarshalling again to ensure ID has not changed
		// the underlying bytes have changed, but ID should be preserved
		newBundle := newKeyBundle(&solanaKeyring{})
		require.NoError(t, newBundle.Unmarshal(newBundleBytes))
		assert.Equal(t, bundle.ID(), newBundle.ID())
	})

	t.Run("Cosmos", func(t *testing.T) {
		// onchain key
		bundle, err := newKeyBundleRand(chaintype.Cosmos, newCosmosKeyring)
		require.NoError(t, err)
		bundleBytes, err := bundle.Marshal()
		require.NoError(t, err)

		// test unmarshalling again to ensure ID has not changed
		// the underlying bytes have changed, but ID should be preserved
		otherBundle := newKeyBundle(&cosmosKeyring{})
		require.NoError(t, otherBundle.Unmarshal(bundleBytes))
		assert.Equal(t, bundle.ID(), otherBundle.ID())
	})

	t.Run("MissingID", func(t *testing.T) {
		// onchain key
		onKey, err := newEVMKeyring(cryptorand.Reader)
		require.NoError(t, err)
		onBytes, err := onKey.Marshal()
		require.NoError(t, err)

		// build key without ID parameter
		oldKey := XXXOldV1GenericKeyBundleRawData{
			ChainType:       chaintype.EVM,
			OffchainKeyring: offBytes,
			Keyring:         onBytes,
		}
		bundleBytes, err := json.Marshal(oldKey)
		require.NoError(t, err)

		// unmarshal first time to generate ID
		bundle := newKeyBundle(&evmKeyring{})
		require.NoError(t, bundle.Unmarshal(bundleBytes))

		// marshal and unmarshal again
		// different bytes generated, ID should not change
		newBundleBytes, err := bundle.Marshal()
		require.NoError(t, err)
		newBundle := newKeyBundle(&evmKeyring{})
		require.NoError(t, newBundle.Unmarshal(newBundleBytes))
		assert.Equal(t, bundle.ID(), newBundle.ID())
	})
}
