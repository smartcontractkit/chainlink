package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestResolver_SetServiceLogLevel(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation SetServicesLogLevels($input: SetServicesLogLevelsInput!) {
			setServicesLogLevels(input: $input) {
				... on SetServicesLogLevelsSuccess {
					config {
						keeper
						headTracker
						fluxMonitor
					}
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`
	input := map[string]interface{}{
		"input": map[string]interface{}{
			"config": map[string]interface{}{
				"headTracker": "INFO",
			},
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: input}, "setServicesLogLevels"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				var lvl zapcore.Level
				err := lvl.UnmarshalText([]byte("info"))
				assert.NoError(t, err)

				f.App.On("SetServiceLogLevel", mock.Anything, logger.HeadTracker, lvl).Return(nil)
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"setServicesLogLevels": {
						"config": {
							"headTracker": "INFO",
							"keeper": null,
							"fluxMonitor": null
						}
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
