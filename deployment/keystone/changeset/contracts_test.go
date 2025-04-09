package changeset_test

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
)

func TestGetOwnableContract(t *testing.T) {
	t.Parallel()

	chain := memory.NewMemoryChain(t, chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector)

	t.Run("finds contract when targetAddr is provided", func(t *testing.T) {
		t.Parallel()

		addrBook := deployment.NewMemoryAddressBook()
		targetAddr := testutils.NewAddress()
		targetAddrStr := targetAddr.String()
		tv := deployment.TypeAndVersion{Type: changeset.CapabilitiesRegistry, Version: deployment.Version1_1_0}
		err := addrBook.Save(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector, targetAddrStr, tv)
		require.NoError(t, err)

		c, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)
		assert.NotNil(t, c)
		contract := *c
		assert.Equal(t, targetAddr, contract.Address())
	})

	t.Run("errors when multiple contracts found without targetAddr", func(t *testing.T) {
		t.Parallel()

		addrBook := deployment.NewMemoryAddressBook()
		targetAddr1 := testutils.NewAddress()
		targetAddrStr1 := targetAddr1.String()
		targetAddr2 := testutils.NewAddress()
		targetAddrStr2 := targetAddr2.String()
		mockAddresses := map[string]deployment.TypeAndVersion{
			targetAddrStr2: {Type: changeset.KeystoneForwarder, Version: deployment.Version1_1_0},
			targetAddrStr1: {Type: changeset.KeystoneForwarder, Version: deployment.Version1_1_0},
		}
		err := addrBook.Save(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector, targetAddrStr1, mockAddresses[targetAddrStr1])
		require.NoError(t, err)
		err = addrBook.Save(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector, targetAddrStr2, mockAddresses[targetAddrStr2])
		require.NoError(t, err)

		// No target address provided
		_, err = changeset.GetOwnableContract[*forwarder.KeystoneForwarder](addrBook, chain, nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple contracts")
		assert.Contains(t, err.Error(), "must provide a `targetAddr`")
	})

	t.Run("errors when no contracts of the requested type exist", func(t *testing.T) {
		t.Parallel()

		targetAddr := testutils.NewAddress()
		targetAddrStr := targetAddr.String()
		addrBook := deployment.NewMemoryAddressBook()
		mockAddresses := map[string]deployment.TypeAndVersion{
			targetAddrStr: {Type: "DifferentType", Version: deployment.Version1_0_0},
		}
		err := addrBook.Save(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector, targetAddrStr, mockAddresses[targetAddrStr])
		require.NoError(t, err)

		// No target address provided
		_, err = changeset.GetOwnableContract[*forwarder.KeystoneForwarder](addrBook, chain, nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no contract of type")
	})

	t.Run("errors when targetAddr not found in address book", func(t *testing.T) {
		t.Parallel()

		targetAddr := testutils.NewAddress()
		targetAddrStr := targetAddr.String()
		nonExistentAddr := testutils.NewAddress()
		nonExistentAddrStr := nonExistentAddr.String()
		addrBook := deployment.NewMemoryAddressBook()
		mockAddresses := map[string]deployment.TypeAndVersion{
			targetAddrStr: {Type: changeset.CapabilitiesRegistry, Version: deployment.Version1_1_0},
		}
		err := addrBook.Save(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector, targetAddrStr, mockAddresses[targetAddrStr])
		require.NoError(t, err)

		_, err = changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &nonExistentAddrStr)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found in address book")
	})
}

func TestGetOwnerTypeAndVersion(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 3,
	}

	t.Run("finds owner in address book", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		chain := env.Chains[env.AllChainSelectors()[0]]
		resp, err := changeset.DeployCapabilityRegistry(env, chain.Selector)
		require.NoError(t, err)
		require.NotNil(t, resp)
		err = env.ExistingAddresses.Merge(resp.AddressBook)
		require.NoError(t, err)
		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)
		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break // Just take the first one
		}

		addrBook := env.ExistingAddresses

		contract, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)
		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)
		mockAddresses := map[string]deployment.TypeAndVersion{
			owner.Hex(): {Type: types.RBACTimelock, Version: deployment.Version1_0_0},
		}
		err = addrBook.Save(chain.Selector, owner.Hex(), mockAddresses[owner.Hex()])
		require.NoError(t, err)

		tv, err := changeset.GetOwnerTypeAndVersion(*contract, addrBook, chain)
		require.NoError(t, err)
		assert.Equal(t, types.RBACTimelock, tv.Type)
		assert.Equal(t, deployment.Version1_0_0, tv.Version)
	})

	t.Run("nil owner when owner not in address book", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		chain := env.Chains[env.AllChainSelectors()[0]]
		resp, err := changeset.DeployCapabilityRegistry(env, chain.Selector)
		require.NoError(t, err)
		require.NotNil(t, resp)
		err = env.ExistingAddresses.Merge(resp.AddressBook)
		require.NoError(t, err)
		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)
		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break // Just take the first one
		}
		addrBook := env.ExistingAddresses
		contract, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)

		ownerTV, err := changeset.GetOwnerTypeAndVersion(*contract, addrBook, chain)

		require.NoError(t, err)
		assert.Nil(t, ownerTV)
	})
}

func TestNewOwnable(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 3,
	}

	t.Run("creates OwnedContract for non-MCMS owner", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		chain := env.Chains[env.AllChainSelectors()[0]]
		resp, err := changeset.DeployCapabilityRegistry(env, chain.Selector)
		require.NoError(t, err)
		require.NotNil(t, resp)
		err = env.ExistingAddresses.Merge(resp.AddressBook)
		require.NoError(t, err)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		addrBook := env.ExistingAddresses

		contract, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)
		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)

		// Setup owner as non-MCMS contract
		mockAddresses := map[string]deployment.TypeAndVersion{
			owner.Hex(): {Type: changeset.CapabilitiesRegistry, Version: deployment.Version1_0_0},
		}
		err = addrBook.Save(chain.Selector, owner.Hex(), mockAddresses[owner.Hex()])
		require.NoError(t, err)

		ownedContract, err := changeset.NewOwnable(*contract, addrBook, chain)
		require.NoError(t, err)

		// Verify the owned contract contains the contract but no MCMS contracts
		assert.Equal(t, (*contract).Address(), ownedContract.Contract.Address())
		assert.Nil(t, ownedContract.McmsContracts)
	})

	t.Run("creates OwnedContract for MCMS owner", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		chain := env.Chains[env.AllChainSelectors()[0]]
		resp, err := changeset.DeployCapabilityRegistry(env, chain.Selector)
		require.NoError(t, err)
		require.NotNil(t, resp)
		err = env.ExistingAddresses.Merge(resp.AddressBook)
		require.NoError(t, err)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		addrBook := env.ExistingAddresses

		contract, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)
		owner, err := (*contract).Owner(nil)
		require.NoError(t, err)

		// Setup owner as timelock contract
		mockAddresses := map[string]deployment.TypeAndVersion{
			owner.Hex(): {Type: types.RBACTimelock, Version: deployment.Version1_0_0},
		}
		err = addrBook.Save(chain.Selector, owner.Hex(), mockAddresses[owner.Hex()])
		require.NoError(t, err)

		ownedContract, err := changeset.NewOwnable(*contract, addrBook, chain)

		require.NoError(t, err)
		assert.Equal(t, (*contract).Address(), ownedContract.Contract.Address())
		assert.NotNil(t, ownedContract.McmsContracts)
	})

	t.Run("no error when owner type lookup fails due to missing address on address book (it is non-MCMS owned)", func(t *testing.T) {
		t.Parallel()

		env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
		chain := env.Chains[env.AllChainSelectors()[0]]
		resp, err := changeset.DeployCapabilityRegistry(env, chain.Selector)
		require.NoError(t, err)
		require.NotNil(t, resp)
		err = env.ExistingAddresses.Merge(resp.AddressBook)
		require.NoError(t, err)

		addrs, err := resp.AddressBook.AddressesForChain(chain.Selector)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// Get the first address from the map
		var targetAddrStr string
		for addr := range addrs {
			targetAddrStr = addr
			break
		}

		addrBook := env.ExistingAddresses

		contract, err := changeset.GetOwnableContract[*capabilities_registry.CapabilitiesRegistry](addrBook, chain, &targetAddrStr)
		require.NoError(t, err)

		// Don't add owner to address book, so lookup will return nil TV and no error

		// Call NewOwnable, should not fail because owner is not in address book, but should return a non-MCMS contract
		ownableContract, err := changeset.NewOwnable(*contract, addrBook, chain)
		require.NoError(t, err)
		assert.NotNil(t, ownableContract)
		assert.Nil(t, ownableContract.McmsContracts)
		assert.Equal(t, (*contract).Address(), ownableContract.Contract.Address())
	})
}
