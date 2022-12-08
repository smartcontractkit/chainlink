package types

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_MercuryConfig(t *testing.T) {
	h := utils.NewHash()
	rawToml := fmt.Sprintf(`
FeedID = "%s"
URL = "http://example.com/reports"
`, h.String())

	var mc MercuryConfig
	err := toml.Unmarshal([]byte(rawToml), &mc)
	require.NoError(t, err)

	assert.Equal(t, h, mc.FeedID)
	assert.Equal(t, "http://example.com/reports", mc.URL.String())

}
