package v2_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/quarantine"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	beholderpb "github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultMock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault/mock"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	workflowEvents "github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/smartcontractkit/cre-sdk-go/cre/testutils/registry"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basicaction"
	basicactionmock "github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basicaction/mock"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

const triggerID = "basic-test-trigger@1.0.0"

func TestEngine_Init(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	initDoneCh := make(chan error)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
	}

	engine, err := v2.NewEngine(cfg)
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
	}, limits.Factory{})
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

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.GlobalExecutionConcurrencyLimiter = sLimiter
	cfg.Hooks = hooks
	var engine1, engine2, engine3, engine4 *v2.Engine

	t.Run("engine 1 inits successfully", func(t *testing.T) {
		engine1, err = v2.NewEngine(cfg)
		require.NoError(t, err)
		require.NoError(t, engine1.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 2 gets rate-limited by per-owner limit", func(t *testing.T) {
		engine2, err = v2.NewEngine(cfg)
		require.NoError(t, err)
		require.NoError(t, engine2.Start(t.Context()))
		initErr := <-initDoneCh
		require.Equal(t, types.ErrPerOwnerWorkflowCountLimitReached, initErr)
	})

	t.Run("engine 3 inits successfully", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerB
		engine3, err = v2.NewEngine(cfg)
		require.NoError(t, err)
		require.NoError(t, engine3.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
	})

	t.Run("engine 4 gets rate-limited by global limit", func(t *testing.T) {
		cfg.WorkflowOwner = testWorkflowOwnerC
		engine4, err = v2.NewEngine(cfg)
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

	cfg := defaultTestConfig(t, nil)
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
		cfg2 := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
			cfg.TriggerSubscriptionLimit.DefaultValue = 1
		})
		cfg2.Module = module
		cfg2.CapRegistry = capreg
		cfg2.Hooks = v2.LifecycleHooks{
			OnInitialized: func(err error) {
				initDoneCh <- err
			},
			OnSubscribedToTriggers: func(triggerIDs []string) {
				subscribedToTriggersCh <- triggerIDs
			},
		}
		engine, err := v2.NewEngine(cfg2)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		servicetest.Run(t, engine)
		var errLimited limits.ErrorBoundLimited[int]
		if assert.ErrorAs(t, <-initDoneCh, &errLimited) {
			assert.Equal(t, "PerWorkflow.TriggerSubscriptionLimit", errLimited.Key)
			assert.Equal(t, settings.ScopeWorkflow, errLimited.Scope)
			assert.Equal(t, "ffffaabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233", errLimited.Tenant)
			assert.Equal(t, 1, errLimited.Limit)
			assert.Equal(t, 2, errLimited.Amount)
		}
	})

	t.Run("trigger capability not found in the registry", func(t *testing.T) {
		engine, err := v2.NewEngine(cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(nil, errors.New("not found")).Once()
		servicetest.Run(t, engine)
		require.ErrorContains(t, <-initDoneCh, "trigger capability not found")
	})

	t.Run("successful trigger registration", func(t *testing.T) {
		engine, err := v2.NewEngine(cfg)
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
		servicetest.Run(t, engine)
		require.NoError(t, <-initDoneCh)
		require.Equal(t, []string{"id_0", "id_1"}, <-subscribedToTriggersCh)
	})

	t.Run("failed trigger registration and rollback", func(t *testing.T) {
		engine, err := v2.NewEngine(cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(2), nil).Once()
		trigger0, trigger1 := capmocks.NewTriggerCapability(t), capmocks.NewTriggerCapability(t)
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger0, nil).Once()
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_1").Return(trigger1, nil).Once()
		tr0Ch := make(chan capabilities.TriggerResponse)
		trigger0.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(tr0Ch, nil).Once()
		trigger1.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(nil, errors.New("failure ABC")).Once()
		trigger0.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
		servicetest.Run(t, engine)
		require.ErrorContains(t, <-initDoneCh, "failed to register trigger: failure ABC")
	})
}

func newTriggerSubs(n int) *sdkpb.ExecutionResult {
	subs := make([]*sdkpb.TriggerSubscription, 0, n)
	for i := range n {
		subs = append(subs, &sdkpb.TriggerSubscription{
			Id:     fmt.Sprintf("id_%d", i),
			Method: "method",
		})
	}
	return &sdkpb.ExecutionResult{
		Result: &sdkpb.ExecutionResult_TriggerSubscriptions{
			TriggerSubscriptions: &sdkpb.TriggerSubscriptionRequest{
				Subscriptions: subs,
			},
		},
	}
}

func TestEngine_OrganizationIdLogger(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	// Create mock org resolver
	mockOrgResolver := &mockOrgResolver{
		orgID: "test-org-123",
		err:   nil,
	}

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.OrgResolver = mockOrgResolver
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	// Setup trigger registration
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	// Mock execution that will retrieve organization ID
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
			// The execution should have called the org resolver and set the organizationID
			// We can't directly access the engine's internal state, but we know this ran successfully
			// if the execution completes without error
		}).
		Return(nil, nil).
		Once()

	// Trigger an execution
	mockTriggerEvent := capabilities.TriggerEvent{
		TriggerType: "basic-trigger@1.0.0",
		ID:          "test_org_id_event",
		Payload:     nil,
	}

	eventCh <- capabilities.TriggerResponse{
		Event: mockTriggerEvent,
	}

	// Wait for execution to finish
	executionID := <-executionFinishedCh
	wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
	require.NoError(t, err)
	require.Equal(t, wantExecID, executionID)

	// Verify that the org resolver was called
	require.True(t, mockOrgResolver.getCalled, "Expected org resolver Get method to be called")
	require.Equal(t, cfg.WorkflowOwner, mockOrgResolver.calledWithOwner, "Expected org resolver to be called with workflow owner")

	require.NoError(t, engine.Close())
}

func TestEngine_OrganizationIdLogger_OrgResolverFailure(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	// Create mock org resolver that returns an error
	mockOrgResolver := &mockOrgResolver{
		orgID: "",
		err:   errors.New("org resolver error"),
	}

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.OrgResolver = mockOrgResolver
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	// Setup trigger registration
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	// Mock execution - should still work even if org resolver fails
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(nil, nil).
		Once()

	// Trigger an execution
	mockTriggerEvent := capabilities.TriggerEvent{
		TriggerType: "basic-trigger@1.0.0",
		ID:          "test_org_id_failure_event",
		Payload:     nil,
	}

	eventCh <- capabilities.TriggerResponse{
		Event: mockTriggerEvent,
	}

	// Wait for execution to finish - should complete successfully even with org resolver failure
	executionID := <-executionFinishedCh
	wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
	require.NoError(t, err)
	require.Equal(t, wantExecID, executionID)

	// Verify that the org resolver was called even though it failed
	require.True(t, mockOrgResolver.getCalled, "Expected org resolver Get method to be called")
	require.Equal(t, cfg.WorkflowOwner, mockOrgResolver.calledWithOwner, "Expected org resolver to be called with workflow owner")

	require.NoError(t, engine.Close())
}

// mockOrgResolver is a test implementation of orgresolver.OrgResolver
type mockOrgResolver struct {
	orgID           string
	err             error
	getCalled       bool
	calledWithOwner string
}

func (m *mockOrgResolver) Get(ctx context.Context, owner string) (string, error) {
	m.getCalled = true
	m.calledWithOwner = owner
	return m.orgID, m.err
}

func (m *mockOrgResolver) Start(ctx context.Context) error {
	return nil
}

func (m *mockOrgResolver) Close() error {
	return nil
}

func (m *mockOrgResolver) HealthReport() map[string]error {
	return map[string]error{m.Name(): nil}
}

func (m *mockOrgResolver) Name() string {
	return "MockOrgResolver"
}

func (m *mockOrgResolver) Ready() error {
	return nil
}

func TestEngine_Execution(t *testing.T) {
	quarantine.Flaky(t, "DX-1725")
	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}
	beholderObserver := beholdertest.NewObserver(t)
	cfg.BeholderEmitter = custmsg.NewLabeler()

	t.Run("successful execution with no capability calls", func(t *testing.T) {
		engine, err := v2.NewEngine(cfg)
		require.NoError(t, err)
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
		trigger := capmocks.NewTriggerCapability(t)
		capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil)
		eventCh := make(chan capabilities.TriggerResponse)
		var capturedTriggerRequest capabilities.TriggerRegistrationRequest
		trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).
			Run(func(ctx context.Context, req capabilities.TriggerRegistrationRequest) {
				capturedTriggerRequest = req
			}).
			Return(eventCh, nil).Once()
		trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

		require.NoError(t, engine.Start(t.Context()))

		require.NoError(t, <-initDoneCh) // successful trigger registration
		require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

		require.Equal(t, fmt.Sprintf("trigger_reg_%s_%d", cfg.WorkflowID, 0), capturedTriggerRequest.TriggerID)
		require.Equal(t, cfg.WorkflowID, capturedTriggerRequest.Metadata.WorkflowID)
		require.Equal(t, cfg.WorkflowOwner, capturedTriggerRequest.Metadata.WorkflowOwner)
		require.Equal(t, cfg.WorkflowName.Hex(), capturedTriggerRequest.Metadata.WorkflowName)
		require.Equal(t, cfg.WorkflowTag, capturedTriggerRequest.Metadata.WorkflowTag)
		require.Equal(t, uint32(0), capturedTriggerRequest.Metadata.WorkflowDonID)
		require.Equal(t, uint32(0), capturedTriggerRequest.Metadata.WorkflowDonConfigVersion)
		require.Equal(t, "trigger_0", capturedTriggerRequest.Metadata.ReferenceID)
		require.Equal(t, "method", capturedTriggerRequest.Method)
		require.Nil(t, capturedTriggerRequest.Payload)
		require.Nil(t, capturedTriggerRequest.Config)

		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "event_012345",
			Payload:     nil,
		}

		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(
				func(_ context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
					wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
					require.NoError(t, err)
					capExec, ok := executor.(*v2.ExecutionHelper)
					require.True(t, ok)
					require.Equal(t, wantExecID, capExec.WorkflowExecutionID)
					require.Equal(t, uint64(0), request.Request.(*sdkpb.ExecuteRequest_Trigger).Trigger.Id)
				},
			).
			Return(nil, nil).
			Once()

		// trigger event with an error should not start an execution
		eventCh <- capabilities.TriggerResponse{
			Err: errors.New("trigger event error"),
		}

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}
		<-executionFinishedCh

		require.NoError(t, engine.Close())

		requireEventsLabels(t, beholderObserver, map[string]string{
			"workflowID":    cfg.WorkflowID,
			"workflowOwner": cfg.WorkflowOwner,
			"workflowName":  cfg.WorkflowName.String(),
		})
		requireEventsMessages(t, beholderObserver, []string{
			"Started",
			"Registering trigger",
			"All triggers registered successfully",
			"Workflow Engine initialized",
			"Workflow execution finished successfully",
		})
	})
}

func TestEngine_ExecutionTimeout(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
		// Set a very short execution timeout (100ms)
		cfg.ExecutionTimeout.DefaultValue = 100 * time.Millisecond
	})
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, status string) {
			// Verify the execution status is timeout
			require.Equal(t, "timeout", status)
			executionFinishedCh <- executionID
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	// Setup trigger registration
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

	// Mock a long-running execution that will exceed the timeout
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
			// Simulate work that takes longer than the 100ms timeout
			select {
			case <-time.After(200 * time.Millisecond):
				// This should not complete due to context timeout
			case <-ctx.Done():
				// Context should be canceled due to timeout
				require.Error(t, ctx.Err())
				require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
			}
		}).
		Return(nil, context.DeadlineExceeded).
		Once()

	// Start the engine and wait for initialization and trigger subscription
	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	// Trigger the execution
	mockTriggerEvent := capabilities.TriggerEvent{
		TriggerType: "basic-trigger@1.0.0",
		ID:          "timeout_test_event",
		Payload:     nil,
	}

	eventCh <- capabilities.TriggerResponse{
		Event: mockTriggerEvent,
	}

	// Wait for execution to finish with timeout status
	executionID := <-executionFinishedCh
	wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
	require.NoError(t, err)
	require.Equal(t, wantExecID, executionID)

	require.NoError(t, engine.Close())
}

func TestEngine_Metering_ValidBillingClient(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	// all tests in this section assume that the billing client returns valid rate cards
	billingClient := setupMockBillingClient(t)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	var logs *observer.ObservedLogs

	cfg := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
		cfg.CapabilityCallTimeout.DefaultValue = 50 * time.Millisecond
	})
	cfg.Lggr, logs = logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, status string) {
			executionFinishedCh <- executionID
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	// Setup trigger registration
	trigger := capmocks.NewTriggerCapability(t)
	eventCh := make(chan capabilities.TriggerResponse)

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	t.Run("incorrect ratios config switches to metering mode", func(t *testing.T) {
		// Setup a metered capability
		capability := capmocks.NewExecutableCapability(t)

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, "metered-capability-1").
			Return(capability, nil).Once()

		capreg.EXPECT().
			ConfigForCapability(mock.Anything, mock.Anything, mock.Anything).
			Return(capabilities.CapabilityConfiguration{}, nil).Once()

		// return some spend types in the Info call
		capability.EXPECT().
			Info(matches.AnyContext).
			Return(capabilities.CapabilityInfo{
				DON: &capabilities.DON{
					ID: 42,
				},
				SpendTypes: []capabilities.CapabilitySpendType{
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()),
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_NETWORK.String()),
				},
			}, nil).Once()

		// verify that spend limits is set and has a length of zero
		capability.EXPECT().
			Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				assert.NotNil(t, req.Metadata.SpendLimits)
				assert.Empty(t, req.Metadata.SpendLimits, 0)
			}).
			Return(capabilities.CapabilityResponse{}, nil).Once()

		// Mock workflow execution that calls the metered capability
		module.EXPECT().
			Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				// Simulate calling the slow capability from within the workflow
				_, errCap := executor.CallCapability(ctx, &sdkpb.CapabilityRequest{
					Id:         "metered-capability-1",
					Method:     "execute",
					CallbackId: 1,
					Payload:    nil,
				})

				require.NoError(t, errCap)
			}).Return(nil, nil).Once()

		// Trigger the execution
		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "metering_capability_test_1",
			Payload:     nil,
		}

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}

		// Wait for execution to finish with error status
		executionID := <-executionFinishedCh
		wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)

		require.NoError(t, err)
		require.Equal(t, wantExecID, executionID)
		capability.AssertExpectations(t)

		logged := logs.TakeAll()
		require.Len(t, logged, 1)
		assert.Contains(t, logged[0].Message, "switching to metering mode")
	})

	t.Run("correct ratios config produces spending limits", func(t *testing.T) {
		// Setup a metered capability
		capability := capmocks.NewExecutableCapability(t)

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, "metered-capability-2").
			Return(capability, nil).Once()

		ratios, _ := values.NewMap(map[string]any{
			metering.RatiosKey: map[string]string{
				billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(): "0.4",
				billing.ResourceType_RESOURCE_TYPE_NETWORK.String(): "0.6",
			},
		})

		capreg.EXPECT().
			ConfigForCapability(mock.Anything, mock.Anything, mock.Anything).
			Return(capabilities.CapabilityConfiguration{RestrictedConfig: ratios}, nil).Once()

		// return some spend types in the Info call
		capability.EXPECT().
			Info(matches.AnyContext).
			Return(capabilities.CapabilityInfo{
				DON: &capabilities.DON{
					ID: 42,
				},
				SpendTypes: []capabilities.CapabilitySpendType{
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()),
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_NETWORK.String()),
				},
			}, nil).Once()

		// verify that spend limits is set and has a length of two
		capability.EXPECT().
			Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				assert.NotNil(t, req.Metadata.SpendLimits)
				assert.Len(t, req.Metadata.SpendLimits, 2)
			}).
			Return(capabilities.CapabilityResponse{
				Metadata: capabilities.ResponseMetadata{
					Metering: []capabilities.MeteringNodeDetail{
						{
							Peer2PeerID: "local",
							SpendUnit:   billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
							SpendValue:  "100",
						},
						{
							Peer2PeerID: "local",
							SpendUnit:   billing.ResourceType_RESOURCE_TYPE_NETWORK.String(),
							SpendValue:  "1000",
						},
					},
				},
			}, nil).Once()

		// Mock workflow execution that calls the metered capability
		module.EXPECT().
			Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				// Simulate calling the slow capability from within the workflow
				_, errCap := executor.CallCapability(ctx, &sdkpb.CapabilityRequest{
					Id:         "metered-capability-2",
					Method:     "execute",
					CallbackId: 1,
					Payload:    nil,
				})

				require.NoError(t, errCap)
			}).Return(nil, nil).Once()

		// Trigger the execution
		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "metering_capability_test_2",
			Payload:     nil,
		}

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}

		// Wait for execution to finish with error status
		executionID := <-executionFinishedCh
		wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)

		require.NoError(t, err)
		require.Equal(t, wantExecID, executionID)
		capability.AssertExpectations(t)

		logged := logs.TakeAll()
		require.Empty(t, logged)
	})

	t.Run("single spend type and no ratios config produces spending limit with no error", func(t *testing.T) {
		// Setup a metered capability
		capability := capmocks.NewExecutableCapability(t)

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, "metered-capability-3").
			Return(capability, nil).Once()

		capreg.EXPECT().
			ConfigForCapability(mock.Anything, mock.Anything, mock.Anything).
			Return(capabilities.CapabilityConfiguration{}, nil).Once()

		// return some spend types in the Info call
		capability.EXPECT().
			Info(matches.AnyContext).
			Return(capabilities.CapabilityInfo{
				DON: &capabilities.DON{
					ID: 42,
				},
				SpendTypes: []capabilities.CapabilitySpendType{
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()),
				},
			}, nil).Once()

		// verify that spend limits is set and has a length of one
		capability.EXPECT().
			Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				assert.NotNil(t, req.Metadata.SpendLimits)
				assert.Len(t, req.Metadata.SpendLimits, 1)
			}).
			Return(capabilities.CapabilityResponse{
				Metadata: capabilities.ResponseMetadata{
					Metering: []capabilities.MeteringNodeDetail{
						{
							Peer2PeerID: "local",
							SpendUnit:   billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
							SpendValue:  "100",
						},
					},
				},
			}, nil).Once()

		// Mock workflow execution that calls the metered capability
		module.EXPECT().
			Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				// Simulate calling the slow capability from within the workflow
				_, errCap := executor.CallCapability(ctx, &sdkpb.CapabilityRequest{
					Id:         "metered-capability-3",
					Method:     "execute",
					CallbackId: 1,
					Payload:    nil,
				})

				require.NoError(t, errCap)
			}).Return(nil, nil).Once()

		// Trigger the execution
		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "metering_capability_test_3",
			Payload:     nil,
		}

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}

		// Wait for execution to finish with error status
		executionID := <-executionFinishedCh
		wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)

		require.NoError(t, err)
		require.Equal(t, wantExecID, executionID)
		capability.AssertExpectations(t)

		logged := logs.TakeAll()
		require.Empty(t, logged)
	})

	t.Run("billing type and capability settle spend type mismatch", func(t *testing.T) {
		// Setup a metered capability
		capability := capmocks.NewExecutableCapability(t)

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, "metered-capability-2").
			Return(capability, nil).Once()

		ratios, _ := values.NewMap(map[string]any{
			metering.RatiosKey: map[string]string{
				billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(): "0.4",
				billing.ResourceType_RESOURCE_TYPE_NETWORK.String(): "0.6",
			},
		})

		capreg.EXPECT().
			ConfigForCapability(mock.Anything, mock.Anything, mock.Anything).
			Return(capabilities.CapabilityConfiguration{RestrictedConfig: ratios}, nil).Once()

		// return some spend types in the Info call
		capability.EXPECT().
			Info(matches.AnyContext).
			Return(capabilities.CapabilityInfo{
				DON: &capabilities.DON{
					ID: 42,
				},
				SpendTypes: []capabilities.CapabilitySpendType{
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()),
					capabilities.CapabilitySpendType(billing.ResourceType_RESOURCE_TYPE_NETWORK.String()),
				},
			}, nil).Once()

		// verify that spend limits is set and has a length of two
		capability.EXPECT().
			Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				assert.NotNil(t, req.Metadata.SpendLimits)
				assert.Len(t, req.Metadata.SpendLimits, 2)
			}).
			Return(capabilities.CapabilityResponse{
				Metadata: capabilities.ResponseMetadata{
					Metering: []capabilities.MeteringNodeDetail{
						{
							Peer2PeerID: "local",
							// SpendUnit does not match units from billing or ratios
							SpendUnit:  "COMPUTE",
							SpendValue: "100",
						},
						{
							Peer2PeerID: "local",
							SpendUnit:   billing.ResourceType_RESOURCE_TYPE_NETWORK.String(),
							SpendValue:  "1000",
						},
					},
				},
			}, nil).Once()

		// Mock workflow execution that calls the metered capability
		module.EXPECT().
			Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				// Simulate calling the slow capability from within the workflow
				_, errCap := executor.CallCapability(ctx, &sdkpb.CapabilityRequest{
					Id:         "metered-capability-2",
					Method:     "execute",
					CallbackId: 1,
					Payload:    nil,
				})

				require.NoError(t, errCap)
			}).Return(nil, nil).Once()

		// Trigger the execution
		mockTriggerEvent := capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "metering_capability_test_2",
			Payload:     nil,
		}

		eventCh <- capabilities.TriggerResponse{
			Event: mockTriggerEvent,
		}

		// Wait for execution to finish with error status
		executionID := <-executionFinishedCh
		wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)

		require.NoError(t, err)
		require.Equal(t, wantExecID, executionID)
		capability.AssertExpectations(t)

		logged := logs.TakeAll()
		require.Len(t, logged, 1)
		assert.Contains(t, logged[0].Message, "metering mode")
	})

	require.NoError(t, engine.Close())
}

func TestEngine_CapabilityCallTimeout(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
		// Set a very short capability call timeout (50ms)
		cfg.CapabilityCallTimeout.DefaultValue = 50 * time.Millisecond
	})
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, status string) {
			// Verify the execution status is errored due to capability timeout
			require.Equal(t, "errored", status)
			executionFinishedCh <- executionID
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	// Setup trigger registration
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()

	// Setup a slow capability that will timeout
	slowCapability := capmocks.NewExecutableCapability(t)
	capreg.EXPECT().GetExecutable(matches.AnyContext, "slow-capability").Return(slowCapability, nil).Once()
	capreg.EXPECT().
		ConfigForCapability(mock.Anything, mock.Anything, mock.Anything).
		Return(capabilities.CapabilityConfiguration{}, nil).
		Once()

	slowCapability.EXPECT().
		Info(matches.AnyContext).
		Return(capabilities.CapabilityInfo{
			DON: &capabilities.DON{
				ID: 42,
			},
		}, nil)
	// Mock capability that takes longer than the 50ms timeout
	slowCapability.EXPECT().Execute(matches.AnyContext, mock.Anything).
		Run(func(ctx context.Context, req capabilities.CapabilityRequest) {
			assert.Equal(t, testWorkflowNameA, req.Metadata.DecodedWorkflowName)
			assert.Equal(t, hashedTestWorkflowNameA, req.Metadata.WorkflowName)
			// Simulate work that takes longer than the 50ms timeout
			select {
			case <-time.After(100 * time.Millisecond):
				// This should not complete due to context timeout
			case <-ctx.Done():
				// Context should be canceled due to timeout
				require.Error(t, ctx.Err())
				require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
			}
		}).
		Return(capabilities.CapabilityResponse{}, context.DeadlineExceeded).
		Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	// Mock workflow execution that calls the slow capability
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, request *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
			// Simulate calling the slow capability from within the workflow
			_, errCap := executor.CallCapability(ctx, &sdkpb.CapabilityRequest{
				Id:         "slow-capability",
				Method:     "execute",
				CallbackId: 1,
				Payload:    nil,
			})
			// Verify that the capability call returns a timeout error
			require.Error(t, errCap)
			require.Contains(t, errCap.Error(), "failed to execute capability")
		}).
		Return(nil, errors.New("capability timeout error")).
		Once()

	// Trigger the execution
	mockTriggerEvent := capabilities.TriggerEvent{
		TriggerType: "basic-trigger@1.0.0",
		ID:          "timeout_capability_test",
		Payload:     nil,
	}

	eventCh <- capabilities.TriggerResponse{
		Event: mockTriggerEvent,
	}

	// Wait for execution to finish with error status
	executionID := <-executionFinishedCh
	wantExecID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, mockTriggerEvent.ID)
	require.NoError(t, err)
	require.Equal(t, wantExecID, executionID)

	require.NoError(t, engine.Close())
}

func TestEngine_WASMBinary_Simple(t *testing.T) {
	cmd := "core/services/workflows/test/wasm/v2/cmd"
	log := logger.TestLogger(t)
	binaryB := wasmtest.CreateTestBinary(cmd, false, t)
	module, err := host.NewModule(&host.ModuleConfig{
		Logger:         log,
		IsUncompressed: true,
	}, binaryB)
	require.NoError(t, err)

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	billingClient := setupMockBillingClient(t)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	resultReceivedCh := make(chan *sdkpb.ExecutionResult, 1)
	executionFinishedCh := make(chan string, 1)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
		OnResultReceived: func(er *sdkpb.ExecutionResult) {
			resultReceivedCh <- er
		},
	}

	basicActionMock := setupExpectedCalls(t)
	wrappedTriggerMock := &TriggerCapabilityWrapper{}
	wrappedActionMock := &MockCapabilityWrapper{
		Capability: basicActionMock,
	}

	t.Run("OK happy path", func(t *testing.T) {
		wantResponse := "Hello, world!"
		engine, err := v2.NewEngine(cfg)
		require.NoError(t, err)

		capreg.EXPECT().
			GetTrigger(matches.AnyContext, triggerID).
			Return(wrappedTriggerMock, nil).
			Once()

		capreg.EXPECT().
			GetExecutable(matches.AnyContext, wrappedActionMock.ID()).
			Return(wrappedActionMock, nil).
			Twice()

		testConf, _ := values.NewMap(map[string]any{
			"spendRatios": map[string]string{
				"spendTypeA": "0.4",
				"spendTypeB": "0.6",
			},
		})

		capreg.EXPECT().
			ConfigForCapability(matches.AnyContext, mock.Anything, mock.Anything).
			Return(capabilities.CapabilityConfiguration{
				RestrictedConfig: testConf,
			}, nil)

		require.NoError(t, engine.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
		require.Equal(t, []string{triggerID}, <-subscribedToTriggersCh)

		// Read the result from the hook and assert that the wanted response was
		// received.
		res := <-resultReceivedCh
		switch output := res.Result.(type) {
		case *sdkpb.ExecutionResult_Value:
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

		execID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, "")
		require.NoError(t, err)

		require.Equal(t, execID, <-executionFinishedCh)
		require.NoError(t, engine.Close())
	})
}

func TestEngine_WASMBinary_With_Config(t *testing.T) {
	cmd := "core/services/workflows/test/wasm/v2/cmd/with_config"
	binaryB := wasmtest.CreateTestBinary(cmd, false, t)

	// Define a custom config to validate against
	giveName := "Foo"
	giveNum := int32(42)
	config := fmt.Appendf(nil, "name: %s\nnumber: %d\n", giveName, giveNum)
	wasmLogger := logger.NewMockLogger(t)
	module, err := host.NewModule(&host.ModuleConfig{
		Logger:         wasmLogger,
		IsUncompressed: true,
	}, binaryB)
	require.NoError(t, err)

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	billingClient := setupMockBillingClient(t)

	cfg := defaultTestConfig(t, nil)
	cfg.WorkflowConfig = config
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	resultReceivedCh := make(chan *sdkpb.ExecutionResult, 1)
	executionFinishedCh := make(chan string, 1)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
		OnResultReceived: func(er *sdkpb.ExecutionResult) {
			resultReceivedCh <- er
		},
	}

	wrappedTriggerMock := &TriggerCapabilityWrapper{
		giveName:   giveName,
		giveNumber: giveNum,
	}
	beholderObserver := beholdertest.NewObserver(t)

	t.Run("OK received expected config", func(t *testing.T) {
		engine, err := v2.NewEngine(cfg)
		require.NoError(t, err)

		capreg.EXPECT().
			GetTrigger(matches.AnyContext, triggerID).
			Return(wrappedTriggerMock, nil).
			Once()

		require.NoError(t, engine.Start(t.Context()))
		require.NoError(t, <-initDoneCh)
		require.Equal(t, []string{triggerID}, <-subscribedToTriggersCh)

		// Read the result from the hook and assert that the wanted response was
		// received.
		res := <-resultReceivedCh
		switch output := res.Result.(type) {
		case *sdkpb.ExecutionResult_Value:
			var value values.Value
			var execErr error
			var unwrapped any

			valuePb := output.Value
			value, execErr = values.FromProto(valuePb)
			require.NoError(t, execErr)
			unwrapped, execErr = value.Unwrap()
			require.NoError(t, execErr)
			require.Equal(t, string(config), unwrapped)
		default:
			t.Fatalf("unexpected response type %T", output)
		}

		execID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, "")
		require.NoError(t, err)

		require.Equal(t, execID, <-executionFinishedCh)
		require.NoError(t, engine.Close())

		requireUserLogs(t, beholderObserver, []string{
			"onTrigger called",
		})
	})
}

func TestSecretsFetcher_Integration(t *testing.T) {
	cmd := "core/services/workflows/test/wasm/v2/cmd/with_secrets"
	binaryB := wasmtest.CreateTestBinary(cmd, false, t)

	// Define a custom config to validate against
	giveName := "Foo"
	giveNum := int32(42)
	config := fmt.Appendf(nil, "name: %s\nnumber: %d\n", giveName, giveNum)
	wasmLogger := logger.NewMockLogger(t)
	module, err := host.NewModule(&host.ModuleConfig{
		Logger:         wasmLogger,
		IsUncompressed: true,
	}, binaryB)
	require.NoError(t, err)

	capreg := regmocks.NewCapabilitiesRegistry(t)
	peer := coreCap.RandomUTF8BytesWord()
	localRegistry := v2.CreateLocalRegistry(t, peer)
	localNode, err := localRegistry.LocalNode(t.Context())
	require.NoError(t, err)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(localNode, nil)
	for _, peerID := range localNode.WorkflowDON.Members {
		node, err2 := localRegistry.NodeByPeerID(t.Context(), peerID)
		require.NoError(t, err2)
		capreg.EXPECT().NodeByPeerID(matches.AnyContext, peerID).Return(node, nil)
	}

	billingClient := setupMockBillingClient(t)
	cfg := defaultTestConfig(t, nil)
	cfg.WorkflowConfig = config
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient

	rawSecret := "Original Secret Text"
	f, n := 2, 3
	_, vaultPublicKey, privateShares, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)

	cipher, err := tdh2easy.Encrypt(vaultPublicKey, []byte(rawSecret))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)

	decryptionShare0, err := tdh2easy.Decrypt(cipher, privateShares[0])
	require.NoError(t, err)
	decryptionShare0Bytes, err := decryptionShare0.Marshal()
	require.NoError(t, err)
	decryptionShare1, err := tdh2easy.Decrypt(cipher, privateShares[1])
	require.NoError(t, err)
	decryptionShare1Bytes, err := decryptionShare1.Marshal()
	require.NoError(t, err)
	decryptionShare2, err := tdh2easy.Decrypt(cipher, privateShares[2])
	require.NoError(t, err)
	decryptionShare2Bytes, err := decryptionShare2.Marshal()
	require.NoError(t, err)

	// Sanity testing that we can decrypt the secret with just 2 shares
	twoDecryptionShares := []*tdh2easy.DecryptionShare{decryptionShare0, decryptionShare1}
	decryptedSecret, err := tdh2easy.Aggregate(cipher, twoDecryptionShares, n)
	require.NoError(t, err)
	assert.Equal(t, rawSecret, string(decryptedSecret))

	// Encrypt the decryption shares with the workflow key. This is the expected output from Vault capability.
	encryptedDecryptionShare0, err := cfg.WorkflowEncryptionKey.Encrypt(decryptionShare0Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare1, err := cfg.WorkflowEncryptionKey.Encrypt(decryptionShare1Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare2, err := cfg.WorkflowEncryptionKey.Encrypt(decryptionShare2Bytes)
	require.NoError(t, err)
	workflowKeyBytes := cfg.WorkflowEncryptionKey.PublicKey()

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Default",
							Owner:     common.HexToAddress("0x" + testWorkflowOwnerA).Hex(),
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedDecryptionShare0),
											hex.EncodeToString(encryptedDecryptionShare2),
											hex.EncodeToString([]byte("blabbermouth")),
											hex.EncodeToString(encryptedDecryptionShare1),
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	capreg.EXPECT().GetExecutable(matches.AnyContext, vault.CapabilityID).Return(mc, nil)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	valueMap, err := values.WrapMap(v2.VaultCapabilityRegistryConfig{
		VaultPublicKey: hex.EncodeToString(vaultPublicKeyBytes),
		Threshold:      1,
	})
	require.NoError(t, err)
	capConfig := capabilities.CapabilityConfiguration{
		DefaultConfig: valueMap,
	}
	capreg.EXPECT().ConfigForCapability(matches.AnyContext, vault.CapabilityID, localNode.WorkflowDON.ID).Return(capConfig, nil)

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	resultReceivedCh := make(chan *sdkpb.ExecutionResult, 1)
	executionFinishedCh := make(chan string, 1)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
		OnResultReceived: func(er *sdkpb.ExecutionResult) {
			resultReceivedCh <- er
		},
	}

	wrappedTriggerMock := &TriggerCapabilityWrapper{
		giveName:   giveName,
		giveNumber: giveNum,
	}

	secretsFetcher := v2.NewSecretsFetcher(
		v2.MetricsLabelerTest(t),
		cfg.CapRegistry,
		cfg.Lggr,
		cfg.LocalLimiters.SecretsConcurrency,
		cfg.WorkflowOwner,
		cfg.WorkflowName.String(),
		cfg.WorkflowID,
		"",
		cfg.WorkflowEncryptionKey,
	)
	cfg.SecretsFetcher = secretsFetcher
	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	capreg.EXPECT().
		GetTrigger(matches.AnyContext, triggerID).
		Return(wrappedTriggerMock, nil).
		Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{triggerID}, <-subscribedToTriggersCh)

	// Read the result from the hook and assert that the wanted response was
	// received.
	res := <-resultReceivedCh
	switch output := res.Result.(type) {
	case *sdkpb.ExecutionResult_Value:
		var value values.Value
		var execErr error
		var unwrapped any

		valuePb := output.Value
		value, execErr = values.FromProto(valuePb)
		require.NoError(t, execErr)
		unwrapped, execErr = value.Unwrap()
		require.NoError(t, execErr)
		require.Equal(t, rawSecret, unwrapped)
	default:
		t.Fatalf("unexpected response type %T: %v", output, output)
	}

	execID, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, "")
	require.NoError(t, err)

	require.Equal(t, execID, <-executionFinishedCh)
	require.NoError(t, engine.Close())
}

// setupMockBillingClient creates a mock billing client with default expectations.
func setupMockBillingClient(t *testing.T) *metmocks.BillingClient {
	billingClient := metmocks.NewBillingClient(t)

	billingClient.EXPECT().
		GetWorkflowExecutionRates(mock.Anything, mock.Anything).
		Return(&billing.GetWorkflowExecutionRatesResponse{
			RateCards: []*billing.RateCard{
				{
					ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
					MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
					UnitsPerCredit:  "0.0001",
				},
				{
					ResourceType:    billing.ResourceType_RESOURCE_TYPE_NETWORK,
					MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_COST,
					UnitsPerCredit:  "0.01",
				},
			},
		}, nil)
	billingClient.EXPECT().
		ReserveCredits(mock.Anything, mock.MatchedBy(func(req *billing.ReserveCreditsRequest) bool {
			return req != nil && req.WorkflowId != "" && req.WorkflowExecutionId != ""
		})).
		Return(&billing.ReserveCreditsResponse{
			Success: true,
			Credits: "10000",
		}, nil)
	billingClient.EXPECT().
		SubmitWorkflowReceipt(mock.Anything, mock.MatchedBy(func(req *billing.SubmitWorkflowReceiptRequest) bool {
			return req != nil && req.WorkflowId != "" && req.WorkflowExecutionId != ""
		})).
		Return(&emptypb.Empty{}, nil).Maybe()
	return billingClient
}

// setupExpectedCalls mocks single call to trigger and two calls to the basic action
// mock capability
func setupExpectedCalls(t *testing.T) *basicactionmock.BasicActionCapability {
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
	return basicAction
}

func requireEventsLabels(t *testing.T, beholderObserver beholdertest.Observer, want map[string]string) {
	msgs := beholderObserver.Messages(t)
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

func requireEventsMessages(t *testing.T, beholderObserver beholdertest.Observer, expected []string) {
	msgs := beholderObserver.Messages(t)
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

func requireUserLogs(t *testing.T, beholderObserver beholdertest.Observer, expectedSubstrings []string) {
	msgs := beholderObserver.Messages(t)
	nextToFind := 0
	for _, msg := range msgs {
		if msg.Attrs["beholder_entity"] == "workflows.v1.UserLogs" {
			var payload events.UserLogs
			require.NoError(t, proto.Unmarshal(msg.Body, &payload))
			if nextToFind >= len(expectedSubstrings) {
				return
			}
			for _, log := range payload.LogLines {
				if strings.Contains(log.Message, expectedSubstrings[nextToFind]) {
					nextToFind++
				}
			}
		}
	}

	if nextToFind < len(expectedSubstrings) {
		t.Errorf("log message not found: %s", expectedSubstrings[nextToFind])
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

type MockCapabilityWrapper struct {
	registry.Capability
}

var _ capabilities.ExecutableCapability = (*MockCapabilityWrapper)(nil)

func (c *MockCapabilityWrapper) RegisterToWorkflow(_ context.Context, _ capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (c *MockCapabilityWrapper) UnregisterFromWorkflow(_ context.Context, _ capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func (c *MockCapabilityWrapper) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	v1Request := capabilitiespb.CapabilityRequestToProto(request)
	v2Request := &sdkpb.CapabilityRequest{
		Id:      v1Request.Metadata.ReferenceId,
		Payload: v1Request.Payload,
		Method:  v1Request.Method,
	}

	v2Response := c.Invoke(ctx, v2Request)
	switch r := v2Response.Response.(type) {
	case *sdkpb.CapabilityResponse_Error:
		return capabilities.CapabilityResponse{}, errors.New(r.Error)
	case *sdkpb.CapabilityResponse_Payload:
		return capabilities.CapabilityResponse{
			Payload: r.Payload,
		}, nil
	default:
		return capabilities.CapabilityResponse{}, fmt.Errorf("unknown capability response type: %T", r)
	}
}

func (c *MockCapabilityWrapper) Info(_ context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		c.ID(), capabilities.CapabilityTypeCombined, "Mock of capability %s"+c.ID())
}

type TriggerCapabilityWrapper struct {
	giveName   string
	giveNumber int32
}

var _ capabilities.TriggerCapability = &TriggerCapabilityWrapper{}

func (c *TriggerCapabilityWrapper) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	ch := make(chan capabilities.TriggerResponse, 1)
	defer close(ch)

	config := &basictrigger.Config{}
	if err := request.Payload.UnmarshalTo(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trigger config: %w", err)
	}

	if c.giveName != "" {
		if config.Name != c.giveName {
			return nil, fmt.Errorf("expected trigger name %s, got %s", c.giveName, config.Name)
		}

		if config.Number != c.giveNumber {
			return nil, fmt.Errorf("expected trigger number %d, got %d", c.giveNumber, config.Number)
		}
	}

	trigger := &basictrigger.Outputs{CoolOutput: "Hello, "}
	payload, err := anypb.New(trigger)
	if err != nil {
		return nil, err
	}
	ch <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: request.TriggerID,
			Payload:     payload,
		},
	}

	return ch, nil
}

func (c *TriggerCapabilityWrapper) UnregisterTrigger(_ context.Context, _ capabilities.TriggerRegistrationRequest) error {
	return nil
}

func (c *TriggerCapabilityWrapper) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		triggerID,
		capabilities.CapabilityTypeTrigger,
		"Mock of trigger capability for testing",
	)
}
