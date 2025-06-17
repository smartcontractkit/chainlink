package src

import (
	"context"
	"fmt"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func provisionCapabillitiesRegistry(env helpers.Environment, nodeSets NodeSets, chainID int64, artefactsDir string) kcr.CapabilitiesRegistryInterface {
	fmt.Printf("Provisioning capabilities registry on chain %d\n", chainID)
	ctx := context.Background()
	reg := deployCR(ctx, artefactsDir, env)
	crProvisioner := NewCapabilityRegistryProvisioner(reg, env)
	streamsTriggerCapSet := NewCapabilitySet(NewStreamsTriggerV1Capability())
	workflowCapSet := NewCapabilitySet(NewOCR3V1ConsensusCapability(), NewEthereumGethTestnetV1WriteCapability())
	workflowDON := nodeKeysToDON(nodeSets.Workflow.Name, nodeSets.Workflow.NodeKeys[1:], workflowCapSet)
	streamsTriggerDON := nodeKeysToDON(nodeSets.StreamsTrigger.Name, nodeSets.StreamsTrigger.NodeKeys[1:], streamsTriggerCapSet)

	crProvisioner.AddCapabilities(ctx, MergeCapabilitySets(streamsTriggerCapSet, workflowCapSet))
	dons := map[string]DON{workflowDON.Name: workflowDON, streamsTriggerDON.Name: streamsTriggerDON}
	nodeOperator := NewNodeOperator(env.Owner.From, "MY_NODE_OPERATOR", dons)
	crProvisioner.AddNodeOperator(ctx, nodeOperator)

	crProvisioner.AddNodes(ctx, nodeOperator, nodeSets.Workflow.Name, nodeSets.StreamsTrigger.Name)

	crProvisioner.AddDON(ctx, nodeOperator, nodeSets.Workflow.Name, true, true)
	crProvisioner.AddDON(ctx, nodeOperator, nodeSets.StreamsTrigger.Name, true, false)

	return reg
}

// nodeKeysToDON converts a slice of NodeKeys into a DON struct with the given name and CapabilitySet.
func nodeKeysToDON(donName string, nodeKeys []NodeKeys, capSet CapabilitySet) DON {
	peers := []peer{}
	for _, n := range nodeKeys {
		p := peer{
			PeerID: n.P2PPeerID,
			Signer: n.OCR2OnchainPublicKey,
		}
		peers = append(peers, p)
	}
	return DON{
		F:             1,
		Name:          donName,
		Peers:         peers,
		CapabilitySet: capSet,
	}
}

func deployCR(ctx context.Context, artefactsDir string, env helpers.Environment) kcr.CapabilitiesRegistryInterface {
	o := LoadOnchainMeta(artefactsDir, env)
	// We always redeploy the capabilities registry to ensure it is up to date
	// since we don't have diffing logic to determine if it has changed
	// if o.CapabilitiesRegistry != nil {
	// 	fmt.Println("CapabilitiesRegistry already deployed, skipping...")
	// 	return o.CapabilitiesRegistry
	// }

	_, tx, capabilitiesRegistry, innerErr := kcr.DeployCapabilitiesRegistry(env.Owner, env.Ec)
	PanicErr(innerErr)
	helpers.ConfirmContractDeployed(ctx, env.Ec, tx, env.ChainID)

	o.CapabilitiesRegistry = capabilitiesRegistry
	WriteOnchainMeta(o, artefactsDir)
	return capabilitiesRegistry
}
