package changeset_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"

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
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// Local constants (same values used in existing tests)
const (
	csaKey              = "4240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c258"
	signer1             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c251"
	signer2             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c252"
	p2pID1              = "p2p_12D3KooWM1111111111111111111111111111111111111111111"
	p2pID2              = "p2p_12D3KooWM1111111111111111111111111111111111111111112"
	p2pID3              = "p2p_12D3KooWM1111111111111111111111111111111111111111113"
	signer3             = "5240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c253"
	encryptionPublicKey = "7240b57854dd1f21c10353ea458eecd8593624d0e0a7cca07c62a4b58df8c254"
)

type updFixture struct {
	rt         *runtime.Runtime
	selector   uint64
	qualifier  string
	address    string
	registry   *capabilities_registry_v2.CapabilitiesRegistry
	donName    string
	capIDs     []string
	isWorkflow bool // whether initial DON.AcceptsWorkflows = true
}

func setupRegistryForUpdateDON(t *testing.T, isWorkflow, useMCMS bool) *updFixture {
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

	// this tests if modifier scaffolding works in the changeset
	// each modifier's code is not tested here
	capWithModifierID := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "aptos:ChainSelector:743186221051783445@1.0.0",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`),
	}
	var aptosWriteChainMeta map[string]any
	require.NoError(t, json.Unmarshal(capWithModifierID.Metadata, &aptosWriteChainMeta))

	nop1 := "test-nop-1"
	nop2 := "test-nop-2"
	nodes := []changeset.CapabilitiesRegistryNodeParams{
		{
			NOP:                 nop1,
			Signer:              signer1,
			P2pID:               p2pID1,
			EncryptionPublicKey: encryptionPublicKey,
			CsaKey:              csaKey,
			CapabilityIDs:       []string{writeChain.CapabilityId, trigger.CapabilityId, capWithModifierID.CapabilityId},
		},
		{
			NOP:                 nop2,
			Signer:              signer2,
			P2pID:               p2pID2,
			EncryptionPublicKey: encryptionPublicKey,
			CsaKey:              csaKey,
			CapabilityIDs:       []string{writeChain.CapabilityId, trigger.CapabilityId, capWithModifierID.CapabilityId},
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
	err = rt.Exec(
		runtime.ChangesetTask(changeset.ConfigureCapabilitiesRegistry{}, changeset.ConfigureCapabilitiesRegistryInput{
			ChainSelector:               selector,
			CapabilitiesRegistryAddress: addr,
			Nops: []changeset.CapabilitiesRegistryNodeOperator{
				{Admin: common.HexToAddress("0x01"), Name: nop1},
				{Admin: common.HexToAddress("0x02"), Name: nop2},
			},
			Capabilities: []changeset.CapabilitiesRegistryCapability{
				{CapabilityID: writeChain.CapabilityId, Metadata: writeChainMeta},
				{CapabilityID: trigger.CapabilityId, Metadata: triggerMeta},
				{CapabilityID: capWithModifierID.CapabilityId, Metadata: aptosWriteChainMeta},
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
						{CapabilityID: capWithModifierID.CapabilityId, Config: cfg},
					},
					Nodes:            nodeSet,
					F:                1,
					IsPublic:         true,
					AcceptsWorkflows: isWorkflow,
				},
			},
		}),
	)
	require.NoError(t, err)

	if !useMCMS {
		return &updFixture{
			rt:         rt,
			selector:   selector,
			qualifier:  qualifier,
			address:    addr,
			registry:   reg,
			donName:    donName,
			capIDs:     []string{writeChain.CapabilityId, trigger.CapabilityId},
			isWorkflow: isWorkflow,
		}
	}

	timelockCfgs := map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
		selector: cldftesthelpers.SingleGroupTimelockConfig(t),
	}

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2), timelockCfgs),
	)
	require.NoError(t, err, "failed to deploy MCMS infrastructure")

	t.Log("MCMS infrastructure deployed successfully")

	t.Log("Transferring ownership to MCMS...")
	err = rt.Exec(
		runtime.ChangesetTask(
			cldf.CreateLegacyChangeSet(keystonechangeset.AcceptAllOwnershipsProposal),
			&keystonechangeset.AcceptAllOwnershipRequest{
				ChainSelector: selector,
				MinDelay:      0,
			},
		),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	)
	require.NoError(t, err, "failed to transfer ownership to MCMS")
	t.Log("Ownership transferred to MCMS successfully")

	return &updFixture{
		rt:         rt,
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
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, false, false)

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

	task := runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
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

	err = fx.rt.Exec(task)
	require.NoError(t, err)

	out := fx.rt.State().Outputs[task.ID()]
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

func TestUpdateDONChangeset_ByName_Direct_Succeeds_MCMS(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, false, true)

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

	newName := fx.donName + "-renamed"

	task := runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName, // required current name
		NewDonName:        newName,    // rename the DON
		CapabilityConfigs: []contracts.CapabilityConfig{
			{Capability: contracts.Capability{CapabilityID: fx.capIDs[0]}, Config: newCfg},
		},
		Force: false,
		MCMSConfig: &crecontracts.MCMSConfig{
			MinDelay: 1 * time.Second,
			TimelockQualifierPerChain: map[uint64]string{
				fx.selector: "",
			},
		},
	})

	err := fx.rt.Exec(task)
	require.NoError(t, err)

	out := fx.rt.State().Outputs[task.ID()]
	assert.NotNil(t, out)
	assert.NotEmpty(t, out.Reports)
	assert.NotEmpty(t, out.MCMSTimelockProposals, "MCMS → proposals must not be empty")
}

// Force override: workflow DON update succeeds when Force=true.
func TestUpdateDONChangeset_ByName_Workflow_Force_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t /*isWorkflow=*/, true, false)

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

	task := runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName, // required
		CapabilityConfigs: []contracts.CapabilityConfig{
			{Capability: contracts.Capability{CapabilityID: fx.capIDs[0]}, Config: newCfg},
		},
		Force: true, // override
	})
	require.NoError(t, err)

	err = fx.rt.Exec(task)
	require.NoError(t, err)

	out := fx.rt.State().Outputs[task.ID()]
	require.NotNil(t, out)
	assert.Empty(t, out.MCMSTimelockProposals, "no MCMS → proposals must be empty")
	require.NotEmpty(t, out.Reports)

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
	require.ErrorContains(t, err, "must provide a non-empty DONName")
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
	require.ErrorContains(t, err, "chain not found for selector")
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
	require.ErrorContains(t, err, "failed to get registry address")
}

// TestUpdateDON_FirstOCR3ConfigCapabilities exercises the first-config gate in
// processOCR3Configs by going through UpdateDON.Apply with ocr3Configs entries.
func TestUpdateDON_FirstOCR3ConfigCapabilities(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t, false, false)
	env := fx.rt.Environment()

	// A minimal ocr3Configs entry matching the proto structure: oracle config
	// params are nested under the "offchainConfig" key.
	ocr3Entry := map[string]any{
		"offchainConfig": map[string]any{
			"deltaProgressMillis":  5000,
			"maxFaultyOracles":     1,
			"transmissionSchedule": []any{2},
		},
	}

	makeInput := func(capID string, firstCaps []string) changeset.UpdateDONInput {
		return changeset.UpdateDONInput{
			RegistryQualifier: fx.qualifier,
			RegistryChainSel:  fx.selector,
			DONName:           fx.donName,
			CapabilityConfigs: []contracts.CapabilityConfig{
				{
					Capability: contracts.Capability{CapabilityID: capID},
					Config: map[string]any{
						"defaultConfig": map[string]any{},
						"ocr3Configs": map[string]any{
							"__default__": ocr3Entry,
						},
					},
				},
			},
			FirstOCR3ConfigCapabilities: firstCaps,
		}
	}

	t.Run("rejects when capability is not in FirstOCR3ConfigCapabilities", func(t *testing.T) {
		_, err := changeset.UpdateDON{}.Apply(env, makeInput(fx.capIDs[0], nil))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "firstOCR3ConfigCapabilities")
	})

	t.Run("passes first-config check when capability is listed", func(t *testing.T) {
		_, err := changeset.UpdateDON{}.Apply(env, makeInput(fx.capIDs[0], []string{fx.capIDs[0]}))
		require.Error(t, err)
		// Should fail later (e.g. ComputeOCR3Config), NOT at the first-config gate.
		assert.NotContains(t, err.Error(), "firstOCR3ConfigCapabilities")
	})

	t.Run("rejects unlisted capability even when another is listed", func(t *testing.T) {
		// List capIDs[0] but provide ocr3Configs for capIDs[1] → should fail
		_, err := changeset.UpdateDON{}.Apply(env, makeInput(fx.capIDs[1], []string{fx.capIDs[0]}))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "firstOCR3ConfigCapabilities")
	})
}

// setupRegistryWith3Nodes creates a registry with 3 NOPs and 3 nodes but a DON that
// initially contains only p2pID1 and p2pID2 (F=1). This lets tests verify that
// UpdateDON can expand the node set to include the third node.
func setupRegistryWith3Nodes(t *testing.T) *updFixture {
	t.Helper()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	qualifier := "update-don-3nodes-tests"

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
		CapabilityId: "write-chain@1.0.1",
		Metadata:     []byte(`{"capabilityType": 3, "responseType": 1}`),
	}
	var writeChainMeta map[string]any
	require.NoError(t, json.Unmarshal(writeChain.Metadata, &writeChainMeta))

	nop1, nop2, nop3 := "test-nop-1", "test-nop-2", "test-nop-3"
	donName := "upd-don-3nodes"

	err = rt.Exec(
		runtime.ChangesetTask(changeset.ConfigureCapabilitiesRegistry{}, changeset.ConfigureCapabilitiesRegistryInput{
			ChainSelector:               selector,
			CapabilitiesRegistryAddress: addr,
			Nops: []changeset.CapabilitiesRegistryNodeOperator{
				{Admin: common.HexToAddress("0x01"), Name: nop1},
				{Admin: common.HexToAddress("0x02"), Name: nop2},
				{Admin: common.HexToAddress("0x03"), Name: nop3},
			},
			Capabilities: []changeset.CapabilitiesRegistryCapability{
				{CapabilityID: writeChain.CapabilityId, Metadata: writeChainMeta},
			},
			Nodes: []changeset.CapabilitiesRegistryNodeParams{
				{
					NOP:                 nop1,
					Signer:              signer1,
					P2pID:               p2pID1,
					EncryptionPublicKey: encryptionPublicKey,
					CsaKey:              csaKey,
					CapabilityIDs:       []string{writeChain.CapabilityId},
				},
				{
					NOP:                 nop2,
					Signer:              signer2,
					P2pID:               p2pID2,
					EncryptionPublicKey: encryptionPublicKey,
					CsaKey:              csaKey,
					CapabilityIDs:       []string{writeChain.CapabilityId},
				},
				{
					NOP:                 nop3,
					Signer:              signer3,
					P2pID:               p2pID3,
					EncryptionPublicKey: encryptionPublicKey,
					CsaKey:              csaKey,
					CapabilityIDs:       []string{writeChain.CapabilityId},
				},
			},
			DONs: []changeset.CapabilitiesRegistryNewDONParams{
				{
					Name: donName,
					Config: map[string]any{
						"defaultConfig": map[string]any{},
					},
					CapabilityConfigurations: []changeset.CapabilitiesRegistryCapabilityConfiguration{
						{CapabilityID: writeChain.CapabilityId, Config: map[string]any{"defaultConfig": map[string]any{}}},
					},
					Nodes:    []string{p2pID1, p2pID2},
					F:        1,
					IsPublic: true,
				},
			},
		}),
	)
	require.NoError(t, err)

	return &updFixture{
		rt:        rt,
		selector:  selector,
		qualifier: qualifier,
		address:   addr,
		registry:  reg,
		donName:   donName,
		capIDs:    []string{writeChain.CapabilityId},
	}
}

func TestUpdateDONChangeset_WithF_OverridesOnChainValue(t *testing.T) {
	t.Parallel()
	// Use the 3-node fixture (F=1 initially). F=2 is valid with 3 nodes (F < N).
	fx := setupRegistryWith3Nodes(t)

	err := fx.rt.Exec(runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier:                 fx.qualifier,
		RegistryChainSel:                  fx.selector,
		DONName:                           fx.donName,
		MergeCapabilityConfigsWithOnChain: true,
		Nodes:                             []string{p2pID1, p2pID2, p2pID3}, // keep 3 nodes so F=2 is valid
		F:                                 2,
	}))
	require.NoError(t, err)

	got, err := fx.registry.GetDONByName(nil, fx.donName)
	require.NoError(t, err)
	assert.Equal(t, uint8(2), got.F)
}

func TestUpdateDONChangeset_WithF_Zero_PreservesExistingF(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t, false, false)

	err := fx.rt.Exec(runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier:                 fx.qualifier,
		RegistryChainSel:                  fx.selector,
		DONName:                           fx.donName,
		MergeCapabilityConfigsWithOnChain: true,
		F:                                 0, // zero value — keep on-chain F
	}))
	require.NoError(t, err)

	got, err := fx.registry.GetDONByName(nil, fx.donName)
	require.NoError(t, err)
	assert.Equal(t, uint8(1), got.F, "F must remain unchanged when input F=0")
}

func TestUpdateDONChangeset_WithNodes_OverridesOnChainNodes(t *testing.T) {
	t.Parallel()
	fx := setupRegistryWith3Nodes(t)

	// DON currently has 2 nodes (p2pID1, p2pID2); expand to all 3.
	err := fx.rt.Exec(runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier:                 fx.qualifier,
		RegistryChainSel:                  fx.selector,
		DONName:                           fx.donName,
		MergeCapabilityConfigsWithOnChain: true,
		Nodes:                             []string{p2pID1, p2pID2, p2pID3},
	}))
	require.NoError(t, err)

	got, err := fx.registry.GetDONByName(nil, fx.donName)
	require.NoError(t, err)
	// The registry stores hashed P2P IDs; membership count is the authoritative assertion.
	assert.Len(t, got.NodeP2PIds, 3, "DON must now have 3 members")
}

func TestUpdateDONChangeset_WithNodes_Empty_PreservesExistingNodes(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t, false, false)

	err := fx.rt.Exec(runtime.ChangesetTask(changeset.UpdateDON{}, changeset.UpdateDONInput{
		RegistryQualifier:                 fx.qualifier,
		RegistryChainSel:                  fx.selector,
		DONName:                           fx.donName,
		MergeCapabilityConfigsWithOnChain: true,
		Nodes:                             nil, // empty — keep on-chain nodes
	}))
	require.NoError(t, err)

	got, err := fx.registry.GetDONByName(nil, fx.donName)
	require.NoError(t, err)
	assert.Len(t, got.NodeP2PIds, 2, "node set must remain unchanged when Nodes is nil")
}

func TestUpdateDONChangeset_WithNodes_InvalidP2PID(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateDON(t, false, false)

	_, err := changeset.UpdateDON{}.Apply(fx.rt.Environment(), changeset.UpdateDONInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DONName:           fx.donName,
		Nodes:             []string{"not-a-valid-p2p-id"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid P2P ID")
}
