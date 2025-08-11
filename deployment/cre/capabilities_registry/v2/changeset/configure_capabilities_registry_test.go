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

	capabilities := []capabilities_registry_v2.CapabilitiesRegistryCapability{
		{
			CapabilityId: "write-chain",
			Metadata:     []byte("metadata"),
		},
		{
			CapabilityId: "read-chain",
			Metadata:     []byte("metadata"),
		},
	}

	nodes := []capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              test32byte(t, "0x01"),
			EncryptionPublicKey: test32byte(t, "0x01"),
			P2pId:               test32byte(t, "0x01"),
			CapabilityIds:       []string{"write-chain"},
			CsaKey:              test32byte(t, "0x01"),
		},
		{
			NodeOperatorId:      uint32(2),
			Signer:              test32byte(t, "0x02"),
			EncryptionPublicKey: test32byte(t, "0x02"),
			P2pId:               test32byte(t, "0x02"),
			CapabilityIds:       []string{"read-chain"},
			CsaKey:              test32byte(t, "0x02"),
		},
	}
	if len(nops) == 0 && len(nodes) == 0 {
		t.Skip("No NOPs or Nodes provided, skipping configuration test")
	}

	configureOutput, err := ConfigureCapabilitiesRegistry{}.Apply(env, ConfigureCapabilitiesRegistryInput{
		ChainSelector:               chainSelector,
		CapabilitiesRegistryAddress: capabilitiesRegistryAddress,
		Nops:                        nops,
		Capabilities:                capabilities,
		Nodes:                       nodes,
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

	registeredCapabilities, err := capabilitiesRegistry.GetCapabilities(nil)
	if err != nil {
		t.Fatalf("failed to get capabilities: %v", err)
	}

	require.Len(t, registeredCapabilities, len(capabilities), "should have registered the correct number of capabilities")
	for _, capability := range capabilities {
		registeredCapability, err := capabilitiesRegistry.GetCapability(nil, capability.CapabilityId)
		require.NoError(t, err, "failed to get registered capability")
		assert.Equal(t, capability.CapabilityId, registeredCapability.CapabilityId, "capability id should match")
		assert.Equal(t, capability.ConfigurationContract, registeredCapability.ConfigurationContract, "capability configuration contract should match")
		assert.Equal(t, capability.Metadata, registeredCapability.Metadata, "capability metadata should match")
	}

	registeredNodes, err := capabilitiesRegistry.GetNodes(nil)
	if err != nil {
		t.Fatalf("failed to get nodes: %v", err)
	}
	require.Len(t, registeredNodes, len(nodes), "should have registered the correct number of nodes")

	for i, node := range nodes {
		got, err := capabilitiesRegistry.GetNode(nil, node.P2pId)
		require.NoError(t, err) // careful here: the err is rpc, contract return empty info if it doesn't find the p2p as opposed to non-exist err.
		assert.Equal(t, node.EncryptionPublicKey, got.EncryptionPublicKey, "mismatch node encryption public key node %d", i)
		assert.Equal(t, node.Signer, got.Signer, "mismatch node signer node %d", i)
		assert.Equal(t, node.NodeOperatorId, got.NodeOperatorId, "mismatch node operator id node %d", i)
		assert.EqualValues(t, node.CapabilityIds, got.CapabilityIds, "mismatch node hashed capability ids node %d", i)
		assert.Equal(t, node.P2pId, got.P2pId, "mismatch node p2p id node %d", i)
	}
}
