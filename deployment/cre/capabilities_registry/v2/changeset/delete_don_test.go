package changeset_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	opscontracts "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
)

type delFixture struct {
	env       cldf.Environment
	selector  uint64
	qualifier string
	address   string
	registry  *capabilities_registry_v2.CapabilitiesRegistry
	donNames  []string
}

func setupRegistryForDeleteDON(t *testing.T, secondDON bool) *delFixture {
	t.Helper()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	qualifier := "delete-don-changeset-tests"

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

	// Capabilities
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

	cfg := map[string]any{
		"defaultConfig": map[string]any{},
	}
	don1 := "del-don-1"
	don2 := "del-don-2"

	// Seed registry with one or two DONs
	params := []changeset.CapabilitiesRegistryNewDONParams{
		{
			Name:        don1,
			DonFamilies: []string{"del-family"},
			Config:      map[string]any{"defaultConfig": map[string]any{}},
			CapabilityConfigurations: []changeset.CapabilitiesRegistryCapabilityConfiguration{
				{CapabilityID: writeChain.CapabilityId, Config: cfg},
			},
			Nodes:            nodeSet,
			F:                1,
			IsPublic:         true,
			AcceptsWorkflows: false,
		},
	}
	if secondDON {
		params = append(params, changeset.CapabilitiesRegistryNewDONParams{
			Name:        don2,
			DonFamilies: []string{"del-family"},
			Config:      map[string]any{"defaultConfig": map[string]any{}},
			CapabilityConfigurations: []changeset.CapabilitiesRegistryCapabilityConfiguration{
				{CapabilityID: writeChain.CapabilityId, Config: cfg},
			},
			Nodes:            nodeSet,
			F:                1,
			IsPublic:         true,
			AcceptsWorkflows: false,
		})
	}

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
		DONs:  params,
	})
	require.NoError(t, err)

	dons := []string{don1}
	if secondDON {
		dons = append(dons, don2)
	}

	return &delFixture{
		env:       rt.Environment(),
		selector:  selector,
		qualifier: qualifier,
		address:   addr,
		registry:  reg,
		donNames:  dons,
	}
}

func TestDeleteDONChangeset_ByName_Direct_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForDeleteDON(t, false)

	_, err := fx.registry.GetDONByName(nil, fx.donNames[0])
	require.NoError(t, err)

	out, err := changeset.DeleteDONs{}.Apply(fx.env, changeset.DeleteDONsInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DonNames:          []string{fx.donNames[0]},
		MCMSConfig:        nil,
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Empty(t, out.MCMSTimelockProposals)

	// Should be gone now
	_, err = fx.registry.GetDONByName(nil, fx.donNames[0])
	require.Error(t, err)
}

func TestDeleteDONChangeset_ByNames_Multi_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForDeleteDON(t, true)

	// Both exist
	for _, name := range fx.donNames {
		_, err := fx.registry.GetDONByName(nil, name)
		require.NoError(t, err)
	}

	out, err := changeset.DeleteDONs{}.Apply(fx.env, changeset.DeleteDONsInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		DonNames:          fx.donNames,
	})
	require.NoError(t, err)
	assert.Empty(t, out.MCMSTimelockProposals)

	// Both should be gone
	for _, name := range fx.donNames {
		_, err := fx.registry.GetDONByName(nil, name)
		require.Error(t, err)
	}
}

func TestDeleteDONChangeset_VerifyPreconditions_EmptyList(t *testing.T) {
	t.Parallel()
	var cs changeset.DeleteDONs
	err := cs.VerifyPreconditions(cldf.Environment{}, changeset.DeleteDONsInput{
		RegistryQualifier: "q",
		RegistryChainSel:  1,
		DonNames:          nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide at least one DON name")
}

func TestDeleteDONChangeset_Apply_MissingDON_Fails(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	qualifier := "delete-don-missing-tests"

	deployTask := runtime.ChangesetTask(changeset.DeployCapabilitiesRegistry{}, changeset.DeployCapabilitiesRegistryInput{
		ChainSelector: selector,
		Qualifier:     qualifier,
	})
	require.NoError(t, rt.Exec(deployTask))

	_, err = changeset.DeleteDONs{}.Apply(rt.Environment(), changeset.DeleteDONsInput{
		RegistryQualifier: qualifier,
		RegistryChainSel:  selector,
		DonNames:          []string{"does-not-exist"},
	})
	require.Error(t, err)
	// Error wording now comes from cldf.DecodeErr(...) and may vary.
	// Avoid asserting specific substrings to keep this robust across ABI/runtime changes.
}

func TestDeleteDONChangeset_ChainNotFound(t *testing.T) {
	t.Parallel()

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	_, err = changeset.DeleteDONs{}.Apply(rt.Environment(), changeset.DeleteDONsInput{
		RegistryQualifier: "anything",
		RegistryChainSel:  0,
		DonNames:          []string{"x"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain not found for selector")
}

func TestDeleteDONChangeset_QualifierNotFound(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	_, err = changeset.DeleteDONs{}.Apply(rt.Environment(), changeset.DeleteDONsInput{
		RegistryQualifier: "missing-qualifier",
		RegistryChainSel:  selector,
		DonNames:          []string{"some-don"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get registry address")
}

func TestDeleteDON_MCMS_Configuration(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForDeleteDON(t, false)

	_, err := fx.registry.GetDONByName(nil, fx.donNames[0])
	require.NoError(t, err)
	mcmsFixture := setupCapabilitiesRegistryWithMCMS(t)
	mcmsFixture.env.OperationsBundle = operations.NewBundle(mcmsFixture.env.GetContext, mcmsFixture.env.Logger, operations.NewMemoryReporter())

	chain, ok := mcmsFixture.env.BlockChains.EVMChains()[mcmsFixture.chainSelector]
	require.True(t, ok, "chain should be found for selector %d", mcmsFixture.chainSelector)

	reg := fx.registry

	mcmsContracts, err := strategies.GetMCMSContracts(mcmsFixture.env, mcmsFixture.chainSelector, *mcmsFixture.configureInput.MCMSConfig)
	require.NoError(t, err, "should be able to get MCMS contracts")
	require.NotNil(t, mcmsContracts, "MCMS contracts should not be nil")

	realStrategy, err := strategies.CreateStrategy(
		chain,
		mcmsFixture.env,
		mcmsFixture.configureInput.MCMSConfig,
		mcmsContracts,
		common.HexToAddress(mcmsFixture.capabilitiesRegistryAddress),
		"test DeleteDON with MCMS",
	)
	require.NoError(t, err, "should be able to create MCMS strategy")

	deps := opscontracts.DeleteDONDeps{
		Env:                  &mcmsFixture.env,
		Strategy:             realStrategy,
		CapabilitiesRegistry: reg,
	}

	input := opscontracts.DeleteDONInput{
		ChainSelector: mcmsFixture.chainSelector,
		DonNames:      []string{fx.donNames[0]},
		MCMSConfig:    mcmsFixture.configureInput.MCMSConfig,
	}

	// Execute the DeleteDON operation with MCMS; this should CREATE a proposal,
	// not execute the removal immediately.
	report, err := operations.ExecuteOperation(
		mcmsFixture.env.OperationsBundle,
		opscontracts.DeleteDON,
		deps,
		input,
	)
	require.NoError(t, err, "DeleteDON with MCMS should succeed (proposal created)")
	require.NotNil(t, report, "operation report should not be nil")

	// Verify operation content mirrors your NOPs test assertions
	require.NotZero(t, report.Output.Operation, "an operation should have been generated")
	require.NotEmpty(t, report.Output.Operation.Transactions, "operation should have transactions")
	assert.Equal(t, []string{fx.donNames[0]}, report.Output.DeletedNames)

	// Since this is only a proposal (NoSend + nonzero GasLimit), the DON must still exist.
	_, err = reg.GetDONByName(nil, fx.donNames[0])
	require.NoError(t, err, "DON should still exist until governance executes the proposal")
}
