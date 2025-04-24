package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

func TestEngine_Init(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)

	initDoneCh := make(chan error)

	cfg := defaultTestConfig(t)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
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

func TestEngine_Start_RateLimited(t *testing.T) {
	t.Parallel()
	sLimiter, err := syncerlimiter.NewWorkflowLimits(logger.TestLogger(t), syncerlimiter.Config{
		Global:   2,
		PerOwner: 1,
	})
	require.NoError(t, err)

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(capabilities.Node{}, nil)
	initDoneCh := make(chan error)
	hooks := v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
	}

	cfg := defaultTestConfig(t)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.GlobalLimits = sLimiter
	cfg.Hooks = hooks
	var engine1, engine2, engine3, engine4 *v2.Engine

	t.Run("engine 1 inits successfully", func(t *testing.T) {
		engine1, err = v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		require.NoError(t, engine1.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 2 gets rate-limited by per-owner limit", func(t *testing.T) {
		engine2, err = v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		require.NoError(t, engine2.Start(t.Context()))
		initErr := <-initDoneCh
		require.Equal(t, types.ErrPerOwnerWorkflowCountLimitReached, initErr)
	})

	t.Run("engine 3 inits successfully", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerB
		engine3, err = v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		require.NoError(t, engine3.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 4 gets rate-limited by global limit", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerC
		engine4, err = v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		require.NoError(t, engine4.Start(t.Context()))
		initErr := <-initDoneCh
		require.Equal(t, types.ErrGlobalWorkflowCountLimitReached, initErr)
	})

	require.NoError(t, engine1.Close())
	require.NoError(t, engine2.Close())
	require.NoError(t, engine3.Close())
	require.NoError(t, engine4.Close())
}
