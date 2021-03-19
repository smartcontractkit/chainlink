package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type testCase struct {
	Description string
	logLevel    string
	logSql      string

	expectedLogLevel zapcore.Level
	expectedLogSql   bool
}

func TestLogController_SetDebug(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	cases := []testCase{
		{
			Description:      "Set log level to debug",
			logLevel:         "debug",
			logSql:           "",
			expectedLogLevel: zapcore.DebugLevel,
		},
		{
			Description:      "Set log level to info",
			logLevel:         "info",
			logSql:           "",
			expectedLogLevel: zapcore.InfoLevel,
		},
		{
			Description:      "Set log level to info and log sql to true",
			logLevel:         "info",
			logSql:           "true",
			expectedLogLevel: zapcore.InfoLevel,
		},
		{
			Description:      "Set log level to warn and log sql to false",
			logLevel:         "warn",
			logSql:           "false",
			expectedLogLevel: zapcore.WarnLevel,
		},
	}

	for _, tc := range cases {
		func() {
			request := web.LogPatchRequest{Level: tc.logLevel, SqlEnabled: tc.logSql}

			requestData, _ := json.Marshal(request)
			buf := bytes.NewBuffer(requestData)

			resp, cleanup := client.Patch("/v2/log", buf)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, http.StatusOK)

			lR := presenters.LogResource{}
			require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &lR))
			if tc.logLevel != "" {
				assert.Equal(t, tc.logLevel, lR.Level)
			}
			if tc.logSql != "" {
				assert.Equal(t, tc.logSql, strconv.FormatBool(lR.SqlEnabled))
			}
			assert.Equal(t, tc.expectedLogLevel.String(), app.GetStore().Config.LogLevel().String())
		}()
	}
}
