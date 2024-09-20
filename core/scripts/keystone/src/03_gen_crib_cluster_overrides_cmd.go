package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type generateCribClusterOverrides struct{}

func NewGenerateCribClusterOverridesCommand() *generateCribClusterOverrides {
	return &generateCribClusterOverrides{}
}

func (g *generateCribClusterOverrides) Name() string {
	return "generate-crib"
}

func (g *generateCribClusterOverrides) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 1337, "chain id")
	outputPath := fs.String("outpath", "../crib", "the path to output the generated overrides")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	nodeList := fs.String("nodes", "", "Custom node list location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")

	templatesDir := "templates"
	err := fs.Parse(args)
	if err != nil || outputPath == nil || *outputPath == "" || chainID == nil || *chainID == 0 {
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

	deployedContracts, err := LoadDeployedContracts(*artefactsDir)
	helpers.PanicErr(err)

	lines := generateCribConfig(*nodeList, *publicKeys, chainID, templatesDir, deployedContracts.ForwarderContract.Hex(), deployedContracts.CapabilityRegistry.Hex())

	cribOverridesStr := strings.Join(lines, "\n")
	err = os.WriteFile(filepath.Join(*outputPath, "crib-cluster-overrides.yaml"), []byte(cribOverridesStr), 0600)
	helpers.PanicErr(err)
}

func generateCribConfig(nodeList string, pubKeysPath string, chainID *int64, templatesDir string, forwarderAddress string, externalRegistryAddress string) []string {
	nca := downloadNodePubKeys(nodeList, *chainID, pubKeysPath)
	nodeAddresses := []string{}
	capabilitiesBootstrapper := fmt.Sprintf("%s@%s:%s", nca[0].P2PPeerID, "app-node1", "6691")

	for _, node := range nca[1:] {
		nodeAddresses = append(nodeAddresses, node.EthAddress)
	}

	lines, err := readLines(filepath.Join(templatesDir, cribOverrideTemplate))
	helpers.PanicErr(err)
	lines = replaceCribPlaceholders(lines, forwarderAddress, nodeAddresses, externalRegistryAddress, capabilitiesBootstrapper)
	return lines
}

func replaceCribPlaceholders(
	lines []string,
	forwarderAddress string,
	nodeFromAddresses []string,
	externalRegistryAddress string,
	capabilitiesBootstrapper string,
) (output []string) {
	for _, l := range lines {
		l = strings.Replace(l, "{{ forwarder_address }}", forwarderAddress, 1)
		l = strings.Replace(l, "{{ node_2_address }}", nodeFromAddresses[0], 1)
		l = strings.Replace(l, "{{ node_3_address }}", nodeFromAddresses[1], 1)
		l = strings.Replace(l, "{{ node_4_address }}", nodeFromAddresses[2], 1)
		l = strings.Replace(l, "{{ node_5_address }}", nodeFromAddresses[3], 1)
		l = strings.Replace(l, "{{ external_registry_address }}", externalRegistryAddress, 1)
		l = strings.Replace(l, "{{ capabilities_bootstrapper }}", capabilitiesBootstrapper, 1)
		output = append(output, l)
	}

	return output
}
