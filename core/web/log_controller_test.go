package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type testCase struct {
	Description string
	logLevel    string
	logSql      *bool
	svcName     string
	svcLevel    string

	expectedLogLevel  zapcore.Level
	expectedLogSQL    bool
	expectedSvcLevel  map[string]zapcore.Level
	expectedErrorCode int
}

func TestLogController_GetLogConfig(t *testing.T) {
	t.Parallel()

	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	logLevel := zapcore.WarnLevel
	cfg.Overrides.LogLevel = &logLevel
	sqlEnabled := true
	cfg.Overrides.LogSQL = null.BoolFrom(sqlEnabled)
	defaultLogLevel := zapcore.WarnLevel
	cfg.Overrides.DefaultLogLevel = &defaultLogLevel

	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()

	resp, err := client.HTTPClient.Get("/v2/log")
	require.NoError(t, err)

	svcLogConfig := presenters.ServiceLogConfigResource{}
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &svcLogConfig))

	require.Equal(t, "warn", svcLogConfig.DefaultLogLevel)

	for i, svcName := range svcLogConfig.ServiceName {

		if svcName == "Global" {
			assert.Equal(t, logLevel.String(), svcLogConfig.LogLevel[i])
		}

		if svcName == "IsSqlEnabled" {
			assert.Equal(t, strconv.FormatBool(sqlEnabled), svcLogConfig.LogLevel[i])
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
		{
			Description: "Set head tracker to info",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.HeadTracker}, ","),
			svcLevel: strings.Join([]string{"info"}, ","),

			expectedLogLevel: zapcore.InfoLevel,
			expectedSvcLevel: map[string]zapcore.Level{logger.HeadTracker: zapcore.InfoLevel},
		},
		{
			Description: "Set flux monitor to warn",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.FluxMonitor}, ","),
			svcLevel: strings.Join([]string{"warn"}, ","),

			expectedLogLevel: zapcore.InfoLevel,
			expectedSvcLevel: map[string]zapcore.Level{logger.FluxMonitor: zapcore.WarnLevel},
		},
		{
			Description: "Set keeper to info",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.Keeper}, ","),
			svcLevel: strings.Join([]string{"info"}, ","),

			expectedLogLevel: zapcore.InfoLevel,
			expectedSvcLevel: map[string]zapcore.Level{logger.Keeper: zapcore.InfoLevel},
		},
		{
			Description: "Set multiple services log levels",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.HeadTracker, logger.FluxMonitor, logger.Keeper}, ","),
			svcLevel: strings.Join([]string{"debug", "warn", "info"}, ","),

			expectedLogLevel: zapcore.InfoLevel,
			expectedSvcLevel: map[string]zapcore.Level{
				logger.HeadTracker: zapcore.DebugLevel,
				logger.FluxMonitor: zapcore.WarnLevel,
				logger.Keeper:      zapcore.InfoLevel,
			},
		},
		{
			Description: "Set incorrect log levels",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.HeadTracker, logger.FluxMonitor, logger.Keeper}, ","),
			svcLevel: strings.Join([]string{"debug", "warning", "infa"}, ","),

			expectedErrorCode: http.StatusBadRequest,
		},
		{
			Description: "Set incorrect service names",
			logLevel:    "info",

			svcName:  strings.Join([]string{logger.HeadTracker, "FLUX-MONITOR", "SHKEEPER"}, ","),
			svcLevel: strings.Join([]string{"debug", "warning", "info"}, ","),

			expectedErrorCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			app := cltest.NewApplicationEVMDisabled(t)
			require.NoError(t, app.Start(testutils.Context(t)))
			client := app.NewHTTPClient()

			request := web.LogPatchRequest{Level: tc.logLevel, SqlEnabled: tc.logSql}

			if tc.svcName != "" {
				svcs := strings.Split(tc.svcName, ",")
				lvls := strings.Split(tc.svcLevel, ",")

				serviceLogLevel := make([][2]string, len(svcs))

				for i, p := range svcs {
					serviceLogLevel[i][0] = p
					serviceLogLevel[i][1] = lvls[i]
				}
				request.ServiceLogLevel = serviceLogLevel
			}

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

					if svcName == logger.HeadTracker {
						assert.Equal(t, tc.expectedSvcLevel[logger.HeadTracker].String(), svcLogConfig.LogLevel[i])
					}

					if svcName == logger.FluxMonitor {
						assert.Equal(t, tc.expectedSvcLevel[logger.FluxMonitor].String(), svcLogConfig.LogLevel[i])
					}
				}
			}
		})
	}
}
