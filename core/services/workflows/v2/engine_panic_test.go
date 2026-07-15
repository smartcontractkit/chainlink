package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

func TestSafeModuleExecute_RethrowsPanic(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(context.Context, *sdkpb.ExecuteRequest, host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			panic("simulated host memory panic")
		})

	lggr, logs := logger.TestObserved(t, zapcore.DPanicLevel)
	e := &Engine{
		cfg:  &EngineConfig{Module: module},
		lggr: logger.Sugared(lggr),
	}

	require.PanicsWithValue(t, "simulated host memory panic", func() {
		_, _ = e.safeModuleExecute(t.Context(), &sdkpb.ExecuteRequest{}, nil)
	})

	require.NotEmpty(t, logs.FilterMessageSnippet("rethrowing for clean restart").All(),
		"panic must be logged with context before rethrow")
}
