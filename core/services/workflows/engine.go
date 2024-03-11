package workflows

import (
	"context"
	"errors"
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
	result, err := e.handleConsensus(ctx, event)
	if err != nil {
		e.logger.Errorf("error in handleConsensus %v", err)
		return
	}
	err = e.handleTargets(ctx, result)
	if err != nil {
		e.logger.Error("error in handleTargets %v", err)
	}
}

func (e *Engine) handleConsensus(ctx context.Context, event capabilities.CapabilityResponse) (values.Value, error) {
	e.logger.Debugw("running consensus", "event", event)
	consensus := e.workflow.Consensus[0]
	cm, err := values.NewMap(consensus.Config)
	if err != nil {
		return nil, err
	}
	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          mockedWorkflowID,
			WorkflowExecutionID: mockedExecutionID,
		},
		Inputs: &values.Map{
			Underlying: map[string]values.Value{
				// each node provides a single observation - outputs of mercury trigger
				"observations": &values.List{
					Underlying: []values.Value{event.Value},
				},
			},
		},
		Config: cm,
	}
	chReports := make(chan capabilities.CapabilityResponse, 10)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = e.consensus.Execute(newCtx, chReports, cr)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-chReports:
		if resp.Err != nil {
			return nil, resp.Err
		}
		return resp.Value, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (e *Engine) handleTargets(ctx context.Context, resp values.Value) error {
	e.logger.Debugw("handle targets")
	inputs := map[string]values.Value{
		"report": resp,
	}

	var combinedErr error
	for i, targetCapability := range e.targets {
		target := e.workflow.Targets[i]
		e.logger.Debugw("sending to target", "target", e.workflow.Targets[i], "inputs", inputs)
		cm, err := values.NewMap(target.Config)
		if err != nil {
			combinedErr = errors.Join(combinedErr, err)
			continue
		}

		tr := capabilities.CapabilityRequest{
			Inputs: &values.Map{Underlying: inputs},
			Config: cm,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          mockedWorkflowID,
				WorkflowExecutionID: mockedExecutionID,
			},
		}
		_, err = capabilities.ExecuteSync(ctx, targetCapability, tr)
		combinedErr = errors.Join(combinedErr, err)
	}
	return combinedErr
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
        - report_data.outputs
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
        - evm_median.outputs.reports
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: [($inputs.report)]
      abi: "receive(report bytes)"
  - type: write_ethereum-testnet-sepolia
    inputs:
      report:
        - evm_median.outputs.reports
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
