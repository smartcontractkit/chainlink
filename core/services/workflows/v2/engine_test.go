package v2_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2/pb"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
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
	module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(0), nil).Once()
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
	module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(0), nil).Times(2)
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

func TestEngine_TriggerSubscriptions(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(capabilities.Node{}, nil)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)

	cfg := defaultTestConfig(t)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
	}

	t.Run("too many triggers", func(t *testing.T) {
		cfg.LocalLimits.MaxTriggerSubscriptions = 1
		engine, err := v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.ErrorContains(t, <-initDoneCh, "too many trigger subscriptions")
		require.NoError(t, engine.Close())
		cfg.LocalLimits.MaxTriggerSubscriptions = 10
	})

	t.Run("trigger capability not found in the registry", func(t *testing.T) {
		engine, err := v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(nil, errors.New("not found")).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.ErrorContains(t, <-initDoneCh, "trigger capability not found")
		require.NoError(t, engine.Close())
	})

	t.Run("successful trigger registration", func(t *testing.T) {
		engine, err := v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		trigger0, trigger1 := capmocks.NewTriggerCapability(t), capmocks.NewTriggerCapability(t)
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger0, nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_1").Return(trigger1, nil).Once()
		tr0Ch, tr1Ch := make(chan capabilities.TriggerResponse), make(chan capabilities.TriggerResponse)
		trigger0.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(tr0Ch, nil).Once()
		trigger1.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(tr1Ch, nil).Once()
		trigger0.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
		trigger1.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
		require.Equal(t, []string{"id_0", "id_1"}, <-subscribedToTriggersCh)
		require.NoError(t, engine.Close())
	})

	t.Run("failed trigger registration and rollback", func(t *testing.T) {
		engine, err := v2.NewEngine(t.Context(), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		trigger0, trigger1 := capmocks.NewTriggerCapability(t), capmocks.NewTriggerCapability(t)
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger0, nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_1").Return(trigger1, nil).Once()
		tr0Ch := make(chan capabilities.TriggerResponse)
		trigger0.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(tr0Ch, nil).Once()
		trigger1.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(nil, errors.New("failure ABC")).Once()
		trigger0.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.ErrorContains(t, <-initDoneCh, "failed to register trigger: failure ABC")
		require.NoError(t, engine.Close())
	})
}

func newTriggerSubs(n int) *wasmpb.ExecutionResult {
	subs := make([]*sdkpb.TriggerSubscription, 0, n)
	for i := range n {
		subs = append(subs, &sdkpb.TriggerSubscription{
			ExecId: fmt.Sprintf("execId_%d", i),
			Id:     fmt.Sprintf("id_%d", i),
			Method: "method",
		})
	}
	return &wasmpb.ExecutionResult{
		Result: &wasmpb.ExecutionResult_TriggerSubscriptions{
			TriggerSubscriptions: &sdkpb.TriggerSubscriptionRequest{
				Subscriptions: subs,
			},
		},
	}
}
