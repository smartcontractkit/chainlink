package strategies_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	capregchangeset "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestCreateStrategy_OwnerCheck(t *testing.T) {
	t.Parallel()

	t.Run("deployer owned, no mcmsConfig", func(t *testing.T) {
		t.Parallel()
		h := test.NewTestHarness(t)
		env := h.Runtime.Environment()
		chain := env.BlockChains.EVMChains()[h.RegistrySelector]

		s, err := strategies.CreateStrategy(chain, env, nil, nil, h.RegistryAddress, "test")
		require.NoError(t, err)
		assert.IsType(t, &strategies.SimpleTransaction{}, s)
	})

	h := test.NewTestHarness(t, test.WithMCMS())
	env := h.Runtime.Environment()
	chain := env.BlockChains.EVMChains()[h.RegistrySelector]
	cfg := &contracts.MCMSConfig{
		TimelockQualifierPerChain: map[uint64]string{h.RegistrySelector: ""},
	}
	mcmsContracts, err := strategies.GetMCMSContracts(env, h.RegistrySelector, *cfg)
	require.NoError(t, err)
	timelock := mcmsContracts.Timelock.Address()

	t.Run("timelock owned, matching mcmsConfig", func(t *testing.T) {
		s, err := strategies.CreateStrategy(chain, env, cfg, mcmsContracts, h.RegistryAddress, "test")
		require.NoError(t, err)
		assert.IsType(t, &strategies.MCMSTransaction{}, s)
	})

	t.Run("timelock owned, no mcmsConfig", func(t *testing.T) {
		_, err := strategies.CreateStrategy(chain, env, nil, nil, h.RegistryAddress, "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), h.RegistryAddress.Hex())
		assert.Contains(t, err.Error(), "is owned by "+timelock.Hex())
		assert.Contains(t, err.Error(), chain.DeployerKey.From.Hex())
		assert.Contains(t, err.Error(), "provide mcmsConfig")
	})

	// A call that is not owner-only (a workflow owner acting on its own workflows) is sent by an
	// account that is not the contract's owner, so the variant without the check must accept it.
	t.Run("timelock owned, no mcmsConfig, call is not owner-only", func(t *testing.T) {
		s, err := strategies.CreateStrategyForNonOwnerCall(chain, env, nil, nil, h.RegistryAddress, "test")
		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("timelock owned, mcmsConfig names a different timelock", func(t *testing.T) {
		other := *mcmsContracts
		otherTimelock := common.HexToAddress("0x00000000000000000000000000000000000000aa")
		// Only the address of the timelock is read by the check, so re-pointing the binding is enough.
		tl, err := bindings.NewRBACTimelock(otherTimelock, chain.Client)
		require.NoError(t, err)
		other.Timelock = tl

		_, err = strategies.CreateStrategy(chain, env, cfg, &other, h.RegistryAddress, "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mcmsConfig names timelock "+otherTimelock.Hex())
		assert.Contains(t, err.Error(), "is owned by "+timelock.Hex())
	})

	t.Run("deployer owned, mcmsConfig", func(t *testing.T) {
		require.NoError(t, h.Runtime.Exec(runtime.ChangesetTask(capregchangeset.DeployCapabilitiesRegistry{},
			capregchangeset.DeployCapabilitiesRegistryInput{ChainSelector: h.RegistrySelector, Qualifier: "deployer-owned"})))
		refs := h.Runtime.Environment().DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(h.RegistrySelector),
			datastore.AddressRefByQualifier("deployer-owned"),
		)
		require.Len(t, refs, 1)
		addr := common.HexToAddress(refs[0].Address)

		_, err := strategies.CreateStrategy(chain, env, cfg, mcmsContracts, addr, "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "owned by the deployer key "+chain.DeployerKey.From.Hex())
		assert.Contains(t, err.Error(), "accept-ownership proposal has not executed")

		// without mcmsConfig the same contract is fine
		_, err = strategies.CreateStrategy(chain, env, nil, nil, addr, "test")
		require.NoError(t, err)
	})

	t.Run("target without owner()", func(t *testing.T) {
		_, err := strategies.CreateStrategy(chain, env, nil, nil, common.HexToAddress("0x01"), "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read owner()")
	})
}
