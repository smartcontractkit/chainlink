package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

func TestConfigureCapabilitiesRegistry(t *testing.T) {
	lggr := logger.Test(t)

	env, chainSelector := BuildMinimalEnvironment(t, lggr)

	// Apply the changeset to deploy the V2 capabilities registry
	t.Log("Running deployment changeset...")
	deployOutput, err := DeployCapabilitiesRegistry{}.Apply(env, DeployCapabilitiesRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "test-capabilities-registry-v2",
	})
	t.Logf("Deployment result: err=%v, output=%v", err, deployOutput)

	if err != nil {
		t.Fatalf("deployment apply failed: %v", err)
	}
	require.NotNil(t, deployOutput, "deployment output should not be nil")
	t.Logf("Deployment applied successfully")

	capabilitiesRegistryAddress := deployOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-capabilities-registry-v2"))[0].Address

	//  Configure the capabilities registry
	t.Log("Starting capabilities registry configuration...")
	nops := []capabilities_registry_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: common.HexToAddress("0x01"),
			Name:  "test nop1",
		},
		{
			Admin: common.HexToAddress("0x02"),
			Name:  "test nop2",
		},
	}

	configureOutput, err := ConfigureCapabilitiesRegistry{}.Apply(env, ConfigureCapabilitiesRegistryInput{
		ChainSelector:               chainSelector,
		CapabilitiesRegistryAddress: capabilitiesRegistryAddress,
		Nops:                        nops,
	})
	t.Logf("Configuration result: err=%v, output=%v", err, configureOutput)
	if err != nil {
		t.Fatalf("configuration apply failed: %v", err)
	}
	require.NotNil(t, configureOutput, "configuration output should not be nil")
	t.Logf("Capabilities registry configured successfully")

	capabilitiesRegistry, err := capabilities_registry_v2.NewCapabilitiesRegistry(common.HexToAddress(capabilitiesRegistryAddress), env.BlockChains.EVMChains()[chainSelector].Client)
	if err != nil {
		t.Fatalf("failed to create CapabilitiesRegistry instance: %v", err)
	}
	t.Logf("CapabilitiesRegistry instance created at address: %s", capabilitiesRegistryAddress)

	// Verify the capabilities registry is configured correctly
	registeredNops, err := capabilitiesRegistry.GetNodeOperators(nil)
	if err != nil {
		t.Fatalf("failed to get node operators: %v", err)
	}

	require.Len(t, registeredNops, len(nops), "should have registered the correct number of node operators")
	for _, nop := range nops {
		assert.Contains(t, registeredNops, nop, "node operator should be registered")
	}
}
