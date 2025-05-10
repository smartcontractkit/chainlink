package utils

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

var (
	LinkTokenEth  = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	linkTokenSol  = "y9MdSjD9Beg9EFaeQGdMpESFWLNdSfZKQKeYLBfmnjJ"
	usdtETH       = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	v1            = semver.MustParse("1.0.0")
	v2            = semver.MustParse("2.1.0")
	chainEthereum = chainselectors.ETHEREUM_MAINNET.Selector
	chainSolana   = chainselectors.SOLANA_MAINNET.Selector
)

func TestAddressBookToDataStore(t *testing.T) {
	t.Parallel()
	ab := cldf.NewMemoryAddressBook()

	// Add addresses to the ethereum
	err := ab.Save(chainEthereum, LinkTokenEth, cldf.TypeAndVersion{
		Type:    "LinkToken",
		Version: *v1,
		Labels:  cldf.NewLabelSet("label1", "label2"),
	})
	require.NoError(t, err)

	err = ab.Save(chainEthereum, usdtETH, cldf.TypeAndVersion{
		Type:    "USDT",
		Version: *v2,
		Labels:  cldf.NewLabelSet("label3"),
	})
	require.NoError(t, err)

	// Add address to the solana
	err = ab.Save(chainSolana, linkTokenSol, cldf.TypeAndVersion{
		Type:    "LinkToken",
		Version: *v1,
		Labels:  cldf.NewLabelSet("testnet"),
	})
	require.NoError(t, err)

	// Convert AddressBook to DataStore
	ds, err := AddressBookToDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata](ab)
	require.NoError(t, err)

	t.Run("verifies expected address count", func(t *testing.T) {
		addressRefs, err := ds.Addresses().Fetch()
		require.NoError(t, err)
		assert.Len(t, addressRefs, 3, "DataStore should contain 3 addresses")
	})

	t.Run("verifies ethereum address 1 conversion", func(t *testing.T) {
		record, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainEthereum, "LinkToken", v1, ""))
		require.NoError(t, err)
		assert.Equal(t, LinkTokenEth, record.Address)
		assert.True(t, record.Labels.Contains("label1"))
		assert.True(t, record.Labels.Contains("label2"))
	})

	t.Run("verifies ethereum address 2 conversion", func(t *testing.T) {
		record, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainEthereum, "USDT", v2, ""))
		require.NoError(t, err)
		assert.Equal(t, usdtETH, record.Address)
		assert.True(t, record.Labels.Contains("label3"))
	})

	t.Run("verifies solana address 1 conversion", func(t *testing.T) {
		record, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainSolana, "LinkToken", v1, ""))
		require.NoError(t, err)
		assert.Equal(t, linkTokenSol, record.Address)
		assert.True(t, record.Labels.Contains("testnet"))
	})
}

func TestDataStoreToAddressBook(t *testing.T) {
	t.Parallel()

	ds := datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()

	// Add to ethereum
	err := ds.Addresses().Upsert(datastore.AddressRef{
		Address:       LinkTokenEth,
		ChainSelector: chainEthereum,
		Type:          "LinkToken",
		Version:       v1,
		Labels:        datastore.NewLabelSet("mainnet", "stable"),
	})
	require.NoError(t, err)

	// Add to solana
	err = ds.Addresses().Upsert(datastore.AddressRef{
		Address:       linkTokenSol,
		ChainSelector: chainSolana,
		Type:          "LinkToken",
		Version:       v2,
		Labels:        datastore.NewLabelSet("mainnet"),
	})
	require.NoError(t, err)

	// Convert DataStore to AddressBook
	ab, err := DataStoreToAddressBook[datastore.DefaultMetadata, datastore.DefaultMetadata](ds.Seal())
	require.NoError(t, err)

	t.Run("verifies expected chain count", func(t *testing.T) {
		addresses, err := ab.Addresses()
		require.NoError(t, err)
		assert.Len(t, addresses, 2, "AddressBook should have 2 chains")
	})

	t.Run("verifies chain 1 conversion", func(t *testing.T) {
		t.Run("chain has expected address count", func(t *testing.T) {
			chain1Addresses, err := ab.AddressesForChain(chainEthereum)
			require.NoError(t, err)
			assert.Len(t, chain1Addresses, 1, "Ethereum address map should have 1 address")
		})

		t.Run("LinkToken address exists with correct properties", func(t *testing.T) {
			ethereumAddresses, err := ab.AddressesForChain(chainEthereum)
			require.NoError(t, err)

			tv1, exists := ethereumAddresses[LinkTokenEth]
			assert.True(t, exists, "Address should exist in ethereum address map")
			assert.Equal(t, cldf.ContractType("LinkToken"), tv1.Type)
			assert.Equal(t, v1.String(), tv1.Version.String())
			assert.True(t, tv1.Labels.Contains("mainnet"))
			assert.True(t, tv1.Labels.Contains("stable"))
		})
	})

	t.Run("verifies solana chain conversion", func(t *testing.T) {
		t.Run("chain has expected address count", func(t *testing.T) {
			solanaAddresses, err := ab.AddressesForChain(chainSolana)
			require.NoError(t, err)
			assert.Len(t, solanaAddresses, 1, "Solana address map should have 1 address")
		})

		t.Run("LinkToken address exists with correct properties", func(t *testing.T) {
			chain2Addresses, err := ab.AddressesForChain(chainSolana)
			require.NoError(t, err)

			tv2, exists := chain2Addresses[linkTokenSol]
			assert.True(t, exists, "Address should exist in solana address map")
			assert.Equal(t, cldf.ContractType("LinkToken"), tv2.Type)
			assert.Equal(t, v2.String(), tv2.Version.String())
			assert.True(t, tv2.Labels.Contains("mainnet"))
		})
	})
}

func TestAddressBookToNewDataStore(t *testing.T) {
	t.Parallel()

	ab := cldf.NewMemoryAddressBook()

	err := ab.Save(chainEthereum, LinkTokenEth, cldf.TypeAndVersion{
		Type:    "LinkToken",
		Version: *v1,
		Labels:  cldf.NewLabelSet("testLabel"),
	})
	require.NoError(t, err)

	// Convert to mutable DataStore
	ds, err := AddressBookToNewDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata](ab)
	require.NoError(t, err)

	t.Run("verifies initial conversion", func(t *testing.T) {
		t.Run("has expected address count", func(t *testing.T) {
			addressRefs, err := ds.Addresses().Fetch()
			require.NoError(t, err)
			assert.Len(t, addressRefs, 1, "Should have 1 address after conversion")
		})

		t.Run("has correct address properties", func(t *testing.T) {
			record, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainEthereum, "LinkToken", v1, ""))
			require.NoError(t, err)
			assert.Equal(t, LinkTokenEth, record.Address)
			assert.True(t, record.Labels.Contains("testLabel"))
		})
	})

	t.Run("verifies mutability", func(t *testing.T) {
		// Add a new address to verify mutability
		err = ds.Addresses().Upsert(datastore.AddressRef{
			Address:       usdtETH,
			ChainSelector: chainEthereum,
			Type:          "USDT",
			Version:       v2,
		})
		require.NoError(t, err)

		t.Run("has increased address count after addition", func(t *testing.T) {
			addressRefs, err := ds.Addresses().Fetch()
			require.NoError(t, err)
			assert.Len(t, addressRefs, 2, "Should have 2 addresses after adding one to mutable store")
		})

		t.Run("new address has correct properties", func(t *testing.T) {
			newRecord, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainEthereum, "USDT", v2, ""))
			require.NoError(t, err)
			assert.Equal(t, usdtETH, newRecord.Address)
		})
	})
}
