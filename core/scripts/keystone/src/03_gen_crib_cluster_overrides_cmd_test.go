package src

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerateCribConfig(t *testing.T) {
	chainID := int64(11155111)
	templatesDir := "../templates"
	forwarderAddress := "0x1234567890abcdef"
	publicKeysPath := "./testdata/PublicKeys.json"

	lines := generateCribConfig(defaultNodeList, publicKeysPath, &chainID, templatesDir, forwarderAddress)

	snaps.MatchSnapshot(t, strings.Join(lines, "\n"))
}
