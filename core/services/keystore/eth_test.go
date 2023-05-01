package keystore_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM evm_key_states")))
		require.NoError(t, keyStore.Unlock(cltest.Password))
	}
	const statesTableName = "evm_key_states"

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
		sql := fmt.Sprintf(`SELECT address, disabled, evm_chain_id, created_at, updated_at from %s LIMIT 1`, statesTableName)
		require.NoError(t, db.Get(&state, sql))
		require.Equal(t, state.Address.Address(), retrievedKeys[0].Address)
		// adds key to db
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(cltest.Password))
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

	t.Run("Delete removes key even if eth_txes are present", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		// ensure at least one state is present
		cltest.AssertCount(t, db, statesTableName, 1)

		// add one eth_tx
		borm := cltest.NewTxStore(t, db, cfg)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 0, 42, key.Address)

		_, err = ethKeyStore.Delete(key.ID())
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(retrievedKeys))
		cltest.AssertCount(t, db, statesTableName, 0)
	})

	t.Run("EnsureKeys / EnabledKeysForChain", func(t *testing.T) {
		defer reset()
		err := ethKeyStore.EnsureKeys(&cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys1, err := ethKeyStore.EnabledKeysForChain(testutils.FixtureChainID)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys1))
		cltest.AssertCount(t, db, statesTableName, 1)

		err = ethKeyStore.EnsureKeys(&cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys2, err := ethKeyStore.EnabledKeysForChain(testutils.FixtureChainID)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys2))
		require.Equal(t, sendingKeys1, sendingKeys2)
	})

	t.Run("EnabledKeysForChain with specified chain ID", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(testutils.FixtureChainID)
		require.NoError(t, err)
		key2, err := ethKeyStore.Create(big.NewInt(1337))
		require.NoError(t, err)

		keys, err := ethKeyStore.EnabledKeysForChain(testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key, keys[0])

		keys, err = ethKeyStore.EnabledKeysForChain(big.NewInt(1337))
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key2, keys[0])

		_, err = ethKeyStore.EnabledKeysForChain(nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "chainID must be non-nil")
	})

	t.Run("EnabledAddressesForChain with specified chain ID", func(t *testing.T) {
		defer reset()
		key, err := ethKeyStore.Create(testutils.FixtureChainID)
		require.NoError(t, err)
		key2, err := ethKeyStore.Create(big.NewInt(1337))
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 2)
		keys, err := ethKeyStore.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 2)

		//get enabled addresses for FixtureChainID
		enabledAddresses, err := ethKeyStore.EnabledAddressesForChain(testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, enabledAddresses, 1)
		require.Equal(t, key.Address, enabledAddresses[0])

		//get enabled addresses for chain 1337
		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(big.NewInt(1337))
		require.NoError(t, err)
		require.Len(t, enabledAddresses, 1)
		require.Equal(t, key2.Address, enabledAddresses[0])

		// /get enabled addresses for nil chain ID
		_, err = ethKeyStore.EnabledAddressesForChain(nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "chainID must be non-nil")

		// disable the key for chain FixtureChainID
		err = ethKeyStore.Disable(key.Address, testutils.FixtureChainID)
		require.NoError(t, err)

		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(testutils.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, enabledAddresses, 0)
		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(big.NewInt(1337))
		require.NoError(t, err)
		assert.Len(t, enabledAddresses, 1)
		require.Equal(t, key2.Address, enabledAddresses[0])
	})
}

func Test_EthKeyStore_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)

	keyStore := cltest.NewKeyStore(t, db, cfg)
	ethKeyStore := keyStore.Eth()

	t.Run("should error when no addresses", func(t *testing.T) {
		_, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.Error(t, err)
	})

	// create keys
	// - key 1
	//   enabled - fixture
	//   enabled - simulated
	// - key 2
	//   enabled - fixture
	//   disabled - simulated
	// - key 3
	//   enabled - simulated
	// - key 4
	//   enabled - fixture
	k1, _ := cltest.MustInsertRandomKey(t, ethKeyStore, []utils.Big{})
	require.NoError(t, ethKeyStore.Enable(k1.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Enable(k1.Address, testutils.SimulatedChainID))

	k2, _ := cltest.MustInsertRandomKey(t, ethKeyStore, []utils.Big{})
	require.NoError(t, ethKeyStore.Enable(k2.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Enable(k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ethKeyStore.Disable(k2.Address, testutils.SimulatedChainID))

	k3, _ := cltest.MustInsertRandomKey(t, ethKeyStore, []utils.Big{})
	require.NoError(t, ethKeyStore.Enable(k3.Address, testutils.SimulatedChainID))

	k4, _ := cltest.MustInsertRandomKey(t, ethKeyStore, []utils.Big{})
	require.NoError(t, ethKeyStore.Enable(k4.Address, testutils.FixtureChainID))

	t.Run("with no address filter, rotates between all enabled addresses", func(t *testing.T) {
		address1, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)
		address2, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)
		address3, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)
		address4, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)
		address5, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)
		address6, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID)
		require.NoError(t, err)

		assert.NotEqual(t, address1, address2)
		assert.NotEqual(t, address2, address3)
		assert.NotEqual(t, address1, address3)
		assert.Equal(t, address1, address4)
		assert.Equal(t, address2, address5)
		assert.Equal(t, address3, address6)
	})

	t.Run("with address filter, rotates between given addresses that match sending keys", func(t *testing.T) {
		{
			// k3 is a disabled address for FixtureChainID so even though it's whitelisted, it will be ignored
			addresses := []common.Address{k4.Address, k3.Address, k1.Address, k2.Address, testutils.NewAddress()}

			address1, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address2, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address3, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address4, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID, addresses...)
			require.NoError(t, err)

			assert.NotEqual(t, k3.Address, address1)
			assert.NotEqual(t, k3.Address, address2)
			assert.NotEqual(t, k3.Address, address3)
			assert.NotEqual(t, address1, address2)
			assert.NotEqual(t, address1, address3)
			assert.NotEqual(t, address2, address3)
			assert.Equal(t, address1, address4)
		}

		{

			// k2 and k4 are disabled address for SimulatedChainID so even though it's whitelisted, it will be ignored
			addresses := []common.Address{k4.Address, k3.Address, k1.Address, k2.Address, testutils.NewAddress()}

			address1, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address2, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address3, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address4, err := ethKeyStore.GetRoundRobinAddress(testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)

			assert.True(t, address1 == k1.Address || address1 == k3.Address)
			assert.True(t, address2 == k1.Address || address2 == k3.Address)
			assert.NotEqual(t, address1, address2)
			assert.Equal(t, address1, address3)
			assert.Equal(t, address2, address4)
		}
	})

	t.Run("with address filter when no address matches", func(t *testing.T) {
		addr := testutils.NewAddress()
		_, err := ethKeyStore.GetRoundRobinAddress(testutils.FixtureChainID, []common.Address{addr}...)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("no sending keys available for chain %s that match whitelist: [%s]", testutils.FixtureChainID.String(), addr.Hex()), err.Error())
	})
}

func Test_EthKeyStore_SignTx(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	ethKeyStore := keyStore.Eth()

	k, _ := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	chainID := big.NewInt(evmclient.NullClientChainID)
	tx := types.NewTransaction(0, testutils.NewAddress(), big.NewInt(53), 21000, big.NewInt(1000000000), []byte{1, 2, 3, 4})

	randomAddress := testutils.NewAddress()
	_, err := ethKeyStore.SignTx(randomAddress, tx, chainID)
	require.EqualError(t, err, fmt.Sprintf("unable to find eth key with id %s", randomAddress.String()))

	signed, err := ethKeyStore.SignTx(k.Address, tx, chainID)
	require.NoError(t, err)

	require.NotEqual(t, tx, signed)
}

func Test_EthKeyStore_E2E(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	err := keyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	ks := keyStore.Eth()
	reset := func() {
		keyStore.ResetXXXTestOnly()
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM evm_key_states")))
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
		ks.XXXTestingOnlyAdd(newKey)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		assert.Equal(t, 1, len(keys))
		_, err = ks.Delete(newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		assert.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		assert.Error(t, err)
		_, err = ks.Delete(newKey.ID())
		assert.Error(t, err)
	})

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		exportedKey := `{"address":"0dd359b4f22a30e44b2fd744b679971941865820","crypto":{"cipher":"aes-128-ctr","ciphertext":"b30af964a3b3f37894e599446b4cf2314bbfcd1062e6b35b620d3d20bd9965cc","cipherparams":{"iv":"58a8d75629cc1945da7cf8c24520d1dc"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c352887e9d427d8a6a1869082619b73fac4566082a99f6e367d126f11b434f28"},"mac":"fd76a588210e0bf73d01332091e0e83a4584ee2df31eaec0e27f9a1b94f024b4"},"id":"a5ee0802-1d7b-45b6-aeb8-ea8a3351e715","version":3}`
		importedKey, err := ks.Import([]byte(exportedKey), "p4SsW0rD1!@#_", &cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Equal(t, "0x0dd359b4f22a30E44b2fD744B679971941865820", importedKey.ID())

		k, err := ks.Import([]byte(exportedKey), cltest.Password, &cltest.FixtureChainID)

		assert.Empty(t, k)
		assert.Error(t, err)
	})

	t.Run("fails to export a non-existent key", func(t *testing.T) {
		k, err := ks.Export("non-existent", cltest.Password)

		assert.Empty(t, k)
		assert.Error(t, err)
	})

	t.Run("getting keys states", func(t *testing.T) {
		defer reset()

		t.Run("returns states for keys", func(t *testing.T) {
			k1, err := ethkey.NewV2()
			require.NoError(t, err)
			k2, err := ethkey.NewV2()
			require.NoError(t, err)
			ks.XXXTestingOnlyAdd(k1)
			ks.XXXTestingOnlyAdd(k2)
			require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))

			states, err := ks.GetStatesForKeys([]ethkey.KeyV2{k1, k2})
			require.NoError(t, err)
			assert.Len(t, states, 1)

			chainStates, err := ks.GetStatesForChain(testutils.FixtureChainID)
			require.NoError(t, err)
			assert.Len(t, chainStates, 2) // one created here, one created above

			chainStates, err = ks.GetStatesForChain(testutils.SimulatedChainID)
			require.NoError(t, err)
			assert.Len(t, chainStates, 0)
		})
	})
}

func Test_EthKeyStore_SubscribeToKeyChanges(t *testing.T) {
	t.Parallel()

	chDone := make(chan struct{})
	defer func() { close(chDone) }()
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()
	chSub, unsubscribe := ks.SubscribeToKeyChanges()
	defer unsubscribe()

	var count atomic.Int32

	assertCountAtLeast := func(expected int32) {
		require.Eventually(
			t,
			func() bool { return count.Load() >= expected },
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

	drainAndReset := func() {
		for len(chSub) > 0 {
			<-chSub
		}
		count.Store(0)
	}

	err := ks.EnsureKeys(&cltest.FixtureChainID)
	require.NoError(t, err)
	assertCountAtLeast(1)

	drainAndReset()

	// Create the key includes a state, triggering notify
	k1, err := ks.Create(testutils.FixtureChainID)
	require.NoError(t, err)
	assertCountAtLeast(1)

	drainAndReset()

	// Enabling the key for a new state triggers the notification callback again
	require.NoError(t, ks.Enable(k1.Address, testutils.SimulatedChainID))
	assertCountAtLeast(1)

	drainAndReset()

	// Disabling triggers a notify
	require.NoError(t, ks.Disable(k1.Address, testutils.SimulatedChainID))
	assertCountAtLeast(1)
}

func Test_EthKeyStore_EnsureKeys(t *testing.T) {
	t.Parallel()

	t.Run("creates one unique key per chain if none exist", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		keyStore := cltest.NewKeyStore(t, db, cfg)
		ks := keyStore.Eth()

		testutils.AssertCount(t, db, "evm_key_states", 0)
		err := ks.EnsureKeys(testutils.FixtureChainID, testutils.SimulatedChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 2)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("does nothing if a key exists for a chain", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		keyStore := cltest.NewKeyStore(t, db, cfg)
		ks := keyStore.Eth()

		// Add one enabled key
		_, err := ks.Create(testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 1)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 1)

		// this adds one more key for the additional chain
		err = ks.EnsureKeys(testutils.FixtureChainID, testutils.SimulatedChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 2)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("does nothing if a key exists but is disabled for a chain", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		keyStore := cltest.NewKeyStore(t, db, cfg)
		ks := keyStore.Eth()

		// Add one enabled key
		k, err := ks.Create(testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 1)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 1)

		// disable the key
		err = ks.Disable(k.Address, testutils.FixtureChainID)
		require.NoError(t, err)

		// this does nothing
		err = ks.EnsureKeys(testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm_key_states", 1)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		assert.Len(t, keys, 1)
		state, err := ks.GetState(k.Address.Hex(), testutils.FixtureChainID)
		require.NoError(t, err)
		assert.True(t, state.Disabled)
	})
}

func Test_EthKeyStore_Reset(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()

	k1, addr1 := cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)
	cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)
	cltest.MustInsertRandomKey(t, ks, testutils.SimulatedChainID)

	newNonce := testutils.NewRandomPositiveInt64()

	t.Run("when no state matches address/chain ID", func(t *testing.T) {
		addr := utils.RandomAddress()
		cid := testutils.NewRandomEVMChainID()
		err := ks.Reset(addr, cid, newNonce)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("key state not found with address %s and chainID %s", addr.Hex(), cid.String()))
	})
	t.Run("when no state matches address", func(t *testing.T) {
		addr := utils.RandomAddress()
		err := ks.Reset(addr, testutils.FixtureChainID, newNonce)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("key state not found with address %s and chainID 0", addr.Hex()))
	})
	t.Run("when no state matches chain ID", func(t *testing.T) {
		cid := testutils.NewRandomEVMChainID()
		err := ks.Reset(addr1, cid, newNonce)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("key state not found with address %s and chainID %s", addr1.Hex(), cid.String()))
	})
	t.Run("resets key with given address and chain ID to the given nonce", func(t *testing.T) {
		err := ks.Reset(k1.Address, testutils.FixtureChainID, newNonce)
		assert.NoError(t, err)

		nonce, err := ks.NextSequence(k1.Address, testutils.FixtureChainID)
		require.NoError(t, err)

		assert.Equal(t, nonce.Int64(), newNonce)

		state, err := ks.GetState(k1.Address.Hex(), testutils.FixtureChainID)
		require.NoError(t, err)
		assert.Equal(t, nonce.Int64(), state.NextNonce)

		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 3)
		states, err := ks.GetStatesForKeys(keys)
		require.NoError(t, err)
		require.Len(t, states, 3)
		for _, state = range states {
			if state.Address.Address() == k1.Address {
				assert.Equal(t, nonce.Int64(), state.NextNonce)
			} else {
				// the other states didn't get updated
				assert.Equal(t, int64(0), state.NextNonce)
			}
		}
	})
}

func Test_NextSequence(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()
	randNonce := testutils.NewRandomPositiveInt64()

	_, addr1 := cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID, randNonce)
	cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)

	nonce, err := ks.NextSequence(addr1, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, randNonce, nonce.Int64())

	_, err = ks.NextSequence(addr1, testutils.SimulatedChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("NextSequence failed: key with address %s is not enabled for chain %s: sql: no rows in result set", addr1.Hex(), testutils.SimulatedChainID.String()))

	randAddr1 := utils.RandomAddress()
	_, err = ks.NextSequence(randAddr1, testutils.FixtureChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("key with address %s does not exist", randAddr1.Hex()))

	randAddr2 := utils.RandomAddress()
	_, err = ks.NextSequence(randAddr2, testutils.NewRandomEVMChainID())
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("key with address %s does not exist", randAddr2.Hex()))
}

func Test_IncrementNextSequence(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()
	randNonce := testutils.NewRandomPositiveInt64()

	_, addr1 := cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID, randNonce)
	evmAddr1 := addr1
	cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)

	err := ks.IncrementNextSequence(evmAddr1, testutils.FixtureChainID, evmtypes.Nonce(randNonce-1))
	assert.ErrorIs(t, err, sql.ErrNoRows)

	err = ks.IncrementNextSequence(evmAddr1, testutils.FixtureChainID, evmtypes.Nonce(randNonce))
	require.NoError(t, err)
	var nonce int64
	require.NoError(t, db.Get(&nonce, `SELECT next_nonce FROM evm_key_states WHERE address = $1 AND evm_chain_id = $2`, addr1, testutils.FixtureChainID.String()))
	assert.Equal(t, randNonce+1, nonce)

	err = ks.IncrementNextSequence(evmAddr1, testutils.SimulatedChainID, evmtypes.Nonce(randNonce+1))
	assert.ErrorIs(t, err, sql.ErrNoRows)

	randAddr1 := utils.RandomAddress()
	err = ks.IncrementNextSequence(randAddr1, testutils.FixtureChainID, evmtypes.Nonce(randNonce+1))
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("key with address %s does not exist", randAddr1.Hex()))

	randAddr2 := utils.RandomAddress()
	err = ks.IncrementNextSequence(randAddr2, testutils.NewRandomEVMChainID(), evmtypes.Nonce(randNonce+1))
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("key with address %s does not exist", randAddr2.Hex()))

	// verify it didnt get changed by any erroring calls
	require.NoError(t, db.Get(&nonce, `SELECT next_nonce FROM evm_key_states WHERE address = $1 AND evm_chain_id = $2`, addr1, testutils.FixtureChainID.String()))
	assert.Equal(t, randNonce+1, nonce)
}

func Test_EthKeyStore_Delete(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()

	randKeyID := utils.RandomAddress().Hex()
	_, err := ks.Delete(randKeyID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("unable to find eth key with id %s", randKeyID))

	_, addr1 := cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)
	_, addr2 := cltest.MustInsertRandomKey(t, ks, testutils.FixtureChainID)
	cltest.MustInsertRandomKey(t, ks, testutils.SimulatedChainID)
	require.NoError(t, ks.Enable(addr1, testutils.SimulatedChainID))

	testutils.AssertCount(t, db, "evm_key_states", 4)
	keys, err := ks.GetAll()
	require.NoError(t, err)
	assert.Len(t, keys, 3)
	_, err = ks.GetState(addr1.Hex(), testutils.FixtureChainID)
	require.NoError(t, err)
	_, err = ks.GetState(addr1.Hex(), testutils.SimulatedChainID)
	require.NoError(t, err)

	deletedK, err := ks.Delete(addr1.String())
	require.NoError(t, err)
	assert.Equal(t, addr1, deletedK.Address)

	testutils.AssertCount(t, db, "evm_key_states", 2)
	keys, err = ks.GetAll()
	require.NoError(t, err)
	assert.Len(t, keys, 2)
	_, err = ks.GetState(addr1.Hex(), testutils.FixtureChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("state not found for eth key ID %s", addr1.Hex()))
	_, err = ks.GetState(addr1.Hex(), testutils.SimulatedChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("state not found for eth key ID %s", addr1.Hex()))
	_, err = ks.GetState(addr2.Hex(), testutils.FixtureChainID)
	require.NoError(t, err)
}

func Test_EthKeyStore_CheckEnabled(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ks := keyStore.Eth()

	// create keys
	// - key 1
	//   enabled - fixture
	//   enabled - simulated
	// - key 2
	//   enabled - fixture
	//   disabled - simulated
	// - key 3
	//   enabled - simulated
	// - key 4
	//   enabled - fixture
	k1, addr1 := cltest.MustInsertRandomKey(t, ks, []utils.Big{})
	require.NoError(t, ks.Enable(k1.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))

	k2, addr2 := cltest.MustInsertRandomKey(t, ks, []utils.Big{})
	require.NoError(t, ks.Enable(k2.Address, testutils.FixtureChainID))
	require.NoError(t, ks.Enable(k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Disable(k2.Address, testutils.SimulatedChainID))

	k3, addr3 := cltest.MustInsertRandomKey(t, ks, []utils.Big{})
	require.NoError(t, ks.Enable(k3.Address, testutils.SimulatedChainID))

	t.Run("enabling the same key multiple times does not create duplicate states", func(t *testing.T) {
		require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(k1.Address, testutils.FixtureChainID))

		states, err := ks.GetStatesForKeys([]ethkey.KeyV2{k1})
		require.NoError(t, err)
		assert.Len(t, states, 2)
		var cids []*big.Int
		for i := range states {
			cid := states[i].EVMChainID.ToInt()
			cids = append(cids, cid)
		}
		assert.Contains(t, cids, testutils.FixtureChainID)
		assert.Contains(t, cids, testutils.SimulatedChainID)

		for _, s := range states {
			assert.Equal(t, addr1, s.Address.Address())
		}
	})

	t.Run("returns nil when key is enabled for given chain", func(t *testing.T) {
		err := ks.CheckEnabled(addr1, testutils.FixtureChainID)
		assert.NoError(t, err)
		err = ks.CheckEnabled(addr1, testutils.SimulatedChainID)
		assert.NoError(t, err)
	})

	t.Run("returns error when key does not exist", func(t *testing.T) {
		addr := utils.RandomAddress()
		err := ks.CheckEnabled(addr, testutils.FixtureChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", addr.Hex()))
	})

	t.Run("returns error when key exists but has never been enabled (no state) for the given chain", func(t *testing.T) {
		err := ks.CheckEnabled(addr3, testutils.FixtureChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("eth key with address %s exists but is has not been enabled for chain 0 (enabled only for chain IDs: 1337)", addr3.Hex()))
	})

	t.Run("returns error when key exists but is disabled for the given chain", func(t *testing.T) {
		err := ks.CheckEnabled(addr2, testutils.SimulatedChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("eth key with address %s exists but is disabled for chain 1337 (enabled only for chain IDs: 0)", addr2.Hex()))
	})
}
