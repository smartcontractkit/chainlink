// This sub CLI acts as a temporary shim for external aptos support

package src

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type Toolkit struct{}

func (t *Toolkit) Name() string {
	return "toolkit"
}

func NewToolkit() *Toolkit {
	return &Toolkit{}
}

func (t *Toolkit) Run(args []string) {
	if len(args) < 1 {
		fmt.Println("Available commands:")
		fmt.Println("  deploy-workflows")
		fmt.Println("  deploy-ocr3-contracts")
		fmt.Println("  deploy-ocr3-jobspecs")
		os.Exit(1)
	}

	command := args[0]
	cmdArgs := args[1:]

	switch command {
	case "get-aptos-keys":
		t.AptosKeys(cmdArgs)
	case "deploy-workflows":
		t.DeployWorkflows(cmdArgs)
	case "deploy-ocr3-contracts":
		t.ProvisionOCR3Contracts(cmdArgs)
	case "deploy-ocr3-jobspecs":
		t.DeployOCR3JobSpecs(cmdArgs)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func (t *Toolkit) AptosKeys(args []string) {
	fs := flag.NewFlagSet("get-aptos-keys", flag.ExitOnError)
	nodesListPath := fs.String("nodes", ".cache/NodesList.txt", "Path to file with list of nodes")
	artefacts := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")
	chainID := fs.Int64("chainid", 1337, "Chain ID")

	if err := fs.Parse(args); err != nil {
		fs.Usage()
		os.Exit(1)
	}

	nodes := mustReadNodesList(*nodesListPath)
	keys := mustFetchNodeKeys(*chainID, nodes, true)

	mustWriteJSON(*artefacts+"/pubnodekeys.json", keys)
}

func (t *Toolkit) ProvisionOCR3Contracts(args []string) {
	fs := flag.NewFlagSet("deploy-ocr3-contracts", flag.ExitOnError)
	ethURL := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "Private key of the deployer account")
	chainID := fs.Int64("chainid", 1337, "Chain ID")
	nodesListPath := fs.String("nodes", ".cache/NodesList.txt", "Path to file with list of nodes")
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")
	ocrConfigFile := fs.String("ocrfile", "ocr_config.json", "Path to OCR config file")

	if err := fs.Parse(args); err != nil || *ethURL == "" || *accountKey == "" {
		fs.Usage()
		os.Exit(1)
	}

	// Set environment variables required by setupenv
	os.Setenv("ETH_URL", *ethURL)
	os.Setenv("ETH_CHAIN_ID", strconv.FormatInt(*chainID, 10))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")

	env := helpers.SetupEnv(false)

	nodes := mustReadNodesList(*nodesListPath)
	nodeKeys := mustFetchNodeKeys(*chainID, nodes, true)

	deployOCR3Contract(nodeKeys, env, *ocrConfigFile, *artefactsDir)
}

func (t *Toolkit) DeployOCR3JobSpecs(args []string) {
	fs := flag.NewFlagSet("deploy-ocr3-jobspecs", flag.ExitOnError)

	ethURL := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "Private key of the deployer account")
	chainID := fs.Int64("chainid", 1337, "Chain ID")
	nodesListPath := fs.String("nodes", ".cache/NodesList.txt", "Path to file with list of nodes")
	p2pPort := fs.Int64("p2pport", 6690, "P2P port")
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")

	if err := fs.Parse(args); err != nil || *ethURL == "" || *accountKey == "" {
		fs.Usage()
		os.Exit(1)
	}

	os.Setenv("ETH_URL", *ethURL)
	os.Setenv("ETH_CHAIN_ID", strconv.FormatInt(*chainID, 10))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")

	env := helpers.SetupEnv(false)

	nodes := mustReadNodesList(*nodesListPath)
	nodeKeys := mustFetchNodeKeys(*chainID, nodes, true)
	o := LoadOnchainMeta(*artefactsDir, env)

	deployOCR3JobSpecs(
		nodes,
		*chainID,
		nodeKeys,
		*p2pPort,
		o,
	)
}

func (t *Toolkit) DeployWorkflows(args []string) {
	fs := flag.NewFlagSet("deploy-workflows", flag.ExitOnError)
	workflowFile := fs.String("workflow", "", "Path to workflow file")
	nodesList := fs.String("nodes", ".cache/NodesList.txt", "Path to file with list of nodes")

	if err := fs.Parse(args); err != nil || *workflowFile == "" {
		fs.Usage()
		os.Exit(1)
	}

	nodesWithCreds := mustReadNodesList(*nodesList)

	for _, node := range nodesWithCreds {
		api := newNodeAPI(node)
		workflowContent, err := os.ReadFile(*workflowFile)
		PanicErr(err)

		upsertJob(api, "workflow", string(workflowContent))
		fmt.Println("Workflow deployed successfully")
	}
}

// Reads in a list of nodes from a file, where each line is in the format:
// http://localhost:50100 http://chainlink.core.1:50100 notreal@fakeemail.ch fj293fbBnlQ!f9vNs
func mustReadNodesList(path string) []NodeWithCreds {
	fmt.Println("Reading nodes list from", path)
	nodesList, err := readLines(path)
	helpers.PanicErr(err)

	nodes := make([]NodeWithCreds, 0, len(nodesList))
	for _, r := range nodesList {
		rr := strings.TrimSpace(r)
		if len(rr) == 0 {
			continue
		}
		s := strings.Split(rr, " ")
		if len(s) != 4 {
			helpers.PanicErr(errors.New("wrong nodes list format"))
		}

		r := SimpleURL{
			Scheme: "http",
			Host:   s[0],
		}
		u := SimpleURL{
			Scheme: "http",
			Host:   s[1],
		}
		remoteURL, err := url.Parse(u.String())
		PanicErr(err)
		nodes = append(nodes, NodeWithCreds{
			URL:       u,
			RemoteURL: r,
			// This is the equivalent of "chainlink.core.1" in our above example
			ServiceName:      remoteURL.Hostname(),
			APILogin:         s[2],
			APIPassword:      s[3],
			KeystorePassword: "",
		})
	}

	return nodes
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
