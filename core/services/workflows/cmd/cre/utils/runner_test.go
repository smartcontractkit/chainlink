package utils

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
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	t.Run("happy path with an empty workflow", func(t *testing.T) {
		t.Parallel()

		// Build before deadline; WASM compile can exceed 5s under CI load.
		binary := wasmtest.CreateTestBinary(t, filepath.Join("core/services/workflows/cmd/cre/examples/v2", "empty"), false)

		duration := 5 * time.Second
		ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(duration))
		defer cancel()

		hooks := DefaultHooks()
		hooks.Finally = func(ctx context.Context, cfg RunnerConfig, registry *capabilities.Registry, svcs []services.Service) {
			for _, service := range svcs {
				err := service.Ready()
				require.ErrorContains(t, err, "Stopped")
			}
		}

		runner := NewRunner(hooks)
		runner.Run(ctx, "", binary, []byte{}, []byte{}, RunnerConfig{
			EnableBeholder:             false,
			EnableBilling:              true,
			EnableStandardCapabilities: false,
			Lggr:                       logger.TestLogger(t),
			LifecycleHooks:             v2.LifecycleHooks{},
		})
	})
}
