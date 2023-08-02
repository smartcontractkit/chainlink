package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

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

type databaseConfig struct {
	config.Database
	logSQL bool
}

func (d *databaseConfig) LogSQL() bool { return d.logSQL }

func TestResolver_SQLLogging(t *testing.T) {
	t.Parallel()

	query := `
		query GetSQLLogging {
			sqlLogging {
				... on SQLLogging {
					enabled
				}
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "sqlLogging"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("Database").Return(&databaseConfig{logSQL: false})
				f.App.On("GetConfig").Return(f.Mocks.cfg)
			},
			query: query,
			result: `
				{
					"sqlLogging": {
						"enabled": false
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}

type log struct {
	config.Log
	level zapcore.Level
}

func (l *log) Level() zapcore.Level {
	return l.level
}

func TestResolver_GlobalLogLevel(t *testing.T) {
	t.Parallel()

	query := `
		query GetGlobalLogLevel {
			globalLogLevel {
				... on GlobalLogLevel {
					level
				}
			}
		}`

	var warnLvl zapcore.Level
	err := warnLvl.UnmarshalText([]byte("warn"))
	assert.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "globalLogLevel"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("Log").Return(&log{level: warnLvl})
				f.App.On("GetConfig").Return(f.Mocks.cfg)
			},
			query: query,
			result: `
				{
					"globalLogLevel": {
						"level": "WARN"
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_SetGlobalLogLevel(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation SetGlobalLogLevel($level: LogLevel!) {
			setGlobalLogLevel(level: $level) {
				... on SetGlobalLogLevelSuccess {
					globalLogLevel {
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
	variables := map[string]interface{}{
		"level": LogLevelError,
	}

	var errorLvl zapcore.Level
	err := errorLvl.UnmarshalText([]byte("error"))
	assert.NoError(t, err)

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "setGlobalLogLevel"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("SetLogLevel", errorLvl).Return(nil)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"setGlobalLogLevel": {
						"globalLogLevel": {
							"level": "ERROR"
						}
					}
				}`,
		},
		{
			name:          "generic error on SetLogLevel",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("SetLogLevel", errorLvl).Return(gError)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"setGlobalLogLevel"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
