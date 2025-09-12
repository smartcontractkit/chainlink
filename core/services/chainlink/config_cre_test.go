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

[CRE]
EnableDKGRecipient = true

[CRE.WorkflowFetcher]
URL = "http://workflow-server.example.com/workflows"
`

	configCREWithFileURL = `
[CRE.WorkflowFetcher]
URL = "file:///path/to/workflows"
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
	assert.True(t, c.EnableDKGRecipient())

	// Test the new WorkflowFetcher URL
	fetcher := c.WorkflowFetcher()
	assert.NotNil(t, fetcher)
	assert.Equal(t, "http://workflow-server.example.com/workflows", fetcher.URL())
}

func TestCREConfigWithFileURL(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{configCREWithFileURL},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	c := cfg.CRE()
	fetcher := c.WorkflowFetcher()
	assert.NotNil(t, fetcher)
	assert.Equal(t, "file:///path/to/workflows", fetcher.URL())
}

func TestEmptyCREConfig(t *testing.T) {
	cfg := creConfig{s: toml.CreSecrets{}, c: toml.CreConfig{}}
	assert.Equal(t, "", cfg.StreamsAPIKey())
	assert.Equal(t, "", cfg.StreamsAPISecret())
	assert.Equal(t, "", cfg.WsURL())
	assert.Equal(t, "", cfg.RestURL())

	// Test empty WorkflowFetcher
	fetcher := cfg.WorkflowFetcher()
	assert.Empty(t, fetcher)
	assert.Empty(t, fetcher.URL(), "Empty WorkflowFetcher should have empty URL")
}

func TestWorkflowFetcherConfig(t *testing.T) {
	testCases := []struct {
		name     string
		config   string
		expected string
	}{
		{
			name: "HTTP URL",
			config: `
[CRE.WorkflowFetcher]
URL = "http://example.com/workflows"
`,
			expected: "http://example.com/workflows",
		},
		{
			name: "HTTPS URL",
			config: `
[CRE.WorkflowFetcher]
URL = "https://secure.example.com/workflows"
`,
			expected: "https://secure.example.com/workflows",
		},
		{
			name: "File URL",
			config: `
[CRE.WorkflowFetcher]
URL = "file:///local/path/to/workflows"
`,
			expected: "file:///local/path/to/workflows",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := GeneralConfigOpts{
				ConfigStrings: []string{tc.config},
			}
			cfg, err := opts.New()
			require.NoError(t, err)

			c := cfg.CRE()
			fetcher := c.WorkflowFetcher()
			assert.NotNil(t, fetcher)
			assert.Equal(t, tc.expected, fetcher.URL())
		})
	}
}
