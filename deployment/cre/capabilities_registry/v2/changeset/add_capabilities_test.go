package changeset_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
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
		CapabilityConfigs: []contracts.CapabilityConfig{{Capability: contracts.Capability{CapabilityID: "cap@1.0.0"}, Config: map[string]any{"k": "v"}}},
	})
	require.NoError(t, err)
}

func TestAddCapabilities_Apply(t *testing.T) {
	// SetupEnvV2 deploys a cap reg v2 and configures it. So no need to do that here, just leverage the existing one.
	fixture := test.SetupEnvV2(t, false)

	// Prepare new capability to add
	newCapID := "new-test-capability@1.0.0"
	newCapMetadata := map[string]any{"capabilityType": float64(0), "responseType": float64(0)}
	newCapConfig := map[string]any{
		"methodConfigs": map[string]any{
			"BalanceAt": map[string]any{
				"remoteExecutableConfig": map[string]any{
					"requestTimeout":            "30s",
					"serverMaxParallelRequests": 10,
				},
			},
			"LogTrigger": map[string]any{
				"remoteTriggerConfig": map[string]any{
					"registrationRefresh":     "20s",
					"registrationExpiry":      "60s",
					"minResponsesToAggregate": 2,
					"messageExpiry":           "120s",
					"maxBatchSize":            25,
					"batchCollectionPeriod":   "0.2s",
				},
			},
			"WriteReport": map[string]any{
				"remoteExecutableConfig": map[string]any{
					"transmissionSchedule":      "OneAtATime",
					"deltaStage":                "38.4s",
					"requestTimeout":            "268.8s",
					"serverMaxParallelRequests": 10,
					"requestHasherType":         "WriteReportExcludeSignatures",
				},
			},
		},
	}

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

	// Here we check that the uptyped input of the changeset was correctly applied on-chain as proto and can be decoded back to the same config
	// encoding to proto bytes is same as in the changeset and decoding to cap cfg is same as in the v2 registry syncer
	capCfg := pkg.CapabilityConfig(newCapConfig)
	configProtoBytes, err := capCfg.MarshalProto() // on chain it is stored as proto bytes
	require.NoError(t, err, "should be able to marshal new capability config to proto bytes")

	expectedConfig := new(pkg.CapabilityConfig) // expected decoded config, to be compared with decoded on-chain config
	err = expectedConfig.UnmarshalProto(configProtoBytes)
	require.NoError(t, err, "should be able to unmarshal new capability config from proto bytes")

	caps, err := pkg.GetCapabilities(nil, capReg)
	require.NoError(t, err)
	var found bool
	for _, c := range caps {
		if c.CapabilityId == newCapID {
			// metadata check
			var gotMeta map[string]any
			require.NoError(t, json.Unmarshal(c.Metadata, &gotMeta))
			assert.Equal(t, newCapMetadata, gotMeta)
			found = true
			break
		}
	}
	require.True(t, found, "new capability should be registered")

	// Nodes should now include new capability id
	nodes, err := pkg.GetNodes(nil, capReg)
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
			got := new(pkg.CapabilityConfig)
			require.NoError(t, got.UnmarshalProto(cfg.Config), "unmarshal capability config proto bytes should not error")
			if diff := cmp.Diff(expectedConfig, got, protocmp.Transform()); diff != "" {
				t.Errorf("capability config proto bytes should match: %s", diff)
			}

			cfgFound = true
		}
	}
	require.True(t, cfgFound, "don should have new capability configuration")
}
