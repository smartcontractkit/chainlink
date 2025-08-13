package changeset

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type testFixture struct {
	env                         cldf.Environment
	chainSelector               uint64
	capabilitiesRegistryAddress string
	nops                        []capabilities_registry_v2.CapabilitiesRegistryNodeOperator
	capabilities                []capabilities_registry_v2.CapabilitiesRegistryCapability
	nodes                       []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	DONs                        []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
	configureInput              ConfigureCapabilitiesRegistryInput
}

func TestConfigureCapabilitiesRegistry(t *testing.T) {
	fixture := setupCapabilitiesRegistryTest(t)

	t.Run("single configuration", func(t *testing.T) {
		t.Log("Starting capabilities registry configuration...")
		configureOutput, err := ConfigureCapabilitiesRegistry{}.Apply(fixture.env, fixture.configureInput)
		t.Logf("Configuration result: err=%v, output=%v", err, configureOutput)
		require.NoError(t, err, "configuration should succeed")
		require.NotNil(t, configureOutput, "configuration output should not be nil")
		t.Logf("Capabilities registry configured successfully")

		// Verify the configuration
		verifyCapabilitiesRegistryConfiguration(t, fixture)
	})

	t.Run("idempotency test - double configuration", func(t *testing.T) {
		t.Log("Starting first capabilities registry configuration...")
		configureOutput1, err := ConfigureCapabilitiesRegistry{}.Apply(fixture.env, fixture.configureInput)
		require.NoError(t, err, "first configuration should succeed")
		require.NotNil(t, configureOutput1, "first configuration output should not be nil")
		t.Logf("First configuration completed successfully")

		t.Log("Starting second capabilities registry configuration (idempotency test)...")
		configureOutput2, err := ConfigureCapabilitiesRegistry{}.Apply(fixture.env, fixture.configureInput)
		require.NoError(t, err, "second configuration should succeed (idempotent)")
		require.NotNil(t, configureOutput2, "second configuration output should not be nil")
		t.Logf("Second configuration completed successfully - idempotency verified")

		// Verify that the final state is still correct
		verifyCapabilitiesRegistryConfiguration(t, fixture)
	})
}

func setupCapabilitiesRegistryTest(t *testing.T) *testFixture {
	lggr := logger.Test(t)
	env, chainSelector := BuildMinimalEnvironment(t, lggr)

	// Apply the changeset to deploy the V2 capabilities registry
	t.Log("Running deployment changeset...")
	deployOutput, err := DeployCapabilitiesRegistry{}.Apply(env, DeployCapabilitiesRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "test-capabilities-registry-v2",
	})
	require.NoError(t, err, "failed to apply deployment changeset")
	require.NotNil(t, deployOutput, "deployment output should not be nil")
	t.Logf("Deployment result: err=%v, output=%v", err, deployOutput)

	capabilitiesRegistryAddress := deployOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-capabilities-registry-v2"))[0].Address

	// Setup test data
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

	writeChainCapability := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`),
	}

	triggerCapability := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "trigger@1.0.0",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 1, "responseType": 1}`),
	}

	capabilities := []capabilities_registry_v2.CapabilitiesRegistryCapability{
		{
			CapabilityId: writeChainCapability.CapabilityId,
			Metadata:     writeChainCapability.Metadata,
		},
		{
			CapabilityId: triggerCapability.CapabilityId,
			Metadata:     triggerCapability.Metadata,
		},
	}

	nodes := []capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              test32byte(t, "0x01"),
			EncryptionPublicKey: test32byte(t, "0x01"),
			P2pId:               test32byte(t, "0x01"),
			CapabilityIds:       []string{writeChainCapability.CapabilityId, triggerCapability.CapabilityId},
			CsaKey:              test32byte(t, "0x01"),
		},
		{
			NodeOperatorId:      uint32(2),
			Signer:              test32byte(t, "0x02"),
			EncryptionPublicKey: test32byte(t, "0x02"),
			P2pId:               test32byte(t, "0x02"),
			CapabilityIds:       []string{writeChainCapability.CapabilityId, triggerCapability.CapabilityId},
			CsaKey:              test32byte(t, "0x02"),
		},
	}

	nodeSet := [][32]byte{}
	for _, n := range nodes {
		nodeSet = append(nodeSet, n.P2pId)
	}

	// Create capability configurations
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	DONs := []capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:        "test-don-1",
			DonFamilies: []string{"don-family-1"},
			Config:      []byte("test-don-v2-config"),
			CapabilityConfigurations: []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: writeChainCapability.CapabilityId,
					Config:       configb,
				},
			},
			Nodes:            nodeSet,
			F:                1,
			IsPublic:         true,
			AcceptsWorkflows: false,
		},
		{
			Name:        "test-don-2",
			DonFamilies: []string{"don-family-2"},
			Config:      []byte("test-don-v2-config"),
			CapabilityConfigurations: []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: triggerCapability.CapabilityId,
					Config:       configb,
				},
			},
			Nodes:            nodeSet,
			F:                1,
			IsPublic:         true,
			AcceptsWorkflows: false,
		},
	}

	configureInput := ConfigureCapabilitiesRegistryInput{
		ChainSelector:               chainSelector,
		CapabilitiesRegistryAddress: capabilitiesRegistryAddress,
		UseMCMS:                     false,
		Nops:                        nops,
		Capabilities:                capabilities,
		Nodes:                       nodes,
		DONs:                        DONs,
	}

	return &testFixture{
		env:                         env,
		chainSelector:               chainSelector,
		capabilitiesRegistryAddress: capabilitiesRegistryAddress,
		nops:                        nops,
		capabilities:                capabilities,
		nodes:                       nodes,
		DONs:                        DONs,
		configureInput:              configureInput,
	}
}

func verifyCapabilitiesRegistryConfiguration(t *testing.T, fixture *testFixture) {
	capabilitiesRegistry, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(fixture.capabilitiesRegistryAddress),
		fixture.env.BlockChains.EVMChains()[fixture.chainSelector].Client,
	)
	require.NoError(t, err, "failed to create CapabilitiesRegistry instance")
	t.Logf("CapabilitiesRegistry instance created at address: %s", fixture.capabilitiesRegistryAddress)

	// Verify node operators
	registeredNops, err := capabilitiesRegistry.GetNodeOperators(nil)
	require.NoError(t, err, "failed to get registered node operators")
	require.Len(t, registeredNops, len(fixture.nops), "should have registered the correct number of node operators")
	for _, nop := range fixture.nops {
		assert.Contains(t, registeredNops, nop, "node operator should be registered")
	}

	// Verify capabilities
	registeredCapabilities, err := capabilitiesRegistry.GetCapabilities(nil)
	require.NoError(t, err, "failed to get registered capabilities")
	require.Len(t, registeredCapabilities, len(fixture.capabilities), "should have registered the correct number of capabilities")
	for _, capability := range fixture.capabilities {
		registeredCapability, err := capabilitiesRegistry.GetCapability(nil, capability.CapabilityId)
		require.NoError(t, err, "failed to get registered capability")
		assert.Equal(t, capability.CapabilityId, registeredCapability.CapabilityId, "capability id should match")
		assert.Equal(t, capability.ConfigurationContract, registeredCapability.ConfigurationContract, "capability configuration contract should match")
		assert.Equal(t, capability.Metadata, registeredCapability.Metadata, "capability metadata should match")
	}

	// Verify nodes
	registeredNodes, err := capabilitiesRegistry.GetNodes(nil)
	require.NoError(t, err, "failed to get registered nodes")
	require.Len(t, registeredNodes, len(fixture.nodes), "should have registered the correct number of nodes")

	for i, node := range fixture.nodes {
		got, err := capabilitiesRegistry.GetNode(nil, node.P2pId)
		require.NoError(t, err) // careful here: the err is rpc, contract return empty info if it doesn't find the p2p as opposed to non-exist err.
		assert.Equal(t, node.EncryptionPublicKey, got.EncryptionPublicKey, "mismatch node encryption public key node %d", i)
		assert.Equal(t, node.Signer, got.Signer, "mismatch node signer node %d", i)
		assert.Equal(t, node.NodeOperatorId, got.NodeOperatorId, "mismatch node operator id node %d", i)
		assert.Equal(t, node.CapabilityIds, got.CapabilityIds, "mismatch node hashed capability ids node %d", i)
		assert.Equal(t, node.P2pId, got.P2pId, "mismatch node p2p id node %d", i)
	}

	// Verify DONs
	registeredDONs, err := capabilitiesRegistry.GetDONs(nil)
	require.NoError(t, err, "failed to get registered DONs")
	require.Len(t, registeredDONs, len(fixture.DONs), "should have registered the correct number of DONs")

	// Verify each expected DON is registered with correct properties
	for _, don := range fixture.DONs {
		var foundDON *capabilities_registry_v2.CapabilitiesRegistryDONInfo
		for _, registeredDON := range registeredDONs {
			if registeredDON.Name == don.Name {
				foundDON = &registeredDON
				break
			}
		}

		require.NotNil(t, foundDON, "DON %s should be found in registered DONs", don.Name)
		assert.Equal(t, don.Name, foundDON.Name, "DON name should match")
		assert.Equal(t, don.DonFamilies, foundDON.DonFamilies, "DON families should match")
		assert.Equal(t, don.Config, foundDON.Config, "DON config should match")
		assert.Equal(t, don.F, foundDON.F, "DON F value should match")
		assert.Equal(t, don.IsPublic, foundDON.IsPublic, "DON isPublic flag should match")
		assert.Equal(t, don.AcceptsWorkflows, foundDON.AcceptsWorkflows, "DON accepts workflows flag should match")
	}

	donsFamilyTwo, err := capabilitiesRegistry.GetDONsInFamily(nil, "don-family-2")
	require.NoError(t, err, "failed to get DONs in family 'don-family-2'")
	require.Len(t, donsFamilyTwo, 1, "should have one DON in family 'don-family-2'")
	assert.Equal(t, big.NewInt(2), donsFamilyTwo[0], "DON ID should match")
}
