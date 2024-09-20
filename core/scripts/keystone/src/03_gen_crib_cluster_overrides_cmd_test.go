package src

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerateCribConfig(t *testing.T) {
	chainID := int64(1337)
	templatesDir := "../templates"
	forwarderAddress := "0x1234567890abcdef"
	externalRegistryAddress := "0xabcdef1234567890"
	publicKeysPath := "./testdata/PublicKeys.json"

	lines := generateCribConfig(defaultNodeList, publicKeysPath, &chainID, templatesDir, forwarderAddress, externalRegistryAddress)

	snaps.MatchSnapshot(t, strings.Join(lines, "\n"))
}
