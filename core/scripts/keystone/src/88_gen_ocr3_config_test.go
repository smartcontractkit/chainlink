package src

import (
	"errors"
	"testing"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerateOCR3Config(t *testing.T) {
	// Generate OCR3 config
	config := generateOCR3Config("./testdata/keyless_node_sets.json", "./testdata/SampleConfig.json", 1337, "./testdata/node_sets.json", 4)

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
