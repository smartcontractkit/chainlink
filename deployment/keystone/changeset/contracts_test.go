package changeset_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

func TestGetOwnableContractV2(t *testing.T) {
	t.Parallel()
	v1 := semver.MustParse("1.1.0")

	chain := memory.NewMemoryChain(t, chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector)

	t.Run("finds contract when targetAddr is provided", func(t *testing.T) {
		t.Parallel()

		// Create a datastore
		ds := datastore.NewMemoryDataStore()
		targetAddr := testutils.NewAddress()
		targetAddrStr := targetAddr.String()

		// Create an address ref
		addrRef := datastore.AddressRef{
			ChainSelector: chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}

		err := ds.AddressRefStore.Add(addrRef)
		require.NoError(t, err)

		c, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)
		assert.NotNil(t, c)
		contract := *c
		assert.Equal(t, targetAddr, contract.Address())
	})

	t.Run("errors when targetAddr not found in datastore", func(t *testing.T) {
		t.Parallel()

		// Create a datastore
		ds := datastore.NewMemoryDataStore()
		targetAddr := testutils.NewAddress()
		targetAddrStr := targetAddr.String()
		nonExistentAddr := testutils.NewAddress()
		nonExistentAddrStr := nonExistentAddr.String()
		v1 := semver.MustParse("1.1.0")

		// Create an address ref for existing address
		addrRef := datastore.AddressRef{
			ChainSelector: chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}

		err := ds.AddressRefStore.Add(addrRef)
		require.NoError(t, err)

		_, err = changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, nonExistentAddrStr)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found in address book")
	})
}

func TestGetOwnerTypeAndVersionV2(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}

	t.Run("finds owner in datastore", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Get the deployed registry address
		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break // Just take the first one
		}

		// Create datastore and save registry address
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		contract, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)

		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)

		// Save owner address to datastore
		v0 := semver.MustParse("1.0.0")
		ownerAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       owner.Hex(),
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       v0,
		}
		err = ds.AddressRefStore.Add(ownerAddrRef)
		require.NoError(t, err)

		tv, err := changeset.GetOwnerTypeAndVersionV2(*contract, ds.Addresses(), chain)
		require.NoError(t, err)
		require.NotNil(t, tv)
		assert.Equal(t, cldf.ContractType(types.RBACTimelock), tv.Type)
		assert.Equal(t, deployment.Version1_0_0, tv.Version)
	})

	t.Run("nil owner when owner not in datastore", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		// Create datastore and save only registry address (not owner)
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		contract, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)

		ownerTV, err := changeset.GetOwnerTypeAndVersionV2(*contract, ds.Addresses(), chain)

		require.NoError(t, err)
		assert.Nil(t, ownerTV)
	})
}

func TestNewOwnableV2(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}

	t.Run("creates OwnedContract for non-MCMS owner", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		// Create datastore and save registry
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		contract, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)

		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)

		// Setup owner as non-MCMS contract
		v0 := semver.MustParse("1.0.0")
		ownerAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       owner.Hex(),
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v0,
		}
		err = ds.AddressRefStore.Add(ownerAddrRef)
		require.NoError(t, err)

		ownedContract, err := changeset.NewOwnableV2(*contract, ds.Addresses(), chain)
		require.NoError(t, err)

		// Verify the owned contract contains the contract but no MCMS contracts
		assert.Equal(t, (*contract).Address(), ownedContract.Contract.Address())
		assert.Nil(t, ownedContract.McmsContracts)
	})

	t.Run("creates OwnedContract for MCMS owner", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		// Create datastore and save registry
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		contract, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)

		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)

		// Setup owner as timelock contract
		v0 := semver.MustParse("1.0.0")
		ownerAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       owner.Hex(),
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       v0,
		}
		err = ds.AddressRefStore.Add(ownerAddrRef)
		require.NoError(t, err)

		ownedContract, err := changeset.NewOwnableV2(*contract, ds.Addresses(), chain)

		require.NoError(t, err)
		assert.Equal(t, (*contract).Address(), ownedContract.Contract.Address())
		assert.NotNil(t, ownedContract.McmsContracts)
	})

	t.Run("no error when owner type lookup fails due to missing address in datastore (it is non-MCMS owned)", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		// Create datastore and save only registry (not owner)
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		contract, err := changeset.GetOwnableContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)

		// Don't add owner to datastore, so lookup will return nil TV and no error

		// Call NewOwnableV2, should not fail because owner is not in datastore, but should return a non-MCMS contract
		ownableContract, err := changeset.NewOwnableV2(*contract, ds.Addresses(), chain)
		require.NoError(t, err)
		assert.NotNil(t, ownableContract)
		assert.Nil(t, ownableContract.McmsContracts)
		assert.Equal(t, (*contract).Address(), ownableContract.Contract.Address())
	})
}

func TestGetOwnedContractV2(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}

	t.Run("successfully creates owned contract", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		resp, err := changeset.DeployCapabilityRegistryV2(env, &changeset.DeployRequestV2{ChainSel: chain.Selector})
		require.NoError(t, err)
		require.NotNil(t, resp)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		// Create datastore and save registry
		ds := datastore.NewMemoryDataStore()
		v1 := semver.MustParse("1.1.0")
		registryAddrRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       targetAddrStr,
			Type:          datastore.ContractType(changeset.CapabilitiesRegistry),
			Version:       v1,
		}
		err = ds.AddressRefStore.Add(registryAddrRef)
		require.NoError(t, err)

		ownedContract, err := changeset.GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, targetAddrStr)
		require.NoError(t, err)
		assert.NotNil(t, ownedContract)
		assert.NotNil(t, ownedContract.Contract)
		// MCMS contracts should be nil since owner is not in datastore
		assert.Nil(t, ownedContract.McmsContracts)
	})

	t.Run("errors when address not found in datastore", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		evmChains := env.BlockChains.EVMChains()
		chain := evmChains[slices.Collect(maps.Keys(evmChains))[0]]

		// Create empty datastore
		ds := datastore.NewMemoryDataStore()
		nonExistentAddr := testutils.NewAddress().String()

		_, err := changeset.GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](ds.Addresses(), chain, nonExistentAddr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found in datastore")
	})
}
