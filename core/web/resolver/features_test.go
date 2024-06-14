package resolver

import (
	"testing"

	configtest2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func Test_ToFeatures(t *testing.T) {
	query := `
	{
		features {
			... on Features {
				csa
				feedsManager
			}	
		}
	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "features"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetConfig").Return(configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
					t, f := true, false
					c.Feature.UICSAKeys = &f
					c.Feature.FeedsManager = &t
				}))
			},
			query: query,
			result: `
			{
				"features": {
					"csa": false,
					"feedsManager": true
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
