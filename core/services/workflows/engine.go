package workflows

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// NOTE: max 32 bytes per ID - consider enforcing exactly 32 bytes?
	mockedTriggerID  = "cccccccccc0000000000000000000000"
	mockedWorkflowID = "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
)

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	logger              logger.Logger
	registry            core.CapabilitiesRegistry
	workflow            *workflow
	executionStates     *inMemoryStore
	pendingStepRequests chan stepRequest
	triggerEvents       chan capabilities.CapabilityResponse
	newWorkerCh         chan struct{}
	stepUpdateCh        chan stepState
	wg                  sync.WaitGroup
	stopCh              services.StopChan
	newWorkerTimeout    time.Duration

	// Used for testing to wait for an execution to complete
	xxxExecutionFinished chan string
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, _ := e.stopCh.NewCtx()

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

LOOP:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			initSuccessful := true
			// Resolve the underlying capability for each trigger
			for _, t := range e.workflow.triggers {
				tg, err := e.registry.GetTrigger(ctx, t.Type)
				if err != nil {
					initSuccessful = false
					e.logger.Errorf("failed to get trigger capability: %s, retrying in %d seconds", err, retrySec)
					continue
				}
				t.trigger = tg
			}
			if !initSuccessful {
				continue
			}

			// Walk the graph and register each step's capability to this workflow
			err := e.workflow.walkDo(keywordTrigger, func(s *step) error {
				// The graph contains a dummy step for triggers, but
				// we handle triggers separately since there might be more than one.
				if s.Ref == keywordTrigger {
					return nil
				}

				// If the capability already exists, that means we've already registered it
				if s.capability != nil {
					return nil
				}

				cp, innerErr := e.registry.Get(ctx, s.Type)
				if innerErr != nil {
					return fmt.Errorf("failed to get capability with ref %s: %s, retrying in %d seconds", s.Type, innerErr, retrySec)
				}

				// We only need to configure actions, consensus and targets here, and
				// they all satisfy the `CallbackExecutable` interface
				cc, ok := cp.(capabilities.CallbackExecutable)
				if !ok {
					return fmt.Errorf("could not coerce capability %s to CallbackExecutable", s.Type)
				}

				if s.config == nil {
					configMap, ierr := values.NewMap(s.Config)
					if ierr != nil {
						return fmt.Errorf("failed to convert config to values.Map: %s", ierr)
					}
					s.config = configMap
				}

				reg := capabilities.RegisterToWorkflowRequest{
					Metadata: capabilities.RegistrationMetadata{
						WorkflowID: e.workflow.id,
					},
					Config: s.config,
				}

				innerErr = cc.RegisterToWorkflow(ctx, reg)
				if innerErr != nil {
					return fmt.Errorf("failed to register to workflow (%+v): %w", reg, innerErr)
				}

				s.capability = cc
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

	t.config = tc

	triggerRegRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: e.workflow.id,
		},
		Config: tc,
		Inputs: triggerInputs,
	}
	eventsCh, err := t.trigger.RegisterTrigger(ctx, triggerRegRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate trigger %s, %s", t.Type, err)
	}

	go func() {
		for event := range eventsCh {
			e.triggerEvents <- event
		}
	}()

	return nil
}

// loop is the synchronization goroutine for the engine, and is responsible for:
//   - dispatching new workers up to the limit specified (default = 100)
//   - starting a new execution when a trigger emits a message on `triggerEvents`
//   - updating the `executionState` with the outcome of a `step`.
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
			e.logger.Debugw("shutting down loop")
			return
		case resp, isOpen := <-e.triggerEvents:
			if !isOpen {
				e.logger.Errorf("trigger events channel is no longer open, skipping")
				continue
			}

			if resp.Err != nil {
				e.logger.Errorf("trigger event was an error; not executing", resp.Err)
				continue
			}

			te := &capabilities.TriggerEvent{}
			err := resp.Value.UnwrapTo(te)
			if err != nil {
				e.logger.Errorf("could not unwrap trigger event", resp.Err)
				continue
			}

			executionID, err := generateExecutionID(e.workflow.id, te.ID)
			if err != nil {
				e.logger.Errorf("could not generate execution ID", resp.Err)
				continue
			}

			err = e.startExecution(ctx, executionID, resp.Value)
			if err != nil {
				e.logger.Errorf("failed to start execution: %w", err)
			}
		case pendingStepRequest := <-e.pendingStepRequests:
			// Wait for a new worker to be available before dispatching a new one.
			// We'll do this up to newWorkerTimeout. If this expires, we'll put the
			// message back on the queue and keep going.
			t := time.NewTimer(e.newWorkerTimeout)
			select {
			case <-e.newWorkerCh:
				e.wg.Add(1)
				go e.workerForStepRequest(ctx, pendingStepRequest)
			case <-t.C:
				e.logger.Errorf("timed out when spinning off worker for pending step request %+v", pendingStepRequest)
				e.pendingStepRequests <- pendingStepRequest
			}
			t.Stop()
		case stepUpdate := <-e.stepUpdateCh:
			// Executed synchronously to ensure we correctly schedule subsequent tasks.
			err := e.handleStepUpdate(ctx, stepUpdate)
			if err != nil {
				e.logger.Errorf("failed to update step state: %+v, %s", stepUpdate, err)
			}
		}
	}
}

func generateExecutionID(workflowID, eventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(eventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}

// startExecution kicks off a new workflow execution when a trigger event is received.
func (e *Engine) startExecution(ctx context.Context, executionID string, event values.Value) error {
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
		workflowID:  e.workflow.id,
		executionID: executionID,
		status:      statusStarted,
	}

	err := e.executionStates.add(ctx, ec)
	if err != nil {
		return err
	}

	// Find the tasks we need to fire when a trigger has fired and enqueue them.
	// This consists of a) nodes without a dependency and b) nodes which depend
	// on a trigger
	triggerDependents, err := e.workflow.dependents(keywordTrigger)
	if err != nil {
		return err
	}

	for _, td := range triggerDependents {
		e.queueIfReady(*ec, td)
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
			err := e.workflow.walkDo(keywordTrigger, func(s *step) error {
				step, ok := state.steps[s.Ref]
				// The step is missing from the state,
				// which means it hasn't been processed yet.
				// Let's mark `workflowCompleted` = false, and
				// continue.
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
				err := e.finishExecution(ctx, state.executionID, statusCompleted)
				if err != nil {
					return err
				}
			}
		}

		for _, sd := range stepDependents {
			e.queueIfReady(state, sd)
		}
	case statusErrored:
		err := e.finishExecution(ctx, state.executionID, statusErrored)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) queueIfReady(state executionState, step *step) {
	// Check if all dependencies are completed for the current step
	var waitingOnDependencies bool
	for _, dr := range step.dependencies {
		stepState, ok := state.steps[dr]
		if !ok {
			waitingOnDependencies = true
			continue
		}

		// Unless the dependency is complete,
		// we'll mark waitingOnDependencies = true.
		// This includes cases where one of the dependent
		// steps has errored, since that means we shouldn't
		// schedule the step for execution.
		if stepState.status != statusCompleted {
			waitingOnDependencies = true
		}
	}

	// If all dependencies are completed, enqueue the step.
	if !waitingOnDependencies {
		e.logger.Debugw("step request enqueued", "ref", step.Ref, "state", copyState(state))
		e.pendingStepRequests <- stepRequest{
			state:   copyState(state),
			stepRef: step.Ref,
		}
	}
}

func (e *Engine) finishExecution(ctx context.Context, executionID string, status string) error {
	e.logger.Infow("finishing execution", "executionID", executionID, "status", status)
	err := e.executionStates.updateStatus(ctx, executionID, status)
	if err != nil {
		return err
	}

	// Signal that an execution has finished in a
	// non-blocking fashion. This is intended for
	// testing purposes only.
	select {
	case e.xxxExecutionFinished <- executionID:
	default:
	}

	return nil
}

func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	defer func() { e.newWorkerCh <- struct{}{} }()
	defer e.wg.Done()

	e.logger.Debugw("executing on a step event", "stepRef", msg.stepRef, "executionID", msg.state.executionID)
	stepState := &stepState{
		outputs:     &stepOutput{},
		executionID: msg.state.executionID,
		ref:         msg.stepRef,
	}

	inputs, outputs, err := e.executeStep(ctx, msg)
	if err != nil {
		e.logger.Errorf("error executing step request: %s", err, "executionID", msg.state.executionID, "stepRef", msg.stepRef)
		stepState.outputs.err = err
		stepState.status = statusErrored
	} else {
		e.logger.Infow("step executed successfully", "executionID", msg.state.executionID, "stepRef", msg.stepRef, "outputs", outputs)
		stepState.outputs.value = outputs
		stepState.status = statusCompleted
	}

	stepState.inputs = inputs

	// Let's try and emit the stepUpdate.
	// If the context is canceled, we'll just drop the update.
	// This means the engine is shutting down and the
	// receiving loop may not pick up any messages we emit.
	// Note: When full persistence support is added, any hanging steps
	// like this one will get picked up again and will be reprocessed.
	select {
	case <-ctx.Done():
		e.logger.Errorf("context canceled before step update could be issued", err, "executionID", msg.state.executionID, "stepRef", msg.stepRef)
	case e.stepUpdateCh <- *stepState:
	}
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
		Config: step.config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          msg.state.workflowID,
			WorkflowExecutionID: msg.state.executionID,
		},
	}

	resp, err := capabilities.ExecuteSync(ctx, step.capability, tr)
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

func (e *Engine) deregisterTrigger(ctx context.Context, t *triggerCapability) error {
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
			WorkflowID: e.workflow.id,
		},
		Inputs: triggerInputs,
		Config: t.config,
	}

	// if t.trigger == nil, then we haven't initialized the workflow
	// yet, and can safely consider the trigger deregistered with
	// no further action.
	if t.trigger != nil {
		return t.trigger.UnregisterTrigger(ctx, deregRequest)
	}

	return nil
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

		err := e.workflow.walkDo(keywordTrigger, func(s *step) error {
			if s.Ref == keywordTrigger {
				return nil
			}

			reg := capabilities.UnregisterFromWorkflowRequest{
				Metadata: capabilities.RegistrationMetadata{
					WorkflowID: e.workflow.id,
				},
				Config: s.config,
			}

			// if capability is nil, then we haven't initialized
			// the workflow yet and can safely consider it deregistered
			// with no further action.
			if s.capability == nil {
				return nil
			}

			innerErr := s.capability.UnregisterFromWorkflow(ctx, reg)
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
	Spec             string
	WorkflowID       string
	Lggr             logger.Logger
	Registry         core.CapabilitiesRegistry
	MaxWorkerLimit   int
	QueueSize        int
	NewWorkerTimeout time.Duration
}

const (
	defaultWorkerLimit      = 100
	defaultQueueSize        = 100000
	defaultNewWorkerTimeout = 2 * time.Second
)

func NewEngine(cfg Config) (engine *Engine, err error) {
	if cfg.MaxWorkerLimit == 0 {
		cfg.MaxWorkerLimit = defaultWorkerLimit
	}

	if cfg.QueueSize == 0 {
		cfg.QueueSize = defaultQueueSize
	}

	if cfg.NewWorkerTimeout == 0 {
		cfg.NewWorkerTimeout = defaultNewWorkerTimeout
	}

	// TODO: validation of the workflow spec
	// We'll need to check, among other things:
	// - that there are no step `ref` called `trigger` as this is reserved for any triggers
	// - that there are no duplicate `ref`s
	// - that the `ref` for any triggers is empty -- and filled in with `trigger`
	// - that the resulting graph is strongly connected (i.e. no disjointed subgraphs exist)
	// - etc.

	workflow, err := Parse(cfg.Spec)
	if err != nil {
		return nil, err
	}

	workflow.id = cfg.WorkflowID

	// Instantiate semaphore to put a limit on the number of workers
	newWorkerCh := make(chan struct{}, cfg.MaxWorkerLimit)
	for i := 0; i < cfg.MaxWorkerLimit; i++ {
		newWorkerCh <- struct{}{}
	}

	engine = &Engine{
		logger:               cfg.Lggr.Named("WorkflowEngine"),
		registry:             cfg.Registry,
		workflow:             workflow,
		executionStates:      newInMemoryStore(),
		pendingStepRequests:  make(chan stepRequest, cfg.QueueSize),
		newWorkerCh:          newWorkerCh,
		stepUpdateCh:         make(chan stepState),
		triggerEvents:        make(chan capabilities.CapabilityResponse),
		stopCh:               make(chan struct{}),
		newWorkerTimeout:     cfg.NewWorkerTimeout,
		xxxExecutionFinished: make(chan string),
	}
	return engine, nil
}
