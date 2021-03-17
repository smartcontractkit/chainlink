package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	Description      string
	enableDug        bool
	expectedLogLevel zapcore.Level
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
			Description:      "Set debug enabled to true",
			enableDug:        true,
			expectedLogLevel: zapcore.DebugLevel,
		},
		{
			Description:      "Set debug enabled to false (info)",
			enableDug:        false,
			expectedLogLevel: zapcore.InfoLevel,
		},
	}

	for _, tc := range cases {
		func() {
			request := web.LoglevelPatchRequest{EnableDebugLog: &tc.enableDug}
			requestData, _ := json.Marshal(request)
			buf := bytes.NewBuffer(requestData)

			resp, cleanup := client.Patch("/v2/log", buf)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, http.StatusOK)

			lR := presenters.LogResource{}
			require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &lR))
			assert.Equal(t, tc.enableDug, lR.DebugEnabled)
			assert.Equal(t, tc.expectedLogLevel.String(), app.GetStore().Config.LogLevel().String())
		}()
	}
}
