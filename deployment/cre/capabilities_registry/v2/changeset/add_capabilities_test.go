package changeset_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

var (
	newCapID       = "new-test-capability@1.0.0"
	anotherCapID   = "another-test-capability@1.0.0"
	newCapMetadata = map[string]any{"capabilityType": float64(0), "responseType": float64(0)}
	newCapConfig   = map[string]any{
		"restrictedConfig": map[string]any{
			"fields": map[string]any{
				"spendRatios": map[string]any{
					"mapValue": map[string]any{
						"fields": map[string]any{
							"RESOURCE_TYPE_COMPUTE": map[string]any{
								"stringValue": "1.0",
							},
						},
					},
				},
			},
		},
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
)

func TestAddCapabilities_VerifyPreconditions(t *testing.T) {
	cs := changeset.AddCapabilities{}

	env := test.SetupEnvV2(t, false)
	chainSelector := env.RegistrySelector

	capCfg := []contracts.CapabilityConfig{{Capability: contracts.Capability{CapabilityID: "cap@1.0.0"}, Config: map[string]any{"k": "v"}}}

	// Empty map
	err := cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:     chainSelector,
		RegistryQualifier:    "qual",
		DonCapabilityConfigs: nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "donCapabilityConfigs must contain at least one DON entry")

	// Empty DON name key
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			"": capCfg,
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty strings")

	// Empty config list for a DON
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			"don-1": {},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one capability config")

	// Valid (single DON)
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			"don-1": capCfg,
		},
	})
	require.NoError(t, err)

	// Valid (multiple DONs)
	err = cs.VerifyPreconditions(*env.Env, changeset.AddCapabilitiesInput{
		RegistryChainSel:  chainSelector,
		RegistryQualifier: "qual",
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			"don-1": capCfg,
			"don-2": capCfg,
		},
	})
	require.NoError(t, err)
}

func addNewCapability(t *testing.T, fixture *test.EnvWrapperV2, capID string) {
	input := changeset.AddCapabilitiesInput{
		RegistryChainSel:  fixture.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			test.DONName: {{
				Capability: contracts.Capability{
					CapabilityID:          capID,
					ConfigurationContract: common.Address{},
					Metadata:              newCapMetadata,
				},
				Config: newCapConfig,
			}},
		},
		Force: true,
	}

	// Preconditions
	err := changeset.AddCapabilities{}.VerifyPreconditions(*fixture.Env, input)
	require.NoError(t, err)

	// Apply
	_, err = changeset.AddCapabilities{}.Apply(*fixture.Env, input)
	require.NoError(t, err)
}

func requireCapability(t *testing.T, fixture *test.EnvWrapperV2, capID string) {
	// Validate on-chain state
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		fixture.RegistryAddress,
		fixture.Env.BlockChains.EVMChains()[fixture.RegistrySelector].Client,
	)
	require.NoError(t, err)

	caps, err := pkg.GetCapabilities(nil, capReg)
	require.NoError(t, err)
	var found bool
	for _, c := range caps {
		if c.CapabilityId == capID {
			// metadata check
			var gotMeta map[string]any
			require.NoError(t, json.Unmarshal(c.Metadata, &gotMeta))
			assert.Equal(t, newCapMetadata, gotMeta)
			found = true
			break
		}
	}
	require.True(t, found, "new capability %s should be registered", capID)

	// Nodes should now include new capability id
	nodes, err := pkg.GetNodes(nil, capReg)
	require.NoError(t, err)
	for _, n := range nodes {
		assert.Contains(t, n.CapabilityIds, capID, "node should have new capability id appended")
	}

	// Here we check that the uptyped input of the changeset was correctly applied on-chain as proto and can be decoded back to the same config
	// encoding to proto bytes is same as in the changeset and decoding to cap cfg is same as in the v2 registry syncer
	capCfg := pkg.CapabilityConfig(newCapConfig)
	configProtoBytes, err := capCfg.MarshalProto() // on chain it is stored as proto bytes
	require.NoError(t, err, "should be able to marshal new capability config to proto bytes")

	expectedConfig := new(pkg.CapabilityConfig) // expected decoded config, to be compared with decoded on-chain config
	err = expectedConfig.UnmarshalProto(configProtoBytes)
	require.NoError(t, err, "should be able to unmarshal new capability config from proto bytes")

	// DON capability configurations should include new capability config
	don, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)
	var cfgFound bool
	for _, cfg := range don.CapabilityConfigurations {
		if cfg.CapabilityId == capID {
			got := new(pkg.CapabilityConfig)
			require.NoError(t, got.UnmarshalProto(cfg.Config), "unmarshal capability config proto bytes should not error")
			if diff := cmp.Diff(expectedConfig, got, protocmp.Transform()); diff != "" {
				t.Errorf("capability config proto bytes should match: %s", diff)
			}

			cfgFound = true
		}
	}
	require.True(t, cfgFound, "expected don to have %s capability configuration", capID)
}

func TestAddCapabilities_Apply(t *testing.T) {
	// SetupEnvV2 deploys a cap reg v2 and configures it. So no need to do that here, just leverage the existing one.
	fixture := test.SetupEnvV2(t, false)

	addNewCapability(t, fixture, newCapID)
	requireCapability(t, fixture, newCapID)

	// add another capability and ensure that both are present
	addNewCapability(t, fixture, anotherCapID)
	requireCapability(t, fixture, newCapID)
	requireCapability(t, fixture, anotherCapID)
}

func TestAddCapabilities_Apply_MCMS(t *testing.T) {
	// SetupEnvV2 deploys a cap reg v2 and configures it. So no need to do that here, just leverage the existing one.
	fixture := test.SetupEnvV2(t, true)

	input := changeset.AddCapabilitiesInput{
		RegistryChainSel:  fixture.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			test.DONName: {{
				Capability: contracts.Capability{
					CapabilityID:          newCapID,
					ConfigurationContract: common.Address{},
					Metadata:              newCapMetadata,
				},
				Config: newCapConfig,
			}},
		},
		Force: true,
		MCMSConfig: &crecontracts.MCMSConfig{
			MinDelay: 1 * time.Second,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		},
	}

	// Preconditions
	err := changeset.AddCapabilities{}.VerifyPreconditions(*fixture.Env, input)
	require.NoError(t, err)

	// Apply
	csOut, err := changeset.AddCapabilities{}.Apply(*fixture.Env, input)
	require.NoError(t, err)

	// Verify the changeset output
	require.NotNil(t, csOut.Reports, "reports should be present")
	require.NotEmpty(t, csOut.MCMSTimelockProposals, "should have MCMS proposals when using MCMS")
}

func aptosTestCapabilityID(aptosChainSelector uint64) string {
	return fmt.Sprintf("aptos:ChainSelector:%d@1.0.0", aptosChainSelector)
}

func addCapabilityWithModifier(t *testing.T, fixture *test.EnvWrapperV2) {
	t.Helper()
	require.NotNil(t, fixture.Env.Offchain, "Aptos add-capabilities needs JD Offchain client")

	capID := aptosTestCapabilityID(fixture.AptosSelector)
	input := changeset.AddCapabilitiesInput{
		RegistryChainSel:  fixture.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			test.DONName: {{
				Capability: contracts.Capability{
					CapabilityID:          capID,
					ConfigurationContract: common.Address{},
					Metadata:              newCapMetadata,
				},
				Config: maps.Clone(newCapConfig),
			}},
		},
		Force: true,
	}

	require.NoError(t, changeset.AddCapabilities{}.VerifyPreconditions(*fixture.Env, input))
	_, err := changeset.AddCapabilities{}.Apply(*fixture.Env, input)
	require.NoError(t, err)
}

func requireCapabilityWithModifier(t *testing.T, fixture *test.EnvWrapperV2) {
	t.Helper()

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		fixture.RegistryAddress,
		fixture.Env.BlockChains.EVMChains()[fixture.RegistrySelector].Client,
	)
	require.NoError(t, err)

	capID := aptosTestCapabilityID(fixture.AptosSelector)
	caps, err := pkg.GetCapabilities(nil, capReg)
	require.NoError(t, err)
	var foundCap bool
	for _, c := range caps {
		if c.CapabilityId == capID {
			foundCap = true
			break
		}
	}
	require.True(t, foundCap, "aptos capability %s should be registered", capID)

	don, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)

	var cfgFound bool
	for _, cfg := range don.CapabilityConfigurations {
		if cfg.CapabilityId == capID {
			got := new(pkg.CapabilityConfig)
			require.NoError(t, got.UnmarshalProto(cfg.Config))
			requireAptosSpecP2PTransmitterMap(t, got)
			cfgFound = true
			break
		}
	}
	require.True(t, cfgFound, "expected don to have %s capability configuration", capID)
}

// requireAptosSpecP2PTransmitterMap checks UnmarshalProto output: specConfig (values.v1.Map)
// contains p2pToTransmitterMap with a non-empty nested map of entries.
func requireAptosSpecP2PTransmitterMap(t *testing.T, cfg *pkg.CapabilityConfig) {
	t.Helper()
	spec, ok := (*cfg)["specConfig"].(map[string]any)
	require.True(t, ok, "specConfig should be present as object")
	fields, ok := spec["fields"].(map[string]any)
	require.True(t, ok, "specConfig should have values.v1.Map fields")
	const p2pKey = "p2pToTransmitterMap"
	raw, ok := fields[p2pKey]
	require.True(t, ok, "specConfig.fields should contain %q", p2pKey)
	p2pVal, ok := raw.(map[string]any)
	require.True(t, ok, "%q should be an object", p2pKey)
	mv, ok := p2pVal["mapValue"].(map[string]any)
	require.True(t, ok, "%q should be a values map (mapValue)", p2pKey)
	inner, ok := mv["fields"].(map[string]any)
	require.True(t, ok, "%q.mapValue should have fields", p2pKey)
	require.NotEmpty(t, inner, "%q should have at least one peer→transmitter entry", p2pKey)
}

func TestAddCapabilities_Apply_Modifier(t *testing.T) {
	fixture := test.SetupEnvV2(t, false)
	addCapabilityWithModifier(t, fixture)
	requireCapabilityWithModifier(t, fixture)
}
