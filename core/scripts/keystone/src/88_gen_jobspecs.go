package src

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type spec []string

func (s spec) ToString() string {
	return strings.Join(s, "\n")
}

type hostSpec struct {
	spec spec
	host string
}

type donHostSpec struct {
	bootstrap hostSpec
	oracles   []hostSpec
}

func genSpecs(
	pubkeysPath string,
	nodeListPath string,
	templatesDir string,
	chainID int64,
	p2pPort int64,
	ocrConfigContractAddress string,
) donHostSpec {
	nodes := downloadNodeAPICredentials(nodeListPath)
	nca := downloadNodePubKeys(chainID, pubkeysPath)
	bootstrapNode := nca[0]

	bootstrapSpecLines, err := readLines(filepath.Join(templatesDir, bootstrapSpecTemplate))
	helpers.PanicErr(err)
	bootHost := nodes[0].url.Host
	bootstrapSpecLines = replacePlaceholders(
		bootstrapSpecLines,
		chainID, p2pPort,
		ocrConfigContractAddress, bootHost,
		bootstrapNode, bootstrapNode,
	)
	bootstrap := hostSpec{bootstrapSpecLines, bootHost}

	oracleSpecLinesTemplate, err := readLines(filepath.Join(templatesDir, oracleSpecTemplate))
	helpers.PanicErr(err)
	oracles := []hostSpec{}
	for i := 1; i < len(nodes); i++ {
		oracleSpecLines := oracleSpecLinesTemplate
		oracleSpecLines = replacePlaceholders(
			oracleSpecLines,
			chainID, p2pPort,
			ocrConfigContractAddress, bootHost,
			bootstrapNode, nca[i],
		)
		oracles = append(oracles, hostSpec{oracleSpecLines, nodes[i].url.Host})
	}

	return donHostSpec{
		bootstrap: bootstrap,
		oracles:   oracles,
	}
}

func replacePlaceholders(
	lines []string,

	chainID, p2pPort int64,
	contractAddress, bootHost string,
	boot, node NodeKeys,
) (output []string) {
	chainIDStr := strconv.FormatInt(chainID, 10)
	bootstrapper := fmt.Sprintf("%s@%s:%d", boot.P2PPeerID, bootHost, p2pPort)
	for _, l := range lines {
		l = strings.Replace(l, "{{ chain_id }}", chainIDStr, 1)
		l = strings.Replace(l, "{{ ocr_config_contract_address }}", contractAddress, 1)
		l = strings.Replace(l, "{{ transmitter_id }}", node.EthAddress, 1)
		l = strings.Replace(l, "{{ ocr_key_bundle_id }}", node.OCR2BundleID, 1)
		l = strings.Replace(l, "{{ bootstrapper_p2p_id }}", bootstrapper, 1)
		output = append(output, l)
	}
	return
}
