package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
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
				"fluxMonitor": "WARN",
			},
		},
	}

	gError := errors.New("error")

	var infoLvl zapcore.Level
	err := infoLvl.UnmarshalText([]byte("info"))
	assert.NoError(t, err)

	var warnLvl zapcore.Level
	err = warnLvl.UnmarshalText([]byte("warn"))
	assert.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: input}, "setServicesLogLevels"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("SetServiceLogLevel", mock.Anything, logger.HeadTracker, infoLvl).Return(nil)
				f.App.On("SetServiceLogLevel", mock.Anything, logger.FluxMonitor, warnLvl).Return(nil)
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"setServicesLogLevels": {
						"config": {
							"headTracker": "INFO",
							"keeper": null,
							"fluxMonitor": "WARN"
						}
					}
				}`,
		},
		{
			name:          "general service log level error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("SetServiceLogLevel", mock.Anything, logger.HeadTracker, infoLvl).Return(nil)
				f.App.On("SetServiceLogLevel", mock.Anything, logger.FluxMonitor, warnLvl).Return(gError)
			},
			query:     mutation,
			variables: input,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"setServicesLogLevels"},
					Message:       "error",
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_SetSQLLogging(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation SetSQLLogging($input: SetSQLLoggingInput!) {
			setSQLLogging(input: $input) {
				... on SetSQLLoggingSuccess {
					sqlLogging {
						enabled
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"enabled": true,
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "setSQLLogging"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("SetLogSQL", true).Return(nil)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"setSQLLogging": {
						"sqlLogging": {
							"enabled": true
						}
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
