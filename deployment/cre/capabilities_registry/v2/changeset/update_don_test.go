package changeset_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

// Local constants (same values used in existing tests)
const (
	csaKey              = "4240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c258"
	signer1             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c251"
	signer2             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c252"
	p2pID1              = "p2p_12D3KooWM1111111111111111111111111111111111111111111"
	p2pID2              = "p2p_12D3KooWM1111111111111111111111111111111111111111112"
	encryptionPublicKey = "7240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c254"
)

type updFixture struct {
	env        cldf.Environment
	selector   uint64
	qualifier  string
	address    string
	registry   *capabilities_registry_v2.CapabilitiesRegistry
	donName    string
	capIDs     []string
	isWorkflow bool // whether initial DON.AcceptsWorkflows = true
}

func setupRegistryForUpdateDON(t *testing.T, isWorkflow bool) *updFixture {
	t.Helper()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	qualifier := "update-don-changeset-tests"

	deployTask := runtime.ChangesetTask(changeset.DeployCapabilitiesRegistry{}, changeset.DeployCapabilitiesRegistryInput{
		ChainSelector: selector,
		Qualifier:     qualifier,
	})
	require.NoError(t, rt.Exec(deployTask))
	deployOutput := rt.State().Outputs[deployTask.ID()]
	require.NotNil(t, deployOutput)

	addr := deployOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(qualifier))[0].Address

	reg, err := capabilities_registry_v2.NewCapabilitiesRegistry(common.HexToAddress(addr), rt.Environment().BlockChains.EVMChains()[selector].Client)
	require.NoError(t, err)

	writeChain := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`),
	}
	var writeChainMeta map[string]any
	require.NoError(t, json.Unmarshal(writeChain.Metadata, &writeChainMeta))

	trigger := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "trigger@1.0.0",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 1, "responseType": 1}`),
	}
	var triggerMeta map[string]any
	require.NoError(t, json.Unmarshal(trigger.Metadata, &triggerMeta))

	nop1 := "test-nop-1"
	nop2 := "test-nop-2"
	nodes := []changeset.CapabilitiesRegistryNodeParams{
		{
			NOP:                 nop1,
			Signer:              signer1,
			P2pID:               p2pID1,
			EncryptionPublicKey: encryptionPublicKey,
			CsaKey:              csaKey,
			CapabilityIDs:       []string{writeChain.CapabilityId, trigger.CapabilityId},
		},
		{
			NOP:                 nop2,
			Signer:              signer2,
			P2pID:               p2pID2,
			EncryptionPublicKey: encryptionPublicKey,
			CsaKey:              csaKey,
			CapabilityIDs:       []string{writeChain.CapabilityId, trigger.CapabilityId},
		},
	}
	nodeSet := []string{p2pID1, p2pID2}

	// Initial DON config (workflow type used for both variants)
	cfg := map[string]any{
		"defaultConfig": map[string]any{},
		"remoteTriggerConfig": map[string]any{
			"registrationRefresh":     "20s",
			"registrationExpiry":      "60s",
			"minResponsesToAggregate": 2,
			"messageExpiry":           "120s",
		},
	}
	donName := "upd-don-v2"

	// Register everything using ConfigureCapabilitiesRegistry (no MCMS)
	_, err = changeset.ConfigureCapabilitiesRegistry{}.Apply(rt.Environment(), changeset.ConfigureCapabilitiesRegistryInput{
		ChainSelector:               selector,
		CapabilitiesRegistryAddress: addr,
		Nops: []changeset.CapabilitiesRegistryNodeOperator{
			{Admin: common.HexToAddress("0x01"), Name: nop1},
			{Admin: common.HexToAddress("0x02"), Name: nop2},
		},
		Capabilities: []changeset.CapabilitiesRegistryCapability{
			{CapabilityID: writeChain.CapabilityId, Metadata: writeChainMeta},
			{CapabilityID: trigger.CapabilityId, Metadata: triggerMeta},
		},
		Nodes: nodes,
		DONs: []changeset.CapabilitiesRegistryNewDONParams{
			{
				Name:        donName,
				DonFamilies: []string{"upd-family"},
				Config: map[string]any{
					"defaultConfig": map[string]any{},
				},
				CapabilityConfigurations: []changeset.CapabilitiesRegistryCapabilityConfiguration{
					{CapabilityID: writeChain.CapabilityId, Config: cfg},
				},
				Nodes:            nodeSet,
				F:                1,
				IsPublic:         true,
				AcceptsWorkflows: isWorkflow,
			},
		},
	})
	require.NoError(t, err)

	return &updFixture{
		env:        rt.Environment(),
		selector:   selector,
		qualifier:  qualifier,
		address:    addr,
		registry:   reg,
		donName:    donName,
		capIDs:     []string{writeChain.CapabilityId, trigger.CapabilityId},
		isWorkflow: isWorkflow,
	}
}

// Happy path: non-workflow DON; also renames the DON; capability config updated; visibility/F preserved.
func TestUpdateDONChangeset_ByName_Direct_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, false)

	// New config to apply
	newCfg := map[string]any{
		"defaultConfig": map[string]any{},
		"remoteTriggerConfig": map[string]any{
			"registrationRefresh":     "25s", // changed value to detect update
			"registrationExpiry":      "60s",
			"minResponsesToAggregate": 2,
			"messageExpiry":           "120s",
		},
	}
	wantProto, err := pkg.CapabilityConfig(newCfg).MarshalProto()
	require.NoError(t, err)

	newName := fx.donName + "-renamed"

	out, err := changeset.UpdateDON{}.Apply(fx.env, changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName, // required current name
		NewDonName:        newName,    // rename the DON
		CapabilityConfigs: []contracts.CapabilityConfig{
			{Capability: contracts.Capability{CapabilityID: fx.capIDs[0]}, Config: newCfg},
		},
		Force:      false,
		MCMSConfig: nil,
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Empty(t, out.MCMSTimelockProposals, "no MCMS → proposals must be empty")
	require.NotEmpty(t, out.Reports)

	// Old name should no longer resolve
	_, err = fx.registry.GetDONByName(nil, fx.donName)
	require.Error(t, err)

	// Verify on-chain state under the new name
	got, err := fx.registry.GetDONByName(nil, newName)
	require.NoError(t, err)
	assert.Equal(t, uint8(1), got.F, "F is preserved from existing DON")
	assert.True(t, got.IsPublic, "IsPublic preserved (we did not toggle it in changeset)")
	require.Len(t, got.CapabilityConfigurations, 1)
	assert.Equal(t, fx.capIDs[0], got.CapabilityConfigurations[0].CapabilityId)
	assert.Equal(t, wantProto, got.CapabilityConfigurations[0].Config)
}

// Safety gate: workflow DON should refuse without Force=true (changeset passes Force through to operation).
func TestUpdateDONChangeset_ByName_Workflow_RefusesWithoutForce(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, true)

	_, err := changeset.UpdateDON{}.Apply(fx.env, changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName, // required
		CapabilityConfigs: []contracts.CapabilityConfig{
			{Capability: contracts.Capability{CapabilityID: fx.capIDs[0]}, Config: map[string]any{"defaultConfig": map[string]any{}}},
		},
		Force: false,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "refusing to update workflow don")
}

// Force override: workflow DON update succeeds when Force=true.
func TestUpdateDONChangeset_ByName_Workflow_Force_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, true)

	// Use a valid protobuf structure with proper fields format
	newCfg := map[string]any{
		"defaultConfig": map[string]any{
			"fields": map[string]any{
				"testKey": map[string]any{
					"stringValue": "testValue",
				},
			},
		},
	}
	wantProto, err := pkg.CapabilityConfig(newCfg).MarshalProto()
	require.NoError(t, err)

	out, err := changeset.UpdateDON{}.Apply(fx.env, changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName, // required
		CapabilityConfigs: []contracts.CapabilityConfig{
			{Capability: contracts.Capability{CapabilityID: fx.capIDs[0]}, Config: newCfg},
		},
		Force: true, // override
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Empty(t, out.MCMSTimelockProposals)

	got, err := fx.registry.GetDONByName(nil, fx.donName)
	require.NoError(t, err)
	require.Len(t, got.CapabilityConfigurations, 1)
	assert.Equal(t, fx.capIDs[0], got.CapabilityConfigurations[0].CapabilityId)
	assert.Equal(t, wantProto, got.CapabilityConfigurations[0].Config)
}

// VerifyPreconditions: empty Name is rejected.
// NOTE: current implementation returns "must provide a non-empty DONName"
func TestUpdateDONChangeset_VerifyPreconditions_EmptyName(t *testing.T) {
	t.Parallel()
	var cs changeset.UpdateDON
	err := cs.VerifyPreconditions(cldf.Environment{}, changeset.UpdateDONInput{
		RegistryQualifier: "q",
		RegistryChainSel:  1,
		DONName:           "", // invalid
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide a non-empty DONName")
}

// Chain not found: Apply should fail early with a clear message.
func TestUpdateDONChangeset_ByName_ChainNotFound(t *testing.T) {
	t.Parallel()

	// Env with no chains (or use a selector not present in env)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	_, err = changeset.UpdateDON{}.Apply(rt.Environment(), changeset.UpdateDONInput{
		RegistryQualifier: "anything",
		RegistryChainSel:  0, // invalid selector for this env
		DONName:           "some-don",
		CapabilityConfigs: nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain not found for selector")
}

// Qualifier not found in DataStore: Apply should fail when it cannot look up the registry address.
func TestUpdateDONChangeset_ByName_QualifierNotFound(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// No deployment under this qualifier → address lookup will fail
	_, err = changeset.UpdateDON{}.Apply(rt.Environment(), changeset.UpdateDONInput{
		RegistryQualifier: "missing-qualifier",
		RegistryChainSel:  selector,
		DONName:           "some-don",
		CapabilityConfigs: nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get registry address")
}
