package changeset_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestAddCapabilities_VerifyPreconditions(t *testing.T) {
	cs := changeset.AddCapabilities{}

	env := test.SetupEnvV2(t, false)
	chainSelector := env.RegistrySelector

	// Missing DON name
	err := cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonName:           "", // invalid
		CapabilityConfigs: []contracts.CapabilityConfig{{Capability: contracts.Capability{CapabilityID: "cap@1.0.0"}}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DONName")

	// Missing capability configs
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonName:           "don-1",
		CapabilityConfigs: nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "capabilityConfigs")

	// Valid
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonName:           "don-1",
		CapabilityConfigs: []contracts.CapabilityConfig{{Capability: contracts.Capability{CapabilityID: "cap@1.0.0"}, Config: map[string]interface{}{"k": "v"}}},
	})
	require.NoError(t, err)
}

func TestAddCapabilities_Apply(t *testing.T) {
	// SetupEnvV2 deploys a cap reg v2 and configures it. So no need to do that here, just leverage the existing one.
	fixture := test.SetupEnvV2(t, false)

	// Prepare new capability to add
	newCapID := "new-test-capability@1.0.0"
	newCapMetadata := map[string]interface{}{"capabilityType": float64(0), "responseType": float64(0)}
	newCapConfig := map[string]interface{}{"newParam": "value"}

	input := changeset.AddCapabilitiesInput{
		RegistryChainSel:  fixture.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonName:           test.DONName,
		CapabilityConfigs: []contracts.CapabilityConfig{{
			Capability: contracts.Capability{
				CapabilityID:          newCapID,
				ConfigurationContract: common.Address{},
				Metadata:              newCapMetadata,
			},
			Config: newCapConfig,
		}},
		Force: true,
	}

	// Preconditions
	err := changeset.AddCapabilities{}.VerifyPreconditions(*fixture.Env, input)
	require.NoError(t, err)

	// Apply
	_, err = changeset.AddCapabilities{}.Apply(*fixture.Env, input)
	require.NoError(t, err)

	// Validate on-chain state
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		fixture.RegistryAddress,
		fixture.Env.BlockChains.EVMChains()[fixture.RegistrySelector].Client,
	)
	require.NoError(t, err)

	caps, err := capReg.GetCapabilities(nil)
	require.NoError(t, err)
	var found bool
	for _, c := range caps {
		if c.CapabilityId == newCapID {
			// metadata check
			var gotMeta map[string]interface{}
			require.NoError(t, json.Unmarshal(c.Metadata, &gotMeta))
			assert.Equal(t, newCapMetadata, gotMeta)
			found = true
			break
		}
	}
	require.True(t, found, "new capability should be registered")

	// Nodes should now include new capability id
	nodes, err := capReg.GetNodes(nil)
	require.NoError(t, err)
	for _, n := range nodes {
		assert.Contains(t, n.CapabilityIds, newCapID, "node should have new capability id appended")
	}

	// DON capability configurations should include new capability config
	don, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)
	var cfgFound bool
	for _, cfg := range don.CapabilityConfigurations {
		if cfg.CapabilityId == newCapID {
			var gotCfg map[string]interface{}
			require.NoError(t, json.Unmarshal(cfg.Config, &gotCfg))
			assert.Equal(t, newCapConfig, gotCfg)
			cfgFound = true
		}
	}
	require.True(t, cfgFound, "don should have new capability configuration")
}
