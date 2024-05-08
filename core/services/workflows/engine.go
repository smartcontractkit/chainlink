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
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	// NOTE: max 32 bytes per ID - consider enforcing exactly 32 bytes?
	mockedTriggerID  = "cccccccccc0000000000000000000000"
	mockedWorkflowID = "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
)

type donInfo struct {
	*capabilities.DON
	PeerID func() *p2ptypes.PeerID
}

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	logger              logger.Logger
	registry            core.CapabilitiesRegistry
	workflow            *workflow
	donInfo             donInfo
	executionStates     *inMemoryStore
	pendingStepRequests chan stepRequest
	triggerEvents       chan capabilities.CapabilityResponse
	newWorkerCh         chan struct{}
	stepUpdateCh        chan stepState
	wg                  sync.WaitGroup
	stopCh              services.StopChan
	newWorkerTimeout    time.Duration

	// testing lifecycle hook to signal when an execution is finished.
	onExecutionFinished func(string)
	// testing lifecycle hook to signal initialization status
	afterInit func(success bool)
	// Used for testing to control the number of retries
	// we'll do when initializing the engine.
	maxRetries int
	// Used for testing to control the retry interval
	// when initializing the engine.
	retryMs int
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

// resolveWorkflowCapabilities does the following:
//
// 1. Resolves the underlying capability for each trigger
// 2. Registers each step's capability to this workflow
func (e *Engine) resolveWorkflowCapabilities(ctx context.Context) error {
	//
	// Step 1. Resolve the underlying capability for each trigger
	//
	triggersInitialized := true
	for _, t := range e.workflow.triggers {
		tg, err := e.registry.GetTrigger(ctx, t.Type)
		if err != nil {
			e.logger.Errorf("failed to get trigger capability: %s", err)
			// we don't immediately return here, since we want to retry all triggers
			// to notify the user of all errors at once.
			triggersInitialized = false
		} else {
			t.trigger = tg
		}
	}
	if !triggersInitialized {
		return fmt.Errorf("failed to resolve triggers")
	}

	// Step 2. Walk the graph and register each step's capability to this workflow
	//
	// This means:
	// - fetching the capability
	// - register the capability to this workflow
	// - initializing the step's executionStrategy
	capabilityRegistrationErr := e.workflow.walkDo(keywordTrigger, func(s *step) error {
		// The graph contains a dummy step for triggers, but
		// we handle triggers separately since there might be more than one
		// trigger registered to a workflow.
		if s.Ref == keywordTrigger {
			return nil
		}

		err := e.initializeCapability(ctx, s)
		if err != nil {
			return err
		}

		return e.initializeExecutionStrategy(s)
	})

	return capabilityRegistrationErr
}

func (e *Engine) initializeCapability(ctx context.Context, s *step) error {
	// If the capability already exists, that means we've already registered it
	if s.capability != nil {
		return nil
	}

	cp, err := e.registry.Get(ctx, s.Type)
	if err != nil {
		return fmt.Errorf("failed to get capability with ref %s: %s", s.Type, err)
	}

	// We configure actions, consensus and targets here, and
	// they all satisfy the `CallbackCapability` interface
	cc, ok := cp.(capabilities.CallbackCapability)
	if !ok {
		return fmt.Errorf("could not coerce capability %s to CallbackCapability", s.Type)
	}

	if s.config == nil {
		configMap, newMapErr := values.NewMap(s.Config)
		if newMapErr != nil {
			return fmt.Errorf("failed to convert config to values.Map: %s", newMapErr)
		}
		s.config = configMap
	}

	registrationRequest := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: e.workflow.id,
		},
		Config: s.config,
	}

	err = cc.RegisterToWorkflow(ctx, registrationRequest)
	if err != nil {
		return fmt.Errorf("failed to register to workflow (%+v): %w", registrationRequest, err)
	}

	s.capability = cc
	return nil
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

	retryErr := retryable(ctx, e.logger, e.retryMs, e.maxRetries, func() error {
		err := e.resolveWorkflowCapabilities(ctx)
		if err != nil {
			return fmt.Errorf("failed to resolve workflow: %s", err)
		}
		return nil
	})

	if retryErr != nil {
		e.logger.Errorf("initialization failed: %s", retryErr)
		e.afterInit(false)
		return
	}

	e.logger.Debug("capabilities resolved, registering triggers")
	for _, t := range e.workflow.triggers {
		err := e.registerTrigger(ctx, t)
		if err != nil {
			e.logger.Errorf("failed to register trigger: %s", err)
		}
	}

	e.logger.Info("engine initialized")
	e.afterInit(true)
}

// initializeExecutionStrategy for `step`.
// Broadly speaking, we'll use `immediateExecution` for non-target steps
// and `scheduledExecution` for targets. If we don't have the necessary
// config to initialize a scheduledExecution for a target, we'll fallback to
// using `immediateExecution`.
func (e *Engine) initializeExecutionStrategy(step *step) error {
	if step.executionStrategy != nil {
		return nil
	}

	// If donInfo has no peerID, then the peer wrapper hasn't been initialized.
	// Let's error and try again next time around.
	if e.donInfo.PeerID() == nil {
		return fmt.Errorf("failed to initialize execution strategy: peer ID %s has not been initialized", e.donInfo.PeerID())
	}

	ie := immediateExecution{}
	if step.CapabilityType != capabilities.CapabilityTypeTarget {
		e.logger.Debugf("initializing step %+v with immediate execution strategy: not a target", step)
		step.executionStrategy = ie
		return nil
	}

	dinfo := e.donInfo
	if dinfo.DON == nil {
		e.logger.Debugf("initializing target step with immediate execution strategy: donInfo %+v", e.donInfo)
		step.executionStrategy = ie
		return nil
	}

	var position *int
	for i, w := range dinfo.Members {
		if w == *dinfo.PeerID() {
			idx := i
			position = &idx
		}
	}

	if position == nil {
		e.logger.Debugf("initializing step %+v with immediate execution strategy: position not found in donInfo %+v", step, e.donInfo)
		step.executionStrategy = ie
		return nil
	}

	step.executionStrategy = scheduledExecution{
		DON:      e.donInfo.DON,
		Position: *position,
		PeerID:   e.donInfo.PeerID(),
	}
	e.logger.Debugf("initializing step %+v with scheduled execution strategy", step)
	return nil
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

	e.onExecutionFinished(executionID)
	return nil
}

func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	defer func() { e.newWorkerCh <- struct{}{} }()
	defer e.wg.Done()

	// Instantiate a child logger; in addition to the WorkflowID field the workflow
	// logger will already have, this adds the `stepRef` and `executionID`
	l := e.logger.With("stepRef", msg.stepRef, "executionID", msg.state.executionID)

	l.Debugw("executing on a step event")
	stepState := &stepState{
		outputs:     &stepOutput{},
		executionID: msg.state.executionID,
		ref:         msg.stepRef,
	}

	inputs, outputs, err := e.executeStep(ctx, l, msg)
	if err != nil {
		l.Errorf("error executing step request: %s", err)
		stepState.outputs.err = err
		stepState.status = statusErrored
	} else {
		l.Infow("step executed successfully", "outputs", outputs)
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
		l.Errorf("context canceled before step update could be issued", err)
	case e.stepUpdateCh <- *stepState:
	}
}

// executeStep executes the referenced capability within a step and returns the result.
func (e *Engine) executeStep(ctx context.Context, l logger.Logger, msg stepRequest) (*values.Map, values.Value, error) {
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

	output, err := step.executionStrategy.Apply(ctx, l, step.capability, tr)
	if err != nil {
		return inputs, nil, err
	}

	return inputs, output, err
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
		e.logger.Info("shutting down engine")
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
	DONInfo          *capabilities.DON
	PeerID           func() *p2ptypes.PeerID

	// For testing purposes only
	maxRetries          int
	retryMs             int
	afterInit           func(success bool)
	onExecutionFinished func(weid string)
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

	if cfg.retryMs == 0 {
		cfg.retryMs = 5000
	}

	if cfg.afterInit == nil {
		cfg.afterInit = func(success bool) {}
	}

	if cfg.onExecutionFinished == nil {
		cfg.onExecutionFinished = func(weid string) {}
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
		logger:   cfg.Lggr.Named("WorkflowEngine").With("workflowID", cfg.WorkflowID),
		registry: cfg.Registry,
		workflow: workflow,
		donInfo: donInfo{
			DON:    cfg.DONInfo,
			PeerID: cfg.PeerID,
		},
		executionStates:     newInMemoryStore(),
		pendingStepRequests: make(chan stepRequest, cfg.QueueSize),
		newWorkerCh:         newWorkerCh,
		stepUpdateCh:        make(chan stepState),
		triggerEvents:       make(chan capabilities.CapabilityResponse),
		stopCh:              make(chan struct{}),
		newWorkerTimeout:    cfg.NewWorkerTimeout,

		onExecutionFinished: cfg.onExecutionFinished,
		afterInit:           cfg.afterInit,
		maxRetries:          cfg.maxRetries,
		retryMs:             cfg.retryMs,
	}
	return engine, nil
}
