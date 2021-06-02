package keystore_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EthKeyStore_CreateNewKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	_, err := ethKeyStore.CreateNewKey()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ethKeyStore.Unlock(cltest.Password)
	assert.NoError(t, err)

	k, err := ethKeyStore.CreateNewKey()
	assert.NoError(t, err)

	has, err := ethKeyStore.HasSendingKeyWithAddress(k.Address.Address())
	require.NoError(t, err)
	assert.True(t, has)

	cltest.AssertCount(t, store, ethkey.Key{}, 1)
}

func Test_EthKeyStore_Unlock(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	k := cltest.MustInsertRandomKey(t, store.DB)

	_, err := ethKeyStore.SendingKeys()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	assert.EqualError(t, ethKeyStore.Unlock("wrong phrase"), fmt.Sprintf("invalid password for account %s; could not decrypt key with given password", k.Address.Hex()))
	assert.NoError(t, ethKeyStore.Unlock(cltest.Password))
}

func Test_EthKeyStore_KeyByAddress(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	_, address := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	key, err := ethKeyStore.KeyByAddress(address)
	require.NoError(t, err)
	require.Equal(t, address, key.Address.Address())

	missingAddress := cltest.NewAddress()
	_, err = ethKeyStore.KeyByAddress(missingAddress)
	require.EqualError(t, err, fmt.Sprintf("address %s not in keystore", missingAddress.Hex()))
}

func Test_EthKeyStore_EnsureFundingKey(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	cltest.AssertCount(t, store, ethkey.Key{}, 0)

	_, _, err := ethKeyStore.EnsureFundingKey()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	require.NoError(t, ethKeyStore.Unlock(cltest.Password))

	k, didExist, err := ethKeyStore.EnsureFundingKey()
	require.NoError(t, err)
	require.False(t, didExist)
	require.True(t, k.IsFunding)

	cltest.AssertCount(t, store, ethkey.Key{}, 1)
}

func Test_EthKeyStore_ImportKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	keyBytes := []byte(`{"address":"72f4f206d41339921570e47409cfef89ad528605","crypto":{"cipher":"aes-128-ctr","ciphertext":"d55d1cf27b464a7262e947fc6b09161c9c56b2efb1a2e6aef8b1ed0c22e02143","cipherparams":{"iv":"ff9effce7ce8318f54029c30e5e60c3a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":2,"p":2,"r":8,"salt":"bdec27593d039aca0fe87047bf425bd603a6eb134b8f04ee993ef090086300f7"},"mac":"5e06e90baef19112fcc301fb708d20577af9220e8b1f72329f9f06a70aade18e"},"id":"ec04d5fc-49ce-4d98-bdce-13d1dfa89eb9","version":3}`)

	_, err := ethKeyStore.ImportKey(keyBytes, cltest.Password)
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ethKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	keys, err := ethKeyStore.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)

	_, err = ethKeyStore.ImportKey(keyBytes, "wrong password")
	require.EqualError(t, err, "EthKeyStore#ImportKey failed to decrypt key: could not decrypt key with given password")

	k, err := ethKeyStore.ImportKey(keyBytes, cltest.Password)
	assert.NoError(t, err)

	keys, err = ethKeyStore.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, k.Address, keys[0].Address)
}

func Test_EthKeyStore_ExportKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	k := cltest.MustInsertRandomKey(t, store.DB)

	_, err := ethKeyStore.ExportKey(cltest.NewAddress(), "some password")
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ethKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	keys, err := ethKeyStore.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	bytes, err := ethKeyStore.ExportKey(k.Address.Address(), "new password")
	require.NoError(t, err)

	var addr struct {
		Address string `json:"address"`
	}
	err = json.Unmarshal(bytes, &addr)
	require.NoError(t, err)

	require.Equal(t, k.Address.Address(), common.HexToAddress("0x"+addr.Address))

	// Check it can be decrypted with new password
	_, err = gethkeystore.DecryptKey(bytes, "new password")
	assert.NoError(t, err)
}

func Test_EthKeyStore_AddKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	ks := ethKeyStore

	key := ethkey.Key{}

	err := ks.AddKey(&key)
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ks.Unlock(cltest.Password)
	require.NoError(t, err)

	err = ks.AddKey(&key)
	assert.EqualError(t, err, "unable to decrypt key JSON with keystore password: unexpected end of JSON input")

	key = cltest.MustGenerateRandomKey(t)

	err = ks.AddKey(&key)
	assert.NoError(t, err)
	assert.Greater(t, key.ID, int32(0))
	assert.True(t, key.CreatedAt.After(time.Time{}))
}

func Test_EthKeyStore_RemoveKey(t *testing.T) {
	t.Run("hard delete", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		_, err := ethKeyStore.RemoveKey(cltest.NewAddress(), false)
		require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

		k := cltest.MustInsertRandomKey(t, store.DB)

		err = ethKeyStore.Unlock(cltest.Password)
		require.NoError(t, err)

		keys, err := ethKeyStore.AllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 1)

		deleted, err := ethKeyStore.RemoveKey(k.Address.Address(), true)
		require.NoError(t, err)

		assert.Equal(t, k.Address, deleted.Address)

		keys, err = ethKeyStore.AllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 0)

		cltest.AssertCount(t, store, ethkey.Key{}, 0)
	})

	t.Run("soft delete", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		_, err := ethKeyStore.RemoveKey(cltest.NewAddress(), false)
		require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

		k := cltest.MustInsertRandomKey(t, store.DB)

		err = ethKeyStore.Unlock(cltest.Password)
		require.NoError(t, err)

		keys, err := ethKeyStore.AllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 1)

		deleted, err := ethKeyStore.RemoveKey(k.Address.Address(), false)
		require.NoError(t, err)

		assert.Equal(t, k.Address, deleted.Address)

		keys, err = ethKeyStore.AllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 0)

		cltest.AssertCount(t, store, ethkey.Key{}, 1)

		// Does not load soft deleted keys on a subsequent unlock
		ks := keystore.New(store.DB, utils.FastScryptParams).Eth
		err = ks.Unlock(cltest.Password)
		require.NoError(t, err)
		keys, err = ks.AllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 0)
	})
}

func Test_EthKeyStore_SignTx(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	k := cltest.MustInsertRandomKey(t, store.DB)

	chainID := big.NewInt(eth.NullClientChainID)
	tx := types.NewTransaction(0, cltest.NewAddress(), big.NewInt(53), 21000, big.NewInt(1000000000), []byte{1, 2, 3, 4})

	_, err := ethKeyStore.SignTx(cltest.NewAddress(), tx, chainID)
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ethKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	randomAddress := cltest.NewAddress()
	_, err = ethKeyStore.SignTx(randomAddress, tx, chainID)
	require.EqualError(t, err, fmt.Sprintf("address %s not in keystore", randomAddress.Hex()))

	signed, err := ethKeyStore.SignTx(k.Address.Address(), tx, chainID)
	require.NoError(t, err)

	assert.NotEqual(t, tx, signed)
}

func Test_EthKeyStore_AllKeys_SendingKeys_FundingKeys_HasSendingKeyWithAddress_GetKeyByAddress(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	sending1 := cltest.MustInsertRandomKey(t, store.DB, false)
	cltest.MustInsertRandomKey(t, store.DB, false)
	funding1 := cltest.MustInsertRandomKey(t, store.DB, true)

	_, err := ethKeyStore.AllKeys()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())
	_, err = ethKeyStore.SendingKeys()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())
	_, err = ethKeyStore.FundingKeys()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())
	_, err = ethKeyStore.HasSendingKeyWithAddress(cltest.NewAddress())
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	err = ethKeyStore.Unlock(cltest.Password)
	assert.NoError(t, err)

	keys, err := ethKeyStore.AllKeys()
	require.NoError(t, err)
	assert.Len(t, keys, 3)
	keys, err = ethKeyStore.SendingKeys()
	require.NoError(t, err)
	assert.Len(t, keys, 2)
	keys, err = ethKeyStore.FundingKeys()
	require.NoError(t, err)
	assert.Len(t, keys, 1)

	has, err := ethKeyStore.HasSendingKeyWithAddress(cltest.NewAddress())
	require.NoError(t, err)
	assert.False(t, has)
	has, err = ethKeyStore.HasSendingKeyWithAddress(funding1.Address.Address())
	require.NoError(t, err)
	assert.False(t, has)
	has, err = ethKeyStore.HasSendingKeyWithAddress(sending1.Address.Address())
	require.NoError(t, err)
	assert.True(t, has)
}

func Test_EthKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	kst := ethKeyStore

	k := []ethkey.Key{
		cltest.MustInsertRandomKey(t, store.DB, true),
		cltest.MustInsertRandomKey(t, store.DB),
		cltest.MustInsertRandomKey(t, store.DB),
		cltest.MustInsertRandomKey(t, store.DB),
	}

	_, err := kst.GetRoundRobinAddress()
	require.EqualError(t, err, keystore.ErrKeyStoreLocked.Error())

	require.NoError(t, kst.Unlock(cltest.Password))

	t.Run("with no address filter, rotates between all sending addresses", func(t *testing.T) {
		address, err := kst.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k[1].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k[2].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k[3].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k[1].Address.Hex(), address.Hex())
	})

	t.Run("with address filter, rotates between given addresses that match sending keys", func(t *testing.T) {
		// k0 is a funding address so even though it's whitelisted, it will be ignored
		addresses := []common.Address{k[0].Address.Address(), k[1].Address.Address(), k[2].Address.Address(), cltest.NewAddress()}

		// Last returned was k[1] so expect k[2] here
		address, err := kst.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k[2].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k[1].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k[2].Address.Hex(), address.Hex())

		address, err = kst.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k[1].Address.Hex(), address.Hex())
	})

	t.Run("with address filter when no address matches", func(t *testing.T) {
		_, err := kst.GetRoundRobinAddress([]common.Address{cltest.NewAddress()}...)
		require.Error(t, err)
		require.Equal(t, "no keys available", err.Error())
	})
}

// Does not require Unlock
func Test_EthKeyStore_HasDBSendingKeys(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	kst := ethKeyStore

	has, err := kst.HasDBSendingKeys()
	require.NoError(t, err)
	require.False(t, has)

	cltest.MustInsertRandomKey(t, store.DB, true)

	has, err = kst.HasDBSendingKeys()
	require.NoError(t, err)
	require.False(t, has)

	cltest.MustInsertRandomKey(t, store.DB, false)

	has, err = kst.HasDBSendingKeys()
	require.NoError(t, err)
	require.True(t, has)

}

func Test_EthKeyStore_ImportKeyFileToDB(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	kst := ethKeyStore

	path := "../../internal/fixtures/keys/7fc66c61f88A61DFB670627cA715Fe808057123e.json"
	k, err := kst.ImportKeyFileToDB(path)
	require.NoError(t, err)
	require.Equal(t, "0x7fc66c61f88A61DFB670627cA715Fe808057123e", k.Address.Hex())

	// importing again simply upserts
	_, err = kst.ImportKeyFileToDB(path)
	require.NoError(t, err)

	var keys []ethkey.Key
	err = store.DB.Find(&keys).Error
	require.NoError(t, err)

	require.Len(t, keys, 1)
	require.Equal(t, "0x7fc66c61f88A61DFB670627cA715Fe808057123e", keys[0].Address.String())
}
