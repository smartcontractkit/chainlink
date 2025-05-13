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
)

func TestCREConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		SecretsStrings: []string{secretsCRE},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	c := cfg.CRE()
	assert.Equal(t, "streams-api-key", c.StreamsApiKey())
	assert.Equal(t, "streams-api-secret", c.StreamsApiSecret())
}

func TestEmptyCREConfig(t *testing.T) {
	cre := toml.CreSecrets{}
	cfg := creConfig{c: cre}
	assert.Equal(t, "", cfg.StreamsApiKey())
	assert.Equal(t, "", cfg.StreamsApiSecret())
}
