package v2_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	beholderpb "github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/protoc/pkg/test_capabilities/basicaction"
	basicactionmock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/protoc/pkg/test_capabilities/basicaction/basic_actionmock"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/protoc/pkg/test_capabilities/basictrigger"
	basictriggermock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/protoc/pkg/test_capabilities/basictrigger/basic_triggermock"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/testutils/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2/pb"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
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
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	initDoneCh := make(chan error)

	cfg := defaultTestConfig(t)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
	}

	engine, err := v2.NewEngine(testutils.Context(t), cfg)
	require.NoError(t, err)

	module.EXPECT().Start().Once()
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(0), nil).Once()
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
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(0), nil).Times(2)
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
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
		engine1, err = v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		require.NoError(t, engine1.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 2 gets rate-limited by per-owner limit", func(t *testing.T) {
		engine2, err = v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		require.NoError(t, engine2.Start(t.Context()))
		initErr := <-initDoneCh
		require.Equal(t, types.ErrPerOwnerWorkflowCountLimitReached, initErr)
	})

	t.Run("engine 3 inits successfully", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerB
		engine3, err = v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		require.NoError(t, engine3.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 4 gets rate-limited by global limit", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerC
		engine4, err = v2.NewEngine(testutils.Context(t), cfg)
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
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

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
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.ErrorContains(t, <-initDoneCh, "too many trigger subscriptions")
		require.NoError(t, engine.Close())
		cfg.LocalLimits.MaxTriggerSubscriptions = 10
	})

	t.Run("trigger capability not found in the registry", func(t *testing.T) {
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(nil, errors.New("not found")).Once()
		require.NoError(t, engine.Start(t.Context()))
		require.ErrorContains(t, <-initDoneCh, "trigger capability not found")
		require.NoError(t, engine.Close())
	})

	t.Run("successful trigger registration", func(t *testing.T) {
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
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
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
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

func TestEngine_Execution(t *testing.T) {
	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

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
		OnExecutionFinished: func(executionID string) {
			executionFinishedCh <- executionID
		},
	}
	beholderTester := tests.Beholder(t)
	cfg.BeholderEmitter = custmsg.NewLabeler()

	t.Run("successful execution with no capability calls", func(t *testing.T) {
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
		trigger := capmocks.NewTriggerCapability(t)
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
		eventCh := make(chan capabilities.TriggerResponse)
		trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
		trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

		require.NoError(t, engine.Start(t.Context()))

		require.NoError(t, <-initDoneCh) // successful trigger registration
		require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "event_012345",
			Payload:     nil,
		}

		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(
				func(_ context.Context, request *wasmpb.ExecuteRequest, executor host.CapabilityExecutor) {
					wantExecID, err := types.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
					require.NoError(t, err)
					capExec, ok := executor.(*v2.CapabilityExecutor)
					require.True(t, ok)
					require.Equal(t, wantExecID, capExec.WorkflowExecutionID)
					require.Equal(t, uint64(0), request.Request.(*wasmpb.ExecuteRequest_Trigger).Trigger.Id)
				},
			).
			Return(nil, nil).
			Once()

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}
		<-executionFinishedCh

		require.NoError(t, engine.Close())

		requireEventsLabels(t, &beholderTester, map[string]string{
			"workflowID":    cfg.WorkflowID,
			"workflowOwner": cfg.WorkflowOwner,
			"workflowName":  cfg.WorkflowName.String(),
		})
		requireEventsMessages(t, &beholderTester, []string{
			"Started",
			"Registering trigger",
			"All triggers registered successfully",
			"Workflow Engine initialized",
			"Workflow execution finished successfully",
		})
	})
}

func TestEngine_MockCapabilityRegistry_NoDAGBinary(t *testing.T) {
	cmd := "core/services/workflows/test/wasm/v2/cmd"
	log := logger.TestLogger(t)
	binaryB := wasmtest.CreateTestBinary(cmd, false, t)
	module, err := host.NewModule(&host.ModuleConfig{
		Logger:         log,
		IsUncompressed: true,
	}, binaryB)
	require.NoError(t, err)

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	cfg := defaultTestConfig(t)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = new(mockBillingClient)

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	resultReceivedCh := make(chan *wasmpb.ExecutionResult, 1)
	executionFinishedCh := make(chan string, 1)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string) {
			executionFinishedCh <- executionID
		},
		OnResultReceived: func(er *wasmpb.ExecutionResult) {
			resultReceivedCh <- er
		},
	}

	triggerMock, basicActionMock := setupExpectedCalls(t)
	wrappedTriggerMock := &registry.CapabilityWrapper{
		Capability: triggerMock,
	}
	wrappedActionMock := &registry.CapabilityWrapper{
		Capability: basicActionMock,
	}

	t.Run("OK happy path", func(t *testing.T) {
		wantResponse := "Hello, world!"
		engine, err := v2.NewEngine(testutils.Context(t), cfg)
		require.NoError(t, err)

		capreg.EXPECT().
			GetTrigger(matches.AnyContext, wrappedTriggerMock.ID()).
			Return(wrappedTriggerMock, nil).
			Once()

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, wrappedActionMock.ID()).
			Return(wrappedActionMock, nil).
			Twice()

		require.NoError(t, engine.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
		require.Equal(t, []string{wrappedTriggerMock.ID()}, <-subscribedToTriggersCh)

		// Read the result from the hook and assert that the wanted response was
		// received.
		res := <-resultReceivedCh
		switch output := res.Result.(type) {
		case *wasmpb.ExecutionResult_Value:
			var value values.Value
			var execErr error
			var unwrapped any

			valuePb := output.Value
			value, execErr = values.FromProto(valuePb)
			require.NoError(t, execErr)
			unwrapped, execErr = value.Unwrap()
			require.NoError(t, execErr)
			require.Equal(t, wantResponse, unwrapped)
		default:
			t.Fatalf("unexpected response type %T", output)
		}

		execID, err := types.GenerateExecutionID(cfg.WorkflowID, "")
		require.NoError(t, err)

		require.Equal(t, execID, <-executionFinishedCh)
		require.NoError(t, engine.Close())
	})
}

// setupExpectedCalls mocks single call to trigger and two calls to the basic action
// mock capability
func setupExpectedCalls(t *testing.T) (
	*basictriggermock.BasicCapability,
	*basicactionmock.BasicActionCapability,
) {
	triggerMock := &basictriggermock.BasicCapability{}
	triggerMock.Trigger = func(ctx context.Context, input *basictrigger.Config) (*basictrigger.Outputs, error) {
		return &basictrigger.Outputs{CoolOutput: "Hello, "}, nil
	}

	basicAction := &basicactionmock.BasicActionCapability{}

	firstCall := true
	callLock := &sync.Mutex{}
	basicAction.PerformAction = func(ctx context.Context, input *basicaction.Inputs) (*basicaction.Outputs, error) {
		callLock.Lock()
		defer callLock.Unlock()
		assert.NotEqual(t, firstCall, input.InputThing, "failed first call assertion")
		firstCall = false
		if input.InputThing {
			return &basicaction.Outputs{AdaptedThing: "!"}, nil
		}
		return &basicaction.Outputs{AdaptedThing: "world"}, nil
	}
	return triggerMock, basicAction
}

func requireEventsLabels(t *testing.T, beholderTester *tests.BeholderTester, want map[string]string) {
	msgs := beholderTester.Messages(t)
	for _, msg := range msgs {
		if msg.Attrs["beholder_entity"] == "BaseMessage" {
			var payload beholderpb.BaseMessage
			require.NoError(t, proto.Unmarshal(msg.Body, &payload))
			for k, v := range want {
				require.Equal(t, v, payload.Labels[k], "label %s does not match", k)
			}
		}
	}
}

func requireEventsMessages(t *testing.T, beholderTester *tests.BeholderTester, expected []string) {
	msgs := beholderTester.Messages(t)
	nextToFind := 0
	for _, msg := range msgs {
		if msg.Attrs["beholder_entity"] == "BaseMessage" {
			var payload beholderpb.BaseMessage
			require.NoError(t, proto.Unmarshal(msg.Body, &payload))
			if nextToFind >= len(expected) {
				return
			}
			if payload.Msg == expected[nextToFind] {
				nextToFind++
			}
		}
	}

	if nextToFind < len(expected) {
		t.Errorf("log message not found: %s", expected[nextToFind])
	}
}

func newNode(t *testing.T) capabilities.Node {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return capabilities.Node{
		PeerID: &peerID,
	}
}
