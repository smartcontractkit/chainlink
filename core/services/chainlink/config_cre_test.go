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
APIKey = "streams-api-key"
APISecret = "streams-api-secret"
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
	assert.Equal(t, "streams-api-key", c.StreamsAPIKey())
	assert.Equal(t, "streams-api-secret", c.StreamsAPISecret())
	assert.Equal(t, "streams.url", c.WsURL())
	assert.Equal(t, "streams.url", c.RestURL())
}

func TestEmptyCREConfig(t *testing.T) {
	cfg := creConfig{s: toml.CreSecrets{}, c: toml.CreConfig{}}
	assert.Equal(t, "", cfg.StreamsAPIKey())
	assert.Equal(t, "", cfg.StreamsAPISecret())
	assert.Equal(t, "", cfg.WsURL())
	assert.Equal(t, "", cfg.RestURL())
}
