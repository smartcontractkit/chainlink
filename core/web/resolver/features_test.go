package resolver

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func Test_ToFeatures(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	query := `
	{
		features {
			... on Features {
				csa
				feedsManager
				multiFeedsManagers
			}	
		}
	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "features"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetConfig").Return(configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
					t, f := true, false
					c.Feature.UICSAKeys = &f
					c.Feature.FeedsManager = &t
					c.Feature.MultiFeedsManagers = &f
				}))
			},
			query: query,
			result: `
			{
				"features": {
					"csa": false,
					"feedsManager": true,
					"multiFeedsManagers": false
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
