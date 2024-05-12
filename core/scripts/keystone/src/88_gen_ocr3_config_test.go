package src

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerateOCR3Config(t *testing.T) {
	// Generate OCR3 config
	config := generateOCR3Config("./testdata/SampleConfig.json", 11155111, "./testdata/PublicKeys.json")
	snaps.MatchJSON(t, config)
}
