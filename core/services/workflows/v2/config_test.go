package v2_test

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func TestEngineConfig_Validate(t *testing.T) {
	t.Parallel()
	cfg := &v2.EngineConfig{
		Lggr:            logger.TestLogger(t),
		Module:          nil,
		CapRegistry:     regmocks.NewCapabilitiesRegistry(t),
		ExecutionsStore: store.NewInMemoryStore(logger.TestLogger(t), clockwork.NewRealClock()),
		WorkflowID:      "test-workflow",
	}
	t.Run("nil module", func(t *testing.T) {
		require.Error(t, cfg.Validate())
	})

	t.Run("success", func(t *testing.T) {
		cfg.Module = modulemocks.NewModuleV2(t)
		require.NoError(t, cfg.Validate())
		require.NotEqual(t, 0, cfg.Limits.CapRegistryAccessRetryIntervalMs)
		require.NotNil(t, cfg.Hooks.OnInitialized)
	})
}
