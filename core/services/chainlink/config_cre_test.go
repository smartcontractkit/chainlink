package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

const (
	secretsCRE = `
[CRE.Streams]
ApiKey = "streams-api-key"
ApiSecret = "streams-api-secret"
`
	configCRE = `
[CRE.Streams]
RestURL = "streams.url"
WsURL = "streams.url"
`
)

func TestCREConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		SecretsStrings: []string{secretsCRE},
		ConfigStrings:  []string{configCRE},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	c := cfg.CRE()
	assert.Equal(t, "streams-api-key", c.StreamsApiKey())
	assert.Equal(t, "streams-api-secret", c.StreamsApiSecret())
	assert.Equal(t, "streams.url", c.WsURL())
	assert.Equal(t, "streams.url", c.RestURL())
}

func TestEmptyCREConfig(t *testing.T) {
	cfg := creConfig{s: toml.CreSecrets{}, c: toml.CreConfig{}}
	assert.Equal(t, "", cfg.StreamsApiKey())
	assert.Equal(t, "", cfg.StreamsApiSecret())
	assert.Equal(t, "", cfg.WsURL())
	assert.Equal(t, "", cfg.RestURL())
}
