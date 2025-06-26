package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	t.Run("happy path with an empty workflow", func(t *testing.T) {
		t.Parallel()

		duration := 5 * time.Second
		ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(duration))
		defer cancel()

		hooks := DefaultHooks()
		hooks.Finally = func(ctx context.Context, cfg RunnerConfig, registry *capabilities.Registry, svcs []services.Service) {
			for _, service := range svcs {
				require.ErrorContains(t, service.Ready(), "Stopped")
			}
		}

		binary := wasmtest.CreateTestBinary(filepath.Join("core/services/workflows/cmd/cre/examples/v2", "empty"), false, t)

		runner := NewRunner(hooks)
		runner.run(ctx, binary, []byte{}, RunnerConfig{
			enableBeholder:             false,
			enableBilling:              true,
			enableStandardCapabilities: false,
			lggr:                       logger.TestLogger(t),
		})
	})
}
