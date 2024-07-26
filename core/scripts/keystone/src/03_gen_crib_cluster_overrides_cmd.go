package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type generateCribClusterOverrides struct {
	PublicKeys string
	NodeList   string
	Artefacts  string
}

func NewGenerateCribClusterOverridesCommand() *generateCribClusterOverrides {
	return &generateCribClusterOverrides{
		PublicKeys: ".cache/PublicKeys.json",
		NodeList:   ".cache/NodeList.txt",
		Artefacts:  artefactsDir,
	}
}

func (g *generateCribClusterOverrides) Name() string {
	return "generate-crib"
}

func (g *generateCribClusterOverrides) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 11155111, "chain id")
	outputPath := fs.String("outpath", "../crib", "the path to output the generated overrides")
	customPublicKeys := fs.String("publickeys", "", "Custom public keys json location")
	customNodeList := fs.String("nodes", "", "Custom node list location")
	customArtefacts := fs.String("artefacts", "", "Custom artefacts directory location")

	templatesDir := "templates"
	err := fs.Parse(args)
	if err != nil || outputPath == nil || *outputPath == "" || chainID == nil || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}

	if *customArtefacts != "" {
		fmt.Printf("Custom  artefacts folder flag detected, using custom path %s", *customArtefacts)
		g.Artefacts = *customArtefacts
	}

	deployedContracts, err := LoadDeployedContracts(g.Artefacts)
	helpers.PanicErr(err)

	if *customPublicKeys != "" {
		fmt.Printf("Custom public keys json override flag detected, using custom path %s", *customPublicKeys)
		g.PublicKeys = *customPublicKeys
	}

	if *customNodeList != "" {
		fmt.Printf("Custom node file override flag detected, using custom node file path %s", *customNodeList)
		g.NodeList = *customNodeList
	}

	lines := generateCribConfig(g.NodeList, g.PublicKeys, chainID, templatesDir, deployedContracts.ForwarderContract.Hex())

	cribOverridesStr := strings.Join(lines, "\n")
	err = os.WriteFile(filepath.Join(*outputPath, "crib-cluster-overrides.yaml"), []byte(cribOverridesStr), 0600)
	helpers.PanicErr(err)
}

func generateCribConfig(nodeList string, pubKeysPath string, chainID *int64, templatesDir string, forwarderAddress string) []string {
	nca := downloadNodePubKeys(nodeList, *chainID, pubKeysPath)
	nodeAddresses := []string{}

	for _, node := range nca[1:] {
		nodeAddresses = append(nodeAddresses, node.EthAddress)
	}

	lines, err := readLines(filepath.Join(templatesDir, cribOverrideTemplate))
	helpers.PanicErr(err)
	lines = replaceCribPlaceholders(lines, forwarderAddress, nodeAddresses)
	return lines
}

func replaceCribPlaceholders(
	lines []string,
	forwarderAddress string,
	nodeFromAddresses []string,
) (output []string) {
	for _, l := range lines {
		l = strings.Replace(l, "{{ forwarder_address }}", forwarderAddress, 1)
		l = strings.Replace(l, "{{ node_2_address }}", nodeFromAddresses[0], 1)
		l = strings.Replace(l, "{{ node_3_address }}", nodeFromAddresses[1], 1)
		l = strings.Replace(l, "{{ node_4_address }}", nodeFromAddresses[2], 1)
		l = strings.Replace(l, "{{ node_5_address }}", nodeFromAddresses[3], 1)
		output = append(output, l)
	}

	return output
}
