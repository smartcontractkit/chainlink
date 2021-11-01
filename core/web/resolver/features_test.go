package resolver

import (
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/stretchr/testify/mock"
	"testing"

	configMocks "github.com/smartcontractkit/chainlink/core/config/mocks"
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

	f := setupFramework(t)
	cfg := &configMocks.GeneralConfig{}
	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, cfg)
	})

	f.App.On("GetConfig").Return(cfg)
	cfg.On("FeatureUICSAKeys").Return(false)
	cfg.On("FeatureUIFeedsManager").Return(true)

	gqltesting.RunTest(t, &gqltesting.Test{
		Context: f.Ctx,
		Schema:  f.RootSchema,
		Query:   query,
		ExpectedResult: `
		{
			"features": {
				"csa": false,
				"feedsManager": true
			}
		}`,
	})
}
