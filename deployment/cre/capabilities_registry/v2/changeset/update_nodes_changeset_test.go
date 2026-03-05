package changeset_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
)

type updNodesFixture struct {
	env       cldf.Environment
	selector  uint64
	qualifier string
	address   string
	registry  *capabilities_registry_v2.CapabilitiesRegistry
	nodes     []changeset.CapabilitiesRegistryNodeParams
	capIDs    []string
}

func setupRegistryForUpdateNodes(t *testing.T) *updNodesFixture {
	t.Helper()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	qualifier := "update-nodes-changeset-tests"

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
	})
	require.NoError(t, err)

	return &updNodesFixture{
		env:       rt.Environment(),
		selector:  selector,
		qualifier: qualifier,
		address:   addr,
		registry:  reg,
		nodes:     nodes,
		capIDs:    []string{writeChain.CapabilityId, trigger.CapabilityId},
	}
}

func TestUpdateNodesChangeset_SignerUpdate_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateNodes(t)

	newSigner := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	out, err := changeset.UpdateNodes{}.Apply(fx.env, changeset.UpdateNodesInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		Nodes: []changeset.CapabilitiesRegistryNodeParams{
			{
				NOP:                 fx.nodes[0].NOP,
				Signer:              newSigner,
				P2pID:               fx.nodes[0].P2pID,
				EncryptionPublicKey: fx.nodes[0].EncryptionPublicKey,
				CsaKey:              fx.nodes[0].CsaKey,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Empty(t, out.MCMSTimelockProposals)
	require.NotEmpty(t, out.Reports)

	// Verify signer was updated on-chain
	p2pBytes, err := p2pkey.MakePeerID(fx.nodes[0].P2pID)
	require.NoError(t, err)
	onChainNodes, err := fx.registry.GetNodesByP2PIds(nil, [][32]byte{p2pBytes})
	require.NoError(t, err)
	require.Len(t, onChainNodes, 1)

	wantSigner, _ := hex.DecodeString(newSigner)
	assert.Equal(t, [32]byte(wantSigner), onChainNodes[0].Signer)

	// Capabilities must be preserved
	assert.ElementsMatch(t, fx.capIDs, onChainNodes[0].CapabilityIds)
}

func TestUpdateNodesChangeset_MultipleNodes_Succeeds(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateNodes(t)

	newSigner1 := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	newSigner2 := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"

	out, err := changeset.UpdateNodes{}.Apply(fx.env, changeset.UpdateNodesInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		Nodes: []changeset.CapabilitiesRegistryNodeParams{
			{
				NOP:                 fx.nodes[0].NOP,
				Signer:              newSigner1,
				P2pID:               fx.nodes[0].P2pID,
				EncryptionPublicKey: fx.nodes[0].EncryptionPublicKey,
				CsaKey:              fx.nodes[0].CsaKey,
			},
			{
				NOP:                 fx.nodes[1].NOP,
				Signer:              newSigner2,
				P2pID:               fx.nodes[1].P2pID,
				EncryptionPublicKey: fx.nodes[1].EncryptionPublicKey,
				CsaKey:              fx.nodes[1].CsaKey,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, out)

	p2p1, err := p2pkey.MakePeerID(fx.nodes[0].P2pID)
	require.NoError(t, err)
	p2p2, err := p2pkey.MakePeerID(fx.nodes[1].P2pID)
	require.NoError(t, err)

	onChain, err := fx.registry.GetNodesByP2PIds(nil, [][32]byte{p2p1, p2p2})
	require.NoError(t, err)
	require.Len(t, onChain, 2)

	want1, _ := hex.DecodeString(newSigner1)
	want2, _ := hex.DecodeString(newSigner2)

	signers := map[[32]byte]bool{
		[32]byte(want1): false,
		[32]byte(want2): false,
	}
	for _, n := range onChain {
		signers[n.Signer] = true
		assert.ElementsMatch(t, fx.capIDs, n.CapabilityIds)
	}
	for s, found := range signers {
		assert.True(t, found, "signer %x not found on-chain", s)
	}
}

func TestUpdateNodesChangeset_VerifyPreconditions_EmptyNodes(t *testing.T) {
	t.Parallel()
	var cs changeset.UpdateNodes
	err := cs.VerifyPreconditions(cldf.Environment{}, changeset.UpdateNodesInput{
		RegistryQualifier: "q",
		RegistryChainSel:  1,
		Nodes:             nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nodes list cannot be empty")
}

func TestUpdateNodesChangeset_VerifyPreconditions_MissingFields(t *testing.T) {
	t.Parallel()

	full := changeset.CapabilitiesRegistryNodeParams{
		P2pID: p2pID1, NOP: "nop", Signer: signer1,
		EncryptionPublicKey: encryptionPublicKey, CsaKey: csaKey,
	}

	tests := []struct {
		name    string
		mutate  func(*changeset.CapabilitiesRegistryNodeParams)
		wantErr string
	}{
		{"empty P2pID", func(n *changeset.CapabilitiesRegistryNodeParams) { n.P2pID = "" }, "empty P2pID"},
		{"empty NOP", func(n *changeset.CapabilitiesRegistryNodeParams) { n.NOP = "" }, "empty NOP"},
		{"empty Signer", func(n *changeset.CapabilitiesRegistryNodeParams) { n.Signer = "" }, "empty Signer"},
		{"empty EncryptionPublicKey", func(n *changeset.CapabilitiesRegistryNodeParams) { n.EncryptionPublicKey = "" }, "empty EncryptionPublicKey"},
		{"empty CsaKey", func(n *changeset.CapabilitiesRegistryNodeParams) { n.CsaKey = "" }, "empty CsaKey"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			node := full
			tc.mutate(&node)
			var cs changeset.UpdateNodes
			err := cs.VerifyPreconditions(cldf.Environment{}, changeset.UpdateNodesInput{
				RegistryQualifier: "q",
				RegistryChainSel:  1,
				Nodes:             []changeset.CapabilitiesRegistryNodeParams{node},
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestUpdateNodesChangeset_ChainNotFound(t *testing.T) {
	t.Parallel()

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	_, err = changeset.UpdateNodes{}.Apply(rt.Environment(), changeset.UpdateNodesInput{
		RegistryQualifier: "anything",
		RegistryChainSel:  0,
		Nodes: []changeset.CapabilitiesRegistryNodeParams{
			{P2pID: p2pID1, NOP: "nop", Signer: signer1, EncryptionPublicKey: encryptionPublicKey, CsaKey: csaKey},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain not found for selector")
}

func TestUpdateNodesChangeset_QualifierNotFound(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	_, err = changeset.UpdateNodes{}.Apply(rt.Environment(), changeset.UpdateNodesInput{
		RegistryQualifier: "missing-qualifier",
		RegistryChainSel:  selector,
		Nodes: []changeset.CapabilitiesRegistryNodeParams{
			{P2pID: p2pID1, NOP: "nop", Signer: signer1, EncryptionPublicKey: encryptionPublicKey, CsaKey: csaKey},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get registry address")
}

func TestUpdateNodesChangeset_UnknownNOP_Fails(t *testing.T) {
	t.Parallel()
	fx := setupRegistryForUpdateNodes(t)

	_, err := changeset.UpdateNodes{}.Apply(fx.env, changeset.UpdateNodesInput{
		RegistryQualifier: fx.qualifier,
		RegistryChainSel:  fx.selector,
		Nodes: []changeset.CapabilitiesRegistryNodeParams{
			{
				NOP:                 "non-existent-nop",
				Signer:              signer1,
				P2pID:               fx.nodes[0].P2pID,
				EncryptionPublicKey: encryptionPublicKey,
				CsaKey:              csaKey,
			},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in contract")
}
