package src

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"

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
	ctx := context.Background()

	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	// create flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 1337, "chain ID of the Ethereum network to deploy to")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	nodeList := fs.String("nodes", "", "Custom node list location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")

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
	if *publicKeys == "" {
		*publicKeys = defaultPublicKeys
	}
	if *nodeList == "" {
		*nodeList = defaultNodeList
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")

	env := helpers.SetupEnv(false)

	reg := getOrDeployCapabilitiesRegistry(ctx, *artefactsDir, env)

	// For now, trigger, target, and workflow DONs are the same node sets, and same don instance
	workflowDON := loadDON(
		*publicKeys,
		*chainID,
		*nodeList,
	)
	crProvisioner := NewCapabilityRegistryProvisioner(reg, env)
	// We're using the default capability set for now
	capSet := NewCapabilitySet()
	crProvisioner.AddCapabilities(ctx, capSet)

	nodeOperator := NewNodeOperator(env.Owner.From, "MY_NODE_OPERATOR", workflowDON)
	crProvisioner.AddNodeOperator(ctx, nodeOperator)

	// Note that both of these calls are simplified versions of the actual calls
	//
	// See the method documentation for more details
	crProvisioner.AddNodes(ctx, nodeOperator, capSet)
	crProvisioner.AddDON(ctx, nodeOperator, capSet, false, true)
}

func loadDON(publicKeys string, chainID int64, nodeList string) []peer {
	nca := downloadNodePubKeys(nodeList, chainID, publicKeys)
	workflowDON := []peer{}
	for _, n := range nca {

		p := peer{
			PeerID: n.P2PPeerID,
			Signer: n.OCR2OnchainPublicKey,
		}
		workflowDON = append(workflowDON, p)
	}
	return workflowDON
}

func getOrDeployCapabilitiesRegistry(ctx context.Context, artefactsDir string, env helpers.Environment) *kcr.CapabilitiesRegistry {
	contracts, err := LoadDeployedContracts(artefactsDir)
	if err != nil {
		fmt.Println("Could not load deployed contracts, deploying new ones")
		// panic(err)
	}

	if contracts.CapabilityRegistry.String() == (common.Address{}).String() {
		_, tx, capabilitiesRegistry, innerErr := kcr.DeployCapabilitiesRegistry(env.Owner, env.Ec)
		if innerErr != nil {
			panic(innerErr)
		}

		helpers.ConfirmContractDeployed(ctx, env.Ec, tx, env.ChainID)
		contracts.CapabilityRegistry = capabilitiesRegistry.Address()
		WriteDeployedContracts(contracts, artefactsDir)
		return capabilitiesRegistry
	} else {
		capabilitiesRegistry, innerErr := kcr.NewCapabilitiesRegistry(contracts.CapabilityRegistry, env.Ec)
		if innerErr != nil {
			panic(innerErr)
		}

		return capabilitiesRegistry
	}
}
