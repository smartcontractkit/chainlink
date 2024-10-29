package src

import (
	"context"
	"flag"
	"fmt"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type provisionCR struct{}

func NewProvisionCapabilitesRegistryCommand() *provisionCR {
	return &provisionCR{}
}

func (c *provisionCR) Name() string {
	return "provision-capabilities-registry"
}

func (c *provisionCR) Run(args []string) {
	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 1337, "Chain ID of the Ethereum network to deploy to")
	accountKey := fs.String("accountkey", "", "Private key of the account to deploy from")
	nodeSetsPath := fs.String("nodesets", "", "Custom node sets location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")
	nodeSetSize := fs.Int("nodesetsize", 5, "Number of nodes in a nodeset")

	err := fs.Parse(args)
	if err != nil ||
		*chainID == 0 || chainID == nil ||
		*ethUrl == "" || ethUrl == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	if *artefactsDir == "" {
		*artefactsDir = defaultArtefactsDir
	}
	if *nodeSetsPath == "" {
		*nodeSetsPath = defaultNodeSetsPath
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")

	env := helpers.SetupEnv(false)

	// We skip the first node in the node set as it is the bootstrap node
	// technically we could do a single addnodes call here if we merged all the nodes together
	nodeSets := downloadNodeSets(
		*chainID,
		*nodeSetsPath,
		*nodeSetSize,
	)
	provisionCapabillitiesRegistry(env, nodeSets, *chainID, *artefactsDir)
}

func provisionCapabillitiesRegistry(env helpers.Environment, nodeSets NodeSets, chainID int64, artefactsDir string) kcr.CapabilitiesRegistryInterface {
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
