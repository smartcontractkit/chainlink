package keystore_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
)

func Test_TronKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.Tron()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(ctx, cltest.Password))
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

	t.Run("creates a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		key, err := ks.Create(ctx)
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		requireEqualKeys(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		key, err := ks.Create(ctx)
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Export("non-existent", cltest.Password)
		require.Error(t, err)
		_, err = ks.Delete(ctx, key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(ctx, exportJSON, cltest.Password)
		require.NoError(t, err)
		_, err = ks.Import(ctx, exportJSON, cltest.Password)
		require.Error(t, err)
		_, err = ks.Import(ctx, []byte(""), cltest.Password)
		require.Error(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		requireEqualKeys(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey, err := tronkey.New()
		require.NoError(t, err)
		err = ks.Add(ctx, newKey)
		require.NoError(t, err)
		err = ks.Add(ctx, newKey)
		require.Error(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
		_, err = ks.Delete(ctx, newKey.ID())
		require.NoError(t, err)
		_, err = ks.Delete(ctx, newKey.ID())
		require.Error(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Empty(t, keys)
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKey(ctx)
		require.NoError(t, err)

		err = ks.EnsureKey(ctx)
		require.NoError(t, err)

		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("sign tx", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey, err := tronkey.New()
		require.NoError(t, err)
		require.NoError(t, ks.Add(ctx, newKey))

		// sign unknown ID
		_, err = ks.Sign(testutils.Context(t), "not-real", nil)
		require.Error(t, err)

		// sign known key

		// Create a mock transaction
		mockTx := createMockTronTransaction(newKey.PublicKeyStr(), "TJRabPrwbZy45sbavfcjinPJC18kjpRTv8", 1000000)
		serializedTx, err := serializeMockTransaction(mockTx)
		require.NoError(t, err)

		hash := sha256.Sum256(serializedTx)
		txHash := hash[:]
		sig, err := ks.Sign(testutils.Context(t), newKey.ID(), txHash)
		require.NoError(t, err)

		directSig, err := newKey.Sign(txHash)
		require.NoError(t, err)

		// signatures should match using keystore sign or key sign
		require.Equal(t, directSig, sig)
	})
}

// MockTronTransaction represents a mock TRON transaction
// This is based on https://developers.tron.network/docs/tron-protocol-transaction
type MockTronTransaction struct {
	RawData struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Amount       int64  `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		Timestamp     int64  `json:"timestamp"`
		FeeLimit      int64  `json:"fee_limit"`
	} `json:"raw_data"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
}

// CreateMockTronTransaction generates a mock TRON transaction for testing
func createMockTronTransaction(ownerAddress, toAddress string, amount int64) MockTronTransaction {
	return MockTronTransaction{
		RawData: struct {
			Contract []struct {
				Parameter struct {
					Value struct {
						Amount       int64  `json:"amount"`
						OwnerAddress string `json:"owner_address"`
						ToAddress    string `json:"to_address"`
					} `json:"value"`
					TypeURL string `json:"type_url"`
				} `json:"parameter"`
				Type string `json:"type"`
			} `json:"contract"`
			RefBlockBytes string `json:"ref_block_bytes"`
			RefBlockHash  string `json:"ref_block_hash"`
			Expiration    int64  `json:"expiration"`
			Timestamp     int64  `json:"timestamp"`
			FeeLimit      int64  `json:"fee_limit"`
		}{
			Contract: []struct {
				Parameter struct {
					Value struct {
						Amount       int64  `json:"amount"`
						OwnerAddress string `json:"owner_address"`
						ToAddress    string `json:"to_address"`
					} `json:"value"`
					TypeURL string `json:"type_url"`
				} `json:"parameter"`
				Type string `json:"type"`
			}{
				{
					Parameter: struct {
						Value struct {
							Amount       int64  `json:"amount"`
							OwnerAddress string `json:"owner_address"`
							ToAddress    string `json:"to_address"`
						} `json:"value"`
						TypeURL string `json:"type_url"`
					}{
						Value: struct {
							Amount       int64  `json:"amount"`
							OwnerAddress string `json:"owner_address"`
							ToAddress    string `json:"to_address"`
						}{
							Amount:       amount,
							OwnerAddress: ownerAddress,
							ToAddress:    toAddress,
						},
						TypeURL: "type.googleapis.com/protocol.TransferContract",
					},
					Type: "TransferContract",
				},
			},
			RefBlockBytes: "1234",
			RefBlockHash:  "abcdef0123456789",
			Expiration:    time.Now().Unix() + 60*60,
			Timestamp:     time.Now().Unix(),
			FeeLimit:      10000000,
		},
	}
}

func serializeMockTransaction(tx MockTronTransaction) ([]byte, error) {
	return json.Marshal(tx)
}
