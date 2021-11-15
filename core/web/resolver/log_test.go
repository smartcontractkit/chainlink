package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestResolver_SetServiceLogLevel(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation SetServicesLogLevels($input: SetServicesLogLevelsInput!) {
			setServicesLogLevels(input: $input) {
				... on SetServicesLogLevelsSuccess {
					logLevels {
						name
						level
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
			"logLevels": []map[string]interface{}{
				{
					"name":  logger.HeadTracker,
					"level": "DEBUG",
				},
			},
		},
	}
	invalidInput := map[string]interface{}{
		"input": map[string]interface{}{
			"logLevels": []map[string]interface{}{
				{
					"name":  logger.HeadTracker,
					"level": "invalid",
				},
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
				err := lvl.UnmarshalText([]byte("debug"))
				assert.NoError(t, err)

				f.App.On("SetServiceLogLevel", f.Ctx, logger.HeadTracker, lvl)
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"setServicesLogLevels": {
						"logLevels": [
							{
								"name": "head_tracker",
								"level": "debug"
							}
						]
					}
				}`,
		},
		{
			name:          "invalid log level",
			authenticated: true,
			query:         mutation,
			variables:     invalidInput,
			result: `
				{
					"setServicesLogLevels": {
						"errors": [
							{
								"path": "head_tracker/invalid",
								"message": "invalid log level",
								"code": "INVALID_INPUT"
							}
						]
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
