package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type generateJobSpecs struct {
}

func NewGenerateJobSpecsCommand() *generateJobSpecs {
	return &generateJobSpecs{}
}

func (g *generateJobSpecs) Name() string {
	return "generate-jobspecs"
}

func (g *generateJobSpecs) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	chainID := fs.Int64("chainid", 80001, "chain id")
	p2pPort := fs.Int64("p2pport", 6690, "p2p port")
	contractAddress := fs.String("contract", "", "oracle contract address")
	err := fs.Parse(args)
	if err != nil || nodesFile == nil || *nodesFile == "" || contractAddress == nil || *contractAddress == "" {
		fs.Usage()
		os.Exit(1)
	}

	nodes := mustReadNodesList(*nodesFile)
	nca := mustFetchNodesKeys(*chainID, nodes)
	bootstrapNode := nca[0]

	lines, err := readLines(filepath.Join(templatesDir, bootstrapSpecTemplate))
	helpers.PanicErr(err)

	bootHost := nodes[0].url.Host
	lines = replacePlaceholders(lines, *chainID, *p2pPort, *contractAddress, bootHost, bootstrapNode, bootstrapNode)
	outputPath := filepath.Join(artefactsDir, bootHost+".toml")
	err = writeLines(lines, outputPath)
	helpers.PanicErr(err)
	fmt.Println("Saved bootstrap node jobspec:", outputPath)

	lines, err = readLines(filepath.Join(templatesDir, oracleSpecTemplate))
	helpers.PanicErr(err)
	for i := 1; i < len(nodes); i++ {
		oracleLines := replacePlaceholders(lines, *chainID, *p2pPort, *contractAddress, bootHost, bootstrapNode, nca[i])
		outputPath := filepath.Join(artefactsDir, nodes[i].url.Host+".toml")
		err = writeLines(oracleLines, outputPath)
		helpers.PanicErr(err)
		fmt.Println("Saved oracle node jobspec:", outputPath)
	}
}

func replacePlaceholders(lines []string, chainID, p2pPort int64, contractAddress, bootHost string, boot *nodeConfig, node *nodeConfig) (output []string) {
	chainIDStr := strconv.FormatInt(chainID, 10)
	bootstrapper := fmt.Sprintf("%s@%s:%d", boot.p2pPeerID, bootHost, p2pPort)
	ts := time.Now().UTC().Format("2006-01-02T15:04")
	for _, l := range lines {
		l = strings.Replace(l, "{{chain_id}}", chainIDStr, 1)
		l = strings.Replace(l, "{{oracle_contract_address}}", contractAddress, 1)
		l = strings.Replace(l, "{{node_eth_address}}", node.ethAddress, 1)
		l = strings.Replace(l, "{{ocr2_key_bundle_id}}", node.ocr2BundleID, 1)
		l = strings.Replace(l, "{{p2p_bootstrapper}}", bootstrapper, 1)
		l = strings.Replace(l, "{{timestamp}}", ts, 1)
		output = append(output, l)
	}
	return
}
