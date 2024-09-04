package src

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenStreamsTriggerJobSpecs(t *testing.T) {
	pubkeysPath := "./testdata/PublicKeys.json"
	nodeListPath := "./testdata/NodeList.txt"
	templatesDir := "../templates"

	feedID := "feed123"
	linkFeedID := "linkfeed123"
	nativeFeedID := "nativefeed123"

	chainID := int64(123456)
	fromBlock := int64(10)

	verifierContractAddress := "verifier_contract_address"
	verifierProxyContractAddress := "verifier_proxy_contract_address"

	output := genStreamsTriggerJobSpecs(
		pubkeysPath,
		nodeListPath,
		templatesDir,
		feedID,
		linkFeedID,
		nativeFeedID,
		chainID,
		fromBlock,
		verifierContractAddress,
		verifierProxyContractAddress,
	)
	prettyOutputs := []string{} 
	for _, o := range output {
		prettyOutputs = append(prettyOutputs, strings.Join(o, "\n"))
	}

	testOutput := strings.Join(prettyOutputs, "\n\n-------------------------------------------------\n\n")
	snaps.MatchSnapshot(t, testOutput)
}
