package keystore_test

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_EthKeyStore(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	err := keyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	ethKeyStore := keyStore.Eth()
	reset := func() {
		keyStore.ResetXXXTestOnly()
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM eth_key_states")))
		keyStore.Unlock(cltest.Password)
	}
	const statesTableName = "eth_key_states"

	t.Run("Create / GetAll / Get", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(retrievedKeys))
		require.Equal(t, key.Address, retrievedKeys[0].Address)
		foundKey, err := ethKeyStore.Get(key.Address.Hex())
		require.NoError(t, err)
		require.Equal(t, key, foundKey)
		// adds ethkey.State
		cltest.AssertCount(t, db, statesTableName, 1)
		var state ethkey.State
		sql := fmt.Sprintf(`SELECT * from %s LIMIT 1`, statesTableName)
		require.NoError(t, db.Get(&state, sql))
		require.Equal(t, state.Address, retrievedKeys[0].Address)
		// adds key to db
		keyStore.ResetXXXTestOnly()
		keyStore.Unlock(cltest.Password)
		retrievedKeys, err = ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(retrievedKeys))
		require.Equal(t, key.Address, retrievedKeys[0].Address)
		// adds 2nd key
		_, err = ethKeyStore.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKeys, err = ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 2, len(retrievedKeys))
	})

	t.Run("GetAll ordering", func(t *testing.T) {
		defer reset()
		var keys []ethkey.KeyV2
		for i := 0; i < 5; i++ {
			key, err := ethKeyStore.Create(&cltest.FixtureChainID)
			require.NoError(t, err)
			keys = append(keys, key)
		}
		retrievedKeys, err := ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 5, len(retrievedKeys))

		sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })

		assert.Equal(t, keys, retrievedKeys)
	})

	t.Run("RemoveKey", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		_, err = ethKeyStore.Delete(key.ID())
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(retrievedKeys))
		cltest.AssertCount(t, db, statesTableName, 0)
	})

	t.Run("EnsureKeys / SendingKeys", func(t *testing.T) {
		defer reset()
		err := ethKeyStore.EnsureKeys(&cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys1, err := ethKeyStore.SendingKeys(nil)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys1))
		cltest.AssertCount(t, db, statesTableName, 2)

		err = ethKeyStore.EnsureKeys(&cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys2, err := ethKeyStore.SendingKeys(nil)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys2))
		require.Equal(t, sendingKeys1, sendingKeys2)
	})

	t.Run("SendingKeys with specified chain ID", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(testutils.FixtureChainID)
		require.NoError(t, err)
		key2, err := ethKeyStore.Create(big.NewInt(1337))
		require.NoError(t, err)

		keys, err := ethKeyStore.SendingKeys(testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key, keys[0])

		keys, err = ethKeyStore.SendingKeys(big.NewInt(1337))
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key2, keys[0])

		keys, err = ethKeyStore.SendingKeys(nil)
		require.NoError(t, err)
		require.Len(t, keys, 2)
	})
}

func Test_EthKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)

	keyStore := cltest.NewKeyStore(t, db, cfg)
	ethKeyStore := keyStore.Eth()

	t.Run("should error when no addresses", func(t *testing.T) {
		_, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.Error(t, err)
	})

	// create 4 keys - 1 funding and 2 sending
	err := ethKeyStore.EnsureKeys(&cltest.FixtureChainID)
	require.NoError(t, err)
	sendingKeys, err := ethKeyStore.SendingKeys(nil)
	assert.NoError(t, err)

	k1 := sendingKeys[0]

	k2, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	// create 1 funding and 1 sending key for a different chain
	err = ethKeyStore.EnsureKeys(testutils.SimulatedChainID)
	require.NoError(t, err)

	sendingKeys, err = ethKeyStore.SendingKeys(nil)
	assert.NoError(t, err)
	require.Equal(t, 3, len(sendingKeys))

	fundingKeys, err := ethKeyStore.FundingKeys()
	assert.NoError(t, err)
	require.Equal(t, 1, len(fundingKeys))

	t.Run("with no address filter, rotates between all sending addresses", func(t *testing.T) {
		address1, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)
		address2, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)
		address3, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)
		address4, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)
		address5, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)
		address6, err := ethKeyStore.GetRoundRobinAddress(nil)
		require.NoError(t, err)

		require.NotEqual(t, address1, address2)
		require.NotEqual(t, address2, address3)
		require.NotEqual(t, address1, address3)
		require.Equal(t, address1, address4)
		require.Equal(t, address2, address5)
		require.Equal(t, address3, address6)
	})

	t.Run("with address filter, rotates between given addresses that match sending keys", func(t *testing.T) {
		// fundingKeys[0] is a funding address so even though it's whitelisted, it will be ignored
		addresses := []common.Address{fundingKeys[0].Address.Address(), k1.Address.Address(), k2.Address.Address(), testutils.NewAddress()}

		address1, err := ethKeyStore.GetRoundRobinAddress(nil, addresses...)
		require.NoError(t, err)
		address2, err := ethKeyStore.GetRoundRobinAddress(nil, addresses...)
		require.NoError(t, err)
		address3, err := ethKeyStore.GetRoundRobinAddress(nil, addresses...)
		require.NoError(t, err)
		address4, err := ethKeyStore.GetRoundRobinAddress(nil, addresses...)
		require.NoError(t, err)

		require.True(t, address1 == k1.Address.Address() || address1 == k2.Address.Address())
		require.True(t, address2 == k1.Address.Address() || address2 == k2.Address.Address())
		require.NotEqual(t, address1, address2)
		require.Equal(t, address1, address3)
		require.Equal(t, address2, address4)
	})

	t.Run("with address filter when no address matches", func(t *testing.T) {
		addr := testutils.NewAddress()
		_, err := ethKeyStore.GetRoundRobinAddress(nil, []common.Address{addr}...)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("no sending keys available that match whitelist: [%s]", addr.Hex()), err.Error())
	})

	t.Run("with non-nil chain ID, filters by chain ID", func(t *testing.T) {
		sendingKeys, err := ethKeyStore.SendingKeys(testutils.SimulatedChainID)
		assert.NoError(t, err)
		require.Len(t, sendingKeys, 1)
		k := sendingKeys[0]
		address1, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID)
		require.NoError(t, err)
		address2, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID)
		require.NoError(t, err)
		address3, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID)
		require.NoError(t, err)

		require.Equal(t, k.Address.Address(), address1)
		require.Equal(t, k.Address.Address(), address2)
		require.Equal(t, k.Address.Address(), address3)
	})
}

func Test_EthKeyStore_SignTx(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	config := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	ethKeyStore := keyStore.Eth()

	k, _ := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	chainID := big.NewInt(evmclient.NullClientChainID)
	tx := types.NewTransaction(0, testutils.NewAddress(), big.NewInt(53), 21000, big.NewInt(1000000000), []byte{1, 2, 3, 4})

	randomAddress := testutils.NewAddress()
	_, err := ethKeyStore.SignTx(randomAddress, tx, chainID)
	require.EqualError(t, err, fmt.Sprintf("unable to find eth key with id %s", randomAddress.Hex()))

	signed, err := ethKeyStore.SignTx(k.Address.Address(), tx, chainID)
	require.NoError(t, err)

	require.NotEqual(t, tx, signed)
}

func Test_EthKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	err := keyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	ks := keyStore.Eth()
	reset := func() {
		keyStore.ResetXXXTestOnly()
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM eth_key_states")))
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
		key, err := ks.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Delete(key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(exportJSON, cltest.Password, &cltest.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		newKey, err := ethkey.NewV2()
		require.NoError(t, err)
		err = ks.Add(newKey, &cltest.FixtureChainID)
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

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		exportedKey := `{"address":"0dd359b4f22a30e44b2fd744b679971941865820","crypto":{"cipher":"aes-128-ctr","ciphertext":"b30af964a3b3f37894e599446b4cf2314bbfcd1062e6b35b620d3d20bd9965cc","cipherparams":{"iv":"58a8d75629cc1945da7cf8c24520d1dc"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c352887e9d427d8a6a1869082619b73fac4566082a99f6e367d126f11b434f28"},"mac":"fd76a588210e0bf73d01332091e0e83a4584ee2df31eaec0e27f9a1b94f024b4"},"id":"a5ee0802-1d7b-45b6-aeb8-ea8a3351e715","version":3}`
		importedKey, err := ks.Import([]byte(exportedKey), cltest.Password, &cltest.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, "0x0dd359b4f22a30E44b2fD744B679971941865820", importedKey.ID())
	})
}

func Test_EthKeyStore_SubscribeToKeyChanges(t *testing.T) {
	chDone := make(chan struct{})
	defer func() { close(chDone) }()
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()
	chSub, unsubscribe := ks.SubscribeToKeyChanges()
	defer unsubscribe()

	count := atomic.NewInt32(0)

	assertCount := func(expected int32) {
		require.Eventually(
			t,
			func() bool { return count.Load() == expected },
			10*time.Second,
			100*time.Millisecond,
			fmt.Sprintf("insufficient number of callbacks triggered. Expected %d, got %d", expected, count.Load()),
		)
	}

	go func() {
		for {
			select {
			case _, ok := <-chSub:
				if !ok {
					return
				}
				count.Add(1)
			case <-chDone:
				return
			}
		}
	}()

	err := ks.EnsureKeys(&cltest.FixtureChainID)
	require.NoError(t, err)
	assertCount(1)
	_, err = ks.Create(&cltest.FixtureChainID)
	require.NoError(t, err)
	assertCount(2)
	newKey, err := ethkey.NewV2()
	require.NoError(t, err)
	err = ks.Add(newKey, &cltest.FixtureChainID)
	require.NoError(t, err)
	assertCount(3)
	_, err = ks.Delete(newKey.ID())
	require.NoError(t, err)
	assertCount(4)
}
