package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

func TestNewStartUpHealthReport(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	ibhr := services.NewStartUpHealthReport(1234, lggr)

	ibhr.Start()
	require.Eventually(t, func() bool { return observed.Len() >= 1 }, time.Second*5, time.Millisecond*100)
	require.Equal(t, "Starting StartUpHealthReport", observed.TakeAll()[0].Message)

	req, err := http.NewRequestWithContext(tests.Context(t), "GET", "http://localhost:1234/health", nil)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	ibhr.Stop()
	require.Eventually(t, func() bool { return observed.Len() >= 1 }, time.Second*5, time.Millisecond*100)
	require.Equal(t, "StartUpHealthReport shutdown complete", observed.TakeAll()[0].Message)
}
