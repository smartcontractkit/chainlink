package src

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type generateCribClusterOverrides struct {
}

func NewGenerateCribClusterOverridesCommand() *generateCribClusterOverrides {
	return &generateCribClusterOverrides{}
}

func (g *generateCribClusterOverrides) Name() string {
	return "generate-crib"
}

func (g *generateCribClusterOverrides) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 11155111, "chain id")
	outputPath := fs.String("outpath", "../crib", "the path to output the generated overrides")

	deployedContracts, err := LoadDeployedContracts()
	helpers.PanicErr(err)
	templatesDir := "templates"
	err = fs.Parse(args)
	if err != nil || outputPath == nil || *outputPath == "" || chainID == nil || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}

	lines := generateCribConfig(".cache/PublicKeys.json", chainID, templatesDir, deployedContracts.ForwarderContract.Hex())

	cribOverridesStr := strings.Join(lines, "\n")
	err = os.WriteFile(filepath.Join(*outputPath, "crib-cluster-overrides.yaml"), []byte(cribOverridesStr), 0600)
	helpers.PanicErr(err)
}

func generateCribConfig(pubKeysPath string, chainID *int64, templatesDir string, forwarderAddress string) []string {
	nca := downloadNodePubKeys(*chainID, pubKeysPath)
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
