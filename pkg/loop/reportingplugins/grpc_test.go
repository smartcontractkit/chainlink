package reportingplugins_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	errorlogtest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog/test"
	pipelinetest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline/test"
	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/reportingplugin/ocr2/test"
	telemetrytest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry/test"
	nettest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net/test"
	reportingplugintest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func newStopCh(t *testing.T) <-chan struct{} {
	stopCh := make(chan struct{})
	if d, ok := t.Deadline(); ok {
		time.AfterFunc(time.Until(d), func() { close(stopCh) })
	}
	return stopCh
}

func PluginGenericTest(t *testing.T, p types.ReportingPluginClient) {
	t.Run("PluginServer", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewReportingPluginFactory(ctx,
			types.ReportingPluginServiceConfig{},
			nettest.MockConn{},
			pipelinetest.PipelineRunner,
			telemetrytest.Telemetry,
			errorlogtest.ErrorLog)
		require.NoError(t, err)

		reportingplugintest.RunFactory(t, factory)
	})
}

func TestGRPCService_MedianProvider(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(
		t,
		ocr2test.MedianID,
		&reportingplugins.GRPCService[types.MedianProvider]{
			PluginServer: ocr2test.MedianProviderServer,
			BrokerConfig: loop.BrokerConfig{
				Logger: logger.Test(t),
				StopCh: stopCh,
			},
		},
		PluginGenericTest,
	)
}

func TestGRPCService_PluginProvider(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(
		t,
		reportingplugins.PluginServiceName,
		&reportingplugins.GRPCService[types.PluginProvider]{
			PluginServer: ocr2test.AgnosticProviderServer,
			BrokerConfig: loop.BrokerConfig{
				Logger: logger.Test(t),
				StopCh: stopCh,
			},
		},
		PluginGenericTest,
	)
}
