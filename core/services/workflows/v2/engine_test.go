package v2_test

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

func TestEngine_Init(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)

	initDoneCh := make(chan error)

	cfg := v2.EngineConfig{
		Lggr:            logger.TestLogger(t),
		Module:          module,
		CapRegistry:     capreg,
		ExecutionsStore: store.NewInMemoryStore(logger.TestLogger(t), clockwork.NewRealClock()),
		WorkflowID:      "test-workflow",
		Limits:          v2.EngineLimits{},
		Hooks: v2.LifecycleHooks{
			OnInitialized: func(err error) {
				initDoneCh <- err
			},
		},
	}
	engine, err := v2.NewEngine(t.Context(), cfg)
	require.NoError(t, err)

	module.EXPECT().Start().Once()
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(capabilities.Node{}, nil).Once()
	require.NoError(t, engine.Start(t.Context()))

	require.NoError(t, <-initDoneCh)

	module.EXPECT().Close().Once()
	require.NoError(t, engine.Close())
}
