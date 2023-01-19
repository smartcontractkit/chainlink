package config

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_PluginConfig(t *testing.T) {
	h := utils.NewHash()
	rawToml := fmt.Sprintf(`
FeedID = "%s"
URL = "http://example.com/reports"
ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
ClientPubKey = "29d2237673697bbc87b67c8500a86cfe55f59c8eb313b7a4e33d0f3f55a1cc84"
`, h.String())

	var mc PluginConfig
	err := toml.Unmarshal([]byte(rawToml), &mc)
	require.NoError(t, err)

	assert.Equal(t, h, mc.FeedID)
	assert.Equal(t, "http://example.com/reports", mc.URL.String())
	assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
	assert.Equal(t, "29d2237673697bbc87b67c8500a86cfe55f59c8eb313b7a4e33d0f3f55a1cc84", mc.ClientPubKey.String())
}
