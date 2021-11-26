package resolver

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver_Config(t *testing.T) {
	t.Parallel()

	query := `
		query GetConfiguration {
			config {
				... on Config {
					allowOrigins
					clientNodeURL
					ethereumSecondaryURLs
					replayFromBlock
				}
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "config"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				parsedURL, err := url.Parse("https://test-url.com")
				assert.NoError(t, err)

				f.Mocks.cfg.On("AllowOrigins").Return("test")
				f.Mocks.cfg.On("ClientNodeURL").Return("some-url")
				f.Mocks.cfg.On("EthereumSecondaryURLs").Return([]url.URL{*parsedURL})
				f.Mocks.cfg.On("ReplayFromBlock").Return(int64(12))
				f.App.On("GetConfig").Return(f.Mocks.cfg)
			},
			query: query,
			result: `
				{
					"config": {
						"allowOrigins": "test",
						"clientNodeURL": "some-url",
						"ethereumSecondaryURLs": ["https://test-url.com"],
						"replayFromBlock": 12
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
