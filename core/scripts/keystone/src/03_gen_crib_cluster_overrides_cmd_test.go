package src

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"gopkg.in/yaml.v3"
)

func TestGeneratePostprovisionConfig(t *testing.T) {
	chainID := int64(1337)
	publicKeysPath := "./testdata/PublicKeys.json"
	nodeListPath := "./testdata/NodeList.txt"
	
	contracts := deployedContracts{
		OCRContract:        [20]byte{0: 1},
		ForwarderContract:  [20]byte{0: 2},
		CapabilityRegistry: [20]byte{0: 3},
		SetConfigTxBlock:   0,
	}

	nodeSetSize := 4

	chart := generatePostprovisionConfig(&nodeListPath, &chainID, &publicKeysPath, contracts, nodeSetSize)

	yamlData, err := yaml.Marshal(chart)
	if err != nil {
		t.Fatalf("Failed to marshal chart: %v", err)
	}

	linesStr := strings.Split(string(yamlData), "\n")
	snaps.MatchSnapshot(t, strings.Join(linesStr, "\n"))
}

func TestGeneratePreprovisionConfig(t *testing.T) {
	nodeSetSize := 4

	chart := generatePreprovisionConfig(nodeSetSize)

	yamlData, err := yaml.Marshal(chart)
	if err != nil {
		t.Fatalf("Failed to marshal chart: %v", err)
	}

	linesStr := strings.Split(string(yamlData), "\n")
	snaps.MatchSnapshot(t, strings.Join(linesStr, "\n"))
}
