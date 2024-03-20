package workflows

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// NOTE: max 32 bytes per ID - consider enforcing exactly 32 bytes?
	mockedWorkflowID  = "aaaaaaaaaa0000000000000000000000"
	mockedExecutionID = "bbbbbbbbbb0000000000000000000000"
	mockedTriggerID   = "cccccccccc0000000000000000000000"
)

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	logger       logger.Logger
	registry     types.CapabilitiesRegistry
	workflow     *workflow
	executionStates        *store
	// NOTE: I do find it confusing that pending step requests are global rather than scoped to a single execution
	pendingStepRequests        *queue[stepRequest]
	triggerEvents   chan capabilities.CapabilityResponse
	newWorkerCh  chan struct{}
	stepUpdateCh chan stepState
	// wg is only used to make sure that in the case of a shutdown,
	// we wait for all pending steps to finish.
	wg           sync.WaitGroup
	stopCh       services.StopChan
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, _ := e.stopCh.NewCtx()

		// queue.start will add to the wg and
		// spin off a goroutine.
		e.pendingStepRequests.start(ctx, &e.wg)

		e.wg.Add(2)
		go e.init(ctx)
		go e.loop(ctx)

		return nil
	})
}

// init does the following:
//
//  1. Resolves the underlying capability for each trigger
//  2. Registers each step's capability to this workflow
//  3. Registers for trigger events now that all capabilities are resolved
//
// Steps 1 and 2 are retried every 5 seconds until successful. 
func (e *Engine) init(ctx context.Context) {
	defer e.wg.Done()

	retrySec := 5
	ticker := time.NewTicker(time.Duration(retrySec) * time.Second)
	defer ticker.Stop()

	initSuccessful := true
LOOP:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Resolve the underlying capability for each trigger
			for _, t := range e.workflow.triggers {
				cp, err := e.registry.GetTrigger(ctx, t.Type)
				if err != nil {
					initSuccessful = false
					e.logger.Errorf("failed to get trigger capability: %s, retrying in %d seconds", err, retrySec)
				} else {
					t.cachedTrigger = cp
				}
			}

			// Walk the graph and register each step's capablity to this workflow
			err := e.workflow.walkDo(keywordTrigger, func(n *step) error {
				// The graph contains a dummy step for triggers, but
				// we handle triggers separately since there might be more than one.
				if n.Ref == keywordTrigger {
					return nil
				}

				// If the capability is already cached, that means we've already registered it
				if n.cachedCapability != nil {
					return nil
				}

				cp, innerErr := e.registry.Get(ctx, n.Type)
				if innerErr != nil {
					return fmt.Errorf("failed to get capability with ref %s: %s, retrying in %d seconds", n.Type, innerErr, retrySec)
				}

				// We only support CallbackExecutable capabilities for now 
				cc, ok := cp.(capabilities.CallbackExecutable)
				if !ok {
					return fmt.Errorf("could not coerce capability %s to CallbackExecutable", n.Type)
				}

				if n.cachedConfig == nil {
					configMap, ierr := values.NewMap(n.Config)
					if innerErr != nil {
						return fmt.Errorf("failed to convert config to values.Map: %s", ierr)
					}
					n.cachedConfig = configMap
				}

				reg := capabilities.RegisterToWorkflowRequest{
					Metadata: capabilities.RegistrationMetadata{
						WorkflowID: mockedWorkflowID,
					},
					Config: n.cachedConfig,
				}

				innerErr = cc.RegisterToWorkflow(ctx, reg)
				if innerErr != nil {
					return fmt.Errorf("failed to register to workflow: %+v", reg)
				}

				n.cachedCapability = cc
				return nil
			})
			if err != nil {
				initSuccessful = false
				e.logger.Error(err)
			}

			if initSuccessful {
				break LOOP
			}
		}
	}

	// We have all needed capabilities, now we can register for trigger events
	for _, t := range e.workflow.triggers {
		err := e.registerTrigger(ctx, t)
		if err != nil {
			e.logger.Errorf("failed to register trigger: %s", err)
		}
	}

	e.logger.Info("engine initialized")
}

// registerTrigger is used during the initialization phase to bind a trigger to this workflow 
func (e *Engine) registerTrigger(ctx context.Context, t *triggerCapability) error {
	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": mockedTriggerID,
		},
	)
	if err != nil {
		return err
	}

	tc, err := values.NewMap(t.Config)
	if err != nil {
		return err
	}

	triggerRegRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Config: tc,
		Inputs: triggerInputs,
	}
	err = t.cachedTrigger.RegisterTrigger(ctx, e.triggerEvents, triggerRegRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate trigger %s, %s", t.Type, err)
	}
	return nil
}

// loop is the synchronization goroutine for the engine, and is responsible for:
//  - dispatching new workers up to the limit specified (default = 100)
//  - starting a new execution when a trigger emits a message on `triggerEvents`
//  - updating the `executionState` with the outcome of a `step`.
//
// Note: `executionState` is only mutated by this loop directly.
//
// This is important to avoid data races, and any accesses of `executionState` by any other
// goroutine should happen via a `stepRequest` message containing a copy of the latest
// `executionState`.
//
// This works because a worker thread for a given step will only
// be spun up once all dependent steps have completed (guaranteeing that the state associated
// with those dependent steps will no longer change). Therefore as long this worker thread only
// accesses data from dependent states, the data will never be stale.
func (e *Engine) loop(ctx context.Context) {
	defer e.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-e.triggerEvents:
			if resp.Err != nil {
				e.logger.Errorf("trigger event was an error; not executing", resp.Err)
				continue
			}

			err := e.startExecution(ctx, resp.Value)
			if err != nil {
				e.logger.Errorf("failed to start execution: %w", err)
			}
		case stepRequest := <-e.pendingStepRequests.dequeue:
			// Wait for a new worker to be available before dispatching a new one.
			<-e.newWorkerCh
			// NOTE: Can we add this to e.workerForStep instead?
			e.wg.Add(1)
			// NOTE: Should we instead add a "process" method to the queue, and do concurrency control there? 
			go e.workerForStepRequest(ctx, stepRequest)
		case stepUpdate := <-e.stepUpdateCh:
			// Executed synchronously to ensure we correctly schedule subsequent tasks.
			err := e.handleStepUpdate(ctx, stepUpdate)
			if err != nil {
				e.logger.Errorf("failed to update step state: %+v, %s", stepUpdate, err)
			}
		}
	}
}

// startExecution kicks off a new workflow execution when a trigger event is received.
func (e *Engine) startExecution(ctx context.Context, event values.Value) error {
	executionID := uuid.New().String()
	e.logger.Debugw("executing on a trigger event", "event", event, "executionID", executionID)
	ec := &executionState{
		steps: map[string]*stepState{
			keywordTrigger: {
				outputs: &stepOutput{
					value: event,
				},
				status: statusCompleted,
			},
		},
		workflowID:  mockedWorkflowID,
		executionID: executionID,
		status:      statusStarted,
	}

	err := e.executionStates.add(ctx, ec)
	if err != nil {
		return err
	}

	// Find the tasks we need to fire when a trigger has fired and enqueue them.
	triggerDependents, err := e.workflow.dependents(keywordTrigger)
	if err != nil {
		return err
	}

	for _, step := range triggerDependents {
		e.logger.Debugw("step request enqueued", "ref", step.Ref, "executionID", executionID)
		e.pendingStepRequests.enqueue <- stepRequest{state: copyState(*ec), stepRef: step.Ref}
	}

	return nil
}

func (e *Engine) handleStepUpdate(ctx context.Context, stepUpdate stepState) error {
	state, err := e.executionStates.updateStep(ctx, &stepUpdate)
	if err != nil {
		return err
	}

	switch stepUpdate.status {
	case statusCompleted:
		stepDependents, err := e.workflow.dependents(stepUpdate.ref)
		if err != nil {
			return err
		}
		// There are no steps left to process in the current path, so let's check if
		// we've completed the workflow.
		// If not, we'll check for any dependents that are ready to process.
		if len(stepDependents) == 0 {
			workflowCompleted := true
			err := e.workflow.walkDo(keywordTrigger, func(n *step) error {
				step, ok := state.steps[n.Ref]
				// Note: Why do we not return an error if !ok?
				if !ok {
					workflowCompleted = false
					return nil
				}

				switch step.status {
				case statusCompleted, statusErrored:
				default:
					workflowCompleted = false
				}
				return nil
			})
			if err != nil {
				return err
			}

			if workflowCompleted {
				err := e.executionStates.updateStatus(ctx, state.executionID, statusCompleted)
				if err != nil {
					return err
				}
			}
		}

		for _, step := range stepDependents {
			// Check if all dependencies are completed for the current step
			var waitingOnDependencies bool
			for _, dr := range step.dependencies {
				stepState, ok := state.steps[dr]
				if !ok {
					return fmt.Errorf("could not locate dependency %s in %+v", dr, state)
				}

				// NOTE: Should we also check for statusErrored?
				if stepState.status != statusCompleted {
					waitingOnDependencies = true
				}
			}

			// If all dependencies are completed, enqueue the step.
			if !waitingOnDependencies {
				e.pendingStepRequests.enqueue <- stepRequest{
					state:   copyState(state),
					stepRef: step.Ref,
				}
			}
		}
	case statusErrored:
		err := e.executionStates.updateStatus(ctx, state.executionID, statusErrored)
		if err != nil {
			return err
		}
	}

	return nil
}

// NOTE: Should this be attached to a step struct instead of the engine?
func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	defer e.wg.Done()

	e.logger.Debugw("executing on a step event", "event", msg, "executionID", msg.state.executionID)
	stepState := &stepState{
		outputs:     &stepOutput{},
		executionID: msg.state.executionID,
		ref:         msg.stepRef,
	}

	inputs, outputs, err := e.executeStep(ctx, msg)
	if err != nil {
		e.logger.Errorf("error executing step request: %w", err, "executionID", msg.state.executionID, "stepRef", msg.stepRef)
		stepState.outputs.err = err
		stepState.status = statusErrored
	} else {
		stepState.outputs.value = outputs
		stepState.status = statusCompleted
		e.logger.Debugw("step executed successfully", "executionID", msg.state.executionID, "stepRef", msg.stepRef, "outputs", outputs)
	}

	stepState.inputs = inputs

	e.stepUpdateCh <- *stepState
	e.newWorkerCh <- struct{}{}
}

// executeStep executes the referenced capability within a step and returns the result. 
func (e *Engine) executeStep(ctx context.Context, msg stepRequest) (*values.Map, values.Value, error) {
	step, err := e.workflow.Vertex(msg.stepRef)
	if err != nil {
		return nil, nil, err
	}

	i, err := findAndInterpolateAllKeys(step.Inputs, msg.state)
	if err != nil {
		return nil, nil, err
	}

	inputs, err := values.NewMap(i.(map[string]any))
	if err != nil {
		return nil, nil, err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: inputs,
		Config: step.cachedConfig,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          msg.state.workflowID,
			WorkflowExecutionID: msg.state.executionID,
		},
	}

	resp, err := capabilities.ExecuteSync(ctx, step.cachedCapability, tr)
	if err != nil {
		return inputs, nil, err
	}

	// `ExecuteSync` returns a `values.List` even if there was
	// just one return value. If that is the case, let's unwrap the
	// single value to make it easier to use in -- for example -- variable interpolation.
	if len(resp.Underlying) > 1 {
		return inputs, resp, err
	}
	return inputs, resp.Underlying[0], err
}

func (e *Engine) deregisterTrigger(_ context.Context, t *triggerCapability) error {
	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": mockedTriggerID,
		},
	)
	if err != nil {
		return err
	}
	deregRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Inputs: triggerInputs,
	}
	return t.cachedTrigger.UnregisterTrigger(context.Background(), deregRequest)
}

func (e *Engine) Close() error {
	return e.StopOnce("Engine", func() error {
		ctx := context.Background()
		// To shut down the engine, we'll start by deregistering
		// any triggers to ensure no new executions are triggered,
		// then we'll close down any background goroutines,
		// and finally, we'll deregister any workflow steps.
		for _, t := range e.workflow.triggers {
			err := e.deregisterTrigger(ctx, t)
			if err != nil {
				return err
			}
		}

		close(e.stopCh)
		e.wg.Wait()

		err := e.workflow.walkDo(keywordTrigger, func(n *step) error {
			if n.Ref == keywordTrigger {
				return nil
			}

			reg := capabilities.UnregisterFromWorkflowRequest{
				Metadata: capabilities.RegistrationMetadata{
					WorkflowID: mockedWorkflowID,
				},
				Config: n.cachedConfig,
			}

			innerErr := n.cachedCapability.UnregisterFromWorkflow(ctx, reg)
			if innerErr != nil {
				return fmt.Errorf("failed to unregister from workflow: %+v", reg)
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
}

type Config struct {
	Spec           string
	Lggr           logger.Logger
	Registry       types.CapabilitiesRegistry
	MaxWorkerLimit int
}

const (
	defaultWorkerLimit = 100
)

func NewEngine(cfg Config) (engine *Engine, err error) {
	if cfg.MaxWorkerLimit == 0 {
		cfg.MaxWorkerLimit = defaultWorkerLimit
	}
	// TODO: validation of the workflow spec
	// We'll need to check, among other things:
	// - that there are no step `ref` called `trigger` as this is reserved for any triggers
	// - that there are no duplicate `ref`s
	// - that the `ref` for any triggers is empty -- and filled in with `trigger`
	// - etc.

	workflow, err := Parse(cfg.Spec)
	if err != nil {
		return nil, err
	}

	// Instantiate semaphore to put a limit on the number of workers
	newWorkerCh := make(chan struct{}, cfg.MaxWorkerLimit)
	for i := 0; i < cfg.MaxWorkerLimit; i++ {
		newWorkerCh <- struct{}{}
	}

	engine = &Engine{
		logger:       cfg.Lggr.Named("WorkflowEngine"),
		registry:     cfg.Registry,
		workflow:     workflow,
		executionStates:        newStore(),
		pendingStepRequests:        newQueue[stepRequest](),
		newWorkerCh:  newWorkerCh,
		stepUpdateCh: make(chan stepState),
		triggerEvents:   make(chan capabilities.CapabilityResponse),
		stopCh:       make(chan struct{}),
	}
	return engine, nil
}
