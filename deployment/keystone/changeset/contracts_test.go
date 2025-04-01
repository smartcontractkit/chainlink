package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
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

		assert.Error(t, err)
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

		assert.Error(t, err)
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

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found in address book")
	})
}
