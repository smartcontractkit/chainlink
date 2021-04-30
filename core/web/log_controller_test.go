package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type testCase struct {
	Description string
	logLevel    string
	logSql      *bool

	expectedLogLevel  zapcore.Level
	expectedLogSQL    bool
	expectedErrorCode int
}

func TestLogController_GetLogConfig(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)

	// Set log config values
	logLevel := "warn"
	sqlEnabled := true
	app.GetStore().Config.Set("LOG_LEVEL", logLevel)
	app.GetStore().Config.Set("LOG_SQL", sqlEnabled)

	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	resp, err := client.HTTPClient.Get("/v2/log")
	require.NoError(t, err)

	lR := presenters.LogResource{}
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &lR))

	assert.Equal(t, lR.SqlEnabled, sqlEnabled)
	assert.Equal(t, lR.Level, logLevel)
}

func TestLogController_PatchLogConfig(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	sqlTrue := true
	sqlFalse := false
	cases := []testCase{
		{
			Description:      "Set log level to debug",
			logLevel:         "debug",
			logSql:           nil,
			expectedLogLevel: zapcore.DebugLevel,
		},
		{
			Description:      "Set log level to info",
			logLevel:         "info",
			logSql:           nil,
			expectedLogLevel: zapcore.InfoLevel,
		},
		{
			Description:      "Set log level to info and log sql to true",
			logLevel:         "info",
			logSql:           &sqlTrue,
			expectedLogLevel: zapcore.InfoLevel,
			expectedLogSQL:   true,
		},
		{
			Description:      "Set log level to warn and log sql to false",
			logLevel:         "warn",
			logSql:           &sqlFalse,
			expectedLogLevel: zapcore.WarnLevel,
			expectedLogSQL:   false,
		},
		{
			Description:       "Send no params to updater",
			expectedErrorCode: http.StatusBadRequest,
		},
		{
			Description:       "Send bad log level request",
			logLevel:          "test",
			expectedErrorCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		request := web.LogPatchRequest{Level: tc.logLevel, SqlEnabled: tc.logSql}

		requestData, _ := json.Marshal(request)
		buf := bytes.NewBuffer(requestData)

		resp, cleanup := client.Patch("/v2/log", buf)
		defer cleanup()

		lR := presenters.LogResource{}
		if tc.expectedErrorCode != 0 {
			cltest.AssertServerResponse(t, resp, tc.expectedErrorCode)
		} else {
			cltest.AssertServerResponse(t, resp, http.StatusOK)
			require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &lR))
			if tc.logLevel != "" {
				assert.Equal(t, tc.expectedLogLevel.String(), lR.Level)
			}
			if tc.logSql != nil {
				assert.Equal(t, tc.logSql, &lR.SqlEnabled)
				assert.Equal(t, &tc.expectedLogSQL, &lR.SqlEnabled)
			}
			assert.Equal(t, tc.expectedLogLevel.String(), app.GetStore().Config.LogLevel().String())
		}
	}
}
