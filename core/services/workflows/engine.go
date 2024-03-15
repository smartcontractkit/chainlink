package workflows

import (
	"context"
	"fmt"
	"time"

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

type Engine struct {
	services.StateMachine
	logger     logger.Logger
	registry   types.CapabilitiesRegistry
	trigger    capabilities.TriggerCapability
	consensus  capabilities.ConsensusCapability
	targets    []capabilities.TargetCapability
	workflow   *Workflow
	callbackCh chan capabilities.CapabilityResponse
	cancel     func()
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, cancel := context.WithCancel(context.Background())
		e.cancel = cancel
		go e.init(ctx)
		go e.triggerHandlerLoop(ctx)
		return nil
	})
}

func (e *Engine) init(ctx context.Context) {
	retrySec := 5
	ticker := time.NewTicker(time.Duration(retrySec) * time.Second)
	defer ticker.Stop()

	// Note: in our hardcoded workflow, there is only one trigger,
	// and one consensus step.
	trigger := e.workflow.Triggers[0]
	consensus := e.workflow.Consensus[0]

	var err error
LOOP:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.trigger, err = e.registry.GetTrigger(ctx, trigger.Type)
			if err != nil {
				e.logger.Errorf("failed to get trigger capability: %s, retrying in %d seconds", err, retrySec)
				break
			}

			e.consensus, err = e.registry.GetConsensus(ctx, consensus.Type)
			if err != nil {
				e.logger.Errorf("failed to get consensus capability: %s, retrying in %d seconds", err, retrySec)
				break
			}
			failed := false
			e.targets = make([]capabilities.TargetCapability, len(e.workflow.Targets))
			for i, target := range e.workflow.Targets {
				e.targets[i], err = e.registry.GetTarget(ctx, target.Type)
				if err != nil {
					e.logger.Errorf("failed to get target capability: %s, retrying in %d seconds", err, retrySec)
					failed = true
					break
				}
			}
			if !failed {
				break LOOP
			}
		}
	}

	// we have all needed capabilities, now we can register for trigger events
	err = e.registerTrigger(ctx)
	if err != nil {
		e.logger.Errorf("failed to register trigger: %s", err)
	}

	// also register for consensus
	cm, err := values.NewMap(consensus.Config)
	if err != nil {
		e.logger.Errorf("failed to convert config to values.Map: %s", err)
	}
	reg := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Config: cm,
	}
	err = e.consensus.RegisterToWorkflow(ctx, reg)
	if err != nil {
		e.logger.Errorf("failed to register consensus: %s", err)
	}

	e.logger.Info("engine initialized")
}

func (e *Engine) registerTrigger(ctx context.Context) error {
	trigger := e.workflow.Triggers[0]
	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": mockedTriggerID,
		},
	)
	if err != nil {
		return err
	}

	tc, err := values.NewMap(trigger.Config)
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
	err = e.trigger.RegisterTrigger(ctx, e.callbackCh, triggerRegRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate mercury_trigger, %s", err)
	}
	return nil
}

func (e *Engine) triggerHandlerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-e.callbackCh:
			go e.handleExecution(ctx, resp)
		}
	}
}

func (e *Engine) handleExecution(ctx context.Context, event capabilities.CapabilityResponse) {
	e.logger.Debugw("executing on a trigger event", "event", event)
	trigger := e.workflow.Triggers[0]
	if event.Err != nil {
		e.logger.Errorf("trigger event was an error; not executing", event.Err)
		return
	}

	ec := &executionState{
		steps: map[string]*stepState{
			trigger.Ref: {
				outputs: &stepOutput{
					value: event.Value,
				},
			},
		},
		workflowID:  mockedWorkflowID,
		executionID: mockedExecutionID,
	}

	consensus := e.workflow.Consensus[0]
	err := e.handleStep(ctx, ec, consensus)
	if err != nil {
		e.logger.Errorf("error in handleConsensus %v", err)
		return
	}

	for _, trg := range e.workflow.Targets {
		err := e.handleStep(ctx, ec, trg)
		if err != nil {
			e.logger.Errorf("error in handleTargets %v", err)
			return
		}
	}
}

func (e *Engine) handleStep(ctx context.Context, es *executionState, node Capability) error {
	stepState := &stepState{
		outputs: &stepOutput{},
	}
	es.steps[node.Ref] = stepState

	// Let's get the capability. If we fail here, we'll bail out
	// and try to handle it at the next execution.
	cp, err := e.registry.Get(ctx, node.Type)
	if err != nil {
		return err
	}

	api, ok := cp.(capabilities.CallbackExecutable)
	if !ok {
		return fmt.Errorf("capability %s must be an action, consensus or target", node.Type)
	}

	i, err := findAndInterpolateAllKeys(node.Inputs, es)
	if err != nil {
		return err
	}

	inputs, err := values.NewMap(i.(map[string]any))
	if err != nil {
		return err
	}

	stepState.inputs = inputs

	config, err := values.NewMap(node.Config)
	if err != nil {
		return err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: inputs,
		Config: config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          es.workflowID,
			WorkflowExecutionID: es.executionID,
		},
	}

	resp, err := capabilities.ExecuteSync(ctx, api, tr)
	if err != nil {
		stepState.outputs.err = err
		return err
	}

	// `ExecuteSync` returns a `values.List` even if there was
	// just one return value. If that is the case, let's unwrap the
	// single value to make it easier to use in -- for example -- variable interpolation.
	if len(resp.Underlying) > 1 {
		stepState.outputs.value = resp
	} else {
		stepState.outputs.value = resp.Underlying[0]
	}
	return nil
}

func (e *Engine) Close() error {
	return e.StopOnce("Engine", func() error {
		defer e.cancel()

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
		return e.trigger.UnregisterTrigger(context.Background(), deregRequest)
	})
}

func NewEngine(lggr logger.Logger, registry types.CapabilitiesRegistry) (engine *Engine, err error) {
	yamlWorkflowSpec := `
triggers:
  - type: "on_mercury_report"
    ref: report_data
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD
        
consensus:
  - type: "offchain_reporting"
    ref: evm_median
    inputs:
      observations:
        - $(report_data.outputs)
    config:
      aggregation_method: data_feeds_2_0
      aggregation_config:
        0x1111111111111111111100000000000000000000000000000000000000000000:
          deviation: "0.001"
          heartbeat: "30m"
        0x2222222222222222222200000000000000000000000000000000000000000000:
          deviation: "0.001"
          heartbeat: "30m"
        0x3333333333333333333300000000000000000000000000000000000000000000:
          deviation: "0.001"
          heartbeat: "30m"
      encoder: EVM
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: write_polygon-testnet-mumbai
    inputs:
      report:
        - $(evm_median.outputs.reports)
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: [($inputs.report)]
      abi: "receive(report bytes)"
  - type: write_ethereum-testnet-sepolia
    inputs:
      report:
        - $(evm_median.outputs.reports)
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(inputs.report)"]
      abi: "receive(report bytes)"
`

	workflow, err := Parse(yamlWorkflowSpec)
	if err != nil {
		return nil, err
	}
	engine = &Engine{
		logger:     lggr.Named("WorkflowEngine"),
		registry:   registry,
		workflow:   workflow,
		callbackCh: make(chan capabilities.CapabilityResponse),
	}
	return engine, nil
}
