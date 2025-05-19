package src

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gkampitakis/go-snaps/snaps"
	"gopkg.in/yaml.v3"
)

func TestGeneratePostprovisionConfig(t *testing.T) {
	chainID := int64(1337)
	capabilitiesP2PPort := int64(6691)
	nodeSetsPath := "./testdata/node_sets.json"
	nodeSetSize := 5
	forwarderAddress := common.Address([20]byte{0: 1}).Hex()
	capabilitiesRegistryAddress := common.Address([20]byte{0: 2}).Hex()
	nodeSets := downloadNodeSets(chainID, nodeSetsPath, nodeSetSize)

	chart := generatePostprovisionConfig(nodeSets, chainID, capabilitiesP2PPort, forwarderAddress, capabilitiesRegistryAddress)

	yamlData, err := yaml.Marshal(chart)
	if err != nil {
		t.Fatalf("Failed to marshal chart: %v", err)
	}

	linesStr := strings.Split(string(yamlData), "\n")
	snaps.MatchSnapshot(t, strings.Join(linesStr, "\n"))
}

func TestGeneratePreprovisionConfig(t *testing.T) {
	nodeSetSize := 5

	chart := generatePreprovisionConfig(nodeSetSize)

	yamlData, err := yaml.Marshal(chart)
	if err != nil {
		t.Fatalf("Failed to marshal chart: %v", err)
	}

	linesStr := strings.Split(string(yamlData), "\n")
	snaps.MatchSnapshot(t, strings.Join(linesStr, "\n"))
}
