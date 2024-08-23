package keystore_test

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

func Test_EthKeyStore(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	err := keyStore.Unlock(testutils.Context(t), cltest.Password)
	require.NoError(t, err)
	ethKeyStore := keyStore.Eth()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
		keyStore.ResetXXXTestOnly()
		require.NoError(t, commonutils.JustError(db.ExecContext(ctx, "DELETE FROM encrypted_key_rings")))
		require.NoError(t, commonutils.JustError(db.ExecContext(ctx, "DELETE FROM evm.key_states")))
		require.NoError(t, keyStore.Unlock(ctx, cltest.Password))
	}
	const statesTableName = "evm.key_states"

	t.Run("Create / GetAll / Get", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ethKeyStore.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(retrievedKeys))
		require.Equal(t, key.Address, retrievedKeys[0].Address)
		foundKey, err := ethKeyStore.Get(ctx, key.Address.Hex())
		require.NoError(t, err)
		require.Equal(t, key, foundKey)
		// adds ethkey.State
		cltest.AssertCount(t, db, statesTableName, 1)
		var state ethkey.State
		sql := fmt.Sprintf(`SELECT address, disabled, evm_chain_id, created_at, updated_at from %s LIMIT 1`, statesTableName)
		require.NoError(t, db.GetContext(ctx, &state, sql))
		require.Equal(t, state.Address.Address(), retrievedKeys[0].Address)
		// adds key to db
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(ctx, cltest.Password))
		retrievedKeys, err = ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(retrievedKeys))
		require.Equal(t, key.Address, retrievedKeys[0].Address)
		// adds 2nd key
		_, err = ethKeyStore.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKeys, err = ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(retrievedKeys))
	})

	t.Run("GetAll ordering", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		var keys []ethkey.KeyV2
		for i := 0; i < 5; i++ {
			key, err := ethKeyStore.Create(ctx, &cltest.FixtureChainID)
			require.NoError(t, err)
			keys = append(keys, key)
		}
		retrievedKeys, err := ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 5, len(retrievedKeys))

		sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })

		assert.Equal(t, keys, retrievedKeys)
	})

	t.Run("RemoveKey", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ethKeyStore.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		_, err = ethKeyStore.Delete(ctx, key.ID())
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(retrievedKeys))
		cltest.AssertCount(t, db, statesTableName, 0)
	})

	t.Run("Delete removes key even if evm.txes are present", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ethKeyStore.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		// ensure at least one state is present
		cltest.AssertCount(t, db, statesTableName, 1)

		// add one eth_tx
		txStore := cltest.NewTestTxStore(t, db)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 42, key.Address)

		_, err = ethKeyStore.Delete(ctx, key.ID())
		require.NoError(t, err)
		retrievedKeys, err := ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(retrievedKeys))
		cltest.AssertCount(t, db, statesTableName, 0)
	})

	t.Run("EnsureKeys / EnabledKeysForChain", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		err := ethKeyStore.EnsureKeys(ctx, &cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys1, err := ethKeyStore.EnabledKeysForChain(ctx, testutils.FixtureChainID)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys1))
		cltest.AssertCount(t, db, statesTableName, 1)

		err = ethKeyStore.EnsureKeys(ctx, &cltest.FixtureChainID)
		assert.NoError(t, err)
		sendingKeys2, err := ethKeyStore.EnabledKeysForChain(ctx, testutils.FixtureChainID)
		assert.NoError(t, err)

		require.Equal(t, 1, len(sendingKeys2))
		require.Equal(t, sendingKeys1, sendingKeys2)
	})

	t.Run("EnabledKeysForChain with specified chain ID", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ethKeyStore.Create(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		key2, err := ethKeyStore.Create(ctx, big.NewInt(1337))
		require.NoError(t, err)

		keys, err := ethKeyStore.EnabledKeysForChain(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key, keys[0])

		keys, err = ethKeyStore.EnabledKeysForChain(ctx, big.NewInt(1337))
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, key2, keys[0])

		_, err = ethKeyStore.EnabledKeysForChain(ctx, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "chainID must be non-nil")
	})

	t.Run("EnabledAddressesForChain with specified chain ID", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ethKeyStore.Create(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		key2, err := ethKeyStore.Create(ctx, big.NewInt(1337))
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 2)
		keys, err := ethKeyStore.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 2)

		//get enabled addresses for FixtureChainID
		enabledAddresses, err := ethKeyStore.EnabledAddressesForChain(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, enabledAddresses, 1)
		require.Equal(t, key.Address, enabledAddresses[0])

		//get enabled addresses for chain 1337
		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(ctx, big.NewInt(1337))
		require.NoError(t, err)
		require.Len(t, enabledAddresses, 1)
		require.Equal(t, key2.Address, enabledAddresses[0])

		// /get enabled addresses for nil chain ID
		_, err = ethKeyStore.EnabledAddressesForChain(ctx, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "chainID must be non-nil")

		// disable the key for chain FixtureChainID
		err = ethKeyStore.Disable(ctx, key.Address, testutils.FixtureChainID)
		require.NoError(t, err)

		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, enabledAddresses, 0)
		enabledAddresses, err = ethKeyStore.EnabledAddressesForChain(ctx, big.NewInt(1337))
		require.NoError(t, err)
		assert.Len(t, enabledAddresses, 1)
		require.Equal(t, key2.Address, enabledAddresses[0])
	})
}

func Test_EthKeyStore_GetRoundRobinAddress(t *testing.T) {
	ctx := testutils.Context(t)
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	t.Run("should error when no addresses", func(t *testing.T) {
		ctx1 := testutils.Context(t)
		_, err := ethKeyStore.GetRoundRobinAddress(ctx1, testutils.FixtureChainID)
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
	k1, _ := cltest.MustInsertRandomKeyNoChains(t, ethKeyStore)
	require.NoError(t, ethKeyStore.Add(ctx, k1.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Add(ctx, k1.Address, testutils.SimulatedChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k1.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k1.Address, testutils.SimulatedChainID))

	k2, _ := cltest.MustInsertRandomKeyNoChains(t, ethKeyStore)
	require.NoError(t, ethKeyStore.Add(ctx, k2.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Add(ctx, k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k2.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ethKeyStore.Disable(ctx, k2.Address, testutils.SimulatedChainID))

	k3, _ := cltest.MustInsertRandomKeyNoChains(t, ethKeyStore)
	require.NoError(t, ethKeyStore.Add(ctx, k3.Address, testutils.SimulatedChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k3.Address, testutils.SimulatedChainID))

	k4, _ := cltest.MustInsertRandomKeyNoChains(t, ethKeyStore)
	require.NoError(t, ethKeyStore.Add(ctx, k4.Address, testutils.FixtureChainID))
	require.NoError(t, ethKeyStore.Enable(ctx, k4.Address, testutils.FixtureChainID))

	t.Run("with no address filter, rotates between all enabled addresses", func(t *testing.T) {
		address1, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		address2, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		address3, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		address4, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		address5, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		address6, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID)
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

			address1, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address2, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address3, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID, addresses...)
			require.NoError(t, err)
			address4, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID, addresses...)
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

			address1, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address2, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address3, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.SimulatedChainID, addresses...)
			require.NoError(t, err)
			address4, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.SimulatedChainID, addresses...)
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
		_, err := ethKeyStore.GetRoundRobinAddress(ctx, testutils.FixtureChainID, []common.Address{addr}...)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("no sending keys available for chain %s that match whitelist: [%s]", testutils.FixtureChainID.String(), addr.Hex()), err.Error())
	})
}

func Test_EthKeyStore_SignTx(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	k, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	chainID := big.NewInt(evmclient.NullClientChainID)
	tx := cltest.NewLegacyTransaction(0, testutils.NewAddress(), big.NewInt(53), 21000, big.NewInt(1000000000), []byte{1, 2, 3, 4})

	randomAddress := testutils.NewAddress()
	_, err := ethKeyStore.SignTx(ctx, randomAddress, tx, chainID)
	require.EqualError(t, err, "Key not found")

	signed, err := ethKeyStore.SignTx(ctx, k.Address, tx, chainID)
	require.NoError(t, err)

	require.NotEqual(t, tx, signed)
}

func Test_EthKeyStore_E2E(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	err := keyStore.Unlock(testutils.Context(t), cltest.Password)
	require.NoError(t, err)
	ks := keyStore.Eth()
	reset := func() {
		ctx := testutils.Context(t)
		keyStore.ResetXXXTestOnly()
		require.NoError(t, commonutils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		require.NoError(t, commonutils.JustError(db.Exec("DELETE FROM evm.key_states")))
		require.NoError(t, keyStore.Unlock(ctx, cltest.Password))
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		keys, err := ks.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("errors when getting non-existent ID", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		_, err := ks.Get(ctx, "non-existent-id")
		require.Error(t, err)
	})

	t.Run("creates a key", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ks.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		retrievedKey, err := ks.Get(ctx, key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		key, err := ks.Create(ctx, &cltest.FixtureChainID)
		require.NoError(t, err)
		exportJSON, err := ks.Export(ctx, key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Delete(ctx, key.ID())
		require.NoError(t, err)
		_, err = ks.Get(ctx, key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(ctx, exportJSON, cltest.Password, &cltest.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(ctx, key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		ctx := testutils.Context(t)
		defer reset()
		newKey, err := ethkey.NewV2()
		require.NoError(t, err)
		ks.XXXTestingOnlyAdd(ctx, newKey)
		keys, err := ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(keys))
		_, err = ks.Delete(ctx, newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, len(keys))
		_, err = ks.Get(ctx, newKey.ID())
		assert.Error(t, err)
		_, err = ks.Delete(ctx, newKey.ID())
		assert.Error(t, err)
	})

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		ctx := testutils.Context(t)
		exportedKey := `{"address":"0dd359b4f22a30e44b2fd744b679971941865820","crypto":{"cipher":"aes-128-ctr","ciphertext":"b30af964a3b3f37894e599446b4cf2314bbfcd1062e6b35b620d3d20bd9965cc","cipherparams":{"iv":"58a8d75629cc1945da7cf8c24520d1dc"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c352887e9d427d8a6a1869082619b73fac4566082a99f6e367d126f11b434f28"},"mac":"fd76a588210e0bf73d01332091e0e83a4584ee2df31eaec0e27f9a1b94f024b4"},"id":"a5ee0802-1d7b-45b6-aeb8-ea8a3351e715","version":3}`
		importedKey, err := ks.Import(ctx, []byte(exportedKey), "p4SsW0rD1!@#_", &cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Equal(t, "0x0dd359b4f22a30E44b2fD744B679971941865820", importedKey.ID())

		k, err := ks.Import(ctx, []byte(exportedKey), cltest.Password, &cltest.FixtureChainID)

		assert.Empty(t, k)
		assert.Error(t, err)
	})

	t.Run("fails to export a non-existent key", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, err := ks.Export(ctx, "non-existent", cltest.Password)

		assert.Empty(t, k)
		assert.Error(t, err)
	})

	t.Run("getting keys states", func(t *testing.T) {
		defer reset()

		t.Run("returns states for keys", func(t *testing.T) {
			ctx := testutils.Context(t)
			k1, err := ethkey.NewV2()
			require.NoError(t, err)
			k2, err := ethkey.NewV2()
			require.NoError(t, err)
			ks.XXXTestingOnlyAdd(ctx, k1)
			ks.XXXTestingOnlyAdd(ctx, k2)
			require.NoError(t, ks.Add(ctx, k1.Address, testutils.FixtureChainID))
			require.NoError(t, ks.Enable(ctx, k1.Address, testutils.FixtureChainID))

			states, err := ks.GetStatesForKeys(ctx, []ethkey.KeyV2{k1, k2})
			require.NoError(t, err)
			assert.Len(t, states, 1)

			chainStates, err := ks.GetStatesForChain(ctx, testutils.FixtureChainID)
			require.NoError(t, err)
			assert.Len(t, chainStates, 2) // one created here, one created above

			chainStates, err = ks.GetStatesForChain(ctx, testutils.SimulatedChainID)
			require.NoError(t, err)
			assert.Len(t, chainStates, 0)
		})
	})
}

func Test_EthKeyStore_SubscribeToKeyChanges(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	chDone := make(chan struct{})
	defer func() { close(chDone) }()
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ks := keyStore.Eth()
	chSub, unsubscribe := ks.SubscribeToKeyChanges(ctx)
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

	err := ks.EnsureKeys(ctx, &cltest.FixtureChainID)
	require.NoError(t, err)
	assertCountAtLeast(1)

	drainAndReset()

	// Create the key includes a state, triggering notify
	k1, err := ks.Create(ctx, testutils.FixtureChainID)
	require.NoError(t, err)
	assertCountAtLeast(1)

	drainAndReset()

	// Enabling the key for a new state triggers the notification callback again
	require.NoError(t, ks.Add(ctx, k1.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(ctx, k1.Address, testutils.SimulatedChainID))
	assertCountAtLeast(1)

	drainAndReset()

	// Disabling triggers a notify
	require.NoError(t, ks.Disable(ctx, k1.Address, testutils.SimulatedChainID))
	assertCountAtLeast(1)
}

func Test_EthKeyStore_Enable(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ks := keyStore.Eth()

	t.Run("already existing disabled key gets enabled", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, _ := cltest.MustInsertRandomKeyNoChains(t, ks)
		require.NoError(t, ks.Add(ctx, k.Address, testutils.SimulatedChainID))
		require.NoError(t, ks.Disable(ctx, k.Address, testutils.SimulatedChainID))
		require.NoError(t, ks.Enable(ctx, k.Address, testutils.SimulatedChainID))
		key, err := ks.GetState(ctx, k.Address.Hex(), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Equal(t, key.Disabled, false)
	})

	t.Run("creates key, deletes it unsafely and then enable creates it again", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, _ := cltest.MustInsertRandomKeyNoChains(t, ks)
		require.NoError(t, ks.Add(ctx, k.Address, testutils.SimulatedChainID))
		_, err := db.Exec("DELETE FROM evm.key_states WHERE address = $1", k.Address)
		require.NoError(t, err)
		require.NoError(t, ks.Enable(ctx, k.Address, testutils.SimulatedChainID))
		key, err := ks.GetState(ctx, k.Address.Hex(), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Equal(t, key.Disabled, false)
	})

	t.Run("creates key and enables it if it exists in the keystore, but is missing from key states db table", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, _ := cltest.MustInsertRandomKeyNoChains(t, ks)
		require.NoError(t, ks.Enable(ctx, k.Address, testutils.SimulatedChainID))
		key, err := ks.GetState(ctx, k.Address.Hex(), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Equal(t, key.Disabled, false)
	})

	t.Run("errors if key is not present in keystore", func(t *testing.T) {
		ctx := testutils.Context(t)
		addrNotInKs := testutils.NewAddress()
		require.Error(t, ks.Enable(ctx, addrNotInKs, testutils.SimulatedChainID))
		_, err := ks.GetState(ctx, addrNotInKs.Hex(), testutils.SimulatedChainID)
		require.Error(t, err)
	})
}

func Test_EthKeyStore_EnsureKeys(t *testing.T) {
	t.Parallel()

	t.Run("creates one unique key per chain if none exist", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		keyStore := cltest.NewKeyStore(t, db)
		ks := keyStore.Eth()

		testutils.AssertCount(t, db, "evm.key_states", 0)
		err := ks.EnsureKeys(ctx, testutils.FixtureChainID, testutils.SimulatedChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 2)
		keys, err := ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("does nothing if a key exists for a chain", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		keyStore := cltest.NewKeyStore(t, db)
		ks := keyStore.Eth()

		// Add one enabled key
		_, err := ks.Create(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 1)
		keys, err := ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 1)

		// this adds one more key for the additional chain
		err = ks.EnsureKeys(ctx, testutils.FixtureChainID, testutils.SimulatedChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 2)
		keys, err = ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("does nothing if a key exists but is disabled for a chain", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		keyStore := cltest.NewKeyStore(t, db)
		ks := keyStore.Eth()

		// Add one enabled key
		k, err := ks.Create(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 1)
		keys, err := ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 1)

		// disable the key
		err = ks.Disable(ctx, k.Address, testutils.FixtureChainID)
		require.NoError(t, err)

		// this does nothing
		err = ks.EnsureKeys(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		testutils.AssertCount(t, db, "evm.key_states", 1)
		keys, err = ks.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, keys, 1)
		state, err := ks.GetState(ctx, k.Address.Hex(), testutils.FixtureChainID)
		require.NoError(t, err)
		assert.True(t, state.Disabled)
	})
}

func Test_EthKeyStore_Delete(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ks := keyStore.Eth()

	randKeyID := utils.RandomAddress().Hex()
	_, err := ks.Delete(ctx, randKeyID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Key not found")

	_, addr1 := cltest.MustInsertRandomKey(t, ks)
	_, addr2 := cltest.MustInsertRandomKey(t, ks)
	cltest.MustInsertRandomKey(t, ks, *ubig.New(testutils.SimulatedChainID))
	require.NoError(t, ks.Add(ctx, addr1, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(ctx, addr1, testutils.SimulatedChainID))

	testutils.AssertCount(t, db, "evm.key_states", 4)
	keys, err := ks.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, keys, 3)
	_, err = ks.GetState(ctx, addr1.Hex(), testutils.FixtureChainID)
	require.NoError(t, err)
	_, err = ks.GetState(ctx, addr1.Hex(), testutils.SimulatedChainID)
	require.NoError(t, err)

	deletedK, err := ks.Delete(ctx, addr1.String())
	require.NoError(t, err)
	assert.Equal(t, addr1, deletedK.Address)

	testutils.AssertCount(t, db, "evm.key_states", 2)
	keys, err = ks.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, keys, 2)
	_, err = ks.GetState(ctx, addr1.Hex(), testutils.FixtureChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("state not found for eth key ID %s", addr1.Hex()))
	_, err = ks.GetState(ctx, addr1.Hex(), testutils.SimulatedChainID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("state not found for eth key ID %s", addr1.Hex()))
	_, err = ks.GetState(ctx, addr2.Hex(), testutils.FixtureChainID)
	require.NoError(t, err)
}

func Test_EthKeyStore_CheckEnabled(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
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
	k1, addr1 := cltest.MustInsertRandomKeyNoChains(t, ks)
	require.NoError(t, ks.Add(ctx, k1.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Add(ctx, k1.Address, testutils.FixtureChainID))
	require.NoError(t, ks.Enable(ctx, k1.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(ctx, k1.Address, testutils.FixtureChainID))

	k2, addr2 := cltest.MustInsertRandomKeyNoChains(t, ks)
	require.NoError(t, ks.Add(ctx, k2.Address, testutils.FixtureChainID))
	require.NoError(t, ks.Add(ctx, k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(ctx, k2.Address, testutils.FixtureChainID))
	require.NoError(t, ks.Enable(ctx, k2.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Disable(ctx, k2.Address, testutils.SimulatedChainID))

	k3, addr3 := cltest.MustInsertRandomKeyNoChains(t, ks)
	require.NoError(t, ks.Add(ctx, k3.Address, testutils.SimulatedChainID))
	require.NoError(t, ks.Enable(ctx, k3.Address, testutils.SimulatedChainID))

	t.Run("enabling the same key multiple times does not create duplicate states", func(t *testing.T) {
		ctx2 := testutils.Context(t)
		require.NoError(t, ks.Enable(ctx2, k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(ctx2, k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(ctx2, k1.Address, testutils.FixtureChainID))
		require.NoError(t, ks.Enable(ctx2, k1.Address, testutils.FixtureChainID))

		states, err := ks.GetStatesForKeys(ctx2, []ethkey.KeyV2{k1})
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
		err := ks.CheckEnabled(ctx, addr1, testutils.FixtureChainID)
		assert.NoError(t, err)
		err = ks.CheckEnabled(ctx, addr1, testutils.SimulatedChainID)
		assert.NoError(t, err)
	})

	t.Run("returns error when key does not exist", func(t *testing.T) {
		addr := utils.RandomAddress()
		err := ks.CheckEnabled(ctx, addr, testutils.FixtureChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", addr.Hex()))
	})

	t.Run("returns error when key exists but has never been enabled (no state) for the given chain", func(t *testing.T) {
		err := ks.CheckEnabled(ctx, addr3, testutils.FixtureChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("eth key with address %s exists but is has not been enabled for chain 0 (enabled only for chain IDs: 1337)", addr3.Hex()))
	})

	t.Run("returns error when key exists but is disabled for the given chain", func(t *testing.T) {
		err := ks.CheckEnabled(ctx, addr2, testutils.SimulatedChainID)
		assert.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("eth key with address %s exists but is disabled for chain 1337 (enabled only for chain IDs: 0)", addr2.Hex()))
	})
}

func Test_EthKeyStore_Disable(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ks := keyStore.Eth()

	t.Run("creates key, deletes it unsafely and then enable creates it again", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, _ := cltest.MustInsertRandomKeyNoChains(t, ks)
		require.NoError(t, ks.Add(ctx, k.Address, testutils.SimulatedChainID))
		_, err := db.Exec("DELETE FROM evm.key_states WHERE address = $1", k.Address)
		require.NoError(t, err)
		require.NoError(t, ks.Disable(ctx, k.Address, testutils.SimulatedChainID))
		key, err := ks.GetState(ctx, k.Address.Hex(), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Equal(t, key.Disabled, true)
	})

	t.Run("creates key and enables it if it exists in the keystore, but is missing from key states db table", func(t *testing.T) {
		ctx := testutils.Context(t)
		k, _ := cltest.MustInsertRandomKeyNoChains(t, ks)
		require.NoError(t, ks.Disable(ctx, k.Address, testutils.SimulatedChainID))
		key, err := ks.GetState(ctx, k.Address.Hex(), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Equal(t, key.Disabled, true)
	})

	t.Run("errors if key is not present in keystore", func(t *testing.T) {
		ctx := testutils.Context(t)
		addrNotInKs := testutils.NewAddress()
		require.Error(t, ks.Disable(ctx, addrNotInKs, testutils.SimulatedChainID))
		_, err := ks.GetState(ctx, addrNotInKs.Hex(), testutils.SimulatedChainID)
		require.Error(t, err)
	})
}
