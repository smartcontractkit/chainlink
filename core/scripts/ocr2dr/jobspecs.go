package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

// reads a template, replaces placeholders, saves under artifacts
func processJobSpecs(cfg *config, nodes []Node) {
	// bootstrap
	lines, err := readLines(filepath.Join(templatesDir, bootstrapSpecTemplate))
	helpers.PanicErr(err)

	lines = replatePlaceholders(lines, cfg, &nodes[0], &nodes[0])
	outputPath := filepath.Join(artefactsDir, nodes[0].Host+".toml")
	err = writeLines(lines, outputPath)
	helpers.PanicErr(err)
	fmt.Printf("Processed: %s\n", outputPath)

	// oracles nodes
	lines, err = readLines(filepath.Join(templatesDir, oracleSpecTemplate))
	helpers.PanicErr(err)
	for i := 1; i < len(nodes); i++ {
		oracleLines := replatePlaceholders(lines, cfg, &nodes[0], &nodes[i])
		outputPath := filepath.Join(artefactsDir, nodes[i].Host+".toml")
		err = writeLines(oracleLines, outputPath)
		helpers.PanicErr(err)
		fmt.Printf("Processed: %s\n", outputPath)
	}
}

func replatePlaceholders(lines []string, cfg *config, boot *Node, node *Node) (output []string) {
	chainIDStr := strconv.FormatInt(cfg.ChainID, 10)
	contractHex := cfg.DONContractAddress.Hex()
	for _, l := range lines {
		l = strings.Replace(l, "{{chain_id}}", chainIDStr, 1)
		l = strings.Replace(l, "{{oracle_contract_address}}", contractHex, 1)
		l = strings.Replace(l, "{{node_eth_address}}", node.ETHKeys[0], 1)
		l = strings.Replace(l, "{{ocr2_key_bundle_id}}", node.OCR2KeyIDs[0], 1)
		l = strings.Replace(l, "{{node_eth_address}}", node.ETHKeys[0], 1)
		l = strings.Replace(l, "{{p2p_bootstrapper}}", boot.P2PPeerIDS[0], 1)
		output = append(output, l)
	}
	return
}
