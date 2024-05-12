package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

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
	cribDir := fs.String("cribdir", "../../../crib", "crib directory")

	deployedContracts, err := LoadDeployedContracts()
	helpers.PanicErr(err)
	templatesDir := "templates"
	err = fs.Parse(args)
	if err != nil || cribDir == nil || *cribDir == "" || chainID == nil || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}
	k8sClient := MustNewK8sClient()

	// Get addresses for all non-bootstrap nodes
	// Create bootstrap node job spec
	lines := generateCribConfig(".cache/PublicKeys.json", chainID, templatesDir, deployedContracts.ForwarderContract.Hex())
	yamlPath := filepath.Join(*cribDir, "devspace.keystone.template.yaml")
	yamlBytes, err := os.ReadFile(yamlPath)
	helpers.PanicErr(err)
	var config map[string]interface{}

	err = yaml.Unmarshal(yamlBytes, &config)
	helpers.PanicErr(err)

	// Specify the path to the element you want to change
	path := "deployments.app.helm.values.chainlink.nodes"
	keys := strings.Split(path, ".")

	// Navigate through the map according to the path to find the node
	lastKey := keys[len(keys)-1]
	m := config
	for _, k := range keys[:len(keys)-1] {
		if m[k] == nil {
			m[k] = make(map[interface{}]interface{})
		}
		m = m[k].(map[string]interface{})
	}

	cribOverridesStr := strings.Join(lines, "\n")
	var cribOverridesMap map[string]interface{}
	err = yaml.Unmarshal([]byte(cribOverridesStr), &cribOverridesMap)
	helpers.PanicErr(err)

	// Replace the value at the specific path with the unmarshaled map
	m[lastKey] = cribOverridesMap["nodes"]

	devspaceFileName := fmt.Sprintf("devspace-generated.%s.yaml", k8sClient.namespace)
	outputPath := filepath.Join(*cribDir, devspaceFileName)
	yamlBytes, err = yaml.Marshal(config)
	helpers.PanicErr(err)
	err = os.WriteFile(outputPath, yamlBytes, 0600)
	helpers.PanicErr(err)

	fmt.Printf("Please apply the generated crib overrides to your crib cluster by running the following within the crib directory: DEVSPACE_CONFIG=%s devspace deploy", devspaceFileName)
}

func generateCribConfig(pubKeysPath string, chainID *int64, templatesDir string, forwarderAddress string) []string {
	nca := downloadNodePubKeys(*chainID, pubKeysPath)
	nodeAddresses := []string{}

	for _, node := range nca[1:] {
		nodeAddresses = append(nodeAddresses, node.EthAddress)
	}

	lines, err := readLines(filepath.Join(templatesDir, cribOverrideTemplate))
	helpers.PanicErr(err)
	lines = replaceCribPlaceholders(
		lines,
		forwarderAddress,
		nodeAddresses,
	)
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
