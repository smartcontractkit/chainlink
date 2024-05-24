package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
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

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Log.Level = ptr(toml.LogLevel(zapcore.WarnLevel))
		c.Database.LogQueries = ptr(true)
	})

	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	resp, clean := client.Get("/v2/log")
	t.Cleanup(clean)

	svcLogConfig := presenters.ServiceLogConfigResource{}
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &svcLogConfig))

	require.Equal(t, "warn", svcLogConfig.DefaultLogLevel)

	for i, svcName := range svcLogConfig.ServiceName {
		if svcName == "Global" {
			assert.Equal(t, zapcore.WarnLevel.String(), svcLogConfig.LogLevel[i])
		}

		if svcName == "IsSqlEnabled" {
			assert.Equal(t, strconv.FormatBool(true), svcLogConfig.LogLevel[i])
		}
	}
}

func TestLogController_PatchLogConfig(t *testing.T) {
	t.Parallel()

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
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			app := cltest.NewApplicationEVMDisabled(t)
			require.NoError(t, app.Start(testutils.Context(t)))
			client := app.NewHTTPClient(nil)

			request := web.LogPatchRequest{Level: tc.logLevel, SqlEnabled: tc.logSql}

			requestData, _ := json.Marshal(request)
			buf := bytes.NewBuffer(requestData)

			resp, cleanup := client.Patch("/v2/log", buf)
			defer cleanup()

			svcLogConfig := presenters.ServiceLogConfigResource{}
			if tc.expectedErrorCode != 0 {
				cltest.AssertServerResponse(t, resp, tc.expectedErrorCode)
			} else {
				cltest.AssertServerResponse(t, resp, http.StatusOK)
				require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &svcLogConfig))

				for i, svcName := range svcLogConfig.ServiceName {
					if svcName == "Global" {
						assert.Equal(t, tc.expectedLogLevel.String(), svcLogConfig.LogLevel[i])
					}

					if svcName == "IsSqlEnabled" {
						assert.Equal(t, strconv.FormatBool(tc.expectedLogSQL), svcLogConfig.LogLevel[i])
					}
				}
			}
		})
	}
}
