package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type provisionKeystone struct{}

func NewProvisionKeystoneCommand() *provisionKeystone {
	return &provisionKeystone{}
}

func (g *provisionKeystone) Name() string {
	return "provision-keystone"
}

func (g *provisionKeystone) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)

	// common flags
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")
	nodeSetSize := fs.Int("nodesetsize", 5, "number of nodes in a nodeset")
	nodeSetsPath := fs.String("nodesets", defaultNodeSetsPath, "Custom node sets location")
	chainID := fs.Int64("chainid", 1337, "chain ID of the Ethereum network to deploy to")

	// preprovisioning flags
	preprovison := fs.Bool("preprovision", false, "Preprovision crib")

	// provisioning flags
	ethURL := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	ocrConfigFile := fs.String("ocrfile", "ocr_config.json", "path to OCR config file")
	p2pPort := fs.Int64("p2pport", 6690, "p2p port")
	capabilitiesP2PPort := fs.Int64("capabilitiesp2pport", 6691, "p2p port for capabilities")
	preprovisionConfigName := fs.String("preprovisionconfig", "crib-preprovision.yaml", "Name of the preprovision config file, stored in the artefacts directory")
	postprovisionConfigName := fs.String("postprovisionconfig", "crib-postprovision.yaml", "Name of the postprovision config file, stored in the artefacts directory")
	// additional flags
	clean := fs.Bool("clean", false, "Clean up resources before provisioning")

	err := fs.Parse(args)

	if err != nil || (!*preprovison && (*ethURL == "" || *accountKey == "")) {
		fs.Usage()
		os.Exit(1)
	}

	if *preprovison {
		fmt.Println()
		fmt.Println()
		fmt.Println("========================")
		fmt.Println("Writing Preprovisioning Config")
		fmt.Println("========================")
		fmt.Println()
		fmt.Println()
		writePreprovisionConfig(*nodeSetSize, filepath.Join(*artefactsDir, *preprovisionConfigName))
		return
	}

	// We always want to start with a clean slate
	/// when it comes to nodesets
	err = os.RemoveAll(*nodeSetsPath)
	PanicErr(err)
	fmt.Println("Collecting node sets...")
	nodeSets := downloadNodeSets(*chainID, *nodeSetsPath, *nodeSetSize)

	if *clean {
		fmt.Println("Cleaning up resources")
		for _, node := range nodeSets.Workflow.Nodes {
			clearJobs(newNodeAPI(node))
		}
		for _, node := range nodeSets.StreamsTrigger.Nodes {
			clearJobs(newNodeAPI(node))
		}
		os.RemoveAll(*artefactsDir)
	}

	// Kinda hacky but it prevents us from refactoring the setupenv function which
	// is used in many other places
	os.Setenv("ETH_URL", *ethURL)
	os.Setenv("ETH_CHAIN_ID", strconv.FormatInt(*chainID, 10))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	env := helpers.SetupEnv(false)

	provisionStreamsDON(
		env,
		nodeSets.StreamsTrigger,
		*chainID,
		*p2pPort,
		*ocrConfigFile,
		*artefactsDir,
	)

	reg := provisionCapabilitiesRegistry(
		env,
		nodeSets,
		*chainID,
		*artefactsDir,
	)

	onchainMeta := provisionWorkflowDON(
		env,
		nodeSets.Workflow,
		*chainID,
		*p2pPort,
		*ocrConfigFile,
		*artefactsDir,
		reg,
	)

	fmt.Println()
	fmt.Println()
	fmt.Println("========================")
	fmt.Println("Writing Postprovision Config")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println()

	writePostProvisionConfig(
		nodeSets,
		*chainID,
		*capabilitiesP2PPort,
		onchainMeta.Forwarder.Address().Hex(),
		onchainMeta.CapabilitiesRegistry.Address().Hex(),
		filepath.Join(*artefactsDir, *postprovisionConfigName),
	)
}

func provisionCapabilitiesRegistry(
	env helpers.Environment,
	nodeSets NodeSets,
	chainID int64,
	artefactsDir string,
) kcr.CapabilitiesRegistryInterface {
	fmt.Println()
	fmt.Println()
	fmt.Println("========================")
	fmt.Println("Provisioning Capabilities Registry DON")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println()
	reg := provisionCapabillitiesRegistry(
		env,
		nodeSets,
		chainID,
		artefactsDir,
	)
	return reg
}

func provisionStreamsDON(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFilePath string,
	artefactsDir string,
) {
	fmt.Println()
	fmt.Println()
	fmt.Println("========================")
	fmt.Println("Provisioning streams DON")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println()
	setupStreamsTrigger(
		env,
		nodeSet,
		chainID,
		p2pPort,
		ocrConfigFilePath,
		artefactsDir,
	)
}

func provisionWorkflowDON(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFile string,
	artefactsDir string,
	reg kcr.CapabilitiesRegistryInterface,
) (onchainMeta *onchainMeta) {
	fmt.Println()
	fmt.Println()
	fmt.Println("========================")
	fmt.Println("Provisioning workflow DON")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println()
	deployForwarder(env, artefactsDir)

	onchainMeta, _ = provisionOCR3(
		env,
		nodeSet,
		chainID,
		p2pPort,
		ocrConfigFile,
		artefactsDir,
	)
	distributeFunds(nodeSet.NodeKeys, env)

	// We don't technically need the capability registry as a dependency
	// as we just use it for a sanity check
	// We could remove it so that we can execute provisioning in parallel
	deployKeystoneWorkflowsTo(nodeSet, reg)

	return onchainMeta
}
