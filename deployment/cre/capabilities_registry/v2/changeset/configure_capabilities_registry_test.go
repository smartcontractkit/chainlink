package changeset

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

type testFixture struct {
	env                         cldf.Environment
	chainSelector               uint64
	capabilitiesRegistryAddress string
	nops                        []CapabilitiesRegistryNodeOperator
	capabilities                []CapabilitiesRegistryCapability
	nodes                       []CapabilitiesRegistryNodeParams
	DONs                        []CapabilitiesRegistryNewDONParams
	configureInput              ConfigureCapabilitiesRegistryInput
}

const (
	csaKey              = "4240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c258"
	signer1             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c251"
	signer2             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c252"
	p2pID1              = "p2p_12D3KooWM1111111111111111111111111111111111111111111"
	p2pID2              = "p2p_12D3KooWM1111111111111111111111111111111111111111112"
	encryptionPublicKey = "7240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c254"
	nodeID1             = "1"
)

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

func TestConfigureCapabilitiesRegistryInput_YAMLSerialization(t *testing.T) {
	originalInput := ConfigureCapabilitiesRegistryInput{
		ChainSelector:               123456789,
		CapabilitiesRegistryAddress: "0x1234567890123456789012345678901234567890",
		UseMCMS:                     true,
		Nops: []CapabilitiesRegistryNodeOperator{
			{
				Admin: common.HexToAddress("0x1111111111111111111111111111111111111111"),
				Name:  "Node Operator 1",
			},
			{
				Admin: common.HexToAddress("0x2222222222222222222222222222222222222222"),
				Name:  "Node Operator 2",
			},
		},
		Capabilities: []CapabilitiesRegistryCapability{
			{
				CapabilityID:          "write-chain@1.0.0",
				ConfigurationContract: common.HexToAddress("0x3333333333333333333333333333333333333333"),
				Metadata: map[string]interface{}{
					"capabilityType": 3,
					"responseType":   0,
				},
			},
			{
				CapabilityID:          "trigger@1.0.0",
				ConfigurationContract: common.Address{}, // Zero address
				Metadata: map[string]interface{}{
					"capabilityType": 0,
					"responseType":   0,
				},
			},
		},
		Nodes: []CapabilitiesRegistryNodeParams{
			{
				NodeOperatorID:      1,
				Signer:              signer1,
				P2pID:               p2pID1,
				EncryptionPublicKey: encryptionPublicKey,
				CsaKey:              csaKey,
				CapabilityIDs:       []string{"write-chain@1.0.0", "trigger@1.0.0"},
			},
		},
		DONs: []CapabilitiesRegistryNewDONParams{
			{
				Name:        "workflow-don-1",
				DonFamilies: []string{"workflow", "test"},
				Config: map[string]interface{}{
					"consensus": "basic",
					"timeout":   "30s",
				},
				CapabilityConfigurations: []CapabilitiesRegistryCapabilityConfiguration{
					{
						CapabilityID: "write-chain@1.0.0",
						Config: map[string]interface{}{
							"targetChain": "ethereum",
						},
					},
					{
						CapabilityID: "trigger@1.0.0",
						Config: map[string]interface{}{
							"schedule": "0 0 * * *",
						},
					},
				},
				Nodes:            []string{nodeID1},
				F:                1,
				IsPublic:         true,
				AcceptsWorkflows: true,
			},
		},
	}

	t.Run("marshal to YAML", func(t *testing.T) {
		yamlData, err := yaml.Marshal(originalInput)
		require.NoError(t, err, "should be able to marshal to YAML")
		require.NotEmpty(t, yamlData, "YAML data should not be empty")

		// Verify the YAML contains expected fields
		yamlStr := string(yamlData)
		assert.Contains(t, yamlStr, "chainSelector:", "should contain chainSelector field")
		assert.Contains(t, yamlStr, "capabilitiesRegistryAddress:", "should contain capabilitiesRegistryAddress field")
		assert.Contains(t, yamlStr, "useMCMS:", "should contain useMCMS field")
		assert.Contains(t, yamlStr, "nops:", "should contain nops field")
		assert.Contains(t, yamlStr, "capabilities:", "should contain capabilities field")
		assert.Contains(t, yamlStr, "nodes:", "should contain nodes field")
		assert.Contains(t, yamlStr, "dons:", "should contain dons field")
	})

	t.Run("unmarshal from YAML", func(t *testing.T) {
		// First marshal to YAML
		yamlData, err := yaml.Marshal(originalInput)
		require.NoError(t, err)

		// Then unmarshal back
		var unmarshaledInput ConfigureCapabilitiesRegistryInput
		err = yaml.Unmarshal(yamlData, &unmarshaledInput)
		require.NoError(t, err, "should be able to unmarshal from YAML")

		// Verify all fields are correctly deserialized
		assert.Equal(t, originalInput.ChainSelector, unmarshaledInput.ChainSelector)
		assert.Equal(t, originalInput.CapabilitiesRegistryAddress, unmarshaledInput.CapabilitiesRegistryAddress)
		assert.Equal(t, originalInput.UseMCMS, unmarshaledInput.UseMCMS)
		assert.Equal(t, originalInput.Nops, unmarshaledInput.Nops)
		assert.Equal(t, originalInput.Capabilities, unmarshaledInput.Capabilities)
		assert.Equal(t, originalInput.Nodes, unmarshaledInput.Nodes)
		assert.Equal(t, originalInput.DONs, unmarshaledInput.DONs)
	})

	t.Run("partial input with omitempty", func(t *testing.T) {
		// Test with minimal input (only required fields)
		minimalInput := ConfigureCapabilitiesRegistryInput{
			ChainSelector:               123456789,
			CapabilitiesRegistryAddress: "0x1234567890123456789012345678901234567890",
			UseMCMS:                     false,
			// Omit optional fields (nops, capabilities, nodes, dons)
		}

		yamlData, err := yaml.Marshal(minimalInput)
		require.NoError(t, err)

		yamlStr := string(yamlData)

		// Should contain required fields
		assert.Contains(t, yamlStr, "chainSelector:")
		assert.Contains(t, yamlStr, "capabilitiesRegistryAddress:")
		assert.Contains(t, yamlStr, "useMCMS:")

		// Should NOT contain optional fields due to omitempty
		assert.NotContains(t, yamlStr, "nops:")
		assert.NotContains(t, yamlStr, "capabilities:")
		assert.NotContains(t, yamlStr, "nodes:")
		assert.NotContains(t, yamlStr, "dons:")

		// Should be able to unmarshal back
		var unmarshaledMinimal ConfigureCapabilitiesRegistryInput
		err = yaml.Unmarshal(yamlData, &unmarshaledMinimal)
		require.NoError(t, err)

		assert.Equal(t, minimalInput.ChainSelector, unmarshaledMinimal.ChainSelector)
		assert.Equal(t, minimalInput.CapabilitiesRegistryAddress, unmarshaledMinimal.CapabilitiesRegistryAddress)
		assert.Equal(t, minimalInput.UseMCMS, unmarshaledMinimal.UseMCMS)
		assert.Empty(t, unmarshaledMinimal.Nops)
		assert.Empty(t, unmarshaledMinimal.Capabilities)
		assert.Empty(t, unmarshaledMinimal.Nodes)
		assert.Empty(t, unmarshaledMinimal.DONs)
	})
}

func TestConfigureCapabilitiesRegistryInput_YAMLFromFile(t *testing.T) {
	yamlConfig := `
chainSelector: 421614
capabilitiesRegistryAddress: "0x1234567890123456789012345678901234567890"
useMCMS: true
nops:
  - admin: "0x1111111111111111111111111111111111111111"
    name: "Node Operator Alpha"
  - admin: "0x2222222222222222222222222222222222222222"
    name: "Node Operator Beta"
capabilities:
  - capabilityID: "write-chain@1.0.0"
    configurationContract: "0x0000000000000000000000000000000000000000"
    metadata: 
      capabilityType: 3
      responseType: 0
  - capabilityID: "trigger@1.0.0"
    configurationContract: "0x0000000000000000000000000000000000000000"
    metadata: 
      capabilityType: 0
      responseType: 1
nodes:
  - nodeOperatorID: 1
    signer: ` + signer1 + `
    p2pID: ` + p2pID1 + `
    encryptionPublicKey: ` + encryptionPublicKey + `
    csaKey: ` + csaKey + `
    capabilityIDs: ["write-chain@1.0.0", "trigger@1.0.0"]
dons:
  - name: "workflow-don-production"
    donFamilies: ["workflow", "production"]
    config:
      consensus: "basic"
    capabilityConfigurations:
      - capabilityID: "write-chain@1.0.0"
        config:
          targetChain: "ethereum"
    nodes: [` + nodeID1 + `]
    f: 1
    isPublic: true
    acceptsWorkflows: true
`

	var input ConfigureCapabilitiesRegistryInput
	err := yaml.Unmarshal([]byte(yamlConfig), &input)
	require.NoError(t, err, "should be able to parse realistic YAML config")

	// Verify the parsed values
	assert.Equal(t, uint64(421614), input.ChainSelector)
	assert.Equal(t, "0x1234567890123456789012345678901234567890", input.CapabilitiesRegistryAddress)
	assert.True(t, input.UseMCMS)

	require.Len(t, input.Nops, 2)
	assert.Equal(t, "Node Operator Alpha", input.Nops[0].Name)
	assert.Equal(t, common.HexToAddress("0x1111111111111111111111111111111111111111"), input.Nops[0].Admin)

	require.Len(t, input.Capabilities, 2)
	assert.Equal(t, "write-chain@1.0.0", input.Capabilities[0].CapabilityID)
	assert.Equal(t, "trigger@1.0.0", input.Capabilities[1].CapabilityID)

	// Verify metadata is decoded properly
	expectedMetadata1 := map[string]interface{}{
		"capabilityType": 3,
		"responseType":   0,
	}
	expectedMetadata2 := map[string]interface{}{
		"capabilityType": 0,
		"responseType":   1,
	}
	assert.Equal(t, expectedMetadata1, input.Capabilities[0].Metadata)
	assert.Equal(t, expectedMetadata2, input.Capabilities[1].Metadata)

	require.Len(t, input.Nodes, 1)
	assert.Equal(t, uint32(1), input.Nodes[0].NodeOperatorID)
	assert.Equal(t, []string{"write-chain@1.0.0", "trigger@1.0.0"}, input.Nodes[0].CapabilityIDs)
	assert.Equal(t, csaKey, input.Nodes[0].CsaKey)

	require.Len(t, input.DONs, 1)
	assert.Equal(t, "workflow-don-production", input.DONs[0].Name)
	assert.Equal(t, []string{"workflow", "production"}, input.DONs[0].DonFamilies)
	assert.True(t, input.DONs[0].IsPublic)
	assert.True(t, input.DONs[0].AcceptsWorkflows)
	assert.Equal(t, uint8(1), input.DONs[0].F)

	// Verify config is decoded properly
	expectedConfig := map[string]interface{}{
		"consensus": "basic",
	}
	assert.Equal(t, expectedConfig, input.DONs[0].Config)

	// Verify capability configuration is decoded properly
	require.Len(t, input.DONs[0].CapabilityConfigurations, 1)
	assert.Equal(t, "write-chain@1.0.0", input.DONs[0].CapabilityConfigurations[0].CapabilityID)
	expectedCapConfig := map[string]interface{}{
		"targetChain": "ethereum",
	}
	assert.Equal(t, expectedCapConfig, input.DONs[0].CapabilityConfigurations[0].Config)
	assert.Equal(t, []string{nodeID1}, input.DONs[0].Nodes, "should contain the correct node IDs")
}

func setupCapabilitiesRegistryTest(t *testing.T) *testFixture {
	lggr := logger.Test(t)
	env, chainSelector := cre.BuildMinimalEnvironment(t, lggr)

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
	nops := []CapabilitiesRegistryNodeOperator{
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
	var writeChainCapabilityMetadata map[string]interface{}
	err = json.Unmarshal(writeChainCapability.Metadata, &writeChainCapabilityMetadata)
	require.NoError(t, err)

	triggerCapability := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "trigger@1.0.0",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 1, "responseType": 1}`),
	}
	var triggerCapabilityMetadata map[string]interface{}
	err = json.Unmarshal(triggerCapability.Metadata, &triggerCapabilityMetadata)
	require.NoError(t, err)

	capabilities := []CapabilitiesRegistryCapability{
		{
			CapabilityID: writeChainCapability.CapabilityId,
			Metadata:     writeChainCapabilityMetadata,
		},
		{
			CapabilityID: triggerCapability.CapabilityId,
			Metadata:     triggerCapabilityMetadata,
		},
	}

	nodes := []CapabilitiesRegistryNodeParams{
		{
			NodeOperatorID:      uint32(1),
			Signer:              signer1,
			EncryptionPublicKey: encryptionPublicKey,
			P2pID:               p2pID1,
			CapabilityIDs:       []string{writeChainCapability.CapabilityId, triggerCapability.CapabilityId},
			CsaKey:              csaKey,
		},
		{
			NodeOperatorID:      uint32(2),
			Signer:              signer2,
			EncryptionPublicKey: encryptionPublicKey,
			P2pID:               p2pID2,
			CapabilityIDs:       []string{writeChainCapability.CapabilityId, triggerCapability.CapabilityId},
			CsaKey:              csaKey,
		},
	}

	nodeSet := []string{}
	for _, n := range nodes {
		nodeSet = append(nodeSet, n.P2pID)
	}

	// Create capability configurations with readable config
	configMap := map[string]interface{}{
		"defaultConfig": map[string]interface{}{},
		"remoteTriggerConfig": map[string]interface{}{
			"registrationRefresh":     "20s",
			"registrationExpiry":      "60s",
			"minResponsesToAggregate": 2,
			"messageExpiry":           "120s",
		},
	}

	DONs := []CapabilitiesRegistryNewDONParams{
		{
			Name:        "test-don-1",
			DonFamilies: []string{"don-family-1"},
			Config: map[string]interface{}{
				"name": "test-don-v2-config",
				"type": "workflow",
			},
			CapabilityConfigurations: []CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityID: writeChainCapability.CapabilityId,
					Config:       configMap,
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
			Config: map[string]interface{}{
				"name": "test-don-v2-config",
				"type": "trigger",
			},
			CapabilityConfigurations: []CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityID: triggerCapability.CapabilityId,
					Config:       configMap,
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
		assert.Contains(t, registeredNops, nop.ToWrapper(), "node operator should be registered")
	}

	// Verify capabilities
	registeredCapabilities, err := capabilitiesRegistry.GetCapabilities(nil)
	require.NoError(t, err, "failed to get registered capabilities")
	require.Len(t, registeredCapabilities, len(fixture.capabilities), "should have registered the correct number of capabilities")
	for _, capability := range fixture.capabilities {
		registeredCapability, err := capabilitiesRegistry.GetCapability(nil, capability.CapabilityID)
		require.NoError(t, err, "failed to get registered capability")
		assert.Equal(t, capability.CapabilityID, registeredCapability.CapabilityId, "capability id should match")
		assert.Equal(t, capability.ConfigurationContract, registeredCapability.ConfigurationContract, "capability configuration contract should match")

		// Convert the struct metadata to bytes for comparison with blockchain data
		expectedMetadataBytes, err := json.Marshal(capability.Metadata)
		require.NoError(t, err, "failed to marshal expected metadata")
		assert.Equal(t, expectedMetadataBytes, registeredCapability.Metadata, "capability metadata should match")
	}

	// Verify nodes
	registeredNodes, err := capabilitiesRegistry.GetNodes(nil)
	require.NoError(t, err, "failed to get registered nodes")
	require.Len(t, registeredNodes, len(fixture.nodes), "should have registered the correct number of nodes")

	for i, node := range fixture.nodes {
		expectedSigner, err := pkg.HexStringTo32Bytes(node.Signer)
		require.NoError(t, err, "failed to convert signer hex string to bytes")

		expectedCsaKey, err := pkg.HexStringTo32Bytes(node.CsaKey)
		require.NoError(t, err, "failed to convert CSA key hex string to bytes")

		bytes32P2pID, err := p2pkey.MakePeerID(node.P2pID)
		require.NoError(t, err, "failed to convert P2P ID string to bytes")

		expectedEncryptionPublicKey, err := pkg.HexStringTo32Bytes(node.EncryptionPublicKey)
		require.NoError(t, err, "failed to convert encryption public key hex string to bytes")

		got, err := capabilitiesRegistry.GetNode(nil, bytes32P2pID)
		require.NoError(t, err) // careful here: the err is rpc, contract return empty info if it doesn't find the p2p as opposed to non-exist err.
		assert.Equal(t, expectedEncryptionPublicKey, got.EncryptionPublicKey, "mismatch node encryption public key node %d", i)
		assert.Equal(t, expectedSigner, got.Signer, "mismatch node signer node %d", i)
		assert.Equal(t, node.NodeOperatorID, got.NodeOperatorId, "mismatch node operator id node %d", i)
		assert.Equal(t, node.CapabilityIDs, got.CapabilityIds, "mismatch node hashed capability ids node %d", i)
		assert.Equal(t, [32]byte(bytes32P2pID), got.P2pId, "mismatch node p2p id node %d", i)
		assert.Equal(t, expectedCsaKey, got.CsaKey, "mismatch node CSA key node %d", i)
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

		// Convert our config map to JSON bytes for comparison
		expectedConfigBytes, err := json.Marshal(don.Config)
		require.NoError(t, err, "failed to marshal expected DON config")
		assert.Equal(t, expectedConfigBytes, foundDON.Config, "DON config should match")

		assert.Equal(t, don.F, foundDON.F, "DON F value should match")
		assert.Equal(t, don.IsPublic, foundDON.IsPublic, "DON isPublic flag should match")
		assert.Equal(t, don.AcceptsWorkflows, foundDON.AcceptsWorkflows, "DON accepts workflows flag should match")
	}

	donsFamilyTwo, err := capabilitiesRegistry.GetDONsInFamily(nil, "don-family-2")
	require.NoError(t, err, "failed to get DONs in family 'don-family-2'")
	require.Len(t, donsFamilyTwo, 1, "should have one DON in family 'don-family-2'")
	assert.Equal(t, big.NewInt(2), donsFamilyTwo[0], "DON ID should match")
}
