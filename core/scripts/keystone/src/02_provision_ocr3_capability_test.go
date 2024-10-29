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
	config := generateOCR3Config(nodeSet.Workflow, "./testdata/SampleConfig.json", 1337)

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

func (d *donHostSpec) ToString() string {
	var result string
	result += "Bootstrap:\n"
	result += "Host: " + d.bootstrap.host + "\n"
	result += d.bootstrap.spec.ToString()
	result += "\n\nOracles:\n"
	for i, oracle := range d.oracles {
		if i != 0 {
			result += "--------------------------------\n"
		}
		result += fmt.Sprintf("Oracle %d:\n", i)
		result += "Host: " + oracle.host + "\n"
		result += oracle.spec.ToString()
		result += "\n\n"
	}
	return result
}

func TestGenSpecs(t *testing.T) {
	nodeSetsPath := "./testdata/node_sets.json"
	chainID := int64(1337)
	p2pPort := int64(6690)
	contractAddress := "0xB29934624cAe3765E33115A9530a13f5aEC7fa8A"
	nodeSet := downloadNodeSets(chainID, nodeSetsPath, 4).Workflow
	specs := generateOCR3JobSpecs(nodeSet, "../templates", chainID, p2pPort, contractAddress)
	snaps.MatchSnapshot(t, specs.ToString())
}
