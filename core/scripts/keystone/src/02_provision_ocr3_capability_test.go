package src

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerateOCR3Config(t *testing.T) {
	// Generate OCR3 config
	nodeSet := downloadNodeSets(1337, "./testdata/node_sets.json", 4)
	nodeKeys := nodeSet.Workflow.NodeKeys
	config := generateOCR3Config(nodeKeys, "./testdata/SampleConfig.json")

	matchOffchainConfig := match.Custom("OffchainConfig", func(s any) (any, error) {
		// coerce the value to a string
		s, ok := s.(string)
		if !ok {
			return nil, errors.New("offchain config is not a string")
		}

		// if the string is not empty
		if s == "" {
			return nil, errors.New("offchain config is empty")
		}

		return "<nonemptyvalue>", nil
	})

	snaps.MatchJSON(t, config, matchOffchainConfig)
}

func TestGenSpecs(t *testing.T) {
	nodeSetsPath := "./testdata/node_sets.json"
	chainID := int64(1337)
	p2pPort := int64(6690)
	contractAddress := "0xB29934624cAe3765E33115A9530a13f5aEC7fa8A"
	nodeSet := downloadNodeSets(chainID, nodeSetsPath, 4).Workflow

	// Create Bootstrap Job Spec
	bootstrapConfig := BootstrapJobSpecConfig{
		JobSpecName:              "ocr3_bootstrap",
		OCRConfigContractAddress: contractAddress,
		ChainID:                  chainID,
	}
	bootstrapSpec := createBootstrapJobSpec(bootstrapConfig)

	// Create Oracle Job Spec
	oracleConfig := OracleJobSpecConfig{
		JobSpecName:              "ocr3_oracle",
		OCRConfigContractAddress: contractAddress,
		OCRKeyBundleID:           nodeSet.NodeKeys[1].OCR2BundleID,
		BootstrapURI:             fmt.Sprintf("%s@%s:%d", nodeSet.NodeKeys[0].P2PPeerID, nodeSet.Nodes[0].ServiceName, p2pPort),
		TransmitterID:            nodeSet.NodeKeys[1].P2PPeerID,
		ChainID:                  chainID,
		AptosKeyBundleID:         nodeSet.NodeKeys[1].AptosBundleID,
	}
	oracleSpec := createOracleJobSpec(oracleConfig)

	// Combine Specs
	generatedSpecs := fmt.Sprintf("%s\n\n%s", bootstrapSpec, oracleSpec)

	snaps.MatchSnapshot(t, generatedSpecs)
}
